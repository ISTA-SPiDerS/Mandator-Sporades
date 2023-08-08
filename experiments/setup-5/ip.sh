# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-13-54-113-84.ap-southeast-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-3-85-119-40.compute-1.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-34-245-58-114.eu-west-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-43-204-237-226.ap-south-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-255-250-143.ap-southeast-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-13-210-35-9.ap-southeast-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-18-212-62-170.compute-1.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-34-247-159-191.eu-west-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-3-111-57-163.ap-south-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-18-143-140-74.ap-southeast-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="13.54.113.84"
replica2_ip="3.85.119.40"
replica3_ip="34.245.58.114"
replica4_ip="43.204.237.226"
replica5_ip="54.255.250.143"

client1_ip="13.210.35.9"
client2_ip="18.212.62.170"
client3_ip="34.247.159.191"
client4_ip="3.111.57.163"
client5_ip="18.143.140.74"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"