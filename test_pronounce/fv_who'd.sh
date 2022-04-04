

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
                        -infolder /Users/test/Downloads/train/us_female_vowels/who\'d_uw/who\'d \
                        -tests /Users/test/Downloads/train/us_female_vowels/who\'d_uw/who\'d/fv_who\'d_input.csv \
                        -expectations /Users/test/Downloads/train/us_female_vowels/who\'d_uw/who\'d/fv_who\'d_expectations.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_23Mar22_v5/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_23Mar22_v5

echo
echo "clean up"
rm /Users/test/Downloads/train/us_female_vowels/who\'d_uw/who\'d/*_fixed*.wav                        
echo
date
echo
echo


