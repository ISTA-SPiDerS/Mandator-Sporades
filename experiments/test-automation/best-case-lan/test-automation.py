import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir)
from execute import *


os.system("/bin/bash experiments/setup-5/setup.sh")

arrivals = [500, 5000, 10000, 15000, 20000, 25000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000]

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
            params["clientWindow"]=clientWindow
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(iteration)
            runPaxosRaft(params)

        # sporades
        for networkNatchTime in [1,2,3]:
            params = {}
            params["scenario"]=scenario
            params["arrival"]=str(arrival)
            params["replicaBatchSize"]=replicaBatchSize
            params["replicaBatchTime"]=replicaBatchTime
            params["clientBatchSize"]=clientBatchSize
            params["clientBatchTime"]=clientBatchTime
            params["clientWindow"]=clientWindow
            params["asyncSimTimeout"]=asyncTimeout
            params["asyncTimeEpochSize"]=asyncTimeEpochSize
            params["benchmarkMode"]=benchmarkMode
            params["viewTimeout"]=viewTimeout
            params["setting"]=setting
            params["networkBatchTime"]=str(networkNatchTime)
            params["pipelineLength"]=pipelineLength
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(iteration)
            runSporades(params)

        # rabia
        params={}
        params["scenario"]=scenario
        params["arrivalRate"]=str(arrival)
        params["ProxyBatchSize"]=str(300)
        params["ProxyBatchTimeout"]=replicaBatchTime
        params["setting"]=setting
        params["ClientBatchSize"]=clientBatchSize
        params["isLeaderKill"]=isLeaderKill
        params["collectClientLogs"]=collectClientLogs
        params["iteration"]=str(iteration)
        runRabia(params)

        # epaxos

        params = {}
        params["scenario"]=scenario
        params["arrival"]=str(arrival)
        params["replicaBatchSize"]=replicaBatchSize
        params["replicaBatchTime"]=replicaBatchTime
        params["setting"]=setting
        params["pipelineLength"]=pipelineLength
        params["conflicts"]=str(2)
        params["clientWindow"]=clientWindow
        params["clientBatchSize"]=clientBatchSize
        params["collectClientLogs"]=collectClientLogs
        params["isLeaderKill"]=isLeaderKill
        params["iteration"] = str(iteration)
        runEPaxos(params)

