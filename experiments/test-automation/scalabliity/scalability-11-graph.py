import matplotlib.pyplot as plt

paxos_throughput=[1094, 5498, 10984, 22011, 55010, 87943, 109965, 137020, 133003, 136217, 146783]
paxos_latency=[229825, 242104, 265761, 304214, 292333, 358795, 506552, 700495, 1088632, 1404785, 1427270]
paxos_tail=[249807, 306005, 354030, 488747, 433436, 568829, 898362, 1111805, 1943239, 1988674, 2000000]

sporades_throughput=[1094, 5496, 10987, 22000, 54960, 88013, 110033, 133758, 144646, 156300]
sporades_latency=[228521, 237412, 267468, 321251, 318611, 344448, 520257, 789798, 951296, 1173938]
sporades_tail=[247933, 261800, 373146, 576130, 450121, 626182, 865983, 1503789, 1954491, 2144599]

mandator_throughput=[1090, 5476, 10968, 21974, 54914, 87916, 109912, 137309, 164782, 192321, 219625, 246483, 274630, 302106, 329589, 351551, 384216]
mandator_latency=[435452, 404399, 423264, 430276, 437869, 454701, 461286, 478741, 484679, 489584, 497830, 509988, 522415, 541618, 545293, 553093, 577111]
mandator_tail=[539102, 528321, 512260, 522554, 547933, 615209, 645090, 1070500, 1243355, 1302697, 1379265, 1878519, 1942122, 1951312, 3288964, 3023897, 3699493]

# cross checked with the csv sept 24 13.10

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
ax.set_xlim([0, 400])
ax.set_ylim([0, 1200])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi-Paxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL-RACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('Median Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/scalability-11/logs/wan_throughput_median.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 320])
ax.set_ylim([0, 2000])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi-Paxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_tail),       'ko-', label="SADL-RACS")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/scalability-11/logs/wan_throughput_tail.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()