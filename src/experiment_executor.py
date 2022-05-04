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



def create_testing_configuration(model_name, output_folder, audio_folder, expectation, input, toml_filename, dict_version=2):
    with open(toml_filename,'w') as f:
        f.write("[test_pronounce_parameters]\n")
        f.write("\n")

        f.write(f"'-dict' = \"./../Dictionaries/art_db_v{dict_version}.dic\"\n")
        f.write(f"'-phdict' = \"./../Dictionaries/art_db_v{dict_version}_inference.phone\"\n")
        f.write(f"'-infolder' = \"{audio_folder}\"\n")
        f.write(f"'-tests' = \"{input}\"\n")
        f.write(f"'-expectations' = \"{expectation}\"\n")
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
                # toml_files.append(os.path.basename(config_file))
                # output_folder = f"./../Test_Output/output_{model_name}_{i+1}"
                # output_folders.append(output_folder)
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


def test_Bare_on_train_data_and_two_expectations():
    experiment_name="Testing_Bare_with_training_data"
    experiment_file=f"./../experiments/{experiment_name}.toml"



    model_names, test_sets = read_model_names(experiment_file)

    expectations=["./../Expectations/train_expectations_rigorous.csv","./../Expectations/train_expectations_lenient.csv"]
    input = "./../Tests/train_inputs.csv"
    audio_folder="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"

    toml_files=[]
    output_folders=[]
    for model_name in model_names:
        for test_set in test_sets:
            for expectation in expectations:
                expectation_type=expectation.split("_")[-1].split(".")[0]
                print("---------------------------------------------------------------------------------------------")
                print(f"Iteration: {expectation}")
                print(f"\nTesting model:\n{model_name}\n on dataset:\n {test_set}\n  ")
                print("---------------------------------------------------------------------------------------------")
                config_file = f"./../testing_configurations/{model_name}_x_{test_set}_{expectation_type}.toml"
                toml_files.append(os.path.basename(config_file))
                output_folder = f"./../Test_Output/output_{model_name}_{test_set}_{expectation_type}"
                output_folders.append(output_folder)
                create_testing_configuration(model_name, output_folder, audio_folder, expectation, input, config_file)
                testing_bare = test(config_file)
                testing_bare.fit()

    #gathering test results
    
    dst_results_folder = "./../Results"
    for i, (toml_file, src_folder) in enumerate(zip(toml_files, output_folders)):
        gather_test_pronounce_results(model_name,test_set, src_folder,dst_results_folder,toml_file)


    #Create report

    report_dir = "./../Reports"
    #experiment_name = "Data_augmentation"
    report_file = os.path.join(report_dir, f"{experiment_name}.xlsx")
    results_folder = "./../Results"
    create_report(report_file, experiment_name, toml_files, results_folder)

def train_Bare(experiment_name, experiment_file):
    model_type="Bare"
    config_folder=f"./../training_configurations/{model_type}"
    model = train(config_folder)
    model_name = model.model_destination
    print("---------------------------------------------------------------------------------------------")
    print(f"\nTraining model:   {model_name} ")
    print("---------------------------------------------------------------------------------------------")
    
    model.fit()
    model.copy_model()

    return model_name


def double_experiment():

    experiment_name="Expanded_revisited"
    experiment_file=f"./../experiments/{experiment_name}.toml"

    model_name = train_Bare(experiment_name, experiment_file)

    #Create experiment file for the record:
    with open(experiment_file, 'w') as f:
        f.write(f"[models]\n")
        bare_name=os.path.basename(model_name).split(".")[0]
        f.write(f"'models' = [\"{bare_name}\"]")
        f.write(f"\n")
        f.write(f"[test_sets]\n")
        f.write(f"'datasets' = [\"Train_set\"]\n")

    return experiment_name, experiment_file    


def test_double_experiment(experiment_name, experiment_file):
    model_names, test_sets = read_model_names(experiment_file)
     
    expectations=["./../Expectations/train_expectations_v3.csv"]
    #expectations=["./../Expectations/expectations_v2.csv"]
    input = "./../Tests/train_inputs_v3.csv"
    #input = "./../Tests/pronounce_input.csv"
    audio_folder="/home/dbarbera/Repositories/art_db/wav/train/art_db_compilation"
    #audio_folder="/home/dbarbera/Data/audio_clips"

    toml_files=[]
    output_folders=[]
    for model_name in model_names:
        for test_set in test_sets:
            for expectation in expectations:
                expectation_type=os.path.basename(expectations[0][:-4])#expectation.split("_")[-1].split(".")[0]
                print("---------------------------------------------------------------------------------------------")
                print(f"Iteration: {expectation}")
                print(f"\nTesting model:\n{model_name}\n on dataset:\n {test_set}\n  ")
                print("---------------------------------------------------------------------------------------------")
                config_file = f"./../testing_configurations/{model_name}_x_{test_set}_{expectation_type}.toml"
                toml_files.append(os.path.basename(config_file))
                output_folder = f"./../Test_Output/output_{model_name}_{test_set}_{expectation_type}"
                output_folders.append(output_folder)
                create_testing_configuration(model_name, output_folder, audio_folder, expectation, input, config_file, dict_version=3)
                testing_bare = test(config_file)
                testing_bare.fit()

    #gathering test results
    
    dst_results_folder = "./../Results"
    for i, (toml_file, src_folder) in enumerate(zip(toml_files, output_folders)):
        gather_test_pronounce_results(model_name,test_set, src_folder,dst_results_folder,audio_folder, toml_file)


    #Create report

    report_dir = "./../Reports"
    #experiment_name = "Data_augmentation"
    report_file = os.path.join(report_dir, f"{experiment_name}_training_and_testing.xlsx")
    results_folder = "./../Results"
    create_report(report_file, experiment_name, toml_files, results_folder)



def main():

    #test_Bare_several_times(10)
    # experiment_name="Testing_Bare_with_training_data"
    # experiment_file=f"./../experiments/{experiment_name}.toml"
    # model_name = train_Bare(experiment_name, experiment_file)
    # test_Bare_on_train_data_and_two_expectations([model_name])
    #test_Bare_on_train_data_and_two_expectations()

    double_experiment()
    experiment_name="Expanded_revisited"
    experiment_file=f"./../experiments/{experiment_name}.toml"
    #test_double_experiment(experiment_name, experiment_file)

    
    print("finished.main")

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")