'''
Train a model using sphinxtrain
'''

import os
import shutil
import subprocess
import time
import datetime
import toml
import stat


class train():

    def __init__(self, config_folder):
        base_name = os.path.basename(config_folder)
        self.cfg = toml.load(os.path.join(config_folder, f"configuration.toml"))
        self.training_config_file = os.path.join(config_folder, "etc", "sphinx_train.cfg")

        print(f"Creating sphinx_train.cfg at {self.training_config_file} .")
        cfg_template_file="./../templates/sphinx_train_template.cfg"
        with open(cfg_template_file, 'r') as f:
            self.template=f.read()

        self.create_training_configuration_file()

        self.base_dir = self.cfg['training_configuration']['__CFG_BASE_DIR__']
        self.etc = os.path.join(self.base_dir,"etc")    
        self.model_name = self.cfg['training_configuration']['__CFG_EXPTNAME__']
        
        
        now = datetime.datetime.now()
        timestamp = now.strftime('%Y-%m-%dT%H:%M:%S') + ('-%03d' % (now.microsecond / 10000))
        self.timestamp_model_name = f"{timestamp}_{self.model_name}"
        self.model_location = os.path.join(self.base_dir, "model_parameters", f"{self.model_name}.ci_cont")

        model_destination_path = os.path.join("./../Models",self.model_name)
        self.model_destination = os.path.join(model_destination_path, f"{self.timestamp_model_name}.ci_cont","model") 
        self.model_etc_destination = os.path.join(model_destination_path, f"{self.timestamp_model_name}.ci_cont","etc") 

        print("Preparing training...")
        self.prepare_training()

        self.log_folder="./../logs"
        log_model_folder = os.path.join(self.log_folder, self.model_name )
        if not os.path.exists(log_model_folder):
            os.mkdir(log_model_folder)
        self.log_file = os.path.join(log_model_folder, "training.log")


        

    def create_directory(self, directory):
        if not os.path.exists(directory):
            os.makedirs(directory)

    def replace_template_variables(self):
        configuration=self.template
        cfg=self.cfg['training_configuration']
        for src in cfg:
            dst = cfg[src]
            configuration=f"{dst}".join(configuration.split(src))

        return configuration


    def create_training_configuration_file(self):
        
        configuration = self.replace_template_variables()

        with  open(self.training_config_file, 'w') as f:
            f.write(configuration)


    def delete_files_from_folder(self, folder):
    # this could go into a utils module
        for filename in os.listdir(folder):
            file_path = os.path.join(folder, filename)
            try:
                if os.path.isfile(file_path): #or os.path.islink(file_path):
                    os.remove(file_path)
                # elif os.path.isdir(file_path):
                #     shutil.rmtree(file_path)
            except Exception as e:
                print(f"Failed to delete {file_path}. Reason: {e}")


    def copy_dictionaries_to_etc(self):
    # this could go into a utils module
        dictionaries = self.cfg['dictionaries_for_training']

        for dictionary in dictionaries:
            src = dictionaries[dictionary]
            name = os.path.basename(src)
            dst = os.path.join(self.etc,name)
            shutil.copy(src, dst)


    def copy_training_files_to_etc(self):
    # this could go into a utils module
        folder = self.cfg['files_for_training']['main_files_folder']

        for filename in os.listdir(folder):
            src = os.path.join(folder, filename)
            dst = os.path.join(self.etc, filename)
            shutil.copy(src, dst)

            
    def copy_src_to_dst_folder(self, src, dst_path):
        filename = os.path.basename(src)
        dst = os.path.join(dst_path, filename)
        shutil.copy(src, dst)


    def copy_other_files_to_etc(self):
        src = self.cfg['features']['feat_params']
        self.copy_src_to_dst_folder(src, self.etc)

        src = self.cfg['testing']['test_fileids']
        self.copy_src_to_dst_folder(src, self.etc)

        src = self.cfg['testing']['test_transcriptions']
        self.copy_src_to_dst_folder(src, self.etc)


    def prepare_training(self):
        #path to current training directory

        #name of the model

        #delete current files in etc/ folder
        print(f"\tDeleting current files at {self.etc}")
        self.delete_files_from_folder(self.etc)
        #copy new training files in etc/ folder
        print(f"\tCopying dictionaries to {self.etc}")
        self.copy_dictionaries_to_etc()

        print(f"\tCopying training files to {self.etc}")
        self.copy_training_files_to_etc()
        #training log files (capturing stdout)
        self.copy_other_files_to_etc()


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
        #call sphinxtrain run
        #command_train=['sphinxtrain','run']
        command_train=['sphinxtrain run']
        cwd = self.base_dir

        print(f"Starting training of model {self.model_name}.")
        self.execute(command_train, cwd, self.log_file)
        
    
    def copytree(self, src, dst, symlinks = False, ignore = None):
        if not os.path.exists(dst):
            os.makedirs(dst)
            shutil.copystat(src, dst)
        lst = os.listdir(src)
        if ignore:
            excl = ignore(src, lst)
            lst = [x for x in lst if x not in excl]
        for item in lst:
            s = os.path.join(src, item)
            d = os.path.join(dst, item)
            if symlinks and os.path.islink(s):
                if os.path.lexists(d):
                    os.remove(d)
                os.symlink(os.readlink(s), d)
                try:
                    st = os.lstat(s)
                    mode = stat.S_IMODE(st.st_mode)
                    os.lchmod(d, mode)
                except:
                    pass # lchmod not available
            elif os.path.isdir(s):
                self.copytree(s, d, symlinks, ignore)
            else:
                shutil.copy2(s, d)


    def copy_model(self):
        print(f"\nModel \"{self.model_name}\" copied from:\n{self.model_location}\n to:\n{os.path.normpath(os.path.join( os.getcwd(), self.model_destination))}\n")
        self.copytree(self.model_location, self.model_destination, symlinks=False, ignore=False)
        print(f"\nTraining files \"{self.model_name}\" copied from:\n{self.etc}\n to:\n{os.path.normpath(os.path.join( os.getcwd(), self.model_etc_destination))}\n")
        self.copytree(self.etc, self.model_etc_destination, symlinks=False, ignore=False)





def main():
    config_folder="./../training_configurations/Bare2"


    training_bare=train(config_folder)

    training_bare.fit()

    training_bare.copy_model()



if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")