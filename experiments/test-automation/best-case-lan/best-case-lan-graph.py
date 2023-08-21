import matplotlib.pyplot as plt
from datetime import datetime

# last updated 2023 august 21 11.12 with the final wan results

paxos_throughput=[2495, 24991, 50011, 74946, 100045, 124987, 150050, 199855, 227984, 228461, 228978]
paxos_latency=[3373, 3609, 3884, 4172, 4437, 4712, 4986, 5443, 5965, 6795, 31448]
paxos_tail=[6684, 6408, 6788, 7727, 342138, 534730, 379815, 784584, 2000000, 2000000, 2000000]

epaxos_exec_throughput=[2515, 25102, 50123, 75117, 100091, 125069, 150132, 200122, 250235, 300204, 350010, 399818, 449024, 497646]
epaxos_exec_latency=[3879, 5711, 5491, 5881, 5706, 5693, 5731, 6161, 6240, 6404, 6811, 7111, 7300, 7595]
epaxos_exec_tail=[12422, 12082, 12111, 12192, 12266, 13599, 18755, 58619, 78779, 165278, 301331, 582026, 517215, 526411]

epaxos_commit_throughput=[2516, 25104, 50126, 75113, 100084, 125078, 150149, 200123, 250235, 300239, 350039, 399681, 449199, 496634]
epaxos_commit_latency=[753, 713, 3225, 3688, 3798, 3964, 4081, 4360, 4491, 4750, 5043, 5339, 5651, 5925]
epaxos_commit_tail=[1964, 3199, 6421, 7084, 7581, 13184, 17035, 27304, 73757, 202131, 235871, 328206, 507771, 367199]

raft_throughput=[2504, 25003, 49994, 74970, 99948, 124998, 149939, 197939, 240606, 274239, 310047, 350601]
raft_latency=[3504, 3714, 4026, 4378, 4578, 4976, 5194, 5525, 6069, 7116, 9773, 14239]
raft_tail=[6932, 6570, 7191, 16214, 29979, 159800, 280311, 1264435, 1751057, 2000000, 2000000, 2000000]

sporades_throughput=[2497, 25021, 50010, 74996, 99995, 124982, 150002, 199737, 247503, 283397, 304189]
sporades_latency=[3451, 3772, 4110, 4571, 4865, 5332, 5545, 6186, 7264, 9304, 576850]
sporades_tail=[6753, 6625, 7242, 176783, 248560, 267550, 419300, 998518, 1153069, 1354980, 1686178]

rabia_throughput=[2516, 25108, 50125, 75120, 100050, 125025, 150031, 199890, 247273, 285206]
rabia_latency=[3512, 3771, 3949, 3959, 4141, 4319, 4269, 4936, 5413, 182651]
rabia_tail=[6570, 6763, 7821, 15438, 46332, 97846, 263192, 415733, 1049593, 1715951]

#  cross checked the above with the csv 2 times august 21 11.42

def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList

plt.figure(figsize=(6, 5))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()
# ax.set_xlim([0, 400])
ax.set_ylim([0, 350])

plt.plot(di_func(rabia_throughput), di_func(rabia_tail), 'y*-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_tail), 'g*-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_tail), 'r*-.', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_tail), 'c*-', label="Epaxos\ncommit")
plt.plot(di_func(raft_throughput), di_func(raft_tail), 'm*-', label="Raft")
# plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_tail), 'k*-', label="Mandator\nSporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_tail), 'b*-', label="Sporades")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case-lan/logs/lan_throughput_tail_'+str(datetime.now())+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(6, 5))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 500])
ax.set_ylim([0, 10])

plt.plot(di_func(rabia_throughput), di_func(rabia_latency), 'y*-', label="Rabia")
plt.plot(di_func(paxos_throughput), di_func(paxos_latency), 'g*-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_latency), 'r*-.', label="Epaxos\nexec")
plt.plot(di_func(epaxos_commit_throughput), di_func(epaxos_commit_latency), 'c*-', label="Epaxos\ncommit")
plt.plot(di_func(raft_throughput), di_func(raft_latency), 'm*-', label="Raft")
# plt.plot(di_func(mandator_async_throughput), di_func(mandator_async_latency), 'k*-', label="Mandator\nSporades")
plt.plot(di_func(sporades_throughput), di_func(sporades_latency), 'b*-', label="Sporades")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(ncols=2)
plt.savefig('experiments/best-case-lan/logs/lan_throughput_median_'+str(datetime.now())+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()

# checked the above 2023 august 21 11.51: graphs look good