arrivalRate=$1
viewTimeoutTime=$2
batchTime=$3
batchSize=$4
pipelineLength=$5
window=$6

replica_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/${arrivalRate}/${viewTimeoutTime}/${batchTime}/${batchSize}/${pipelineLength}/${window}/"

rm -r ${output_path}
mkdir -p ${output_path}

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client

echo "Killed previously running instances"

nohup ./${replica_path} --name 1 --batchSize "${batchSize}" --batchTime "${batchTime}"   --debugOn      --debugLevel 1 --viewTimeout "${viewTimeoutTime}" --pipelineLength "${pipelineLength}" --logFilePath ${output_path} --asyncbatchTime 3 >${output_path}1.log &
nohup ./${replica_path} --name 2 --batchSize "${batchSize}" --batchTime "${batchTime}"   --debugOn      --debugLevel 1 --viewTimeout "${viewTimeoutTime}" --pipelineLength "${pipelineLength}" --logFilePath ${output_path} --asyncbatchTime 3  >${output_path}2.log &
nohup ./${replica_path} --name 3 --batchSize "${batchSize}" --batchTime "${batchTime}"   --debugOn      --debugLevel 1 --viewTimeout "${viewTimeoutTime}" --pipelineLength "${pipelineLength}" --logFilePath ${output_path} --asyncbatchTime 3  >${output_path}3.log &
nohup ./${replica_path} --name 4 --batchSize "${batchSize}" --batchTime "${batchTime}"   --debugOn      --debugLevel 1 --viewTimeout "${viewTimeoutTime}" --pipelineLength "${pipelineLength}" --logFilePath ${output_path} --asyncbatchTime 3  >${output_path}4.log &
nohup ./${replica_path} --name 5 --batchSize "${batchSize}" --batchTime "${batchTime}"   --debugOn      --debugLevel 1 --viewTimeout "${viewTimeoutTime}" --pipelineLength "${pipelineLength}" --logFilePath ${output_path} --asyncbatchTime 3  >${output_path}5.log &

echo "Started 5 replicas"

sleep 5

./${ctl_path} --name 21 --logFilePath ${output_path} --requestType status --operationType 1   --debugLevel 0 >${output_path}status1.log

sleep 5

echo "sent initial status"

./${ctl_path} --name 22 --logFilePath ${output_path} --requestType status --operationType 3   --debugLevel 0 >${output_path}status3.log

sleep 15

echo "sent consensus start up status"

echo "starting clients"

nohup ./${ctl_path} --name 21 --logFilePath ${output_path} --requestType request  --debugLevel 0 --batchSize  "${batchSize}" --batchTime "${batchTime}" --arrivalRate "${arrivalRate}" --window "${window}" >${output_path}21.log &
nohup ./${ctl_path} --name 22 --logFilePath ${output_path} --requestType request  --debugLevel 0 --batchSize  "${batchSize}" --batchTime "${batchTime}" --arrivalRate "${arrivalRate}" --window "${window}" >${output_path}22.log &
nohup ./${ctl_path} --name 23 --logFilePath ${output_path} --requestType request  --debugLevel 0 --batchSize  "${batchSize}" --batchTime "${batchTime}" --arrivalRate "${arrivalRate}" --window "${window}" >${output_path}23.log &
nohup ./${ctl_path} --name 24 --logFilePath ${output_path} --requestType request  --debugLevel 0 --batchSize  "${batchSize}" --batchTime "${batchTime}" --arrivalRate "${arrivalRate}" --window "${window}" >${output_path}24.log &
nohup ./${ctl_path} --name 25 --logFilePath ${output_path} --requestType request  --debugLevel 0 --batchSize  "${batchSize}" --batchTime "${batchTime}" --arrivalRate "${arrivalRate}" --window "${window}" >${output_path}25.log &

sleep 80

echo "finished running clients"


nohup ./${ctl_path} --name 21 --logFilePath ${output_path} --requestType status --operationType 2  --debugOn --debugLevel 0 >${output_path}status2.log &


echo "sent status to print logs"

sleep 30

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client

python3 integration-test/python/overlay-test.py ${output_path}/1-consensus.txt ${output_path}/2-consensus.txt ${output_path}/3-consensus.txt ${output_path}/4-consensus.txt ${output_path}/5-consensus.txt

echo "Killed instances"

echo "Finish test"
