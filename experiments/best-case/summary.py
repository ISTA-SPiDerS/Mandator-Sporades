import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

setting = sys.argv[1]
numIter = sys.argv[2]

replicaBatchSize = 1000
replicaBatchTime = 5000

iterations = list(range(1, int(numIter) + 1))
pipelines = [10]

arrivals = [500, 10000, 25000, 35000, 40000, 50000, 60000, 75000, 100000]


def getEPaxosSummary():
    l_records = []
    for arrival in arrivals:
        for pipeline in pipelines:
            record = ["epaxos-exec", str(pipeline), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/best-case/logs/epaxos/" +str(arrival)+"/"+str(replicaBatchSize)+"/"+str(replicaBatchTime)+"/"+setting+"/"+str(pipeline)+ "/"+str(iteration)+"/execution/"
                t, l, n, e = getEPaxosPaxosPerformance(root, 7, 5)
                throughput.append(t)
                latency.append(l)
                nine9.append(n)
                err.append(e)
            record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
            record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
            l_records.append(record)

            record = ["epaxos-commit", str(pipeline), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/best-case/logs/epaxos/" +str(arrival)+"/"+str(replicaBatchSize)+"/"+str(replicaBatchTime)+"/"+setting+"/"+str(pipeline)+ "/"+str(iteration)+"/commit/"
                t, l, n, e = getEPaxosPaxosPerformance(root, 7, 5)
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
        for pipeline in pipelines:
            record = ["multi-paxos", str(pipeline), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/best-case/logs/paxos/" + str(arrival)+"/"+str(replicaBatchSize) + "/"+str(replicaBatchTime) +"/" +setting +"/" + str(pipeline) + "/" +str(iteration) +"/execution/"
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


        record = ["raft", str(1), str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/raft/" + str(arrival)+"/"+str(replicaBatchSize) + "/"+str(replicaBatchTime) +"/" +setting +"/" + str(pipeline) + "/" +str(iteration) +"/execution/"
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

def getMandatorSummary():
    l_records = []
    for arrival in arrivals:
        for algo in ["paxos", "async"]:
            record = ["mandator-"+algo, str(1), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/best-case/logs/mandator/"  + algo + "/"+str(arrival)+ "/" + str(replicaBatchSize)+"/"+str(replicaBatchTime) +"/"+setting + "/3/" + str(iteration) +"/execution/"
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

def getSporadesSummary():
    l_records = []
    for arrival in arrivals:
        for pipeline in pipelines:
            record = ["sporades", str(pipeline), str(arrival * 5)]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments/best-case/logs/sporades/" + str(arrival)+"/"+str(replicaBatchSize) + "/" + str(replicaBatchTime) + "/" + setting +"/3/" + str(pipeline) + "/" + str(iteration)+ "/execution/"
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
        record = ["rabia", str(1), str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/rabia/" + str(arrival) +"/" + str(replicaBatchSize)+ "/"+ setting + "/"+str(iteration)+ "/execution/"
            t, l, n, e = getRabiaPerformance(root, 5, 5)
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


headers = ["algo", "pipeline", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

ePaxosSummary = getEPaxosSummary()
paxos_raftSummary = getPaxosRaftSummary()
mandatorSummary = getMandatorSummary()
sporadesSummary = getSporadesSummary()
rabiaSummary = getRabiaSummary()


records = records + ePaxosSummary + paxos_raftSummary + mandatorSummary + sporadesSummary + rabiaSummary

import csv

with open("experiments/best-case/logs/summary.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)