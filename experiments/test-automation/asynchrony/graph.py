import matplotlib.pyplot as plt

paxos_throughput=[53, 283, 442, 1000, 2879]
paxos_latency=[141112, 190001, 400311, 600154, 2000000]

raft_throughput=[48, 156, 322, 1070, 2466]
raft_latency=[230892, 450982, 659032, 981203, 2000000]

sporades_throughput=[463, 1658, 1897, 8288, 11378, 11121, 12103, 15819, 23845, 24975, 25829, 28100]
sporades_latency=[163357, 180703, 193014, 267203, 279567, 370970, 430280, 441208, 661075, 666952, 998613, 1982044]

mandator_throughput=[492, 2487, 4959, 24890, 37242, 49671, 62055, 74388, 99177, 122937, 143258, 196510]
mandator_latency=[1218977, 1232065, 1370635, 1434295, 1544424, 1606062, 1619296, 1650134, 1771920, 1963338, 2135589, 2176353]

# cross checked with the csv 2023 september 20 12.13

def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
# ax.set_xlim([0, 700])
# ax.set_ylim([0, 800])
# plt.xscale("log")

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi-Paxos")
plt.plot(di_func(raft_throughput),             di_func(raft_latency),           'm*-', label="Raft")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL\nRACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('Median Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc="lower right")
plt.savefig('experiments/asynchrony/logs/asynchrony.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()