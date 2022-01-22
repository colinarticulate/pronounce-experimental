//go:build testCase
// +build testCase

package scanScheduler

import (
	"fmt"
	"io"
	//"log"
	"os"
	pathpkg "path"
	"sort"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(">>>>>> scanScheduler !!!!  ", e)
		//log.Fatal(e)
		//panic(e)
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
	parts := strings.Split(audiofile, "_fixed")
	test_case_name := parts[0]

	return test_case_name
}

func test_case_name_from_batch(ctl_file string) string {
	audiofile := strings.Split(ctl_file, "ctl_")[1]
	parts := strings.Split(audiofile, "_fixed")
	test_case_name := parts[0]

	return test_case_name
}

func create_test_case_folder_from_batch(param_values map[string]string, word string) string {
	outfolder := pathpkg.Clean(pathpkg.Join(pathpkg.Dir(param_values["-ctl"]), "/.."))

	tc_folder := pathpkg.Join(outfolder, "test_cases")

	create_no_overwrite(tc_folder)
	ctl_file := pathpkg.Base(param_values["-ctl"])

	test_case_name := test_case_name_from_batch(ctl_file)

	test_case_folder := pathpkg.Join(tc_folder, test_case_name+"_"+word)

	create_no_overwrite(test_case_folder)

	return test_case_folder
}

func create_test_case_folder(param_values map[string]string, word string) string {
	outfolder := pathpkg.Dir(param_values["-infile"])

	tc_folder := pathpkg.Join(outfolder, "test_cases")

	create_no_overwrite(tc_folder)
	audio_file := pathpkg.Base(param_values["-infile"])

	test_case_name := test_case_name_from_audiofile(audio_file)

	test_case_folder := pathpkg.Join(tc_folder, test_case_name+"_"+word)

	create_no_overwrite(test_case_folder)

	return test_case_folder
}

func copy_file(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if !os.IsNotExist(err) {
		if err != nil {
			return 0, err
		}

		if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf(">>>>> scanScheduler: %s is not a regular file", src)
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

func create_param_file(file string, folder string, params map[string]string, jsgf_file string, log_file string, audio_file string) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	//Order the dictionary by key to make it easier to inspect visually
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//Rearranging pathfiles (jsgf, log and audio) to new place within the test cases folder
	for _, k := range keys {
		switch k {
		case "-jsgf":
			_, errw := f.WriteString(k + " " + pathpkg.Join(folder, jsgf_file) + "\n")
			check(errw)
		case "-logfn":
			_, errw := f.WriteString(k + " " + pathpkg.Join(folder, log_file) + "\n")
			check(errw)
		case "-infile":
			_, errw := f.WriteString(k + " " + pathpkg.Join(folder, audio_file) + "\n")
			check(errw)
		default:
			_, errw := f.WriteString(k + " " + params[k] + "\n")
			check(errw)
		}
	}

}

func create_param_file_batch(file string, folder string, params map[string]string, log_file string, ctl_file string, cep_folder string) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	//Order the dictionary by key to make it easier to inspect visually
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//Rearranging pathfiles (jsgf, log and audio) to new place within the test cases folder
	for _, k := range keys {
		switch k {
		case "-logfn":
			_, errw := f.WriteString(k + " " + pathpkg.Join(folder, log_file) + "\n")
			check(errw)
		case "-ctl":
			_, errw := f.WriteString(k + " " + pathpkg.Join(folder, ctl_file) + "\n")
			check(errw)
		case "-cepdir":
			_, errw := f.WriteString(k + " " + cep_folder + "\n")
			check(errw)
		default:
			_, errw := f.WriteString(k + " " + params[k] + "\n")
			check(errw)
		}
	}
}

func copy_files(test_case_folder string, frate string, param_values map[string]string) {

	//src_tmp_folder := pathpkg.Dir(param_values["-jsgf"])

	//filenames without path
	jsgf_file := pathpkg.Base(param_values["-jsgf"])
	log_file := pathpkg.Base(param_values["-logfn"])
	audio_file := pathpkg.Base(param_values["-infile"])

	//folder to keep debug data to be used from the tmp folder
	debug_folder := pathpkg.Join(test_case_folder, "debug")
	create_no_overwrite(debug_folder)

	scanScheduler_folder := pathpkg.Join(debug_folder, "scanScheduler")
	create_no_overwrite(scanScheduler_folder)

	ps_folder := pathpkg.Join(scanScheduler_folder, "continuous")
	create_no_overwrite(ps_folder)

	grammar_folder := pathpkg.Join(test_case_folder, "grammar")
	log_folder := pathpkg.Join(test_case_folder, "log")
	audio_folder := pathpkg.Join(test_case_folder, "audio")
	create_no_overwrite(grammar_folder)
	create_no_overwrite(log_folder)
	create_no_overwrite(audio_folder)

	// jsgf_debug_path := pathpkg.Join(grammar_folder, jsgf_file)
	log_debug_path := pathpkg.Join(log_folder, log_file)
	// audio_debug_path := pathpkg.Join(audio_folder, audio_file)

	//Debug info:
	variant_id := strings.Split(strings.Split(string(pathpkg.Base(param_values["-jsgf"])), "-")[0], "_")[2] //Not relevant to identify variants, just to distinguish ps_calls

	_, log_file = pathpkg.Split(param_values["-logfn"])
	log_file = strings.Split(log_file, ".")[0] + "_" + variant_id + "_frate_" + frate + "_.log"

	jsgf_debug_path := pathpkg.Join(grammar_folder, jsgf_file)
	log_debug_path = pathpkg.Join(log_folder, log_file)
	audio_debug_path := pathpkg.Join(audio_folder, audio_file)

	//Saving data for test case
	_, err := copy_file(param_values["-jsgf"], jsgf_debug_path)
	check(err)
	_, err = copy_file(param_values["-logfn"], log_debug_path)
	check(err)
	_, err = copy_file(param_values["-infile"], audio_debug_path)
	check(err)

	rel_data_folder := pathpkg.Join("test_cases", pathpkg.Base(test_case_folder))

	//Debug info:
	// variant_id := strings.Split(strings.Split(string(pathpkg.Base(param_values["-jsgf"])), "-")[0], "_")[2] //Not relevant to identify variants, just to distinguish ps_calls
	// fmt.Println("Variant id: ")
	// fmt.Println(variant_id)
	ps_call_folder := pathpkg.Join(ps_folder, "ps_call-"+variant_id+"_frate_"+frate)
	create_no_overwrite(ps_call_folder)

	//files with debug info:
	ps_args_file := pathpkg.Join(ps_call_folder, "ps_args.txt")
	//result_file := pathpkg.Join(ps_call_folder, "result.txt")

	log_file = pathpkg.Join("log", log_file)
	jsgf_file = pathpkg.Join("grammar", jsgf_file)
	audio_file = pathpkg.Join("audio", audio_file)

	create_param_file(ps_args_file, rel_data_folder, param_values, jsgf_file, log_file, audio_file)

}

func copy_files_batch(test_case_folder string, frate string, param_values map[string]string, logfn string, cmnvec []string) {

	//filenames without path
	//jsgf_file := pathpkg.Base(param_values["-jsgf"])
	//log_file := pathpkg.Base(param_values["-logfn"])
	log_file := pathpkg.Base(logfn)
	//log_file_ := pathpkg.Base(param_values["-logfn"])
	ctl_file := pathpkg.Base(param_values["-ctl"])

	//folder to keep debug data to be used from the tmp folder
	debug_folder := pathpkg.Join(test_case_folder, "debug")
	create_no_overwrite(debug_folder)

	scanScheduler_folder := pathpkg.Join(debug_folder, "scanScheduler")
	create_no_overwrite(scanScheduler_folder)

	ps_folder := pathpkg.Join(scanScheduler_folder, "batch")
	create_no_overwrite(ps_folder)

	//grammar_folder := pathpkg.Join(test_case_folder, "grammar")
	log_folder := pathpkg.Join(test_case_folder, "log")
	ctl_folder := pathpkg.Join(test_case_folder, "ctl")

	//create_no_overwrite(grammar_folder)
	create_no_overwrite(log_folder)
	create_no_overwrite(ctl_folder)

	// jsgf_debug_path := pathpkg.Join(grammar_folder, jsgf_file)
	//log_debug_path := pathpkg.Join(log_folder, log_file)
	// audio_debug_path := pathpkg.Join(audio_folder, audio_file)

	//Debug info:
	// Can't access variant id info when batch scan: it doesn't use -jsgf grammar files.
	// variant_id := strings.Split(strings.Split(string(pathpkg.Base(param_values["-jsgf"])), "-")[0], "_")[2] //Not relevant to identify variants, just to distinguish ps_calls
	// fmt.Println("Variant id: ")
	// fmt.Println(variant_id)

	log_file = pathpkg.Base(logfn)
	log_file = strings.Split(log_file, ".")[0] + "_BATCH_.log"

	//jsgf_debug_path := pathpkg.Join(grammar_folder, jsgf_file)
	log_debug_path := pathpkg.Join(log_folder, log_file)
	ctl_debug_path := pathpkg.Join(ctl_folder, ctl_file)

	//Saving data for test case
	//_, err := copy_file(param_values["-jsgf"], jsgf_debug_path)
	//check(err)
	_, err := copy_file(logfn, log_debug_path)
	check(err)
	_, err = copy_file(param_values["-ctl"], ctl_debug_path)
	check(err)

	rel_data_folder := pathpkg.Join("test_cases", pathpkg.Base(test_case_folder))

	//Debug info:
	// variant_id := strings.Split(strings.Split(string(pathpkg.Base(param_values["-jsgf"])), "-")[0], "_")[2] //Not relevant to identify variants, just to distinguish ps_calls
	// fmt.Println("Variant id: ")
	// fmt.Println(variant_id)
	ps_call_folder := pathpkg.Join(ps_folder, "frate_"+frate)
	create_no_overwrite(ps_call_folder)

	//files with debug info:
	ps_args_file := pathpkg.Join(ps_call_folder, "ps_args.txt")
	//result_file := pathpkg.Join(ps_call_folder, "result.txt")

	log_file = pathpkg.Join("log", log_file)
	//jsgf_file = pathpkg.Join("grammar", jsgf_file)
	ctl_file = pathpkg.Join("ctl", ctl_file)
	cepdir_folder := pathpkg.Join(rel_data_folder, "audio")

	ps_cmnvec_file := pathpkg.Join(ps_call_folder, "cmnvec.txt")

	create_param_file_batch(ps_args_file, rel_data_folder, param_values, log_file, ctl_file, cepdir_folder)
	create_batch_file(ps_cmnvec_file, cmnvec)

}

func create_batch_file(file string, cmnvec []string) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	_, errw := f.WriteString(cmnvec[0])
	check(errw)

	for _, k := range cmnvec[1:] {
		_, errw = f.WriteString("," + k)
		check(errw)
	}

}

func readParams(args []string) map[string]string {

	param_values := make(map[string]string)

	//Create dictionary: key:value == -ps_option:value
	for i := 0; i < len(args)-1; i = i + 2 {
		param_values[args[i]] = args[i+1]
	}

	//Order the dictionary by key to make it easier to inspect visually
	keys := make([]string, 0, len(param_values))
	for k := range param_values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return param_values
}

func testCaseItContinuous(args []string, word string) {

	param_values := readParams(args)
	//Create new test case folder (ex: "allowed1_philip/")
	//It will create a folder for all test cases if not already created
	//Test cases are created within the 'output' folder in cli_pron: usually the 'audio_clips' folder used by testPronounce
	test_case_folder := create_test_case_folder(param_values, word)

	//copy all the files pocketsphinx_continuous have used after a call (jsgf, audio and logfn)
	//we add frate in the filename (except audio) to make it easier to inspect things.

	frate := param_values["-frate"]
	copy_files(test_case_folder, frate, param_values)
}

func testCaseItBatch(args []string, word string, logfn string, cmnvec []string) {

	param_values := readParams(args)
	//Create new test case folder (ex: "allowed1_philip/")
	//It will create a folder for all test cases if not already created
	//Test cases are created within the 'output' folder in cli_pron: usually the 'audio_clips' folder used by testPronounce
	test_case_folder := create_test_case_folder_from_batch(param_values, word)

	//copy all the files pocketsphinx_continuous have used after a call (jsgf, audio and logfn)
	//we add frate in the filename (except audio) to make it easier to inspect things.

	frate := param_values["-frate"]
	copy_files_batch(test_case_folder, frate, param_values, logfn, cmnvec)
}
