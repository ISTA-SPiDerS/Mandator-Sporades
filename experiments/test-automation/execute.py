import os

def runPaxosRaft(params):
    os.system("/bin/bash experiments/bash/paxos_raft.sh " +
    params["scenario"]+" "+
    params["arrival"]+" "+
    params["replicaBatchSize"]+" "+
    params["replicaBatchTime"]+" "+
    params["clientBatchSize"]+" "+
    params["clientBatchTime"]+" "+
    params["setting"]+" "+
    params["pipelineLength"]+" "+
    params["algo"]+" "+
    params["asyncTimeout"]+" "+
    params["benchmarkMode"]+" "+
    params["asyncTimeEpochSize"]+" "+
    params["viewTimeout"]+" "+
    params["clientWindow"]+" "+
    params["collectClientLogs"]+" "+
    params["isLeaderKill"]+" "+
    params["iteration"])


def runEPaxos(params):
    os.system("/bin/bash experiments/bash/epaxos.sh "+
    params["scenario"]+" "+
    params["arrival"]+" "+
    params["replicaBatchSize"]+" "+
    params["replicaBatchTime"]+" "+
    params["setting"]+" "+
    params["pipelineLength"]+" "+
    params["conflicts"]+" "+
    params["clientWindow"]+" "+
    params["clientBatchSize"]+" "+
    params["collectClientLogs"]+" "+
    params["isLeaderKill"]+" "+
    params["iteration"])


def runSporades(params):
    os.system("/bin/bash experiments/bash/sporades.sh "+
    params["scenario"]+" "+
    params["arrival"]+" "+
    params["replicaBatchSize"]+" "+
    params["replicaBatchTime"]+" "+
    params["clientBatchSize"]+" "+
    params["clientBatchTime"]+" "+
    params["clientWindow"]+" "+
    params["asyncSimTimeout"]+" "+
    params["asyncTimeEpochSize"]+" "+
    params["benchmarkMode"]+" "+
    params["viewTimeout"]+" "+
    params["setting"]+" "+
    params["networkBatchTime"]+" "+
    params["pipelineLength"]+" "+
    params["collectClientLogs"]+" "+
    params["isLeaderKill"]+" "+
    params["iteration"])

def runMandator(params):
    os.system("/bin/bash experiments/bash/mandator.sh "+
    params["scenario"]+" "+
    params["arrival"]+" "+
    params["replicaBatchSize"]+" "+
    params["replicaBatchTime"]+" "+
    params["setting"]+" "+
    params["algo"]+" "+
    params["networkBatchTime"]+" "+
    params["clientWindow"]+" "+
    params["asyncSimTime"]+" "+
    params["clientBatchSize"]+" "+
    params["clientBatchTime"]+" "+
    params["benchmarkMode"]+" "+
    params["broadcastMode"]+" "+
    params["asyncTimeEpochSize"]+" "+
    params["viewTimeout"]+" "+
    params["collectClientLogs"]+" "+
    params["isLeaderKill"]+" "+
    params["iteration"])

def runRabia(params):
    os.system("/bin/bash experiments/bash/rabia.sh "+
    params["scenario"]+" "+
    params["arrivalRate"]+" "+
    params["ProxyBatchSize"]+" "+
    params["ProxyBatchTimeout"]+" "+
    params["setting"]+" "+
    params["ClientBatchSize"]+" "+
    params["isLeaderKill"]+" "+
    params["collectClientLogs"]+" "+
    params["iteration"])
