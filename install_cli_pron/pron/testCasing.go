//go:build testCase
// +build testCase

package pron

import (
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"sort"
	"strconv"
	"strings"
	//"time"

	//"os"
	//"xyz"
	"github.com/colinarticulate/scanScheduler"
)

func check(e error) {
	if e != nil {
		fmt.Println(">>>>>> pron !!!! ", e)
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
	parts := []string{""}
	if strings.Contains(audiofile, "_fixed") {
		parts = strings.Split(audiofile, "_fixed")
	} else {
		parts = strings.Split(audiofile, ".")
	}
	test_case_name := parts[0]

	return test_case_name
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
			return 0, fmt.Errorf(">>>>>> pron: !!!! %s is not a regular file", src)
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

func create_result_file(file string, result []psPhonemeDatum) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	n := len(result)

	for i := 0; i < n; i++ {
		_, errw := f.WriteString(string(result[i].phoneme) + " " + strconv.Itoa(result[i].start) + " " + strconv.Itoa(result[i].end) + "\n")
		check(errw)
	}

}

func extract_variant(result []psPhonemeDatum) string {
	variant := "_"

	n := len(result)
	for i := 0; i < n; i++ {
		variant += string(result[i].phoneme) + "_" //+ " " + strconv.Itoa(result[i].start) + " " + strconv.Itoa(result[i].end) + "\n")
	}

	return variant
}

func copy_files(test_case_folder string, frate string, param_values map[string]string, result []psPhonemeDatum) {

	//src_tmp_folder := pathpkg.Dir(param_values["-jsgf"])

	//filenames without path
	jsgf_file := pathpkg.Base(param_values["-jsgf"])
	log_file := pathpkg.Base(param_values["-logfn"])
	audio_file := pathpkg.Base(param_values["-infile"])

	//tmp folder used by cli_pron
	//tmp_folder := pathpkg.Base(src_tmp_folder)
	//tmp folder for our test case
	//dst_tmp_folder := pathpkg.Join(test_case_folder, tmp_folder)
	//create_no_overwrite(dst_tmp_folder)
	//folder to keep debug data to be used from the tmp folder
	debug_folder := pathpkg.Join(test_case_folder, "debug")
	create_no_overwrite(debug_folder)

	pron_folder := pathpkg.Join(debug_folder, "pron")
	create_no_overwrite(pron_folder)

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
	ps_call_folder := pathpkg.Join(pron_folder, "ps_call-"+variant_id+"_frate_"+frate)
	create_no_overwrite(ps_call_folder)

	//files with debug info:
	ps_args_file := pathpkg.Join(ps_call_folder, "ps_args.txt")
	result_file := pathpkg.Join(ps_call_folder, "result.txt")

	log_file = pathpkg.Join("log", log_file)
	jsgf_file = pathpkg.Join("grammar", jsgf_file)
	audio_file = pathpkg.Join("audio", audio_file)

	create_param_file(ps_args_file, rel_data_folder, param_values, jsgf_file, log_file, audio_file)
	create_result_file(result_file, result)

}

func internalInit(param_values map[string]string) []scanScheduler.PsParam {

	params := []scanScheduler.PsParam{}
	for k := range param_values {
		param := scanScheduler.PsParam{
			Flag:  k,
			Value: param_values[k],
		}
		params = append(params, param)
	}

	return params
}

func extracParameters(ps_params []scanScheduler.PsParam) map[string]string {
	param_values := make(map[string]string)

	//Create dictionary: key:value == -ps_option:value
	for i := range ps_params {
		param_values[ps_params[i].Flag] = ps_params[i].Value
	}

	//Order the dictionary by key to make it easier to inspect visually
	keys := make([]string, 0, len(param_values))
	for k := range param_values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return param_values
}

func testCaseIt(ps_params []scanScheduler.PsParam, result []psPhonemeDatum, word string) {

	param_values := extracParameters(ps_params)
	//Create new test case folder (ex: "allowed1_philip/")
	//It will create a folder for all test cases if not already created
	//Test cases are created within the 'output' folder in cli_pron: usually the 'audio_clips' folder used by testPronounce
	test_case_folder := create_test_case_folder(param_values, word)

	//copy all the files pocketsphinx_continuous have used after a call (jsgf, audio and logfn)
	//we add frate in the filename (except audio) to make it easier to inspect things.
	frate := param_values["-frate"]
	copy_files(test_case_folder, frate, param_values, result)

}

func testCaseAudio(originalwavFile string, word string) {
	originalwav := pathpkg.Base(originalwavFile)
	audios_dir := pathpkg.Dir(originalwavFile)

	wav_fixed := strings.Split(originalwav, ".")[0] + "_fixed.wav"
	wav_fixed_lowpass := strings.Split(originalwav, ".")[0] + "_fixed_lowpass.wav"
	wav_fixed_trimmed := strings.Split(originalwav, ".")[0] + "_fixed_trimmed.wav"

	tc_folder := pathpkg.Join(audios_dir, "test_cases")

	create_no_overwrite(tc_folder)

	test_case_name := test_case_name_from_audiofile(originalwav)
	test_case_folder := pathpkg.Join(tc_folder, test_case_name+"_"+word)

	create_no_overwrite(test_case_folder)

	audiofolder := pathpkg.Join(test_case_folder, "audio")
	create_no_overwrite(audiofolder)

	src_fixed := pathpkg.Join(audios_dir, wav_fixed)
	src_fixed_lowpass := pathpkg.Join(audios_dir, wav_fixed_lowpass)
	src_fixed_trimmed := pathpkg.Join(audios_dir, wav_fixed_trimmed)

	dst_original := pathpkg.Join(audiofolder, originalwav)
	dst_fixed := pathpkg.Join(audiofolder, wav_fixed)
	dst_fixed_lowpass := pathpkg.Join(audiofolder, wav_fixed_lowpass)
	dst_fixed_trimmed := pathpkg.Join(audiofolder, wav_fixed_trimmed)

	_, err := copy_file(originalwavFile, dst_original)
	check(err)
	_, err = copy_file(src_fixed, dst_fixed)
	check(err)
	_, err = copy_file(src_fixed_lowpass, dst_fixed_lowpass)
	check(err)
	_, err = copy_file(src_fixed_trimmed, dst_fixed_trimmed)
	check(err)

}
