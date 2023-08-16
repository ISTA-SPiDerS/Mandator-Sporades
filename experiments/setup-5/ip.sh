# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-54-236-43-18.compute-1.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-34-207-211-15.compute-1.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-100-25-150-20.compute-1.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-54-196-44-135.compute-1.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-52-4-250-212.compute-1.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-18-206-206-192.compute-1.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-54-196-47-233.compute-1.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-34-239-253-172.compute-1.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-100-25-35-14.compute-1.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-157-148-240.compute-1.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="54.236.43.18"
replica2_ip="34.207.211.15"
replica3_ip="100.25.150.20"
replica4_ip="54.196.44.135"
replica5_ip="52.4.250.212"

client1_ip="18.206.206.192"
client2_ip="54.196.47.233"
client3_ip="34.239.253.172"
client4_ip="100.25.35.14"
client5_ip="54.157.148.240"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

