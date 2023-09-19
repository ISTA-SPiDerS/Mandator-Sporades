import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *

arrivals = [100, 500, 1000, 5000, 7500, 10000, 12500, 15000, 20000, 25000, 30000, 40000, 50000]

scenario="asynchrony"
replicaBatchSize=str(3000)
replicaBatchTime=str(5000)
clientBatchSize=str(50)
clientBatchTime=str(1000)
setting="WAN"
pipelineLength=str(1)
asyncTimeout=str(500)
benchmarkMode=str(0)
asyncTimeEpochSize=str(sys.argv[1])
viewTimeout=str(300000)
clientWindow=str(10000)
collectClientLogs="no"
isLeaderKill="no"

iterations = [1,2,3]
headers = ["algo", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

def getPaxosRaftSummary():
    l_records = []
    for arrival in arrivals:
        for algo in ["paxos", "raft"]:
            record = [algo, str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments"+"/"+scenario+"/logs/paxos_raft/"+str(arrival)+"/"+ replicaBatchSize+"/"+ replicaBatchTime +"/" +clientBatchSize+ "/"+ clientBatchTime+"/"+setting +"/" +pipelineLength +"/" + algo +"/" +asyncTimeout+"/"+ benchmarkMode+ "/"+ asyncTimeEpochSize +"/"+ viewTimeout +"/"+ str(1000) +"/" + str(iteration) +"/"+"execution/"
                t, l, n, e = getPaxosRaftPerformance(root, 21, 5)
                throughput.append(t)
                latency.append(l)
                nine9.append(n)
                err.append(e)
            record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
            l_records.append(record)

    return l_records


def getSporadesSummary():
    l_records = []
    for arrival in arrivals:
        for networkNatchTime in [0]:
            record = ["sporades-"+str(networkNatchTime), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/" +scenario+"/logs/sporades/"+str(arrival)+ "/"+ replicaBatchSize+ "/" + replicaBatchTime +"/"+ clientBatchSize +"/"+ clientBatchTime +"/"+ str(1000) +"/"+asyncTimeout+"/"+asyncTimeEpochSize+"/"+benchmarkMode+"/"+viewTimeout+"/"+setting+"/"+str(networkNatchTime)+"/"+pipelineLength+"/"+str(iteration)+"/"+"execution/"
                t, l, n, e = getManatorSporadesPerformance(root, 21, 5)
                throughput.append(t)
                latency.append(l)
                nine9.append(n)
                err.append(e)
            record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
            l_records.append(record)

    return l_records


def getMandatorSummary():
    l_records = []
    for arrival in arrivals:
        for networkNatchTime in [30]:
            record = ["mandator-sporades-"+str(networkNatchTime), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/" +scenario+"/logs/mandator/"+str(arrival)+ "/"+ replicaBatchSize+ "/" + replicaBatchTime +"/"+setting+"/async/"+ str(networkNatchTime) +"/"+ str(10000) +"/"+ asyncTimeout +"/"+clientBatchSize+"/"+clientBatchTime+"/"+benchmarkMode+"/"+str(1)+"/"+asyncTimeEpochSize+"/"+str(viewTimeout)+"/"+str(iteration)+"/execution/"
                t, l, n, e = getManatorSporadesPerformance(root, 21, 5)
                throughput.append(t)
                latency.append(l)
                nine9.append(n)
                err.append(e)
            record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
            l_records.append(record)

    return l_records


paxos_raftSummary = getPaxosRaftSummary()
sporadesSummary = getSporadesSummary()
mandatorSummary = getMandatorSummary()


records = records + paxos_raftSummary + sporadesSummary + mandatorSummary

import csv

with open("experiments/"+scenario+"/logs/"+scenario+"-"+asyncTimeEpochSize+"-summary"+str(datetime.now())+".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)