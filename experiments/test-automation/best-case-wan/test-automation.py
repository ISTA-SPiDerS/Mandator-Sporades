import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir)
from execute import *


os.system("/bin/bash experiments/setup-5/setup.sh")

arrivals = [500, 10000, 20000, 30000, 35000, 40000, 45000, 50000, 60000, 70000, 80000, 90000]

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


for iteration in [1,2,3]:
    for arrival in arrivals:
        # paxos

        params = {}
        params["scenario"]=scenario
        params["arrival"]=str(arrival)
        params["replicaBatchSize"]=replicaBatchSize
        params["replicaBatchTime"] = replicaBatchTime
        params["clientBatchSize"]=clientBatchSize
        params["clientBatchTime"]=clientBatchTime
        params["setting"]=setting
        params["pipelineLength"]=pipelineLength
        params["algo"]="paxos"
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
        params["networkBatchTime"]=str(5)
        params["pipelineLength"]=pipelineLength
        params["collectClientLogs"]=collectClientLogs
        params["isLeaderKill"]=isLeaderKill
        params["iteration"]=str(iteration)
        runSporades(params)

        # epaxos

        params = {}
        params["scenario"]=scenario
        params["arrival"]=str(arrival)
        params["replicaBatchSize"]=replicaBatchSize
        params["replicaBatchTime"]=replicaBatchTime
        params["setting"]=setting
        params["pipelineLength"]=pipelineLength
        params["conflicts"]=str(2)
        params["clientWindow"]=str(1000)
        params["clientBatchSize"]=clientBatchSize
        params["collectClientLogs"]=collectClientLogs
        params["isLeaderKill"]=isLeaderKill
        params["iteration"] = str(iteration)
        runEPaxos(params)

        # mandator
        params={}
        params["scenario"]=scenario
        params["arrival"]=str(arrival)
        params["replicaBatchSize"]=replicaBatchSize
        params["replicaBatchTime"]=replicaBatchTime
        params["setting"]=setting
        params["algo"]="async"
        params["networkBatchTime"]=str(30)
        params["clientWindow"]=clientWindow
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

