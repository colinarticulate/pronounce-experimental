//go:build testCase
// +build testCase

package main

import (
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(">>>>>> cli_pron !!!", e)
		panic(e)
	}
}

// exists returns whether the given file or directory exists
func create_no_overwrite(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err_mk := os.Mkdir(path, 0777)
		check(err_mk)
	}

}

func test_case_name_from_audiofile(audiofile string) string {
	parts := strings.Split(audiofile, ".wav") //what cli_pron sends to pron
	test_case_name := parts[0]

	return test_case_name
}

func create_test_case_directories_cli_pron(outfolder string, audiofile string, word string) (string, string, string) {

	//out_folder := pathpkg.Clean(pathpkg.Join(pathpkg.Dir(audiofile), "/.."))
	out_folder := pathpkg.Dir(audiofile)

	tc_folder := pathpkg.Join(out_folder, "test_cases")

	create_no_overwrite(tc_folder)

	audio_file := pathpkg.Base(audiofile)

	test_case_name := test_case_name_from_audiofile(audio_file) // parts[0] + "_" + strings.Split(parts[1], ".")[0]
	test_case_folder := pathpkg.Join(tc_folder, test_case_name+"_"+word)
	//fmt.Println(test_case_folder)
	create_no_overwrite(test_case_folder)

	// tmp_folder := pathpkg.Join(test_case_folder, pathpkg.Base(outfolder))
	// create_no_overwrite(tmp_folder)

	debug_folder := pathpkg.Join(test_case_folder, "debug")
	create_no_overwrite(debug_folder)

	cli_folder := pathpkg.Join(debug_folder, "cli_pron")
	create_no_overwrite(cli_folder)

	rel_data_folder := pathpkg.Join("test_cases", test_case_name) //, pathpkg.Base(outfolder))

	return cli_folder, test_case_folder, rel_data_folder
}

func create_param_file(file string, folder string, audiofolder string, audiofile string, word string, dictfile string, phdictfile string, featureparams string) {

	f, err := os.Create(file)
	check(err)
	defer f.Close()

	_, errw := f.WriteString("outfolder" + " " + folder + "\n")
	check(errw)

	_, audiofilename := pathpkg.Split(audiofile)
	_, errw = f.WriteString("audiofile" + " " + pathpkg.Join(audiofolder, audiofilename) + "\n")
	check(errw)

	_, errw = f.WriteString("word" + " " + word + "\n")
	check(errw)
	_, errw = f.WriteString("dictfile" + " " + dictfile + "\n")
	check(errw)
	_, errw = f.WriteString("phdictfile" + " " + phdictfile + "\n")
	check(errw)
	_, errw = f.WriteString("featureparams" + " " + featureparams + "\n")
	check(errw)

}

func create_result_file(file string, result string) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	_, errw := f.WriteString(result)
	check(errw)
}

func copy_file(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if !os.IsNotExist(err) {
		if err != nil {
			return 0, err
		}

		if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf(" >>>> cli_pron: %s is not a regular file", src)
		}

		source, err := os.Open(src)
		if err != nil {
			return 0, err
		}
		defer source.Close()

		destination, err := os.Create(dst)
		if err != nil {
			return 0, err
		}
		defer destination.Close()
		nBytes, err := io.Copy(destination, source)
		return nBytes, err
	}
	return 0, err
}

func copy_audio(test_case_folder string, audiofile string, rel_data_folder string) string {

	audio_file := pathpkg.Base(audiofile)
	audio_folder := pathpkg.Join(test_case_folder, "audio")
	create_no_overwrite(audio_folder)
	audio_path := pathpkg.Join(audio_folder, audio_file)

	//Saving data for test case
	_, err := copy_file(audiofile, audio_path)
	check(err)

	audiofolder := pathpkg.Join(rel_data_folder, "audio")

	return audiofolder

}

func testCaseIt(outfolder string, audiofile string, word string, dictfile string, phdictfile string, featureparams string, result string) {
	cli_folder, test_case_folder, rel_data_folder := create_test_case_directories_cli_pron(outfolder, audiofile, word)

	rel_audio_folder := copy_audio(test_case_folder, audiofile, rel_data_folder)

	params_file := pathpkg.Join(cli_folder, "params.txt")

	create_param_file(params_file, rel_data_folder, rel_audio_folder, audiofile, word, dictfile, phdictfile, featureparams)

	result_file := pathpkg.Join(cli_folder, "result.txt")
	create_result_file(result_file, result)

}
