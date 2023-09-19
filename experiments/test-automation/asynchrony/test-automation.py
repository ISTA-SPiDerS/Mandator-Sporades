import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir)
from execute import *


os.system("/bin/bash experiments/setup-5/setup.sh")

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


for iteration in [1,2,3]:
    for arrival in arrivals:

        # paxos and raft

        for algo in ["paxos", "raft"]:
            params = {}
            params["scenario"]=scenario
            params["arrival"]=str(arrival)
            params["replicaBatchSize"]=replicaBatchSize
            params["replicaBatchTime"] = replicaBatchTime
            params["clientBatchSize"]=clientBatchSize
            params["clientBatchTime"]=clientBatchTime
            params["setting"]=setting
            params["pipelineLength"]=pipelineLength
            params["algo"]=algo
            params["asyncTimeout"]=asyncTimeout
            params["benchmarkMode"]=benchmarkMode
            params["asyncTimeEpochSize"]=asyncTimeEpochSize
            params["viewTimeout"]=viewTimeout
            params["clientWindow"]=str(1000)
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(iteration)
            runPaxosRaft(params)

        # sporades

        params = {}
        params["scenario"]=scenario
        params["arrival"]=str(arrival)
        params["replicaBatchSize"]=replicaBatchSize
        params["replicaBatchTime"]=replicaBatchTime
        params["clientBatchSize"]=clientBatchSize
        params["clientBatchTime"]=clientBatchTime
        params["clientWindow"]=str(1000)
        params["asyncSimTimeout"]=asyncTimeout
        params["asyncTimeEpochSize"]=asyncTimeEpochSize
        params["benchmarkMode"]=benchmarkMode
        params["viewTimeout"]=viewTimeout
        params["setting"]=setting
        params["networkBatchTime"]=str(0)
        params["pipelineLength"]=pipelineLength
        params["collectClientLogs"]=collectClientLogs
        params["isLeaderKill"]=isLeaderKill
        params["iteration"]=str(iteration)
        runSporades(params)

        for algo in ["async"]:
            # mandator
            params={}
            params["scenario"]=scenario
            params["arrival"]=str(arrival)
            params["replicaBatchSize"]=replicaBatchSize
            params["replicaBatchTime"]=replicaBatchTime
            params["setting"]=setting
            params["algo"]=algo
            params["networkBatchTime"]=str(30)
            params["clientWindow"]=str(10000)
            params["asyncSimTime"]=asyncTimeout
            params["clientBatchSize"]=clientBatchSize
            params["clientBatchTime"]=clientBatchTime
            params["benchmarkMode"]=benchmarkMode
            params["broadcastMode"]=str(1)
            params["asyncTimeEpochSize"]=asyncTimeEpochSize
            params["viewTimeout"]=viewTimeout
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(iteration)
            runMandator(params)

