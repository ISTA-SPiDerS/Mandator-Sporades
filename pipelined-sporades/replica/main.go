package main

import (
	"flag"
	"fmt"
	"os"
	"pipelined-sporades/configuration"
	"pipelined-sporades/replica/src"
)

func main() {
	name := flag.Int64("name", 1, "name of the replica as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int("batchSize", 50, "batch size")
	batchTime := flag.Int("batchTime", 5000, "maximum time to wait for collecting a batch of requests in micro seconds")
	debugOn := flag.Bool("debugOn", false, "false or true")
	isAsyncSim := flag.Bool("isAsyncSim", false, "to turn on asynchronous simulations")
	debugLevel := flag.Int("debugLevel", 0, "debug level")
	viewTimeout := flag.Int("viewTimeout", 2000000, "view timeout in micro seconds")
	keyLen := flag.Int("keyLen", 8, "key length")
	valLen := flag.Int("valLen", 8, "value length")
	benchmarkMode := flag.Int("benchmarkMode", 0, "0: resident store, 1: redis")
	pipelineLength := flag.Int("pipelineLength", 1, "pipeline length")
	networkbatchTime := flag.Int("networkbatchTime", 10, "artificial delay for sporades messages in ms")
	asyncSimTimeout := flag.Int("asyncSimTimeout", 10, "artificial delay in ms to simulate asynchrony")
	timeEpochSize := flag.Int("timeEpochSize", 500, "in ms the length of a time epoch for attacks")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		panic(err)
	}

	rp := src.New(int32(*name), cfg, *logFilePath, *batchSize, *batchTime, *debugOn, *debugLevel, *viewTimeout, *benchmarkMode, *keyLen, *valLen, *pipelineLength, *networkbatchTime, *isAsyncSim, *asyncSimTimeout, *timeEpochSize)

	rp.WaitForConnections()
	rp.StartOutgoingLinks()
	rp.Run() // this is run in main thread

}
