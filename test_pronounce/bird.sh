

#!/bin/bash

echo
echo
echo
echo "Using the new test_pronounce from Github"
echo "Start:"
echo
#


time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v2_inference.phone \
                        -infolder /Users/test/Downloads/train/SpeechCommands/bird/ \
                        -tests /Users/test/Downloads/train/SpeechCommands/bird/bird_input.csv \
                        -expectations /Users/test/Downloads/train/SpeechCommands/bird/bird_expectations.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-14T15.50.13-019_Bare.ci_cont/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-14T15.50.13-019_Bare.ci_cont

echo
echo "clean up"
rm /Users/test/Downloads/train/SpeechCommands/bird/*_fixed*.wav                        
echo
date
echo
echo


