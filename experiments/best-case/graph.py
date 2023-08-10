import matplotlib.pyplot as plt

# ninty latency

plt.figure(figsize=(15, 14))
# plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()


rabia_throughput=[102, 52, 119, 157, 182, 217]
rabia_latency=[2000000, 2000000, 2000000, 2000000, 2000000, 2000000]
rabia_tail=[2000000, 2000000, 2000000, 2000000, 2000000, 2000000]

paxos_throughput=[2485, 50006, 124968, 173720, 198444, 237306]
paxos_latency=[220864, 239465, 225486, 287760, 303336, 292446]
paxos_tail=[259133, 332395, 461070, 717855, 723246, 1626308]

epaxos_exec_throughput=[2515, 49795, 104533, 158625, 174866, 206779]
epaxos_exec_latency=[167731, 489013, 592759, 880419, 895231, 943564]
epaxos_exec_tail=[311939, 520000, 840000, 1643021, 1785452, 1884837]

epaxos_commit_throughput=[2516, 50118, 125059, 175095, 200104, 242175]
epaxos_commit_latency=[116021, 127235, 129506, 139080, 139325, 321572]
epaxos_commit_tail=[121179, 182326, 179452, 220374, 243053, 427266]

raft_throughput=[2505, 40539, 91689, 109605, 116568, 141334]
raft_latency=[258062, 569332, 615637, 704929, 737593, 739410]
raft_tail=[317447, 2000000, 2000000, 2000000, 2000000, 2000000]

mandator_paxos_throughput=[2493, 49945, 119878, 160896, 184670, 225381]
mandator_paxos_latency=[361856, 389905, 399112, 454121, 427513, 460564]
mandator_paxos_tail=[440470, 527506, 1133513, 1447991, 1500674, 1732333]

mandator_async_throughput=[2493, 49950, 123268, 171198, 195365, 234402]
mandator_async_latency=[365825, 352061, 376185, 389624, 393502, 406085]
mandator_async_tail=[444906, 439194, 850272, 942351, 1026126, 1629484]

sporades_throughput=[2495, 50016, 125099, 174902, 199655, 243456]
sporades_latency=[110914, 119553, 118578, 117037, 119988, 117931]
sporades_tail=[125718, 135193, 292255, 373543, 456020, 1605724]


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
plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_tail), 'y-', label="Mandator Paxos")
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
plt.plot(di_func(mandator_paxos_throughput), di_func(mandator_paxos_latency), 'y-', label="Mandator Paxos")
plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_latency), 'k-', label="Mandator Sporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_latency), 'b-', label="Sporades")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case/logs/throughput_median.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()