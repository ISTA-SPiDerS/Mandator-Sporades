import matplotlib.pyplot as plt

paxos_throughput=[299, 1500, 3007, 5996, 15003, 30019, 59983, 119967, 179993, 238174, 288670, 235878, 233859, 219918]
paxos_latency=[203552, 190979, 207502, 207367, 209578, 197911, 199306, 210431, 215379, 247428, 268607, 462822, 536568, 616380]
paxos_tail=[229827, 225524, 249660, 255564, 273425, 256656, 297417, 398225, 429641, 680611, 843074, 1030870, 1058185, 1114122]

epaxos_exec_throughput=[304, 1509, 3019, 6026, 15056, 30070, 55189, 110869, 158158, 198572, 244941,270431, 316565, 336257, 341590, 375737]
epaxos_exec_latency=[178636, 191860, 216862, 231590, 247960, 308951, 587186, 605119, 746773, 862534, 1004385, 1391329, 1454957, 1661759, 1867526, 2042793]
epaxos_exec_tail=[251885, 241054, 266016, 299107, 443709, 762516, 800000, 1435359, 1547178, 1880677, 1969020, 2058342, 2194049, 2282642, 2338791, 2447572]

epaxos_commit_throughput=[305, 1510, 3020, 6025, 15062, 30073, 60049, 120070, 179989, 238320, 296301, 350389, 408853, 446380, 482180, 531516]
epaxos_commit_latency=[57309, 118685, 123275, 114320, 132898, 130847, 140348, 135787, 147019, 138854, 155177, 176342, 184192, 208789, 203590, 250340]
epaxos_commit_tail=[82386, 118324, 128748, 121548, 181626, 183925, 191394, 200439, 349603, 421369, 442835, 528394, 544540, 564143, 673660, 684040]

sporades_throughput=[298, 1498, 2998, 6006, 15006, 29997, 59967, 120007, 179330, 233078, 234069, 234057, 247228, 264894, 274155, 276221]
sporades_latency=[183340, 188159, 205811, 191083, 195550, 192851, 199870, 209506, 258040, 267574, 354412, 522336, 591942, 668576, 671644, 682646]
sporades_tail=[201338, 215575, 238566, 248024, 257830, 274553, 319089, 526959, 738807, 1359847, 1359948, 1448310, 1464329, 1499675, 1599209, 1650951]

mandator_throughput=[296, 1503, 2983, 5992, 14982, 29998, 59952, 119954, 179895, 239249, 299534, 359543, 447242]
mandator_latency=[405445, 397581, 423952, 427381, 426554, 433591, 414449, 462796, 449473, 470922, 530500, 561203, 555327]
mandator_tail=[501758, 481837, 514335, 518779, 516671, 533096, 571572, 987157, 1176432, 1677324, 1724752, 2227140, 2713965]

# cross checked with the csv sept 24 22.59

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
ax.set_xlim([0, 700])
ax.set_ylim([0, 650])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_latency),    'ro-', label="Epaxos\nexec")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL\nRACS")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_latency),  'c--', label="Epaxos\ncommit")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc='lower right')
plt.savefig('experiments/scalability-3/logs/wan_throughput_median-3.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()


plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()
# ax.set_xlim([0, 650])
ax.set_ylim([0, 1250])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(epaxos_exec_throughput),      di_func(epaxos_exec_tail),    'ro-', label="Epaxos\nexec")
plt.plot(di_func(mandator_throughput),         di_func(mandator_tail),       'ko-', label="SADL\nRACS")
plt.plot(di_func(epaxos_commit_throughput),    di_func(epaxos_commit_tail),  'c--', label="Epaxos\ncommit")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0, loc='center right')
plt.savefig('experiments/scalability-3/logs/wan_throughput_tail-3.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()