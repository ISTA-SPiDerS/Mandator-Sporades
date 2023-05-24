#     1. Build the project
#     2. Copy client/ configuration/ experiments/ logs/ replica/ worker/ directories to logs/async/ directory
#     2. Zip logs/async/ and create async.zip
#     3. Copy asycn.zip to all machines

rm client/bin/client
rm replica/bin/replica
rm worker/bin/child

echo "Removed old binaries"
mage generate && mage build
echo "Built Async"

output_path="logs/"

rm -r ${output_path}; mkdir ${output_path} 

reset_directory="rm -r /home/pasindu/async; mkdir /home/pasindu/async"
kill_insstances="pkill replica ; pkill client; pkill child"
unzip_async="cd /home/pasindu/async && unzip async.zip"

replica1=pasindu@dedis-140.icsil1.epfl.ch
replica1_cert="/home/pasindu/Pictures/pasindu_rsa"
replica2=pasindu@dedis-141.icsil1.epfl.ch
replica2_cert="/home/pasindu/Pictures/pasindu_rsa"
replica3=pasindu@dedis-142.icsil1.epfl.ch
replica3_cert="/home/pasindu/Pictures/pasindu_rsa"


client1=pasindu@dedis-145.icsil1.epfl.ch
client1_cert="/home/pasindu/Pictures/pasindu_rsa"
client2=pasindu@dedis-146.icsil1.epfl.ch
client2_cert="/home/pasindu/Pictures/pasindu_rsa"
client3=pasindu@dedis-147.icsil1.epfl.ch
client3_cert="/home/pasindu/Pictures/pasindu_rsa"

mkdir temp

cp -r logs/ temp/
cp -r client/ temp/
cp -r configuration/ temp/
cp -r experiments/ temp/
cp -r replica/ temp/
cp -r worker/ temp/

zip -r logs/async.zip temp/
rm -r temp/

local_zip_path="logs/async.zip"
replica_home_path="/home/pasindu/async/"

echo "Replica 1"
sshpass ssh ${replica1} -i ${replica1_cert} ${reset_directory}
sshpass ssh ${replica1} -i ${replica1_cert} ${kill_insstances}
scp -i ${replica1_cert} ${local_zip_path} ${replica1}:${replica_home_path}
sshpass ssh ${replica1} -i ${replica1_cert} ${unzip_async}

echo "Replica 2"
sshpass ssh ${replica2} -i ${replica2_cert} ${reset_directory}
sshpass ssh ${replica2} -i ${replica2_cert} ${kill_insstances}
scp -i ${replica2_cert} ${local_zip_path} ${replica2}:${replica_home_path}
sshpass ssh ${replica2} -i ${replica2_cert} ${unzip_async}

echo "Replica 3"
sshpass ssh ${replica3} -i ${replica3_cert} ${reset_directory}
sshpass ssh ${replica3} -i ${replica3_cert} ${kill_insstances}
scp -i ${replica3_cert} ${local_zip_path} ${replica3}:${replica_home_path}
sshpass ssh ${replica3} -i ${replica3_cert} ${unzip_async}

echo "Client 1"
sshpass ssh ${client1} -i ${client1_cert} ${reset_directory}
sshpass ssh ${client1} -i ${client1_cert} ${kill_insstances}
scp -i ${client1_cert} ${local_zip_path} ${client1}:${replica_home_path}
sshpass ssh ${client1} -i ${client1_cert} ${unzip_async}

echo "Client 2"
sshpass ssh ${client2} -i ${client2_cert} ${reset_directory}
sshpass ssh ${client2} -i ${client2_cert} ${kill_insstances}
scp -i ${client2_cert} ${local_zip_path} ${client2}:${replica_home_path}
sshpass ssh ${client2} -i ${client2_cert} ${unzip_async}

echo "Client 3"
sshpass ssh ${client3} -i ${client3_cert} ${reset_directory}
sshpass ssh ${client3} -i ${client3_cert} ${kill_insstances}
scp -i ${client3_cert} ${local_zip_path} ${client3}:${replica_home_path}
sshpass ssh ${client3} -i ${client3_cert} ${unzip_async}

rm ${local_zip_path}

echo "setup complete"