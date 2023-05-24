Client implementation supports two operations.

(1) Send a ```status``` request to replicas

(2) Send SMR ```request```s to replicas

To send a status request ```./client/bin/client --name 11 --requestType status --operationType [1, 2, 3]```

```OperationType 1``` for server bootstrapping, ```OperationType 2``` for server log printing, and ```OperationType 3``` for starting consensus layer

To send client requests with minimal options ```./client/bin/client --name 11 --defaultReplica 1 --requestType request```

Supported options

```--name```: name of the client as specified in the ```configuration.yml```

```--config```: configuration file 

```--logFilePath```: log file path 

```--batchSize```: client batch size 

```--batchTime```: maximum time to wait for collecting a batch of requests in micro seconds 

```--defaultReplica```: default replica to send requests to 

```--requestSize```: request size in bytes 

```--testDuration```: test duration in seconds

```--arrivalRate```: poisson arrival rate in requests per second

```--requestType```: "request type: [```status``` , ```request```] 

```--operationType```: Type of operation for a status request: ```1``` (bootstrap server), ```2```: (print log) 

```--debugOn```: ```false``` or ```true```

```--debugLevel```: debug level 

```--keyLen```: key length 

```--valLen```: value length 