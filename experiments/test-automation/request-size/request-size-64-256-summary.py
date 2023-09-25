import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *

scenario="request-size"

arrivals = [100, 200, 500,  1000, 2000, 5000, 10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000]


iterations = [1,2,3]
headers = ["algo", "size", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]


def getPaxosRaftSummary():
    l_records = []
    for arrival in arrivals:
        for size in [64, 256]:

            if size == 256 and arrival > 10000:
                continue
            else:
                record = ["paxos", str(size), str(arrival * int(5))]
                throughput, latency, nine9, err = [], [], [], []
                for iteration in iterations:
                    root = "experiments"+"/"+scenario+"/logs/paxos/"+str(arrival)+"/"+str(size)+"/"+ str(iteration) +"/"
                    t, l, n, e = getPaxosRaftPerformance(root, 21, int(5))
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
        for size in [64, 256]:

            if size == 256 and arrival > 10000:
                continue
            else:
                record = ["sporades", str(size), str(arrival * int(5))]
                throughput, latency, nine9, err = [], [], [], []
                for iteration in iterations:
                    root = "experiments"+"/"+scenario+"/logs/sporades/"+str(arrival)+"/"+str(size)+"/"+ str(iteration) +"/"
                    t, l, n, e = getManatorSporadesPerformance(root, 21, int(5))
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
        for size in [64, 256]:
            record = ["mandator", str(size), str(arrival * int(5))]
            throughput, latency, nine9, err = [], [], [], []
            for iteration in iterations:
                root = "experiments"+"/"+scenario+"/logs/mandator/"+str(arrival)+"/"+str(size)+"/"+ str(iteration) +"/"
                t, l, n, e = getManatorSporadesPerformance(root, 21, int(5))
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

with open("experiments/"+scenario+"/logs/"+scenario+"-summary-"+str(datetime.now())+".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)