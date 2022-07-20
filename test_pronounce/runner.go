package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type testIn struct {
	audiofile, word string
}

type testsIn struct {
	infolder string
	tests    []testIn
}

type verdict int

const (
	good verdict = iota
	possible
	missing
	surprise
)

var verdicts = []string{
	"good",
	"possible",
	"missing",
	"surprise",
}

type result struct {
	phon       string
	goodBadEtc string
}

type results []result

type namedResult struct {
	name       string
	wordResult []result
}

type runner struct {
	schedule chan scheduledTest
	closing  chan chan error
	pending  []scheduledTest
}

type scheduledTest struct {
	args    []string
	replyTo chan []byte
}

func new() runner {
	chSch := make(chan scheduledTest)
	chC := make(chan chan error)

	r := runner{
		chSch, chC,
		[]scheduledTest{},
	}
	go r.loop()
	return r
}

func (r *runner) loop() {
	//Throttling
	maxRunningTests := 150 //100 //10 default
	runningTests := 0
	nms := 7000         //8500         // 2000
	ms := nms * 1000000 //in milliseconds
	//ms := float64(nms) * 0.001
	fmt.Println(time.Duration(ms))         //0.001=1ms
	limit := time.Tick(time.Duration(nms)) //5 //500 default
	//limit := time.Tick(1 * time.Millisecond)

	maxTest := strconv.Itoa(maxRunningTests)
	millisecons := strconv.Itoa(nms)
	filename := "timings_ipaul_x100_" + maxTest + "_" + millisecons + ".txt"
	fmt.Println(filename)
	f, err := os.Create(filename)
	check(err)
	f.WriteString("Id Waiting Response\n")
	// fms := float64(time.Millisecond)

	// // starts := make(map[int]time.Duration)
	// // ends := make(map[int]time.Duration)

	// start_test := time.Now()
	// count := 0
	for {
		select {
		case <-limit:
			// Run the next test, if there are any to run...
			if len(r.pending) > 0 && runningTests < maxRunningTests {
				go func(test scheduledTest) {
					runningTests++
					fmt.Println("running new test...")
					//out, err := exec.Command("cli_pron", test.args...).Output()
					os.Setenv("GODEBUG", "cgocheck=0")
					out, err := exec.Command("cli_pron", test.args...).Output()
					if err != nil {
						fmt.Println("testPronounce:Oops, err is", err, "Check args maybe?...", test.args)
						// What do we do here?
						fmt.Println(string(out))
					}
					runningTests--
					fmt.Println("test run...")
					test.replyTo <- out

					// //Timing Code:
					// count = count + 1
					// start := time.Now()
					// start_time := start.Sub(start_test)

					// fmt.Println("Start test(", count, ") =", start_time)
					// runningTests++
					// //out, err := exec.Command("cli_pron", test.args...).Output()
					// os.Setenv("GODEBUG", "cgocheck=0")
					// out, err := exec.Command("cli_pron", test.args...).Output()
					// if err != nil {
					// 	fmt.Println("testPronounce:Oops, err is", err, "Check args maybe?...", test.args)
					// 	// What do we do here?
					// 	fmt.Println(string(out))
					// }
					// runningTests--
					// end := time.Now()
					// elapsed := end.Sub(start)

					// fs := float64(start_time) / fms
					// fe := (float64(start_time) + float64(elapsed)) / fms
					// ss := strconv.FormatFloat(fs, 'f', -1, 64)
					// se := strconv.FormatFloat(fe, 'f', -1, 64)
					// file := filepath.Base(test.args[1])
					// idx := strings.Split(file, "_")[0]
					// f.WriteString(idx + " " + ss + " " + se + "\n")

					// fmt.Println("End test(", count-maxRunningTests+1, "), time taken =", elapsed)
					// test.replyTo <- out
				}(r.pending[0])

				r.pending = r.pending[1:]
			}
		case schTest := <-r.schedule:
			r.pending = append(r.pending, schTest)
		case errc := <-r.closing:
			errc <- nil
			close(r.schedule)
			//f.Close()
			return
		}
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (r *runner) scheduleTest(args []string, ch chan []byte) {
	schT := scheduledTest{
		args,
		ch,
	}
	r.schedule <- schT
}

func (r runner) Close() error {
	errc := make(chan error)
	r.closing <- errc
	return <-errc
}

func contains(s []string, t string) bool {
	for i := 0; i < len(s); i++ {
		if t == s[i] {
			return true
		}
	}
	return false
}

type expectationNode struct {
	expected      expectation
	nextExpecteds []expectationNode
}

func (eN *expectationNode) add(es []expectation) {
	if len(eN.nextExpecteds) == 0 {
		// This is a leaf node so add the expectations, es
		for _, e := range es {
			new := expectationNode{
				e,
				[]expectationNode{},
			}
			eN.nextExpecteds = append(eN.nextExpecteds, new)
		}
		return
	}
	for _, e := range eN.nextExpecteds {
		e.add(es)
	}
}

type expectation interface {
	pass(r results) (results, bool)
}

// If len(data) > 1 then the phonemes should be interpreted as 'any of them
// will do' for the test to be marked pass
// If len(goodBadEtc) > 1 then any of the verdicts found in the actual result
// will count as a pass
type simpleExpect struct {
	phons      []string
	goodBadEtc []string
}

func (ex simpleExpect) pass(r results) (results, bool) {
	if len(r) == 0 {
		// We've run out of results so the only way we can return true here is if
		// the verdict contains absent
		if contains(ex.goodBadEtc, "absent") {
			return r, true
		}
		return r, false
	}
	workingR := r
	for i, res := range workingR {
		if contains(ex.phons, res.phon) {
			if contains(ex.goodBadEtc, res.goodBadEtc) {
				// This is the easy bit. We've got a match so this is a pass
				return workingR[i+1:], true
			} else {
				if contains(ex.goodBadEtc, "absent") {
					// If the verdict is absent we should probably continue and try the
					// next expectation
					return r, true
				}
				return r, false
			}
		} else {
			// Well, if we get this far then r.phon is not contained in e[j].phons
			// so check to see if the simpleExpect phoneme is wild-carded
			if contains(ex.phons, "*") {
				// A wild-carded expected phoneme only  makes sense for a surprise or
				// absent verdict. We've already handled the absent case so there's
				// only the surprise case left
				if res.goodBadEtc == "surprise" && contains(ex.goodBadEtc, res.goodBadEtc) {
					return workingR[i+1:], true
				}
			}
			// Or is the verdict absent? In which case try the next expectation
			if contains(ex.goodBadEtc, "absent") {
				return r, true
			}
			return r, false
		}
	}
	// Does this always work? In the case where the last expectation has an
	// absent we should pass back true...
	return results{}, true
}

type notExpect struct {
	not simpleExpect
}

func (ex notExpect) pass(r results) (results, bool) {
	workingR := r
	for i, res := range workingR {
		if !contains(ex.not.phons, res.phon) {
			if contains(ex.not.goodBadEtc, res.goodBadEtc) {
				// This is the easy bit. We've got a match so this is a pass
				return workingR[i+1:], true
			} else {
				if contains(ex.not.goodBadEtc, "absent") {
					// If the verdict is absent we should probably continue and try the
					// next expectation
					return r, true
				}
			}
		} else {
			// Well, if we get this far then r.phon is contained in e[j].phons
			if !contains(ex.not.goodBadEtc, res.goodBadEtc) {
				return workingR[i+1:], true
			}
			return r, false
		}
	}
	// Does this always work? In the case where the last expectation has an
	// absent we should pass back true...
	return results{}, true
}

type orExpect struct {
	or [][]simpleExpect
}

func (orE orExpect) pass(r results) (results, bool) {
	for _, exs := range orE.or {
		workingR := r
		failed := false
		for _, e := range exs {
			// These are and'ed together so stop as soon as an expectation
			// doesn't pass
			r1, passed := e.pass(workingR)
			if !passed {
				failed = true
				break
			} else {
				workingR = r1
			}
		}
		if !failed {
			return workingR, true
		}
	}
	return r, false
}

func (orE orExpect) allPass(r results) ([]results, bool) {
	remResults := []results{}
	for _, exs := range orE.or {
		workingR := r
		failed := false
		for _, e := range exs {
			// These are and'ed together so stop as soon as an expectation
			// doesn't pass
			r1, passed := e.pass(workingR)
			if !passed {
				failed = true
				break
			} else {
				workingR = r1
			}
		}
		if !failed {
			remResults = append(remResults, workingR)
		}
	}
	return remResults, !(len(remResults) == 0)
}

func (e simpleExpect) isZero() bool {
	return len(e.phons) == 0 && len(e.goodBadEtc) == 0
}

type testOut struct {
	name     string
	actual   results
	expected []expectation
	pass     bool
}

type testsOut []testOut

func runTests(dictfile, phdictfile, infolder, testsToRun, expectations, outfolder, featparams, hmm string) int {
	testfile := testsToRun //filepath.Join(infolder, testsToRun)
	tests, err := fetchTests(testfile)
	if err != nil {
		fmt.Println("Failed to fetch tests, err =", err)
		return 0
	}
	testInData := testsIn{
		infolder,
		tests,
	}
	expectFile := expectations //filepath.Join(infolder, expectations)
	testOutData, err := fetchExpectations(expectFile)
	if err != nil {
		fmt.Println("Failed to fetch expectations, err =", err)
		return 0
	}

	r := new()
	defer r.Close()

	var wg sync.WaitGroup
	allResults := make(chan namedResult)

	for _, test := range testInData.tests {
		wavPath := filepath.Join(testInData.infolder, test.audiofile+".wav")
		// Create ctlfile
		ctlfile, err := SaveCtlfile(test.audiofile)
		if err != nil {
			fmt.Println("Error saving ctl file", err.Error())
			continue
		}
		args := []string{
			"-audio", wavPath, "-featparams", featparams, "-word", test.word, "-dict", dictfile, "-phdict", phdictfile, "-hmm", hmm,
		}
		wg.Add(1)
		go func(test testIn) {
			defer wg.Done()

			replyTo := make(chan []byte)
			r.scheduleTest(args, replyTo)
			out := <-replyTo

			name := testName(test.audiofile, test.word)
			result := namedResult{
				name,
				parseResult(findResult(out)),
			}
			allResults <- result

			saveToDisk(out, filepath.Join(outfolder, name+".txt"))
			_ = os.Remove(ctlfile)
		}(test)
	}

	go func() {
		wg.Wait()
		close(allResults)
		// fmt.Println("len(allResults) =", len(allResults))
	}()

	for result := range allResults {
		fmt.Println("result =", result)
		for i, test := range testOutData {
			if result.name == test.name {
				testOutData[i].actual = result.wordResult
			}
		}
	}

	for i, test := range testOutData {
		// Check result and update summary file
		testOutData[i].pass = test.checkResult(hmm)
	}
	accuracy := testOutData.summary(outfolder)
	clean_ctls("./tmp")

	return accuracy
}

func fetchTests(file string) ([]testIn, error) {
	tests := []testIn{}

	f, err := os.Open(file)
	if err != nil {
		return tests, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Each line should be a ,-separated set of strings (filename and word)
		// with some possible whitespace
		tokens := strings.FieldsFunc(line, func(arg1 rune) bool {
			return arg1 == ',' || unicode.IsSpace(arg1)
		})
		if len(tokens) != 2 {
			fmt.Println("Found bad test entry. Entry is ", line)
			continue
		}
		test := testIn{
			tokens[0], // + ".wav",
			tokens[1],
		}
		tests = append(tests, test)
	}

	// Check for duplicate tests

	duplicateCheck := map[testIn]bool{}
	for _, test := range tests {
		if _, ok := duplicateCheck[test]; ok {
			fmt.Println("Duplicate test found. Test is", test.audiofile+" "+test.word)
		} else {
			duplicateCheck[test] = true
		}
	}
	fmt.Println("Number of tests =", len(tests))
	return tests, nil
}

func testName(audiofile, word string) string {
	return audiofile + "_" + word
}

func resultFilename(audiofile, word, outfolder string) string {
	file := filepath.Base(audiofile)
	ext := filepath.Ext(file)
	file = file[:len(file)-len(ext)]

	return filepath.Join(outfolder, file+"_"+word)
}

func saveToDisk(b []byte, toPath string) {
	f, err := os.Create(toPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.Write(b)

	w.Flush()
}

func fetchExpectations(file string) (testsOut, error) {
	tests := []testOut{}

	f, err := os.Open(file)
	if err != nil {
		return tests, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.FieldsFunc(line, func(arg1 rune) bool {
			return arg1 == ',' || unicode.IsSpace(arg1)
		})
		// We expect to at least see an audiofile name, a phoneme and a verdict so
		// move onto the next line if there aren't enough tokens
		if len(tokens) < 3 {
			continue
		}
		expect, err := parseTestExpectation(tokens[1:])
		if err != nil {
			// Report error
		} else {
			test := testOut{
				tokens[0],
				[]result{}, // Fill this in later when we run the test
				expect,
				false,
			}
			tests = append(tests, test)
		}
	}
	return tests, nil
}

func parseTestExpectation(tokens []string) ([]expectation, error) {
	expects := []expectation{}
	for len(tokens) > 0 {
		var err error
		var exp expectation
		var remainingTokens []string
		switch tokens[0] {
		case "~":
			exp, remainingTokens, err = parseNotExpectation(tokens)
			expects = append(expects, exp)
		case "(":
			if exp, remainingTokens, err = parseOrExpectation(tokens); err == nil {
				expects = append(expects, exp)
			}
		default:
			exp, remainingTokens, err = parseSimpleExpectation(tokens)
			expects = append(expects, exp)
		}
		if err != nil {
			return expects, err
		}
		tokens = remainingTokens
	}
	return expects, nil
}

func parseOrExpectation(tokens []string) (expectation, []string, error) {
	expects := orExpect{}
	for len(tokens) > 0 {
		switch tokens[0] {
		case "(":
			expects.or = append(expects.or, []simpleExpect{})
			tokens = tokens[1:]
		case "||":
			expects.or = append(expects.or, []simpleExpect{})
			tokens = tokens[1:]
		case ")":
			return expects, tokens[1:], nil
		default:
			if len(expects.or) == 0 {
				//If we get here then we found an expectation before a '(' or '||'
				return expects, []string{}, parseError
			}
			e, remainingTokens, err := parseSimpleExpectation(tokens)
			if err != nil {
				return expects, []string{}, parseError
			}
			expects.or[len(expects.or)-1] = append(expects.or[len(expects.or)-1], e)
			tokens = remainingTokens
		}
	}
	// If we get here we haven't already returned so something's up with
	// the parse
	return expects, []string{}, parseError
}

func parseNotExpectation(tokens []string) (notExpect, []string, error) {
	expect := notExpect{}
	for len(tokens) > 0 {
		switch tokens[0] {
		case "~":
			tokens = tokens[1:]
		default:
			simple, remaining, err := parseSimpleExpectation(tokens)
			if err != nil {
				return expect, []string{}, err
			}
			expect.not = simple
			return expect, remaining, nil
		}
	}
	return expect, []string{}, parseError
}

func parseSimpleExpectation(tokens []string) (simpleExpect, []string, error) {
	isVerdict := func(v string) bool {
		verdicts := []string{
			"good", "possible", "missing", "surprise", "absent",
		}
		for _, verdict := range verdicts {
			if v == verdict {
				return true
			}
		}
		return false
	}
	verdictLastRead := false
	expect := simpleExpect{
		[]string{}, []string{},
	}
	for i, t := range tokens {
		if isVerdict(t) {
			if len(expect.phons) == 0 {
				// Something's wrong. We've read a verdict but not added any
				// phonemes to expect
				return expect, []string{}, parseError
			}
			expect.goodBadEtc = append(expect.goodBadEtc, t)
			verdictLastRead = true
		} else {
			if verdictLastRead {
				// We're looking at something other than a verdict so we've run on to
				// the next expectation
				return expect, tokens[i:], nil
			}
			if t == "(" || t == ")" || t == "||" {
				return expect, []string{}, parseError
			}
			expect.phons = append(expect.phons, t)
		}
	}
	// We haven't already returned so this must be the last expectation. Return
	// the expectation now
	return expect, []string{}, nil
}

func findResult(b []byte) []string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		if tokens[0] == "testPronounce" {
			return tokens[1:]
		}
	}
	return []string{}
}

func parseResult(tokens []string) []result {
	// We're expecting an even number of tokens in the result (the result's a
	// sequence of phoneme-verdict pairs)
	if len(tokens)%2 != 0 {
		return []result{}
	}
	results := []result{}
	for i := 0; i < len(tokens); i += 2 {
		result := result{
			tokens[i],
			tokens[i+1],
		}
		results = append(results, result)
	}
	return results
}

type state struct {
	i   int
	rem []results
}

type stateStack struct {
	states []state
}

func (s *stateStack) push(new state) {
	s.states = append(s.states, new)
}

func (s *stateStack) pop() (state, bool) {
	if len(s.states) == 0 {
		// Nothing to pop
		return state{}, false
	}
	retState := s.states[len(s.states)-1]
	s.states = s.states[:len(s.states)-1]

	return retState, true
}

func (s *stateStack) backtrack() (int, results, bool) {
	if len(s.states) == 0 {
		// Cant backtrack any further
		return 0, results{}, false
	}
	popped, ok := s.pop()
	if !ok {
		// There's nothing left to pop so we can't backtrack so we've reached the
		// end of the line
		return 0, results{}, false
	}
	if len(popped.rem) == 0 {
		// Should probably error here. This should never happen!
		return 0, results{}, false
	}
	// We're going to return the last set of results in popped and push the rest
	// back on they stack. We should only do tat though if there are any
	// results left to push
	if len(popped.rem) > 1 {
		new := state{
			popped.i,
			popped.rem[:len(popped.rem)-1],
		}
		s.push(new)
	}
	return popped.i, popped.rem[len(popped.rem)-1], true
}

func (r testOut) checkResult(hmm string) bool {
	stack := stateStack{}
	workingR := r.actual
	testCaseIt(r.name, r.actual, hmm)
	i := 0
	for i < len(r.expected) {
		var r1 results
		var allR []results
		var passed bool
		e := r.expected[i]
		if eOr, ok := e.(orExpect); ok {
			allR, passed = eOr.allPass(workingR)
			if passed {
				new := state{
					i + 1,
					allR[1:],
				}
				stack.push(new)
				r1 = allR[0]
			}
		} else {
			r1, passed = e.pass(workingR)
		}
		if !passed {
			// Backtrack to see if we can and try to find another expectation that
			// works
			j, r2, backtracked := stack.backtrack()
			if !backtracked {
				return false
			}
			i = j
			r1 = r2
		} else {
			i++
		}
		workingR = r1
	}
	if len(workingR) != 0 {
		// We haven't swallowed all of the actual result yet. If there's some
		// unprocessed result left over (like a surprise) then this must be a fail
		return false
	}
	return true
}

func (ts testsOut) summary(outfolder string) int {
	if err := os.Mkdir(outfolder, 0777); err != nil && !os.IsExist(err) {
		return 0
	}
	summaryFile := filepath.Join(outfolder, "000__summary__000.txt")
	f, err := os.Create(summaryFile)
	if err != nil {
		log.Fatal()
	}
	defer f.Close()

	content := ""
	passCount := 0
	runCount := 0
	for _, t := range ts {
		if len(t.actual) != 0 {
			content += t.name + ", "
			if t.pass {
				content += "PASS\n"
				passCount++
			} else {
				content += "FAIL\n"
			}
			runCount++
		}
	}
	passRate := 0
	if runCount != 0 {
		passRate = passCount * 100 / runCount
	}
	content += "Pass rate = " + strconv.Itoa(passRate) + "%\n"
	_, err = f.WriteString(content)
	if err != nil {
		// Not sure there's much else we can do if the write fails
		//
		log.Fatal(err)
	}

	return passRate
}

// func check(e error) error {
// 	if e != nil {
// 		// fmt.Println(">>>> test_pronounce : ", e)
// 		// panic(e)
// 		return e
// 	}
// 	return e
// }

// exists returns whether the given file or directory exists
func create_no_overwrite_(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err_mk := os.Mkdir(path, 0777)
		if err_mk != nil {
			return "", err_mk
		}
		return path, nil
	}
	return path, nil
}

func SaveCtlfile(audio string) (string, error) {
	create_no_overwrite_("./tmp")
	filename := filepath.Join("./tmp/", "ctl_"+audio+".txt")
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(audio)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func clean_ctls(folder string) {
	os.RemoveAll(folder)
}
