import numpy as np
from tqdm import tqdm

interm = np.loadtxt("./interm.csv", delimiter=',')
waves = np.loadtxt("./wave.csv", delimiter=',')

wavesT = waves.T
intermT = interm.T

table = []
for k_idx in tqdm(range(256)):
    l = []
    for t in range(len(waves[0])):
        l.append(abs(np.corrcoef(wavesT[t], intermT[k_idx])[0][1]))
    table.append(l)

np.savetxt("./py-output.csv", table, delimiter=',', fmt='%.8f')
