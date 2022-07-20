package main

import (
	"os"
	pathpkg "path"
	"testing"

	"github.com/colinarticulate/pron"
	"github.com/google/uuid"
)

func TestPron(t *testing.T) {
	err := guard()
	if err != nil {
		returnJSON(pron.ToJSON([]pron.LettersVerdict{}, err))
		return
	}

	audiofile := "/home/dbarbera/Data/audio_clips/abandoned1_paul.wav"
	outfolder := pathpkg.Dir(audiofile)
	// test_folder := "test_output"
	// outfolder := pathpkg.Clean(pathpkg.Join(pathpkg.Dir(audiofile), "/.."+test_folder))

	outfolder = pathpkg.Join(outfolder, "Temp_"+uuid.New().String())
	err = os.Mkdir(outfolder, 0777)
	if err != nil {

	}
	word := "abandoned"
	dictfile := "/home/dbarbera/Repositories/articulate-pocketsphinx-go/caller_plus/test_data/dictionaries/art_db.dic"
	phdictfile := "/home/dbarbera/Repositories/articulate-pocketsphinx-go/caller_plus/test_data/dictionaries/art_db.phone"
	featureparams := "/home/dbarbera/Repositories/articulate-pocketsphinx-go/caller_plus/test_data/model/art-en-us/en-us/feat.params"
	hmm := "/home/dbarbera/Repositories/articulate-pocketsphinx-go/caller_plus/test_data/model/art-en-us/en-us"
	results, err := pron.Pronounce(outfolder, audiofile, word, dictfile, phdictfile, featureparams, hmm)

}
