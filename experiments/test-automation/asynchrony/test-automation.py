import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir)
from execute import *
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *


os.system("/bin/bash experiments/setup-5/setup.sh")

scenario="asynchrony"
replicaBatchSize=str(3000)
replicaBatchTime=str(5000)
clientBatchSize=str(50)
clientBatchTime=str(1000)
setting="WAN"
pipelineLength=str(10)
asyncTimeout=str(0)
benchmarkMode=str(0)
asyncTimeEpochSize=str(500)
viewTimeout=str(300000000)
clientWindow=str(1000)
collectClientLogs="no"
isLeaderKill="no"



def simulatePaxosRaft():
    MIN_ARRIVAL = 5
    MAX_ARRIVAL = 80000
    INIT_GUAGE = MAX_ARRIVAL - MIN_ARRIVAL
    throughputs = []
    for iteration in [1,2,3]:
        start = MIN_ARRIVAL
        gauge = INIT_GUAGE
        last_throughput = -1
        found = False
        iter_num = 0
        while gauge > 10:
            iter_num = iter_num +1
            arrival = start
            # run paxos/raft for this configuration
            params = {}
            params["scenario"]=scenario
            params["arrival"]=str(arrival)
            params["replicaBatchSize"]=replicaBatchSize
            params["replicaBatchTime"] = replicaBatchTime
            params["clientBatchSize"]=clientBatchSize
            params["clientBatchTime"]=clientBatchTime
            params["setting"]=setting
            params["pipelineLength"]=pipelineLength
            params["algo"]="paxos"
            params["asyncTimeout"]=asyncTimeout
            params["benchmarkMode"]=benchmarkMode
            params["asyncTimeEpochSize"]=asyncTimeEpochSize
            params["viewTimeout"]=viewTimeout
            params["clientWindow"]=clientWindow
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(iteration)
            runPaxosRaft(params)
            throughput = getPaxosRaftPerformance("experiments/"+scenario+"/logs/paxos_raft/" + str(arrival) + "/"+ replicaBatchSize + "/"+ replicaBatchTime + "/"+ clientBatchSize + "/"+ clientBatchTime + "/"+ setting + "/"+ pipelineLength + "/"+ "paxos" + "/"+ asyncTimeout + "/"+ benchmarkMode + "/"+ asyncTimeEpochSize + "/"+ viewTimeout + "/"+ clientWindow + "/"+ str(iteration) + "/execution/", 21, 5)
            print("Multi-Paxos iteration : " + str(iter_num)+", throughput: "+str(throughput))
            sys.stdout.flush()
            if throughput >= 0.9 * (arrival *5):
                if arrival >= MAX_ARRIVAL:
                    throughputs.append(throughput)
                    found = True
                    break
                else:
                    last_throughput = throughput
                    start = start + gauge
                    guage = guage / 2
                    continue
            if throughput < 0.9 * (arrival *5):
                if arrival < MIN_ARRIVAL:
                    throughputs.append(throughput)
                    found = True
                    break
                else:
                    last_throughput = throughput
                    start = start - guage
                    guage = guage / 2
                    continue

        if not found:
                throughputs.append(last_throughput)

    print("Multi-Paxos: throughputs-:"+str(throughputs)+" , average throughput = "+str(sum(throughputs)/3))
    sys.stdout.flush()


simulatePaxosRaft()