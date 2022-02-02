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


def read_model_names(file):
    cfg=toml.load(file)

    model_names = cfg['models']['models']
    test_sets = cfg['test_sets']['datasets']

    return model_names, test_sets


def main():
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
            testing_bare = test(config_file)
            testing_bare.fit()

    
    print("finished.main")

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")