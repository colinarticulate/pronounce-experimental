package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

type Utt struct {
	Text       string
	Start, End int32
}

//--------------------------------
type PsParams map[string]string

type PsParam struct {
	Flag, Value string
}

type psError struct {
	args []string
}

func (p psError) Error() string {
	return fmt.Sprintf("Check pocketsphinx settings? args are %v\n", p.args)
}

type UttResp struct {
	Utts []Utt
	Err  error
}

type PsScan struct {
	Settings     PsParams
	defaults     PsParams
	ContextFlags []string
	RespondTo    chan UttResp
}

type Config struct {
	Model string `json:"model"`
}

type defaultParams map[string]string

func getDefaultParams() PsParams {
	defaults := make(PsParams)

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		// Not much I can do here. This should never happen but...
		fmt.Println("Failed to read config.json")
		return defaults
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Failed to parse config.json")
		return defaults
	}

	defaults["-hmm"] = path.Join(config.Model, "en-us/en-us")
	defaults["-lm"] = path.Join(config.Model, "en-us/en-us.lm.bin")
	defaults["-dict"] = path.Join(config.Model, "en-us/cmudict-en-us.dict")

	return defaults
}

func process_result() {
	test := []string{"sil kl ay m v b sil (-3641)*word,start,end*sil,3,90*(NULL),90,90*kl,91,109*(NULL),109,109*ay,110,147*(NULL),147,147*m,148,165*(NULL),165,165*v,166,172*b,173,176*sil,177,245**0000000000000000000000000000000000000"}

	raw := strings.Split(test[0], "**")

	fmt.Printf("%T", raw)
	fields := strings.Split(raw[0], "*")

	fmt.Println(fields)
	hyp := fields[0]
	header := fields[1]

	fmt.Println(hyp)
	fmt.Println(strings.Split(header, ","))
	utts := []Utt{}
	//var utts = make([]Utt, len(fields)-2)

	for i := 0; i < len(fields)-2; i++ {
		parts := strings.Split(fields[2:][i], ",")
		phoneme := parts[0]
		text_start := parts[1]
		text_end := parts[2]
		start, serr := strconv.Atoi(text_start)
		end, eerr := strconv.Atoi(text_end)

		if phoneme != "(NULL)" {
			fmt.Println(phoneme, start, end)
			utts = append(utts, Utt{phoneme, int32(start), int32(end)})

			if serr != nil || eerr != nil {
				fmt.Println(serr, eerr)
			}
		}
	}
	fmt.Println()
	//var number int32

	tnum := "1234"
	tnonu := "aabb"
	number, err := strconv.Atoi(tnum)
	notnumber, nnerr := strconv.Atoi(tnonu)

	fmt.Printf(" in text: %s (%T), in integer: %d (%T)\n", tnum, tnum, number, int32(number))
	fmt.Println("err: ", err)
	fmt.Println(tnonu, notnumber)
	fmt.Println("err: ", nnerr)

	fmt.Println("Done")
}

func extract_params(scan PsScan, bmnVec string) []string {
	params := []string{"pocketsphinx_continuous"}

	//var keys = fmtParams.keys()
	for key := range scan.defaults {
		fmt.Println(key, scan.defaults[key])
		params = append(params, key, scan.defaults[key])
	}

	for key := range scan.Settings {
		fmt.Println(key, scan.Settings[key])
		params = append(params, key, scan.Settings[key])
	}

	params = append(params, "-cmninit", bmnVec)

	return params
}

// Main function
func main_() {
	var fmtParams = getDefaultParams()

	fmt.Println(fmtParams)
	params := []string{"pocketsphinx_continuous"}

	//var keys = fmtParams.keys()
	for key := range fmtParams {
		fmt.Println(key, fmtParams[key])
		params = append(params, key, fmtParams[key])
	}

}

// //scheduler as in main branch
// func (s *Scheduler) doScan(scan PsScan, bScan batchScan) {
// 	// Set up the arguments for pocketsphinx_continuous
// 	args := []string{}
// 	for _, setting := range scan.Settings {
// 		value := setting.Value
// 		if setting.Flag == "-cmninit" {
// 			// Add in the cmn vector from the batch scan...
// 			value = strings.Join(bScan.cmnVec, ",")
// 		}
// 		args = append(args, setting.Flag, value)
// 	}

// 	// The output bytes are useless! So just return back to the caller
// 	// indictating an error - if there was one
// 	_, err := exec.Command("pocketsphinx_continuous", args...).Output()
// 	if err != nil {
// 		err = psError{
// 			args,
// 		}
// 	}
// 	scan.RespondTo <- err
// }

// //scheduler as in ps-go branch:
// //ps-xyz: this function needs to be changed
// func (s *Scheduler) doScan(scan PsScan, bScan batchScan) {
// 	var i int
// 	i, s.scanId = s.scanId, s.scanId+1
// 	debug("i, s.scanId =", i, s.scanId)
// 	// debug("And continuing...,", i)
// 	// debug("doScan->")
// 	// pocketsphinx can crash under some configurations so catch it
// 	// here
// 	defer func() {
// 		if r := recover(); r != nil {
// 			scan.RespondTo <- UttResp{
// 				[]Utt{},
// 				errors.New("pocketsphinx crashed!"),
// 			}
// 			// scan.RespondTo <- errors.New("pocketsphinx crashed!")
// 		}
// 	}()

// 	//ps-xyz: from here on:
// 	//1. []string with all the params as if it where the command line for pocketsphinx
// 	//2. load jsgf file
// 	//3. load audio file
// 	//4. ps_xyz call (jsgf, audio, params)
// 	//5. result in utterance struct
// 	//
// 	// ... Something like this:
// 	// jsgf, audio, params := converter(PsScan)
// 	// uttResp := xyz.ps_call(jsgf, audio, params)
// 	// scan.RespondTo <- uttResp
// 	//
// 	// where response (uttResp is)
// 	//
// 	// type UttResp struct {
// 	//   Utts []Utt // array of Utt
// 	//   Err  error
// 	// }
// 	//
// 	// type Utt struct {
// 	//   Text string
// 	//   Start, End int32
// 	// }
// 	//
// 	//

// 	// Get the audiofile now.
// 	audiofile, ok := scan.Settings["-infile"]
// 	if ok {
// 		// pocketsphinx doesn't want the audiofile as part of the config
// 		// settings
// 		delete(scan.Settings, "-infile")
// 	} else {
// 		// We've got no audio to play!
// 		scan.RespondTo <- UttResp{
// 			[]Utt{},
// 			errors.New("No audiofile provided"),
// 		}
// 		// scan.RespondTo <- errors.New("No audiofile provided")
// 		debug("No audiofile provided,", i)
// 		return
// 	}
// 	// debug("Calling NewConfig with cmnVec,", strings.Join(bScan.cmnVec, ","))
// 	debug("Calling NewConfig,", i)
// 	cfg := sphinx.NewConfig(
// 		options(
// 			scan.Settings,
// 			scan.defaults,
// 			strings.Join(bScan.cmnVec, ","),
// 		)...,
// 	)
// 	debug("Returned from NewConfig,", i)
// 	if cfg == nil {
// 		debug("cfg == nil!")
// 	}
// 	debug("Calling NewDecoder,", i)
// 	dec, err := NewDecoder(cfg)
// 	debug("Returned from NewDecoder,", i)
// 	if err != nil {
// 		scan.RespondTo <- UttResp{
// 			[]Utt{},
// 			err,
// 		}
// 		// scan.RespondTo <- err
// 		debug("Decoder initialisation failed,", i)
// 		return
// 	}

// 	//ps-xyz: new ps (ps_xyz) requires all the params in one go (jsgf,audio,params)
// 	artDec := artDecoder{
// 		dec, cfg, scan.Settings["-logfn"], i,
// 	}
// 	debug("Calling decodeFromFile,", i)
// 	uttResp := artDec.decodeFromFile(
// 		audiofile,
// 	)
// 	debug("Returned from decodeFromFile,", i)
// 	scan.RespondTo <- uttResp
//----------------------------------
//}
