import time

def get_dictionary( file ):
    with open(file, 'r') as f:
        lines = f.readlines()

    dictionary={}
    for line in lines:
        parts = line.strip("\n").split(" ")
        entry = parts[0]

        if len(parts[1:])>1:
            transcription = " ".join(parts[1:])
        else:
            transcription = parts[1]

        dictionary[entry]=transcription

    return dictionary


def convert_transcription_to_dummy_train(dictionary, transcriptions_in, transcriptions_out):

    with open(transcriptions_in, "r") as fr:
        raw =fr.read()

    lines=raw.strip("\n").split("\n")

    dictionary['<sil>']='<sil>'
    dictionary['<s>']='<s>'
    dictionary['</s>']='</s>'

    new_lines=[]
    with open(transcriptions_out, "w") as fw:
        for i,line in enumerate(lines):
            parts = line.split("\t")
            file= parts[-1]
            transcription = parts[0]
            transcription_words = transcription.split(" ")

            new_transcription_words = []
            for word in transcription_words:
                new_word = "_".join(dictionary[word].split(" "))
                new_transcription_words.append(new_word)

            new_transcription = " ".join(new_transcription_words)

            new_line = f"{new_transcription}\t{file}\n"
            fw.write(new_line)

            


    #         entry = line.split("<sil>")[1].strip(" ")
            
    #         entries = entry.split(" ")
    #         if len(entries) > 1:
    #             all_new_entries=[]
    #             for entry_ in entries:
    #                 new_entry_= "_".join(dictionary[entry_].split(" "))
    #                 all_new_entries.append(new_entry_)
    #             new_entry = " ".join(all_new_entries)
    #         else:
    #             new_entry = "_".join(dictionary[entry].split(" "))

    #         parts = line.split(entry)
    #         new_line = "".join([parts[0], new_entry, parts[1]])+"\n"
    #         print(f"{i+1}\t\t{new_line[:-1]}")
    #         fw.write(new_line)
    #         new_lines.append(new_line)

    # print(f"Train Length: {len(new_lines)}")


def convert_transcription_to_dummy_test(dictionary, transcriptions_in, transcriptions_out):

    # with open(transcriptions_in, "r") as fr:
    #     raw =fr.read()

    # lines=raw.strip("\n").split("\n")

        
    # with open(transcriptions_out, "w") as fw:
    #     for i,line in enumerate(lines):
    #         a = line.split("<s>")[1].strip(" ")
    #         entry = a.split("</s>")[0].strip(" ")
            
    #         entries = entry.split(" ")
    #         if len(entries) > 1:
    #             all_new_entries=[]
    #             for entry_ in entries:
    #                 new_entry_= "_".join(dictionary[entry_].split(" "))
    #                 all_new_entries.append(new_entry_)
    #             new_entry = " ".join(all_new_entries)
    #         else:
    #             new_entry = "_".join(dictionary[entry].split(" "))

    #         parts = line.split(entry)
    #         new_= "".join([parts[0], new_entry, parts[1]])
    #         new_line=f"{new_}\n"
    #         fw.write(new_line)
    with open(transcriptions_in, "r") as fr:
        raw =fr.read()

    lines=raw.strip("\n").split("\n")


    new_lines=[]
    with open(transcriptions_out, "w") as fw:
        for i,line in enumerate(lines):
            
            entry = line.split("<sil>")[1].strip(" ")
            
            entries = entry.split(" ")
            if len(entries) > 1:
                all_new_entries=[]
                for entry_ in entries:
                    new_entry_= "_".join(dictionary[entry_].split(" "))
                    all_new_entries.append(new_entry_)
                new_entry = " ".join(all_new_entries)
            else:
                new_entry = "_".join(dictionary[entry].split(" "))

            parts = line.split(entry)
            new_line = "".join([parts[0], new_entry, parts[1]])+"\n"
            print(f"{i+1}\t\t{new_line[:-1]}")
            fw.write(new_line)
            new_lines.append(new_line)

    print(f"Test Length: {len(new_lines)}")



def create_dummy_dictionary(dictionary, new_dictionary_file):

    dummy=[]
    for transcript in list(dictionary.values()):
        entry="_".join(transcript.split(" "))
        line=f"{entry} {transcript}\n"
        dummy.append(line)

    new_dummy=sorted(list(set(dummy)))
    #new_dummy=sorted(dummy)

    print(f"Original dictionary # entries: {len(dictionary.values())}")
    print(f"Dummy # dictionary entries: {len(dummy)}, # dummy entries: {len(new_dummy)}")

    with open(new_dictionary_file, 'w') as f:
        for entry in new_dummy:
            f.write(entry)



def main():
    
    #dictionary_file="art_db_oldest.dic"
    #New versions:
    dictionary_file="./../../Dictionaries/art_db_v2.dic"
    transcriptions_file="data/art_db_new_train_noDummy.transcription"

    #Need updating from above files:
    new_transcriptions_file="data/art_db_new_train_dummy.transcription"
    new_dictionary_file="data/art_db_v2_dummy.dic"

    #No need at the moment until we want to obtain phone error rate
    # test_transcription_file="art_db_test.transcription"
    # new_test_transcription_file="art_db_test_dummy.transcription"

    dictionary=get_dictionary(dictionary_file)

    convert_transcription_to_dummy_train(dictionary, transcriptions_file, new_transcriptions_file)
    #convert_transcription_to_dummy_test(dictionary, test_transcription_file, new_test_transcription_file)

    create_dummy_dictionary(dictionary,new_dictionary_file)


    


if __name__ == '__main__':

    start=time.time()
    main()
    stop=time.time()

    print("Finished.")
    print(f"Time: {stop-start} seconds.")