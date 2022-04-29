package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	pathpkg "path"
	"time"

	"github.com/colinarticulate/pron"
	"github.com/google/uuid"
)

type pronError int

const (
	badCall pronError = iota
	missingWord
)

var pronErrors = []string{
	"File missing. Usage: ./pronounce -audiofolder folder -ctlfile filename -featparams filename -hmm hmm_folder -word word -dict filename -phdict filename\n",
	"Word missing",
}

func (p pronError) Error() string {
	return pronErrors[p]
}

var (
	audiofile, featureparams, hmm, word, dictfile, phdictfile string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	err := guard()
	if err != nil {
		returnJSON(pron.ToJSON([]pron.LettersVerdict{}, err))
		return
	}

	outfolder := pathpkg.Dir(audiofile)
	// test_folder := "test_output"
	// outfolder := pathpkg.Clean(pathpkg.Join(pathpkg.Dir(audiofile), "/.."+test_folder))

	outfolder = pathpkg.Join(outfolder, "Temp_"+uuid.New().String())
	err = os.Mkdir(outfolder, 0777)
	if err != nil {

	}
	start := time.Now()
	results, err := pron.Pronounce(outfolder, audiofile, word, dictfile, phdictfile, featureparams, hmm)
	elapsed := time.Since(start)
	fmt.Printf("Total: %s\n", elapsed)
	result := pron.ToJSON(results, err)
	returnJSON(result)
	//testCaseIt(outfolder, audiofile, word, dictfile, phdictfile, featureparams, string(result))
}

func init() {
	flag.StringVar(&audiofile, "audio", "", "The URL of the audio file.")
	flag.StringVar(&featureparams, "featparams", "", "File containing feature extraction parameters.")
	flag.StringVar(&hmm, "hmm", "", "Directory containing acoustic model files.")
	flag.StringVar(&word, "word", "", "The word to be checked for pronunciation")
	flag.StringVar(&dictfile, "dict", "", "The dictionary to be used")
	flag.StringVar(&phdictfile, "phdict", "", "The phonemes.dict file to be used")
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

// func audiofileName(audiofolder, ctlfile string) (string, error) {
// 	audiofile, err := ioutil.ReadFile(ctlfile)
// 	if err != nil {
// 		return "", err
// 	}
// 	return pathpkg.Join(audiofolder, string(audiofile)+".wav"), nil
// }
