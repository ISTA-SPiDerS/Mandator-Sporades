# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-3-25-175-180.ap-southeast-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-115-110-197.ap-northeast-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-3-39-224-81.ap-northeast-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-13-208-172-205.ap-northeast-3.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-255-248-3.ap-southeast-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-3-106-116-85.ap-southeast-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-43-207-67-32.ap-northeast-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-3-34-189-11.ap-northeast-2.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-15-152-39-234.ap-northeast-3.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-179-97-87.ap-southeast-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="3.25.175.180"
replica2_ip="13.115.110.197"
replica3_ip="3.39.224.81"
replica4_ip="13.208.172.205"
replica5_ip="54.255.248.3"

client1_ip="3.106.116.85"
client2_ip="43.207.67.32"
client3_ip="3.34.189.11"
client4_ip="15.152.39.234"
client5_ip="54.179.97.87"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

