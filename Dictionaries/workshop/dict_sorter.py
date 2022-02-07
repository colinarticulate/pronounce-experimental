import os
import time


def sort_dictionary(dict):
    with open(dict, 'r') as f:
        raw= f.read()

    entries=raw.strip("\n").split("\n")
    sorted_entries=sorted(entries)

    with open(f"{dict[:-4]}_sorted.dic", 'w') as f:
        f.write("\n".join(sorted_entries))

    dictionary={}

    for entry in sorted_entries:
        parts = entry.split(" ")
        word = parts[0]
        transcription = " ".join(parts[1:])
        dictionary[word]=transcription

    return dictionary


def search_both_fields(dictionary, field1, field2):

    results={}
    for entry in dictionary:
        if field1 in entry and field2 in dictionary[entry]:
            results[entry]=dictionary[entry]

    return results


def save_dictionary(file, dictionary):
    str_dictionary = [f"{entry} {dictionary[entry]}" for entry in dictionary]
    str_dictionary = "\n".join(str_dictionary)
    with open(file, 'w') as f:
        f.write(str_dictionary)


def main():
    pass

    dict_file="./data/art_db_v2_changing.dic"

    dictionary = sort_dictionary(dict_file)

    # results = search_both_fields(dictionary, "AIN", "EH")
    # save_dictionary("data/results.txt", results)

    save_dictionary("data/art_db_v2_changed_sorted.dic", dictionary)

    #7print(results)

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")