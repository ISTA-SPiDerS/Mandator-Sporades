scenario="scalablity-7"
arrival=$1
replicaBatchSize=3000
replicaBatchTime=5000
clientBatchSize=50
clientBatchTime=1000
clientWindow=1000
asyncSimTimeout=0
asyncTimeEpochSize=500
benchmarkMode=0
viewTimeout=300000000
setting="WAN" # LAN or WAN
networkBatchTime=0
pipelineLength=10
collectClientLogs="no"
isLeaderKill="no"
iteration=1


pwd=$(pwd)
. "${pwd}"/experiments/setup-7/ip.sh

remote_algo_path="/mandator/binary/pipe_replica"
remote_ctl_path="/mandator/binary/pipe_client"


remote_config_path="/home/${user_name}/mandator/binary/mandator-sporades.yml"

echo "Starting test"

output_path="${pwd}/experiments/${scenario}/logs/sporades/${arrival}/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert}  -n -f "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

sleep 2

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 1  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}1.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 2  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}2.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 3  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}3.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 4  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}4.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 5  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}5.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica6}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 6  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}6.log  &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica7}  ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 7  --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}7.log  &


sleep 10

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log &
echo "Sent initial status"

sleep 15

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 3" >${output_path}status3.log &
echo "Sent consensus start up"

sleep 30

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 21 --requestType request  --window ${clientWindow}" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 22 --requestType request  --window ${clientWindow}" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 23 --requestType request  --window ${clientWindow}" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 24 --requestType request  --window ${clientWindow}" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 25 --requestType request  --window ${clientWindow}" >${output_path}25.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client6}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 26 --requestType request  --window ${clientWindow}" >${output_path}26.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client7}  ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 27 --requestType request  --window ${clientWindow}" >${output_path}27.log &


sleep 115

echo "Completed Client[s]"

echo "Finish test"
