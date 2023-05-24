Workers create mini mempool of client request batches, send to majority of other workers, and handle client responses

To start a worker with minimal options

```./worker/bin/child --name 21```

Supported options

```--name```: name of the worker as specified in the ```configuration.yml``` 

```--config```: configuration file 

```--workerMapconfig```: worker to replica matching configuration file 

```--logFilePath```: log file path 

```--batchSize```: worker batch size 

```--batchTime```: maximum time to wait for collecting a batch of requests in micro seconds 

```--debugOn```: ```false``` or ```true``` 

```--mode```: ```1``` for all to all broadcast. ```2``` for selective broadcast 

```--debugLevel```: debug level int 

```--window```: window for abortable broadcast