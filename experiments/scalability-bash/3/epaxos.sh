scenario="scalability-3"
arrival=$1
replicaBatchSize=3000
replicaBatchTime=5000
setting="WAN" # LAN or WAN
pipelineLength=10
conflicts=2
clientWindow=1000
clientBatchSize=50
iteration=$2

pwd=$(pwd)
. "${pwd}"/experiments/setup-3/ip.sh

remote_algo_path="/mandator/binary/epaxos_server"
remote_ctl_path="/mandator/binary/epaxos_client"
remote_master_path="/mandator/binary/epaxos_master"

echo "Starting execution latency test"

output_path="${pwd}/experiments/${scenario}/logs/epaxos/${arrival}/${iteration}/execution/"
rm -r "${output_path}" ; mkdir -p "${output_path}"
echo "Removed old local log files"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

sleep 2

echo "starting master"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 3 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -exec  -dreply -e" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -exec  -dreply -e" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -exec  -dreply -e" >${output_path}4.log &
sleep 5


echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} -name 7     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 0 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} -name 8     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 1 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} -name 9     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 2 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}23.log &


sleep 120

echo "Completed Client[s]"

echo "finished execution latency test"

######################
pwd=$(pwd)
. "${pwd}"/experiments/setup-3/ip.sh

echo "Starting commit latency test"

output_path="${pwd}/experiments/${scenario}/logs/epaxos/${arrival}/${iteration}/commit/"
rm -r "${output_path}" ; mkdir -p "${output_path}"
echo "Removed old local log files"


for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert}  -n -f "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

sleep 2

echo "starting master"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 3 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -e" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -e" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}    -e" >${output_path}4.log &
sleep 5


echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} -name 7     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 0 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} -name 8     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 1 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} -name 9     -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 2 -logFilePath ${remote_log_path} --window ${clientWindow}" >${output_path}23.log &

sleep 150

echo "Completed Client[s]"

echo "finished commit latency test"

