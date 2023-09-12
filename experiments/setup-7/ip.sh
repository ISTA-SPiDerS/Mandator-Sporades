# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-13-211-126-127.ap-southeast-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-54-199-52-226.ap-northeast-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-43-201-6-85.ap-northeast-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-13-208-185-9.ap-northeast-3.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-18-143-131-248.ap-southeast-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

replica6_name=ec2-13-238-128-160.ap-southeast-2.compute.amazonaws.com
replica6=${user_name}@${replica6_name}

replica7_name=ec2-43-207-171-135.ap-northeast-1.compute.amazonaws.com
replica7=${user_name}@${replica7_name}

client1_name=ec2-13-210-220-164.ap-southeast-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-13-112-28-246.ap-northeast-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-3-38-191-224.ap-northeast-2.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-15-168-39-2.ap-northeast-3.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-251-167-217.ap-southeast-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

client6_name=ec2-13-54-7-28.ap-southeast-2.compute.amazonaws.com
client6=${user_name}@${client6_name}

client7_name=ec2-52-192-70-166.ap-northeast-1.compute.amazonaws.com
client7=${user_name}@${client7_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${replica6} ${replica7} ${client1} ${client2} ${client3} ${client4} ${client5} ${client6} ${client7})
echo "ip addresses loaded"

replica1_ip="13.211.126.127"
replica2_ip="54.199.52.226"
replica3_ip="43.201.6.85"
replica4_ip="13.208.185.9"
replica5_ip="18.143.131.248"
replica6_ip="13.238.128.160"
replica7_ip="43.207.171.135"

client1_ip="13.210.220.164"
client2_ip="13.112.28.246"
client3_ip="3.38.191.224"
client4_ip="15.168.39.2"
client5_ip="54.251.167.217"
client6_ip="13.54.7.28"
client7_ip="52.192.70.166"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

