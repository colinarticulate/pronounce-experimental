'''
Collect Results. Gathers all results from all experiments with test_pronounce output format.
Test_pronounce output format is a folder with a summary file and text files for each test.
This is to transform the output of test pronunce into one single '.toml' file
'''

import os
import shutil
import subprocess
import time
import datetime
import toml
import stat
import numpy as np
from collections import defaultdict

from train import train
from test import test


def get_test_input_result(file):
    #from listen.go in Pron:
    pron_verdict = {
                0:"good",
                1:"possible",
                2:"missing",
                3:"surprise"
    }

    with open(file, 'r') as f:
        contents = f.read()

    result=[]
    # if "lettersVerdicts = [" in contents and "]\npublish<-" in contents:
    #     part = contents.split("lettersVerdicts = [")[1]
    #     raw_result = part.split("]\npublish<-")[0].strip("\n")
    #     entries = raw_result[1:-1].split("} {")

    #     for entry in entries:
    #         parts1 = entry.split(" [")
    #         parts2 = parts1[1].split("] ")
    #         letter = parts1[0]
    #         phoneme_transcription = parts2[0]
    #         int_verdict = parts2[1]

    #         verdict = pron_verdict[int(int_verdict)]

    #         phonemes = phoneme_transcription.split(" ")

    #         for phoneme in phonemes:
            
    #             result.append([phoneme,verdict])
    if "testPronounce" in contents and "publish->" in contents:
        raw_result=contents.split("testPronounce ")[1].split("publish->")[0].strip("\n")
        result_list=raw_result.split(" ")
        phonemes=[]
        verdicts=[]
        for i, r in enumerate(result_list):
            if i%2==0:
                phonemes.append(r)
            else:
                verdicts.append(r)

        result=[]
        for phoneme, verdict in zip(phonemes,verdicts):
            result.append([phoneme,verdict])

     
    else:
        print(f"Error: file {file} has no explicit result written. Please check.")
        result.append(" ")
    

    return result



def get_summary_data(file, results):
    with open(file, 'r') as f:
        raw = f.read()

    lines = raw.strip("\n").split("\n")

    txt_predictions={}

    for line in lines:
        if "Pass rate = " in line:
            accuracy = line.split("Pass rate = ")[1]

        else:
            parts = line.split(",")
            test_name = parts[0]
            prediction = parts[1].strip(" ")
            txt_predictions[test_name] =  prediction

    #Need to do this, otherwise if a txt file was corrupted, its entry won't be in the summary file.
    txt_predictions_keys=list(txt_predictions.keys())
    results_keys = list(results.keys())
    predictions={}
    for result_key in results_keys:
        if result_key not in txt_predictions_keys:
            predictions[result_key] = ""
        else:
            predictions[result_key]=txt_predictions[result_key]
            

    return accuracy, predictions


def gather_test_pronounce_results(model_name, dataset_name, src_results_folder, dst_results_folder, toml_filename=None):

    if toml_filename == None:
        toml_filename = f"{model_name}_x_{dataset_name}.toml"
    else:
        toml_filename = toml_filename

    #with open(os.path.join(dst_results_folder, toml_filename), 'w') as f:

    files = [ f for f in os.listdir(src_results_folder) if "000__summary__000" not in f]
    files = sorted(files) #Sorted alphabetically
    summary_file = "000__summary__000.txt" # [ f for f in os.listdir(src_results_folder) if "000__summary__000" in f][0]

    results={}#defaultdict(list)
    accuracy=0
    predictions={}
    for file in files[:]:
          
        result = get_test_input_result(os.path.join(src_results_folder,file))
        results[file[:-4]] = result

    accuracy, predictions = get_summary_data(os.path.join(src_results_folder, summary_file), results)

    
    with open(os.path.join(dst_results_folder,toml_filename),'w') as f:
        f.write(f"[info]\n")
        f.write(f"\"model_name\" = \"{model_name}\"\n")
        f.write(f"\"dataset\" = \"{dataset_name}\"\n")
        f.write("\n")
        f.write("[performance]\n")
        f.write(f"\"accuracy\" = \"{accuracy}\"\n")
        f.write("\n")
        f.write(f"[predictions]\n")
        for prediction in predictions:
            f.write(f"\"{prediction}\" = \"{predictions[prediction]}\"\n")


        f.write("\n")
        f.write("[results]\n")
        for result in results:
            f.write(f"\"{result}\" = {results[result]}\n")

  

def main():
    #--------------------------------------------------------------------------------------
    # model_name = "Bare"
    # dataset_name = "Test_Harness"
    # output_folder = "./../Test_Output"
    # #src_results_folder = "./../Test_Output"
    # dst_results_folder = "./../Results"

    # #gather_test_pronounce_results(model_name, dataset_name, src_results_folder, dst_results_folder)

    # folders = [f for f in os.listdir(output_folder) if os.path.isdir(os.path.join(output_folder,f))]

    # for folder in folders:

    #     model_name = folder.split("output_")[-1]
    #     dataset_name = "Test_Harness"
    #     output_folder = "./../Test_Output"
    #     src_results_folder= os.path.join(output_folder,folder)
    #     gather_test_pronounce_results(model_name, dataset_name, src_results_folder, dst_results_folder)
    #--------------------------------------------------------------------------------------------


    file="/home/dbarbera/Repositories/pronounce-experimental/Results/2022-02-15T10:29:55-062_Bare_x_Train_set_train_expectations_v3_.toml"
    output="/home/dbarbera/Repositories/pronounce-experimental/Test_Output/output_2022-02-15T10:29:55-062_Bare_Train_set_train_expectations_v3_"
    summary_file="/home/dbarbera/Repositories/pronounce-experimental/Test_Output/output_2022-02-15T10:29:55-062_Bare_Train_set_train_expectations_v3_/000__summary__000.txt"

    with open(summary_file) as f:
        raw=f.read()

    lines=raw.strip("\n").split("\n")

    pass_rate=lines[-1]
    predictions=[f.split(",")[0] for f in lines[:-1]]

    

    files = [f[:-4] for f in os.listdir(output) if f.endswith(".txt") and "000__summary__000" not in f]

    print(len(predictions), len(files))

    diff1= list(set(predictions)-set(files))
    diff2=list(set(files)-set(predictions))

    pset=np.array(predictions)
    fset=np.array(files)

    ndiff1=np.setdiff1d(pset,fset)
    ndiff2=np.setdiff1d(fset,pset)

    print(ndiff2)
    print(len(list(set(predictions))))

    counts={}
    for prediction in predictions:
        if prediction in counts.keys():
            counts[prediction]=counts[prediction]+1
        else:
            counts[prediction]=0

    print(len(counts.keys()))

    print("finished.main")  

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")