package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	dictfile, phdictfile, infolder, tests, expectations, outfolder, featparams, hmm string
)

type testError int

const (
	badCall testError = iota
	parseError
)

var testErrors = []string{
	"File missing. Usage: ./testPronounce -dict filename -phdict filename -infolder filename -tests filename -expectations filename -outfolder foldername -featparams filename -hmm foldername\n",
}

func (p testError) Error() string {
	return testErrors[p]
}

func main() {
	fmt.Println("Starting test pronounce...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	clean(outfolder)
	err := mkDir(outfolder)
	if os.IsPermission(err) {
		fmt.Println("Error making outfolder", err.Error())
		return
	}
	if err = guard(); err != nil {
		fmt.Println("here", err)
		fmt.Println("Error in input parameters", err.Error())
		return
	}
	t0 := time.Now()
	runTests(dictfile, phdictfile, infolder, tests, expectations, outfolder, featparams, hmm)
	t1 := time.Now()
	fmt.Println("Elapsed time =", t1.Sub(t0))
}

func init() {
	flag.StringVar(&dictfile, "dict", "", "The dictionary to be used")
	flag.StringVar(&phdictfile, "phdict", "", "The phonemes dictionary to be used")
	flag.StringVar(&infolder, "infolder", "", "The folder containing audio clips to be tested")
	flag.StringVar(&tests, "tests", "", "The file containing the tests to run")
	flag.StringVar(&expectations, "expectations", "", "The file containing the expected results for the tests")
	flag.StringVar(&outfolder, "outfolder", "", "The folder to write test results out to")
	flag.StringVar(&featparams, "featparams", "", "The URL of the modified feat.params file")
	flag.StringVar(&hmm, "hmm", "", "Folder with acoustic model files.")
}

func guard() error {
	params := []string{dictfile, phdictfile, infolder, tests, expectations, outfolder, featparams, hmm}
	fmt.Println(params)
	for _, param := range params {
		if param == "" {
			fmt.Println("badcall1 -> param: ", param)
			return badCall
		}
	}
	filepaths := []string{dictfile, phdictfile, infolder, tests, expectations, outfolder, featparams, hmm}
	for _, path := range filepaths {
		fmt.Println("parsing ", path)
		if path == "" {
			fmt.Println("badcall2")
			return badCall
		}
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Print("Problems with ", path)
			return err
		}
	}
	// files := []string{tests, expectations}
	// for _, file := range files {
	// 	path := filepath.Join(infolder, file)
	// 	_, err := os.Stat(path)
	// 	if os.IsNotExist(err) {
	// 		return err
	// 	}
	// }
	return nil
}

func mkDir(dirname string) error {
	err := os.Mkdir(dirname, 0777)
	if os.IsPermission(err) {
		return err
	}
	if os.IsExist(err) {
		return err
	}
	return nil
}

func clean(folder string) {
	err := os.RemoveAll(outfolder)
	if err != nil {
		fmt.Println("Doh! Error on removing folder", outfolder, "Error =", err)
	}
}
