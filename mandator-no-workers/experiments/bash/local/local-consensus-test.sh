# A local test that tests the consensus layer by sending client requests and printing the mempool and consensus logs while simulating asynchrony using SIGSTOP and SIGCONT
arrivalRate=$1
algo=$2
asyncBatchTime=$3
attackTime=$4
viewTimeoutTime=$5

mage generate && mage build

replica_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/"

rm -r ${output_path}
mkdir ${output_path}

pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client

echo "Killed previously running instances"

nohup ./${replica_path} --name 1 --consAlgo "${algo}"  --window 1000   --debugOn --debugLevel 5 --viewTimeout "${viewTimeoutTime}"  --asyncBatchTime "${asyncBatchTime}" >${output_path}1.log &
nohup ./${replica_path} --name 2 --consAlgo "${algo}"  --window 1000   --debugOn --debugLevel 5 --viewTimeout "${viewTimeoutTime}"  --asyncBatchTime "${asyncBatchTime}" >${output_path}2.log &
nohup ./${replica_path} --name 3 --consAlgo "${algo}"  --window 1000   --debugOn --debugLevel 5 --viewTimeout "${viewTimeoutTime}"  --asyncBatchTime "${asyncBatchTime}" >${output_path}3.log &

echo "Started 3 replicas"

sleep 5

./${ctl_path} --name 11 --requestType status --operationType 1  --debugOn --debugLevel 1 >${output_path}status1.log

sleep 5

echo "sent initial status"

./${ctl_path} --name 11 --requestType status --operationType 3  --debugOn --debugLevel 1 >${output_path}status3.log

sleep 5

echo "sent consensus start up status"

nohup python3 experiments/python/crash-recovery-test.py logs/1.log logs/2.log logs/3.log "${attackTime}"> ${output_path}python-consensus-crash.log &
echo "started crash recovery script"

echo "starting clients"

nohup ./${ctl_path} --name 11 --requestType request --defaultReplica 1  --debugOn --debugLevel 5 --arrivalRate "${arrivalRate}" >${output_path}11.log &
nohup ./${ctl_path} --name 12 --requestType request --defaultReplica 2  --debugOn --debugLevel 5 --arrivalRate "${arrivalRate}" >${output_path}12.log &
nohup ./${ctl_path} --name 13 --requestType request --defaultReplica 3  --debugOn --debugLevel 5 --arrivalRate "${arrivalRate}" >${output_path}13.log &

sleep 180

echo "finished running clients"


./${ctl_path} --name 11 --requestType status --operationType 2  --debugOn --debugLevel 1 >${output_path}status2.log

echo "sent status to print logs"

sleep 30

pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client

python3 experiments/python/overlay-test.py logs/1-mem-pool.txt  logs/2-mem-pool.txt  logs/3-mem-pool.txt  > ${output_path}python-mem-pool.log
python3 experiments/python/overlay-test.py logs/1-consensus.txt logs/2-consensus.txt logs/3-consensus.txt > ${output_path}python-consensus.log
python3 experiments/python/commit-log-test.py logs/1.log logs/2.log logs/3.log

echo "Killed instances"

echo "Finish test"