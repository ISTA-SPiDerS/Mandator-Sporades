# This test only tests the correctness of the committed logs
import sys
import csv

files = []
for i in range(1, len(sys.argv)):
    lines = []
    with open(sys.argv[i]) as file:
        lines = [line.rstrip() for line in file]

    blocks = []
    for l in lines:
        if l.startswith("Committed "):
            blocks.append(l.split(" ")[4])
    files.append(blocks)

    print("Number of committed blocks in " + sys.argv[i] + " is " + str(len(blocks)))

minLength = min(len(f) for f in files)

f = open('logs/committed-blocks.csv', 'w')
# create the csv writer
writer = csv.writer(f)
for i in range(0, minLength):
    # write a row to the csv file
    row = [files[0][i], files[1][i], files[2][i],files[0][i]==files[1][i] and files[1][i]==files[2][i]]
    if not(files[0][i]==files[1][i] and files[1][i]==files[2][i]):
        print("Error")
    writer.writerow(row)

# close the file
f.close()