//go:build testCase
// +build testCase

package main

import (
	"fmt"
	//"io"
	"os"
	pathpkg "path"
	//"sort"
	//"strconv"
	"encoding/json"
	"strings"
	//"time"
	//"os"
	//"xyz"
	//"github.com/colinarticulate/scanScheduler"
)

func check(e error) {
	if e != nil {
		fmt.Println(">>>> test_pronounce : ", e)
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

func create_test_case_directories_test_pronounce(name string, outfolder string) string {

	tc_folder := pathpkg.Join(outfolder, "test_cases")

	//fmt.Println(tc_folder)
	create_no_overwrite(tc_folder)

	case_folder := pathpkg.Join(tc_folder, name)
	create_no_overwrite(case_folder)

	return case_folder
}

func create_model_results_directories(outfolder, hmm string) string {

	models_folder := pathpkg.Join(outfolder, "models_results")

	//fmt.Println(tc_folder)
	create_no_overwrite(models_folder)
	model_name := pathpkg.Base(hmm)
	model_folder := pathpkg.Join(models_folder, model_name)
	create_no_overwrite(model_folder)

	return model_folder
}

func create_result_file(file string, result string) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	_, errw := f.WriteString(result)
	check(errw)
}

func ToJSON(results []result) []byte {
	type JSON_result struct {
		//Letters  string `json:"letters"`
		Phon    string `json:"phonemes"`
		Verdict string `json:"verdict"`
	}
	type JSON_results struct {
		Results []JSON_result `json:"results"`
		//ErrorMsg *string       `json:"err"`
	}
	jResults := []JSON_result{}
	for _, result := range results {
		// phons := []string{}
		// for _, phon := range result.Phons {
		// 	phons = append(phons, cmubetToIpa[phon])
		// }
		jResults = append(jResults, JSON_result{
			//result.Letters,
			//strings.Join(phons, " "),
			result.phon,
			result.goodBadEtc, //verdicts[result.GoodBadEtc],
		})
	}
	// All this just so I can get Go to print null in the JSON when there's no
	// error to report
	//
	// var errStr *string
	// if err != nil {
	// 	temp := err.Error()
	// 	errStr = &temp
	// }
	out := JSON_results{
		jResults,
		//errStr,
	}
	j, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Print("toJSON: call to MarshalIndent failed. err =", err)
		//log.Panic()
	}
	return j
}

func standardise_name(name string) string {
	keywords := []string{"male", "female", "deven", "tressa", "paul", "colin", "tomo", "hossein", "khurrum", "philip", "ashwin", "paulo", "andrei", "svetlana", "rishka"}
	parts := strings.Split(name, "_")
	last := len(parts) - 1
	for _, subs := range keywords {
		if strings.Contains(parts[last], subs) {
			return name
		}
	}
	test_case_name := parts[0]
	for i := 1; i < last; i++ {
		test_case_name = test_case_name + "_" + parts[i]
	}
	return test_case_name

}

func testCaseIt(name string, result []result, hmm string) {
	outfolder := "/home/dbarbera/Repositories/test_pronounce/audio_clips/"
	//"/home/dbarbera/Repositories/test_pronounce/audio_clips/"

	//standard_name := standardise_name(name)
	standard_name := name
	case_folder := create_test_case_directories_test_pronounce(standard_name, outfolder)

	result_folder := pathpkg.Join(case_folder, "result")
	create_no_overwrite(result_folder)
	result_file := pathpkg.Join(result_folder, standard_name+".txt")

	create_result_file(result_file, string(ToJSON(result)))

	model_folder := create_model_results_directories(outfolder, hmm)
	result_file = pathpkg.Join(model_folder, standard_name+".txt")
	create_result_file(result_file, string(ToJSON(result)))

}
