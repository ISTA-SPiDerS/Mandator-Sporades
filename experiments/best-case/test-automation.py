import os
import sys
from datetime import datetime

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

setting = sys.argv[1]
numIter = sys.argv[2]

replicaBatchSize = 3000
replicaBatchTime = 5000

iterations = list(range(1, int(numIter) + 1))
pipelines = [1, 10]

arrivals = [500, 25000, 35000, 40000, 60000, 80000, 100000]

for iteration in iterations:
    for arrival in arrivals:
        for pipeline in pipelines:
            print(str(datetime.now()) + ": iteration: " + str(iteration) + ", arrival:" + str(
                arrival) + ", pipeline:" + str(pipeline) + "\n")
            os.system(
                "/bin/bash experiments/best-case/epaxos.sh " + str(int(arrival)) + " "
                + str(replicaBatchSize) + " "
                + str(replicaBatchTime) + " "
                + setting + " "
                + str(pipeline) + " "
                + str(iteration))
            for algo in ["async", "paxos"]:
                if pipeline != 1:
                    continue
                os.system(
                    "/bin/bash experiments/best-case/mandator.sh " + str(int(arrival)) + " "
                    + str(replicaBatchSize) + " "
                    + str(replicaBatchTime) + " "
                    + setting + " " + algo + " " + str(3) + " "
                    + str(iteration))
            for algo in ["paxos", "raft"]:
                if algo == "raft" and pipeline != 1:
                    continue
                os.system(
                    "/bin/bash experiments/best-case/paxos_raft.sh " + str(int(arrival)) + " "
                    + str(replicaBatchSize) + " "
                    + str(replicaBatchTime) + " "
                    + setting + " " + str(pipeline) + " " + algo + " " +
                    str(iteration))

            if pipeline == 1:
                os.system(
                    "/bin/bash experiments/best-case/rabia.sh " + str(int(arrival)) + " "
                    + str(replicaBatchSize) + " "
                    + setting + " "
                    + str(iteration))

            os.system(
                "/bin/bash experiments/best-case/sporades.sh " + str(int(arrival)) + " "
                + str(replicaBatchSize) + " "
                + str(replicaBatchTime) + " "
                + setting + " " + str(3) + " " + str(pipeline) + " "
                + str(iteration))
