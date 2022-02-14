import time
from transcriber import get_dictionary

def parse_out_symbols(word_list, dictionary):
    words=[]
    for w in word_list:
        phonemes=dictionary[w]
        if phonemes!="":
            words.append(w)

    return words


def variant_format(word):
    variant=word

    # if "'" in variant:
    #     variant=variant.strip("'")

    if "(" in variant and ")" in variant:
        number=variant.split("(")[1].split(")")[0]
        raw_word = variant.split("(")[0]
        variant=f"{raw_word}_{number}"

    return variant

def strip_variant_number(variant):
    word = variant
    if "(" in variant:
        word= variant.split("(")[0]

    return word


def obtain_transcriptions_data(transcription_file, dictionary):
#
#    We are assuming no repetitions in the transcriptions (i.e., no data augmentation)
#
    with open(transcription_file, 'r') as f:
        raw=f.read()
    raw_transcriptions=raw.strip("\n").split("\n")

    dictionary["<sil>"] = ""
    dictionary["<s>"] = ""
    dictionary["</s>"] = ""

    audiofiles=[]
    phonetic_transcriptions=[]
    words=[]
    m=1
    for i, raw_transcription in enumerate(raw_transcriptions):
        parts=raw_transcription.split("\t")
        transcription=parts[0]
        audiofile= parts[1][1:-1] #because of the parenthesis
        word_list = transcription.split(" ")
        raw_words = parse_out_symbols(word_list, dictionary)
        raw_translation=[]
        word_entries=[]
        for w in raw_words:
            raw_translation.append(dictionary[w])
            #variant=variant_format(w)
            w_bare=strip_variant_number(w)
            word_entries.append(w_bare.lower())

        phonemes=" ".join(raw_translation).strip(" ").lower()    

        
        if len(word_entries)>1:
            print(f"{m}\t{i}\t{raw_transcription}")
            m=m+1
        elif len(word_entries)==1:    
            audiofiles.append(audiofile)
            phonetic_transcriptions.append(phonemes)
            words.append("-".join(word_entries))
        else:
            print("!!! Error: wrong number of words.  --------------------------------------------------")

    return audiofiles, phonetic_transcriptions, words 


def create_inputs_file(inputs_file, audiofiles, words):

    with open(inputs_file, 'w') as f:
        for audiofile, word in zip(audiofiles, words):
            f.write(f"{audiofile},{word}\n")


def create_exepectation_files(expectations_file, audiofiles, phonetic_transcriptions, words):
    #rigorous
    with open(expectations_file[:-4]+"_rigorous.csv", 'w') as f:
        for audiofile, phonetic_transcription, word in zip(audiofiles, phonetic_transcriptions, words):
            f.write(f"{audiofile}_{word},")
            for phoneme in phonetic_transcription.split(" "):
                f.write(f"{phoneme.lower()},good,")
            f.write("\n")

    #lenient
    with open(expectations_file[:-4]+"_lenient.csv", 'w') as f:
        for audiofile, phonetic_transcription, word in zip(audiofiles, phonetic_transcriptions, words):
            f.write(f"{audiofile}_{word},")
            for phoneme in phonetic_transcription.split(" "):
                f.write(f"{phoneme.lower()},good,possible,")
            f.write("\n")


def create_expectations_from_transcriptions(transcription_file, dictionary, expectations_file, inputs_file):

        audiofiles, phonetic_transcriptions, words = obtain_transcriptions_data(transcription_file, dictionary)

        create_inputs_file(inputs_file, audiofiles, words)

        create_exepectation_files(expectations_file, audiofiles, phonetic_transcriptions, words)


def main():

    word_based_transcription_file="data/art_db_new_train_noDummy.transcription"
    expectations_file="data/train_expectations.csv"
    inputs_file="data/train_inputs.csv"
    dictionary_file="./../../Dictionaries/art_db_v2.dic"

    dictionary = get_dictionary(dictionary_file)
    create_expectations_from_transcriptions(word_based_transcription_file, dictionary, expectations_file, inputs_file)

    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")