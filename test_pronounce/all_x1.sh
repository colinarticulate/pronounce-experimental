

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
                        -infolder /Users/test/go/src/Phonemes/audio_clips/ \
                        -tests /Users/test/go/src/Phonemes/audio_clips/pronounce_input.csv \
                        -expectations /Users/test/go/src/Phonemes/audio_clips/expectations.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-07T14.11.46-092_Bare_with_UWs.ci_cont/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-07T14.11.46-092_Bare_with_UWs.ci_cont
                        
echo
date
echo
echo


