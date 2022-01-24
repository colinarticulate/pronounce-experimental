# pronounce-experimental
Here we test new experimental features for Pronounce, keep track of new dictionaries, transcriptions for training, new models, tests, and expectations.  
It is meant to work accross platforms and it should give the same results or at least very similar.  
  
Current version of test_pronounce and cli_pron accepts -hmm as an option.  
  
    
    

##  Requirements, Installation and Execution  
Pronounce is currently written in Go. Therefore, Go must be installed in your system

1. 'sox' must be installed in your system.  http://sox.sourceforge.net/   
        (make sure 'sox' and 'soxi' are available system-wide)
2. 'pocketsphinx' must be installed in your system.  https://github.com/cmusphinx/pocketsphinx  
        (make sure 'pocketsphinx_continuous' and 'pocketsphinx_batch' are available system-wide)
3. cli_pron must be installed in your system (see test_pronounce directory for instructions)
4. test_pronounce must be build to test Pronounce (follow instructions in test_pronounce directory)
5. You must have somewhere in your system a folder with all the audios (.wav) for the test harness required to test Pronounce.  


By default the system is configured to work given the current directory structure. However, the path to the audios for the test harness must be reconfigured depending where your have them located in your system. The following must be changed:  

1. if you use test_pronounce from the cli, you must change the -infolder option accordingly. See the example in README.md within test_pronounce folder.  
2. if you use test_pronunce debugging from vscode, you must change the -infolder option in the 'args' field of the file launch.json.  
3. if you want to use the python tool to clean up temporary files that test_prononunce might have left (due to a crash or stopping execution) you must change 'dir_target' accordingly in the file 'remove_tmp.py':    
```
#Currently set to "~/Data/audio_clips", where '~' represents '/home/user'  
dir_target=os.path.join(home,"Data","audio_clips")  

#Change it to "~/wherever/is/my-audios-folder"  
dir_target=os.path.join(home,"wherever","is","my-audio-folder") 
```  
4. if you want to build test_pronounce with the tag "testCase", you must change the path in testCasing.go line 134:
```  
func testCaseIt(name string, result []result, hmm string) {
	outfolder := "/home/dbarbera/Repositories/test_pronounce/audio_clips/"  
```  

and use go build -tags "noSend debug testCase" to build test_pronounce. This will store the results of Pronounce in the "/test_cases" folder within the audios folder.  If you install cli_pron with the "testCase" folder, executing test_pronounce will also store all the information passed to pocketsphinx during Pronounce execution (logs, pocketsphinx parameters for each call, audios and pocketsphinx's results)  
 

