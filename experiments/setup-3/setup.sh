pwd=$(pwd)
. "${pwd}"/experiments/setup-3/ip.sh

reset_directory="rm -r /home/${user_name}/mandator/ ; mkdir -p /home/${user_name}/mandator/logs"
local_binary_path="${pwd}/experiments/binary/"
replica_home_path="/home/${user_name}/mandator/"

# generate config files

python3 experiments/python/paxos-config.py               3 3 ${replica1_ip} ${replica2_ip} ${replica3_ip} ${client1_ip} ${client2_ip} ${client3_ip} > experiments/binary/paxos.yml
python3 experiments/python/mandator-sporades-config.py   3 3 ${replica1_ip} ${replica2_ip} ${replica3_ip} ${client1_ip} ${client2_ip} ${client3_ip} > experiments/binary/mandator-sporades.yml


for i in "${machines[@]}"
do
   echo "copying files to ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" "$i" -i ${cert} "${reset_directory}"
   scp -r -i ${cert} ${local_binary_path} "$i":${replica_home_path}
done


rm experiments/binary/mandator-sporades.yml
rm experiments/binary/paxos.yml

echo "3 replica setup complete"

