import os

arrivalRate = [10, 1000, 10000]
viewTimeoutTime = [1000, 10000, 15000, 300000000]
batchTime = [1000, 5000]
batchSize = [1, 50]
networkBatchTime = [3, 5, 10]
asyncSimTime = [1, 5, 10]
consAlgo = ["paxos", "async"]
mode = [1, 2]

for arr in arrivalRate:
    for viewTimeout in viewTimeoutTime:
        for batchTi in batchTime:
            for batchSi in batchSize:
                for networkBatchTim in networkBatchTime:
                    for asyncSimTim in asyncSimTime:
                        for cons in consAlgo:
                            for m in mode:
                                os.system("/bin/bash integration-test/safety_test.sh " + str(arr) + " " + str(
                                    viewTimeout) + " " + str(batchTi) + " " + str(batchSi) + " " + str(
                                    networkBatchTim) + " " + str(asyncSimTim) + " " + str(cons) + " " + str(m))
