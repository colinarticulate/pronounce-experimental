

#!/bin/bash

echo
echo
echo
echo "Using the new test_pronounce from Github"
echo "Start:"
echo
#


time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
                        -infolder /Users/test/Downloads/train/SpeechCommands/bird/bird/ \
                        -tests /Users/test/Downloads/train/SpeechCommands/bird/bird/bird_input.csv \
                        -expectations /Users/test/Downloads/train/SpeechCommands/bird/bird/bird_expectations.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3Mar22_v3/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3Mar22_v3

echo
echo "clean up"
rm /Users/test/Downloads/train/SpeechCommands/bird/bird/*_fixed*.wav                        
echo
date
echo
echo


