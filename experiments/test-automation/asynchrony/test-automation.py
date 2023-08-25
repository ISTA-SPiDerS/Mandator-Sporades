import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir)
from execute import *
grandParentdir = os.path.dirname(parentdir)
sys.path.append(grandParentdir + "/python")
from performance_extract import *


os.system("/bin/bash experiments/setup-5/setup.sh")


for asyncTimeout in [str(500)]:
    for asyncTimeEpochSize in [str(200), str(400)]:
        scenario="asynchrony"
        replicaBatchSize=str(3000)
        replicaBatchTime=str(5000)
        clientBatchSize=str(50)
        clientBatchTime=str(1000)
        setting="WAN"
        pipelineLength=str(10)
        benchmarkMode=str(0)
        viewTimeout=str(500000)
        clientWindow=str(10000)
        collectClientLogs="no"
        isLeaderKill="no"


        MIN_ARRIVAL = 100
        MAX_ARRIVAL = 40000
        INIT_GUAGE = int((MAX_ARRIVAL - MIN_ARRIVAL)/2)

        # multi-paxos and raft

        for algo in ["paxos", "raft"]:

            print("Starting "+algo)
            sys.stdout.flush()

            final_throughput = 0

            start = MAX_ARRIVAL
            gauge = INIT_GUAGE
            last_throughput = -1
            iter_num = 0
            while gauge > 400:
                iter_num = iter_num +1
                arrival = int(start)

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
                params["iteration"]=str(1)
                runPaxosRaft(params)
                throughput = getPaxosRaftPerformance("experiments/"+scenario+"/logs/paxos_raft/" + str(arrival) + "/"+ replicaBatchSize + "/"+ replicaBatchTime + "/"+ clientBatchSize + "/"+ clientBatchTime + "/"+ setting + "/"+ pipelineLength + "/"+ algo + "/"+ asyncTimeout + "/"+ benchmarkMode + "/"+ asyncTimeEpochSize + "/"+ viewTimeout + "/"+ clientWindow + "/"+ str(1) + "/execution/", 21, 5)[0]
                print( "async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize)+ " --- " + algo + " - iteration: " + str(iter_num)+", throughput: "+str(throughput)+", arrival: "+str(arrival*5))
                sys.stdout.flush()

                if throughput > last_throughput:
                    last_throughput = throughput

                if float(throughput) >= 0.8 * (arrival *5):
                    start = start + gauge
                    gauge = gauge / 2
                    continue
                if float(throughput) < 0.8 * (arrival *5):
                    start = start - gauge
                    if start <= 0:
                        start = 1000
                    gauge = gauge / 2
                    continue

            final_throughput  = last_throughput

            print("async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize)+"---" + algo + " final throughput - "+str(final_throughput))
            sys.stdout.flush()

        # mandator
        print("Starting Mandator Sporades")
        sys.stdout.flush()
        final_throughput = 0

        start = MAX_ARRIVAL
        gauge = INIT_GUAGE
        last_throughput = -1
        iter_num = 0
        while gauge > 400:
            iter_num = iter_num +1
            arrival = int(start)
            # run mandator for this configuration
            params = {}
            params["scenario"]=scenario
            params["arrival"]=str(arrival)
            params["replicaBatchSize"]=replicaBatchSize
            params["replicaBatchTime"]=replicaBatchTime
            params["setting"]=setting
            params["algo"]="async"
            params["networkBatchTime"]=str(10)
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
            params["iteration"]=str(1)
            runMandator(params)
            throughput = getManatorSporadesPerformance("experiments/"+scenario+"/logs/mandator"+"/" + str(arrival)+"/"+ replicaBatchSize+"/"+ replicaBatchTime+"/"+ setting+"/"+ "async"+"/"+ str(10)+"/"+ clientWindow+"/"+ asyncTimeout+"/"+ clientBatchSize+"/"+ clientBatchTime+"/"+ benchmarkMode+"/"+ str(1)+"/"+ asyncTimeEpochSize+"/"+ viewTimeout+"/"+ str(1)+"/execution/", 21, 5)[0]
            print("async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize)+ "Mandator Sporades - iteration: " + str(iter_num)+", throughput: "+str(throughput)+",  arrival: "+str(arrival*5))
            sys.stdout.flush()

            if throughput > last_throughput:
                last_throughput = throughput

            if float(throughput) >= 0.8 * (arrival *5):
                start = start + gauge
                gauge = gauge / 2
                continue
            if float(throughput) < 0.8 * (arrival *5):
                start = start - gauge
                if start <= 0:
                    start = 1000
                gauge = gauge / 2
                continue

        final_throughput  = last_throughput

        print("async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize)+ " --- mandator sporades final: throughput- "+str(final_throughput))
        sys.stdout.flush()

        # sporades
        print("Starting Pipelined Sporades")
        sys.stdout.flush()

        final_throughput = 0

        start = MAX_ARRIVAL
        gauge = INIT_GUAGE
        last_throughput = -1
        iter_num = 0
        while gauge > 400:
            iter_num = iter_num +1
            arrival = int(start)
            # run sporades for this configuration
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
            params["networkBatchTime"]=str(0)
            params["pipelineLength"]=pipelineLength
            params["collectClientLogs"]=collectClientLogs
            params["isLeaderKill"]=isLeaderKill
            params["iteration"]=str(1)
            runSporades(params)
            throughput = getManatorSporadesPerformance("experiments/"+scenario+"/logs/sporades/"+str(arrival)+"/"+replicaBatchSize+"/"+replicaBatchTime+"/"+clientBatchSize+"/"+clientBatchTime+"/"+clientWindow+"/"+asyncTimeout+"/"+asyncTimeEpochSize+"/"+benchmarkMode+"/"+viewTimeout+"/"+setting+"/"+str(0)+"/"+pipelineLength+"/"+str(1)+"/execution/", 21, 5)[0]
            print("async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize) + "Pipelined Sporades - iteration: " + str(iter_num)+", throughput: "+str(throughput)+",  arrival: "+str(arrival*5))
            sys.stdout.flush()

            if throughput > last_throughput:
                last_throughput = throughput

            if float(throughput) >= 0.8 * (arrival *5):
                start = start + gauge
                gauge = gauge / 2
                continue
            if float(throughput) < 0.8 * (arrival *5):
                start = start - gauge
                if start <= 0:
                    start = 1000
                gauge = gauge / 2
                continue

        final_throughput  = last_throughput

        print("async timeout:" + str(asyncTimeout)+", asyncTimeEpochSize:" + str(asyncTimeEpochSize)+"--- pipelined sporades  final: throughput- "+str(final_throughput))
        sys.stdout.flush()


print("-----Experiment Finished-----")