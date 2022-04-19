#!/bin/bash

echo
echo
echo
echo " This test Judith user experience"
echo " Start:"
echo



time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
                        -infolder /Users/test/Downloads/train/Shtooka/judith/judith/ \
                        -tests /Users/test/Downloads/train/Shtooka/judith/judith/judith_input.csv  \
                        -expectations /Users/test/Downloads/train/Shtooka/judith/judith/judith_expectations.csv  \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_4Apr22/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_4Apr22
                        


echo "clean up"
rm /Users/test/Downloads/train/Shtooka/judith/judith/*_fixed*.wav


echo
date
echo
echo
