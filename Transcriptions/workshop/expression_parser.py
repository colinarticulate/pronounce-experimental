'''
Parsing utilities for grammar-like patterns
'''

import os
import time

import numpy as np


def combinations(mx):
    '''
        A way to count where each dimension might have different length.
        mx: an array with max lengths per dimension
        returns: all possible combinations 
    '''
    m = mx.prod()
    stack=[]
    k=1
    for i in range(len(mx)):
        a = np.arange( start=0, stop=mx[i], step=1)
        k=k*mx[i]
        
        b = np.repeat(np.tile(a,int(k/mx[i])),int(m/k))
        
        stack.append(b)

    result = np.vstack(stack).T
    
    return result


def creates_combinations_given_lengths():
    idx = np.array([0,1,2,3])
    mx = np.array([1,2,1,3])

    m=mx.prod()

    a=np.arange(start=0, stop=2, step=1)

    pattern = int(m/len(a))
    b=np.tile(a, pattern)
    c1=np.repeat(a,8)
    d0=np.repeat(np.tile(a,1),8)
    d1=np.repeat(np.tile(a,2),4)
    d2=np.repeat(np.tile(a,4),2)
    d3=np.repeat(np.tile(a,8),1)


    idx = np.array([0,1,2,3])
    mx2 = np.array([2,2,2,2])

    result2 = combinations(mx2)

    result = combinations(mx)


def join_list(alist,join_pattern):

    result=[]
    for item in alist[:-1]:
        if item!='':
            result.append([item.strip(" ")])
        result.append(join_pattern)
    if alist[-1]!='':
        result.append([alist[-1].strip(" ")])

    return tuple(result)


def main():

    transcription = [["K Y UW T IH K AX L"]]
    r1=["K axL", "K AX L", "K L"]

    for i,item in enumerate(transcription):
        if len(item)==1:
            for pattern in r1:
                if pattern in item[0]:
                    parts = item[0].split(pattern)
                    #tr = f"{pattern}".join(parts)
                    result = join_list(parts, r1)
                    transcription[i]= result


    print("working on it.")




    



if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
