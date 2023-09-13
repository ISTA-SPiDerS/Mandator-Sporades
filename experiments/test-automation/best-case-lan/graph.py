import matplotlib.pyplot as plt

# last updated 2023 September 13 14.51

paxos_throughput=[2504, 50015, 100011, 199995, 250030, 300031, 348066, 394746, 440371, 472538, 535065]
paxos_latency=[3321, 3719, 4234, 5100, 5294, 5698, 6226, 6589, 6918, 7378, 18889]
paxos_tail=[6670, 6588, 7598, 229320, 358267, 495844, 779152, 856529, 885923, 1122147, 1139849]

epaxos_exec_throughput=[2515, 50124, 100077, 200119, 250242, 300214, 350057, 400118, 450025, 499205, 587605, 701962, 795922]
epaxos_exec_latency=[3900, 5842, 5423, 5449, 5599, 5827, 5949, 6126, 6272, 6481, 6930, 7788, 8227]
epaxos_exec_tail=[10768, 11997, 12087, 38354, 83671, 334035, 309345, 360171, 428785, 497779, 730944, 932062, 875030]

epaxos_commit_throughput=[2516, 50125, 100093, 200123, 250258, 300218, 350126, 400059, 449979, 499707, 588312, 701316, 789652]
epaxos_commit_latency=[753, 2973, 3575, 3940, 4084, 4239, 4471, 4688, 4823, 5012, 5532, 6484, 7679]
epaxos_commit_tail=[1335, 6176, 6692, 24165, 26305, 224399, 427931, 335983, 414122, 475371, 673233, 809543, 813350]

raft_throughput=[2493, 50013, 100001, 199927, 244891, 289754, 332428, 376083, 416825, 453922, 472206]
raft_latency=[3425, 3843, 4278, 5112, 5555, 5930, 6213, 6518, 7411, 9995, 16155]
raft_tail=[6764, 6763, 18599, 409195, 1638669, 1484706, 1481646, 1447762, 1496754, 2000000, 2000000]

sporades_throughput=[2492, 50003, 100034, 200023, 247980, 295300, 343689, 385311, 430412, 458529]
sporades_latency=[3405, 4115, 4651, 5976, 6534, 6946, 7827, 9257, 9904, 47244]
sporades_tail=[6728, 7242, 11217, 316377, 1091804, 1137046, 969826, 1004098, 917585, 1442212]

rabia_throughput=[2516, 50124, 100050, 199891, 249928, 299602, 349987, 398472, 445729, 478877, 560958, 693295, 785540]
rabia_latency=[3409, 3498, 3658, 3720, 3904, 3784, 4081, 3895, 4023, 4268, 4823, 4606, 4696]
rabia_tail=[6366, 6598, 6865, 140505, 413509, 344466, 346905, 709583, 578927, 1222574, 1372473, 1335164, 1299787]

# cross checked with csv september 13 15.34

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
ax.set_xlim([0, 420])
ax.set_ylim([0, 100])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_tail),    'ro-', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_tail),  'c--', label="Epaxos\ncommit")
plt.plot(di_func(raft_throughput),             di_func(raft_tail),           'm*-', label="Raft")
plt.plot(di_func(rabia_throughput),            di_func(rabia_tail),          'y.-', label="Rabia")
# plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_tail), 'ko-', label="SADL-RACS")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc='lower right')
plt.savefig('experiments/best-case-lan/logs/lan_throughput_tail.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 800])
ax.set_ylim([0, 18])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi-Paxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_latency),    'ro-', label="Epaxos-exec")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_latency),  'c--', label="Epaxos-commit")
plt.plot(di_func(raft_throughput),             di_func(raft_latency),           'm*-', label="Raft")
plt.plot(di_func(rabia_throughput),            di_func(rabia_latency),          'y.-', label="Rabia")
# plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_latency), 'ko-', label="SADL-RACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/best-case-lan/logs/lan_throughput_median.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()

# 2023 september 13 15.37 verified