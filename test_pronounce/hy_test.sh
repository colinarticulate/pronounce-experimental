#!/bin/bash

echo
echo
echo
echo " This test Hossein's user experience"
echo " Start:"
echo



#time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
# 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
#                        -infolder /Users/test/Downloads/train/scrape/hossein/HY_test/ \
#                        -tests /Users/test/Downloads/train/scrape/hossein/HY_test/0_hy_test_input.csv  \
#                        -expectations /Users/test/Downloads/train/scrape/hossein/HY_test/0_HY_expectations.csv  \
#                        -outfolder /Users/test/Downloads/audio_temp \
#                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_13Apr22_v2/feat.params \
#                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_13Apr22_v2

time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22_mini.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_64only_inference.phone \
                        -infolder /Users/test/Downloads/train/scrape/hossein/HY_test/ \
                        -tests /Users/test/Downloads/train/scrape/hossein/HY_test/0_hy_test_input.csv  \
                        -expectations /Users/test/Downloads/train/scrape/hossein/HY_test/0_HY_expectations.csv  \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_20Apr22/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_20Apr22


#time ./test_pronounce -dict /Users/test/Downloads/pronunce_production/dictionary/art_db.dic \
# 						-phdict /Users/test/Downloads/pronunce_production/dictionary/art_db.phone \
#                        -infolder /Users/test/Downloads/train/scrape/hossein/HY_test/ \
#                        -tests /Users/test/Downloads/train/scrape/hossein/HY_test/0_hy_test_input.csv  \
#                        -expectations /Users/test/Downloads/train/scrape/hossein/HY_test/0_HY_expectations.csv  \
#                        -outfolder /Users/test/Downloads/audio_temp \
#                        -featparams /Users/test/Downloads/pronunce_production/acoustic_model/feat.params \
#                        -hmm /Users/test/Downloads/pronunce_production/acoustic_model

echo "clean up"
rm /Users/test/Downloads/train/scrape/hossein/HY_test/*_fixed*.wav


echo
date
echo
echo
