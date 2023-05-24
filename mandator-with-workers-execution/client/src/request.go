package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

/*
	Upon receiving a client response, add the request to the received requests map
*/

func (cl *Client) handleClientResponseBatch(batch *proto.ClientBatch) {
	cl.receivedResponses[batch.UniqueId] = requestBatch{
		batch: *batch,
		time:  time.Now(), // record the time when the response was received
	}
	workerIndex := cl.workerArrayIndex[batch.Sender]
	cl.lastSeenTimeMutexes[workerIndex].Lock()
	cl.lastSeenTimes[workerIndex] = time.Now() // mark the last time a response was received
	cl.lastSeenTimeMutexes[workerIndex].Unlock()
	common.Debug("Added response Batch from "+strconv.Itoa(int(batch.Sender))+" to received array", 0, cl.debugLevel, cl.debugOn)
	common.Debug("Response Batch contains "+fmt.Sprintf("%v", batch.Requests), 0, cl.debugLevel, cl.debugOn)
}

/*
	start the poisson arrival process (put arrivals to arrivalTimeChan) in a separate thread
	start request generation processes  (get arrivals from arrivalTimeChan and generate batches and send them) in separate threads, and send them to the default worker, and write batch to the correct array in sentRequests
	start failure detector that checks the time since the last response was received, and update the default worker -- currently disabled
	start the scheduler that schedules new requests
	the thread sleeps for test duration and then starts processing the responses. This is to handle inflight responses after the test duration
*/

func (cl *Client) SendRequests() {
	cl.generateArrivalTimes()
	cl.startRequestGenerators()
	//cl.startFailureDetector()
	cl.startScheduler() // this is sync, main thread waits for this to finish

	// end of test

	time.Sleep(time.Duration(cl.testDuration*5) * time.Second) // additional sleep duration to make sure that all the in-flight responses are received
	fmt.Printf("Finish sending requests \n")
	cl.computeStats()
}

/*
	Each request generator generates requests by generating string requests, forming batches, send batches and save them in the correct sent array
*/

func (cl *Client) startRequestGenerators() {
	for i := 0; i < numRequestGenerationThreads; i++ { // i is the thread number
		go func(threadNumber int) {
			localCounter := 0
			lastSent := time.Now() // used to get how long to wait
			for true {             // this runs forever
				numRequests := 0
				var requests []*proto.ClientBatch_SingleOperation
				// this loop collects requests until the minimum batch size is met OR the batch time is timeout
				for !(numRequests >= cl.clientBatchSize || (time.Now().Sub(lastSent).Microseconds() > int64(cl.clientBatchTime) && numRequests > 0)) {
					_ = <-cl.arrivalChan // keep collecting new requests arrivals
					requests = append(requests, &proto.ClientBatch_SingleOperation{
						Id: strconv.Itoa(int(cl.clientName)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter) + "." + strconv.Itoa(numRequests),
						Command: fmt.Sprintf("%d%v%v", rand.Intn(2),
							cl.RandString(cl.keyLen),
							cl.RandString(cl.valueLen)),
					})
					numRequests++
				}
				defaultWorkerIndex := threadNumber % len(cl.defaultWorkers)
				cl.defaultWorkerMutexes[defaultWorkerIndex].RLock()
				defaultWorker := cl.defaultWorkers[defaultWorkerIndex]
				cl.defaultWorkerMutexes[defaultWorkerIndex].RUnlock()
				// create a new client batch
				batch := proto.ClientBatch{
					Sender:   cl.clientName,
					Receiver: int32(defaultWorker),
					UniqueId: strconv.Itoa(int(cl.clientName)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter), // this is a unique string id,
					Type:     1,                                                                                                      // client request batch
					Note:     "",
					Requests: requests,
				}
				common.Debug("Sent "+strconv.Itoa(int(cl.clientName))+"."+strconv.Itoa(threadNumber)+"."+strconv.Itoa(localCounter)+" batch size "+strconv.Itoa(len(requests)), 0, cl.debugLevel, cl.debugOn)
				localCounter++
				rpcPair := common.RPCPair{
					Code: cl.messageCodes.ClientBatchRpc,
					Obj:  &batch,
				}

				cl.sendMessage(defaultWorker, rpcPair)
				lastSent = time.Now()
				cl.sentRequests[threadNumber] = append(cl.sentRequests[threadNumber], requestBatch{
					batch: batch,
					time:  time.Now(),
				})
			}

		}(i)
	}

}

/*
	After the request arrival time is arrived, inform the request generators
*/

func (cl *Client) startScheduler() {
	cl.startTime = time.Now()
	for time.Now().Sub(cl.startTime).Nanoseconds() < int64(cl.testDuration*1000*1000*1000) {
		nextArrivalTime := <-cl.arrivalTimeChan

		for time.Now().Sub(cl.startTime).Nanoseconds() < nextArrivalTime {
			// busy waiting until the time to dispatch this request arrives
		}
		cl.arrivalChan <- true
	}
}

/*
	Generate Poisson arrival times
*/

func (cl *Client) generateArrivalTimes() {
	go func() {
		lambda := float64(cl.arrivalRate) / (1000.0 * 1000.0 * 1000.0) // requests per nano second
		arrivalTime := 0.0

		for true {
			// Get the next probability value from Uniform(0,1)
			p := rand.Float64()

			//Plug it into the inverse of the CDF of Exponential(_lamnbda)
			interArrivalTime := -1 * (math.Log(1.0-p) / lambda)

			// Add the inter-arrival time to the running sum
			arrivalTime = arrivalTime + interArrivalTime

			cl.arrivalTimeChan <- int64(arrivalTime)
		}
	}()
}

/*
	Monitors the time the last response was received from each worker in the default workers array. If the default workers fail to send a response before a timeout, change the default worker
-- currently this is implemented, but not used because EPaxos and Rabia clients do not implement this
*/

func (cl *Client) startFailureDetector() {
	go func() {
		common.Debug("Starting failure detector", 0, cl.debugLevel, cl.debugOn)
		for true {
			time.Sleep(time.Duration(cl.workerTimeout) * time.Second)
			for i := 0; i < len(cl.defaultWorkers); i++ {
				workerIndex := cl.workerArrayIndex[cl.defaultWorkers[i]]
				cl.lastSeenTimeMutexes[workerIndex].Lock()
				if time.Now().Sub(cl.lastSeenTimes[workerIndex]).Seconds() > float64(cl.workerTimeout) {
					common.Debug("Default worker "+strconv.Itoa(int(cl.defaultWorkers[i]))+" has time out, setting a new replica at "+fmt.Sprintf("%v", time.Now().Sub(cl.startTime)), 4, cl.debugLevel, cl.debugOn)
					// change the default replica
					cl.defaultWorkerMutexes[i].Lock()
					cl.defaultWorkers[i] = cl.getRandomWorkerNode()
					cl.defaultWorkerMutexes[i].Unlock()
					common.Debug("New default worker "+fmt.Sprintf("%v", cl.defaultWorkers[i])+" at "+fmt.Sprintf("%v", time.Now().Sub(cl.startTime)), 4, cl.debugLevel, cl.debugOn)
					cl.lastSeenTimes[workerIndex] = time.Now()
				}
				cl.lastSeenTimeMutexes[workerIndex].Unlock()
			}
		}
	}()
}

/*
	Returns a random worker node
*/

func (cl *Client) getRandomWorkerNode() int32 {
	keys := make([]int32, len(cl.workerAddrList))
	i := 0
	for k := range cl.workerAddrList {
		keys[i] = k
		i++
	}
	return keys[rand.Intn(len(keys))]
}

/*
	random string generation adapted from the Rabia SOSP 2021 code base https://github.com/haochenpan/rabia/
*/

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // low conflict
	letterIdxBits = 6                                                      // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
)

/*
	generate a random string of length n
*/

func (cl *Client) RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
