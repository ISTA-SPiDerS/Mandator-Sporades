# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-18-144-3-160.us-west-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-3-101-42-144.us-west-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-18-144-133-32.us-west-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-13-52-80-58.us-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-177-36-115.us-west-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-54-153-75-177.us-west-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-54-193-197-45.us-west-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-184-169-242-225.us-west-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-54-215-251-137.us-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-13-52-102-42.us-west-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="18.144.3.160"
replica2_ip="3.101.42.144"
replica3_ip="18.144.133.32"
replica4_ip="13.52.80.58"
replica5_ip="54.177.36.115"

client1_ip="54.153.75.177"
client2_ip="54.193.197.45"
client3_ip="184.169.242.225"
client4_ip="54.215.251.137"
client5_ip="13.52.102.42"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

