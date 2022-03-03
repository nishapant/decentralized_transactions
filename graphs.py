import matplotlib.pyplot as plt
import math 
import numpy as np

data1 = [1, 2, 3, 4, 5, 1, 4, 3, 6, 7, 8, 9, 10, 1, 5, 6, 1, 1, 2, 1, 4]
data2 = [1, 2, 3, 4, 5, 1, 4, 3, 6, 7, 8, 9, 10, 1]
data3 = []
data4 = []
data5 = []
data6 = []
data7 = []
data8 = []

total_data = data + data1

# CDF Plotting
data_sorted = np.sort(total_data)

# calculate the proportional values of samples
p = 1. * np.arange(len(total_data)) / (len(total_data) - 1)

# plot the sorted data:
fig = plt.figure()
# ax1 = fig.add_subplot(121)
# ax1.plot(p, data_sorted)
# ax1.set_xlabel('$p$')
# ax1.set_ylabel('$x$')

ax2 = fig.add_subplot()
ax2.plot(data_sorted, p)
ax2.set_xlabel('$x$')
ax2.set_ylabel('$p$')

plt.show()