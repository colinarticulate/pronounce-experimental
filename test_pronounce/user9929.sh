

#!/bin/bash

echo
echo
echo
echo "Using the new test_pronounce from Github"
echo "Start:"
echo


time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
 						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
                        -infolder /Users/test/Downloads/user9929 \
                        -tests /Users/test/Downloads/user9929/user9929_input.csv \
                        -expectations /Users/test/Downloads/user9929/user9929_expectations1.csv \
                        -outfolder /Users/test/Downloads/audio_temp \
                        -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_13Apr22_v2/feat.params \
                        -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_13Apr22_v2


#time ./test_pronounce -dict /Users/test/Downloads/pronunce_production/dictionary/art_db.dic \
# 						-phdict /Users/test/Downloads/pronunce_production/dictionary/art_db.phone \
#                        -infolder /Users/test/Downloads/user9929 \
#                        -tests /Users/test/Downloads/user9929/user9929_input.csv \
#                        -expectations /Users/test/Downloads/user9929/user9929_expectations.csv \
#                        -outfolder /Users/test/Downloads/audio_temp \
#                        -featparams /Users/test/Downloads/pronunce_production/acoustic_model/feat.params \
#                        -hmm /Users/test/Downloads/pronunce_production/acoustic_model


echo
echo "clean up"
rm /Users/test/Downloads/user9929/*_fixed*.wav                        
echo
date
echo
echo
