scenario=$1
arrival=$2
replicaBatchSize=$3
replicaBatchTime=$4
setting=$5 # LAN or WAN
algo=$6
networkBatchTime=$7
clientWindow=$8
asyncSimTime=$9
clientBatchSize=${10}
clientBatchTime=${11}
benchmarkMode=${12} # 0 or 1
broadcastMode=${13} # 1 or 2
asyncTimeEpochSize=${14}
viewTimeout=${15}
collectClientLogs=${16}
isLeaderKill=${17}
iteration=${18}

pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

remote_algo_path="/mandator/binary/man_replica"
remote_ctl_path="/mandator/binary/man_client"


remote_config_path="/home/${user_name}/mandator/binary/mandator-sporades.yml"

echo "Starting test"

output_path="${pwd}/experiments/${scenario}/logs/mandator/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${setting}/${algo}/${networkBatchTime}/${clientWindow}/${asyncSimTime}/${clientBatchSize}/${clientBatchTime}/${benchmarkMode}/${broadcastMode}/${asyncTimeEpochSize}/${viewTimeout}/${iteration}/execution/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert}  -n -f "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done


nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path} --asyncSimTime ${asyncSimTime} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --consAlgo ${algo} --isAsync --logFilePath ${remote_log_path}  --mode ${broadcastMode} --name 1  --networkBatchTime  ${networkBatchTime} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}1.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path} --asyncSimTime ${asyncSimTime} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --consAlgo ${algo} --isAsync --logFilePath ${remote_log_path}  --mode ${broadcastMode} --name 2  --networkBatchTime  ${networkBatchTime} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}2.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path} --asyncSimTime ${asyncSimTime} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --consAlgo ${algo} --isAsync --logFilePath ${remote_log_path}  --mode ${broadcastMode} --name 3  --networkBatchTime  ${networkBatchTime} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}3.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path} --asyncSimTime ${asyncSimTime} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --consAlgo ${algo} --isAsync --logFilePath ${remote_log_path}  --mode ${broadcastMode} --name 4  --networkBatchTime  ${networkBatchTime} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}4.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path} --asyncSimTime ${asyncSimTime} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --benchmarkMode ${benchmarkMode} --config ${remote_config_path} --consAlgo ${algo} --isAsync --logFilePath ${remote_log_path}  --mode ${broadcastMode} --name 5  --networkBatchTime  ${networkBatchTime} --timeEpochSize ${asyncTimeEpochSize}  --viewTimeout ${viewTimeout} " >${output_path}5.log &


sleep 10

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log &
echo "Sent initial status"

sleep 5

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 3" >${output_path}status3.log &
echo "Sent consensus start up"

sleep 30

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --designatedReplica 1 --logFilePath ${remote_log_path} --name 21  --requestType request --window ${clientWindow} " >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --designatedReplica 2 --logFilePath ${remote_log_path} --name 22  --requestType request --window ${clientWindow} " >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --designatedReplica 3 --logFilePath ${remote_log_path} --name 23  --requestType request --window ${clientWindow} " >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --designatedReplica 4 --logFilePath ${remote_log_path} --name 24  --requestType request --window ${clientWindow} " >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path}  --arrivalRate ${arrival} --batchSize ${clientBatchSize} --batchTime ${clientBatchTime} --config ${remote_config_path} --designatedReplica 5 --logFilePath ${remote_log_path} --name 25  --requestType request --window ${clientWindow} " >${output_path}25.log &

sleep 15

if [[ "${isLeaderKill}" == "yes" ]]
then
  echo "killing the first leader"
  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "${replica1}" "${kill_command};${kill_command}"
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
