import numpy as np
import scipy.stats as st
import matplotlib.pyplot as plt


num = 100
lambda1 = num * 0.04 
#n = np.arange(,10)
n = [5,10,15,20,25,30,35,40,45,50,55,60,65,70,75,80,85,90,95,100]

print n

#y = st.poisson.pmf(n, lambda1)

y = [30,31,32,33,34,35,36,38,40,42,44,46,49,53,59,66,76,93,133,330]
y2 = [37,40,43,46,48,51,54,57,61,65,70,75,81,89,98,110,128,154,200,330]
y3 = [32,33,34,35,36,37,38,39,40,42,43,45,47,49,51,55,60,67,81,306]

ay = [55,55,55,55,55,55,55,55,55,55,55,55,55,55,55,55,55,55,55,55]
ay2 = [83,83,83,83,83,83,83,83,83,83,83,83,83,83,83,83,83,83,83,83]
ay3 = [46,46,46,46,46,46,46,46,46,46,46,46,46,46,46,46,46,46,46,46]
print y




fig, ax = plt.subplots()
ax.plot(n, y, ls = 'dashed',lw=2, c='r', label='1 copy')
ax.plot(n, y2, ls = 'dashed',lw=2, c='g', label='3 copy')
ax.plot(n, y3, ls = 'dashed',lw=2, c='b', label='F/2 +1 = 2')
ax.plot(n, ay, lw=2, c='r', label='1 copy avg')
ax.plot(n, ay2, lw=2, c='g', label=' 3 copy avg')
ax.plot(n, ay3, lw=2, c='b', label='F/2 +1 = 2 avg')
#ax.plot(n, y, ls='dashed', lw=2, c='r', 'k--', label='simulate zipf dist1')
#ax.plot(n, y2, ls='dashed', lw=2, c='g', 'k:', label='simulate zipf dist2')
#ax.plot(n, y3, ls='dashed', lw=2, c='b', 'k', label='simulate zipf dist3')

# The frame is matplotlib.patches.Rectangle instance surrounding the legend.

legend = ax.legend(loc='upper center', shadow=True)

plt.ylabel('# simulate value ')
plt.xlabel('# distribution percent')

plt.xlim([0, 100])       
plt.ylim([0, 300])
plt.grid(True)
plt.title('Zipf Distribution Analysis')


plt.show()
#print rv.pmf(10)
#num_years = [4, 10, 7, 5, 4, 0, 0, 1]
#x = range(10)
#ax.bar(np.array(x)-.4, num_years, label='Observed instances')
#ax.plot(x, sum(num_years)*rv.pmf(x), ls='dashed', 
#ax.plot(x, rv.pmf(x)), ls='dashed', lw=2, c='r', label='Poisson distribution\n$(\lambda=4.0)$')

#ax.xlim([1, 10])       
#ax.ylim([0.0, 1])
#ax.xlabel('# disk fail in 1hour')
#ax.ylabel('# rate')
#ax.legend(loc='best')
#ax.show()

