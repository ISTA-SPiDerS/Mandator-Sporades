import os
import sys

setting = sys.argv[1]  # LAN or WAN
numIter = sys.argv[2]

if setting != "LAN" and setting != "WAN":
    exit("wrong input, input should be LAN/WAN")

if int(numIter) < 3:
    exit("at least 3 iterations needed")

os.system("/bin/bash experiments/setup-5/setup.sh")

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

for iteration in iterations:
    for arrival in arrivals:
        os.system(
            "/bin/bash experiments/best-case/epaxos.sh " + str(int(arrival)) + " "
            + str(replicaBatchSize) + " "
            + str(replicaBatchTime) + " "
            + setting + " "
            + str(iteration))

        os.system(
            "/bin/bash experiments/best-case/paxos-v1.sh " + str(int(arrival)) + " "
            + str(replicaBatchSize) + " "
            + str(replicaBatchTime) + " "
            + setting + " "
            + str(iteration))

        os.system(
            "/bin/bash experiments/best-case/paxos-v2.sh " + str(int(arrival)) + " "
            + str(replicaBatchSize) + " "
            + str(replicaBatchTime) + " "
            + setting + " "
            + str(iteration))

        os.system(
            "/bin/bash experiments/best-case/quepaxa.sh " + str(int(arrival)) + " "
            + str(replicaBatchSize) + " "
            + str(replicaBatchTime) + " "
            + setting + " "
            + str(iteration) + " "
            + str(0) + " " + str(0))


        os.system(
            "/bin/bash experiments/best-case/quepaxa.sh " + str(int(arrival)) + " "
            + str(replicaBatchSize) + " "
            + str(replicaBatchTime) + " "
            + setting + " "
            + str(iteration) + " "
            + str(1) + " " + str(0))

        os.system(
            "/bin/bash experiments/best-case/rabia.sh " + str(int(arrival)) + " "
            + str(300) + " "
            + setting + " "
            + str(iteration))
