package scanScheduler

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/davidbarbera/articulate-pocketsphinx-go/xyz_plus"
)

type batchState int

const (
	notRun batchState = iota
	running
	hasRun
)

type batchScan struct {
	id     batchId
	state  batchState
	cmnVec []string
}

type batchId string

type Scheduler struct {
	outfolder, cepdir, ctl, dict string
	batchResults                 map[batchId]*batchScan
	pending                      map[batchId][]PsScan
	newScan                      chan PsScan
	closing                      chan chan error
}

func New(outfolder, audiofolder, dict string) Scheduler {
	// The ctl.txt file should contain the audiofile name (minus any extension)
	cepdir, audiofile := path.Split(audiofolder)
	audiobase := strings.TrimSuffix(audiofile, path.Ext(audiofile))

	ctl := path.Join(outfolder, "ctl_"+audiobase+".txt")
	f, err := os.Create(ctl)
	defer f.Close()

	_, err = f.WriteString(audiobase)
	if err != nil {
		debug("NewBatchScan: failed to write to file. err =", err)
		log.Fatal()
	}

	batchResults := make(map[batchId]*batchScan)
	pending := make(map[batchId][]PsScan)
	newScan := make(chan PsScan)
	closing := make(chan chan error)

	sch := Scheduler{
		outfolder,
		cepdir,
		ctl,
		dict,
		batchResults,
		pending,
		newScan,
		closing,
	}
	go sch.loop()

	return sch
}

// Schedule scan requests as they arrive and run them when we're ready
func (s *Scheduler) loop() {
	type batchResult struct {
		id     batchId
		cmnVec []string
	}
	batchScanDone := make(chan batchResult)
	for {
		select {
		case sc := <-s.newScan:
			// Create batchId
			id := sc.getBatchId()
			if batch, ok := s.batchResults[id]; ok {
				// We have a batchScan
				if batch.state == hasRun {
					go func(sc PsScan, bSc batchScan) {
						s.doScan(sc, bSc)
					}(sc, *batch)
				} else {
					// Run it later when the batch scan is complete
					s.pending[id] = append(s.pending[id], sc)
				}
			} else {
				// Kick off a batch scan now and add this scan to a pending queue
				batch := batchScan{
					id,
					running,
					[]string{},
				}
				s.batchResults[id] = &batch
				go func(sc PsScan, cepdir, ctl, dict string) {
					batchScanDone <- batchResult{
						id, s.batchResults[id].doBatchScan(sc, cepdir, ctl, dict),
					}
				}(sc, s.cepdir, s.ctl, s.dict)
				s.pending[id] = append(s.pending[id], sc)
			}
		case bRes := <-batchScanDone:
			s.batchResults[bRes.id].cmnVec = bRes.cmnVec
			s.batchResults[bRes.id].state = hasRun

			// Now check the pending queue and peel off any waiting scans
			for _, item := range s.pending[bRes.id] {
				go func(sc PsScan) {
					s.newScan <- sc
				}(item)
			}
		case errc := <-s.closing:
			errc <- nil
			close(s.newScan)
			return
		}
	}
}

type PsParam struct {
	Flag, Value string
}

type psError struct {
	args []string
}

type Utt struct {
	Text       string
	Start, End int32
}

type UttResp struct {
	Utts []Utt
	Err  error
}

func (p psError) Error() string {
	return fmt.Sprintf("Check pocketsphinx settings? args are %v\n", p.args)
}

type PsScan struct {
	Settings     []PsParam
	ContextFlags []string
	//RespondTo    chan error
	RespondTo    chan []xyz_plus.Utt
	Jsgf_buffer  []byte
	Audio_buffer []byte
	Parameters   []string
}

func (s Scheduler) ScheduleScan(sc PsScan) {
	s.newScan <- sc
}

func (p PsScan) getBatchId() batchId {
	contains := func(ss []string, t string) bool {
		for _, s := range ss {
			if s == t {
				return true
			}
		}
		return false
	}
	id := ""
	// The context we want is a string made from the array of flags and
	// corresponding values
	/*
	  // This can result in lengthy file names particularly if the number of
	  // settings gets large resulting in a failure to create a file containing
	  // a cmn vector
	  for _, setting := range p.Settings {
	    if contains(p.ContextFlags, setting.Flag) {
	      id += "_" + setting.Flag + "_" + setting.Value
	    }
	  }
	*/
	// For now create the batch id using only the -frate flag and its value
	for _, setting := range p.Settings {
		if contains(p.ContextFlags, setting.Flag) {
			if setting.Flag == "-frate" {
				id += "_" + setting.Flag + "_" + setting.Value
				break
			}
		}
	}
	return batchId(id)
}

func (s *Scheduler) doScan(scan PsScan, bScan batchScan) {
	// Set up the arguments for pocketsphinx_continuous
	//result := []xyz_plus.Utt
	args := []string{"pocketsphinx_continuous"}
	word := "word"

	for _, setting := range scan.Settings {
		value := setting.Value
		if setting.Flag == "-word" {
			word = value
			continue
		}
		if setting.Flag == "-cmninit" {
			// Add in the cmn vector from the batch scan...
			value = strings.Join(bScan.cmnVec, ",")
		}
		args = append(args, setting.Flag, value)
	}

	scan.Parameters = args

	// The output bytes are useless! So just return back to the caller
	// indictating an error - if there was one

	// _, err := exec.Command("pocketsphinx_continuous", args...).Output()

	testCaseItContinuous(args, word) //After the call so we can add the log file from pocketsphinx_continuous to the test case.
	// if err != nil {
	// 	err = psError{
	// 		args,
	// 	}
	// }
	// scan.RespondTo <- err

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		scan.RespondTo <- UttResp{
	// 			[]Utt{},
	// 			errors.New("pocketsphinx crashed!"),
	// 		}
	// 		// scan.RespondTo <- errors.New("pocketsphinx crashed!")
	// 	}
	// }()
	//start := time.Now()
	scan.RespondTo <- xyz_plus.Ps_plus_call(scan.Jsgf_buffer, scan.Audio_buffer, args)
	// elapsed := time.Since(start)
	// fmt.Printf("continuous: %s\n", elapsed)
}

func (s *Scheduler) DoScan(scan PsScan) {
	s.newScan <- scan
}

func (b batchScan) doBatchScan(scan PsScan, cepdir, ctl, dict string) []string {

	contains := func(ss []string, s string) bool {
		for _, t := range ss {
			if t == s {
				return true
			}
		}
		return false
	}
	// Get the logfn value out of the scan and modify it (that is, add the batch
	// id) for the batch scan
	logfn := ""
	word := "word"
	for _, setting := range scan.Settings {
		if setting.Flag == "-word" {
			word = setting.Value
			continue
		}
		if setting.Flag == "-logfn" {
			dir, file := path.Split(setting.Value)
			ext := path.Ext(file)
			file = file[:len(file)-len(ext)] + string(b.id) + ext
			logfn = path.Join(dir, file)
		}
	}

	args := []string{"pocketsphinx_batch", //required for the xyz_plus API
	//args := []string{"xyzpocketsphinx_batch", //   <----- Should this be used instead of the above?  PE 20/Apr/21
		"-adcin", "yes", "-cepdir", cepdir, "-cepext", ".wav", "-ctl", ctl, "-dict", dict, "-logfn", logfn,
	}
	// Add any other parameters. What happens if a setting is already included in
	// args above?
	for _, setting := range scan.Settings {
		if contains(scan.ContextFlags, setting.Flag) {
			args = append(args, setting.Flag, setting.Value)
		}
	}

	// _, err := exec.Command("pocketsphinx_batch", args...).Output()

	// if err != nil {
	// 	debug("Oops, check pocketsphinx settings? args are...", args)
	// 	return []string{}
	// }
	// cmnVec := b.getCmnVec(logfn)
	//start := time.Now()
	cmnVec := xyz_plus.Ps_batch_plus_call(scan.Audio_buffer, args)
	// elapsed := time.Since(start)
	// fmt.Printf("batch: %s\t%s\n", elapsed, cmnVec)
	testCaseItBatch(args, word, logfn, cmnVec)

	// We're done with the batch file so remove it now
	// if err := os.Remove(logfn); err != nil {
	// 	// Not much I can do here - I failed to remove the logfile...
	// }
	return cmnVec
}

func (s *Scheduler) Close() {
	// _ = os.RemoveAll(s.outfolder)
}
