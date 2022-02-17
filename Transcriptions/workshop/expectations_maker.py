import time
from transcriber import get_dictionary, get_dictionary_inv
from expression_parser import extract_rules, parser, generate_expectation

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


def strip_variant(variant):
    word=variant
    if '(' in variant:
        word = variant.split("(")[0]

    return word

def extract_word(phonemes, word_dictionary, audiofile):
    temptative=audiofile.split("-")[-1].split("_")[0].upper()
    if temptative in word_dictionary.values():
        word = temptative
    else:
        word = strip_variant(word_dictionary[phonemes])

    return word

def obtain_transcriptions_data_v3(transcription_file, dictionary, word_dictionary):
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
            word_entries.append(w_bare)

        phonemes=" ".join(raw_translation).strip(" ")    

        
        if len(word_entries)>1:
            print(f"{m}\t{i}\t{raw_transcription}")
            m=m+1
        elif len(word_entries)==1:    
            audiofiles.append(audiofile)
            
            word = extract_word(phonemes, word_dictionary, audiofile)
            #word = strip_variant(word_dictionary[phonemes])

            phonetic_transcriptions.append(phonemes)
            words.append(word.lower())
        else:
            print("!!! Error: wrong number of words.  --------------------------------------------------")

    return audiofiles, phonetic_transcriptions, words 


def create_inputs_file(inputs_file, audiofiles, words):

    with open(inputs_file, 'w') as f:
        for audiofile, word in zip(audiofiles, words):
            f.write(f"{audiofile},{word}\n")


def create_exepectation_files_v3(expectations_file, audiofiles, phonetic_transcriptions, words, rules):
    #rigorous
    with open(expectations_file, 'w') as f:
        for audiofile, phonetic_transcription, word in zip(audiofiles, phonetic_transcriptions, words):
            multi_transcript=parser(phonetic_transcription, rules)
            expectation=generate_expectation(multi_transcript)
            f.write(f"{audiofile}_{word},{expectation}\n")


def create_file_with_lines(file, lines):
    with open(file,'w') as f:
        contents = "\n".join(lines)+"\n"
        f.write(contents)


def create_expectations_and_inputs(expectations_file, inputs_file, audiofiles, phonetic_transcriptions, words, rules):

    unique_expectations={}
    unique_inputs={}
    for a,p,w in zip(audiofiles, phonetic_transcriptions, words):
        unique_inputs[f"{a},{w}"]=f"{a},{w}"
        multi_transcript=parser(p, rules)
        expectation=generate_expectation(multi_transcript)
        unique_expectations[f"{a}_{w}"]=f"{a}_{w},{expectation.lower()}"

    create_file_with_lines(inputs_file, list(unique_inputs.values()))    
    create_file_with_lines(expectations_file, list(unique_expectations.values()))    


def create_expectations_from_transcriptions_v3(transcription_file, dictionary, word_dictionary, expectations_file, inputs_file, rules_file):

        audiofiles, phonetic_transcriptions, words = obtain_transcriptions_data_v3(transcription_file, dictionary, word_dictionary)

        #create_inputs_file(inputs_file, audiofiles, words)

        rules=extract_rules(rules_file)
        #create_exepectation_files_v3(expectations_file, audiofiles, phonetic_transcriptions, words, rules)
        create_expectations_and_inputs(expectations_file, inputs_file, audiofiles, phonetic_transcriptions, words, rules)


def main():

    transcription_file="./data/art_db_Bare_train_Expanded.transcription" 
    expectations_file="./../../Expectations/train_expectations_v3.csv"
    inputs_file="./../../Tests/train_inputs_v3.csv"

    # transcription_file="./data/art_db_Bare_train_Expanded_debug.transcription" 
    # expectations_file="./../../Expectations/debug_train_expectations_v3.csv"
    # inputs_file="./../../Tests/debug_train_inputs_v3.csv"

    dictionary_file="./../../Dictionaries/art_db_v3_dummy.dic"
    word_dictionary_file="./../../Dictionaries/art_db_v3.dic"
    rules_file="./data/rules.toml"

    dictionary = get_dictionary(dictionary_file)
    word_dictionary = get_dictionary_inv(word_dictionary_file) #From transcript to word

    create_expectations_from_transcriptions_v3(transcription_file, dictionary, word_dictionary, expectations_file, inputs_file, rules_file)

    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")