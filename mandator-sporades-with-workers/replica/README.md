Replica implements the overlay for mem-blocks and consensus logic

To run a replica with minimum options

```./replica/bin/replica --name 0```

Supported parameters

```--name```: name of the replica as specified in the ```configuration.yml```

```--config```: configuration file

```--consAlgo```: consensus algo [```async```, ```paxos```]

```--workerMapconfig```: worker to replica matching configuration file 

```--logFilePath```: log file path 

```--batchSize```: worker batch size 

```--batchTime```: maximum time to wait for collecting a batch of requests in micro seconds 

```--debugOn```: ```false``` or ```true```

```--mode```: ```1``` for all to all broadcast, and ```2``` for selective broadcast 

```--debugLevel```: debug level int 

```--viewTimeout```: view timeout in micro seconds 

```--window```: window for abortable broadcast 

```--asyncBatchTime```: async batch time in ms 

```--keyLen```: key length 

```--valLen```: value length 

```--benchmarkMode``` ```0```- resident store, ```1```- redis 