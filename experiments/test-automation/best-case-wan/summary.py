import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *

arrivals = [500, 10000, 30000, 35000, 40000, 45000, 50000, 60000, 80000, 100000]

scenario="best-case-wan"
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
clientWindow=str(10000)
collectClientLogs="no"
isLeaderKill="no"

iterations = [1,2,3]
headers = ["algo", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

def getEPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["epaxos-exec", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/" + scenario+"/logs/epaxos/" +str(arrival)+"/"+replicaBatchSize+"/"+replicaBatchTime+"/"+setting+"/"+pipelineLength+"/"+str(2)+"/"+str(1000)+ "/" +clientBatchSize+"/" +str(iteration)+"/execution/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 21, 5)
            throughput.append(t)
            latency.append(l)
            nine9.append(n)
            err.append(e)
        record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
        l_records.append(record)

        record = ["epaxos-commit", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/"+scenario+"/logs/epaxos/" +str(arrival)+"/"+replicaBatchSize+"/"+replicaBatchTime+"/"+setting+"/"+pipelineLength+"/"+str(2)+"/"+str(1000)+ "/" +clientBatchSize+"/" +str(iteration)+"/commit/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 21, 5)
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


def getPaxosRaftSummary():
    l_records = []
    for arrival in arrivals:
        for algo in ["paxos"]:
            record = [algo, str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments"+"/"+scenario+"/logs/paxos_raft/"+str(arrival)+"/"+ replicaBatchSize+"/"+ replicaBatchTime +"/" +clientBatchSize+ "/"+ clientBatchTime+"/"+setting +"/" +pipelineLength +"/" + algo +"/" +asyncTimeout+"/"+ benchmarkMode+ "/"+ asyncTimeEpochSize +"/"+ viewTimeout +"/"+ clientWindow +"/" + str(iteration) +"/"+"execution/"
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
                root = "experiments/" +scenario+"/logs/sporades/"+str(arrival)+ "/"+ replicaBatchSize+ "/" + replicaBatchTime +"/"+ clientBatchSize +"/"+ clientBatchTime +"/"+ clientWindow +"/"+asyncTimeout+"/"+asyncTimeEpochSize+"/"+benchmarkMode+"/"+viewTimeout+"/"+setting+"/"+str(networkNatchTime)+"/"+pipelineLength+"/"+str(iteration)+"/"+"execution/"
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
                root = "experiments/" +scenario+"/logs/mandator/"+str(arrival)+ "/"+ replicaBatchSize+ "/" + replicaBatchTime +"/"+setting+"/async/"+ str(networkNatchTime) +"/"+ clientWindow +"/"+ asyncTimeout +"/"+clientBatchSize+"/"+clientBatchTime+"/"+benchmarkMode+"/"+str(1)+"/"+asyncTimeEpochSize+"/"+str(viewTimeout)+"/"+str(iteration)+"/execution/"
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

            record = ["mandator-paxos-"+str(networkNatchTime), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/" +scenario+"/logs/mandator/"+str(arrival)+ "/"+ replicaBatchSize+ "/" + replicaBatchTime +"/"+setting+"/paxos/"+ str(networkNatchTime) +"/"+ clientWindow +"/"+ asyncTimeout +"/"+clientBatchSize+"/"+clientBatchTime+"/"+benchmarkMode+"/"+str(1)+"/"+asyncTimeEpochSize+"/"+str(viewTimeout)+"/"+str(iteration)+"/execution/"
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


ePaxosSummary = getEPaxosSummary()
paxos_raftSummary = getPaxosRaftSummary()
sporadesSummary = getSporadesSummary()
mandatorSummary = getMandatorSummary()


records = records + ePaxosSummary + paxos_raftSummary + sporadesSummary + mandatorSummary

import csv

with open("experiments/"+scenario+"/logs/"+scenario+"-summary"+str(datetime.now())+".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)