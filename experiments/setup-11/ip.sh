# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023
user_name="ubuntu"

replica1_name=ec2-3-27-149-75.ap-southeast-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-43-207-113-180.ap-northeast-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-13-125-5-113.ap-northeast-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-13-208-176-141.ap-northeast-3.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-179-124-230.ap-southeast-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

replica6_name=ec2-54-252-196-189.ap-southeast-2.compute.amazonaws.com
replica6=${user_name}@${replica6_name}

replica7_name=ec2-54-95-0-176.ap-northeast-1.compute.amazonaws.com
replica7=${user_name}@${replica7_name}

replica8_name=ec2-3-35-217-93.ap-northeast-2.compute.amazonaws.com
replica8=${user_name}@${replica8_name}

replica9_name=ec2-13-208-242-126.ap-northeast-3.compute.amazonaws.com
replica9=${user_name}@${replica9_name}

replica10_name=ec2-13-231-178-2.ap-northeast-1.compute.amazonaws.com
replica10=${user_name}@${replica10_name}

replica11_name=ec2-3-36-87-60.ap-northeast-2.compute.amazonaws.com
replica11=${user_name}@${replica11_name}

client1_name=ec2-13-238-120-105.ap-southeast-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-18-183-212-228.ap-northeast-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-3-36-119-2.ap-northeast-2.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-15-152-49-72.ap-northeast-3.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-255-213-54.ap-southeast-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

client6_name=ec2-54-206-96-204.ap-southeast-2.compute.amazonaws.com
client6=${user_name}@${client6_name}

client7_name=ec2-13-231-109-193.ap-northeast-1.compute.amazonaws.com
client7=${user_name}@${client7_name}

client8_name=ec2-52-78-193-110.ap-northeast-2.compute.amazonaws.com
client8=${user_name}@${client8_name}

client9_name=ec2-15-168-8-47.ap-northeast-3.compute.amazonaws.com
client9=${user_name}@${client9_name}

client10_name=ec2-43-207-143-219.ap-northeast-1.compute.amazonaws.com
client10=${user_name}@${client10_name}

client11_name=ec2-43-201-62-173.ap-northeast-2.compute.amazonaws.com
client11=${user_name}@${client11_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${replica6} ${replica7} ${replica8} ${replica9} ${replica10} ${replica11} ${client1} ${client2} ${client3} ${client4} ${client5} ${client6} ${client7} ${client8} ${client9} ${client10} ${client11})
echo "ip addresses loaded"

replica1_ip="3.27.149.75"
replica2_ip="43.207.113.180"
replica3_ip="13.125.5.113"
replica4_ip="13.208.176.141"
replica5_ip="54.179.124.230"
replica6_ip="54.252.196.189"
replica7_ip="54.95.0.176"
replica8_ip="3.35.217.93"
replica9_ip="13.208.242.126"
replica10_ip="13.231.178.2"
replica11_ip="3.36.87.60"

client1_ip="13.238.120.105"
client2_ip="18.183.212.228"
client3_ip="3.36.119.2"
client4_ip="15.152.49.72"
client5_ip="54.255.213.54"
client6_ip="54.206.96.204"
client7_ip="13.231.109.193"
client8_ip="52.78.193.110"
client9_ip="15.168.8.47"
client10_ip="43.207.143.219"
client11_ip="43.201.62.173"

kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill man_client; pkill man_replica;  pkill pa_ra_replica ; pkill pa_ra_client; pkill pipe_client; pkill pipe_replica; pkill rabia"
remote_log_path="/home/${user_name}/mandator/logs/"
reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}; ${kill_command}"
done

