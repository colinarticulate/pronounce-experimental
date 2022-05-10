

#!/bin/bash

echo
echo
echo
echo "Using the new test_pronounce from Github"
echo "Start:"
echo
#
#time /Users/test/keep_away_from_go/art/cmd/testPronounce/testPronounce -dict /Users/test/Downloads/Test_stt_training/Out/sourceFiltered_new.dict -phdict /Users/test/Downloads/Test_stt_training/Out/art_db_v3.phone -infolder /Users/test/go/src/Phonemes/audio_clips/ -tests pronounce_input.csv -expectations expectations.csv -outfolder ./testout -featparams /usr/local/share/pocketsphinx/model/en-us/en-us/feat.params

time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
                        -infolder /Users/test/Downloads/2264/pronounce \
                        -tests /Users/test/Downloads/2264/pronounce/2264_input.csv \
                        -expectations /Users/test/Downloads/2264/pronounce/2264_expectation.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                       -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_9May22_v2/model/feat.params \
                       -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_9May22_v2/model
                        
echo
echo "clean up"
rm /Users/test/Downloads/2264/pronounce/*_fixed*.wav                        
echo
date
echo
echo

