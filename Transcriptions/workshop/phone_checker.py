'''
Given a normal word-based .transcription file and a dictionary. Returns .fileids file and .transcription based on dummy word (joined phonemes)
and also phonemes-based for different modalities of training.
'''

import os
import time
import shutil

from transcriber import get_dictionary, create_dummy_dictionary



def discard_symobls(word_list):
    accepted = []
    for word in word_list:
        if word != '<s>' and word != '<sil>' and word != '</s>':
            accepted.append(word)

    return accepted


def extract_all_phones_from_dummy_transcription(transcription_file):

    with open(transcription_file, 'r') as f:
        raw=f.read()

    transcriptions = raw.strip("\n").split("\n")

    phones=[]
    for i, transcription in enumerate(transcriptions):
        raw = transcription.split("\t")[0].split(" ")

        dummy_words = discard_symobls(raw)

        for dummy_word in dummy_words:
            phones_list = dummy_word.split("_")
            if 'tSH' in phones_list:
                print(i, transcription)
            phones=phones+phones_list

        
        
        phones=list(set(phones))
        
    return sorted(phones)


def compare_phone_lists(file1, file2):
    with open(file1, 'r') as f:
        raw=f.read()
    phones1=raw.strip("\n").split("\n")

    with open(file2, 'r') as f:
        raw=f.read()
    phones2=raw.strip("\n").split("\n") 

    print("phones from phones2 not in phones1")
    for phone in phones2:
        if phone not in phones1:
            print(phone)  

    print("phones from phones1 not in phones2")
    for phone in phones1:
        if phone not in phones2:
            print(phone)  


def save_phones(file, phones):
    with open(file, 'w') as f:
        f.write("\n".join(phones))
        f.write("\n")


def main():

    transcription_file="./data/art_db_Bare_train_Double.transcription"
    phones = extract_all_phones_from_dummy_transcription(transcription_file)
    print(phones, len(phones))

    file="./data/phones_from_transcription.txt"
    save_phones(file,phones)


    file_reference="./data/art_db_v2.phone"
    compare_phone_lists(file_reference, file)
    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")