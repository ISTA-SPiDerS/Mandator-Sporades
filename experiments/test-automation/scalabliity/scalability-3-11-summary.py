import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *

numReplicas = int(sys.argv[1])

scenario="scalability-"+str(numReplicas)

arrivals = []

if numReplicas == 3:
    # upto 120k per client
    arrivals = [100, 500, 1000, 2000, 5000, 10000, 20000, 40000, 60000, 80000, 100000, 120000, 150000, 180000, 200000, 250000]

if numReplicas == 11:
    # upto 20k per client and then upto 35k for mandator
    arrivals = [100, 500, 1000, 2000, 5000, 8000, 10000, 12500, 15000, 17500, 20000, 22500, 25000, 27500, 30000, 32000, 35000]

iterations = [1,2,3]
headers = ["algo", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

def getEPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["epaxos-exec", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/" + scenario+"/logs/epaxos/" +str(arrival)+"/" +str(iteration)+"/execution/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 21, int(numReplicas))
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
            root = "experiments/"+scenario+"/logs/epaxos/" +str(arrival)+"/"+str(iteration)+"/commit/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 21, int(numReplicas))
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
                root = "experiments"+"/"+scenario+"/logs/paxos/"+str(arrival)+"/"+ str(iteration) +"/"
                t, l, n, e = getPaxosRaftPerformance(root, 21, int(numReplicas))
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
                root = "experiments/" +scenario+"/logs/sporades/"+str(arrival)+ "/" +str(iteration)+"/"
                t, l, n, e = getManatorSporadesPerformance(root, 21, int(numReplicas))
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
                root = "experiments/" +scenario+"/logs/mandator/"+str(arrival)+ "/"+str(iteration)+"/"
                t, l, n, e = getManatorSporadesPerformance(root, 21, int(numReplicas))
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


ePaxosSummary = []
if int(numReplicas) == 3:
    ePaxosSummary = getEPaxosSummary()

paxos_raftSummary = getPaxosRaftSummary()
sporadesSummary = getSporadesSummary()
mandatorSummary = getMandatorSummary()


records = records + ePaxosSummary + paxos_raftSummary + sporadesSummary + mandatorSummary

import csv

with open("experiments/"+scenario+"/logs/"+scenario+"-summary-"+str(numReplicas)+"-"+str(datetime.now())+".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)