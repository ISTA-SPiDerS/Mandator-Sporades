# a test script to run 3 replicas. 6 workers and 3 clients in a LAN

arrivalRate=$1
asyncBatchTime=$2
algo=$3

remote_algo_path="/async/temp/replica/bin/replica"
remote_ctl_path="/async/temp/client/bin/client"
remote_worker_path="/async/temp/worker/bin/child"
remote_log_path="/home/pasindu/async/temp/logs/"
remote_config_path="/home/pasindu/async/temp/configuration/remote/configuration.yml"
remote_config_worker_map_path="/home/pasindu/async/temp/configuration/remote/workermapconfiguration.yml"
output_path="logs/"

replica1=pasindu@dedis-140.icsil1.epfl.ch
replica1_cert="/home/pasindu/Pictures/pasindu_rsa"
replica2=pasindu@dedis-141.icsil1.epfl.ch
replica2_cert="/home/pasindu/Pictures/pasindu_rsa"
replica3=pasindu@dedis-142.icsil1.epfl.ch
replica3_cert="/home/pasindu/Pictures/pasindu_rsa" 

client1=pasindu@dedis-145.icsil1.epfl.ch
client1_cert="/home/pasindu/Pictures/pasindu_rsa"
client2=pasindu@dedis-146.icsil1.epfl.ch
client2_cert="/home/pasindu/Pictures/pasindu_rsa"
client3=pasindu@dedis-147.icsil1.epfl.ch
client3_cert="/home/pasindu/Pictures/pasindu_rsa"

rm -r ${output_path}; mkdir ${output_path}

echo "Removed old log files"

sshpass ssh -i ${replica1_cert} ${replica1} "rm ${remote_log_path}1-mem-pool.txt; rm ${remote_log_path}1-consensus.txt"
sshpass ssh -i ${replica2_cert} ${replica2} "rm ${remote_log_path}2-mem-pool.txt; rm ${remote_log_path}2-consensus.txt"
sshpass ssh -i ${replica3_cert} ${replica3} "rm ${remote_log_path}3-mem-pool.txt; rm ${remote_log_path}3-consensus.txt"

sshpass ssh -i ${client1_cert} ${client1} "rm ${remote_log_path}11.txt"
sshpass ssh -i ${client2_cert} ${client2} "rm ${remote_log_path}12.txt"
sshpass ssh -i ${client3_cert} ${client3} "rm ${remote_log_path}13.txt"

sleep 5
echo "Removed past log files in remote servers"

kill_command="pkill replica ; pkill client; pkill child; pkill child"

sshpass ssh -i ${replica1_cert} ${replica1} "${kill_command}"
sshpass ssh -i ${replica2_cert} ${replica2} "${kill_command}"
sshpass ssh -i ${replica3_cert} ${replica3} "${kill_command}"

sshpass ssh -i ${client1_cert} ${client1} "${kill_command}"
sshpass ssh -i ${client2_cert} ${client2} "${kill_command}"
sshpass ssh -i ${client3_cert} ${client3} "${kill_command}"

echo "killed previous running instances"

sleep 5

echo "starting replicas"

nohup sshpass ssh -i ${replica1_cert} -n -f ${replica1} ".${remote_algo_path} --name 1 --consAlgo ${algo} --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 --asyncBatchTime ${asyncBatchTime} --benchmarkMode 1" >${output_path}1.log &
nohup sshpass ssh -i ${replica2_cert} -n -f ${replica2} ".${remote_algo_path} --name 2 --consAlgo ${algo} --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 --asyncBatchTime ${asyncBatchTime} --benchmarkMode 1" >${output_path}2.log &
nohup sshpass ssh -i ${replica3_cert} -n -f ${replica3} ".${remote_algo_path} --name 3 --consAlgo ${algo} --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 --asyncBatchTime ${asyncBatchTime} --benchmarkMode 1" >${output_path}3.log &

echo "starting workers"

nohup sshpass ssh -i ${replica1_cert} -n -f ${replica1} ".${remote_worker_path} --name 21  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}21.log &
nohup sshpass ssh -i ${replica1_cert} -n -f ${replica1} ".${remote_worker_path} --name 22  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}22.log &
nohup sshpass ssh -i ${replica2_cert} -n -f ${replica2} ".${remote_worker_path} --name 23  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}23.log &
nohup sshpass ssh -i ${replica2_cert} -n -f ${replica2} ".${remote_worker_path} --name 24  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}24.log &
nohup sshpass ssh -i ${replica3_cert} -n -f ${replica3} ".${remote_worker_path} --name 25  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}25.log &
nohup sshpass ssh -i ${replica3_cert} -n -f ${replica3} ".${remote_worker_path} --name 26  --debugLevel 4 --window 1000 --config ${remote_config_path} --workerMapconfig ${remote_config_worker_map_path}  --logFilePath ${remote_log_path} --batchSize 100 --batchTime 5000 " >${output_path}26.log &

echo "Started servers"

sleep 5

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 11 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log
echo "Sent initial status"

sleep 5

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 11 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 3" >${output_path}status3.log
echo "Sent consensus start up"

sleep 5

echo "Starting client[s]"

nohup sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 11 --defaultReplicas 21,22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 100 --batchTime 5000 " >${output_path}11.log &
nohup sshpass ssh -i ${client2_cert} ${client2} ".${remote_ctl_path} --name 12 --defaultReplicas 23,24 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 100 --batchTime 5000 " >${output_path}12.log &
      sshpass ssh -i ${client3_cert} ${client3} ".${remote_ctl_path} --name 13 --defaultReplicas 25,26 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 100 --batchTime 5000 " >${output_path}13.log

sleep 20

echo "Completed Client[s]"

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 11 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 2" >${output_path}status2.log
echo "Sent log printing"

sleep 50

scp -i ${replica1_cert} ${replica1}:${remote_log_path}1-mem-pool.txt ${output_path}1-mem-pool.txt
scp -i ${replica1_cert} ${replica1}:${remote_log_path}1-consensus.txt ${output_path}1-consensus.txt

scp -i ${replica2_cert} ${replica2}:${remote_log_path}2-mem-pool.txt ${output_path}2-mem-pool.txt
scp -i ${replica2_cert} ${replica2}:${remote_log_path}2-consensus.txt ${output_path}2-consensus.txt

scp -i ${replica3_cert} ${replica3}:${remote_log_path}3-mem-pool.txt ${output_path}3-mem-pool.txt
scp -i ${replica3_cert} ${replica3}:${remote_log_path}3-consensus.txt ${output_path}3-consensus.txt

scp -i ${client1_cert} ${client1}:${remote_log_path}11.txt ${output_path}11.txt
scp -i ${client2_cert} ${client2}:${remote_log_path}12.txt ${output_path}12.txt
scp -i ${client3_cert} ${client3}:${remote_log_path}13.txt ${output_path}13.txt

echo "Copied all the files to local machine"

sleep 5

python3 experiments/python/overlay-test.py logs/1-mem-pool.txt  logs/2-mem-pool.txt  logs/3-mem-pool.txt  > ${output_path}python-mem-pool.log
python3 experiments/python/overlay-test.py logs/1-consensus.txt logs/2-consensus.txt logs/3-consensus.txt > ${output_path}python-consensus.log

sshpass ssh -i ${replica1_cert} ${replica1} "pkill replica; pkill child; pkill child; "
sshpass ssh -i ${replica2_cert} ${replica2} "pkill replica; pkill child; pkill child; "
sshpass ssh -i ${replica3_cert} ${replica3} "pkill replica; pkill child; pkill child; "

sshpass ssh -i ${client1_cert} ${client1} "pkill client"
sshpass ssh -i ${client2_cert} ${client2} "pkill client"
sshpass ssh -i ${client3_cert} ${client3} "pkill client"

echo "killed  instances"

echo "Finish test"