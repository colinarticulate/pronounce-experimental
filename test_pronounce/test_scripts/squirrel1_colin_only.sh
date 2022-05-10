

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
                       -infolder /Users/test/go/src/Phonemes/audio_clips/ \
                       -tests /Users/test/go/src/Phonemes/audio_clips/squirrel1_colin_only.csv \
                       -expectations /Users/test/go/src/Phonemes/audio_clips/expectations.csv \
                       -outfolder /Users/test/Downloads/audio_temp \
                       -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3May22/model/feat.params \
                       -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3May22/model

# time ./test_pronounce -dict /Users/test/Downloads/sourceFiltered_exp3Mar22.dict \
#  						-phdict /Users/test/Documents/GitHub/pronounce-experimental/Dictionaries/art_db_v3_inference.phone \
#                         -infolder /Users/test/go/src/Phonemes/audio_clips/ \
#                         -tests /Users/test/go/src/Phonemes/audio_clips/i1_paul_only_x100.csv \
#                         -expectations /Users/test/go/src/Phonemes/audio_clips/expectations.csv \
#                         -outfolder /Users/test/Downloads/audio_temp \
#                         -featparams /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3May22/model/feat.params \
#                         -hmm /Users/test/Documents/GitHub/pronounce-experimental/Models/Bare.ci_cont_3May22/model

echo
echo "clean up"
rm /Users/test/go/src/Phonemes/audio_clips/*_fixed*.wav                        
echo
date
echo
echo
