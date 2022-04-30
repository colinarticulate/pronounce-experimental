

#!/bin/bash

echo
echo
echo
echo "Using the new test_pronounce from Github"
echo "Start:"
echo


#time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
# 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
#                        -infolder /Users/test/Downloads/2264/pronounce \
#                        -tests /Users/test/Downloads/2264/pronounce/2264_input.csv \
#                        -expectations /Users/test/Downloads/2264/pronounce/2264_expectation.csv \
#                        -outfolder /Users/test/Downloads/audio_temp \
#                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_14Apr22/feat.params \
#                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_14Apr22

time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22_mini.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_64only_inference.phone \
                        -infolder /Users/test/Downloads/2264/pronounce \
                        -tests /Users/test/Downloads/2264/pronounce/2264_input.csv \
                        -expectations /Users/test/Downloads/2264/pronounce/2264_expectation.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_20Apr22/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_20Apr22


#time ./test_pronounce -dict /Users/test/Downloads/pronunce_production/dictionary/art_db.dic \
# 						-phdict /Users/test/Downloads/pronunce_production/dictionary/art_db.phone \
#                        -infolder /Users/test/Downloads/2264/pronounce \
#                        -tests /Users/test/Downloads/2264/pronounce/2264_input.csv \
#                        -expectations /Users/test/Downloads/2264/pronounce/2264_expectation.csv \
#                        -outfolder /Users/test/Downloads/audio_temp \
#                        -featparams /Users/test/Downloads/pronunce_production/acoustic_model/feat.params \
#                        -hmm /Users/test/Downloads/pronunce_production/acoustic_model


echo
echo "clean up"
rm /Users/test/Downloads/2264/pronounce/*_fixed*.wav                        
echo
date
echo
echo

