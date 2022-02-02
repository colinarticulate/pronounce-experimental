'''
Experiment builder to test how data augmentation might help improve accuracy
'''

import os
import shutil
import subprocess
import time
from sklearn import model_selection
import toml
import stat
import numpy as np

from exec import execute


class Experiment():

    def __init__(self,  original_model_name, transforms_parameters, augmentation_sets, 
                        wav_dir, audios_dir, 
                        fileids_file, transcription_file, 
                        training_configurations_dir, training_template_file, 
                        testing_configurations_dir, testing_template_file,
                        test_set_folder, test_set_name,
                        experiments_folder, experiment_name):

        self.experiments_folder = experiments_folder
        self.experiment_name = experiment_name

        self.training_configurations_dir = training_configurations_dir
        self.training_template_file = training_template_file
        self.testing_template_file = testing_template_file
        self.testing_configurations_dir  = testing_configurations_dir

        self.test_set_folder = test_set_folder
        self.test_set_name = test_set_name

        self.original_model_name=original_model_name

        self.transforms=list(transforms_parameters.keys())
        self.parameters=transforms_parameters
        self.augmentation_sets=augmentation_sets

        self.wav_dir = wav_dir # the base dir where our audios are stored
        self.audios_dir = audios_dir # where our audios are store relative to wav_dir

        #self.original_audios = [f for f in os.listdir(audios_dir)]
        with open(fileids_file,'r') as f:
            raw=f.read()
            self.original_fileids = raw.strip("\n").split("\n")
        with open(transcription_file,'r') as f:
            raw=f.read()
            self.original_transcription=raw.strip("\n").split("\n")

        


    def generate_transform_file_extension(self, transform, p):

        initial=transform[0]
        if (p+1)/int(p+1)>1:
            coded_integer=f"{int(100*abs(p))}".zfill(3)
        else:
            coded_integer = f"p{abs(int(p))}" if p >=0 else f"n{abs(int(p))}" # p for positive and n for negative
                                
        extension=f"{initial}{coded_integer}"
        return extension


    def sox_call(self, src, dst, pipeline):
        
        commands=[f"sox {src} -p"]
        for pipe in pipeline:
            transform = pipe[0]
            p=pipe[1]
            commands.append(f"sox - -p {transform} {p}")
        commands.append(f"sox - {dst}")

        command_piped = " | ".join(commands)
               
        cws=os.getcwd()
        execute(command_piped, cws)


    def process_audio(self, src, dst, pipeline):
        if not os.path.exists(dst):
            self.sox_call(src, dst, pipeline)


    def prepare_augmentation(self, transform_set):
        transforms = transform_set.split("x")

        transformation_sets=[]
        #test_sets=[]
        for transform in transforms:
            transformations=[]
            #set_i=[]
            for i,p in enumerate(self.parameters[transform]):
                transformations.append(f"{transform}_{p}")
                #transformations.append([transform,p])
                #set_i.append(i)
            transformation_sets.append(transformations)
            #test_sets.append(np.array(set_i))


        transformations_per_file = np.array(np.meshgrid(*tuple(transformation_sets))).T.reshape(-1,len(transforms))

        transformation_pipelines=[]
        file_extensions=[]

        for transforms in transformations_per_file:
            pipeline=[]
            file_extension=["aug"]
            for transformation in transforms:
                parts=transformation.split("_")
                transform=parts[0]
                parameter=parts[1]

                extension = self.generate_transform_file_extension(transform, float(parameter))
                pipe = [transform,parameter]

                pipeline.append(pipe)
                file_extension.append(extension)

            transformation_pipelines.append(pipeline)
            file_extensions.append("_"+"_".join(file_extension))


        return transformation_pipelines, file_extensions



    def create_augmented_audios(self, filename, transform_set):
        transformation_pipelines, file_extensions = self.prepare_augmentation(transform_set)
        
        src = os.path.join(self.wav_dir, filename + ".wav")

        new_file_names=[]
        for pipeline, extension in zip(transformation_pipelines, file_extensions):
            new_name = filename + extension
            new_file_names.append(new_name)

            dst = os.path.join(self.wav_dir, new_name + ".wav")
            self.process_audio(src,dst,pipeline)

        return new_file_names
        


            

    def create_data(self, initial_fileids, initial_transcription, transform_sets):
        #Audio: check first if has been already created
        new_fileids=[]
        new_transcriptions=[]
        for transform_set in transform_sets:
            for fileid, transcription in zip(initial_fileids, initial_transcription):
                new_fileids.append(fileid)
                new_transcriptions.append(transcription)
                transformed_filenames = self.create_augmented_audios(fileid, transform_set)
                filename = os.path.basename(fileid)
                for new_fileid in transformed_filenames:
                    new_fileids.append(new_fileid)
                    new_filename = os.path.basename(new_fileid)
                    new_transcription=f"{new_filename}".join(transcription.split(filename))
                    new_transcriptions.append(new_transcription)

        return new_fileids, new_transcriptions


    def add_fileids_transcriptions(self, fileids, transcriptions, etc_dir, filename_fileids, filename_transcriptions):
        with open(os.path.join(etc_dir,filename_fileids),'w') as f:
            f.write("\n".join(fileids)+"\n")
        with open(os.path.join(etc_dir,filename_transcriptions),'w') as f:
            f.write("\n".join(transcriptions)+"\n")


    def add_training_configuration(self, cfg_file, model_name):
        with open(self.training_template_file, 'r') as f:
            contents=f.read()

        new_contents=f"{model_name}".join(contents.split("$(__MODEL_NAME__)"))

        with open(cfg_file, 'w') as f:
            f.write(new_contents)



    def prepare_training(self, fileids, transcriptions, transform_sets):
        #decide model_name
        model_name = self.original_model_name + "_" + "_".join(transform_sets)
        model_dir=os.path.join(self.training_configurations_dir, model_name)
        etc_dir = os.path.join(model_dir, "etc")
        #create folder
        if not os.path.exists(model_dir):
            os.mkdir(model_dir)

        if not os.path.exists(etc_dir):
            os.mkdir(etc_dir)

        self.add_fileids_transcriptions(fileids, transcriptions, etc_dir, f"art_db_{model_name}_train.fileids", f"art_db_{model_name}_train.transcription")
        self.add_training_configuration(os.path.join(model_dir,"configuration.toml"), model_name)
                
        return model_name 

    
    def prepare_model_testing(self, model_name, test_set_folder, test_set_name):
        with open(self.testing_template_file,'r') as f:
            contents = f.read()

            testing_output_folder = f"./../Test_Output/output_{model_name}"

            template_keys = {
                '-dict' : ["__DICT__", "./../Dictionaries/art_db_v2.dic"],
                '-phdict' : ["__PHDICT__", "./../Dictionaries/art_db_v2_inference.phone"],
                '-infolder' : ["__AUDIOS_FOLDER__", test_set_folder], #Test_harness 
                '-tests' : ["__INPUTS__", "./../Tests/pronounce_input.csv"],
                '-expectations' : ["__EXPECTATIONS__", "./../Expectations/expectations_v2.csv"],
                '-outfolder' : ["__OUTPUT_FOLDER__", testing_output_folder],
                '-featparams' : ["__FEAT_PARAMS__", f"./../Models/{model_name}.ci_cont/feat.params"],
                '-hmm' : ["__HMM__", f"./../Models/{model_name}.ci_cont"]
            }

            if not os.path.exists(testing_output_folder):
                os.mkdir(testing_output_folder)

        
            for k in template_keys:
                contents=f"{template_keys[k][1]}".join(contents.split(f"{template_keys[k][0]}"))

            with open(os.path.join(self.testing_configurations_dir, f"{model_name}_x_{test_set_name}.toml"), 'w') as f:
                f.write(contents)

  
    def create_experiment_configuration(self):
        with open(os.path.join(self.experiments_folder,f"{self.experiment_name}.toml"), 'w') as f:
            models_header="[models]\n"
            datasets_header="[test_sets]\n"

            f.write(models_header)
            if self.model_names!=None or self.model_names!=[]:
                f.write("'models' = [")
                for model_name in self.model_names[:-1]:
                    f.write(f"\"{model_name}\", ")
                f.write(f"\"{self.model_names[-1]}\" ]\n")

            f.write("\n")

            f.write(datasets_header)
            #At the moment, only one dataset (test harness)
            f.write(f"'datasets' = [")
            f.write(f"\"{self.test_set_name}\" ]")


    def create_experiment(self):
        print("Creating augmented audios and corresponding training files for:\n")
        self.augmented_fileids=self.original_fileids
        self.augmented_transcription=self.original_transcription

        self.model_names=[]
        for i,transform_sets in enumerate(self.augmentation_sets[:]):
            model_name=self.original_model_name+"_"+"_".join(transform_sets)
            self.model_names.append(model_name)

            print(f"{i+1}\t\t{model_name}")
            
            new_fileids, new_transcriptions = self.create_data(self.original_fileids, self.original_transcription, transform_sets)
            model_name_set = self.prepare_training(new_fileids, new_transcriptions, transform_sets)

            self.prepare_model_testing(model_name_set, self.test_set_folder, self.test_set_name )

        self.create_experiment_configuration()

        return self.model_names


def main():
    #config_test_folder="./../testing_configurations/data_augmentation_experiment.toml"

    config_folder="./../training_configurations/Bare"
    fileids= os.path.join(config_folder,"etc","art_db_Bare_train.fileids")
    transcription= os.path.join(config_folder,"etc","art_db_Bare_train.transcription")
    wav_dir = "/home/dbarbera/Repositories/art_db/wav"
    audios_dir = "train/art_db_compilation"
    training_configurations_dir = "/home/dbarbera/Repositories/pronounce-experimental/training_configurations"
    training_template_file = "/home/dbarbera/Repositories/pronounce-experimental/templates/model_training_template.toml"
    testing_configurations_dir = "/home/dbarbera/Repositories/pronounce-experimental/testing_configurations"
    testing_template_file = "/home/dbarbera/Repositories/pronounce-experimental/templates/model_testing_template.toml"
    test_set_folder = "/home/dbarbera/Data/audio_clips"
    test_set_name = "Test_Harness"
    experiments_folder = "./../experiments"
    experiment_name = "Testing_Training_Data_Augmentation_no_pitch"

    audios = [ f for f in os.listdir(os.path.join(wav_dir,audios_dir))]
    

    parameters={'pitch':[-200, 200],
                'loudness':[-10],
                'speed':[0.90, 0.93, 0.95, 1.03, 1.05, 1.10],
                'tempo':[1.50, 1.25]
                }

    augmentation_sets=[ ['pitch'],
                        ['loudness'],
                        ['speed'],
                        ['tempo'],
                        ['pitch','loudness'],
                        ['pitch','speed'],
                        ['pitch','tempo'],
                        ['loudness', 'speed'],
                        ['loudness', 'tempo'],
                        ['pitch','loudness','speed'],
                        ['pitch','loudness','tempo'],
                        ['loudness','speed','tempo'],
                        ['pitch','loudness','speed','tempo'],
                        ['speed','tempo'],
                        ['pitchxloudness'],
                        ['pitchxloudness','speed'],
                        ['pitchxloudness','tempo'],
                        ['pitchxloudness','speed','tempo'],
                        ['pitchxloudnessxspeed'],
                        ['pitchxloudnessxspeed','tempo'],
                        ['pitchxloudnessxtempo'],
                        ['pitchxloudnessxtempo','speed'],
                        ['pitch','loudness','speed','tempo','pitchxloudness', 'pitchxloudnessxspeed','pitchxloudnessxtempo']
                        ]

    augmentation_sets=[ 
                        ['loudnessxspeed'],
                        ['loudnessxtempo'],
                        ['speedxtempo'],
                        ['loudnessxspeed','loudnessxtempo'],
                        ['speed','tempo','speedxtempo'],
                        ['speed','loudnessxspeed'],
                        ['tempo','loudnessxtempo'],
                        ['loudness','speed','tempo','speedxtempo','loudnessxspeed', 'loudnessxtempo','loudnessxspeedxtempo']
                        ]

    augmentation=Experiment("Bare", parameters, augmentation_sets, wav_dir, audios_dir, fileids, transcription, 
                                    training_configurations_dir, training_template_file, 
                                    testing_configurations_dir, testing_template_file,
                                    test_set_folder, test_set_name,
                                    experiments_folder, experiment_name)

    model_names = augmentation.create_experiment()

    

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")