# test_pronounce  
A test harness to test performance of our current system.  

##  Requirements, Installation and Execution  

test_pronounce requires cli_pron to be installed in your system.  
  
  1. Clone the cli_pron repo:  
  ``` $ git clone  https://github.com/colinarticulate/cli_pron.git ```  
  2. Install cli_pron in your system:  
  ``` $ cd cli_pron ```  
  ``` $ go install -tags "noSend debug" ```  
  3. Clone test_pronounce repo in a different place:  
  ``` $ git clone https://github.com/colinarticulate/test_pronounce.git ```  
  4. Build test_pronounce:  
  ``` $ cd test_pronounce ```  
  ``` $ go build -tags "noSend debug" ```  
  5. Install the audio files in a folder of your choice together with the expectations files and the test inputs files.  
  6. pocketsphinx needs to be installed in your system and the following programs available system-wide:  
        - pocketsphinx_continuous  
        - pocketsphinx_batch  
  7. Update the MODELDIR folder (/usr/loca/share/pocketsphinx/model) with Articulate's current configuration in use.  
  8. Locate the -dict and -phdict files in your system with Articulate's current versions in use.  
  9. Execute test_pronounce with all the options from the command line:  
  ``` $ ./test_pronounce -dict <dict-file>.phone -phdict <phdict-file>.phone -infolder <folder-containing-the-audios-expectations-and-input-test-files> -test <file-with-the-test-inputs(inside infolder)>.csv -outfolder <folder-to-output-the-results-of-the-test>  -featparams /usr/local/share/pocketsphinx/model/en-us/en-us/feat.params  ```  
  10. Alternatively, make your own srcript invoking the above command line in point 9. You can time your code with 'time'.  
  11. Results are a list of .txt files in the output folder (one per audio file tested) together with  the 'summary.txt' file in which you can check the final performance accuracy.  




