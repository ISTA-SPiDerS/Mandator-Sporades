import matplotlib.pyplot as plt
import numpy as np

# scalability summary
plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()


# set width of bar
barWidth = 0.20

# set height of bar
mandator = [110 * 3, 50 * 7, 30 * 11]

sporades = [90 * 3, 27 * 7, 12.5 * 11]

paxos = [100 * 3, 30 * 7,13 * 11]

# Set position of bar on X axis
br1 = np.arange(len(mandator))
br2 = [x + barWidth for x in br1]
br3 = [x + barWidth for x in br2]

# Make the plot
plt.bar(br1, mandator, color='k', width=barWidth, edgecolor='grey', label='SADL\nRACS')
plt.bar(br2, sporades, color='b', width=barWidth, edgecolor='grey', label='RACS')
plt.bar(br3, paxos,    color='g', width=barWidth, edgecolor='grey', label='Multi\nPaxos')
# Adding Xticks
plt.ylabel('Throughput (x1k cmd/sec)')
plt.xticks([r + barWidth for r in range(len(mandator))],
           ['N=3','N=7', "N=11"])

# plt.legend(ncol= 3,bbox_to_anchor =(1.10, 1.15),fontsize="9.30")
plt.legend()
ax.set_xlim([-0.5, 3.5])

plt.savefig('experiments/scalability/scalability.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()