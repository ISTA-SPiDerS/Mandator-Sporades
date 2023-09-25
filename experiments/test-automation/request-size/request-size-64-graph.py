import matplotlib.pyplot as plt

paxos_throughput=[497, 995, 2501, 5002, 9995, 25009, 49978, 68479, 61246]
paxos_latency=[202678, 204167, 226555, 213327, 228835, 231051, 223540, 1415093, 1581400]
paxos_tail=[226484, 244639, 283384, 244612, 300354, 388200, 396293, 1988278, 2000000]

sporades_throughput=[502, 995, 2507, 4999, 9983, 25008, 50024, 74626, 102626, 108442, 123029]
sporades_latency=[195990, 201971, 222856, 213362, 229409, 211545, 245283, 1270686, 1453757, 1865038, 1857231]
sporades_tail=[213091, 236820, 254395, 246217, 297441, 281042, 398928, 1830995, 2394096, 2837560, 2576122]

mandator_throughput=[494, 994, 2501, 4981, 9997, 24975, 49946, 99908, 149829, 199729]
mandator_latency=[412421, 431311, 428915, 431582, 437724, 418398, 448879, 441643, 465428, 523336]
mandator_tail=[508389, 524524, 517811, 523045, 530635, 522533, 676794, 1610982, 1722198, 2199235]

# cross checked with the csv 2023 september 25 19.56

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
ax.set_ylim([0, 1989])

plt.plot(di_func(sporades_throughput),         di_func(sporades_tail),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_tail),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_tail),       'ko-', label="SADL\nRACS")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/request-size/logs/wan_throughput_tail-64.pdf', bbox_inches='tight', pad_inches=0)
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
ax.set_ylim([0, 750])

plt.plot(di_func(sporades_throughput),         di_func(sporades_latency),       'b*-', label="RACS")
plt.plot(di_func(paxos_throughput),            di_func(paxos_latency),          'g.-', label="Multi\nPaxos")
plt.plot(di_func(mandator_throughput),         di_func(mandator_latency),       'ko-', label="SADL\nRACS")


plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('median Latency (ms)')
plt.legend(fancybox=True, framealpha=0)
plt.savefig('experiments/request-size/logs/wan_throughput_median-64.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()