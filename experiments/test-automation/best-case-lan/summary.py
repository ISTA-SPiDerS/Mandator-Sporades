import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *

arrivals = [500, 10000, 20000, 40000, 50000, 60000, 70000, 80000]

scenario="best-case-lan"
replicaBatchSize=str(2000)
replicaBatchTime=str(5000)
clientBatchSize=str(50)
clientBatchTime=str(1000)
setting="LAN"
pipelineLength=str(1)
asyncTimeout=str(0)
benchmarkMode=str(0)
asyncTimeEpochSize=str(500)
viewTimeout=str(300000000)
clientWindow=str(1000)
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
            root = "experiments/" + scenario+"/logs/epaxos/" +str(arrival)+"/"+replicaBatchSize+"/"+replicaBatchTime+"/"+setting+"/"+pipelineLength+"/"+str(2)+"/"+clientWindow+ "/" +clientBatchSize+"/" +str(iteration)+"/execution/"
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
            root = "experiments/"+scenario+"/logs/epaxos/" +str(arrival)+"/"+replicaBatchSize+"/"+replicaBatchTime+"/"+setting+"/"+pipelineLength+"/"+str(2)+"/"+clientWindow+ "/" +clientBatchSize+"/" +str(iteration)+"/commit/"
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
        for algo in ["paxos", "raft"]:
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


def getRabiaSummary():
    l_records = []
    for arrival in arrivals:
        record = ["rabia", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments" +"/"+ scenario+"/logs/rabia/" +str(arrival)+"/"+str(300)+"/"+str(5)+"/"+setting+"/"+clientBatchSize+"/"+str(iteration)+"/"+"execution/"
            t, l, n, e = getRabiaPerformance(root, 21, 5)
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
rabiaSummary = getRabiaSummary()


records = records + ePaxosSummary + paxos_raftSummary + sporadesSummary + rabiaSummary

import csv

with open("experiments/"+scenario+"/logs/"+scenario+"-summary-"+ str(datetime.now())+".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)