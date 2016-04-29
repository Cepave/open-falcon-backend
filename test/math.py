#!/usr/bin/python 

import numpy as np

#{"www.google.com", "13.24", "38.90", "19.62", "9.48", "13.62"},
#{"www.yahoo.com", "6.72", "29.08", "8.55", "7.40", "-", "6.26"},


a=[13.24, 38.90, 19.62, 9.48, 13.62] 
data_a = np.array(a)
print "For array", a
print "median", np.median(data_a)
print "dev", np.std(data_a)
print "max", max(a)
print "min", min(a)
print "average", np.mean(data_a)


b=[6.72, 29.08, 8.55, 7.40, 6.26]
data_b = np.array(b) 
print "For array", b
print "median", np.median(data_b)
print "dev", np.std(data_b)
print "max", max(b)
print "min", min(b)
print "average", np.mean(data_b)
