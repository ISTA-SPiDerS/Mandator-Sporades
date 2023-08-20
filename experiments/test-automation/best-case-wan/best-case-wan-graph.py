import matplotlib.pyplot as plt
from datetime import datetime

# last updated 2023 august 20 14.52 with the final wan results

paxos_throughput=[2491, 50030, 100004, 150031, 174982, 199945, 225055, 240181, 227403]
paxos_latency=[207882, 232379, 237292, 226951, 246589, 244059, 299160, 300526, 658649]
paxos_tail=[240377, 317369, 502695, 574341, 614192, 661468, 926448, 2000000, 2000000]

epaxos_exec_throughput=[2515, 49555, 80284, 138909, 159857, 178875, 195866, 213970, 250939]
epaxos_exec_latency=[156685, 398825, 1145606, 953674, 844149, 797040, 741863, 699373, 630086]
epaxos_exec_tail=[309212, 956961, 980000, 1384757, 1564793, 1582607, 1634220, 1742918, 1882497]

epaxos_commit_throughput=[2516, 50119, 100080, 150088, 175049, 200093, 225160, 250165, 300158, 349600, 398170, 445707]
epaxos_commit_latency=[106480, 129374, 142055, 139056, 136025, 132331, 131238, 133650, 146201, 188382, 249271, 263934]
epaxos_commit_tail=[111936, 183753, 202666, 205696, 203625, 210963, 206268, 219069, 339982, 531478, 558393, 575160]

mandator_async_throughput=[2498, 49920, 99874, 149831, 174754, 199692, 224749, 246414, 287452, 328957, 366869, 412827]
mandator_async_latency=[431094, 442795, 433611, 447527, 481953, 491409, 464044, 479652, 530972, 515180, 573499, 559205]
mandator_async_tail=[521588, 545704, 596956, 1057418, 1211447, 1342458, 1377030, 1946132, 2000000, 2000000, 2000000, 2000000]

sporades_throughput=[2497, 50005, 99930, 149930, 175021, 199902, 224816, 226676]
sporades_latency=[204103, 232927, 243127, 233075, 333449, 355796, 399411, 1306364]
sporades_tail=[229151, 321681, 389958, 514558, 649526, 836412, 1360596, 2000000]

#  cross checked the above with the csv 2 times august 20 15.08

def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList

plt.figure(figsize=(6, 5))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 400])
ax.set_ylim([0, 1750])

# plt.plot(di_func(rabia_throughput), di_func(rabia_tail), 'b*-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_tail), 'g*-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_tail), 'r*-.', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_tail), 'c*-', label="Epaxos\ncommit")
# plt.plot(di_func(raft_throughput), di_func(raft_tail), 'm*-', label="Raft (no pipline)")
# plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_tail), 'y*-', label="Mandator Paxos")
plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_tail), 'k*-', label="Mandator\nSporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_tail), 'b*-', label="Sporades")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case-wan/logs/wan_throughput_tail_'+str(datetime.now())+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(6, 5))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 400])
ax.set_ylim([0, 600])

# plt.plot(di_func(rabia_throughput), di_func(rabia_latency), 'b*-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_latency), 'g*-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_latency), 'r*-.', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_latency), 'c*-', label="Epaxos\ncommit")
# plt.plot(di_func(raft_throughput), di_func(raft_latency), 'm*-', label="Raft (no pipeline)")
# plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_latency), 'y*-', label="Mandator Paxos")
plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_latency), 'k*-', label="Mandator\nSporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_latency), 'b*-', label="Sporades")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case-wan/logs/wan_throughput_median_'+str(datetime.now())+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()

# checked the above 2023 august 20 15.11