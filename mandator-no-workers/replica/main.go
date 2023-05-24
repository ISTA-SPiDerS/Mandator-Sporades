package main

import (
	"async-consensus/configuration"
	"async-consensus/replica/src"
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.Int64("name", 1, "name of the replica as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "configuration file")
	consAlgo := flag.String("consAlgo", "async", "consensus algo [async, paxos]")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int("batchSize", 50, "batch size")
	batchTime := flag.Int("batchTime", 5000, "maximum time to wait for collecting a batch of requests in micro seconds")
	debugOn := flag.Bool("debugOn", false, "false or true")
	mode := flag.Int("mode", 2, "1 for all to all broadcast, and 2 for selective broadcast")
	debugLevel := flag.Int("debugLevel", 0, "debug level")
	viewTimeout := flag.Int("viewTimeout", 2000000, "view timeout in micro seconds")
	window := flag.Int("window", 10, "window for abortable broadcast")
	asyncBatchTime := flag.Int("asyncBatchTime", 10, "async batch time in ms")
	keyLen := flag.Int("keyLen", 8, "key length")
	valLen := flag.Int("valLen", 8, "value length")
	benchmarkMode := flag.Int("benchmarkMode", 0, "0: resident store, 1: redis")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		panic(err)
	}

	rp := src.New(int32(*name), cfg, *logFilePath, *batchSize, *batchTime, *debugOn, *mode, *debugLevel, *viewTimeout, *window, *asyncBatchTime, *consAlgo, *benchmarkMode, *keyLen, *valLen)

	rp.WaitForConnections()
	rp.StartOutgoingLinks()
	rp.Run() // this is run in main thread

}
