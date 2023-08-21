# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-54-252-64-134.ap-southeast-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-231-59-24.ap-northeast-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-13-125-106-14.ap-northeast-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-13-208-172-203.ap-northeast-3.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-255-226-190.ap-southeast-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-13-211-26-191.ap-southeast-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-52-199-134-27.ap-northeast-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-3-35-236-52.ap-northeast-2.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-15-168-13-110.ap-northeast-3.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-13-212-174-190.ap-southeast-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="54.252.64.134"
replica2_ip="13.231.59.24"
replica3_ip="13.125.106.14"
replica4_ip="13.208.172.203"
replica5_ip="54.255.226.190"

client1_ip="13.211.26.191"
client2_ip="52.199.134.27"
client3_ip="3.35.236.52"
client4_ip="15.168.13.110"
client5_ip="13.212.174.190	"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

