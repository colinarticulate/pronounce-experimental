# A Python tool to clean up test_pronounce temporary files  

Test_pronounce will write all sort of temporary folders and files within the infolder directory where the audios for testing are located.  

This is just a tool to clean up all the temporary files test_pronounce leaves behind when interrupting its execution.

##  Usage  
  
From this folder open terminal and type:  
  ``` $ python remove_tmp.py ```  
  

You can change the directory were your audio_clips are located in the script.  
The line below is equivalent to the path "/home/you-user/Data/audio_clips"  

``` dir_target=os.path.join(home,"Data","audio_clips") ```

