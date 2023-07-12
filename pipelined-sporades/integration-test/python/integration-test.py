import os

# case 1: no batching, no pipelining, sync

arrivalRate = 10
viewTimeoutTime = 30000000
batchTime = 1000
batchSize = 1
pipelineLength = 1
networkbatchTime = 1
asyncSimTimeout = 0

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))

# case 2: batching no pipelining sync

arrivalRate = 10000
viewTimeoutTime = 30000000
batchTime = 5000
batchSize = 50
pipelineLength = 1
networkbatchTime = 3
asyncSimTimeout = 0

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))

# case 3: no batching, pipelining, sync

arrivalRate = 30
viewTimeoutTime = 30000000
batchTime = 1000
batchSize = 1
pipelineLength = 10
networkbatchTime = 1
asyncSimTimeout = 0

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))

# case 4: batching, pipelining, sync

arrivalRate = 10000
viewTimeoutTime = 30000000
batchTime = 5000
batchSize = 50
pipelineLength = 5
networkbatchTime = 3
asyncSimTimeout = 0

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))

# case 5: batching, pipelining, async-medium

arrivalRate = 1000
viewTimeoutTime = 15000
batchTime = 5000
batchSize = 50
pipelineLength = 5
networkbatchTime = 3
asyncSimTimeout = 6

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))

# case 6: batching, pipelining, async

arrivalRate = 100
viewTimeoutTime = 10000
batchTime = 5000
batchSize = 50
pipelineLength = 5
networkbatchTime = 3
asyncSimTimeout = 6

os.system("/bin/bash integration-test/safety_test.sh " + str(arrivalRate) + " " + str(viewTimeoutTime) + " " + str(
    batchTime) + " " + str(batchSize) + " " + str(pipelineLength) + " " + str(networkbatchTime) + " " + str(
    asyncSimTimeout))
