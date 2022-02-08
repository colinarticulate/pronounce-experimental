'''
Experiment Executor. Executes training and testing to obtain result for model comparison and test hypothesis.
'''

import os
import shutil
import subprocess
import time
import toml
import stat
import numpy as np

from train import train
from test import test
from collect_results import gather_test_pronounce_results
from report import create_report


def read_model_names(file):
    cfg=toml.load(file)

    model_names = cfg['models']['models']
    test_sets = cfg['test_sets']['datasets']

    return model_names, test_sets


def data_augmentation_experiment():
    experiment_name="Testing_Training_Data_Augmentation_no_pitch"
    experiment_file=f"./../experiments/{experiment_name}.toml"

    model_names, test_sets = read_model_names(experiment_file)

    for i,model_name in enumerate(model_names):
        print("---------------------------------------------------------------------------------------------")
        print(f"\nTraining {i+1}\tmodel:   {model_name} ")
        print("---------------------------------------------------------------------------------------------")
        config_folder=f"./../training_configurations/{model_name}"
        model = train(config_folder)
        model.fit()
        model.copy_model()

        for test_set in test_sets:
            print("---------------------------------------------------------------------------------------------")
            print(f"\nTesting model:\n{model_name}\n on dataset:\n {test_set}\n  ")
            print("---------------------------------------------------------------------------------------------")
            config_file = f"./../testing_configurations/{model_name}_x_{test_set}.toml"
            testing_bare = test(config_file, i+1)
            testing_bare.fit()



def create_traning_configuration(model_name, output_folder, toml_filename):
    with open(toml_filename,'w') as f:
        f.write("[test_pronounce_parameters]\n")
        f.write("\n")

        f.write("'-dict' = \"./../Dictionaries/art_db_v2.dic\"\n")
        f.write("'-phdict' = \"./../Dictionaries/art_db_v2_inference.phone\"\n")
        f.write("'-infolder' = \"/home/dbarbera/Data/audio_clips\"\n")
        f.write("'-tests' = \"./../Tests/pronounce_input.csv\"\n")
        f.write("'-expectations' = \"./../Expectations/expectations_v2.csv\"\n")
        f.write(f"'-outfolder' = \"{output_folder}\"\n")
        f.write(f"'-featparams' = \"./../Models/Bare/{model_name}.ci_cont/feat.params\"\n")
        f.write(f"'-hmm' = \"./../Models/Bare/{model_name}.ci_cont\"\n")



def test_Bare_several_times(n):
    experiment_name="Testing_new_Bare_with_UWs"
    experiment_file=f"./../experiments/{experiment_name}.toml"

    model_names, test_sets = read_model_names(experiment_file)

    toml_files=[]
    output_folders=[]
    for model_name in model_names:
        for test_set in test_sets:
            for i in range(n-2):
                print("---------------------------------------------------------------------------------------------")
                print(f"Iteration: {i+1}")
                print(f"\nTesting model:\n{model_name}\n on dataset:\n {test_set}\n  ")
                print("---------------------------------------------------------------------------------------------")
                config_file = f"./../testing_configurations/{model_name}_x_{test_set}_{i+1}.toml"
                toml_files.append(os.path.basename(config_file))
                output_folder = f"./../Test_Output/output_{model_name}_{i+1}"
                output_folders.append(output_folder)
                # create_traning_configuration(model_name, output_folder, config_file)
                # testing_bare = test(config_file)
                # testing_bare.fit()

    #gathering test results
    
    dst_results_folder = "./../Results"
    for i, (toml_file, src_folder) in enumerate(zip(toml_files, output_folders)):
        if i < 9:
            gather_test_pronounce_results(model_name,test_set, src_folder,dst_results_folder,toml_file)


    #Create report

    report_dir = "./../Reports"
    #experiment_name = "Data_augmentation"
    report_file = os.path.join(report_dir, f"{experiment_name}.xlsx")
    results_folder = "./../Results"

    #results_files = [f for f in os.listdir(results_folder) if f.endswith(".toml")]
    #results_files = ['Bare_loudness_speed_x_Test_Harness.toml']
    

    create_report(report_file, experiment_name, toml_files, results_folder)



def main():

    test_Bare_several_times(10)

    
    print("finished.main")

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")