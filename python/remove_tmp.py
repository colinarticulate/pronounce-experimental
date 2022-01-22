import os 
import shutil

dir_file=os.path.dirname(__file__)
home=os.path.expanduser("~")
dir_target=os.path.join(home,"Data","audio_clips")
dir_target1="../"

wavfiles = [os.remove(os.path.join(dir_target,f)) for f in os.listdir(dir_target) if (f.endswith(".wav") and "_fixed" in f)]

temp_dirs = [shutil.rmtree(os.path.join(dir_target,f)) for f in os.listdir(dir_target) if "Temp_" in f]

ctl_files = [os.remove(os.path.join(dir_target1,f)) for f in os.listdir(dir_target1) if "ctl_" in f]

#print(wavfiles)
#print(temp_dirs)
#print(ctl_files)
print(f"{len(wavfiles)+len(temp_dirs)+len(ctl_files)} files have been removed.")


print("Finished!!")