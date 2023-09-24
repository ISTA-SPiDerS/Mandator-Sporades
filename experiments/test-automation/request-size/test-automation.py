import os

os.system("/bin/bash experiments/setup-5/setup.sh")

arrivals = [100, 200, 500,  1000, 2000, 5000, 10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000]

for iteration in [1, 2, 3]:
    for size in [64, 256]:
        for arrival in arrivals:
            # mandator
            os.system("/bin/bash experiments/request-size-bash/mandator.sh "+str(arrival)+" "+str(size)+" "+str(iteration))

            if size == 256 and arrival >10000:
                continue
            else:
                # paxos
                os.system("/bin/bash experiments/request-size-bash/paxos_raft.sh "+str(arrival)+" "+str(size)+" "+str(iteration))

                # sporades
                os.system("/bin/bash experiments/request-size-bash/sporades.sh "+str(arrival)+" "+str(size)+" "+str(iteration))
