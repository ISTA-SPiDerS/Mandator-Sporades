import matplotlib.pyplot as plt
import numpy as np

# asynchrony summary
plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 13.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
ax = plt.gca()
ax.grid()


# set width of bar
barWidth = 0.20

# set height of bar
mandator = [226117/1000, 242148/1000]
sporades = [54630/1000, 57474/1000]
paxos = [18495/1000, 36456/1000]
raft = [21756/1000, 30870/1000]

# Set position of bar on X axis
br1 = np.arange(len(mandator))
br2 = [x + barWidth for x in br1]
br3 = [x + barWidth for x in br2]
br4 = [x + barWidth for x in br3]

# Make the plot
plt.bar(br1, mandator, color='k', width=barWidth, edgecolor='grey', label='SADL\nRACS')
plt.bar(br2, sporades, color='b', width=barWidth, edgecolor='grey', label='RACS')
plt.bar(br3, paxos,    color='g', width=barWidth, edgecolor='grey', label='Multi\nPaxos')
plt.bar(br4, raft,     color='m', width=barWidth, edgecolor='grey', label='Raft')
# Adding Xticks
plt.ylabel('Throughput (x1k cmd/sec)')
plt.xticks([r + barWidth for r in range(len(mandator))],
           ['Epoch=500','Epoch=2000'])

# plt.legend(ncol= 3,bbox_to_anchor =(1.10, 1.15),fontsize="9.30")
plt.legend()
# ax.set_ylim([0, 750])

plt.savefig('experiments/asynchrony/logs/asynchrony.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()