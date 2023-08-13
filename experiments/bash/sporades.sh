scenario=$1
arrival=$1
replicaBatchSize=$2
replicaBatchTime=$3
clientBatchSize=$4
clientBatchTime=$5
clientWindow=$6
asyncSimTimeout=$7
asyncTimeEpochSize=$8
benchmarkMode=$9 # 0 or 1
viewTimeout=${10}
setting=${11} # LAN or WAN
networkBatchTime=${12}
pipelineLength=${13}
collectClientLogs=${14}
isLeaderKill=${15}
iteration=${16}


pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

remote_algo_path="/mandator/binary/pipe_replica"
remote_ctl_path="/mandator/binary/pipe_client"


remote_config_path="/home/${user_name}/mandator/binary/mandator-sporades.yml"

echo "Starting test"

output_path="${pwd}/experiments/${scenario}/logs/sporades/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${clientBatchSize}/${clientBatchTime}/${clientWindow}/${asyncSimTimeout}/${asyncTimeEpochSize}/${benchmarkMode}/${viewTimeout}/${setting}/${networkBatchTime}/${pipelineLength}/${iteration}/execution/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

sleep 2

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 1 --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}1.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 2 --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}2.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 3 --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}3.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 4 --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}4.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path} --asyncSimTimeout ${asyncSimTimeout} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --isAsyncSim  --logFilePath ${remote_log_path} --name 5 --networkbatchTime  ${networkBatchTime} --pipelineLength ${pipelineLength} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}5.log &

sleep 10

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log &
echo "Sent initial status"

sleep 5

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 3" >${output_path}status3.log &
echo "Sent consensus start up"

sleep 20

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 21 --requestType request  --window ${clientWindow}" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 22 --requestType request  --window ${clientWindow}" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 23 --requestType request  --window ${clientWindow}" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 24 --requestType request  --window ${clientWindow}" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --logFilePath ${remote_log_path} --name 25 --requestType request  --window ${clientWindow}" >${output_path}25.log &

sleep 15

if [[ "${isLeaderKill}" == "yes" ]]
then
  echo "killing the first leader"
  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "${replica1}" "${kill_command}"
fi

sleep 110

echo "Completed Client[s]"

if [[ "${collectClientLogs}" == "yes" ]]
then
  echo "collecting client logs"
  scp -i ${cert} ${client1}:${remote_log_path}21.txt ${output_path}21.txt
  scp -i ${cert} ${client2}:${remote_log_path}22.txt ${output_path}22.txt
  scp -i ${cert} ${client3}:${remote_log_path}23.txt ${output_path}23.txt
  scp -i ${cert} ${client4}:${remote_log_path}24.txt ${output_path}24.txt
  scp -i ${cert} ${client5}:${remote_log_path}25.txt ${output_path}25.txt
fi

echo "Finish test"
