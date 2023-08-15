/bin/bash experiments/bash/epaxos.sh best-case 20000 3000 5000 WAN 10 2 10000 50 yes no 1

/bin/bash experiments/bash/mandator.sh best-case 20000 3000 5000 WAN paxos 10 10000 0 50 1000 0 1 200 3000000000 yes no 1

/bin/bash experiments/bash/mandator.sh best-case 20000 3000 5000 WAN async 10 10000 0 50 1000 0 1 200 3000000000 yes no 1

/bin/bash experiments/bash/paxos_raft.sh best_case 20000 3000 5000 50 1000 WAN 10 paxos 0 0 500 300000000 10000 yes no 1

/bin/bash experiments/bash/paxos_raft.sh best_case 20000 3000 5000 50 1000 WAN 10 raft 0 0 500 300000000 10000 yes no 1

/bin/bash experiments/bash/rabia.sh best-case  20000 3000 5000 WAN 50 no yes 1

/bin/bash experiments/bash/sporades.sh best-case 20000 3000 5000 50 1000 10000 0 500 0 300000000 WAN 10 10 yes no 1