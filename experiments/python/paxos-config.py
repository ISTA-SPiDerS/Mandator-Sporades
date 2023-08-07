import sys

if len(sys.argv) < 3:
    exit("not enough arguments")

numReplicas = int(sys.argv[1])
numClients = int(sys.argv[2])
replicaIPs = []
clientIPs = []

if len(sys.argv) < 3 + numReplicas + numClients:
    exit("not enough arguments")

argC = 3
for i in range(argC, argC + numReplicas, 1):
    replicaIPs.append(sys.argv[i])

argC = argC + numReplicas

for i in range(argC, argC + numClients, 1):
    clientIPs.append(sys.argv[i])


def print_peers(replicaIPs):
    print("peers:")
    for j in range(1, 1 + numReplicas, 1):
        print("   - name: " + str(j))
        print("     address: " + str(replicaIPs[j - 1])+str(":")+str(10000+j))
        print("     gaddress: " + str(replicaIPs[j - 1])+str(":")+str(11000+j))


def print_clients(clientIPs):
    print("clients:")
    for j in range(1, 1 + numClients, 1):
        print("   - name: " + str(j+20))
        print("     address: " + str(clientIPs[j - 1])+str(":")+str(12000+j))
        print("     gaddress: " + str(clientIPs[j - 1])+str(":")+str(13000+j))


print_peers(replicaIPs)
print()
print_clients(clientIPs)
