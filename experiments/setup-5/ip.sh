# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-54-183-11-112.us-west-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-57-19-239.us-west-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-13-52-76-87.us-west-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-3-101-29-195.us-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-183-243-244.us-west-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-54-193-123-199.us-west-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-54-183-89-87.us-west-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-13-57-202-186.us-west-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-54-67-47-171.us-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-215-125-189.us-west-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="54.183.11.112"
replica2_ip="13.57.19.239"
replica3_ip="13.52.76.87"
replica4_ip="3.101.29.195"
replica5_ip="54.183.243.244"

client1_ip="54.193.123.199"
client2_ip="54.183.89.87"
client3_ip="13.57.202.186"
client4_ip="54.67.47.171"
client5_ip="54.215.125.189"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

