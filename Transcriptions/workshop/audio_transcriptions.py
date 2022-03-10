'''
Given a normal word-based .transcription file and a dictionary. Returns .fileids file and .transcription based on dummy word (joined phonemes)
and also phonemes-based for different modalities of training.
'''

import os
import time
import shutil

from transcriber import get_dictionary, create_dummy_dictionary, create_file_ids

def copy_audios_and_create_transcription( folder, word, dst_audio, dst_transcription):

    audio_files = [f for f in os.listdir(folder) if f.endswith(".wav")]


    fileids=[]
    transcriptions=[]
    for file in sorted(audio_files):
        src_file = os.path.join(folder,file)
        new_file_name = f"{word.lower()}-{file[:-4]}"
        dst_file = os.path.join(dst_audio,new_file_name+".wav")
        shutil.copy(src_file, dst_file)

        fileids.append(f"train/art_db_compilattion/{new_file_name}")
        transcriptions.append(f"<s> <sil> {word.upper()} <sil> </s>\t({new_file_name})")

    with open(os.path.join(dst_transcription,f"{word.lower()}.fileids"), 'w' ) as f:
        f.write("\n".join(fileids))
        f.write("\n")


    with open(os.path.join(dst_transcription,f"{word.lower()}.transcription"), 'w' ) as f:
        f.write("\n".join(transcriptions))
        f.write("\n")


def discard_symobls(word_list):
    accepted = []
    for word in word_list:
        if word != '<s>' and word != '<sil>' and word != '</s>':
            accepted.append(word)

    return accepted


def create_dummy_dictionary_entries( transcriptions):

    dictionary_entries={}

    for entry in transcriptions:

        transcription = entry.split("\t")[0]
        word_list = transcription.split(" ")
        dummy_words = discard_symobls(word_list)
        for dummy_word in dummy_words:
            dictionary_entries[dummy_word] = " ".join(dummy_word.split("_"))
    
    return dictionary_entries


def merge_dummy_dictionaries( dict1, dict2):

    dict_merged={}
    for entry in dict1:
        dict_merged[entry] = dict1[entry]

    for entry in dict2:
        dict_merged[entry] = dict2[entry]

    return dict_merged

def create_fileids_from_transcription( transcription_file, fileids_file, audio_folder):

    with open(transcription_file, "r") as f:
        raw=f.read()

    transcriptions=raw.strip("\n").split("\n")

    fileids=[]
    with open(fileids_file, 'w') as f:
        for transcription in transcriptions:
            fileid = transcription.split("\t")[1][1:-1]
            f.write(f"{audio_folder}{fileid}\n")
            fileids.append(fileid)

    return transcriptions, fileids
    
def load_dummy(dummy_dict_file):

    with open(dummy_dict_file, 'r') as f:
        raw=f.read()

    entries = raw.strip("\n").split("\n")

    dictionary = {}

    for entry in entries:
        parts= entry.split(" ")
        dummy_word=parts[0]
        transcription = " ".join(parts[1:])
        dictionary[dummy_word]=transcription

    return dictionary

def save_dummy_dict(filename, dictionary):
    with open(filename, 'w') as f:
        for dummy_word in dictionary:
            f.write(f"{dummy_word} {dictionary[dummy_word]}\n")


def discard_entries_from_training_files(to_discard_file, fileids, transcriptions):
    with open(to_discard_file, 'r') as f:
        raw=f.read()

    missing_files = raw.strip("\n").split("\n")

    clean_fileids=[]
    clean_transcriptions=[]
    for fileid, transcription in zip(fileids, transcriptions):
        filename = transcription.split("\t")[1][-1:1]
        if filename not in missing_files:
            clean_fileids.append(fileid)
            clean_transcriptions.append(transcription)

    return clean_fileids, clean_transcriptions


def given_dummy_transcriptions_create_fileids_and_an_update_general_dummy_dict():
    # folder = "/home/dbarbera/Data/"
    # word = "TWO"   

    # dst_audio="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation" 
    # dst_transcription="./data"

    # words = ['TWO', 'THREE', 'EIGHT']

    # for word in words:
    #     word_folder = os.path.join(folder, f"{word.lower()}")

    #     copy_audios_and_create_transcription( word_folder, word, dst_audio, dst_transcription)
   

    # transcription_file="./data/art_db_Bare_train_Double.transcription"
    # fileids_file="./data/art_db_Bare_train_Double.fileids"
    transcription_file="./data/art_db_Bare_train_Expanded.transcription"
    fileids_file="./data/art_db_Bare_train_Expanded.fileids"
    audio_folder="train/art_db_compilation/"

    transcriptions, fileids = create_fileids_from_transcription( transcription_file, fileids_file, audio_folder)
    create_file_ids(transcriptions, fileids_file, audio_folder)
    #raw_transcriptions, raw_fileids = create_fileids_from_transcription( transcription_file, fileids_file, audio_folder)

    #to_discard_file="./data/missing_not_found.txt"
    #transcriptions, fileids = discard_entries_from_training_files(to_discard_file, raw_fileids, raw_transcriptions)
    
    dummy_entries = create_dummy_dictionary_entries(transcriptions)

    dummy_dict_file="./../../Dictionaries/art_db_v3_dummy.dic"
    dictionary = get_dictionary("./../../Dictionaries/art_db_v3.dic") #create dummy again
    create_dummy_dictionary(dictionary, dummy_dict_file)

    
    current_dummy = load_dummy(dummy_dict_file)

    merged_dummy_entries = merge_dummy_dictionaries(current_dummy, dummy_entries)

    filename="./data/art_db_v3_dummy_new.dic"
    save_dummy_dict(filename, merged_dummy_entries)



def check_and_create_missing_audios(missing_audios_file, src_path, dst_path):

    with open(missing_audios_file, 'r') as f:
        raw=f.read()

    audiofiles=raw.strip("\n").split("\n")

    with open("./data/missing_not_found.txt",'w') as f:
        for audiofile in audiofiles:
            src = "/".join(audiofile.split("-"))
            src_file = os.path.join(src_path, src+".wav")
            dst_file = os.path.join(dst_path, audiofile+".wav")
            if os.path.exists(src_file):
                shutil.copy(src_file, dst_file)
            else:
                f.write(f"{audiofile}\n")
                
        

def fix_naming_audios():

    replace_dict={
        "mustn_t":"mustn't",
        "i_m":"i'm",
        "don_t":"don't",
        "i_ve":"i've",
        "it_s":"it's",
        "let_s":"let's",
        "that_s":"that's",
        "they_re":"they're",
        "you_re":"you're"
    }

    audios_dir="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"
    transcription_file="./data/art_db_Bare_train_Expanded.transcription"
    fileids_file="./data/art_db_Bare_train_Expanded.fileids"

    audios = [f for f in os.listdir(audios_dir) if f.endswith(".wav")]
    
    for audio in audios:
        for r in replace_dict.keys():
            if r in audio:
                fixed_audio = replace_dict[r].join(audio.split(r))
                shutil.copy(os.path.join(audios_dir,audio),os.path.join(audios_dir, fixed_audio))
                break

    with open(transcription_file,'r') as f:
        raw=f.read()
    transcriptions=raw.strip("\n").split("\n")

    fixed_transcriptions=[]
    for transcription in transcriptions:
        fixed_transcription = transcription
        for r in replace_dict.keys():
            if r in transcription:
                fixed_transcription = replace_dict[r].join(transcription.split(r))

        fixed_transcriptions.append(fixed_transcription)

    with open(transcription_file, 'w') as f:
        f.write("\n".join(fixed_transcriptions)+"\n")


    with open(fileids_file,'r') as f:
        raw=f.read()
    fileids=raw.strip("\n").split("\n")

    fixed_fileids=[]
    for fileid in fileids:
        fixed_fileid = fileid
        for r in replace_dict.keys():
            if r in fileid:
                fixed_fileid = replace_dict[r].join(fileid.split(r))

        fixed_fileids.append(fixed_fileid)

    with open(fileids_file, 'w') as f:
        f.write("\n".join(fixed_fileids)+"\n")


def fix_audiofilenames_in_transcriptions_and_fileids():
  
    transcription_file="./data/art_db_Bare_train_Expanded.transcription"
    fileids_file="./data/art_db_Bare_train_Expanded.fileids"

 
    with open(transcription_file,'r') as f:
        raw=f.read()
    transcriptions=raw.strip("\n").split("\n")

    locator="scrape-paul2-paul2-"

    fixed_transcriptions=[]
    for transcription in transcriptions:
        fixed_transcription = transcription
        if locator in transcription:
            fixed_transcription = "_".join(transcription.split("_")[:-1]) +")"

        fixed_transcriptions.append(fixed_transcription)

    with open(transcription_file, 'w') as f:
        f.write("\n".join(fixed_transcriptions)+"\n")


    with open(fileids_file,'r') as f:
        raw=f.read()
    fileids=raw.strip("\n").split("\n")

    fixed_fileids=[]
    for fileid in fileids:
        fixed_fileid = fileid
        if locator in fileid:
            fixed_fileid = "_".join(fileid.split("_")[:-1])

        fixed_fileids.append(fixed_fileid)

    with open(fileids_file, 'w') as f:
        f.write("\n".join(fixed_fileids)+"\n")

def delete_invalid_audios():

    path_dir = "/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    files = [f for f in os.listdir(path_dir) if "_f90" in f]

    for file in files:
        os.remove(os.path.join(path_dir, file))



def main():

    # missing_audios_file="./data/missing_audios.txt"
    # src_path="/home/dbarbera/Repositories/art_db/wav/train"
    # dst_path="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    # check_and_create_missing_audios(missing_audios_file, src_path, dst_path)

    given_dummy_transcriptions_create_fileids_and_an_update_general_dummy_dict()
    #delete_invalid_audios()

    #fix_naming_audios()


    #This a one-of fix
    #fix_audiofilenames_in_transcriptions_and_fileids()    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
