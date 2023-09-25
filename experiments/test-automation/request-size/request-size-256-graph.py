import matplotlib.pyplot as plt

paxos_throughput=[494, 992, 2492, 4995, 9997, 23120]
paxos_latency=[201327, 207917, 208717, 231765, 217076, 860773]
paxos_tail=[237793, 258612, 256400, 323927, 379371, 1945283]

sporades_throughput=[494, 994, 2507, 5001, 9994, 23790, 29613]
sporades_latency=[214990, 221126, 222672, 211953, 211086, 859755, 1839512]
sporades_tail=[259839, 273758, 259697, 252782, 294708, 1670834, 2431089]

mandator_throughput=[495, 1001, 2511, 4980, 9974, 25006, 49946]
mandator_latency=[438597, 408952, 411276, 441855, 414717, 450564, 452711]
mandator_tail=[541989, 496325, 498574, 540861, 626300, 1088832, 1931188]

# cross checked with the csv 2023 september 25 20.08

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
# ax.set_xlim([0, 650])
ax.set_ylim([0, 2000])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_tail),       'ko-', label="SADL\nRACS")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/request-size/logs/wan_throughput_tail-256.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
# ax.set_xlim([0, 700])
ax.set_ylim([0, 1000])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL\nRACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/request-size/logs/wan_throughput_median-256.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()