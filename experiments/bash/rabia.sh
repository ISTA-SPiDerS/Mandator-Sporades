scenario=$1
arrivalRate=$2
ProxyBatchSize=$3
ProxyBatchTimeout=$4 # ms
setting=$5 # LAN or WAN
ClientBatchSize=$6
isLeaderKill=$7
collectClientLogs=$8
iteration=$9

pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

Controller=${replica1_ip}:8070
RCFolder="/home/${user_name}/mandator/"
NServers=5
NFaulty=2
NClients=5

ClientTimeout=60 # test duration
RC_Peers_N="${replica1_ip}:10090,${replica2_ip}:10091,${replica3_ip}:10092,${replica4_ip}:10093,${replica5_ip}:10094"

output_path="${pwd}/experiments/${scenario}/logs/rabia/${arrivalRate}/${ProxyBatchSize}/${ProxyBatchTimeout}/${setting}/${ClientBatchSize}/${iteration}/execution/"
rm -r "${output_path}" ; mkdir -p "${output_path}"
echo "Removed old log files"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert}   -n -f "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

Rabia_Path="/mandator/binary/rabia"

export_command="export LogFilePath=${remote_log_path} RC_Ctrl=${Controller} RC_Folder=${RCFolder} RC_LLevel="warn" Rabia_ClosedLoop=false Rabia_NServers=${NServers} Rabia_NFaulty=${NFaulty} Rabia_NClients=${NClients} Rabia_NConcurrency=1 Rabia_ClientBatchSize=${ClientBatchSize} Rabia_ClientTimeout=${ClientTimeout} Rabia_ClientThinkTime=0 Rabia_ClientNRequests=0 Rabia_ClientArrivalRate=${arrivalRate} Rabia_ProxyBatchSize=${ProxyBatchSize} Rabia_ProxyBatchTimeout=${ProxyBatchTimeout} Rabia_NetworkBatchSize=0 Rabia_NetworkBatchTimeout=0 RC_Peers=${RC_Peers_N} Rabia_StorageMode=0"

echo "starting replicas"

svr_export="export RC_Role=svr RC_Index=0 RC_SvrIp="${replica1_ip}" RC_PPort="9090" RC_NPort="10090""
nohup sshpass ssh -o "StrictHostKeyChecking no"  -i ${cert} -n -f ${replica1} "${svr_export} && ${export_command} &&  .${Rabia_Path}" >${output_path}0.log &

svr_export="export RC_Role=svr RC_Index=1 RC_SvrIp="${replica2_ip}" RC_PPort="9091" RC_NPort="10091""
nohup sshpass ssh -o "StrictHostKeyChecking no"  -i ${cert} -n -f ${replica2} "${svr_export} && ${export_command} &&  .${Rabia_Path}" >${output_path}1.log &

svr_export="export RC_Role=svr RC_Index=2 RC_SvrIp="${replica3_ip}" RC_PPort="9092" RC_NPort="10092""
nohup sshpass ssh -o "StrictHostKeyChecking no"  -i ${cert} -n -f ${replica3} "${svr_export} && ${export_command} &&  .${Rabia_Path}" >${output_path}2.log &

svr_export="export RC_Role=svr RC_Index=3 RC_SvrIp="${replica4_ip}" RC_PPort="9093" RC_NPort="10093""
nohup sshpass ssh -o "StrictHostKeyChecking no"  -i ${cert} -n -f ${replica4} "${svr_export} && ${export_command} &&  .${Rabia_Path}" >${output_path}3.log &

svr_export="export RC_Role=svr RC_Index=4 RC_SvrIp="${replica5_ip}" RC_PPort="9094" RC_NPort="10094""
nohup sshpass ssh -o "StrictHostKeyChecking no"  -i ${cert} -n -f ${replica5} "${svr_export} && ${export_command} &&  .${Rabia_Path}" >${output_path}4.log &

sleep 5

echo "Starting client[s]"

cli_export="export RC_Role=cli RC_Index=0 RC_Proxy="${replica1_ip}:9090""
nohup sshpass ssh -o "StrictHostKeyChecking no"   -i ${cert} -n -f ${client1} "${cli_export} && ${export_command} && .${Rabia_Path}" >${output_path}21.log &

cli_export="export RC_Role=cli RC_Index=1 RC_Proxy="${replica2_ip}:9091""
nohup sshpass ssh -o "StrictHostKeyChecking no"   -i ${cert} -n -f ${client2} "${cli_export} && ${export_command} && .${Rabia_Path}" >${output_path}22.log &

cli_export="export RC_Role=cli RC_Index=2 RC_Proxy="${replica3_ip}:9092""
nohup sshpass ssh -o "StrictHostKeyChecking no"   -i ${cert} -n -f ${client3} "${cli_export} && ${export_command} && .${Rabia_Path}" >${output_path}23.log &

cli_export="export RC_Role=cli RC_Index=3 RC_Proxy="${replica4_ip}:9093""
nohup sshpass ssh -o "StrictHostKeyChecking no"   -i ${cert} -n -f ${client4} "${cli_export} && ${export_command} && .${Rabia_Path}" >${output_path}24.log &

cli_export="export RC_Role=cli RC_Index=4 RC_Proxy="${replica5_ip}:9094""
nohup sshpass ssh -o "StrictHostKeyChecking no"   -i ${cert} -n -f ${client5} "${cli_export} && ${export_command} && .${Rabia_Path}" >${output_path}25.log &

echo "starting controller"

crl_export="export RC_Role=ctrl"
nohup sshpass ssh -o "StrictHostKeyChecking no"    -i ${cert} -n -f ${replica1} "${crl_export} && ${export_command} && .${Rabia_Path}" >${output_path}10.log &

sleep 10

if [[ "${isLeaderKill}" == "yes" ]]
then
  echo "killing the first leader"
  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "${replica1}" "${kill_command}; ${kill_command}"
fi

sleep 100

echo "Completed Client[s]"

if [[ "${collectClientLogs}" == "yes" ]]
then
  echo "collecting client logs"
  scp -i ${cert} ${client1}:${remote_log_path}0.txt ${output_path}21.txt
  scp -i ${cert} ${client2}:${remote_log_path}1.txt ${output_path}22.txt
  scp -i ${cert} ${client3}:${remote_log_path}2.txt ${output_path}23.txt
  scp -i ${cert} ${client4}:${remote_log_path}3.txt ${output_path}24.txt
  scp -i ${cert} ${client5}:${remote_log_path}4.txt ${output_path}25.txt
fi

echo "Finish test"