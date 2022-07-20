package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	pathpkg "path"

	"github.com/colinarticulate/pron"
	"github.com/google/uuid"
)

type pronError int

const (
	badCall pronError = iota
	missingWord
)

var pronErrors = []string{
	"File missing. Usage: ./pronounce -audiofolder folder -ctlfile filename -featparams filename -word word -dict filename -phdict filename\n",
	"Word missing",
}

func (p pronError) Error() string {
	return pronErrors[p]
}

var (
	audiofile, featureparams, word, dictfile, phdictfile, hmm string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			returnJSON(
				pron.ToJSON(
					word,
					[]pron.LettersVerdict{},
					errors.New("pron panicked, now recovered"),
				),
			)
		}
	}()

	err := guard()
	if err != nil {
		returnJSON(pron.ToJSON(word, []pron.LettersVerdict{}, err))
		return
	}

	outfolder := pathpkg.Dir(audiofile)

	outfolder = pathpkg.Join(outfolder, "Temp_"+uuid.New().String())
	err = os.Mkdir(outfolder, 0777)
	if err != nil {

	}

	fmt.Println("Calling Pronounce with args: outfolder", outfolder, "audiofile", audiofile, "word", word, "dictfile", dictfile, "phdictfile", phdictfile, "featureparams", featureparams, "hmm", hmm)
	results, err := pron.Pronounce(outfolder, audiofile, word, dictfile, phdictfile, featureparams, hmm)
	returnJSON(pron.ToJSON(word, results, err))
}

func init() {
	flag.StringVar(&audiofile, "audio", "", "The URL of the audio file.")
	flag.StringVar(&featureparams, "featparams", "", "The URL of the modified feat.params file")
	flag.StringVar(&word, "word", "", "The word to be checked for pronunciation")
	flag.StringVar(&dictfile, "dict", "", "The dictionary to be used")
	flag.StringVar(&phdictfile, "phdict", "", "The phonemes.dict file to be used")
	flag.StringVar(&hmm, "hmm", "", "The URL of the model files")
}

func guard() error {
	paths := []string{audiofile, dictfile, phdictfile, featureparams}
	for _, path := range paths {
		if path == "" {
			return badCall
		}
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return err
		}
	}
	if word == "" {
		return missingWord
	}
	return nil
}

func returnJSON(json []byte) {
	fmt.Println(string(json))
}
