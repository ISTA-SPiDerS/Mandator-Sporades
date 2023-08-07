import sys

import numpy as np


def getPaxosRaftPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient + numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            sys.exit("Error in " + file_name + "\n")

        with f:
            content = f.readlines()
        if len(content) < 8:
            sys.exit("Error in " + file_name + "\n")

        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[4].strip().split(" ")[0] == "Throughput" and content[5].strip().split(" ")[0] == "Median"):
            sys.exit("Error in " + file_name + "\n")

        throughputs.append(float(content[4].strip().split(" ")[5]))
        medians.append(float(content[5].strip().split(" ")[3]))
        ninety9s.append(float(content[6].strip().split(" ")[4]))
        errors.append(float(content[7].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]


def getEPaxosPaxosPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient + numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            sys.exit("Error in " + file_name + "\n")

        with f:
            content = f.readlines()
        if len(content) < 5:
            sys.exit("Error in " + file_name + "\n")

        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[1].strip().split(" ")[0] == "Throughput" and content[2].strip().split(" ")[0] == "Median"):
            sys.exit("Error in " + file_name + "\n")

        throughputs.append(float(content[1].strip().split(" ")[5]))
        medians.append(float(content[2].strip().split(" ")[3]))
        ninety9s.append(float(content[3].strip().split(" ")[4]))
        errors.append(float(content[4].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]


def getRabiaPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient + numClients, 1)):
        file_name = root + str(cl) + ".log"

        try:
            f = open(file_name, 'r')
        except OSError:
            sys.exit("Error in " + file_name + "\n")

        with f:
            content = f.readlines()
        if len(content) < 5:
            sys.exit("Error in " + file_name + "\n")

        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[1].strip().split(" ")[0] == "Throughput" and content[2].strip().split(" ")[0] == "Median"):
            sys.exit("Error in " + file_name + "\n")

        throughputs.append(float(content[1].strip().split(" ")[2]))
        medians.append(float(content[2].strip().split(" ")[3]))
        ninety9s.append(float(content[3].strip().split(" ")[4]))
        errors.append(float(content[4].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]

def remove_farthest_from_median(A, n):
    median = np.median(A)

    # Calculate distances from the median
    distances = np.abs(A - median)

    # Sort the distances in descending order
    sorted_indices = np.argsort(distances)[::-1]

    # Remove n items farthest from the median
    cleaned_A = np.delete(A, sorted_indices[:n])

    return cleaned_A
