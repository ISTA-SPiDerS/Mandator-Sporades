package main

import (
	"async-consensus/client/src"
	"async-consensus/configuration"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	name := flag.Int64("name", 11, "name of the client as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int("batchSize", 50, "client batch size")
	batchTime := flag.Int("batchTime", 5000, "maximum time to wait for collecting a batch of requests in micro seconds")
	defaultReplica := flag.Int64("defaultReplica", 1, "default replica to send requests to")
	requestSize := flag.Int("requestSize", 8, "request size in bytes")
	testDuration := flag.Int("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int("arrivalRate", 1000, "poisson arrival rate in requests per second")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int("operationType", 1, "Type of operation for a status request: 1 (bootstrap server), 2: (print log)")
	debugOn := flag.Bool("debugOn", false, "false or true")
	debugLevel := flag.Int("debugLevel", -1, "debug level int")
	keyLen := flag.Int("keyLen", 8, "key length")
	valLen := flag.Int("valLen", 8, "value length")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		panic(err)
	}

	cl := src.New(int32(*name), cfg, *logFilePath, *batchSize, *batchTime, int32(*defaultReplica), *requestSize, *testDuration, *arrivalRate, *requestType, *operationType, *debugOn, *debugLevel, *keyLen, *valLen)

	cl.WaitForConnections()
	cl.Run()
	cl.StartOutgoingLinks()
	cl.ConnectToReplicas()

	time.Sleep(10 * time.Second)

	if cl.RequestType == "status" {
		cl.SendStatus(cl.OperationType)
	} else if cl.RequestType == "request" {
		cl.SendRequests()
	}

}
