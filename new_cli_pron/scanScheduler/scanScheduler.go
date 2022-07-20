package scanScheduler

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/davidbarbera/xyz_plus/v2"
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
	ready                        []PsScan
	newScan                      chan PsScan
	maxScans                     int
	tryScan                      chan bool
	closing                      chan chan error
	logIt                        chan string
}

func New(outfolder, audiofolder, dict string) Scheduler {
	// The ctl.txt file should contain the audiofile name (minus any extension)
	cepdir, audiofile := path.Split(audiofolder)
	audiobase := strings.TrimSuffix(audiofile, path.Ext(audiofile))

	ctl := path.Join(outfolder, "ctl_"+audiobase+".txt")
	f, _ := os.Create(ctl)

	_, err := f.WriteString(audiobase)
	if err != nil {
		debug("NewBatchScan: failed to write to file. err =", err)
		log.Fatal()
	}
	defer f.Close()

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
		[]PsScan{},
		newScan,
		4,
		make(chan bool),
		closing,
		make(chan string),
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
	runningScans := 0
	l := newLogger(path.Join(s.outfolder, "..", "..", "scheduler.log"))
	for {
		select {
		case sc := <-s.newScan:
			// Create batchId
			id := sc.getBatchId()
			if batch, ok := s.batchResults[id]; ok {
				// We have a batchScan
				if batch.state == hasRun {
					go func(sc PsScan, bSc batchScan, l logger) {
						s.doScan(sc, l)
					}(sc, *batch, l)
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
				go func(sc PsScan, cepdir, ctl, dict string, s *Scheduler) {
					batchScanDone <- batchResult{
						id, doBatchScan(sc, cepdir, ctl, dict, s),
					}
				}(sc, s.cepdir, s.ctl, s.dict, s)
				s.pending[id] = append(s.pending[id], sc)
			}
		case bRes := <-batchScanDone:
			s.batchResults[bRes.id].cmnVec = bRes.cmnVec
			s.batchResults[bRes.id].state = hasRun

			// Add all the scans pending this batch scan to the ready queue...
			for _, item := range s.pending[bRes.id] {
				item.prepare(s.batchResults[bRes.id].cmnVec)
				s.ready = append(s.ready, item)
			}
			// ... and try to run them
			go func() {
				s.tryScan <- false
			}()
		case scanDone := <-s.tryScan:
			if scanDone {
				runningScans--
			}
			// Now check if we can do any more scans. These are throttled by
			// maxScans...
			if s.canRunScans(runningScans) {
				runningScans++
				go func(sc PsScan, l logger) {
					s.doScan(sc, l)
				}(s.ready[0], l)
				s.ready = s.ready[1:]
				go func() {
					s.tryScan <- false
				}()
			}
		case entry := <-s.logIt:
			l.addEntry(entry)
		case errc := <-s.closing:
			errc <- nil
			close(s.newScan)
			return
		}
	}
}

func (s Scheduler) canRunScans(runningScans int) bool {
	if len(s.ready) == 0 {
		// There's nothing to run
		return false
	}
	if s.maxScans == -1 {
		// We're unthrottled
		return true
	}
	return s.maxScans > runningScans
}

type PsParam struct {
	Flag, Value string
}

type PsError struct {
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

func toUttResp(xyzUtts xyz_plus.UttResp) UttResp {
	utts := []Utt{}
	for _, utt := range xyzUtts.Utts {
		utts = append(utts, Utt{
			utt.Text,
			utt.Start,
			utt.End,
		})
	}
	return UttResp{
		utts,
		xyzUtts.Err,
	}
}

func (p PsError) Error() string {
	return fmt.Sprintf("Check xyz_plus.Ps_plus_call settings? args are %v\n", p.args)
}

type PsScan struct {
	Settings     []PsParam
	ContextFlags []string
	RespondTo    chan UttResp
	Jsgf_buffer  []byte
	Audio_buffer []byte
	Parameters   []string
}

func (s Scheduler) ScheduleScan(sc PsScan) {
	s.newScan <- sc
}

func (ps *PsScan) prepare(cmnVec []string) {
	args := []string{"pocketsphinx_continuous"}

	for _, setting := range ps.Settings {
		value := setting.Value
		if setting.Flag == "-cmninit" {
			// Add in the cmn vector from the batch scan...
			value = strings.Join(cmnVec, ",")
		}
		args = append(args, setting.Flag, value)
	}
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

func (s *Scheduler) doScan(scan PsScan, l logger) {
	var args []string
	defer func() {
		if r := recover(); r != nil {
			// Log this
			s.logIt <- fmt.Sprintf(
				"xyz_plus.Ps_plus_call crashed with args, %v",
				args,
			)
			scan.RespondTo <- UttResp{
				[]Utt{},
				PsError{args},
			}
		}
	}()

	// Set up the arguments for pocketsphinx_continuous
	//result := []xyz_plus.Utt
	args = []string{"pocketsphinx_continuous"}
	word := "word"

	for _, setting := range scan.Settings {
		value := setting.Value
		if setting.Flag == "-word" {
			word = value
			continue
		}
		args = append(args, setting.Flag, value)
	}

	testCaseItContinuous(args, word) //After the call so we can add the log file from pocketsphinx_continuous to the test case.

	result := xyz_plus.Ps_plus_call(scan.Jsgf_buffer, scan.Audio_buffer, args)
	scan.RespondTo <- toUttResp(result)

	s.tryScan <- true
}

func (s *Scheduler) DoScan(scan PsScan) {
	s.newScan <- scan
}

func doBatchScan(scan PsScan, cepdir, ctl, dict string, s *Scheduler) (defVec []string) {
	var args []string
	contains := func(ss []string, s string) bool {
		for _, t := range ss {
			if t == s {
				return true
			}
		}
		return false
	}
	// Catch any panic caused by Ps_batch_plus_call crashing
	defer func() {
		if r := recover(); r != nil {
			// Should we also log this?...
			s.logIt <- fmt.Sprintf(
				"xyz_plus.Ps_batch_plus_call crashed with args, %v",
				args,
			)
			// Provide a default cmninit value
			defVec = []string{
				"41.00", "-5.29", "-0.12", "5.09", "2.48", "-4.07", "-1.37", "-1.78", "-5.08", "-2.05", "-6.45", "-1.42", "1.17",
			}
		}
	}()

	word := "word"
	for _, setting := range scan.Settings {
		if setting.Flag == "-word" {
			word = setting.Value
			continue
		}
	}

	args = []string{"pocketsphinx_batch", //required for the xyz_plus API
		"-adcin", "yes", "-cepdir", cepdir, "-cepext", ".wav", "-ctl", ctl, "-dict", dict,
	}
	// Add any other parameters. What happens if a setting is already included in
	// args above?
	for _, setting := range scan.Settings {
		if contains(scan.ContextFlags, setting.Flag) {
			args = append(args, setting.Flag, setting.Value)
		}
	}

	res := xyz_plus.Ps_batch_plus_call(scan.Audio_buffer, args)
	cmnVec := res.Cmn
	if res.Err != nil {
		cmnVec = []string{
			"41.00", "-5.29", "-0.12", "5.09", "2.48", "-4.07", "-1.37", "-1.78", "-5.08", "-2.05", "-6.45", "-1.42", "1.17",
		}
	}

	// Adding a dummy logfn value to keep the call to testCaseItBatch from
	// complaining
	logfn := ""
	testCaseItBatch(args, word, logfn, cmnVec)

	return cmnVec
}

func (s *Scheduler) Close() {
	// _ = os.RemoveAll(s.outfolder)
}
