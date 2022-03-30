package pron

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/colinarticulate/dictionary"
	"github.com/colinarticulate/scanScheduler"
	"github.com/davidbarbera/articulate-pocketsphinx-go/xyz_plus"
	"github.com/google/uuid"
)

type psError int

const (
	psAborted psError = iota
)

var psErrors = []string{
	"pocketsphinx aborted?",
}

func (p psError) Error() string {
	return psErrors[p]
}

type phoneme string

const (
	aa  = "aa"
	ae  = "ae"
	ah  = "ah"
	ao  = "ao"
	aw  = "aw"
	ax  = "ax"
	ay  = "ay"
	b   = "b"
	ch  = "ch"
	d   = "d"
	dh  = "dh"
	eh  = "eh"
	ehr = "ehr"
	er  = "er"
	ey  = "ey"
	f   = "f"
	g   = "g"
	hh  = "hh"
	ih  = "ih"
	ing = "ing"
	iy  = "iy"
	jh  = "jh"
	k   = "k"
	l   = "l"
	m   = "m"
	n   = "n"
	ng  = "ng"
	oh  = "oh"
	ow  = "ow"
	oy  = "oy"
	p   = "p"
	r   = "r"
	s   = "s"
	sh  = "sh"
	sil = "sil"
	ss  = "ss"
	t   = "t"
	th  = "th"
	uh  = "uh"
	uw  = "uw"
	v   = "v"
	w   = "w"
	y   = "y"
	z   = "z"
	zh  = "zh"

	axl = "axl"
	axm = "axm"
	axn = "axn"
	axr = "axr"
	ks  = "ks"
	kw  = "kw"
	dz  = "dz"

	uwl = "uwl"
	uwn = "uwn"
	uwm = "uwm"

	ts = "ts"
	kl = "kl"
	pl = "pl"
	bl = "bl"

	kr = "kr"
	gr = "gr"
	tr = "tr"
	pr = "pr"

	thr = "thr"
	ihl = "ihl"
	yuw = "yuw"
	sts = "sts"
	st  = "st"
	ehl = "ehl"
	kt  = "kt"
	fl  = "fl"
)

var cmubetToIpa = map[phoneme]string{
	aa:  "ɑ",
	ae:  "æ",
	ah:  "ʌ",
	ao:  "ɔ",
	aw:  "ɑʊ",
	ax:  "ə",
	ay:  "ɑɪ",
	b:   "b",
	ch:  "ʧ",
	d:   "d",
	dh:  "ð",
	eh:  "ɛ",
	ehr: "ɛː",
	er:  "ɜɹ",
	ey:  "eɪ",
	f:   "f",
	g:   "g",
	hh:  "h",
	ih:  "ɪ",
	ing: "ɪŋ",
	iy:  "iː",
	jh:  "ʤ",
	k:   "k",
	l:   "l",
	m:   "m",
	n:   "n",
	//nd: "nd",
	ng: "ŋ",
	oh: "ɒ",
	ow: "oʊ",
	oy: "ɔɪ",
	p:  "p",
	r:  "ɹ",
	s:  "s",
	//st: "st",
	sh:  "ʃ",
	sil: "sil",
	ss:  "ss",
	t:   "t",
	th:  "θ",
	uh:  "ʊ",
	uw:  "u",
	v:   "v",
	w:   "w",
	y:   "j",
	z:   "z",
	zh:  "ʒ",

	axl: "axl",
	axm: "axm",
	axn: "axn",
	axr: "əɹ",
	ks:  "ks",
	kw:  "kw",
	dz:  "dz",

	uwl: "uwl",
	uwn: "uwn",
	uwm: "uwm",

	ts: "ts",
	kl: "kl",
	pl: "pl",
	bl: "bl",

	kr: "kr",
	gr: "gr",
	tr: "tr",
	pr: "pr",

	thr: "thr",
	ihl: "ihl",
	yuw: "yuw",
	sts: "sts",
	st:  "st",
	ehl: "ehl",
	kt:  "kt",
	fl:  "fl",
}

type neighbours map[phoneme][]phoneme

func (n neighbours) neighbours(p phoneme) []phoneme {
	return n[p]
}

func (n neighbours) isNeighbour(this, to phoneme) bool {
	neighboursOfThis := n[this]
	for _, neighbour := range neighboursOfThis {
		if to == neighbour {
			return true
		}
	}
	return false
}

var neighbourRules = neighbours{
	aa: {
		aa,
		//ah,
		//er,
		//ao,
	},
	ae: {
		ae,
		//eh,
		//er,
		//ah,
	},
	ah: {
		ah,
		//ae,
		//er,
		//aa,
	},
	ao: {
		ao,
		//aa,
		//er,
		//uh,
	},
	ax: {
		ax,
	},
	axr: {
		axr,
	},
	aw: {
		aw,
		//aa,
		//uh,
		//ow,
	},
	ay: {
		ay,
		//aa,
		//iy,
		//oy,
		//ey,
	},
	b: {
		b,
		//p,
		//d,
	},
	ch: {
		ch,
		//sh,
		//jh,
		//t,
	},
	dh: {
		dh,
		//th,
		//z,
		//v,
	},
	d: {
		d,
		//t,
		//jh,
		//g,
		//b,
	},
	eh: {
		eh,
		//ih,
		//er,
		//ae,
	},
	ehr: {
		ehr,
	},
	er: {
		er,
		//eh,
		//ah,
		//ao,
	},
	ey: {
		ey,
		//eh,
		//iy,
		//ay,
	},
	f: {
		f,
		//hh,
		//th,
		//v,
	},
	g: {
		g,
		//k,
		//d,
	},
	hh: {
		hh,
		//th,
		//f,
		//p,
		//t,
		//k,
	},
	ih: {
		ih,
		//iy,
		//eh,
	},
	ing: {
		ing,
	},
	iy: {
		iy,
		//ih,
	},
	jh: {
		jh,
		//ch,
		//zh,
		//d,
	},
	k: {
		k,
		//g,
		//t,
		//hh,
	},
	l: {
		l,
		//hh,
		//r,
		//w,
	},
	m: {
		m,
		//n,
	},
	ng: {
		ng,
		//n,
	},

	//nd: {
	//  nd,
	//},
	n: {
		n,
		//m,
		//ng,
	},
	oh: {
		oh,
	},
	ow: {
		ow,
		//ao,
		//uh,
		//aw,
	},
	oy: {
		oy,
		//ao,
		//iy,
		//ay,
	},
	p: {
		p,
		//t,
		//b,
		//hh,
	},
	r: {
		r,
		//y,
		//l,
	},
	s: {
		s,
		//sh,
		//z,
		//th,
	},
	sh: {
		sh,
		//s,
		//zh,
		//ch,
	},
	//st: {
	//  st,
	//},
	t: {
		t,
		//ch,
		//k,
		//d,
		//p,
		//hh,
	},
	th: {
		th,
		//s,
		//dh,
		//f,
		//hh,
	},
	uh: {
		uh,
		//ao,
		//uw,
	},
	uw: {
		uw,
		//uh,
	},
	v: {
		v,
		//f,
		//dh,
	},
	w: {
		w,
		//l,
		//y,
	},
	y: {
		y,
		//w,
		//r,
	},
	z: {
		z,
		//s,
		//dh,
	},
	zh: {
		zh,
		//sh,
		//z,
		//jh,
	},

	axl: {
		axl,
	},
	axm: {
		axm,
	},
	axn: {
		axn,
	},
	ks: {
		ks,
	},
	kw: {
		kw,
	},
	dz: {
		dz,
	},
	uwl: {
		uwl,
	},
	uwm: {
		uwm,
	},
	uwn: {
		uwn,
	},
	ts: {
		ts,
	},
	kl: {
		kl,
	},
	pl: {
		pl,
	},
	bl: {
		bl,
	},

	kr: {
		kr,
	},
	gr: {
		gr,
	},
	tr: {
		tr,
	},
	pr: {
		pr,
	},

	thr: {
		thr,
	},
	ihl: {
		ihl,
	},
	yuw: {
		yuw,
	},
	sts: {
		sts,
	},
	st: {
		st,
	},
	ehl: {
		ehl,
	},
	kt: {
		kt,
	},
	fl: {
		fl,
	},
}

type psFlag string

const (
	alpha           = "-alpha"
	backtrace       = "-backtrace"
	beam            = "-beam"
	bestpath        = "-bestpath"
	cmn             = "-cmn"
	cmninit         = "-cmninit"
	dict            = "-dict"
	dither          = "-dither"
	doublebw        = "-doublebw"
	featparams      = "-featparams"
	frate           = "-frate"
	fsgusefiller    = "-fsgusefiller"
	fwdflat         = "-fwdflat"
	infile          = "-infile"
	jsgf            = "-jsgf"
	logfn           = "-logfn"
	lw              = "-lw"
	maxhmmpf        = "-maxhmmpf"
	nfft            = "-nfft"
	nwpen           = "-nwpen"
	pbeam           = "-pbeam"
	pip             = "-pip"
	remove_dc       = "-remove_dc"
	remove_noise    = "-remove_noise"
	remove_silence  = "-remove_silence"
	topn            = "-topn"
	vad_threshold   = "-vad_threshold"
	vad_startspeech = "-vad_startspeech"
	vad_postspeech  = "-vad_postspeech"
	vad_prespeech   = "-vad_prespeech"
	wbeam           = "-wbeam"
	wip             = "-wip"
	wlen            = "-wlen"
	pl_beam         = "-pl_beam"     //		1e-10		Beam width applied to phone loop search for lookahead
	pl_pbeam        = "-pl_pbeam"    //		1e-10		Beam width applied to phone loop transitions for lookahead
	pl_pip          = "-pl_pip"      //		1.0		Phone insertion penalty for phone loop
	pl_weight       = "-pl_weight"   //		3.0		Weight for phoneme lookahead penalties
	pl_window       = "-pl_window"   //		5		Phoneme lookahead window size, in frames
	allphone_ci     = "-allphone_ci" // Perform phoneme decoding with phonetic lm and context-independent units only

	topn_beam  = "-topn_beam"
	nfilt      = "-nfilt"
	transform  = "-transform"
	compallsen = "-compallsen"
	maxwpf     = "-maxwpf"
	lpbeam     = "-lpbeam"     //Beam width applied to last phone in words
	lponlybeam = "-lponlybeam" //Beam width applied to last phone in single-phone words

	silprob = "-silprob" //		0.005		Silence word transition probability

	// New acoustic model
	hmm = "-hmm"
	lm  = "-lm"
)

type psPhonemeSettings map[psFlag]string

/*
// Default pockectsphinx settings
var template = psPhonemeSettings {
  alpha: "0.97",
  backtrace: "yes",
  beam: "1e-48",
  bestpath: "no",					// default is yes, however, all online stuff and papers turn to no
  dict: "",
  dither: "no",
  fsgusefiller: "no",
  fwdflat: "yes",
  infile: "",
  jsgf: "",
  lw: "6.5",
  maxhmmpf: "30000",
  nfft: "512",
  nwpen: "1",
  pbeam: "1e-48",
  pip: "1",
  remove_dc: "no",
  remove_noise: "yes",
  remove_silence: "yes",
  topn: "4",
  wbeam: "7e-29",
  wip: "0.65",
  wlen: "0.025625",
}
*/

/*
Preemphasizer

Implements a high-pass filter that compensates for attenuation in the audio data. Speech signals have an attenuation (a decrease in intensity of a signal) of 20 dB/dec.
It increases the relative magnitude of the higher frequencies with respect to the lower frequencies.
The Preemphasizer takes a Dataobject that usually represents audio data as input, and outputs the same Dataobject, but with preemphasis applied.
For each value X[i] in the input Data object X, the following formula is applied to obtain the output Data object Y:

Y[i] = X[i] - (X[i-1] * preemphasisFactor)

where 'i' denotes time.

The preemphasis factor has a default defined by PROP_PREEMPHASIS_FACTOR_DEFAULT. A common value for this factor is something around 0.97.

Other Dataobjects are passed along unchanged through this Preemphasizer.

The Preemphasizer emphasizes the high frequency components, because they usually contain much less energy than lower frequency components, even though they are still important for speech recognition.
It is a high-pass filter because it allows the high frequency components to "pass through", while weakening or filtering out the low frequency components.


  pl_window: "0",					//		5		Phoneme lookahead window size, in frames
  									// According to https://cmusphinx.github.io/wiki/pocketsphinxhandhelds/   larger values gives fast decode but lower accuracy
  lpbeam: "1e-10000",				//1e-40		Beam width applied to last phone in words
  lponlybeam: "1e-10000",			//7e-29		Beam width applied to last phone in single-phone words

*/

type psSuite []psPhonemeSettings

//type psSuite []configSettings

var settingsCatalgue = []psPhonemeSettings{

	//======================================================================
	//   ___         _          ___  __
	//  | __| _ __ _| |_ ___   / _ \/  \
	//  | _| '_/ _` |  _/ -_)  \_, / () |
	//  |_||_| \__,_|\__\___|   /_/ \__/
	//
	//======================================================================

	{
		// hmm:        "/public_storage/pronounce/model",
		cmn:        "live",
		cmninit:    "",
		featparams: "",
		hmm:        "",

		backtrace:    "yes",
		bestpath:     "no",
		dict:         "",
		fsgusefiller: "no",
		fwdflat:      "no",
		infile:       "",
		jsgf:         "",
		nwpen:        "1",

		remove_dc:      "no",
		remove_noise:   "yes",
		remove_silence: "no", // Turning VAD on at this low Frate seems to create a mess???

		alpha: "0.97",
		//dither: "yes",
		//dither: "no",

		//maxhmmpf: "3",
		//maxwpf: "3",

		maxhmmpf: "-1",
		maxwpf:   "-1",

		beam:  "1e-10000",
		pbeam: "1e-10000",
		wbeam: "1e-10000",

		frate: "125",
		wlen:  "0.016",
		nfft:  "256",

		lw:       "6",
		pip:      "1",
		wip:      "0.5",
		topn:     "4",
		dither:   "no",
		doublebw: "no",
		/*
		   //lw: "99",
		   lw: "17",     //17
		   //lw: "15",
		   //lw: "9",
		   //lw: "6.5",
		   //lw: "5",
		   //lw: "4.2",
		   //lw: "3.14159",
		   //lw: "3.3",
		   //lw: "2.0",
		   //lw: "1.5",
		   //lw: "5",
		   //lw: "5",

		   //pip: "1e-12",
		   //pip: "1e-6",
		   //pip: "1e-5",
		   //pip: "1e-4",
		   //pip: "1e-3",
		   //pip: "1e-2",
		   //pip: "1e-1",
		   //pip: "9090",
		   pip: "1.15",
		   //pip: "1",

		   //wip: "1e-5",
		   //wip: "1e-4",
		   //wip: "1e-3",
		   //wip: "1e-2",
		   //wip: "1e-1",
		   wip: "0.25",
		   //wip: "1",
		   //wip: "9000",
		   //wip: "1",

		   //topn: "4",			//8
		   topn: "6",
		   //topn: "16",

		   doublebw: "yes",
		   //doublebw: "no",
		*/

		pl_window: "0",

		lpbeam:     "1e-10000",
		lponlybeam: "1e-10000",

		// These numbers have been adjusted from the default by a factor of 100(default)/80(new Frate)=1.25
		vad_postspeech:  "25", // Default = 50   Num of silence frames to keep after from speech to silence.  	(50x1.25=62.5)  (50/1.25=40)
		vad_prespeech:   "5",  // Default = 20   Num of speech frames to keep before silence to speech.		(20x1.25=25)    (20/1.25=16)
		vad_startspeech: "8",  // Default = 10   Num of speech frames to trigger vad from silence to speech.	(10x1.25=12.5)  (10/1.25=8)
		vad_threshold:   "1",  // Default is 3   Threshold for decision between noise and silence frames. Log-ratio between signal level and noise level.  (Using 3.9 causes a lot of crashes - nothing returned)

		// One frame = 512 samples or 0.032 seconds  (window length)
		// Start speech of 10 frames * 0.032 = 0.320 seconds            (@F190 - Start speech of 21 frames * 0.016 = 0.336 seconds)
		// Post speech - gap for things like T -  16*0.032=0.512                      (@F190 - 30*0.016=0.48 seconds)

	},

	//======================================================================
	//   ___         _          _ _  __
	//  | __| _ __ _| |_ ___   / / |/  \
	//  | _| '_/ _` |  _/ -_)  | | | () |
	//  |_||_| \__,_|\__\___|  |_|_|\__/
	//
	//======================================================================
	{
		// hmm:        "/public_storage/pronounce/model",
		cmn:        "live",
		cmninit:    "",
		featparams: "",
		hmm:        "",

		backtrace:    "yes",
		bestpath:     "no",
		dict:         "",
		fsgusefiller: "no",
		fwdflat:      "no",
		infile:       "",
		jsgf:         "",
		nwpen:        "1",

		remove_dc:      "no",
		remove_noise:   "yes",
		remove_silence: "no",

		alpha: "0.97",
		//dither: "yes",

		maxhmmpf: "-1",
		maxwpf:   "-1",

		beam:  "1e-10000",
		pbeam: "1e-10000",
		wbeam: "1e-10000",

		frate: "105", //105
		wlen:  "0.020",
		nfft:  "512",

		lw:       "6",
		pip:      "1",
		wip:      "0.5",
		topn:     "4",
		dither:   "no",
		doublebw: "no",

		/*
		     //lw: "25",
		     lw: "17",
		     //lw: "9",
		     //lw: "6.5",
		     //lw: "5.0",
		     //lw: "4.5",
		     //lw: "3.14159",
		     //lw: "2.0",
		     //lw: "1.5",
		     //lw: "1.05",
		     //lw: "5",

		     //pip: "1e-12",
		     //pip: "1e-6",
		     //pip: "1e-5",
		     //pip: "1e-4",
		     //pip: "1e-3",
		     //pip: "1e-2",
		     //pip: "5e-1",
		     //pip: "1e-1",
		     //pip: "9.9999",
		     //pip: "1.5",
		     pip: "1.15",     //1.15

		     //wip: "1e-5",
		     //wip: "1e-4",
		     //wip: "1e-3",
		     //wip: "1e-2",
		     //wip: "1e-1",
		     wip: "0.25",
		     //wip: "0.1",
		     //wip: "1",
		     //wip: "0.5",

		     topn: "7",
		     //topn: "64",


		   doublebw: "yes",
		   //doublebw: "no",
		*/

		pl_window: "0",

		lpbeam:     "1e-10000",
		lponlybeam: "1e-10000",

		// These numbers have been adjusted from the default by a factor of 100(default)/105(new Frate)
		vad_postspeech:  "20",  // default =50
		vad_prespeech:   "5",   //Num of speech frames to keep before silence to speech. default is 20.
		vad_startspeech: "5",   // default = 10
		vad_threshold:   "1.5", // Default is 3   Threshold for decision between noise and silence frames. Log-ratio between signal level and noise level.  (Using 3.9 causes a lot of crashes - nothing returned)

	},

	//======================================================================
	//  ___         _          _ ___ __
	// | __| _ __ _| |_ ___   / |_  )  \
	// | _| '_/ _` |  _/ -_)  | |/ / () |
	// |_||_| \__,_|\__\___|  |_/___\__/
	//
	//======================================================================
	{
		// hmm:        "/public_storage/pronounce/model",
		cmn:        "live",
		cmninit:    "",
		featparams: "",
		hmm:        "",

		backtrace:    "yes",
		bestpath:     "no",
		dict:         "",
		fsgusefiller: "no",
		fwdflat:      "no",
		infile:       "",
		jsgf:         "",
		nwpen:        "1",

		remove_dc:    "no",
		remove_noise: "yes",
		//remove_silence: "yes",
		remove_silence: "no",

		alpha: "0.97",
		//dither: "yes",

		maxhmmpf: "-1",
		maxwpf:   "-1",

		beam:  "1e-10000",
		pbeam: "1e-10000",
		wbeam: "1e-10000",

		frate: "91", //91
		wlen:  "0.024",
		nfft:  "512",

		lw:       "6",
		pip:      "1",
		wip:      "0.5",
		topn:     "4",
		dither:   "no",
		doublebw: "no",

		/*
		     //lw: "25",
		     //lw: "19",
		     lw: "17",
		     //lw: "7",
		     //lw: "5",
		     //lw: "3.14159",
		     //lw: "2.0",
		     //lw: "1.5",
		     //lw: "5",

		     //pip: "1e-12",
		     //pip: "1e-6",
		     //pip: "1e-5",
		     //pip: "1e-4",
		     //pip: "1e-3",
		     //pip: "5e-2",
		     //pip: "9e-2",
		     //pip: "1e-2",
		     //pip: "1e-1",
		     pip: "1.15",
		     //pip: "1.0",
		     //pip: "0.7",
		     //pip: "0.6",
		     //pip: "9",

		     //wip: "1e-5",
		     //wip: "1e-4",
		     //wip: "1e-3",
		     //wip: "1e-2",
		     //wip: "1e-1",
		     wip: "0.25",
		     //wip: "1",
		     //wip: "0.5",

		     //topn: "4",
		     topn: "6",



		   doublebw: "yes",
		   //doublebw: "no",
		*/

		pl_window: "0",

		lpbeam:     "1e-10000",
		lponlybeam: "1e-10000",

		// These numbers have been adjusted from the default by a factor of 100(default)/190(new Frate)
		vad_postspeech:  "20",  // Default = 50   Num of silence frames to keep after from speech to silence.  	(50x0.526=26.3)  (50/0.526=95.06)
		vad_prespeech:   "5",   // Default = 20   Num of speech frames to keep before silence to speech.		(20x0.526=10.53)    (20/0.526=38.02)
		vad_startspeech: "5",   // Default = 10   Num of speech frames to trigger vad from silence to speech.	(10x0.526=5.25)  (10/1.25=19.01)
		vad_threshold:   "0.5", // Default is 3   Threshold for decision between noise and silence frames. Log-ratio between signal level and noise level.  (Using 3.9 causes a lot of crashes - nothing returned)

	},

	//======================================================================
	//   ___         _           _ _ _  ___
	//  | __| _ __ _| |_ ___    / | | || __|
	//  | _| '_/ _` |  _/ -_)   | |_  _|__ \
	//  |_||_| \__,_|\__\___|   |_| |_||___/
	//
	//======================================================================
	{
		// hmm:        "/public_storage/pronounce/model",
		cmn:        "live",
		cmninit:    "",
		featparams: "",
		hmm:        "",

		backtrace:    "yes",
		bestpath:     "no",
		dict:         "",
		fsgusefiller: "no",
		fwdflat:      "no",
		infile:       "",
		jsgf:         "",
		nwpen:        "1",

		remove_dc:      "no",
		remove_noise:   "yes",
		remove_silence: "yes",

		alpha:  "0.97",
		dither: "yes",

		//maxhmmpf: "3",
		//maxwpf: "-1",

		maxhmmpf: "-1",
		maxwpf:   "-1",

		beam:  "1e-10000",
		pbeam: "1e-10000",
		wbeam: "1e-10000",

		frate: "80",
		wlen:  "0.028",
		nfft:  "512",

		//lw: "25",
		lw: "6",
		//lw: "9.0",
		//lw: "6.5",
		//lw: "5.1",
		//lw: "3.14159",
		//lw: "2.0",
		//lw: "1.5",

		//pip: "1e-12",
		//pip: "1e-6",
		//pip: "1e-5",
		//pip: "1e-4",
		//pip: "1e-3",
		//pip: "1e-2",
		//pip: "1e-1",
		pip: "1.15",

		//wip: "1e-5",
		//wip: "1e-4",
		//wip: "1e-3",
		//wip: "5e-2",
		//wip: "1e-1",
		wip: "0.25",
		//wip: "1",

		//topn: "4",
		topn: "6",

		pl_window: "0",

		lpbeam:     "1e-10000",
		lponlybeam: "1e-10000",

		//doublebw: "no",
		doublebw: "yes",

		vad_postspeech:  "20",  // Default = 50   Num of silence frames to keep after from speech to silence.  	(50x0.526=26.3)  (50/0.526=95.06)
		vad_prespeech:   "5",   // Default = 20   Num of speech frames to keep before silence to speech.		(20x0.526=10.53)    (20/0.526=38.02)
		vad_startspeech: "5",   // Default = 10   Num of speech frames to trigger vad from silence to speech.	(10x0.526=5.25)  (10/1.25=19.01)
		vad_threshold:   "1.5", // Default is 3   Threshold for decision between noise and silence frames. Log-ratio between signal level and noise level.  (Using 3.9 causes a lot of crashes - nothing returned)

	},

	//======================================================================
	//  ___         _          _ ___  __
	// | __| _ __ _| |_ ___   / / _ \/  \
	// | _| '_/ _` |  _/ -_)  | \_, / () |
	// |_||_| \__,_|\__\___|  |_|/_/ \__/
	//
	//======================================================================
	{
		// hmm:        "/public_storage/pronounce/model",
		cmn:        "live",
		cmninit:    "",
		featparams: "",
		hmm:        "",

		backtrace:    "yes",
		bestpath:     "no",
		dict:         "",
		fsgusefiller: "no",
		fwdflat:      "no",
		infile:       "",
		jsgf:         "",
		nwpen:        "1",

		remove_dc:      "no",
		remove_noise:   "yes",
		remove_silence: "yes",

		alpha:  "0.97",
		dither: "yes",

		//maxhmmpf: "3",
		//maxwpf: "-1",

		maxhmmpf: "-1",
		maxwpf:   "-1",

		beam:  "1e-10000",
		pbeam: "1e-10000",
		wbeam: "1e-10000",

		frate: "72",
		wlen:  "0.032",
		nfft:  "512",

		// One frame = 256 samples (audio sample rate * window length = 16kHz * 0.016 = 256 samples long)

		//lw: "25",
		lw: "6",
		//lw: "9.0",
		//lw: "7.0",
		//lw: "6.5",
		//lw: "4.0",
		//lw: "3.14159",
		//lw: "2.0",
		//lw: "1.5",
		//lw: "2",

		//pip: "1e-12",
		//pip: "1e-6",
		//pip: "1e-5",
		//pip: "1e-4",
		//pip: "1e-3",
		//pip: "1e-2",
		//pip: "1e-1",
		pip: "1.15",
		//pip: "9000",

		//wip: "1e-6",
		//wip: "1e-5",
		//wip: "1e-4",
		//wip: "1e-3",
		//wip: "1e-2",
		//wip: "1e-1",
		//wip: "150",
		wip: "0.25",
		//wip: "1",

		//topn: "64",
		topn: "6",

		pl_window: "0",

		lpbeam:     "1e-10000", // Beam width applied to last phone in words  default 1e-40
		lponlybeam: "1e-10000",

		//doublebw: "no",
		doublebw: "yes",
		//nfilt: "0",

		// These numbers have been adjusted from the default by a factor of 100(default)/190(new Frate)
		vad_postspeech:  "20",  // Default = 50   Num of silence frames to keep after from speech to silence.  	(50x0.526=26.3)  (50/0.526=95.06)
		vad_prespeech:   "5",   // Default = 20   Num of speech frames to keep before silence to speech.		(20x0.526=10.53)    (20/0.526=38.02)
		vad_startspeech: "5",   // Default = 10   Num of speech frames to trigger vad from silence to speech.	(10x0.526=5.25)  (10/1.25=19.01)
		vad_threshold:   "1.5", // Default is 3   Threshold for decision between noise and silence frames. Log-ratio between signal level and noise level.  (Using 3.9 causes a lot of crashes - nothing returned)

		// One frame = 256 samples or 0.016 seconds
		// Start speech of 21 frames * 0.016 = 0.336 seconds
		// Post speech - gap for things like T - 30*0.016=0.48 seconds

	},

	// #########################################################################################################################################
	// #########################################################################################################################################
	// #########################################################################################################################################

	/*

	   //======================================================================
	   //   ___         _          ___  __
	   //  | __| _ __ _| |_ ___   / _ \/  \
	   //  | _| '_/ _` |  _/ -_)  \_, / () |
	   //  |_||_| \__,_|\__\___|   /_/ \__/
	   //
	   //======================================================================


	     {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "140",
	     wlen: "0.025",
	     nfft: "512",


	     topn: "6",
	     //topn: "2",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     },




	   //======================================================================
	   //   ___         _          _ _  __
	   //  | __| _ __ _| |_ ___   / / |/  \
	   //  | _| '_/ _` |  _/ -_)  | | | () |
	   //  |_||_| \__,_|\__\___|  |_|_|\__/
	   //
	   //======================================================================
	     {
	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",


	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "120",
	     wlen: "0.025",
	     nfft: "512",

	     topn: "7",
	     //topn: "64",

	     pl_window: "0",
	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     },






	   //======================================================================
	   //  ___         _          _ ___ __
	   // | __| _ __ _| |_ ___   / |_  )  \
	   // | _| '_/ _` |  _/ -_)  | |/ / () |
	   // |_||_| \__,_|\__\___|  |_/___\__/
	   //
	   //======================================================================
	     {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "100",
	     wlen: "0.025",
	     nfft: "512",



	     topn: "6",      			// 4 and time/climb works but others fail (pest, greatest, blair) ..... 6 and time/climb fails

	     pl_window: "0",
	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",
	     },




	   //======================================================================
	   //   ___         _           _ _ _  ___
	   //  | __| _ __ _| |_ ___    / | | || __|
	   //  | _| '_/ _` |  _/ -_)   | |_  _|__ \
	   //  |_||_| \__,_|\__\___|   |_| |_||___/
	   //
	   //======================================================================
	   {
	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "80",
	     wlen: "0.025",
	     nfft: "512",


	     topn: "6",

	     pl_window: "0",
	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",
	     },





	   //======================================================================
	   //  ___         _          _ ___  __
	   // | __| _ __ _| |_ ___   / / _ \/  \
	   // | _| '_/ _` |  _/ -_)  | \_, / () |
	   // |_||_| \__,_|\__\___|  |_|/_/ \__/
	   //
	   //======================================================================
	   {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "60",
	     wlen: "0.025",
	     nfft: "512",


	     topn: "6",

	     pl_window: "0",
	     lpbeam: "1e-10000",		// Beam width applied to last phone in words  default 1e-40
	     lponlybeam: "1e-10000",

	     },
	*/

	// #########################################################################################################################################
	// #########################################################################################################################################
	// #########################################################################################################################################

	/*

	   //======================================================================
	   //   ___         _          ___  __
	   //  | __| _ __ _| |_ ___   / _ \/  \
	   //  | _| '_/ _` |  _/ -_)  \_, / () |
	   //  |_||_| \__,_|\__\___|   /_/ \__/
	   //
	   //======================================================================
	     {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",
	     nwpen: "1",

	     remove_dc: "no",
	     remove_noise: "yes",
	     remove_silence: "yes",   // Turning VAD on at this low Frate seems to create a mess???

	     alpha: "0.97",
	     dither: "yes",

	     //maxhmmpf: "3",
	     //maxwpf: "-1",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "90",
	     wlen: "0.032",
	     nfft: "512",

	     //lw: "30",
	     //lw: "15",
	     //lw: "10",
	     //lw: "6.5",
	     //lw: "4.2",
	     lw: "3.14159",
	     //lw: "2.0",
	     //lw: "1.5",
	     //lw: "1.05",
	     //lw: "1",

	     //pip: "1e-12",
	     //pip: "1e-6",
	     //pip: "1e-5",
	     //pip: "1e-4",
	     //pip: "1e-3",
	     //pip: "1e-2",
	     pip: "1e-2",
	     //pip: "1",
	     //pip: "0.5",

	     //wip: "1e-5",
	     //wip: "1e-4",
	     //wip: "1e-3",
	     //wip: "1e-2",
	     wip: "1e-2",
	     //wip: "0.6",
	     //wip: "1",

	     //topn: "32",
	     topn: "7",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     doublebw: "yes",
	     //doublebw: "no",

	     },





	   //======================================================================
	   //   ___         _          _ _  __
	   //  | __| _ __ _| |_ ___   / / |/  \
	   //  | _| '_/ _` |  _/ -_)  | | | () |
	   //  |_||_| \__,_|\__\___|  |_|_|\__/
	   //
	   //======================================================================
	     {
	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",
	     nwpen: "1",

	     remove_dc: "no",
	     remove_noise: "yes",
	     remove_silence: "no",

	     alpha: "0.97",
	     dither: "yes",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "105",
	     wlen: "0.025",
	     nfft: "512",

	     //lw: "15",
	     //lw: "11",
	     //lw: "6.5",
	     //lw: "5.0",
	     //lw: "4.5",
	     lw: "3.14159",
	     //lw: "2.0",
	     //lw: "1.5",
	     //lw: "1.05",
	     //lw: "1",

	     //pip: "1e-12",
	     //pip: "1e-6",
	     //pip: "1e-5",
	     //pip: "1e-4",
	     //pip: "1e-3",
	     pip: "1e-2",
	     //pip: "1e-1",
	     //pip: "1",
	     //pip: "0.9",

	     //wip: "1e-5",
	     //wip: "1e-4",
	     //wip: "1e-3",
	     wip: "1e-2",
	     //wip: "1e-1",
	     //wip: "0.1",
	     //wip: "1",
	     //wip: "0.9",

	     topn: "7",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     //doublebw: "yes",
	     doublebw: "no",


	     // These numbers have been adjusted from the default by a factor of 100(default)/105(new Frate)
	     //vad_postspeech: "48",		// default =50
	     //vad_prespeech: "20",      //Num of speech frames to keep before silence to speech. default is 20.
	     //vad_startspeech: "10",	   // default = 10

	     },




	   //======================================================================
	   //  ___         _          _ ___ __
	   // | __| _ __ _| |_ ___   / |_  )  \
	   // | _| '_/ _` |  _/ -_)  | |/ / () |
	   // |_||_| \__,_|\__\___|  |_/___\__/
	   //
	   //======================================================================
	     {
	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",
	     nwpen: "1",

	     remove_dc: "no",
	     remove_noise: "yes",
	     remove_silence: "yes",

	     alpha: "0.97",
	     dither: "no",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "120",
	     wlen: "0.025",
	     nfft: "512",

	     //lw: "15",
	     //lw: "11",
	     //lw: "6.5",
	     lw: "5.0",
	     //lw: "4.5",
	     //lw: "3.14159",
	     //lw: "2.0",
	     //lw: "1.5",
	     //lw: "1.05",
	     //lw: "1",

	     //pip: "1e-12",
	     //pip: "1e-6",
	     //pip: "1e-5",
	     //pip: "1e-4",
	     //pip: "1e-3",
	     //pip: "1e-2",
	     //pip: "5e-1",
	     //pip: "1e-1",
	     //pip: "1",
	     pip: "0.5",

	     //wip: "1e-5",
	     //wip: "1e-4",
	     //wip: "1e-3",
	     //wip: "1e-2",
	     //wip: "1e-1",
	     wip: "0.5",
	     //wip: "1",
	     //wip: "0.9",

	     topn: "16",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     doublebw: "yes",
	     //doublebw: "no",


	     // These numbers have been adjusted from the default by a factor of 100(default)/105(new Frate)
	     //vad_postspeech: "48",		// default =50
	     //vad_prespeech: "20",      //Num of speech frames to keep before silence to speech. default is 20.
	     //vad_startspeech: "10",	   // default = 10

	     },



	   //======================================================================
	   //   ___         _           _ _ _  ___
	   //  | __| _ __ _| |_ ___    / | | || __|
	   //  | _| '_/ _` |  _/ -_)   | |_  _|__ \
	   //  |_||_| \__,_|\__\___|   |_| |_||___/
	   //
	   //======================================================================
	     {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",
	     nwpen: "1",

	     remove_dc: "no",
	     remove_noise: "yes",
	     remove_silence: "no",

	     alpha: "0.97",
	     dither: "no",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "143",
	     wlen: "0.025",
	     nfft: "512",

	     //lw: "11",
	     //lw: "9",
	       lw: "8.5",
	     //lw: "6.5",
	     //lw: "3.14159",
	     //lw: "2.0",
	     //lw: "1.5",
	     //lw: "1.05",

	     //pip: "1e-12",
	     //pip: "1e-6",
	     //pip: "1e-5",
	     //pip: "1e-4",
	     //pip: "1e-3",
	     //pip: "1e-2",
	     pip: "9e-2",
	     //pip: "1e-1",
	     //pip: "1",
	     //pip: "0.5",
	     //pip: "0.75",

	     //wip: "1e-5",
	     //wip: "1e-4",
	     //wip: "1e-3",
	     wip: "1e-2",
	     //wip: "1e-1",
	     //wip: "0.5",
	     //wip: "5e-2",
	     //wip: "1",

	     //topn: "4",
	     topn: "32",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     doublebw: "yes",
	     //doublebw: "no",



	     // These numbers have been adjusted from the default by a factor of 100(default)/120(new Frate)
	     //vad_postspeech: "42",		// default =50
	     //vad_prespeech: "17",      //Num of speech frames to keep before silence to speech. default is 20.
	     //vad_startspeech: "9",	   // default = 10


	     },


	   //======================================================================
	   //  ___         _          _ ___ __
	   // | __| _ __ _| |_ ___   / |_  )  \
	   // | _| '_/ _` |  _/ -_)  | |/ / () |
	   // |_||_| \__,_|\__\___|  |_/___\__/
	   //
	   //======================================================================
	     {

	     cmn: "live",
	     cmninit: "",
	     featparams: "",

	     backtrace: "yes",
	     bestpath: "no",
	     dict: "",
	     fsgusefiller: "no",
	     fwdflat: "no",
	     infile: "",
	     jsgf: "",
	     nwpen: "1",

	     remove_dc: "no",
	     remove_noise: "yes",
	     remove_silence: "no",

	     alpha: "0.97",
	     dither: "no",

	     maxhmmpf: "-1",
	     maxwpf: "-1",

	     beam: "1e-10000",
	     pbeam: "1e-10000",
	     wbeam: "1e-10000",

	     frate: "190",
	     wlen: "0.016",
	     nfft: "256",

	     //lw: "11",
	     lw: "9",
	     //lw: "6.5",
	     //lw: "3.14159",
	     //lw: "2.0",
	     //lw: "1.5",
	     //lw: "1.05",

	     //pip: "1e-12",
	     //pip: "1e-6",
	     //pip: "1e-5",
	     //pip: "1e-4",
	     //pip: "1e-3",
	     //pip: "5e-2",
	     pip: "1e-2",
	     //pip: "1e-1",
	     //pip: "1",
	     //pip: "0.5",
	     //pip: "0.75",

	     //wip: "1e-5",
	     //wip: "1e-4",
	     //wip: "1e-3",
	     wip: "1e-2",
	     //wip: "1e-1",
	     //wip: "0.5",
	     //wip: "5e-2",
	     //wip: "1",

	     //topn: "8",
	     topn: "32",

	     pl_window: "0",

	     lpbeam: "1e-10000",
	     lponlybeam: "1e-10000",

	     doublebw: "yes",
	     //doublebw: "no",



	     // These numbers have been adjusted from the default by a factor of 100(default)/120(new Frate)
	     //vad_postspeech: "42",		// default =50
	     //vad_prespeech: "17",      //Num of speech frames to keep before silence to speech. default is 20.
	     //vad_startspeech: "9",	   // default = 10


	     },

	*/

}

var defaultSuite = psSuite{
	settingsCatalgue[0],
	settingsCatalgue[1],
	settingsCatalgue[2],
	settingsCatalgue[3],
	settingsCatalgue[4],
}

type psPhonemeDatum struct {
	phoneme
	start, end int
}

type psPhonemeResults struct {
	frate int
	data  []psPhonemeDatum
}

type psConfig struct {
	frates   []int
	settings psPhonemeSettings
	word     string
	phonemes [][]phoneme
	jsgfData neighbours
	tempDir  string
}

type newPsConfig struct {
	settings psSuite
	word     string
	phonemes [][]phoneme
	tRule    R_target
	jsgfData neighbours
	tempDir  string
}

type phonemeVerdict struct {
	psPhonemeDatum
	index      int
	goodBadEtc verdict
}

type phonVerdict struct {
	phon       phoneme
	goodBadEtc verdict
}

//=====================================================================
//  _____    _         _  _     _
// |_   _| _(_)_ __   | \| |___(_)______
//   | || '_| | '  \  | .` / _ \ (_-< -_)
//   |_||_| |_|_|_|_| |_|\_\___/_/__|___|
//
//=====================================================================

func trimNoise(results psPhonemeResults) psPhonemeResults {
	if len(results.data) == 0 {
		// There are no results so just return
		//
		return results
	}
	type nonSilBlock struct {
		from, to int
	}

	type phonemeBlock struct {
		ph       phoneme
		from, to int
	}

	step1Data := []psPhonemeDatum{}
	for _, result := range results.data {

		// Replace really short phonemes with silence
		if result.end-result.start < 2 {
			new := psPhonemeDatum{
				sil,
				result.start,
				result.end,
			}
			step1Data = append(step1Data, new)
		} else {
			step1Data = append(step1Data, result)
		}
	}
	step1aData := []phonemeBlock{}
	// Join adjacent sils
	silReadLast := false
	for i, result := range step1Data {
		if result.phoneme == sil {
			if silReadLast {
				// If silLastRead then step1aData can't be empty so we're safe to
				// access the last element
				step1aData[len(step1aData)-1].to = i
			} else {
				new := phonemeBlock{
					result.phoneme,
					i,
					i,
				}
				step1aData = append(step1aData, new)
			}
			silReadLast = true
		} else {
			new := phonemeBlock{
				result.phoneme,
				i,
				i,
			}
			step1aData = append(step1aData, new)
			silReadLast = false
		}
	}
	step2Data := []phonemeBlock{}
	for _, result := range step1aData {

		// Remove really short silences
		if result.ph == sil && results.data[result.to].end-results.data[result.from].start < 3 {
			continue
		}
		// Also remove long consonants
		//if isConsonant(result.ph) && results.data[result.to].end - results.data[result.from].start > 30 {
		//  continue
		//}

		step2Data = append(step2Data, result)
	}
	step3Data := []phonemeBlock{}
	phLastRead := false
	// Join phonemes
	for _, result := range step2Data {
		if result.ph != sil {
			if phLastRead {
				step3Data[len(step3Data)-1].to = result.to
			} else {
				new := phonemeBlock{
					result.ph,
					result.from,
					result.to,
				}
				step3Data = append(step3Data, new)
			}
			phLastRead = true
		} else {
			step3Data = append(step3Data, result)
			phLastRead = false
		}
	}
	step3aData := []phonemeBlock{}
	for _, result := range step3Data {
		// Remove all phoneme blocks that are < 4 long
		if results.data[result.to].end-results.data[result.from].start < 4 {
			continue
		}
		step3aData = append(step3aData, result)
	}
	if len(step3aData) == 0 {
		// We didn't find any non-sil blocks
		return psPhonemeResults{}
	}
	// Join nonSilBlocks if they're < 10 apart
	nonSilBlocks := []nonSilBlock{}
	for _, block := range step3aData {
		if block.ph == sil {
			continue
		}
		if len(nonSilBlocks) == 0 {
			new := nonSilBlock{
				block.from,
				block.to,
			}
			nonSilBlocks = append(nonSilBlocks, new)
		} else {
			if results.data[block.from].start-results.data[nonSilBlocks[len(nonSilBlocks)-1].to].end < 4 {
				nonSilBlocks[len(nonSilBlocks)-1].to = block.to
			} else {
				new := nonSilBlock{
					block.from,
					block.to,
				}
				nonSilBlocks = append(nonSilBlocks, new)
			}
		}
	}
	// Look for the largest non-sil block
	largest := nonSilBlocks[0]
	size := results.data[largest.to].end - results.data[largest.from].start
	for _, block := range nonSilBlocks {
		if size < results.data[block.to].end-results.data[block.from].start {
			largest = block
			size = results.data[block.to].end - results.data[block.from].start
		}
	}
	return psPhonemeResults{
		results.frate,
		results.data[largest.from : largest.to+1],
	}
}

func (ps psConfig) normalise(results psPhonemeResults) psPhonemeResults {
	ret := psPhonemeResults{
		results.frate,
		[]psPhonemeDatum{},
	}
	if len(results.data) == 0 {
		// There are no results so just return
		//
		return results
	}
	contains := func(xx [][]phoneme, y phoneme) bool {
		for _, x := range xx {
			for _, x1 := range x {
				if x1 == y {
					return true
				}
			}
		}
		return false
	}
	data := []psPhonemeDatum{}
	// First let's gather up all adjacent like phonemes into one big phoneme
	// removing any NULLs as we go
	//
	lastDatum := results.data[0]
	for _, result := range results.data[1:] {
		if result.phoneme == "(NULL)" {
			continue
		}
		if lastDatum.phoneme == result.phoneme {
			lastDatum.end = result.end
			continue
		}
		if ps.jsgfData.isNeighbour(lastDatum.phoneme, result.phoneme) && !contains(ps.phonemes, result.phoneme) {
			lastDatum.end = result.end
			continue
		}
		data = append(data, lastDatum)
		lastDatum = result
	}
	data = append(data, lastDatum)

	factor := 100.0 / float64(ret.frate)
	for _, datum := range data {
		if datum.end-datum.start < 3 {
			// If pocketsphinx_continuous returns the phoneme length as <= 2 then
			// dump the phoneme now is it's likely not really there
			continue
		}
		start := int(math.Round(float64(datum.start) * factor))
		end := int(math.Round(float64(datum.end) * factor))
		newDatum := psPhonemeDatum{
			datum.phoneme,
			start,
			end,
		}
		ret.data = append(ret.data, newDatum)
	}
	return ret
}

func timeNormalise(results psPhonemeResults) psPhonemeResults {
	if len(results.data) == 0 {
		// Nothing to normalise
		return results
	}
	ret := psPhonemeResults{
		results.frate,
		[]psPhonemeDatum{},
	}
	factor := 100.0 / float64(ret.frate)
	lastDatum := results.data[0]
	for _, datum := range results.data {
		start := datum.start
		end := datum.end
		if end-start < 3 {
			// If the duration is < 3 the phoneme's probably not really there.
			// pocketsphinx tends to create short phonemes like this when the grammar
			// insists on it being present but pocketsphinx hasn't really found it
			end = start
		}
		start = int(math.Round(float64(start) * factor))
		end = int(math.Round(float64(end) * factor))
		if datum.phoneme == lastDatum.phoneme &&
			datum.phoneme == sil {

			lastDatum.end = datum.end
		} else {
			ret.data = append(ret.data, lastDatum)
			lastDatum = psPhonemeDatum{
				datum.phoneme,
				start,
				end,
			}
		}
		// ret.data = append(ret.data, newDatum)
	}
	ret.data = append(ret.data, lastDatum)
	return ret
}

func (ps newPsConfig) normalise(results psPhonemeResults, withDiphthongs bool) psPhonemeResults {
	ret := psPhonemeResults{
		results.frate,
		[]psPhonemeDatum{},
	}
	if len(results.data) == 0 {
		// There are no results so just return
		//
		return results
	}
	contains := func(xx [][]phoneme, y phoneme) bool {
		for _, x := range xx {
			for _, x1 := range x {
				if x1 == y {
					return true
				}
			}
		}
		return false
	}
	data := []psPhonemeDatum{}
	// First let's gather up all adjacent like phonemes into one big phoneme
	// removing any NULLs as we go
	//
	lastDatum := results.data[0]
	for _, result := range results.data[1:] {
		if result.phoneme == "(NULL)" {
			continue
		}
		if lastDatum.phoneme == result.phoneme {
			lastDatum.end = result.end
			continue
		}
		if ps.jsgfData.isNeighbour(lastDatum.phoneme, result.phoneme) && !contains(ps.phonemes, result.phoneme) {
			lastDatum.end = result.end
			continue
		}
		data = append(data, lastDatum)
		lastDatum = result
	}
	data = append(data, lastDatum)

	factor := 100.0 / float64(ret.frate)
	for _, datum := range data {
		start := int(math.Round(float64(datum.start) * factor))
		end := int(math.Round(float64(datum.end) * factor))
		newDatum := psPhonemeDatum{
			datum.phoneme,
			start,
			end,
		}
		ret.data = append(ret.data, newDatum)
	}
	return ret
}

func tryMakeDiphthong(phonemes []psPhonemeDatum, firstPh int) (phoneme, bool) {
	if firstPh >= len(phonemes)-1 {
		// Just return a phoneme. It shouldn't be used by the caller becasue we're
		// also returning false
		return aa, false
	}
	diphthong := [2]phoneme{
		phonemes[firstPh].phoneme, phonemes[firstPh+1].phoneme,
	}
	for key, value := range diphthongs {
		if value == diphthong {
			return key, true
		}
	}
	return aa, false
}

func (r psPhonemeResults) clean() psPhonemeResults {
	// Remove NULLs and recombine any diphthongs
	ret := psPhonemeResults{
		r.frate,
		[]psPhonemeDatum{},
	}
	if len(r.data) == 0 {
		// There are no results so just return
		//
		return ret
	}
	data := []psPhonemeDatum{}
	// First let's gather up all adjacent like phonemes into one big phoneme
	// removing any NULLs as we go
	//
	i := 0
	for i < len(r.data) {
		if r.data[i].phoneme == "(NULL)" {
			continue
		}
		if diphthong, ok := tryMakeDiphthong(r.data, i); ok {
			datum := psPhonemeDatum{
				diphthong,
				r.data[i].start, r.data[i+1].end,
			}
			data = append(data, datum)
			i += 2
			continue
		}
		data = append(data, r.data[i])
		i++
	}
	// Normalise the duration of each phoneme later when I'm ready to adjust
	// other areas of code.
	/*
	  factor := 100.0 / float64(ret.frate)
	  for _, datum := range data {
	    if datum.end - datum.start < 3 {
	      // If pocketsphinx_continuous returns the phoneme length as <= 2 then
	      // dump the phoneme now is it's likely not really there
	      continue
	    }
	    start := int(math.Round(float64(datum.start) * factor))
	    end := int(math.Round(float64(datum.end) * factor))
	    newDatum := psPhonemeDatum{
	      datum.phoneme,
	      start,
	      end,
	    }
	    ret.data = append(ret.data, newDatum)
	  }
	*/
	ret.data = data
	return ret
}

type pruneAssist interface {
	resolve(l link, s linkSet) link
}

func (p psPhonemeResults) resolve(l link, s linkSet) link {
	debug("resolve->: l =", l, "s =", s)
	if len(s.links) == 0 {
		// There's nothing to resolve!
		//
		debug("resolve<-:", l)
		return l
	}
	// A local function to determine whether this is a to conflict.
	//
	to := func(x link, y linkSet) bool {
		for y1 := range y.links {
			if y1.to != l.to {
				return false
			}
		}
		return true
	}
	// Decide what kind of conflict needs to be resolved
	//
	if to(l, s) {
		// Pick the earliest expected phoneme - that is the smallest from value
		//
		ret := l
		for k := range s.links {
			if k.from < ret.from {
				ret = k
			}
		}
		debug("resolve<-:", ret)
		return ret
	}
	// if it's not a to conflict then this is either a from conflict or a
	// crossover conflict. Either way pick the link with the largest to range
	//
	ret := l
	diff := p.data[ret.to].end - p.data[ret.to].start
	for k := range s.links {
		thisDiff := p.data[k.to].end - p.data[k.to].start
		if thisDiff > diff {
			ret = k
			diff = thisDiff
		}
	}
	debug("resolve<-:", ret)
	return ret
}

/*
func bestOfThree(wordPhonemes []phoneme, i int, results []psPhonemeResults, linkSets []linkSet) phonemeVerdict {
   // return these lines to return to normal
  ret := phonemeVerdict{}
  if len(results) != 3 {
    return ret
  }
  wordPhoneme := wordPhonemes[i]
  starts := []int{}
  ends := []int{}
  for j := 0; j < len(linkSets); j++ {
    if l, ok := linkSets[j].linkWithFrom(i); ok {
      starts = append(starts, results[j].data[l.to].start)
      ends = append(ends, results[j].data[l.to].end)
    }
  }
  var verdict phonemeVerdict
  switch len(starts) {
  case 0:
    // Should never happen because we insist that pocketsphinx_continuous
    // finds the phonemes we're looking for
    //
    verdict = phonemeVerdict{
      psPhonemeDatum{
        wordPhoneme,
        0,
        0,
      },
      i,
      missing,
    }
  case 1:
    verdict = phonemeVerdict{
      psPhonemeDatum{
        wordPhoneme,
        starts[0],
        ends[0],
      },
      i,
      missing,
    }
  case 2:
    start := max(starts[0], starts[1:]...)
    end := min(ends[0], ends[1:]...)
    diff := end - start
    if diff <= 0 {
      verdict = phonemeVerdict{
        psPhonemeDatum{
          wordPhoneme,
          0,
          0,
        },
        i,
        missing,
      }
    } else {
      verdict = phonemeVerdict{
        psPhonemeDatum{
          wordPhoneme,
          start,
          end,
        },
        i,
        possible,
      }
    }
  case 3:
    start := max(starts[0], starts[1:]...)
    end := min(ends[0], ends[1:]...)
    diff := end - start
    if diff > 0 {
      verdict = phonemeVerdict{
        psPhonemeDatum{
          wordPhoneme,
          start,
          end,
        },
        i,
        good,
      }
    } else {
      starts_k := []int{}
      ends_k := []int{}
      for k := 0; k < len(starts); k++ {
        switch k {
        case 0:
          starts_k = append(starts_k, max(starts[1], starts[2:]...))
          ends_k = append(ends_k, min(ends[1], ends[2:]...))
        case 1:
          starts_k = append(starts_k, max(starts[0], starts[2:]...))
          ends_k = append(ends_k, min(ends[0], ends[2:]...))
        case 2:
          starts_k = append(starts_k, max(starts[0], starts[1:2]...))
          ends_k = append(ends_k, min(ends[0], ends[1:2]...))
        default:
          log.Panic()
        }
      }
      maxDiff := ends_k[0] - starts_k[0]
      maxDiffIndex := 0
      for k := 1; k < len(starts); k++ {
        diff := ends_k[k] - starts_k[k]
        if diff > maxDiff {
          maxDiff = diff
          maxDiffIndex = k
        }
      }
      if maxDiff > 0 {
        verdict = phonemeVerdict{
          psPhonemeDatum{
            wordPhoneme,
            starts_k[maxDiffIndex],
            ends_k[maxDiffIndex],
          },
          i,
          possible,
        }
      } else {
        verdict = phonemeVerdict{
          psPhonemeDatum{
            wordPhoneme,
            0,
            0,
          },
          i,
          missing,
        }
      }
    }
  default:
  }
  return verdict
}

func bestOfFive(wordPhonemes []phoneme, i int, results []psPhonemeResults, linkSets []linkSet) phonemeVerdict {
  // Try to get the best verdict out of the 5 scans for the ith phoneme we're
  // interested in
  candidates := [][]int{
    []int{
      0, 1, 2, 3, 4,
    },
    []int{
      0, 1, 2, 3,
    },
    []int{
      0, 1, 2, 4,
    },
    []int{
      0, 1, 3, 4,
    },
    []int{
      0, 2, 3, 4,
    },
    []int{
      1, 2, 3, 4,
    },
    []int{
      0, 1, 2,
    },
    []int{
      0, 1, 3,
    },
    []int{
      0, 1, 4,
    },
    []int{
      0, 2, 3,
    },
    []int{
      0, 2, 4,
    },
    []int{
      0, 3, 4,
    },
    []int{
      1, 2, 3,
    },
    []int{
      1, 2, 4,
    },
    []int{
      1, 3, 4,
    },
    []int{
      2, 3, 4,
    },
    []int{
      0, 1,
    },
    []int{
      0, 2,
    },
    []int{
      0, 3,
    },
    []int{
      0, 4,
    },
    []int{
      1, 2,
    },
    []int{
      1, 3,
    },
    []int{
      1, 4,
    },
    []int{
      2, 3,
    },
    []int{
      2, 4,
    },
    []int{
      3, 4,
    },
  }
  if len(results) != 5 {
    return phonemeVerdict{}
  }
  wordPhoneme := wordPhonemes[i]
  starts := []int{}
  ends := []int{}
  for j := 0; j < len(linkSets); j++ {
    if l, ok := linkSets[j].linkWithFrom(i); ok {
      starts = append(starts, results[j].data[l.to].start)
      ends = append(ends, results[j].data[l.to].end)
    }
  }
  for _, candidate := range candidates {
    starts_i := []int{}
    ends_i := []int{}
    badCandidate := false
    for j, _ := range candidate {
      if candidate[j] > len(starts) - 1 {
        badCandidate = true
        break
      }
      starts_i = append(starts_i, starts[candidate[j]])
      ends_i = append(ends_i, ends[candidate[j]])
    }
    if badCandidate {
      continue
    }
    start := max(starts_i[0], starts_i[1:]...)
    end := min(ends_i[0], ends_i[1:]...)
    if end - start > 0 {
      var thisVerdict verdict
      if len(candidate) >= 3 {
        thisVerdict = good
      }
      if len(candidate) == 2 {
        thisVerdict = possible
      }
      // Should never happen
      if len(candidate) <= 1 {
        thisVerdict = missing
      }
      return phonemeVerdict{
        psPhonemeDatum{
          wordPhoneme,
          start,
          end,
        },
        i,
        thisVerdict,
      }
    }
  }
  return phonemeVerdict{
    psPhonemeDatum{
      wordPhoneme,
      0,
      0,
    },
    i,
    missing,
  }
}
*/

func phonemeAt(i int, results psPhonemeResults) psPhonemeDatum {
	for _, result := range results.data {
		if i <= result.end && i >= result.start {
			return result
		}
	}
	return psPhonemeDatum{}
}

func surpriseInrange(from, to int, results []psPhonemeResults) []phonemeVerdict {
	ret := []phonemeVerdict{}
	i := from
	for i < to {
		starts := []int{}
		ends := []int{}
		phons := []psPhonemeDatum{}
		for j := 0; j < len(results); j++ {
			datum := phonemeAt(i, results[j])
			if (psPhonemeDatum{}) == datum {
				continue
			}
			phons = append(phons, datum)
			starts = append(starts, datum.start)
			ends = append(ends, datum.end)
		}
		if len(phons) != len(results) {
			i++
			continue
		}
		start := max(starts[0], starts[1:]...)
		end := min(ends[0], ends[1:]...)
		firstPhon := phons[0]
		if firstPhon.phoneme == "sil" {
			// Not interested in sils so move on
			//
			i = end + 1
			continue
		}
		same := true
		for j := 1; j < len(phons); j++ {
			if phons[j].phoneme != firstPhon.phoneme {
				same = false
				break
			}
		}
		if !same {
			i = end + 1
			continue
		}
		// So the phonemes are the same, all we have to do now is check the
		// interval is big enough
		//
		if end-start >= 3 {
			ret = append(ret, phonemeVerdict{
				psPhonemeDatum{
					phons[0].phoneme,
					start,
					end,
				},
				-1,
				surprise,
			})
		}
		i = end + 1
	}
	return ret
}

func insertSurprises(results []psPhonemeResults, verdicts []phonemeVerdict) []phonemeVerdict {
	ret := []phonemeVerdict{}
	if len(results) == 0 {
		return ret
	}
	from := 0
	for _, verdict := range verdicts {
		if verdict.goodBadEtc == missing && verdict.start == 0 && verdict.end == 0 {
			ret = append(ret, verdict)
			continue
		}
		to := verdict.start - 1
		surprises := surpriseInrange(from, to, results)
		ret = append(ret, surprises...)
		// Now append the verdict
		//
		ret = append(ret, verdict)
		from = verdict.end + 1
	}
	// Look for anything surprising after the last verdict
	//
	// Look for last phoneme in all results
	//
	ends := []int{}
	for _, result := range results {
		// Handle crash at next line. If len(result.data) == 0 this causes a
		// crash!
		//
		dataLen := len(result.data)
		if dataLen == 0 {
			continue
		}
		ends = append(ends, result.data[dataLen-1].end)
	}
	var to int
	switch len(ends) {
	case 0:
		return ret
	case 1:
		to = ends[0]
	default:
		to = min(ends[0], ends[1:]...)
	}
	surprises := surpriseInrange(from, to, results)
	ret = append(ret, surprises...)

	return ret
}

type finalVerdict struct {
	phons []phonemeVerdict
	err   error
}

type LettersVerdict struct {
	Letters string
	Phons   []phoneme
	// index int
	GoodBadEtc verdict
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

func getVerdict(phToAbcEntry phonToAlphas, verdicts []phonVerdict) []LettersVerdict {
	ret := []LettersVerdict{}
	if len(verdicts) == 0 {
		debug("getVerdict<-:len(verdicts) == 0")
		return ret
	}
	// This is straightforward. There's only one phoneme so return the corresponding
	// letter(s), with the corrseponding verdict.
	if len(phToAbcEntry.phons) == 1 {
		new := LettersVerdict{
			phToAbcEntry.alphas,
			phToAbcEntry.phons,
			verdicts[0].goodBadEtc,
		}
		debug("getVerdict<-:new", new)
		return []LettersVerdict{new}
	} else {
		// So here we have a number of phonemes corresponding to a letter (or number of
		// letters). What the hell does the index mean though? It's the index of the phoneme
		// in the word (expected) phonemes
		// We need to loop through the phoneme verdicts until we've gone through
		// all the phonemes in phToAbcEntry
		last := len(phToAbcEntry.phons) - 1
		expected := LettersVerdict{
			phToAbcEntry.alphas,
			phToAbcEntry.phons,
			good,
		}
		surprises := []LettersVerdict{}
		j := 0
		k := 0
		for j < len(verdicts) {
			v := verdicts[j]
			if v.goodBadEtc == surprise {
				// There's a surprise so mark the phoneme as a surprise (it doesn't correspond to
				// any letter so add the string-ified version of the phoneme as the letter) and
				// mark the expected letter(s) as missing
				new := LettersVerdict{
					string(v.phon),
					[]phoneme{
						v.phon,
					},
					surprise,
				}
				surprises = append(surprises, new)
				expected.GoodBadEtc = missing
			} else {
				// Update expected if the verdict is worse than what we already have for
				// expected. What this means is if we have a bunch of letters in this
				// phToAbcEntry the overall verdict is the worst of the lot
				if v.goodBadEtc > expected.GoodBadEtc {
					expected.GoodBadEtc = v.goodBadEtc
				}
				if k == last {
					// We're done
					break
				}
				k++
			}
			j++
		}
		ret = append(ret, expected)
		ret = append(ret, surprises...)
	}
	return ret
}

type internalError struct {
	file string
	line int
}

func newInternalError() internalError {
	ret := internalError{}
	if _, f, l, ok := runtime.Caller(1); ok {
		ret = internalError{
			f,
			l,
		}
	}
	return ret
}

func (i internalError) Error() string {
	return fmt.Sprintf("internal error: %s, %d", i.file, i.line)
}

func publish(word string, phons []phoneme, verdicts []phonVerdict) ([]LettersVerdict, error) {
	debug("publish->")
	// First a quick check to see if more than 50% of the phonemes are bad
	//
	badCount := 0
	for _, verdict := range verdicts {
		if verdict.goodBadEtc == missing {
			badCount++
		}
	}
	phonsToAlphas, err := mapPhToA(phons, word)
	if err != nil {
		return []LettersVerdict{}, err
	}

	debug("phonsToAlphas =", phonsToAlphas)
	lettersVerdicts := []LettersVerdict{}
	i := 0
	j := 0
	for i < len(verdicts) {
		verdict := verdicts[i]
		var newVerdict LettersVerdict
		if verdict.goodBadEtc != surprise {
			if j < len(phonsToAlphas) {
				newVerdicts := getVerdict(phonsToAlphas[j], verdicts[i:])
				lettersVerdicts = append(lettersVerdicts, newVerdicts...)
				j++
				k := 0
				for _, nv := range newVerdicts {
					k += len(nv.Phons)
				}
				i += k
			} else {
				return []LettersVerdict{}, newInternalError()
			}
		} else {
			newVerdict = LettersVerdict{
				string(verdict.phon),
				[]phoneme{
					verdict.phon,
				},
				// i,
				verdict.goodBadEtc,
			}
			lettersVerdicts = append(lettersVerdicts, newVerdict)
			i++
		}
	}
	debug("lettersVerdicts =", lettersVerdicts)

	// Send results off
	//
	// json := toJSON(lettersVerdicts, nil)
	//
	// c := communicator.New()
	// _ = c.Post(json, "pronounce_ingress.php")
	//

	debug("publish<-")
	return lettersVerdicts, nil
}

func ToJSON(results []LettersVerdict, err error) []byte {
	type JSON_result struct {
		Letters  string `json:"letters"`
		Phonemes string `json:"phonemes"`
		Verdict  string `json:"verdict"`
	}
	type JSON_results struct {
		Results  []JSON_result `json:"results"`
		ErrorMsg *string       `json:"err"`
	}
	jResults := []JSON_result{}
	for _, result := range results {
		phons := []string{}
		for _, phon := range result.Phons {
			phons = append(phons, cmubetToIpa[phon])
		}
		jResults = append(jResults, JSON_result{
			result.Letters,
			strings.Join(phons, " "),
			verdicts[result.GoodBadEtc],
		})
	}
	// All this just so I can get Go to print null in the JSON when there's no
	// error to report
	//
	var errStr *string
	if err != nil {
		temp := err.Error()
		errStr = &temp
	}
	out := JSON_results{
		jResults,
		errStr,
	}
	j, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		debug("toJSON: call to MarshalIndent failed. err =", err)
		log.Panic()
	}
	return j
}

// func postResults(c communicator.CommunicatorT, json []byte) {
// 	err := c.Post(json, "pronounce_ingress.php")
// 	if err != nil {
// 		debug("postResults: call to Post failed. err =", err)
// 		log.Panic(err)
// 	}
// }

func min(i int, rest ...int) int {
	min := i
	for j := 0; j < len(rest); j++ {
		if rest[j] < min {
			min = rest[j]
		}
	}
	return min
}

func max(i int, rest ...int) int {
	max := i
	for j := 0; j < len(rest); j++ {
		if rest[j] > max {
			max = rest[j]
		}
	}
	return max
}

func phonemesForWord(dict dictionary.Dictionary, word string) ([]phoneme, error) {
	ret := []phoneme{}

	phStrs, err := dict.Phonemes(word)
	if err != nil {
		return ret, err
	}
	for _, phStr := range phStrs {
		ret = append(ret, phoneme(phStr))
	}
	return ret, err
}

func allMissing(verdicts []LettersVerdict) bool {
	for _, v := range verdicts {
		if v.GoodBadEtc != missing {
			return false
		}
	}
	return true
}

func score(verdicts []phonVerdict) int {
	debug("score->")
	ret := 0
	allGood := true
	for _, verdict := range verdicts {
		switch verdict.goodBadEtc {
		case good:
			debug("+1 for", verdict.phon)
			ret += 1 // ret += 2
		case possible:
			debug("+1 for", verdict.phon)
			ret += 1
		case missing:
			debug("-1 for", verdict.phon)
			allGood = false
			ret += -1
		case surprise:
			debug("-1 for", verdict.phon)
			allGood = false
			ret += -1
		}
	}
	if allGood {
		ret += 5
	}
	debug("score<-:", ret)
	return ret
}

type span struct {
	start, end int
}

func diphthongsInWord(phons []phoneme) bool {
	for _, ph := range phons {
		if isDiphthong(ph) {
			return true
		}
	}
	return false
}

// var trim span

//  _____    _         ___         _
// |_   _| _(_)_ __   | __| _ __ _| |_ ___
//   | || '_| | '  \  | _| '_/ _` |  _/ -_)
//   |_||_| |_|_|_|_| |_||_| \__,_|\__\___|
//

func trimResults(results []psPhonemeResults, trim span) []psPhonemeResults {
	ret := []psPhonemeResults{}
	for _, result := range results {
		data := result.data
		trimmed := []psPhonemeDatum{}
		for _, datum := range data {
			if datum.end < trim.start || datum.start > trim.end {
				// Throw the whole datum away
				continue
			}
			if datum.start < trim.start {
				// Trim the start of the datum
				newDatum := psPhonemeDatum{
					datum.phoneme,
					trim.start,
					datum.end,
				}
				trimmed = append(trimmed, newDatum)
				continue
			}

			// Should work regardless of Trim
			//   ___                           _                           _                 _                  _         __                      _
			//   | _ \___ _ __  _____ _____    | |___ _ _  __ _    _ _  ___(_)___ ___    __ _| |_    ___ _ _  __| |   ___ / _|  __ __ _____ _ _ __| |
			//   |   / -_) '  \/ _ \ V / -_)   | / _ \ ' \/ _` |  | ' \/ _ \ (_-</ -_)  / _` |  _|  / -_) ' \/ _` |  / _ \  _|  \ V  V / _ \ '_/ _` |
			//   |_|_\___|_|_|_\___/\_/\___|   |_\___/_||_\__, |  |_||_\___/_/__/\___|  \__,_|\__|  \___|_||_\__,_|  \___/_|     \_/\_/\___/_| \__,_|
			//                                              /_/

			/*
				         if datum.phoneme == b && (datum.end - datum.start) > 10 {
				          continue
				        }


						 if datum.phoneme == l && datum.end - datum.start > 20 {
				          continue
				        }

				         if datum.phoneme == ow && datum.end - datum.start > 25 {
				          continue
				        }

				      	 if datum.phoneme == k && datum.end - datum.start > 20 {
				          continue
				        }
			*/

			if datum.end > trim.end {
				// A bit of a hack but Paul wants to kill off long p inserts at the
				// end of an array of phonemes
				if isConsonant(datum.phoneme) && datum.end-datum.start > 20 {
					continue
				}
				// Trim the end of the datum
				newDatum := psPhonemeDatum{
					datum.phoneme,
					datum.start,
					trim.end,
				}
				trimmed = append(trimmed, newDatum)
				continue
			}
			// No need to trim this datum
			trimmed = append(trimmed, datum)
		}
		trimmedResult := result
		trimmedResult.data = trimmed
		ret = append(ret, trimmedResult)
	}
	return ret
}

type variantResult struct {
	score int
	// verdict []phonemeVerdict
	verdict []phonVerdict
	phons   []phoneme
}

func runPs(logfile string, withSettings psPhonemeSettings) []psPhonemeDatum {
	args := []string{}
	for k, v := range withSettings {
		args = append(args, string(k))
		args = append(args, v)
	}
	args = append(args, string(logfn))
	args = append(args, logfile)
	// The output bytes are useless! What we want is in the logfile, so ignore
	// output bytes and parse the logfile
	//
	_, err := exec.Command("pocketsphinx_continuous", args...).Output()
	if err != nil {
		debug("Oops, check pocketsphinx settings? args are...", args)
		return []psPhonemeDatum{}
	}
	return parsePsData(logfile)
}

// NOTE:
// For reasons best known to the pocketsphinx developers the data in the log
// file shows that pocketsphinx sometimes has more than one go at getting the
// phonemes in a word. So once we've found a word we need to carry on looking
// through the file in case there are more attempts.
func parsePsData(filename string) []psPhonemeDatum {
	phonemeData := []psPhonemeDatum{}
	candidateData := [][]psPhonemeDatum{}

	f, err := os.Open(filename)
	if err != nil {
		debug("parsePsData: failed to open file,", filename, ". err =", err)
		log.Panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	phonemesFound := false
	for s.Scan() {
		l := s.Text()
		tokens := strings.Fields(l)
		if len(tokens) == 0 {
			continue
		}
		if tokens[0] == "word" {
			phonemeData = []psPhonemeDatum{}
			phonemesFound = true
			continue
		}
		if phonemesFound {
			if tokens[0] == "INFO:" {
				candidateData = append(candidateData, phonemeData)
				phonemesFound = false
				continue
			} else {
				if tokens[0] == "(NULL)" {
					// Not interested in NULLs so move to the next line
					//
					continue
				}
				start, err := strconv.Atoi(tokens[1])
				if err != nil {
					// Something's gone wrong so just return nothing
					//
					return []psPhonemeDatum{}
				}
				end, err := strconv.Atoi(tokens[2])
				if err != nil {
					return []psPhonemeDatum{}
				}
				datum := psPhonemeDatum{
					phoneme(tokens[0]),
					start,
					end,
				}
				phonemeData = append(phonemeData, datum)
			}
		}
	}
	if phonemesFound == true {
		// The logfile ended with the last phoneme of the word (so no INFO: line)
		// so add the last result to candidateData
		candidateData = append(candidateData, phonemeData)
	}
	// Now work out which of the possibly several blocks of phoneme data is the
	// one we want. For now, we pick the longest (in time) one.
	if len(candidateData) == 1 {
		return candidateData[0]
	}
	bestDataIndex := -1
	maxDuration := 0
	for i, psData := range candidateData {
		duration := 0
		for _, psDatum := range psData {
			if psDatum.phoneme != sil {
				duration += psDatum.end - psDatum.start
			}
		}
		if duration > maxDuration {
			maxDuration = duration
			bestDataIndex = i
		}
	}
	if bestDataIndex != -1 {
		return candidateData[bestDataIndex]
	}
	return []psPhonemeDatum{}
}

func mkDir(dirname string) {
	err := os.Mkdir(dirname, 0777)
	if os.IsPermission(err) {
		debug("mkDir: failed to make directory. err =", err)
		log.Panic()
	}
	if os.IsExist(err) {
		/*
		   // Clear out the existing contents of the folder.
		   d, err := os.Open(dirname)
		   if err != nil {
		     log.Panic()
		   }
		   defer d.Close()

		   names, err := d.Readdirnames(-1)
		   if err != nil {
		     debug("Error", err, "found on reading files and folders")
		   }
		   for _, name := range names {
		     err = os.RemoveAll(path.Join(dirname, name))
		     if err != nil {
		       debug("Error,", err, "removing file/folder,", name)
		     }
		   }
		*/
	}
}

var nearestNeighbours = map[phoneme]phoneme{
	aa:  ih, // was UW, changing to IH
	ae:  iy,
	ah:  eh,
	ao:  ih,
	ax:  ax,
	axr: ax,
	aw:  uw,
	ay:  " ", //ae,			// b - works in time for Paul but ... Not sure if this is a good idea or another vowel should be used to tear apart false AY words?
	// complicated Ay is a transition --> AA IY,  the ending used for other IY is AE but
	b:   " ", // f      was p, but may not be strong enough
	ch:  sh,
	dh:  t,
	d:   v,  // was z but not working in world - tomo, trying v instead.
	eh:  aa, // was iy, which seems ok, but trying aa to avoid creating a diphthong and maybe improve vests
	ehr: aa, // Added 8th March 2021
	er:  ih,
	ey:  ih,
	f:   v,
	g:   k,
	hh:  th,
	ih:  ae, // Seems to work well at splitting voice "seat" text = "sit" to bring out the IY
	ing: ng,
	iy:  ae,
	jh:  ch,
	k:   v,  // originally g, works in conversation
	l:   hh, // originally r, but didn't work for Paulo 'last' ..... k seems ok, but may not be strong enough (p works for e-school but doesn't in Khurrum animal) hh required for animal
	m:   s,  // was ch, which seems ok
	n:   p,  // was m
	//nd: d,
	ng: n,
	oh: oh,
	ow: aa,
	oy: uw,
	p:  t,
	r:  b,
	s:  d, // changed from z to p to try and bring clarity in vests (tomo) but the proceeding t is also plosive, so causes issues, w seems ok-ish
	sh: zh,
	//st: d,
	t:  v,
	th: dh,
	uh: ah,
	uw: aa,
	v:  k, // was v, moved to k in an attempt to recover V in vests for Tomo recording
	w:  b,
	y:  ae, // was w, which seems ok, but trying something strong to try and break up Khurrum - year
	z:  f,  // k seems ok

	axl: l,
	axm: m,
	axn: n,
	ks:  s,
	kw:  k,
	dz:  z,
	uwl: l,
	uwn: uw,
	uwm: uw,
	ts:  s,
	kl:  k,
	pl:  p,
	bl:  b,

	kr: r,
	gr: r,
	tr: r,
	pr: r,

	thr: r,
	ihl: l,
	yuw: y,
	sts: s,
	st:  s,
	ehl: l,
	kt:  t,
	fl:  l,
}

type phonemePair struct {
	p1, p2 phoneme
}

var inserts = map[phonemePair]string{

	// Potential rules
	// =================
	// alternative vowels need sentinel guards
	// alternative consonants do not need guards
	// voiced plosives can be challenged by voiceless plosives without a guard (but not voiced ones without a guard)
	// voiceless plosives should not be challenged by anything
	// soft fade required when word end in ion

	//opening
	//=======

	{sil, aa}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant>  n]|[n <any_consonant>] )             ((",
	{sil, ae}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant>  t]|[t <any_consonant>] )             ((",
	{sil, ah}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant>  n]|[n <any_consonant>] )             ((",
	{sil, ao}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant>  n]|[n <any_consonant>] )             ((",
	{sil, ay}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant>  n]|[n <any_consonant>] )             ((",
	{sil, b}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]       ([<any_consonant_X_b>  n]|[n <any_consonant_X_b>] )     zh           ((",
	{sil, ch}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_ch> n]|[n <any_consonant_X_ch>]) ((",
	{sil, d}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_d>  n]|[n <any_consonant_X_d>] ) ((",
	{sil, dh}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_dh> n]|[n <any_consonant_X_dh>]) ((",
	{sil, eh}: "sil [<any_consonant>]        ((",
	{sil, er}: "sil [<any_consonant>]        ((",
	{sil, ey}: "sil [<any_consonant>]        ((",
	{sil, f}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_f>  n]|[n <any_consonant_X_f>] ) ((",
	{sil, g}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_g>  n]|[n <any_consonant_X_g>] ) ((",
	{sil, hh}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_hh> n]|[n <any_consonant_X_hh>])         t        ((", // extra t added to prevent hh in paulo hardly
	{sil, ih}: "sil [<any_consonant>]        ((",
	{sil, iy}: "sil [<any_consonant>]        ((",
	{sil, jh}: "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_k>  f]|[f <any_consonant_X_k>] ) ((",
	{sil, k}:  "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_k>  n]|[n <any_consonant_X_k>] )           s  ((",
	{sil, l}:  "sil+ [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]        ([<any_consonant_X_l>  n]|[n <any_consonant_X_l>] ) ((",
	{sil, m}:  "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_m>  f]|[f <any_consonant_X_m>] ) ((",
	{sil, n}:  "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_n>  f]|[f <any_consonant_X_n>] ) ((",
	{sil, ng}: "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_ng> f]|[f <any_consonant_X_ng>]) ((",
	{sil, ow}: "sil [<any_consonant>]        ((",
	{sil, oy}: "sil [<any_consonant>]        ((",
	{sil, p}:  "sil [(f (<any_vowel_noSlide>))|((<any_vowel_noSlide>) f)]         ([<any_consonant_X_r>  sh]|[sh <any_consonant_X_r>]) ((",
	{sil, r}:  "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_r>  k]|[n <any_consonant_X_k>] )         b     ((",
	{sil, s}:  "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_s>  n]|[n <any_consonant_X_s>] )         p     ((",
	{sil, sh}: "sil [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ([<any_consonant_X_sh> n]|[n <any_consonant_X_sh>]) ((",
	{sil, t}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_t>  n]|[n <any_consonant_X_t>] ) ((",
	{sil, th}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_th> l]|[l <any_consonant_X_th>]) ((",
	{sil, uh}: "sil [<any_consonant>]        ((",
	{sil, uw}: "sil [<any_consonant>]        ((",
	{sil, v}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_v>  n]|[n <any_consonant_X_v>] ) ((",
	{sil, w}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_w>  n]|[n <any_consonant_X_w>] ) ((",
	{sil, y}:  "sil [<any_consonant>]        ((",
	{sil, z}:  "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_z>  n]|[n <any_consonant_X_z>] ) ((",
	{sil, zh}: "sil [(n (<any_vowel_noSlide>))|((<any_vowel_noSlide>) n)]         ([<any_consonant_X_zh> p]|[p <any_consonant_X_zh>]) ((",

	//ending
	//======

	{aa, sil}: ")|(<any_vowel_Naa> n)|(n <any_vowel_Naa>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{ah, sil}: ")|(<any_vowel_Nah> n)|(n <any_vowel_Nah>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{ay, sil}: ")|(<any_vowel_Nay> f)|(f <any_vowel_Nay>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{d, sil}: ")|(<any_consonant_X_d> ch)|(ch <any_consonant_X_d>) )     	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_d>  f]|[f <any_consonant_X_d>] )       [<soft_fade>] sil",
	{er, sil}: ")|(<any_vowel_Ner> n)|(n <any_vowel_Ner>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{iy, sil}: ")|(<any_vowel_Nih> n)|(n <any_vowel_Nih>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{k, sil}: ")|(<any_consonant_X_k> n)|(n <any_consonant_X_k>) )       	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_k>  f]|[f <any_consonant_X_k>] )       n [<soft_fade>] sil",
	{l, sil}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )     	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_l>  f]|[f <any_consonant_X_l>] )       [<soft_fade>] sil",
	{m, sil}: ")|(<any_consonant_X_m> t)|(t <any_consonant_X_m>) )       	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_m>  f]|[f <any_consonant_X_m>] )       [<soft_fade>] sil",
	{n, sil}: ")|(n n)|(n hh)|(<any_consonant_X_n> t)|(t <any_consonant_X_n>) )       ([<any_vowel_noSlide> ch]|[ch <any_vowel_noSlide>]) ([<any_consonant_X_n>  f]|[f <any_consonant_X_n>] )     k  [<soft_fade>] sil",
	{ow, sil}: ")|(<any_vowel_Now> n)|(n <any_vowel_Now>) )               	([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{r, sil}: ")|(<any_consonant_X_r> n)|(n <any_consonant_X_r>) )        	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_r>  f]|[f <any_consonant_X_r>] )       [<soft_fade>] sil",
	{s, sil}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )        	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_s>  f]|[f <any_consonant_X_s>] )       [<soft_fade>] sil",
	{t, sil}: ")|(<any_consonant_X_t> n)|(n <any_consonant_X_t>) )        	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_t>  f]|[f <any_consonant_X_t>] )       n sil ", //"<soft_fade> sil",           // not clear but sil may need to be removed for Paulo last?
	{y, sil}: ")|(<any_vowel_Ny> n)|(n <any_vowel_Ny>) )         			([<any_consonant>  f]|[f <any_consonant>] )  ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>])      		[<soft_fade>] sil",
	{z, sil}: ")|(<any_consonant_X_z> n)|(n <any_consonant_X_z>) )        	([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_z>  f]|[f <any_consonant_X_z>] )       [<soft_fade>] sil",

	//a into something
	//================
	{aa, ae}: ")|(<any_vowel_Naa_ae> s)|(s <any_vowel_Naa_ae>) )         ((",
	{aa, ah}: ")|(<any_vowel_Naa_ah> s)|(s <any_vowel_Naa_ah>) )         ((",
	{aa, ao}: ")|(<any_vowel_Naa_ao> s)|(s <any_vowel_Naa_ao>) )         ((",
	{aa, aw}: ")|(<any_vowel_Naa_aw> s)|(s <any_vowel_Naa_aw>) )         ((",
	{aa, ay}: ")|(<any_vowel_Naa_ay> s)|(s <any_vowel_Naa_ay>) )         ((",

	{aa, b}:  ")|(<any_vowel_Naa> f)|(f <any_vowel_Naa>) )         ([<any_consonant_X_b>  f]|[f <any_consonant_X_b>] ) ((",
	{aa, ch}: ")|(<any_vowel_Naa> f)|(f <any_vowel_Naa>) )         ([<any_consonant_X_ch> n]|[n <any_consonant_X_ch>]) ((",
	{aa, d}:  ")|(<any_vowel_Naa> s)|(s <any_vowel_Naa>) )         ([<any_consonant_X_d>  f]|[f <any_consonant_X_d>] ) ((",
	{aa, dh}: ")|(<any_vowel_Naa> s)|(s <any_vowel_Naa>) )         ([<any_consonant_X_dh> f]|[f <any_consonant_X_dh>]) ((",

	{aa, l}: ")|(<any_vowel_Naa> s)|(s <any_vowel_Naa>) )         ([<any_consonant_X_l>  f]|[f <any_consonant_X_l>] ) ((",
	{aa, n}: ")|(<any_vowel_Naa> s)|(s <any_vowel_Naa>) )         ([<any_consonant_X_n>  f]|[f <any_consonant_X_n>] ) ((",
	{aa, r}: ")|(<any_vowel_Naa> s)|(s <any_vowel_Naa>) )         ([<any_consonant_X_r>  f]|[f <any_consonant_X_r>] ) ((",
	{aa, s}: ")|(<any_vowel_Naa> n)|(n <any_vowel_Naa>) )         ([<any_consonant_X_s>  f]|[f <any_consonant_X_s>] ) ((", // used to be b removed to add find s in paulo last  (try, in last-paul m, n, jh, b, ng, th, z, r, k, g, v - all ok  f - bad)
	{aa, t}: ")|(<any_vowel_Naa> n)|(n <any_vowel_Naa>) )         ([<any_consonant_X_t>  f]|[f <any_consonant_X_t>] ) ((",
	{aa, z}: ")|(<any_vowel_Naa> m)|(m <any_vowel_Naa>) )         ([<any_consonant_X_z>  f]|[f <any_consonant_X_z>] ) ((",

	{ae, k}: ")|(<any_vowel_Nae> n)|(n <any_vowel_Nae>) )         ([<any_consonant_X_k>  f]|[f <any_consonant_X_k>] ) ((",
	{ae, n}: ")|(<any_vowel_Nae> s)|(s <any_vowel_Nae>) )         ([<any_consonant_X_n>  f]|[f <any_consonant_X_n>] )        f    ((",
	{ae, s}: ")|(<any_vowel_Nae> m)|(m <any_vowel_Nae>) )         ([<any_consonant_X_s>  f]|[f <any_consonant_X_s>] ) ((", // used to be b
	{ae, t}: ")|(<any_vowel_Nae> m)|(m <any_vowel_Nae>) )         ([<any_consonant_X_t>  f]|[f <any_consonant_X_t>] ) ((",

	{ah, b}:  ")|(<any_vowel_Nah> f)|(f <any_vowel_Nah>) )         ([<any_consonant_X_b>  f]|[f <any_consonant_X_b>] ) ((",
	{ah, n}:  ")|(<any_vowel_Nah> f)|(f <any_vowel_Nah>) )         ([<any_consonant_X_n>  f]|[f <any_consonant_X_n>] ) ((",
	{ah, l}:  ")|(<any_vowel_Nah> s)|(s <any_vowel_Nah>) )         ([<any_consonant_X_l>  n]|[n <any_consonant_X_l>] ) ((",
	{ah, t}:  ")|(<any_vowel_Nah> n)|(n <any_vowel_Nah>) )         ([<any_consonant_X_t>  n]|[n <any_consonant_X_t>] ) ((",
	{ah, s}:  ")|(<any_vowel_Nah> p)|(p <any_vowel_Nah>) )         ([<any_consonant_X_s>  n]|[n <any_consonant_X_s>] ) ((",
	{ah, sh}: ")|(<any_vowel_Nah> m)|(m <any_vowel_Nah>) )         ([<any_consonant_X_sh> n]|[n <any_consonant_X_sh>]) ((",

	{ao, l}:  ")|(<any_vowel_Nao> p)|(p <any_vowel_Nao>) )         ([<any_consonant_X_l>  n]|[n <any_consonant_X_l>] ) ((",
	{ao, n}:  ")|(<any_vowel_Nao> p)|(p <any_vowel_Nao>) )         ([<any_consonant_X_n>  f]|[f <any_consonant_X_r>] ) ((",
	{ao, m}:  ")|(<any_vowel_Nao> p)|(p <any_vowel_Nao>) )         ([<any_consonant_X_n>  f]|[f <any_consonant_X_r>] ) ((",
	{ao, s}:  ")|(<any_vowel_Nao> p)|(p <any_vowel_Nao>) )         ([<any_consonant_X_s>  n]|[n <any_consonant_X_s>] ) ((",
	{ao, t}:  ")|(<any_vowel_Nao> n)|(n <any_vowel_Nao>) )         ([<any_consonant_X_t>  n]|[n <any_consonant_X_t>] ) ((",
	{ao, uh}: ")|(<any_vowel_Nao_uh> s)|(s <any_vowel_Nao_uh>) )         ((",
	{ao, z}:  ")|(<any_vowel_Nao> p)|(p <any_vowel_Nao>) )      [sil]     ([<any_consonant_X_z>  n]|[n <any_consonant_X_z>] ) ((", // not n - Alveolar like Z, not k - too weak, not b - voiced
	//{ao, z}:  ")|(<any_vowel_Nao> w)|(w <any_vowel_Nao>) )             ([<any_consonant_X_z>  n]|[n <any_consonant_X_z>] ) ((",

	{ay, ah}: ")|(<any_vowel_Nay_ah> ch)|(ch <any_vowel_Nay_ah>) )         ((",
	{ay, d}:  ")|(<any_vowel_Nay> ch)|(ch <any_vowel_Nay>) )        ([<any_consonant_X_d>  n]|[n <any_consonant_X_d>] ) ((",
	{ay, m}:  ")|(<any_vowel_Nay> f)|(f <any_vowel_Nay>) )        ([<any_consonant_X_m>  n]|[n <any_consonant_X_m>] ) ((",

	//Earlier concepts
	/*
			{aa, t}: ")| (n ah) |(n ae)|(n er)|(n ao)|(ae n)|(er n)|(ao n)  ) [<any_vowel_NaaN> n])                (((",
			{ah, dh}: "( (n t) | (n d) | (n sh) | (n zh) | th |",							// tried n, t, jh, hh, s   -- this may require something on the other side of the dh as well - other-paul2 .... didn't work:- v,
		 	{ah, t}: ")|(n ah)|(n ae)|(n er)|(n aa)|(ae n)|(er n)|(aa n)  ) [<any_vowel_NahN> n])                (((",
			{ah, t}: ")|(ow)|(ae)|(er)|(aa)  ) | (<any_vowel_NahN> n) )                ((",
			{ao, t}: ")|(aa)|(ae)|(ah)|(er)|(uh)|(ow)  ) [<any_vowel_NaoN> n])                ((",
			{ao, t}: ")|(n aa)|(n ae)|(n ah)|(n er)|(n uh)|(n ow)|(aa n)|(ae n)|(ah n)|(er n)|(uh n)|(ow n)  ) [<any_vowel_NaoN> n])                ((",
			{ao, t}: ")|(<any_vowel_NaoN> n)|(n <any_vowel_NaoN>)|(ao <any_vowel_NaoN> n)|(n <any_vowel_NaoN> ao) )               ((",
			{ay, m}: ")        ([<any_consonant_X_mN>] ((",              												// tried t, m, n, p, k, hh, l, n, ng, ch, b, d, dh, r, jh, y, s, t but none work in Climbed - Paul
	*/

	//b into something
	//================

	//{b, ih}: ")|(b b)|(b hh)|p|k|t|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )          ([<any_vowel_Nih> f]|[f <any_vowel_Nih>]) ((",
	//{b, iy}: ")|(b b)|(b hh)|p|k|t|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )          ([<any_vowel_Niy> f]|[f <any_vowel_Niy>]) ((",

	{b, l}: ") | (<any_consonant_X_b_l> jh)|(jh <any_consonant_X_b_l>) )                ( (",

	{b, ih}: ")|(b b)|(b hh)|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )         ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_b>  f]|[f <any_consonant_X_b>] )          ([<any_vowel_Nih> f]|[f <any_vowel_Nih>]) ((",
	{b, iy}: ")|(b b)|(b hh)|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )         ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_b>  f]|[f <any_consonant_X_b>] )          ([<any_vowel_Niy> f]|[f <any_vowel_Niy>]) ((",

	//{b, ih}: ")|(b b)|(b hh)|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )                  ([<any_vowel_Nih> f]|[f <any_vowel_Nih>]) ((",
	//{b, iy}: ")|(b b)|(b hh)|(<any_consonant_X_b> sh)|(sh <any_consonant_X_b>) )                 ([<any_vowel_Niy> f]|[f <any_vowel_Niy>]) ((",

	//ch into something
	//================
	{ch, ah}: ")|(t)| (<any_consonant_X_ch> n)|(n <any_consonant_X_ch>) )       ([<any_vowel_Nah> f]|[f <any_vowel_Nah>])  ( (",
	{ch, er}: ")|(t)| (<any_consonant_X_ch> n)|(n <any_consonant_X_ch>) )       ([<any_vowel_Ner> f]|[f <any_vowel_Ner>])  ( (",

	//d into something
	//================
	{d, ae}: ") | (<any_consonant_X_d> p)|(p <any_consonant_X_d>) )             ([<any_vowel_Nae> f]|[f <any_vowel_Nae>])    ( (",
	{d, ah}: ") | (<any_consonant_X_d> p)|(p <any_consonant_X_d>) )             ([<any_vowel_Nah> f]|[f <any_vowel_Nah>])    ( (",
	{d, l}:  "|(d d)|(d hh)|(d t)|(d s)|(d sil) ) | (<any_consonant_X_d> jh)|(jh <any_consonant_X_d>) )                ( (",

	{dh, aa}: ") | (<any_consonant_X_dh> p)|(p <any_consonant_X_dh>) )             ([<any_vowel_Naa> f]|[f <any_vowel_Naa>])    ( (",
	{dh, ah}: ") | (<any_consonant_X_dh> p)|(p <any_consonant_X_dh>) )             ([<any_vowel_Nah> f]|[f <any_vowel_Nah>])    ( (",
	{dh, er}: ") | (<any_consonant_X_d> jh)|(jh <any_consonant_X_d>) )             ([<any_vowel_Ner> f]|[f <any_vowel_Ner>])    ( (",
	{dh, iy}: ") | (<any_consonant_X_d> jh)|(jh <any_consonant_X_d>) )             ([<any_vowel_Niy> f]|[f <any_vowel_Niy>])    ( (",

	//e into something
	//================
	{eh, s}: ")|(<any_vowel_Neh> n)|(n <any_vowel_Neh>) )               ((",
	{eh, l}: ")|(<any_vowel_Ner> n)|(n <any_vowel_Ner>) )               ((",
	{er, l}: ")|(<any_vowel_Ner> n)|(n <any_vowel_Ner>) )               ((",
	{er, s}: ")|(<any_vowel_Ner> n)|(n <any_vowel_Ner>) )               ((",

	{ey, sh}: ")|(<any_vowel_Ney> n)|(n <any_vowel_Ney>) )               ((",
	{ey, t}:  ")|(<any_vowel_Ney> n)|(n <any_vowel_Ney>) )     [<any_vowel_Ney>] [<any_consonant_X_t>]          ((",
	{ey, z}:  ")|(<any_vowel_Ney> n)|(n <any_vowel_Ney>) )               ((",

	//f into something
	//================
	{f, l}: ")|(<any_consonant_X_f_l> n)|(n <any_consonant_X_f_l>) )                ( (",

	//g into something
	//================
	{g, aa}: "| <gg>)|k|p|t |(<any_consonant_X_g> n)|(n <any_consonant_X_g>) )            ([<any_vowel_Naa> f]|[f <any_vowel_Naa>])    ( (",
	{g, ah}: "| <gg>)|k|p|t |(<any_consonant_X_g> n)|(n <any_consonant_X_g>) )            ([<any_vowel_Nah> f]|[f <any_vowel_Nah>])    ( (",
	{g, ao}: "| <gg>)|k|p|t |(<any_consonant_X_g> n)|(n <any_consonant_X_g>) )            ([<any_vowel_Nao> f]|[f <any_vowel_Nao>])    ( (",

	{g, ih}: "| <gg>)|k|p|t |(<any_consonant_X_g> n)|(n <any_consonant_X_g>) )            ([<any_vowel_Nih> f]|[f <any_vowel_Nih>])    ( (",

	//{g, r}:	"| <gg>)|k|(<any_consonant_X_g> n)|(n <any_consonant_X_g>)  )               ( (",
	//{g, r}:	"| <gg>)|(sil k)|(k k)|(g k)|(<any_consonant> n)|(n <any_consonant>)  )               ( (",
	//{g, r}:	")|(sil k)|(k k)|(g k)|(<any_consonant> n)|(n <any_consonant>)  )               ( (",
	{g, r}: ")|k|(<any_consonant> <any_consonant_X_g_r>)  )      [<any_vowel>] [<any_consonant_X_g_r>]         ( (",
	{g, l}: "| <gg>)|k|p|t |(<any_consonant_X_g_l> n)|(n <any_consonant_X_g_l>)  )               ( (",

	//h into something
	//================
	{hh, aa}: ") | (<any_consonant_X_hh> n)|(n <any_consonant_X_hh>) )                ( (",
	{hh, ay}: ") | (<any_consonant_X_hh> n)|(n <any_consonant_X_hh>) )                ( (",
	{hh, eh}: ") | (<any_consonant_X_hh> n)|(n <any_consonant_X_hh>) )                ( (",
	{hh, t}:  ") | (<any_consonant_X_hh> n)|(n <any_consonant_X_hh>) )                ( (",

	//i into something
	//================
	{ih, b}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_b>  n]|[n <any_consonant_X_b>] ) ((",
	{ih, d}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_d>  n]|[n <any_consonant_X_d>] ) ((",
	{ih, f}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_f>  n]|[n <any_consonant_X_f>] ) ((",
	{ih, g}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_g>  n]|[n <any_consonant_X_g>] ) ((",
	{ih, k}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_k>  n]|[n <any_consonant_X_k>] ) ((",
	{ih, l}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_l>  n]|[n <any_consonant_X_l>] ) ((",
	{ih, m}: ")|(<any_vowel_Nih> s)|(s <any_vowel_Nih>) )               ((",
	{ih, p}: ")|(<any_vowel_Nih> sh)|(sh <any_vowel_Nih>) )               ((",
	{ih, r}: ")|(<any_vowel_Nih> n)|(n <any_vowel_Nih>) )               ((",
	{ih, s}: ")|(<any_vowel_Nih> n)|(n <any_vowel_Nih>) )               ([<any_consonant_X_s>  n]|[n <any_consonant_X_s>] )   ((",
	{ih, t}: ")|(<any_vowel_Nih> n)|(n <any_vowel_Nih>) )               ((",
	{ih, v}: ")|(<any_vowel_Nih> v)|(v <any_vowel_Nih>) )               ([<any_consonant_X_v>  n]|[n <any_consonant_X_v>] ) ((",

	{iy, ch}: ")|(<any_vowel_Niy> n)|(n <any_vowel_Niy>) )              ([<any_consonant_X_ch>  n]|[n <any_consonant_X_ch>] ) ((",
	{iy, g}:  ")|(<any_vowel_Niy> v)|(v <any_vowel_Niy>) )               ([<any_consonant_X_g>  n]|[n <any_consonant_X_g>] ) ((",
	{iy, k}:  ")|(<any_vowel_Niy> v)|(v <any_vowel_Niy>) )               ([<any_consonant_X_k>  n]|[n <any_consonant_X_k>] ) ((",
	{iy, m}:  ")|(<any_vowel_Niy> s)|(s <any_vowel_Niy>) )               ((",
	{iy, p}:  ")|(<any_vowel_Niy> sh)|(sh <any_vowel_Niy>) )               ((",
	{iy, s}:  ")|(<any_vowel_Niy> n)|(n <any_vowel_Niy>)|(p <any_vowel_Nih> ih)|(ih <any_vowel_Nih> p) )               ((",
	{iy, t}:  ")|(<any_vowel_Niy> n)|(n <any_vowel_Niy>) )               ((",

	//***** vowel into vowel, probably needs more work! ....... may not exist
	{iy, y}: ")|(<any_vowel_Niy> n)|(n <any_vowel_Niy>) )               ((",

	//Earlier concepts
	/*
			//{ih, m}: ")          ([<any_consonant_X_mN>] (",
			//{iy, p}: ")          ([<any_consonant_X_pN>]    ((sil p) | (hh p) | (p f) | (p p) | (p hh) |  (",
			//{ih, p}: ")          ([<any_consonant_X_pN>]    ((sil p) | (hh p) | (p f) | (p p) | (p hh) |  (",
		    //{ih, t}: " | (aa f) | (ae f) | (ah f) | (ao f) | (eh f) | (uh f) | (uw f)) ((",    // tried trialing p and n ---- f works, in that it stops the answer being 'p eh' becoming 'p ih'
		    //{iy, t}: "p))         [<any_consonant_X_tN>]  ( (n ch n) | (n k n) | (n d n) | (n p n) | (n",         //{iy, p}: ")      [<any_consonant_X_pN>]    ((sil p) | (hh p) | (sh p) | (v p) | (z p) | (p b) | (p p) | (p hh) | (f p) | (p hh) | (p hh p) | (p hh hh) | (",
	*/

	//k into something
	//================
	{k, aa}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Naa> f]|[f <any_vowel_Naa>])  ( (",
	{k, ae}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Nae> f]|[f <any_vowel_Nae>])  ( (",
	{k, ah}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Nah> f]|[f <any_vowel_Nah>])  ( (",
	{k, ao}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Nao> f]|[f <any_vowel_Nao>])  ( (",
	{k, ih}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Nih> f]|[f <any_vowel_Nih>])  ( (",
	{k, l}:  "| <kk>)|p|t| (<any_consonant_X_k_l> n)|(n <any_consonant_X_k_l>) )         ( (",
	{k, ow}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Now> n]|[n <any_vowel_Now>])  ( (", // had to remove the p/t challenge in "cot", "caught", "cut" when mixing then up
	{k, uw}: "| <kk>) | (<any_consonant_X_k> n)|(n <any_consonant_X_k>) )                ([<any_vowel_Nuw> f]|[f <any_vowel_Nuw>])  ( (",
	{k, t}:  "| <kk>) | (<any_consonant_X_k_t> n)|(n <any_consonant_X_k_t>) )                ( (",
	{k, r}:  "| <kk>) | (<any_consonant_X_k_r> n)|(n <any_consonant_X_k_r>) )                ( (",
	{k, w}:  "| <kk>) | (<any_consonant_X_k_w> n)|(n <any_consonant_X_k_w>) )                ( (",

	{l, aa}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Naa> f]|[f <any_vowel_Naa>])                   ((",
	{l, ae}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Nae> n]|[n <any_vowel_Nae>])                   ((",
	{l, ah}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Nah> f]|[f <any_vowel_Nah>]) ((",
	{l, ay}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Nay> f]|[f <any_vowel_Nay>]) ((",
	{l, d}:  ")|(<any_consonant_X_l_d> ch)|(ch <any_consonant_X_l_d>) )          ((",
	{l, ih}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Nih> f]|[f <any_vowel_Nih>]) ((",

	{l, iy}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Niy> dh]|[dh <any_vowel_Niy>]) ((",
	{l, ow}: ")|(<any_consonant_X_l> ch)|(ch <any_consonant_X_l>) )          ([<any_vowel_Now> f]|[f <any_vowel_Now>]) ((",

	{l, t}: ")|(<any_consonant_X_l_t> ch)|(ch <any_consonant_X_l_t>) )          ((",
	{l, z}: ")|(<any_consonant_X_l_z> ch)|(ch <any_consonant_X_l_z>) )          ((",

	{m, aa}: "|(m m)|(m m m)|(m hh))|(<any_consonant_X_m> s)|(s <any_consonant_X_m>) )           ([<any_vowel_Naa> f]|[f <any_vowel_Naa>]) ((",
	{m, ah}: "|(m m)|(m m m)|(m hh))|(<any_consonant_X_m> s)|(s <any_consonant_X_m>) )           ([<any_vowel_Nah> f]|[f <any_vowel_Nah>]) ((",
	{m, d}:  ")|(<any_consonant_X_m_d> s)|(s <any_consonant_X_m_d>) )            ((",
	{m, l}:  "|(m m)|(m m m)|(m hh))|(<any_consonant_X_m_l> s)|(s <any_consonant_X_m_l>) )          s  [(p (<any_vowel_noSlide>))|((<any_vowel_noSlide>) p)]         ((",
	{m, s}:  ")|(<any_consonant_X_m_s> p)|(p <any_consonant_X_m_s>) )            ((",
	{m, v}:  ")|(<any_consonant_X_m_v> s)|(s <any_consonant_X_m_v>) )            ((",
	{m, z}:  ")|(<any_consonant_X_m_z> s)|(s <any_consonant_X_m_z>) )            ((",

	{n, ah}: ")|(<any_consonant_X_n_d> s)|(s <any_consonant_X_n_d>) )           ((",
	{n, d}:  ")|(<any_consonant_X_n_d> s)|(s <any_consonant_X_n_d>) )           ((",
	{n, m}:  ")|(<any_consonant_X_m_n> s)|(s <any_consonant_X_m_n>) )           ((",
	{n, ih}: "|(n n)|(n hh)  )|(<any_consonant_X_n> s)|(s <any_consonant_X_n>) )           f  ([<any_vowel_Nah> f]|[f <any_vowel_Nah>]) ((",
	{n, t}:  ")|(<any_consonant_X_n_t> f)|(f <any_consonant_X_n_t>) )           ((",
	{n, v}:  ")|(n n)|(<any_consonant_X_n_v> s)|(s <any_consonant_X_n_v>) )           ((",

	{ow, dh}: " ( (n t) | (n d) | (n sh) | (n zh) |", // tried t, m, p-X, th, hh, d, f, r, l, t, z  ..... jh, v, b, w, s, k, g
	{ow, hh}: "))",

	{ow, n}: ") | (<any_vowel_NowN> f)| (f <any_vowel_NowN>) )         ((",
	//{ow, t}: ")| (n ao) |(n uh)|(n aw)|(ao n)|(uh n)|(aw n)  ) [<any_vowel_NowN> n])          (((",         //*********** keep for a while, coat, caught, cot, very difficult to distinguish
	{ow, t}: ") | (<any_vowel_NowN> n)| (n <any_vowel_NowN>) )       [<any_consonant_X_t>]     ((",

	{p, ah}: "| <pp>) | (<any_consonant_X_p> n)|(n <any_consonant_X_p>) )                ([<any_vowel_Nah> n]|[n <any_vowel_Nah>]) ( (",
	{p, b}:  "| <pp>) | (<any_consonant_X_p_b> n)|(n <any_consonant_X_p_b>) )                ( (",
	{p, l}:  "| <pp>) | (<any_consonant_X_p_l> n)|(n <any_consonant_X_p_l>) )         jh       ( (",
	{p, ih}: "| <pp>) | (<any_consonant_X_p> n)|(n <any_consonant_X_p>) )                ([<any_vowel_Nih> n]|[n <any_vowel_Nih>]) ( (",

	{r, aa}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )               ( (",
	{r, ae}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )               ( (",
	{r, ah}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )               ( (",
	{r, d}:  ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )               ( (",
	{r, eh}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )               ( (",
	//{r, ey}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )     [<any_vowel_Ney>] [<any_consonant_X_r>]          ( (",
	{r, ey}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )     ([<any_vowel_Ney> n]|[n <any_vowel_Ney>]) ([<any_consonant_X_r>  f]|[f <any_consonant_X_r>] )           ( (",

	{r, ih}: ") | (<any_consonant_X_r> z) | (z <any_consonant_X_r>)  )               ( (",
	{r, iy}: ") | (<any_consonant_X_r> v) | (v <any_consonant_X_r>)  )               ( (",

	{r, l}: ") | (<any_consonant_X_r> n) | (n <any_consonant_X_r>)  )  f             ( (",

	{r, t}: ") | (<any_consonant_X_r_t> n) | (n <any_consonant_X_r_t>)  )              ( (",

	//{r, ih}: "[<any_consonant_X_rN>]) | (r r [<any_consonant_X_rN>]) | (r l [<any_consonant_X_rN>]) | (r w [<any_consonant_X_rN>]) ) ))       ( ([<any_consonant_X_rN>] p (<any_vowel>) p) | ([<any_consonant_X_rN>] ih) | ",
	//{r, iy}: "[<any_consonant_X_rN>]) | (r r [<any_consonant_X_rN>]) | (r w [<any_consonant_X_rN>]) ) ))       ( ([<any_consonant_X_rN>] p (<any_vowel>) p) | ([<any_consonant_X_rN>] iy) | ",

	{s, ae}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )           ((",
	{s, ah}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )           ((",
	{s, d}:  ")|(<any_consonant_X_s_d> n)|(n <any_consonant_X_s_d>) )           ((",
	{s, ey}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )           ((",
	{s, ih}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )           ((",
	{s, iy}: ")|(<any_consonant_X_s> n)|(n <any_consonant_X_s>) )           ((",
	{s, k}:  ")|(<any_consonant_X_s_k> n)|(n <any_consonant_X_s_k>) )           ((",
	{s, n}:  ")|(<any_consonant_X_s_n> p)|(p <any_consonant_X_s_n>) )           ((",
	{s, p}:  ")|(<any_consonant_X_s_p> n)|(n <any_consonant_X_s_p>) )           ((",
	{s, t}:  ")|(<any_consonant_X_s_t> n)|(n <any_consonant_X_s_t>) )    jh ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_s_t>  f]|[f <any_consonant_X_s_t>] )      ((",

	{sh, ah}: ")|(<any_consonant_X_sh> k)|(k <any_consonant_X_sh>) )           ((",
	{sh, k}:  ")|(<any_consonant_X_sh> n)|(n <any_consonant_X_sh>) )           ((",
	{sh, n}:  ")|(<any_consonant_X_sh> n)|(n <any_consonant_X_sh>) )           ((",

	{t, ah}: ")|(<any_consonant_X_t> m)|(m <any_consonant_X_t>) )       [<any_vowel_Nah>] [<any_consonant_X_t>]        ((",
	{t, ay}: ")|(<any_consonant_X_t> m)|(m <any_consonant_X_t>) )           ((",
	{t, er}: ")|(<any_consonant_X_t> m)|(m <any_consonant_X_t>) )           ((",
	{t, ih}: ")|(<any_consonant_X_t> m)|(m <any_consonant_X_t>) )       [<any_vowel_Nih>] [<any_consonant_X_s>]    ((",
	{t, s}:  ")|(<any_consonant_X_s_t> m)|(m <any_consonant_X_s_t>) )       jh    ([<any_vowel_noSlide> n]|[n <any_vowel_noSlide>]) ([<any_consonant_X_s_t>  f]|[f <any_consonant_X_s_t>] )      ((",

	{uh, r}: ")|(<any_vowel_Nuw> m)|(m <any_vowel_Nuw>) )               ((",
	{uw, l}: ")|(<any_vowel_Nuw> m)|(m <any_vowel_Nuw>) )               ((",

	{v, eh}: ")|(<any_consonant_X_v> m)|(m <any_consonant_X_v>) )           ((",
	{v, er}: ")|(<any_consonant_X_v> m)|(m <any_consonant_X_v>) )           ((",
	{v, ih}: ")|(<any_consonant_X_v> m)|(m <any_consonant_X_v>) )           ((",
	{v, z}:  ")|(<any_consonant_X_v_z> d)|(d <any_consonant_X_v_z>) )           ((",
	{v, l}:  ")|(<any_consonant_X_v_l> d)|(d <any_consonant_X_v_l>) )           ((",

	//{w, ao}: "|v|(v hh)|(v v)|(f v)|(v f)|(v hh v)|(v v hh)|(hh v) )|(<any_consonant_X_w> m)|(m <any_consonant_X_w>) )           ((",
	{w, ao}: "| v ) |(<any_consonant_X_w> m)|(m <any_consonant_X_w>) )           ([<any_vowel_Nao> p]|[p <any_vowel_Nao>])  ( (",
	{w, er}: ")|(<any_consonant_X_w> m)|(m <any_consonant_X_w>) )           ((",
	{w, ih}: ")|(<any_consonant_X_w> m)|(m <any_consonant_X_w>) )           ((",

	{y, ih}: ")|(<any_vowel_Ny_ih> m)|(m <any_vowel_Ny_ih>) )               ((",
	{z, d}:  ")|(<any_consonant_X_z> n)|(n <any_consonant_X_z>) )           ((",
	{z, v}:  ")|(<any_consonant_X_z> d)|(d <any_consonant_X_z>) )           ((",
}

func trimInitialPlosives(results []psPhonemeResults) []psPhonemeResults {
	ret := []psPhonemeResults{}
	for _, result := range results {
		if len(result.data) == 0 {
			ret = append(ret, result)
			continue
		}
		phons := result.data
		if len(phons) == 0 {
			ret = append(ret, result)
			continue
		}
		if isPlosive(phons[0].phoneme) && phons[0].end-phons[0].start > 10 {
			newResult := result
			newResult.data = newResult.data[1:]
			ret = append(ret, newResult)
		} else {
			ret = append(ret, result)
		}
	}
	return ret
}

func couldBeBetter(result variantResult) bool {
	for _, v := range result.verdict {
		if v.goodBadEtc == missing {
			return true
		}
	}
	return false
}

func bestVerdict(v1, v2 phonVerdict) phonVerdict {
	switch v1.goodBadEtc {
	case good:
		// good is as good as it gets so return v1
		return v1
	case possible:
		if v2.goodBadEtc == good {
			return v2
		} else {
			return v1
		}
	case missing:
		//It doesn't get any worse than this so return v2. At the very worst it's
		// also missing
		return v2
	case surprise:
		if v2.goodBadEtc == missing {
			return phonVerdict{
				v1.phon,
				possible,
			}
		} else {
			return v2
		}
	default:
		log.Panic()
	}
	// Go's lack of enums means I have to return something here
	return v1
}

func phonsEqual(phons1, phons2 []phonVerdict) bool {
	if phons1 == nil || phons2 == nil {
		return false
	}
	if len(phons1) != len(phons2) {
		return false
	}
	for i := range phons1 {
		if phons1[i].phon != phons2[i].phon {
			return false
		}
	}
	return true
}

func verdictsEqual(verd1, verd2 []phonVerdict) bool {
	if verd1 == nil || verd2 == nil {
		return false
	}
	if len(verd1) != len(verd2) {
		return false
	}
	for i := range verd1 {
		if verd1[i] != verd2[i] {
			return false
		}
	}
	return true
}

// pocketsphinx can, at times, return apparently inconsistent results. For
// example k, good; ao, good; t, missing; but given an alternate phonetic
// spelling pocketsphinx will find the missing t (for example k, good; ao,
// good; r, missing; t, possible;) indicating there is a t present.
// searchForBest and following functions is an attempt to find apparently
// missing phonemes that appear fine in an alternate phonetic spelling.

func searchForBest(results []variantResult, spellings [][]phoneme) (variantResult, bool) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	for _, result := range results {
		if better, found := searchForBetter(result, spellings); found {
			return better, found
		}
	}
	return variantResult{}, false
}

func searchForBetter(result variantResult, spellings [][]phoneme) (variantResult, bool) {
	for _, spelling := range spellings {
		if better, found := trySpelling(result.verdict, spelling); found {
			newPhonemes := []phoneme{}
			for _, v := range result.verdict {
				if v.goodBadEtc != missing {
					newPhonemes = append(newPhonemes, v.phon)
				}
			}
			betterResult := variantResult{
				result.score,
				better,
				spelling,
			}
			return betterResult, true
		}
	}
	return variantResult{}, false
}

func trySpelling(verdict []phonVerdict, spelling []phoneme) ([]phonVerdict, bool) {
	if len(spelling) == 0 {
		// We've got to the end of the spelling so we check to see if there's
		// anything left in the verdict that needs to be returned
		returnVerdict := []phonVerdict{}
		for _, v := range verdict {
			if v.goodBadEtc == missing {
				continue
			}
			newVerdict := phonVerdict{
				v.phon,
				surprise,
			}
			returnVerdict = append(returnVerdict, newVerdict)
		}
		return returnVerdict, true
	}
	ph := spelling[0]
	indices := lookFor(ph, verdict)
	if len(indices) == 0 {
		return []phonVerdict{}, false
	}
	for j := 0; j < len(indices); j++ {
		betterResult, found := trySpelling(verdict[indices[j]+1:], spelling[1:])
		if !found {
			// Try starting at the next index
			continue
		}
		returnResult := []phonVerdict{}

		for k := 0; k < indices[j]; k++ {
			if verdict[k].goodBadEtc == missing {
				continue
			}
			returnResult = append(returnResult, verdict[k])
		}
		verdict_j := verdict[indices[j]]
		if verdict_j.goodBadEtc == surprise {
			verdict_j.goodBadEtc = possible
		}
		returnResult = append(returnResult, verdict_j)
		returnResult = append(returnResult, betterResult...)
		return returnResult, true
	}
	return []phonVerdict{}, false
}

// lookFor looks for instance of the phoneme, ph, in the variantResult, in,
// stopping as soon as it finds a different phoneme that is either good or
// possible. The index of any matching phoneme, which is not missing, found
// along the way is returned
func lookFor(ph phoneme, in []phonVerdict) []int {
	foundIndices := []int{}
	for i := 0; i < len(in); i++ {
		phVerdict := in[i]
		if phVerdict.goodBadEtc == good || phVerdict.goodBadEtc == possible {
			//return []int{i}    // breaks animal2_khurrum and probably others
			if phVerdict.phon != ph {
				return foundIndices
			}
		}
		if ph == phVerdict.phon && phVerdict.goodBadEtc != missing {
			foundIndices = append(foundIndices, i)
		}
	}
	return foundIndices
}

func trimDupSurprises(psV []phonVerdict) []phonVerdict {
	trimmedVerdict := []phonVerdict{}
	for i, pV := range psV {
		if i != 0 {
			if pV.goodBadEtc == surprise && pV.phon == psV[i-1].phon {
				continue
			}
		}
		trimmedVerdict = append(trimmedVerdict, pV)
	}
	return trimmedVerdict
}

func fixAudioFile(audiofile string) string {
	dir, file := filepath.Split(audiofile)
	ext := filepath.Ext(file)
	fixedfile := filepath.Join(dir, file[:len(file)-len(ext)]+"_fixed"+ext)
	_, err := exec.Command("sox", audiofile, "-r", "16000", "-c", "1", "-b", "16", fixedfile).Output()
	if err != nil {
		debug("fixAudioFile: call to sox failed. err =", err)
		log.Panic()
	}
	return fixedfile
}

func clean(outfolder string) {
	err := os.RemoveAll(outfolder)
	if err != nil {
		debug("Doh! Error on removing folder, Temp. Error =", err)
	}
}

func cleanWavFiles(originalWavFile string, word string) {
	testCaseAudio(originalWavFile, word)
	// wavDir, _ := filepath.Split(originalWavFile)
	// _ = filepath.Walk(wavDir, func(path string, info os.FileInfo, err error) error {
	// 	_, wavFile := filepath.Split(path)
	// 	if matched, _ := pathpkg.Match("*fixed*.wav", wavFile); matched {

	// 		removeFromDisk(path)
	// 	}
	// 	return nil
	// })
}

//=====================================================================
//   ___                                                   _   _
//  | _ \_  _ _ _      __ ___ _ _  __ _  _ _ _ _ _ ___ _ _| |_| |_  _
//  |   / || | ' \    / _/ _ \ ' \/ _| || | '_| '_/ -_) ' \  _| | || |
//  |_|_\\_,_|_||_|   \__\___/_||_\__|\_,_|_| |_| \___|_||_\__|_|\_, |
//                                                               |__/
//=====================================================================

func Pronounce(outfolder string,
	audiofile string,
	word string,
	dictfile, phdictfile string,
	featparams string,
	hmm string,
) ([]LettersVerdict, error) {

	defer clean(outfolder)
	defer cleanWavFiles(audiofile, word)

	verdict, err := pronounce(outfolder, audiofile, word, dictfile, phdictfile, featparams, hmm)

	return verdict, err
}

func pronounce(outfolder string, audiofile string, word string, dictfile, phdictfile string, featparams string, hmm string) ([]LettersVerdict, error) {
	dict := dictionary.Create(dictfile)
	variantWords, err := dict.LookupWord(word)
	if err != nil {
		return []LettersVerdict{}, err
	}
	// trimAudio := trimAudio(audiofile)

	variantPhons := [][]phoneme{}
	for _, variantWord := range variantWords {
		phons, err := phonemesForWord(dict, variantWord)
		if err != nil {
			return []LettersVerdict{}, err
		}
		variantPhons = append(variantPhons, phons)
	}

	if len(variantPhons) > 0 {
		audiofile = fixAudioFile(audiofile)
		audiofile = trimAudio(audiofile, variantPhons[0])
	}

	sch := scanScheduler.New(outfolder, audiofile, phdictfile)

	d1 := make(chan bool)

	var results = []variantResult{}
	var variantResults = []variantResult{}
	runTestVariantScans(sch, outfolder, audiofile, phdictfile, dict, featparams, hmm, word, variantPhons, defaultSuite, func(arg1 []variantResult) {
		results = arg1
		go func() {
			d1 <- true
		}()
	})

	// <- d

	// d2 := make(chan bool)

	// var diphthongResults []variantResult
	// runTestDiphthongScans(sch, outfolder, audiofile, phdictfile, dict, featparams, word, variantPhons, defaultSuite, func(arg1 []variantResult) {
	// 	diphthongResults = arg1
	// 	go func() {
	// 		d2 <- true
	// 	}()
	// })

	<-d1
	results = append(results, variantResults...)

	// <-d2
	// results = append(results, diphthongResults...)

	if len(results) == 0 {
		return []LettersVerdict{}, psAborted
	}
	bestResult := results[0]
	for _, result := range results {
		if result.score > bestResult.score {
			bestResult = result
		}
	}

	// But is this really the best result? pocketsphinx can be a bit inconsistent
	// at times for instance claiming that the R is missing in HH AA R D L IY
	// but when presented with the (phonetic) spelling HH AA D L IY then says the
	// D is missing! Clearly HH AA D L IY should be the result to go for...
	if couldBeBetter(bestResult) {
		debug("bestResult =", bestResult)
		// if result, found := searchForBetter(results); found {
		if result, found := searchForBest(results, variantPhons); found {
			bestResult = result
			debug("Updated bestResult =", bestResult)
		}
	}

	// Print bestResult to the console for testPronounce to pick up
	testResult(bestResult.verdict)

	ret, err := publish(word, bestResult.phons, bestResult.verdict)
	if err != nil {
		return []LettersVerdict{}, err
	}
	if allMissing(ret) {
		return ret, psAborted
	}
	return ret, nil
}

func testResult(bestResult []phonVerdict) {
	str := "testPronounce"
	for _, result := range bestResult {
		str += " " + string(result.phon) + " " + verdicts[result.goodBadEtc]
	}
	debug(str)
}

func runTestVariantScans(sch scanScheduler.Scheduler, outfolder, audiofile, phdictfile string, dict dictionary.Dictionary, featparams string, hmm string, word string, variants [][]phoneme, suiteToRun psSuite, f func([]variantResult)) {
	c := make(chan variantResult)
	var wg sync.WaitGroup

	for _, variant := range variants {

		//if !diphthongsInWord(variant) || len(variant) > 3 {

		// Only run a non-diphthong scan if the word has no diphthong vowels or
		// is longer than 3 phonemes
		wg.Add(1)
		go runVariantScan(c, &wg, sch, outfolder, audiofile, phdictfile, dict, featparams, hmm, word, variant, suiteToRun)

		//}
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	results := []variantResult{}
	for d := range c {
		results = append(results, d)
	}
	f(results)
}

func updateVerdict(ruleVerdict []phonVerdict, timeAlignPhons []timeAlignedPhoneme) []phonVerdict {
	prev := func(i int, phs []phonVerdict) (phoneme, bool) {
		for j := i - 1; j >= 0; j-- {
			if phs[j].goodBadEtc != surprise {
				return phs[j].phon, true
			}
		}
		return aa, false
	}
	next := func(i int, phs []phonVerdict) (phoneme, bool) {
		for j := i + 1; j < len(phs); j++ {
			if phs[j].goodBadEtc != surprise {
				return phs[j].phon, true
			}
		}
		return aa, false
	}
	// Start off with the existing verdict
	newVerdict := ruleVerdict
	// Now check it against the time-aligned phonemes
	if len(timeAlignPhons) == 0 {
		// Nothing to update so return now
		return newVerdict
	}
	iNew := 0
	iTime := 0
	for i := 0; i < len(ruleVerdict); i++ {
		if iTime >= len(timeAlignPhons) {
			// We've run out of time-aligned phonemes so whatever's in newVerdict is
			// as good as it can get
			return newVerdict
		}
		if iNew >= len(newVerdict) {
			// We've got to the end of newVerdict so can't make any more changes so
			// return newVerdict now
			return newVerdict
		}
		switch ruleVerdict[i].goodBadEtc {
		case good, possible:
			if ruleVerdict[i].phon != timeAlignPhons[iTime].phoneme {
				// There's disagreement between the rule-aligned and time-aligned
				// verdicts
				// Search forward to see if we can
				switch ruleVerdict[i].phon {
				case b:
					_, ok := prev(i, ruleVerdict)
					if ok {
						// b is not the opening phoneme so change the verdict to missing
						newVerdict[iNew].goodBadEtc = missing
					}
					iNew++

				case t:
					prevPh, okPrev := prev(i, ruleVerdict)
					nextPh, okNext := next(i, ruleVerdict)
					if okPrev && okNext && prevPh == s && nextPh == s {
					} else {
						// The t is not sandwiched between two s's so change the verdict
						// for t to missing
						newVerdict[iNew].goodBadEtc = missing
					}
					if okPrev && okNext && prevPh == s && nextPh == r {
					} else {
						// The t is not sandwiched between an s an an r so change the
						// verdict for t to missing
						newVerdict[iNew].goodBadEtc = missing
					}
					iNew++
				default:
					// In general the time-aligned result wins out so change the verdict
					// to missing
					// Actually we need to search forward to see if we can find a
					// time-aligned phoneme that matches the current rule-aligned
					// phoneme. If we find one we could add each of the intervening
					// phonemes as a surprise and leave the rule-aligned phoneme in
					// place. If we don't find a matching time-aligned phoneme then we
					// should
					newVerdict[iNew].goodBadEtc = missing
					iNew++
				}
			} else {
				iNew++
				iTime++
			}
		case missing:
			if ruleVerdict[i].phon == timeAlignPhons[iTime].phoneme {
				newVerdict[iNew].goodBadEtc = possible
				iNew++
				iTime++
			} else {
				iNew++
			}
		case surprise:
			if iTime >= len(timeAlignPhons) || ruleVerdict[i].phon != timeAlignPhons[iTime].phoneme {
				// Not sure the surprise is showing up in timeAlignPhons so perhaps
				// remove this from newVerdict
				newVerdict = append(newVerdict[:iNew], newVerdict[iNew+1:]...)
			} else {
				iNew++
				iTime++
			}
		}
	}
	return newVerdict
}

type resultWithConfig struct {
	psPhonemeResults
	newPsConfig
}

func runVariantScan(c chan<- variantResult, wg *sync.WaitGroup, sch scanScheduler.Scheduler, outfolder, audiofile, phdictfile string, dict dictionary.Dictionary, featparams string, hmm string, word string, variant []phoneme, suiteToRun psSuite) {
	defer wg.Done()

	c1 := make(chan resultWithConfig)
	var wg1 sync.WaitGroup

	frates := []int{
		65, 100, 256,
	}
	for _, scan := range suiteToRun {
		debug("Running scan for frate =", scan["-frate"])
		wg1.Add(1)
		builderConfig := new_jsgfStandard()
		suiteOfOne := psSuite{
			scan,
		}
		/*
		   frate, err := strconv.Atoi(scan["-frate"])
		   if err != nil {
		     debug("Failed to get frate for batch scan")
		   }
		*/
		config := TestConfig(outfolder, audiofile, phdictfile, dict, featparams, hmm, word, variant, frates, suiteOfOne, targetRuleForWord, &builderConfig)
		go runVariantScanWithConfig(c1, &wg1, sch, config, word)
	}

	go func() {
		wg1.Wait()
		close(c1)
	}()

	results := []resultWithConfig{}
	for d := range c1 {
		results = append(results, d)
	}

	sortedResults := results
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].psPhonemeResults.frate < sortedResults[j].psPhonemeResults.frate
	})
	str := "\nresults =\n"
	for _, result := range sortedResults {
		str += fmt.Sprintln(result.psPhonemeResults)
	}
	debug(str)

	combiner := newCombiner(variant)
	str = "normalised results =\n"
	for _, result := range sortedResults {
		normResult := timeNormalise(result.psPhonemeResults)
		str += fmt.Sprintln(normResult)
		combiner.addResult(normResult, result.newPsConfig.tRule)
	}
	debug(str)

	combiner.parse()
	ruleVerdicts := combiner.ruleAlign()
	debug("ruleAligned =", ruleVerdicts)
	for i := range combiner.results {
		combiner.timeAligner.AddResult(combiner.results[i].psPhonemeResults)
	}
	debug("")

	str = "results (after rule-alignment) =\n"
	for _, result := range combiner.results {
		str += fmt.Sprintln(result.psPhonemeResults)
	}
	debug(str)

	combiner.timeAlignedPhonemes = combiner.timeAligner.timeAlign()
	timeVerdicts := combiner.timeAlignedPhonemes
	debug("timeAligned =", timeVerdicts)
	debug("")
	combinedVerdict := combiner.integrate()
	debug("combinedVerdict    =", combinedVerdict)
	ruleAlignedVerdict := mapLinkedVerdicts(ruleVerdicts)
	debug("ruleAlignedVerdict =", ruleAlignedVerdict)
	debug("")

	verdict := combinedVerdict
	// verdict := ruleAlignedVerdict
	verdict = trimDupSurprises(verdict)

	result := variantResult{
		score(verdict),
		verdict,
		variant,
	}

	c <- result
}

////////////////////////////////////////////////////////////
// Remove the extra P
// If changing the put the function below back in

/*
func runVariantScanWithConfig(c chan <-resultWithConfig, wg *sync.WaitGroup, config newPsConfig) {
  defer wg.Done()

  scheduler := NewScanScheduler()
  scheduler.RunNewScan(config, func(arg1 []psPhonemeResults) {
    result := resultWithConfig{}
    // There should only be one result
    if len(arg1) == 1 {
      psRes := psPhonemeResults{
        arg1[0].frate,
        arg1[0].data[1:],
      }
      result = resultWithConfig{
        psRes,
        config,
      }
    }
    c <- result
  })
}
*/
// type Utt struct {
// 	Text       string
// 	Start, End int32
// }

func toPhonemeData(resp []xyz_plus.Utt) []psPhonemeDatum {
	phonemeData := []psPhonemeDatum{}

	for _, utt := range resp {
		phonemeData = append(phonemeData, psPhonemeDatum{
			phoneme(utt.Text),
			int(utt.Start),
			int(utt.End),
		})
	}
	return phonemeData
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func doRunScan(s scanScheduler.Scheduler, config newPsConfig, word string, f func([]psPhonemeResults)) {
	// Need to set up arguments for acall to s.DoScan(scan PsScan)
	params := []scanScheduler.PsParam{}
	// Really need to sort this out. it's a mess... Right now there should only
	// ever be one set of settings so let's throw an error if that's not so
	if len(config.settings) != 1 {
		log.Panic()
	}
	for flag, value := range config.settings[0] {
		param := scanScheduler.PsParam{
			Flag:  string(flag),
			Value: value,
		}
		params = append(params, param)
	}
	// Set a logfile name
	// logfile := path.Join(config.tempDir, pseudo_uuid() + ".log")
	logfile := path.Join(config.tempDir, uuid.New().String()+".log")
	param := scanScheduler.PsParam{
		Flag:  "-logfn",
		Value: logfile,
	}
	params = append(params, param)
	//This is just so we can debug scanScheduler's call to pocketsphinx
	//so we can right a proper test case name (sometimes same audio is processed with a different target word
	//to test for mispronunciations, like in minimal pairs)
	paramword := scanScheduler.PsParam{
		Flag:  "-word",
		Value: word,
	}
	params = append(params, paramword)

	//Locate jsgf and audio files paths
	var jsgf_file string
	var audio_file string
	for _, param := range params {
		if param.Flag == "-infile" {
			audio_file = param.Value
		}
		if param.Flag == "-jsgf" {
			jsgf_file = param.Value
		}
	}
	var jsgf_buffer []byte
	var audio_buffer []byte
	var err error
	jsgf_buffer, err = os.ReadFile(jsgf_file)
	check(err)
	audio_buffer, err = os.ReadFile(audio_file)
	check(err)
	//  ___                            _           _          _      _
	// | _ \__ _ _ _ __ _ _ __  ___   | |_ ___    | |__  __ _| |_ __| |_      ___ __ __ _ _ _
	// |  _/ _` | '_/ _` | '  \(_-<   |  _/ _ \   | '_ \/ _` |  _/ _| ' \    (_-</ _/ _` | ' \
	// |_| \__,_|_| \__,_|_|_|_/__/    \__\___/   |_.__/\__,_|\__\__|_||_|   /__/\__\__,_|_||_|
	// Parameters for use with batch scan

	/*
	   The following notes are derived from running command line experiments. They detail the parameter which affect the cmninit values output by the batch scan

	   Makes no difference
	   ===================
	   - maxhmmpf
	   - maxwpf
	   - beam
	   - fwdflat
	   - bestpath

	   Has an impact to CMNinit values
	   ===============================
	   - alpha
	   - frate
	   - dither
	   - doublebw
	   - nfft
	   - wlen

	   Not understood why but -nfft and -wlen matching (between batch and continuous seemed to have a negative effect??)
	*/

	//context := []string{"-frate", "-lw", "-nfft", "-wlen", "-alpha", "-dither", "-doublebw", "-maxhmmpf", "-maxwpf", "-beam", "-wbeam", "-pbeam", "-fwdflat", "-bestpath", "-wip", "-pip"}
	//context := []string{"-frate",  "-remove_noise", "-remove_silence", "-vad_postspeech", "-vad_prespeech", "-vad_startspeech", "-vad_threshold",  "-topn", "-pl_window", "-lpbeam", "-lponlybeam"}
	//context := []string{"-frate", "-nfft", "-wlen", "-alpha", "-dither", "-doublebw", "-remove_silence", "-vad_postspeech", "-vad_prespeech", "-vad_startspeech", "-vad_threshold"}

	context := []string{"-hmm", "-frate", "-lw", "-nfft", "-wlen", "-alpha", "-dither", "-doublebw", "-maxhmmpf", "-maxwpf", "-beam", "-wbeam", "-pbeam", "-fwdflat", "-bestpath", "-wip", "-pip", "-remove_noise", "-remove_silence", "-vad_postspeech", "-vad_prespeech", "-vad_startspeech", "-vad_threshold", "-topn", "-pl_window", "-lpbeam", "-lponlybeam"}
	//context := []string{"-frate"}
	//ch := make(chan error)
	ch := make(chan []xyz_plus.Utt, 1)

	// Now create a scan, send it and wait on ch for a reply
	psScan := scanScheduler.PsScan{
		Settings:     params,
		ContextFlags: context,
		RespondTo:    ch,
		Jsgf_buffer:  jsgf_buffer,
		Audio_buffer: audio_buffer,
	}
	frate, _ := strconv.Atoi((config.settings[0]["-frate"]))

	s.DoScan(psScan)
	//err := <-ch // All that's returned is an error. It's the log file I need
	response := <-ch
	// if err != nil {
	// 	debug("DoScan returned err =", err)
	// }

	utts := toPhonemeData(response)
	results := psPhonemeResults{
		frate,
		//parsePsData(logfile),
		utts,
	}
	testCaseIt(params, results.data, word)
	f([]psPhonemeResults{results})
}

// Not an ideal solution but implemented to get a solution together in as
// short a time as possible and with minimal change to the existing code
func runVariantScanWithConfig(c chan<- resultWithConfig, wg *sync.WaitGroup, sch scanScheduler.Scheduler, config newPsConfig, word string) {
	defer wg.Done()

	// scheduler := NewScanScheduler()
	// scheduler.RunNewScan(config, func(arg1 []psPhonemeResults) {
	doRunScan(sch, config, word, func(arg1 []psPhonemeResults) {
		result := resultWithConfig{}
		// There should only be one result
		if len(arg1) == 1 {
			/*
			   psRes := psPhonemeResults{
			     arg1[0].frate,
			     arg1[0].data[1:],
			   }
			   result = resultWithConfig{
			     psRes,
			     config,
			   }
			*/
			result = resultWithConfig{
				arg1[0],
				config,
			}
		}
		c <- result
	})
}

func runTestDiphthongScans(sch scanScheduler.Scheduler, outfolder, audiofile, phdictfile string, dict dictionary.Dictionary, featparams string, hmm string, word string, variants [][]phoneme, suiteToRun psSuite, f func([]variantResult)) {
	c := make(chan variantResult)
	var wg sync.WaitGroup

	for _, variant := range variants {
		if diphthongsInWord(variant) {
			wg.Add(1)
			go runDiphthongScan(c, &wg, sch, outfolder, audiofile, phdictfile, dict, featparams, hmm, word, variant, suiteToRun)
		}
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	results := []variantResult{}
	for d := range c {
		results = append(results, d)
	}
	f(results)
}

func runDiphthongScan(c chan<- variantResult, wg *sync.WaitGroup, sch scanScheduler.Scheduler, outfolder, audiofile, phdictfile string, dict dictionary.Dictionary, featparams string, hmm string, word string, variant []phoneme, suiteToRun psSuite) {
	defer wg.Done()

	c1 := make(chan resultWithConfig)
	var wg1 sync.WaitGroup

	frates := []int{
		65, 100, 256,
	}
	for _, scan := range suiteToRun {
		wg1.Add(1)
		builderConfig := new_jsgfDiphthong()
		suiteOfOne := psSuite{
			scan,
		}
		/*
		   frate, err := strconv.Atoi(scan["-frate"])
		   if err != nil {
		     debug("Failed to get frate for batch scan")
		   }
		*/
		config := TestConfig(outfolder, audiofile, phdictfile, dict, featparams, hmm, word, variant, frates, suiteOfOne, targetRuleForWord, &builderConfig)
		go runVariantScanWithConfig(c1, &wg1, sch, config, word)
	}

	go func() {
		wg1.Wait()
		close(c1)
	}()

	results := []resultWithConfig{}
	for d := range c1 {
		results = append(results, d)
	}

	// Print out some debug stuff
	sortedResults := results
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].psPhonemeResults.frate < sortedResults[j].psPhonemeResults.frate
	})
	str := "\ndiphthong results =\n"
	for _, result := range sortedResults {
		str += fmt.Sprintln(result.psPhonemeResults)
	}
	debug(str)

	combiner := newCombiner(variant)
	str = "normalised results =\n"
	for _, result := range sortedResults {
		normResult := timeNormalise(result.psPhonemeResults)
		str += fmt.Sprintln(normResult)
		combiner.addResult(normResult, result.newPsConfig.tRule)
	}
	debug(str)

	combiner.parse()
	ruleVerdicts := combiner.ruleAlign()
	debug("ruleAligned =", ruleVerdicts)
	for i := range combiner.results {
		combiner.timeAligner.AddResult(combiner.results[i].psPhonemeResults)
	}
	debug("")

	str = "results (after rule-alignment) =\n"
	for _, result := range combiner.results {
		str += fmt.Sprintln(result.psPhonemeResults)
	}
	debug(str)

	combiner.timeAlignedPhonemes = combiner.timeAligner.timeAlign()
	timeVerdicts := combiner.timeAlignedPhonemes
	debug("timeAligned =", timeVerdicts)
	debug("")
	combinedVerdict := combiner.integrate()
	debug("combinedVerdict =", combinedVerdict)
	ruleAlignedVerdict := mapLinkedVerdicts(ruleVerdicts)
	debug("ruleAlignedVerdict =", ruleAlignedVerdict)
	verdict := combinedVerdict

	// verdict := ruleAlignedVerdict
	verdict = verdictWithDiphthongs(variant, verdict)
	verdict = trimDupSurprises(verdict)

	result := variantResult{
		score(verdict),
		verdict,
		variant,
	}
	c <- result
}

func minVerdict(v1, v2 verdict) verdict {
	return verdict(min(int(v1), int(v2)))
}

func maxVerdict(v1, v2 verdict) verdict {
	return verdict(max(int(v1), int(v2)))
}

func verdictWithDiphthongs(phonemes []phoneme, verdict []phonVerdict) []phonVerdict {
	searchForPhoneme := func(ph phoneme, v1 []phonVerdict) (int, bool) {
		for i := range v1 {
			if ph == v1[i].phon && v1[i].goodBadEtc != surprise {
				return i, true
			}
		}
		return 0, false
	}
	removeVerdict := func(v1 []phonVerdict, i int) []phonVerdict {
		if i < len(v1) {
			return append(v1[:i], v1[i+1:]...)
		}
		return v1
	}
	insertVerdict := func(v1 []phonVerdict, i int, v phonVerdict) []phonVerdict {
		if i < len(v1) {
			return append(v1[:i], append([]phonVerdict{v}, v1[i:]...)...)
		}
		return v1
	}
	diphthongVerdict := []phonVerdict{}
	v := 0
	i := 0
	for v < len(verdict) {
		// A surprise is a surprise - diphthongs or no diphthongs
		if verdict[v].goodBadEtc == surprise {
			diphthongVerdict = append(diphthongVerdict, verdict[v])
			v++
			continue
		}
		if verdict[v].phon == phonemes[i] {
			// Easy peasy, the phonemes match, so just add the verdict unchanged
			diphthongVerdict = append(diphthongVerdict, verdict[v])
			v++
			i++
			continue
		}
		// If we get here the phonemes don't match and we're not looking at a
		// surprise so this had better be part of a diphthong or something's gone
		// wrong
		if diphthong, ok := diphthongs[phonemes[i]]; ok {
			// The parts of the diphthong might have been split up by a surprise
			// so search for the second part
			if diphthong[0] != verdict[v].phon {
				debug("What the hell's going on? Expecting dipthong", diphthong, "but verdict phone is", verdict[v].phon)
				v++
			} else {
				if j, ok := searchForPhoneme(diphthong[1], verdict[v:]); ok {
					// We're searching from the first part phoneme of the diphthong so
					// if the two parts are adjacent j == 1
					if j == 1 {
						// Good, the phonemes are adjacent in the verdict
						newVerdict := phonVerdict{
							phonemes[i],
							maxVerdict(verdict[v].goodBadEtc, verdict[v+j].goodBadEtc),
						}
						diphthongVerdict = append(diphthongVerdict, newVerdict)
					} else {
						// The diphthong was split by a surprise (possibly many surprises)
						// so record the diphthong as missing
						newVerdict := phonVerdict{
							phonemes[i],
							missing,
						}
						diphthongVerdict = append(diphthongVerdict, newVerdict)
						// Bring the two parts of the diphthong together
						verdictToInsert := verdict[j+1]
						verdict = removeVerdict(verdict, v+j)
						verdict = insertVerdict(verdict, v+1, verdictToInsert)
					}
					// We've effectively processed both phonemes in the diphthong so set
					// v to be the following phoneme
					v += 2
					i++
					continue
				}
				debug("Failed to find the second part of dipthong", diphthong)
				v++
			}
		} else {
			// It's not a diphthong! What to do?
			debug("Hmm... Thought I'd found a diphthong but it's not. Phoneme is", phonemes[i])
			// Better increment v just so we don't get stuck...
			v++
			continue
		}
	}
	if v != len(verdict) || i != len(phonemes) {
		debug("v, len(verdict), i, len(phonemes) =", v, len(verdict), i, len(phonemes))
	}
	return diphthongVerdict
}

//=====================================================================
//  ___                           _      _ _
// | _ \_  _ _ _      ___ ___ _ _(_)__ _| | |_  _
// |   / || | ' \    (_-</ -_) '_| / _` | | | || |
// |_|_\\_,_|_||_|   /__/\___|_| |_\__,_|_|_|\_, |
//                                           |__/
//=====================================================================

/*
func testPronounce(audiofile string, word string, dictfile, phdictfile string) ([]LettersVerdict, error) {
  dict := dictionary.Create(dictfile)
  variantWords, err := dict.LookupWord(word)
  if err != nil {
    return []LettersVerdict{}, err
  }
  trimAudio := trimAudio(audiofile)

  // d := make(chan bool)

  variantPhons := [][]phoneme{}
  for _, variantWord := range variantWords {
    phons, err := phonemesForWord(dict, variantWord)
    if err != nil {
      return []LettersVerdict{}, err
    }
    variantPhons = append(variantPhons, phons)
  }
  var results = []variantResult{}
  runSerialTestVariantScans(trimAudio, phdictfile, dict, word, variantPhons, defaultSuite, func(arg1 []variantResult) {
    results = arg1
  })

  var diphthongResults []variantResult
  runSerialTestDiphthongScans(trimAudio, phdictfile, dict, word, variantPhons, defaultSuite, func(arg1 []variantResult) {
    diphthongResults = arg1
  })

  results = append(results, diphthongResults...)

  if len(results) == 0 {
    return []LettersVerdict{}, psAborted
  }
  bestResult := results[0]
  for _, result := range results {
    if result.score > bestResult.score {
      bestResult = result
    }
  }
  // Print bestResult to the console for testPronounce to pick up
  testResult(bestResult.verdict)
  ret, err := publish(word, bestResult.phons, bestResult.verdict)
  if err != nil {
    return []LettersVerdict{}, err
  }
  if allMissing(ret) {
    return ret, psAborted
  }
  return ret, nil
}

func testResult(bestResult []phonemeVerdict) {
  str := "testPronounce"
  for _, result := range bestResult {
    str += " " + string(result.psPhonemeDatum.phoneme) + " " + verdicts[result.goodBadEtc]
  }
  debug(str)
}

func runSerialTestVariantScans(audiofile, phdictfile string, dict dictionary.Dictionary, word string, variants [][]phoneme, suiteToRun psSuite, f func([]variantResult)) {
  results := []variantResult{}
  for _, variant := range variants {
    results = append(results, runSerialTestVariantScan(audiofile, phdictfile, dict, word, variant, suiteToRun))
  }
  f(results)
}

func runSerialTestVariantScan(audiofile, phdictfile string, dict dictionary.Dictionary, word string, variant []phoneme, suiteToRun psSuite) variantResult {
  scheduler := NewScanScheduler()
  frates := []int{
    99, 100, 101,
  }
  builderConfig := new_jsgfStandard()
  config := TestConfig(audiofile, phdictfile, dict, word, variant, template, frates, suiteToRun, targetRuleForWord, &builderConfig)
  return scheduler.RunSerialNewScan(config, func(arg1 []psPhonemeResults) variantResult {
    // trimmedResults := trimResults(arg1, trim)
    // Put them in frate order for printing to the console
    sort.Slice(arg1, func (i,j int) bool {
      return arg1[i].frate < arg1[j].frate
    })
    debug()
    debug("trimmedResults =")
    for _, result := range arg1 {
      debug(result)
    }
    debug()
    verdict := process(word, variant, arg1)
    result := variantResult{
      score(verdict),
      verdict,
      variant,
    }
    return result
  })
}

func runSerialTestDiphthongScans(audiofile, phdictfile string, dict dictionary.Dictionary, word string, variants [][]phoneme, suiteToRun psSuite, f func([]variantResult)) {
  results := []variantResult{}
  for _, variant := range variants {
    results = append(results, runSerialTestDiphthongScan(audiofile, phdictfile, dict, word, variant, template, suiteToRun))
  }
  f(results)
}

func runSerialTestDiphthongScan(audiofile, phdictfile string, dict dictionary.Dictionary, word string, variant []phoneme, template psPhonemeSettings, suiteToRun psSuite) variantResult {
  scheduler := NewScanScheduler()
  frates := []int{
    99, 100, 101,
  }
  diphthongBuilderConfig := new_jsgfDiphthong()
  config := TestConfig(audiofile, phdictfile, dict, word, variant, template, frates, suiteToRun, targetRuleWithDiphthongs, &diphthongBuilderConfig)
  return scheduler.RunSerialNewScan(config, func(arg1 []psPhonemeResults) variantResult {
    debug("diphthong results = ", arg1)
    diphthongisedResults := diphthongise(arg1)
    debug("diphthongised results = ", diphthongisedResults)
    // trimmedResults := trimResults(diphthongisedResults, trim)
    // Put them in frate order for printing to the console
    sort.Slice(diphthongisedResults, func (i,j int) bool {
      return diphthongisedResults[i].frate < diphthongisedResults[j].frate
    })
    debug()
    debug("trimmedResults =")
    for _, result := range diphthongisedResults {
      debug(result)
    }
    debug()
    // Amalgamate vowels into diphthongs

    verdict := process(word, variant, diphthongisedResults)
    result := variantResult{
      score(verdict),
      verdict,
      variant,
    }
    return result
  })
}
*/

type candidatePhoneme struct {
	psPhonemeDatum
	confidence int
}

func (p candidatePhoneme) incrementEnd() candidatePhoneme {
	p.end += 1
	return p
}

/*
func merge(results []psPhonemeResults) []psPhonemeDatum {
  // Make structure to hold merged results
  max := 0
  for _, result := range results {
    end := result.data[len(result.data) - 1].end
    if end > max {
      max = end
    }
  }
  mergedResults := make([]map[phoneme]int, max + 1)
  for i := 0; i < max + 1; i++ {
    mergedResults[i] = make(map[phoneme]int)
  }
  // Now merge the results
  for _, result := range results {
    for _, datum := range result.data {
      start := datum.start
      end := datum.end
      for j := start; j <= end; j++ {
        if v, ok := mergedResults[j][datum.phoneme]; ok {
          mergedResults[j][datum.phoneme] = v + 1
        } else {
          mergedResults[j][datum.phoneme] = 1
        }
      }
    }
  }
  // Now filter the merged results so we only have phonemes of interest
  // - that is either a single phoneme with a confidence of 3 or two
  // phonemes each with a count of 2
  for i, mergedResult := range mergedResults {
    max := 0
    for _, count := range mergedResult {
      if count > max {
        max = count
      }
    }
    for ph, count := range mergedResult {
      if count < max {
        delete(mergedResults[i], ph)
      }
    }
  }
  ret := []psPhonemeDatum{}
  scratchPad := map[phoneme]psPhonemeDatum{}

  for i, mergedResult := range mergedResults {
    // Tidy up scratchpad
    for ph, datum := range scratchPad {
      confidence, ok := mergedResult[ph]
      if !ok || confidence < 2 {
        // Add the phoneme datum in scratchPad to ret and delete the entry
        // from scratchPad
        ret = append(ret, datum)
        delete(scratchPad, ph)
      }
    }
    for ph, count := range mergedResult {
      if count < 2 {
        continue
      }
      if _, ok := scratchPad[ph]; !ok {
        d := psPhonemeDatum{
          ph,
          i,
          i,
        }
        scratchPad[ph] = d
      } else {
        scratchPad[ph] = scratchPad[ph].incrementEnd()
      }
    }
  }
  // Finally copy anything left in the scratchPad into ret
  for _, datum := range scratchPad {
    ret = append(ret, datum)
  }
  debug("ret =", ret)
  return ret
}
*/

func merge(results []psPhonemeResults) []candidatePhoneme {
	// Make structure to hold merged results
	max := 0
	for _, result := range results {
		end := result.data[len(result.data)-1].end
		if end > max {
			max = end
		}
	}
	mergedResults := make([]map[phoneme]int, max+1)
	for i := 0; i < max+1; i++ {
		mergedResults[i] = make(map[phoneme]int)
	}
	// Now merge the results
	for _, result := range results {
		for _, datum := range result.data {
			start := datum.start
			end := datum.end
			for j := start; j <= end; j++ {
				if v, ok := mergedResults[j][datum.phoneme]; ok {
					mergedResults[j][datum.phoneme] = v + 1
				} else {
					mergedResults[j][datum.phoneme] = 1
				}
			}
		}
	}
	// Now filter the merged results so we only have phonemes of interest
	// - that is either a single phoneme with a confidence of 3 or two
	// phonemes each with a count of 2
	for i, mergedResult := range mergedResults {
		max := 0
		for _, count := range mergedResult {
			if count > max {
				max = count
			}
		}
		for ph, count := range mergedResult {
			if count < max {
				delete(mergedResults[i], ph)
			}
		}
	}
	ret := []candidatePhoneme{}
	scratchPad := map[phoneme]candidatePhoneme{}

	for i, mergedResult := range mergedResults {
		// Tidy up scratchpad
		for ph, candidate := range scratchPad {
			confidence, ok := mergedResult[ph]
			if !ok || confidence < 2 {
				// Add the phoneme datum in scratchPad to ret and delete the entry
				// from scratchPad
				ret = append(ret, candidate)
				delete(scratchPad, ph)
			}
		}
		for ph, count := range mergedResult {
			if count < 2 {
				continue
			}
			if _, ok := scratchPad[ph]; !ok {
				candidate := candidatePhoneme{
					psPhonemeDatum{
						ph,
						i,
						i,
					},
					count,
				}
				scratchPad[ph] = candidate
			} else {
				candidate := scratchPad[ph]
				candidate.end = i
				if count > candidate.confidence {
					candidate.confidence = count
				}
				scratchPad[ph] = candidate
			}
		}
	}
	// Finally copy anything left in the scratchPad into ret
	for _, datum := range scratchPad {
		ret = append(ret, datum)
	}
	return ret
}

func (p candidatePhoneme) powerset() []candidatePhoneme {
	ret := []candidatePhoneme{
		p,
	}
	if diff := p.end - p.start; diff > 1 {
		for i := p.start; i <= p.end-(diff-1); i++ {
			q := candidatePhoneme{
				psPhonemeDatum{
					p.phoneme,
					i,
					i + diff - 1,
				},
				p.confidence,
			}
			ret = append(ret, q.powerset()...)
		}
	}
	return ret
}

type candidateData []candidatePhoneme

func (ps candidateData) powerset() []candidatePhoneme {
	data := []candidatePhoneme{}
	for _, p := range ps {
		data = append(data, p.powerset()...)
	}
	// Remove duplicates
	keys := make(map[candidatePhoneme]bool)
	ret := []candidatePhoneme{}
	for _, datum := range data {
		if _, value := keys[datum]; !value {
			keys[datum] = true
			ret = append(ret, datum)
		}
	}
	return ret
}

func (c candidateData) resolve(l link, s linkSet) link {
	debug("resolve->: l =", l, "s =", s)
	if len(s.links) == 0 {
		// There's nothing to resolve!
		//
		debug("resolve<-:", l)
		return l
	}
	// A local function to determine whether this is a to conflict.
	//
	to := func(x link, y linkSet) bool {
		for y1 := range y.links {
			if y1.to != l.to {
				return false
			}
		}
		return true
	}
	// Decide what kind of conflict needs to be resolved
	//
	if to(l, s) {
		// Pick the earliest expected phoneme - that is the smallest from value
		//
		ret := l
		for k := range s.links {
			if k.from < ret.from {
				ret = k
			}
		}
		debug("resolve<-:", ret)
		return ret
	}
	// if it's not a to conflict then this is either a from conflict or a
	// crossover conflict. Either way pick the link with the largest to range
	//
	ret := l
	diff := c[ret.to].end - c[ret.to].start
	// diff := p.data[ret.to].end - p.data[ret.to].start
	for k := range s.links {
		thisDiff := c[k.to].end - c[k.to].start
		// thisDiff := p.data[k.to].end - p.data[k.to].start
		if thisDiff > diff {
			ret = k
			diff = thisDiff
		}
	}
	debug("resolve<-:", ret)
	return ret
}

func judgement(phons []phoneme, candidates []candidatePhoneme, links linkWithConfidenceSet) []phonemeVerdict {
	debug("judgement->: candidates", candidates)
	verdicts := []phonemeVerdict{}
	if len(candidates) == 0 {
		return verdicts
	}
	findLink := func(from int, links linkWithConfidenceSet) (linkWithConfidence, bool) {
		for link := range links {
			if link.from == from {
				return link, true
			}
		}
		return linkWithConfidence{}, false
	}
	for i, phon := range phons {
		if link, ok := findLink(i, links); ok {
			v := phonemeVerdict{
				candidates[link.to].psPhonemeDatum,
				i,
				good,
			}
			if link.confidence == 2 {
				v.goodBadEtc = possible
			}
			debug("Adding verdict,", v)
			verdicts = append(verdicts, v)
		} else {
			//  Can't do this! There is no candidate phoneme datum because the
			// phoneme wasn't found
			v := phonemeVerdict{
				psPhonemeDatum{
					phon,
					0,
					0,
				},
				i,
				missing,
			}
			debug("Adding verdict,", v)
			verdicts = append(verdicts, v)
		}
	}
	debug("judgement<-")
	return verdicts
}

func unexpectedInRange(from, to int, candidates []candidatePhoneme) []phonemeVerdict {
	unexpected := []phonemeVerdict{}
	for _, candidate := range candidates {
		if candidate.start >= from && candidate.end <= to && candidate.confidence >= 3 {
			v := phonemeVerdict{
				candidate.psPhonemeDatum,
				-1,
				surprise,
			}
			unexpected = append(unexpected, v)
		}
	}
	return unexpected
}

func insertUnexpected(candidates []candidatePhoneme, verdicts []phonemeVerdict) []phonemeVerdict {
	finalVerdicts := []phonemeVerdict{}
	if len(candidates) == 0 {
		return verdicts
	}
	from := 0
	for _, verdict := range verdicts {
		if verdict.goodBadEtc == missing && verdict.start == 0 && verdict.end == 0 {
			finalVerdicts = append(finalVerdicts, verdict)
			continue
		}
		to := verdict.start
		unexpected := unexpectedInRange(from, to, candidates)
		finalVerdicts = append(finalVerdicts, unexpected...)
		finalVerdicts = append(finalVerdicts, verdict)
		from = verdict.end
	}
	// Look for anything surprising after the last verdict
	//
	// Look for last phoneme in all results
	//
	to := candidates[len(candidates)-1].end
	unexpected := unexpectedInRange(from, to, candidates)
	finalVerdicts = append(finalVerdicts, unexpected...)

	return finalVerdicts
}

/*
func process(word string, phons []phoneme, results []psPhonemeResults) []phonemeVerdict {
  candidates := merge(results)
  links := createLinkWithConfidenceSet(candidates, phons)
  remainingLinks := pruneWithConfidence(links, candidateData(candidates))
  debug("remainingLinks =", remainingLinks)
  verdicts := judgement(phons, candidates, remainingLinks)
  debug("verdicts =", verdicts)
  verdicts = insertUnexpected(candidates, verdicts)
  debug("verdicts(final) =", verdicts)
  return verdicts
}
*/
