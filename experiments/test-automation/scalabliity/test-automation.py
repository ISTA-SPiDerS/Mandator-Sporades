import os
import sys

numReplicas = sys.argv[1]

os.system("/bin/bash experiments/setup-"+str(numReplicas)+ "/setup.sh")

arrivals = []

if numReplicas == 3:
    # upto 120k per client
    arrivals = [100, 500, 1000, 2000, 5000, 10000, 20000, 40000, 60000, 80000, 100000, 120000, 150000]

if numReplicas == 11:
    # upto 20k per client and then upto 35k for mandator
    arrivals = [100, 500, 1000, 2000, 5000, 8000, 10000, 12500, 15000, 17500, 20000, 22500, 25000, 27500, 30000, 32000, 35000]

for iteration in [1,2,3]:
    for arrival in arrivals:
        # paxos
        os.system("/bin/bash experiments/scalability/"+str(numReplicas)+"/paxos_raft.sh "+str(arrival))

        # mandator
        os.system("/bin/bash experiments/scalability/"+str(numReplicas)+"/mandator.sh "+str(arrival))

        # sporades
        os.system("/bin/bash experiments/scalability/"+str(numReplicas)+"/sporades.sh "+str(arrival))





