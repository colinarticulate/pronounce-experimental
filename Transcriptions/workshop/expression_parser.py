'''
Parsing utilities for grammar-like patterns
'''

import os
import time
import toml

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


def example_creating_combinations_given_lengths():
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

    return result

def remove_unidirectional(rule):
    for i,r in enumerate(rule):
        if ">>" in r:
            rule[i]=r.strip(">>")
    
    return rule

def parser(transcription_str, rules_set):

    transcription = [[transcription_str]]
    for rule in rules_set:
        for i,item in enumerate(transcription):
            for pattern in rule:
            
                if len(item)==1:

                    if pattern in item[0] and ">>" not in pattern:
                        parts = item[0].split(pattern)
                        uni_rule=remove_unidirectional(rule)
                        result = join_list(parts, uni_rule)
                        before=transcription[:i]
                        after = transcription[i+1:]
                        transcription = before+result+after
                        break
 
    return transcription


def parser_example():
    transcription = [["K Y UW T IH K AX L"]]
    r1=["K axL", "K AX L", "K L"]
    r2=["Y UW","Yuw"]
    rules_set=[r1,r2]

    valid_transcriptions = parser(transcription, rules_set)
    print("transcription: ",transcription)
    print("all valid transcriptions: ", valid_transcriptions)


def generate_valid_transcriptions( valid_possibilities ):
    lengths = []
    for item in valid_possibilities:
        lengths.append(len(item))

    mx=np.array(lengths)
    idxs = combinations(mx)

    m, n = idxs.shape

    all_transcriptions=[]
    for i in range(m):
        transcription=[]
        for j in range(n):
            phone = valid_possibilities[j][idxs[i,j]]
            transcription.append(phone)
        all_transcriptions.append(" ".join(transcription))

    return all_transcriptions

def example_combinations_of_valid_possibilities():
    valid_transcriptions=[['K'], ['Y UW', 'Yuw'], ['T IH'], ['K axL', 'K AX L', 'K L']]

    all_transcriptions = generate_valid_transcriptions( valid_transcriptions)

    print(len(all_transcriptions),all_transcriptions)

#def example_all_together():

def generate_expectation(multi_transcript):
    sep=","
    verdicts=["good","possible"]
    expectation=[]
    for item in multi_transcript:
        if len(item)==1:
            for phone in item[0].split(" "):
                expectation.append(phone)
                for v in verdicts:
                    expectation.append(v)
        else:
            expectation.append("(")
            for subitem in item[:-1]:
                for phone in subitem.split(" "):
                    expectation.append(phone)
                    for v in verdicts:
                        expectation.append(v)
                expectation.append("||")
            for phone in item[-1].split(" "):
                expectation.append(phone)
                for v in verdicts:
                    expectation.append(v)
            expectation.append(")")

    merge1=" ".join(expectation)
    merge2=merge1.split(" ")
    merge3=sep.join(merge2)+sep
    
    return merge3



def extract_rules(rules_file):
    cfg=toml.load(rules_file)
    
    sets=cfg['Sets']
    metarules=cfg['Meta-Rules']
    rules=cfg['Rules']
    
    all_rules=[]
    for mr in metarules:
        for s in sets:
            if s in metarules[mr][0]:
                
                for phone in sets[s]:
                    rule_instance=[]
                    for pattern in metarules[mr]:
                    
                        instance = pattern.replace(s,phone)
                        rule_instance.append(instance)
                    all_rules.append(rule_instance)
    
    for rule in rules:
        all_rules.append(rules[rule])

    return all_rules


def example_for_one_transcription():
    rules_file="./data/rules.toml"

    rules=extract_rules(rules_file)

    transcription = "K Y UW T IH K AX L"
    transcription = "A B S AX N T"
    transcription = "EHR M AX N"
    #transcription = "K EH M IH S T"
    transcription = "IH N V AY AX R axN M axN T"
    #transcription = "AE R M"


    multi_transcript = parser(transcription, rules)

    valid_transcriptions = generate_valid_transcriptions(multi_transcript)
    
    expectation = generate_expectation(multi_transcript)    
    print(transcription)
    print(multi_transcript)
    print(expectation)
    print(valid_transcriptions)


def extract_word(fileid):
    word = fileid.split("-")[-1].split("_")[0]
    return word

def extract_transcriptions_inputs( transcription_file ):

    with open(transcription_file, "r") as f:
        raw=f.read()

    transcriptions=raw.strip("\n").split("\n")

    inputs=[]
    
    for transcription in transcriptions:
        fileid = transcription.split("\t")[1][1:-1]
        word = extract_word(fileid)
        inputs.append([fileid,word])

    return transcriptions, inputs

def main():
    ## To add to unit testing:
    #example_creating_combinations_given_lengths()
    #parser_example()
    #example_combinations_of_valid_possibilities()
    example_for_one_transcription()
            

    #transcription_file="./data/art_db_Bare_train_Expanded.transcription"   

        

    print("working on it.")

   



if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
