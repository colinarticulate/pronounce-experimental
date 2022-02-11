#!/bin/bash

echo
echo
echo
echo " This test Hossein's user experience"
echo " Start:"
echo



time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v2_inference.phone \
                        -infolder /Users/test/Downloads/train/scrape/hossein/HY_test/ \
                        -tests /Users/test/Downloads/train/scrape/hossein/HY_test/0_hy_test_input.csv  \
                        -expectations /Users/test/Downloads/train/scrape/hossein/HY_test/0_HY_expectations.csv  \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-11T11.25.24-056_Bare.ci_cont/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/2022-02-11T11.25.24-056_Bare.ci_cont
                        

#
#time /Users/test/keep_away_from_go/art/cmd/testPronounce/testPronounce 
#						-dict /Users/test/Downloads/temp_accoustic/art_db.dic 
#						-phdict /Users/test/Downloads/temp_accoustic/art_db.phone 
#						-infolder /Users/test/Downloads/train/scrape/hossein/HY_test/ 
#						-tests 0_hy_test_input.csv 
#						-expectations 0_HY_expectations.csv 
#						-outfolder ./testout -featparams /usr/local/share/pocketsphinx/model/en-us/en-us/feat.params

echo "clean up"
rm /Users/test/Downloads/train/scrape/hossein/HY_test/*_fixed*.wav


echo
date
echo
echo
