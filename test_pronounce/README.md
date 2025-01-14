# test_pronounce  
A test harness to test performance of our current system.  

##  Requirements, Installation and Execution  

test_pronounce requires cli_pron to be installed in your system.  
  
   1. Go to the folder cli_pron and install it in your system:  
  ``` $ cd cli_pron ```  
  ``` $ go install -tags "noSend debug" ```  
  2. Go to the folder test_pronounce and build test_pronounce:  
  ``` $ cd test_pronounce ```  
  ``` $ go build -tags "noSend debug" ```  
  3. Install the audio files in a folder of your choice.  
  4. pocketsphinx needs to be installed in your system and the following programs available system-wide:  
        - pocketsphinx_continuous  
        - pocketsphinx_batch  
  5. Locate the -dict and -phdict files in your system with Articulate's current versions in use in the folder /Dictionaries.  
  6. Execute test_pronounce with all the options from the command line:  
  ``` $ ./test_pronounce -dict <dict-file>.dic -phdict <phdict-file>.phone -infolder <folder-containing-the-audios-expectations-and-input-test-files> -test <file-with-the-test-inputs(inside infolder)>.csv -outfolder <folder-to-output-the-results-of-the-test>  -featparams <your-model-folder>/feat.params -hmm <your-model-folder> ```  
  7. Alternatively, make your own srcript invoking the above command line in point 6. You can time your code with 'time'.  
  8. Results are a list of .txt files in the output folder (one per audio file tested) together with  the 'summary.txt' file in which you can check the final performance accuracy.  

  Example:  

 ```
$ ./test_pronounce     -dict ./../Dictionaries/art_db_v2.dic \  
                        -phdict ./../Dictionaries/art_db_v2_inference.phone \  
                        -infolder /home/Data/audio_clips \
                        -tests ./../Tests/pronounce_inputs.csv \  
                        -expectations ./../Expectations/expectations_v2.csv \  
                        -outfolder ./output_25600  \  
                        -featparams ./../Models/25600.ci_cont/feat.params \  
                        -hmm ./../Models/25600.ci_cont 
``` 




