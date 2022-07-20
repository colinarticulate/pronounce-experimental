package pron

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"

	"github.com/colinarticulate/dictionary"
	"github.com/google/uuid"
)

func newGrammarFromTemplate(filename string, word string, phons []phoneme) {
	f, err := os.Open("template.jsgf")
	if err != nil {
		debug("newGrammarFromTemplate: failed to open file. err =", err)
		log.Panic(err)
	}
	defer f.Close()

	outlines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "public <word>") {
			targetRule := "public <" + word + ">" + " = (sil+ [<any>] sil+)"
			for i, ph := range phons {
				if i == len(phons)-1 && ph == er {
					targetRule += " [<check>] " + "<er_ending>+"
					continue
				}
				targetRule += " [<check>] " + string(ph) + "+"
			}
			targetRule += " [<check>] (sil+ [<any>] sil+);"
			outlines = append(outlines, targetRule)
		} else {
			outlines = append(outlines, line)
		}
	}
	g, err := os.Create(filename)
	if err != nil {
		debug("newGrammarFromTemplate: failed to create file. err =", err)
		log.Panic()
	}
	defer g.Close()

	for _, line := range outlines {
		_, err = g.WriteString(line + "\n\n")
		if err != nil {
			debug("newGrammarFromTemplate: failed to write to file. err =", err)
			log.Panic()
		}
	}
}

func newGrammarForVariantPhons(filename string, word string, varPhons [][]phoneme) {
	f, err := os.Open("template.jsgf")
	if err != nil {
		debug("newGrammarForVariantPhons: failed to open file. err =", err)
		log.Panic(err)
	}
	defer f.Close()

	outlines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "public <word>") {
			fred := " (sil+ [<any>] sil+) "
			//targetRule := "public <" + word + ">" + " = (sil+ [<any>]* sil+)"
			targetRule := "public <" + word + ">" + " = "
			var targetForVariant []string
			for _, phons := range varPhons {
				targetRuleForVariant := "[<check>]"
				for i, ph := range phons {
					if i == len(phons)-1 && ph == er {
						targetRuleForVariant += " <er_ending>+ [<check>]"
						continue
					}
					targetRuleForVariant += " " + string(ph) + "+ [<check>]"
				}
				//targetForVariant = append(targetForVariant, "(" + targetRuleForVariant + ")" )
				targetForVariant = append(targetForVariant, "("+fred+targetRuleForVariant+fred+")")
			}
			targetRule += strings.Join(targetForVariant, " | ")
			//targetRule += " (sil+ [<any>] sil+);"
			targetRule += " ;"
			outlines = append(outlines, targetRule)
		} else {
			outlines = append(outlines, line)
		}
	}
	g, err := os.Create(filename)
	if err != nil {
		debug("newGrammarForVariantPhons: failed to create file. err =", err)
		log.Panic()
	}
	defer g.Close()

	for _, line := range outlines {
		_, err = g.WriteString(line + "\n\n")
		if err != nil {
			debug("newGrammarForVariantPhons: failed to write to file. err =", err)
			log.Panic()
		}
	}
}

var approximateVowel = []phoneme{
	l, r, w, y,
}

var voicelessPlosives = []phoneme{
	p, t,
}

var voicedPlosives = []phoneme{
	b, d, g, k,
}

var insertBefore_r = []phoneme{
	b,
}

var allowedInsert = map[phoneme][]phoneme{
	r: {
		b, f, g, m, p, v,
	},
}

func contains(list []phoneme, ph phoneme) bool {
	for _, l := range list {
		if l == ph {
			return true
		}
	}
	return false
}

func nextPhoneme(phons []phoneme, currPh int) (phoneme, bool) {
	if currPh < 0 || currPh+1 >= len(phons) {
		return "", false
	}
	return phons[currPh+1], true
}

func targetRuleForWord(word string, phons []phoneme) string {
	//end := "(sil+ [<any>] sil* [<any>] sil+)"
	//targetRule := "public <" + word + ">" + " = " + end + " "

	targetRule := "public <" + word + ">" + " = " + " "
	for i, ph := range phons {
		targetRule += " "
		switch i {
		case 0:
			if v, ok := inserts[phonemePair{sil, ph}]; ok {
				targetRule += v + " "
			}
			targetRule += string(ph) + " "
			// This could be a one phoneme word
			if len(phons)-1 == 0 {
				if v, ok := inserts[phonemePair{ph, sil}]; ok {
					targetRule += v + " "
				}
			}
		case len(phons) - 1:
			if v, ok := inserts[phonemePair{phons[i-1], ph}]; ok {
				targetRule += v + " "
			}
			targetRule += string(ph) + " "
			if v, ok := inserts[phonemePair{ph, sil}]; ok {
				targetRule += v + " "
			}
		default:
			if v, ok := inserts[phonemePair{phons[i-1], ph}]; ok {
				targetRule += v + " "
			}
			targetRule += string(ph) + "+ "
		}
	}
	//targetRule += "[<check>] " + end + ";"

	targetRule += " " + ";"

	return targetRule
}

var diphthongs = map[phoneme][2]phoneme{
	aw: {
		aa, uh,
	},
	ay: {
		aa, iy,
	},
	/*
	  er: {
	    uh, r,
	  },
	*/
	ey: {
		eh, iy,
	},
	/*
	  ow: {
	    ao, uh,
	  },
	*/
	oy: {
		ao, iy,
	},
}

func isDiphthong(ph phoneme) bool {
	/*
	  diphthongs := []phoneme{
	    aw, ay, oy, // er, ey, ow,
	  }
	  for _, d := range diphthongs {
	    if ph == d {
	      return true
	    }
	  }
	*/
	for key := range diphthongs {
		if key == ph {
			return true
		}
	}
	return false
}

func isSemiVowel(ph phoneme) bool {
	semiVowels := []phoneme{
		l, r, w, y,
	}
	for _, sv := range semiVowels {
		if ph == sv {
			return true
		}
	}
	return false
}

func targetRuleWithDiphthongs(word string, phons []phoneme, dict dictionary.Dictionary) string {
	targetRule := "public <" + word + ">" + " = " + " "

	for i, ph := range phons {
		phonemeString := string(ph)
		if isDiphthong(ph) {
			phonemeString = "<" + phonemeString + "_dip> "
		}
		targetRule += " "
		switch i {
		case 0:
			if v, ok := inserts[phonemePair{sil, ph}]; ok {
				targetRule += v + " "
			}
			targetRule += phonemeString + " "
			// This could be a one phoneme word
			if len(phons)-1 == 0 {
				if v, ok := inserts[phonemePair{ph, sil}]; ok {
					targetRule += v + " "
				}
			}
		case len(phons) - 1:
			if v, ok := inserts[phonemePair{phons[i-1], ph}]; ok {
				targetRule += v + " "
			}
			targetRule += phonemeString + " "
			if v, ok := inserts[phonemePair{ph, sil}]; ok {
				targetRule += v + " "
			}
		default:
			if v, ok := inserts[phonemePair{phons[i-1], ph}]; ok {
				targetRule += v + " "
			}
			targetRule += phonemeString + " "
		}
	}
	targetRule += " " + ";"

	debug("Diphthong target rule =", targetRule)
	return targetRule
}

func targetRuleForTrim(word string, phons []phoneme, dict dictionary.Dictionary) string {
	/*
	  end := "(sil+ [<any>] sil* [<any>] sil+)"
	  targetRule := "public <" + word + ">" + " = " + end + " [<check>] "
	  for i, ph := range phons {
	    targetRule += " " + string(ph) + " "
	    if isSemiVowel(ph) && i == len(phons) - 1 {
	      targetRule += string(p) + " "
	    }
	  }
	  targetRule += " [<check>] " + end + ";"
	*/
	targetRule := "public <" + word + ">" + " = sil+"
	for i := 0; i < 20; i++ {
		targetRule += " <any> sil+"
	}
	targetRule += ";"
	debug("Trim target rule =", targetRule)
	return targetRule
}

// Version of targetRuleForTrim from Paul. I don't think we're ever calling this though
// so changes haven't been merged in
/*
func targetRuleForTrim(word string, phons []phoneme, dict dictionary.Dictionary) string {

	//original
	//=========
	// 	end := "(sil+ [<any>] sil+)"

	// 	targetRule := "public <" + word + ">" + " = " + end + " "

	// 	for _, ph := range phons {
	// 		targetRule += " <check> " + string(ph)
	// 	}

	// 	targetRule += " " + end + " ;"


	//concept 1
	//=========
	// 	end := "(sil+ [<any>] sil+ [<any>] sil+)  <any> "
	// 	end2 := " <any>  (sil+ [<any>] sil+ [<any>] sil+)"

	// 	targetRule := "public <" + word + ">" + " = " + end + " "

	// 	for _, ph := range phons {
	// 		targetRule += " ( " + string(ph) + "  [<any>] ) <check_123> "
	// 	}

	// 	targetRule += " " + end2 + " ;"

	//Concept 2	      -- This will fail if the audio is too short, effects the clipped audio samples, those on the phone (or extended) are ok
	//==========
		idea := " <find_noise> "
		targetRule := "public <" + word + ">" + " = " + idea + " "
		targetRule += " " + " ;"


  fmt.Println("Trim target rule =", targetRule)
  return targetRule
}
*/

func NewGrammarFile(filename string, word string, phons []phoneme, dict dictionary.Dictionary, targetRule func(string, []phoneme, dictionary.Dictionary) string) {
	f, err := os.Open("template.jsgf")
	if err != nil {
		debug("NewGrammarFile: failed to open file. err =", err)
		log.Panic(err)
	}
	defer f.Close()

	outlines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "public <word>") {
			targetLine := targetRule(word, phons, dict)
			outlines = append(outlines, targetLine)
		} else {
			outlines = append(outlines, line)
		}
	}
	g, err := os.Create(filename)
	if err != nil {
		debug("NewGrammarFile: failed to create file. err =", err)
		log.Panic(err)
	}
	defer g.Close()

	for _, line := range outlines {
		_, err = g.WriteString(line + "\n")
		if err != nil {
			debug("NewGrammarFile: failed to write to file. err =", err)
			log.Panic(err, line)
		}
	}
}

func NewConfig(
	audiofile, phdictfile string,
	cmuDict dictionary.Dictionary,
	cmnVec []string,
	fparams string,
	word string,
	phons []phoneme,
	template psPhonemeSettings,
	frates []int,
	targetRule func(string, []phoneme, dictionary.Dictionary) string,
) psConfig {

	settings := make(psPhonemeSettings)
	// What the hell!? go does NOT provide a means of generating a UUID
	// citing it's a vague concept - what a cop out. So for now use Temp as
	// a folder name.
	//
	tempDir := "Temp"
	mkDir(tempDir)
	for k, v := range template {
		switch k {
		case jsgf:
			jsgfFilename := path.Join(tempDir, "forced_align_"+uuid.New().String()+"_"+word+".jsgf")
			NewGrammarFile(jsgfFilename, word, phons, cmuDict, targetRule)
			// targetRule(word, phons, cmuDict)
			settings[jsgf] = jsgfFilename
		case infile:
			settings[infile] = audiofile
		case cmninit:
			settings[cmninit] = strings.Join(cmnVec, ",")
		case dict:
			settings[dict] = phdictfile
		case featparams:
			settings[featparams] = fparams
		default:
			settings[k] = v
		}
	}
	return psConfig{
		/*
		   []int{
		     133,
		     137,
		     143,
		   },
		*/
		frates,
		settings,
		word,
		[][]phoneme{
			phons,
		},
		neighbourRules,
		"Temp",
	}
}

func new_testGrammarFile(config jsgfConfig) jsgfGrammar {
	new := jsgfGrammar{
		"#JSGF V1.0",
		"grammar test",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		config,
	}
	return new
}

/*
func TestConfig(
  audiofile, phdictfile string,
  cmuDict dictionary.Dictionary,
  word string,
  phons []phoneme,
  template psPhonemeSettings,
  frates []int,
  suiteToRun psSuite,
  targetRule func(string, []phoneme, dictionary.Dictionary) string,
  config jsgfConfig,
) psConfig {

  settings := make(psPhonemeSettings)
  // What the hell!? go does NOT provide a means of generating a UUID
  // citing it's a vague concept - what a cop out. So for now use Temp as
  // a folder name.
  //
  tempDir := "Temp"
  mkDir(tempDir)
  for k, v := range template {
    switch k {
    case jsgf:
      jsgfFilename := path.Join(tempDir, "forced_align_" + pseudo_uuid() + "_" + word + ".jsgf")
      g := new_testGrammarFile(config)
      g.new_R_target(word, phons)
      g.SaveToDisk(jsgfFilename)
      // NewGrammarFile(jsgfFilename, word, phons, cmuDict, targetRule)
      // targetRule(word, phons, cmuDict)
      settings[jsgf] = jsgfFilename
    case infile:
      settings[infile] = audiofile
    case dict:
      settings[dict] = phdictfile
    default:
      settings[k] = v
      }
  }
  return psConfig{
    []int{
      133,
      137,
      143,
    },
    frates,
    settings,
    word,
    [][]phoneme{
        phons,
    },
    neighbourRules,
    "Temp",
  }
}
*/

func TestConfig(
	outfolder, audiofile, phdictfile string,
	// cmnVec []string,
	fparams string,
	hidden_mm string,
	word string,
	phons []phoneme,
	// template psPhonemeSettings,
	frates []int,
	suiteToRun psSuite,
	// targetRule func(string, []phoneme, dictionary.Dictionary) string,
	config jsgfConfig,
) (newPsConfig, []byte) {

	settings := []psPhonemeSettings{}
	// What the hell!? go does NOT provide a means of generating a UUID
	// citing it's a vague concept - what a cop out. So for now use Temp as
	// a folder name.
	//
	tempDir := outfolder
	mkDir(tempDir)
	g := new_testGrammarFile(config)
	g.new_R_target(word, phons)
	_, file := path.Split(audiofile)
	jsgfFilename := path.Join(tempDir, "forced_align_"+uuid.New().String()+"_"+file+".jsgf")
	//g.SaveToDisk(jsgfFilename)// the only time this is written to a file
	for _, s := range suiteToRun {
		runSettings := make(psPhonemeSettings)
		for k, v := range s {
			switch k {
			case jsgf:
				runSettings[jsgf] = jsgfFilename
			case infile:
				runSettings[infile] = audiofile
			/*
			   case cmninit:
			     runSettings[cmninit] = strings.Join(cmnVec, ",")
			*/
			case dict:
				runSettings[dict] = phdictfile
			case featparams:
				runSettings[featparams] = fparams
			case hmm:
				runSettings[hmm] = hidden_mm
			default:
				if k == frate {
					debug("frate, target rule =", v, g.target.generate())
					//g.target.generate()
				}
				runSettings[k] = v
			}
		}
		settings = append(settings, runSettings)
	}
	return newPsConfig{
		settings,
		word,
		[][]phoneme{
			phons,
		},
		g.target,
		neighbourRules,
		outfolder,
	}, g.SaveToByteSlice()
}
