package main

import (
	"async-consensus/client/src"
	"async-consensus/configuration"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	name := flag.Int64("name", 5, "name of the client as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int("batchSize", 50, "client batch size")
	batchTime := flag.Int("batchTime", 5000, "maximum time to wait for collecting a batch of requests in micro seconds")
	defaultReplicas := flag.String("defaultReplicas", "11,12", "default workers to send requests to")
	workerTimeout := flag.Int("workerTimeout", 2, "worker timeout in seconds")
	requestSize := flag.Int("requestSize", 8, "request size in bytes")
	testDuration := flag.Int("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int("arrivalRate", 1000, "poisson arrival rate in requests per second")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int("operationType", 1, "Type of operation for a status request: 1 (bootstrap server), 2: (print log)")
	debugOn := flag.Bool("debugOn", false, "false or true")
	debugLevel := flag.Int("debugLevel", 1, "debug level int")
	keyLen := flag.Int("keyLen", 8, "key length")
	valLen := flag.Int("valLen", 8, "value length")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		panic(err)
	}

	defaultReplicaArray := getDefaultReplicas(*defaultReplicas)

	cl := src.New(int32(*name), cfg, *logFilePath, *batchSize, *batchTime, defaultReplicaArray, *workerTimeout, *requestSize, *testDuration, *arrivalRate, *requestType, *operationType, *debugOn, *debugLevel, *keyLen, *valLen)

	cl.WaitForConnections()
	cl.Run()
	cl.StartOutgoingLinks()
	cl.ConnectToWorkers()

	time.Sleep(10 * time.Second)

	if cl.RequestType == "status" {
		cl.SendStatus(cl.OperationType)
	} else if cl.RequestType == "request" {
		cl.SendRequests()
	}

}

/*
	Return an int32 array given a string by splitting and casting
*/

func getDefaultReplicas(replicas string) []int32 {
	stArray := strings.Split(replicas, ",")
	var intArray []int32
	for i := 0; i < len(stArray); i++ {
		j, err := strconv.ParseInt(stArray[i], 10, 32)
		if err != nil {
			panic(err)
		}
		intArray = append(intArray, int32(j))
	}

	return intArray
}
