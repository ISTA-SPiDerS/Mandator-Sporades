import matplotlib.pyplot as plt

paxos_throughput=[2501, 50003, 150018, 174933, 199986, 224885, 241286, 228180]
paxos_latency=[221007, 236813, 239064, 294928, 301635, 310872, 320896, 652695]
paxos_tail=[257864, 322307, 571398, 627575, 674817, 703623, 900000, 1900000]

epaxos_exec_throughput=[2515, 49120, 136922, 158166, 176807, 194876, 214472, 252427]
epaxos_exec_latency=[153646, 449910, 539513, 645937, 716051, 735543, 882764, 916124]
epaxos_exec_tail=[293570, 493000, 1094807, 1334711, 1593062, 1731443, 1910471, 2116239]

epaxos_commit_throughput=[2516, 50095, 150081, 175092, 200100, 225130, 250178, 300129, 397689, 492254]
epaxos_commit_latency=[111386, 131614, 136820, 130989, 133854, 130932, 133190, 135651, 149602, 147134]
epaxos_commit_tail=[116595, 187432, 199996, 198526, 201351, 222266, 237983, 425684, 462100, 578212]

sporades_throughput=[2505, 50006, 150041, 174950, 199938, 224930, 249911, 299950]
sporades_latency=[218707, 230459, 253958, 265279, 303011, 369823, 412370, 1672533]
sporades_tail=[251876, 307322, 456644, 496294, 993407, 1380851, 2013499, 3456320]

mandator_throughput=[2488, 49955, 149878, 174744,199824, 224630, 249567, 299647, 399471, 492422]
mandator_latency=[429321, 420584, 449218, 453608, 459840, 466334, 483094, 496336, 536729, 790723]
mandator_tail=[519452, 518677, 1061630, 1165082, 1209596, 1268825, 1447247, 1902584, 2724231, 4834552]

# cross checked with the csv 2023 september 13 16.59

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
ax.set_xlim([0, 500])
ax.set_ylim([0, 3000])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_tail),    'ro-', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_tail),  'c--', label="Epaxos\ncommit")
# plt.plot(di_func(raft_throughput),             di_func(raft_tail),           'm*-', label="Raft")
plt.plot(di_func(mandator_throughput),         di_func(mandator_tail),       'ko-', label="SADL\nRACS")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc='center right')
plt.savefig('experiments/best-case-wan/logs/wan_throughput_tail.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 700])
ax.set_ylim([0, 800])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_latency),    'ro-', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_latency),  'c--', label="Epaxos\ncommit")
# plt.plot(di_func(raft_throughput),             di_func(raft_latency),           'm*-', label="Raft")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL\nRACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc="lower right")
plt.savefig('experiments/best-case-wan/logs/wan_throughput_median.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()