import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

setting = sys.argv[1]  # LAN or WAN
numIter = sys.argv[2]

if setting != "LAN" and setting != "WAN":
    exit("wrong input, input should be LAN/WAN")

if int(numIter) < 3:
    exit("at least 3 iterations needed")

replicaBatchSize = 2000
replicaBatchTime = 4000

if setting == "WAN":
    replicaBatchSize = 3000
    replicaBatchTime = 5000

iterations = list(range(1, int(numIter) + 1))
arrivals = []

if setting == "LAN":
    arrivals = [500, 1000, 2000, 5000, 10000, 20000, 30000, 40000, 50000, 80000, 100000, 110000, 112000, 115000, 120000,
                130000, 150000, 180000, 200000]

if setting == "WAN":
    arrivals = [200, 1000, 5000, 10000, 15000, 20000, 25000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000]

propTime = 0


def getEPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["epaxos-exec", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/epaxos/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/execution/"
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

        record = ["epaxos-commit", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/epaxos/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/commit/"
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

def getRabiaSummary():
    l_records = []
    for arrival in arrivals:
        record = ["rabia", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/rabia/" + str(arrival) + "/" + str(int(300)) \
                   + "/" + setting + "/" + str(iteration) + "/execution/"
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


def getPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["paxos-v2", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/paxos-v2/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/execution/"
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


headers = ["algo", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

ePaxosSummary = getEPaxosSummary()
paxosV1Summary = getPaxosV1Summary()
paxosV2Summary = getPaxosV2Summary()
quePaxaSummary = getQuePaxaSummary()
rabiaSummary = getRabiaSummary()

records = records + ePaxosSummary + paxosV1Summary + paxosV2Summary + quePaxaSummary + rabiaSummary

import csv

with open("experiments/best-case/logs/summary.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)