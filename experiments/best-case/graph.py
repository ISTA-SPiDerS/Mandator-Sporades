import matplotlib.pyplot as plt

# ninty latency

plt.figure(figsize=(15, 14))
# plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()


rabia_throughput=[109, 59, 116, 173, 156, 192, 230, 235, 256]
rabia_latency=[2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000]
rabia_tail=[2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000]

paxos_throughput=[2496, 49972, 124904, 174481, 199046, 239684, 272594, 290303, 288709]
paxos_latency=[207587, 218421, 246655, 235635, 259089, 292495, 404121, 592633, 711625]
paxos_tail=[238793, 295619, 485151, 595120, 625348, 1356380, 1561161, 1752709, 1999799]

epaxos_exec_throughput=[2515, 49650, 106973, 156397, 177068, 210729, 237084, 272212, 311803]
epaxos_exec_latency=[160797, 426939, 647134, 846429, 973362, 929382, 988007, 955007, 974429]
epaxos_exec_tail=[273655, 1950410, 2000000, 1642412, 1490445, 1275933, 1199610, 1100752, 1040970]

epaxos_commit_throughput=[2516, 50118, 125032, 175079, 200064, 242292, 285092, 330024, 401299]
epaxos_commit_latency=[110476, 131083, 133676, 132880, 134635, 318661, 425677, 433331, 444617]
epaxos_commit_tail=[116170, 188710, 184195, 207824, 225371, 424121, 548796, 592643, 733799]

raft_throughput=[2498, 40265, 93399, 105691, 119207, 128778, 149730, 180128, 190477]
raft_latency=[277673, 565866, 603433, 729443, 715104, 784782, 793317, 779202, 824111]
raft_tail=[343238, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000, 2000000]

mandator_paxos_throughput=[2497, 49921, 119874, 164931, 185586, 226556, 240850, 249202, 282185]
mandator_paxos_latency=[361667, 388580, 420192, 412452, 461984, 448021, 486718, 553603, 565564]
mandator_paxos_tail=[439592, 519058, 1099325, 1359606, 1427399, 1834772, 1999293, 2000000, 2000000]

mandator_async_throughput=[2489, 49947, 123267, 171438, 194321, 234402, 270278, 323577, 391299]
mandator_async_latency=[337109, 374795, 370951, 382836, 419436, 439415, 453381, 472524, 466547]
mandator_async_tail=[408835, 466267, 817876, 892314, 978684, 1458373, 1582406, 1701333, 1786463]

sporades_throughput=[2498, 50012, 124729, 174114, 197471, 244535, 271476, 292130, 304695]
sporades_latency=[208603, 242525, 271662, 277005, 364665, 337112, 452590, 611986, 682459]
sporades_tail=[230394, 334095, 518389, 682385, 750929, 822099, 1508541, 1908574, 1653197]


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


# plt.plot(di_func(rabia_throughput), di_func(rabia_tail), 'b-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_tail), 'g*-', label="Multi-Paxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_tail), 'r-.', label="Epaxos-exec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_tail), 'c-', label="Epaxos-commit")
plt.plot(di_func(raft_throughput), di_func(raft_tail), 'm-', label="Raft (no pipline)")
# plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_tail), 'y-', label="Mandator Paxos")
plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_tail), 'k-', label="Mandator Sporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_tail), 'b-', label="Sporades")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case/logs/throughput_tail.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


# plt.plot(di_func(rabia_throughput), di_func(rabia_latency), 'b-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_latency), 'g*-', label="Multi-Paxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_latency), 'r-.', label="Epaxos-exec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_latency), 'c-', label="Epaxos-commit")
plt.plot(di_func(raft_throughput), di_func(raft_latency), 'm-', label="Raft (no pipeline)")
# plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_latency), 'y-', label="Mandator Paxos")
plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_latency), 'k-', label="Mandator Sporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_latency), 'b-', label="Sporades")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case/logs/throughput_median.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()