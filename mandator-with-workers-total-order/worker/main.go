package main

import (
	"async-consensus/configuration"
	"async-consensus/worker/src"
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.Int64("name", 5, "name of the worker as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "configuration file")
	workerMapConfigFile := flag.String("workerMapconfig", "configuration/local/workermapconfiguration.yml", "worker to replica matching configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int("batchSize", 50, "worker batch size")
	batchTime := flag.Int("batchTime", 5000, "maximum time to wait for collecting a batch of requests in micro seconds")
	debugOn := flag.Bool("debugOn", false, "false or true")
	mode := flag.Int("mode", 2, "1 for all to all broadcast. 2 for selective broadcast")
	debugLevel := flag.Int("debugLevel", 0, "debug level int")
	window := flag.Int("window", 100, "window for abortable broadcast")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		panic(err)
	}

	mapCfg, err := configuration.NewReplicaConfig(*workerMapConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load map config: %v\n", err)
		panic(err)
	}

	wr := src.New(int32(*name), cfg, mapCfg, *logFilePath, *batchSize, *batchTime, *debugOn, *mode, *debugLevel, *window)

	wr.WaitForConnections()
	wr.StartOutgoingLinks()
	wr.Run() // this is run in main thread

}
