#!/usr/bin/python

import time
import sys
import editdistance

INS = 1
DEL = 2
MATCH = 4
SUB = 3
undef = '-'

def initialize(ref_words, hyp_words):
    align_matrix = [[0 for i in range(hyp_words + 1)] for j in range(ref_words + 1)]
    backtrace_matrix = [[0 for i in range(hyp_words + 1)] for j in range(ref_words + 1)]

    for j in range(0, hyp_words + 1):
        align_matrix[0][j] = j
        backtrace_matrix[0][j] = INS

    for i in range(0, ref_words + 1):
        align_matrix[i][0] = i
        backtrace_matrix[i][0] = DEL

    return align_matrix, backtrace_matrix

def align(refs, hyps):
    ref_words = len(refs)
    hyp_words = len(hyps)
    align_matrix, backtrace_matrix = initialize(ref_words, hyp_words)

    for i in range(1, ref_words + 1):
        for j in range (1, hyp_words + 1):
            
            if refs[i - 1] != hyps[j - 1]:
                cost = 1
            else:
                cost = 0

            ins = align_matrix[i][j - 1] + 1
            dels = align_matrix[i - 1][j] + 1
            substs = align_matrix[i - 1][j - 1] + cost
            m = min(ins, dels, substs)
            align_matrix[i][j] = m

            if m == substs:
                backtrace_matrix[i][j] = MATCH + cost
            elif m == ins:
                backtrace_matrix[i][j] = INS
            elif m == dels:
                backtrace_matrix[i][j] = DEL

    backtrace(refs, hyps, align_matrix, backtrace_matrix)
    return
    
def backtrace(refs, hyps, align_matrix, backtrace_matrix):
    
    alignment = []
        
    inspen, delpen, substpen, match = (0, 0, 0, 0)

    i = len(refs)
    j = len(hyps)

    while not (i == 0 and j == 0):
        pointer = backtrace_matrix[i][j]
        if pointer == INS:
            alignment.append((undef, hyps[j - 1].upper()))
            inspen = inspen + 1
            j = j - 1
        elif pointer == DEL:
            alignment.append((refs[i - 1].upper(), undef))
            delpen = delpen + 1
            i = i - 1
        elif pointer == MATCH:
            alignment.append((refs[i - 1].lower(), hyps[j - 1].lower()))
            match = match + 1
            j = j - 1
            i = i - 1
        else:
            alignment.append((refs[i - 1].upper(), hyps[j - 1].upper()))
            substpen = substpen + 1
            j = j - 1
            i = i - 1

    print("INS", inspen)
    print("DEL", delpen)
    print("SUBST", substpen)
    print("TOTAL", len(refs))
    print("WER", float(inspen + delpen + substpen) / len(refs) * 100)

def extract_hypothesis(file):
    with open(file,  'r') as f:
        raw=f.read()

    lines = raw.strip("\n").split("\n")

    hypothesis=[]
    for line in lines:
        hyp=line.split(" (")[0].split(" ")
        hypothesis.append(hyp)

    return hypothesis

def extract_words(file):
    with open(file,  'r') as f:
        raw=f.read()

    lines = raw.strip("\n").split("\n")

    words=[]
    for line in lines:
        word=line.split(" <sil> </s>")[0].split("<s> <sil> ")[1].split(" ")
        words.append(["SIL"]+word+["SIL"])

    return words

def main():
    hyp_file="/home/dbarbera/Repositories/art_db/__workshop/tests_art_db_1-to-1.match"
    #hyp_file="/home/dbarbera/Repositories/art_db/__workshop/tests_en-us-us_1-to-1.match"
    ref_file="/home/dbarbera/Repositories/art_db/__workshop/art_db_test_phonemes.transcription"

    # fn1 = open(sys.argv[1], "r")
    # fn2 = open(sys.argv[2], "r")

    # fn1 = open(file1, "r")
    # fn2 = open(file2, "r")

    # words1 = " ".join(fn1.readlines()).split()
    # words2 = " ".join(fn2.readlines()).split()

    words1 = extract_hypothesis(hyp_file)
    words2 = extract_words(ref_file)

    #align(words1, words2)
    count=0
    relative_error=0
    for i, (word1, word2) in enumerate(zip(words1,words2)):
        ed=editdistance.eval(word1,word2)
        
        w1=" ".join(word1)
        w2=" ".join(word2)
        print(f"\t{w1}")
        print(f"\t{w2}")
        print(f"{i} {ed} {ed/len(word2)*100:.2f}%\n")
        count=count+ed
        relative_error=relative_error+ed/len(word2)

    print(count, count/len(words1)*100)
    print("PER: ", relative_error/len(words2)*100)
 


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
