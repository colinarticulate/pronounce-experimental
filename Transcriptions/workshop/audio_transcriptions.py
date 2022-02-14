'''
Given a normal word-based .transcription file and a dictionary. Returns .fileids file and .transcription based on dummy word (joined phonemes)
and also phonemes-based for different modalities of training.
'''

import os
import time
import shutil


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


def main():

    folder = "/home/dbarbera/Data/"
    word = "TWO"   

    dst_audio="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation" 
    dst_transcription="./data"

    words = ['TWO', 'THREE', 'EIGHT']

    for word in words:
        word_folder = os.path.join(folder, f"{word.lower()}")

        copy_audios_and_create_transcription( word_folder, word, dst_audio, dst_transcription)
   



    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")
