

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
                        -infolder /Users/test/Downloads/train/SpeechCommands/bed/ \
                        -tests /Users/test/Downloads/train/SpeechCommands/bed/bed_input.csv \
                        -expectations /Users/test/Downloads/train/SpeechCommands/bed/bed_expectations.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-14T15.50.13-019_Bare.ci_cont/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-14T15.50.13-019_Bare.ci_cont

echo
echo "clean up"
rm /Users/test/Downloads/train/SpeechCommands/bed/*_fixed*.wav                        
echo
date
echo
echo


