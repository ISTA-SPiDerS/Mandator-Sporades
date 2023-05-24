import os
import sys

# arrivalRate=$1
# algo=$2
# asyncBatchTime=$3
# attackTime=$4
# viewTimeout=$5

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 paxos 10 0 50000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 paxos 10 0 20000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 async 10 0 50000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 async 10 0 20000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 async 10 0 10000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 paxos 10 2 30000")

os.system("/bin/bash /home/pasindu/Documents/async-consensus/experiments/bash/local/local-consensus-test.sh 20000 async 10 2 30000")

sys.stdout.flush()

print("Test completed")