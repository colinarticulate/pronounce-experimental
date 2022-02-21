'''
Given a normal word-based .transcription file and a dictionary. Returns .fileids file and .transcription based on dummy word (joined phonemes)
and also phonemes-based for different modalities of training.
'''

import os
import time
import shutil

from transcriber import get_dictionary, create_dummy_dictionary


def replace_audio_filenames_with_vowels_entry_for_dictionary():
    files_folder="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    files = [ f for f in os.listdir(files_folder) if "scrape-vowels-vowels-" in f and "_aug_" not in f and "_fixed" not in f]

    for i,file in enumerate(files):
        src = os.path.join(files_folder, file)

        vowel = file.split("-")[-1].split("_")[0]

        dict_vowel = f"-'{vowel}_"

        start = file.split("-")[:-1]
        end = file.split("-")[-1].split("_")[1:]

        new_filename = "-".join(start)+dict_vowel+"_".join(end)

        dst = os.path.join(files_folder, new_filename)
        shutil.copy(src,dst)
        print(f"{i}\t{file}\t\t{new_filename}")

def uw_case_to_reflect_dictionary_entry():
    files_folder="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    files = [ f for f in os.listdir(files_folder) if ("feb22__16000_mono_16bit" in f or "feb22_16000_mono_16bit" in f )and ("_aug_" not in f and "_fixed_" not in f and "_fixed" not in f)]

    for i,file in enumerate(files):
        src = os.path.join(files_folder, file)

    
        new_filename = f"'uw_{file}"

        dst = os.path.join(files_folder, new_filename)
        shutil.copy(src,dst)
        print(f"{i}\t{file}\t\t{new_filename}")


def adding_word_correctly_for_speech_commands_audios():
    files_folder="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    files = [ f for f in os.listdir(files_folder) if "SpeechCommands-" in f and ("_aug_" not in f and "_fixed_" not in f and "_fixed" not in f)]

    for i,file in enumerate(files):
        src = os.path.join(files_folder, file)

        word = file.split("-")[1]
        formatted_word = f"-{word}_"

        start = file.split("-")[:-1]
        end = file.split("-")[-1].split("_")

        new_filename = "-".join(start)+formatted_word+"_".join(end)

        dst = os.path.join(files_folder, new_filename)
        shutil.copy(src,dst)
        print(f"{i}\t{file}\t\t{new_filename}")

def read_lines(file):
    with open(file, 'r') as f:
        raw=f.read()
    lines=raw.strip("\n").split("\n")

    return lines

def write_list_in_lines(file, items_list):
    with open(file, 'w') as f:
       
        contents="\n".join(items_list)+"\n"
        f.write(contents)
   

def replace_damage_audios_from_transcript_and_fileid():
    transcription_file="./data/art_db_Bare_train_Expanded.transcription.old"
    fileids_file="./data/art_db_Bare_train_Expanded.fileids.old"
    removal_file="./data/removal_from_training.txt"

    transcriptions=read_lines(transcription_file)
    fileids=read_lines(fileids_file)
    removal=read_lines(removal_file)

    new_transcriptions=[]
    new_fileids=[]
        
    count=1
    for fileid,transcription in zip(fileids,transcriptions):
        filename=os.path.basename(fileid)
        if filename in removal:
            print(count,filename)
            count=count+1
        else:
            new_transcriptions.append(transcription)
            new_fileids.append(fileid)

#write in to a file...
    write_list_in_lines(transcription_file[:-4], new_transcriptions)
    write_list_in_lines(fileids_file[:-4], new_fileids)

    print("working")


def main():

    #uw_case_to_reflect_dictionary_entry()
    #adding_word_correctly_for_speech_commands_audios()
    replace_damage_audios_from_transcript_and_fileid()
    



if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
