'''
Train a model using test_pronounce
'''

import os
import shutil
import subprocess
import time
import toml
import stat


class test():

    def __init__(self, config_file):
        self.cfg = toml.load(config_file)
        self.test_pronounce_location="./../test_pronounce"

        self.log_folder="./../logs"
        log_model_folder = os.path.join(self.log_folder, os.path.basename(config_file[:-4]) )
        if not os.path.exists(log_model_folder):
            os.mkdir(log_model_folder)
        self.log_file = os.path.join(log_model_folder, "testing.log")
    
    def execute(self, command, cwd, file):
    # should be in a utils module
        with open(file, 'w') as f:
            process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, cwd=cwd, universal_newlines=True)

            while True:
                output = process.stdout.readline()
                print(output.strip())
                f.write(output)
                # Do something else
                return_code = process.poll()
                if return_code is not None:
                    print('\n>>> RETURN CODE', return_code)
                    f.write(f"\n>>> RETURN CODE {return_code}\n")
                    # Process has finished, read rest of the output 
                    for output in process.stdout.readlines():
                        print(output.strip())
                        f.write(output)
                    break

    def fit(self):

        cfg=self.cfg['test_pronounce_parameters']

        command=["./test_pronounce"]
        for param in cfg:
            command.append(param)
            command.append(cfg[param])


        command=" ".join(command)
        cwd = self.test_pronounce_location
        
        print("Started testing ...")
        self.execute(command, cwd, self.log_file)
        print("Testing finished.")
        print(f"Check results in the folder: \n{os.path.normpath(os.path.join(self.test_pronounce_location,self.cfg['test_pronounce_parameters']['-outfolder']))}\n")





def main():
    #config_file= os.path.normpath(os.path.join(os.getcwd(), "./../testing_configurations/Bare_pitch_x_Test_Harness.toml"))
    config_file = "./../testing_configurations/2022-02-07T14:11:46-092_Bare_x_Test_Harness.toml"

    testing_bare = test(config_file)

    testing_bare.fit()

    


if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")