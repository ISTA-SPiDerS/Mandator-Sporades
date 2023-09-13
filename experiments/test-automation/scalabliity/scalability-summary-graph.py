import matplotlib.pyplot as plt
import numpy as np

# scalability

# set width of bar
barWidth = 0.20
plt.figure(figsize=(5, 3))
plt.rcParams.update({'font.size': 9.30})
plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42

# set height of bar
quepaxa = [584935/1000, 582105/1000, 572427/1000, 571279/1000, 517638/1000, 467337/1000]
paxos = [434919/1000, 414942/1000, 411128/1000, 404734/1000, 385958/1000, 380000/1000]
epaxos = [614732/1000, 674802/1000, 0,0,0,0]

# Set position of bar on X axis
br1 = np.arange(len(quepaxa))
br2 = [x + barWidth for x in br1]
br3 = [x + barWidth for x in br2]

# Make the plot
plt.bar(br1, quepaxa, color='b', width=barWidth, edgecolor='grey', label='QuePaxa')
plt.bar(br2, paxos, color='c', width=barWidth, edgecolor='grey', label='Multi-Paxos')
plt.bar(br3, epaxos, color='m', width=barWidth, edgecolor='grey', label='EPaxos-commit')
# Adding Xticks
plt.ylabel('Throughput (x1k cmd/sec)')
plt.xticks([r + barWidth for r in range(len(quepaxa))],
           ['N=3','N=5', 'N=7', 'N=9', 'N=11', 'N=13'])

# plt.legend(ncol= 3,bbox_to_anchor =(1.10, 1.15),fontsize="9.30")
plt.legend()
ax = plt.gca()
ax.grid()
ax.set_ylim([0, 750])

plt.savefig('scalability.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()
