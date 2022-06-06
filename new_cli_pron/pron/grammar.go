package pron

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/colinarticulate/dictionary"
)

type noGuardsError struct {
	context       string
	this          phoneme
	previous      phoneme
	next          phoneme
	lastPenalties []phoneme
	grammar       jsgfGrammar
}

func (g noGuardsError) Error() string {
	str := fmt.Sprintln("No guards! " + g.context)
	if g.this != "" {
		str += fmt.Sprintln("This phoneme =", g.this)
	}
	if g.previous != "" {
		str += fmt.Sprintln("Previous phoneme =", g.previous)
	}
	if g.next != "" {
		str += fmt.Sprintln("Next phoneme =", g.next)
	}
	if len(g.lastPenalties) != 0 {
		str += fmt.Sprintln("Last penalties =", g.lastPenalties)
	}
	str += fmt.Sprintln("Grammar =", g.grammar)
	return str
}

type parserError int

const (
	parseFailed parserError = iota
)

func (p parserError) Error() string {
	return "Parser error"
}

type rule interface {
	generate() string
}

type parseResult struct {
	start, end   int
	phonemeFound bool
}

type parsableRule interface {
	rule
	parse(data []psPhonemeDatum, i int) ([]parseResult, error)
}

type ruleName string

const (
	r_aa  = "<aa>"
	r_ae  = "<ae>"
	r_ah  = "<ah>"
	r_ao  = "<ao>"
	r_aw  = "<aw>"
	r_ax  = "<ax>" // 18 Nov 2020
	r_ay  = "<ay>"
	r_b   = "<b>"
	r_ch  = "<ch>"
	r_dh  = "<dh>"
	r_d   = "<d>"
	r_eh  = "<eh>"
	r_ehr = "<ehr>" // 8 March 2021
	r_er  = "<er>"
	r_ey  = "<ey>"
	r_f   = "<f>"
	r_g   = "<g>"
	r_hh  = "<hh>"
	r_ih  = "<ih>"
	r_iy  = "<iy>"
	r_jh  = "<jh>"
	r_k   = "<k>"
	r_l   = "<l>"
	r_m   = "<m>"
	r_ng  = "<ng>"
	r_n   = "<n>"
	r_oh  = "<oh>" //18 Nov 2020
	r_ow  = "<ow>"
	r_oy  = "<oy>"
	r_p   = "<p>"
	r_r   = "<r>"
	r_ss  = "<ss>" // Is this really an s phoneme??
	r_sh  = "<sh>"
	r_sil = "<sil>"
	r_t   = "<t>"
	r_th  = "<th>"
	r_uh  = "<uh>"
	r_uw  = "<uw>"
	r_v   = "<v>"
	r_w   = "<w>"
	r_y   = "<y>"
	r_z   = "<z>"
	r_zh  = "<zh>" // axL, axM, axN, kS, Kw, dZ, tS all added 18th March 2021

	r_axl = "<axl>"
	r_axn = "<axn>"
	r_axm = "<axm>"
	//r_nd = "<nd>"  	// 9th March 2021
	r_ks = "<ks>"
	r_kw = "<kw>"
	r_dz = "<dz>"
	r_ts = "<ts>"
	r_kl = "<kl>"
	r_pl = "<pl>"
	r_bl = "<bl>"

	r_uwl = "<uwl>"
	r_uwn = "<uwn>"
	r_uwm = "<uwm>"
	r_axr = "<axr>"

	r_kr  = "<kr>"
	r_pr  = "<pr>"
	r_gr  = "<gr>"
	r_tr  = "<tr>"
	r_ing = "<ing>"

	// Added 14th Feb 2022
	r_ihl = "<ihl>"
	r_ehl = "<ehl>"
	r_fl  = "<fl>"

	//remove for 64only
	r_thr = "<thr>"
	r_yuw = "<yuw>"
	r_sts = "<sts>"
	r_st  = "<st>"
	r_kt  = "<kt>"

	// commented out for the 64Only phoneme set
	//St
	//StS
	//Yuw
	//axM
	//kT
	//thR
	//uwM

	r_any_vowel_noSlide = "<any_vowel_noSlide>"
	r_check             = "<check>"
	r_wild_consonant    = "<wild_consonant>"
	r_wild_vowel        = "<wild_vowel>"
)

type blahRule struct {
	name ruleName
	exp  []phonOrRul
}

var blah = blahRule{
	r_aa,
	[]phonOrRul{
		{
			phonT,
			aa,
			"",
		},
		{
			ruleT,
			"",
			r_aa,
		},
	},
}

type phonOrRuleTag string

const (
	phonT = "phon"
	ruleT = "rule"
)

type phonOrRul struct {
	tag  phonOrRuleTag
	phon phoneme
	exp  ruleName
}

type optOpRul struct {
	exp []phonOrRul
}

func (r optOpRul) generate() string {
	return ""
}

type targetRule struct {
	name  string
	rules []parsableRule
}

func (r targetRule) generate() string {
	return ""
}

type orRule struct {
	name  string
	rules []rule
}

func (r orRule) generate() string {
	return ""
}

type seqRule struct {
	name ruleName
	exp  []phonOrRul
}

func (r seqRule) generate() string {
	return ""
}

type phonRule struct {
	ph phoneme
}

func new_R_ph(p phoneme) phonRule {
	return phonRule{
		p,
	}
}

func (r phonRule) generate() string {
	return fmt.Sprintf(string(r.ph))
}

type jsgfGrammar struct {
	header, grammar string
	target          R_target
	rules           map[string]namedRule
	config          jsgfConfig
}

type nullRule struct {
}

func (r nullRule) generate() string {
	return ""
}

func (r nullRule) parse(data []psPhonemeDatum, i int) ([]parseResult, error) {
	result := parseResult{
		i, i, false,
	}
	return []parseResult{
		result,
	}, nil
}

func (g jsgfGrammar) getRuleNamed(name string) (namedRule, bool) {
	if r, ok := g.rules[name]; ok {
		return r, true
	}
	return namedRule{"", nullRule{}}, false
}

func newTargetRule(word string, dict dictionary.Dictionary) rule {

	return targetRule{}
}

func (g jsgfGrammar) SaveToDisk(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		_debug("SaveToDisk: failed to create file. err =", err)
		log.Panic()
	}
	defer f.Close()

	content := g.header + ";\n"
	content += g.grammar + ";\n"
	content += g.target.generate() + ";\n"
	for _, rule := range g.rules {
		content += rule.derivation() + ";\n"
	}
	_, err = f.WriteString(content)
	// _debug(g.target.generate())
	_debug("")
	if err != nil {
		// Not sure there's much else we can do if the write fails
		//
		_debug("SaveToDisk: failed to write to file. err =", err)
		log.Panic()
	}
}

func (g jsgfGrammar) SaveToByteSlice() []byte {

	content := g.header + ";\n"
	content += g.grammar + ";\n"
	content += g.target.generate() + ";\n"
	for _, rule := range g.rules {
		content += rule.derivation() + ";\n"
	}

	return []byte(content)
}

type consonantGroup int

const (
	bilabial consonantGroup = iota
	labiodental
	dental
	alveolar
	post_alveolar
	retroflex
	palatal
	velar
	glottal

	plosive
	nasal
	fricative
	approximant

	fricative_voiced
	fricative_voiceless

	plosive_voiced
	plosive_voiceless

	sss
	zzz
)

var consonantsByGroup = map[consonantGroup][]phoneme{

	bilabial: {
		b, m, p,
	},
	labiodental: {
		f, v,
	},
	dental: {
		//t, d, n, l, r, th, dh,
		t, d, n, l, r, th, dh, thr, kt, st, sts,
		//64only
		//t, d, n, l, r, th, dh,
	},
	alveolar: {
		//d, l, n, s, t, z, r, dz, ts,
		d, l, n, s, t, z, r, dz, ts, st, sts, kt, thr,
		//64only
		//d, l, n, s, t, z, r, dz, ts,
	},
	post_alveolar: {
		//t, d, n, l, r, sh, zh,
		t, d, n, l, r, sh, zh, kt, st, sts,
	},
	retroflex: {
		//t, d, ng, s, z, r, l, dz, ts, ing,
		t, d, ng, s, z, r, l, dz, ts, ing, sts, st, fl, kt,
		//64only
		//t, d, ng, s, z, r, l, dz, ts, ing, fl,
	},
	palatal: {
		ch, jh, // what about ch and jh?   .... removing y, placing in vowels
	},
	velar: {
		g, k, ng, w, kw, ks, ing, // x is missing from CMUBET
	},
	glottal: {
		hh,
	},

	plosive: {
		//b, d, g, jh, ch, k, p, t, kl, pl, bl, kr, gr, tr, pr,

		b, d, g, jh, ch, k, p, t, kl, pl, bl, kr, gr, tr, pr, kt, st, sts,
		//64only
		//b, d, g, jh, ch, k, p, t, kl, pl, bl, kr, gr, tr, pr,
	},
	nasal: {
		m, n, ng, ing,
	},
	fricative: {
		//dh, hh, v, z, zh, f, hh, s, sh, th, // should jh be in this list?

		dh, hh, v, z, zh, f, hh, s, sh, th, thr, fl,
		//64only
		//dh, hh, v, z, zh, f, hh, s, sh, th, fl,
	},
	approximant: {
		l, r, w, // moving y to vowel
	},

	fricative_voiced: {
		dh, hh, v, z, zh,
	},
	fricative_voiceless: {
		f, hh, s, sh, th, // x is missing from CMUBET
	},

	plosive_voiced: {
		b, d, g, jh,
	},
	plosive_voiceless: {
		ch, k, p, t,
	},

	sss: {
		s,
	},
	zzz: {
		z,
	},
}

func isPlosive(ph phoneme) bool {
	plosives := consonantsByGroup[plosive]
	for _, el := range plosives {
		if el == ph {
			return true
		}
	}
	return false
}

/*
var groupsByConsonant = map[phoneme][]consonantGroup{
  b: {
    labial, plosive_voiced,
  },
  ch: {
    plosive_voiceless, post_alveolar,
  },
  d: {
    alveolar, plosive_voiced,
  },
  dh: {
    dental, fricative_voiced,
  },
  f: {
    labial, fricative_voiceless,
  },
  g: {
    velar, plosive_voiced,
  },
  hh: {
    glottal, fricative_voiced, fricative_voiceless,
  },
  jh: {
    plosive_voiced, post_alveolar,
  },
  k: {
    plosive_voiceless, velar,
  },
  l: {
    alveolar, approximant,
  },
  m: {
    labial, nasal,
  },
  n: {
    alveolar, nasal,
  },
  ng: {
    nasal, velar,
  },
  p: {
    labial, plosive_voiceless,
  },
  r: {
    approximant, post_alveolar,
  },
  s: {
    alveolar, fricative_voiceless, sss,
  },
  sh: {
    fricative_voiceless, post_alveolar,
  },
  t: {
    alveolar, plosive_voiceless,
  },
  th: {
    dental, fricative_voiceless,
  },
  v: {
    fricative_voiced, labial,
  },
  w: {
    approximant, velar,
  },
  y: {
    approximant, palatal,
  },
  z: {
    alveolar, fricative_voiced,
  },
  zh: {
    fricative_voiced, post_alveolar,
  },
}
*/

var groupsByConsonant = map[phoneme][]consonantGroup{
	b: {
		bilabial, plosive_voiced,
	},
	ch: {
		plosive, palatal,
	},
	d: {
		alveolar, plosive,
	},
	dh: {
		dental, fricative,
	},
	f: {
		labiodental, fricative,
	},
	g: {
		velar, plosive,
	},
	hh: {
		glottal, fricative,
	},
	jh: {
		palatal, fricative,
	},
	k: {
		plosive, velar,
	},
	l: {
		//alveolar, approximant,
		dental, alveolar, post_alveolar, approximant,
	},
	m: {
		bilabial, nasal, labiodental, //  not really labiodental but may be useful to block out guards
	},
	n: {
		//alveolar, nasal,
		dental, alveolar, post_alveolar, nasal,
	},
	//nd: {
	//  dental, alveolar, post_alveolar, nasal,
	//},
	ng: {
		//nasal, velar,
		nasal, velar, retroflex,
	},
	p: {
		bilabial, plosive,
	},
	r: {
		//approximant, alveolar,
		dental, alveolar, post_alveolar, approximant,
	},
	s: {
		alveolar, fricative, sss,
	},
	sh: {
		fricative, post_alveolar,
	},
	t: {
		dental, alveolar, post_alveolar, plosive,
	},
	th: {
		dental, fricative,
	},
	v: {
		fricative, labiodental,
	},
	w: {
		approximant, velar,
	},
	y: {
		approximant, palatal,
	},
	z: {
		alveolar, fricative,
	},
	zh: {
		fricative, post_alveolar,
	},

	dz: {
		alveolar, retroflex,
	},
	ts: {
		alveolar, retroflex,
	},

	ks: {
		velar,
	},
	kw: {
		velar,
	},
	kl: {
		plosive,
	},
	pl: {
		plosive,
	},
	bl: {
		plosive,
	},
	kr: {
		plosive,
	},
	pr: {
		plosive,
	},
	gr: {
		plosive,
	},
	tr: {
		plosive,
	},
	ing: {
		//nasal, velar,
		nasal, velar, retroflex,
	},

	sts: {
		alveolar, retroflex, plosive, dental,
	},
	st: {
		alveolar, retroflex, plosive, dental,
	},

	thr: {
		fricative, alveolar, dental,
	},
	kt: {
		plosive, dental,
	},
	fl: {
		retroflex, fricative,
	},
}

/*
  axl: {
    alveolar,
  },
  //axm: {
  //  alveolar,
  //},
  axn: {
    alveolar,
  },
  axm: {
    alveolar,
  },
  uwl: {
    alveolar,
  },
  uwn: {
    alveolar,
  },
  uwm: {
    alveolar,
  },
*/

type R_or struct {
	rs []rule
}

func new_R_or(rs ...rule) R_or {
	return R_or{
		rs,
	}
}

func (r R_or) generate() string {
	if len(r.rs) == 0 {
		// Should never happen
		return ""
	}
	str := r.rs[0].generate()
	for _, r1 := range r.rs[1:] {
		str += " | " + r1.generate()
	}
	if len(r.rs) != 1 {
		// Don't overdo the brackets. We only need an outer pair of brackets if
		// there's more than one rule to or
		str = "(" + str + ")"
	}
	return fmt.Sprintf(str)
}

type R_opt struct {
	r1 rule
}

func new_R_opt(r rule) R_opt {
	return R_opt{
		r,
	}
}

func (r R_opt) generate() string {
	return fmt.Sprintf("[" + r.r1.generate() + "]")
}

type R_and struct {
	rs []rule
}

func new_R_and(rs ...rule) R_and {
	return R_and{
		rs,
	}
}

func (r R_and) generate() string {
	if len(r.rs) == 0 {
		// Should never happen
		return ""
	}
	var str = ""
	str += r.rs[0].generate()
	for _, r1 := range r.rs[1:] {
		str += " " + r1.generate()
	}
	if len(r.rs) != 1 {
		// Don't overdo the brackets. We only need an outer pair of brackets if
		// there's more than one rule to and
		str = "(" + str + ")"
	}
	return str
}

type namedRule struct {
	name string // Also used to find lookup the rule definition
	rule
}

func (g *jsgfGrammar) new_R_named(name string, ruleGen func() rule) namedRule {
	if r, ok := g.rules[name]; ok {
		return r
	}
	new := namedRule{
		name,
		ruleGen(),
	}
	g.rules[name] = new
	return new
}

func (r namedRule) generate() string {
	return fmt.Sprintf("<" + r.name + ">")
}

func (r namedRule) derivation() string {
	str := r.generate()
	str += " = "
	str += r.rule.generate()

	return str
}

func newConsonants() map[phoneme]bool {
	consonants := map[phoneme]bool{
		//b: true, ch: true, d: true, dh: true, f: true, g: true, hh: true, jh: true, k: true, l: true, m: true, n: true, nd: true, ng: true, p: true, r: true, s: true, sh: true, st: true, t: true, th: true, v: true, w: true, z: true, zh: true, dz: true, ts: true, axn: true, axm: true, axl: true, kw: true, ks: true,
		//b: true, ch: true, d: true, dh: true, f: true, g: true, hh: true, jh: true, k: true, l: true, m: true, n: true, ng: true, p: true, r: true, s: true, sh: true, t: true, th: true, v: true, w: true, z: true, zh: true, dz: true, ts: true, axn: true, axl: true, kw: true, ks: true, uwl: true, sts: true,
		//b: true, ch: true, d: true, dh: true, f: true, g: true, hh: true, jh: true, k: true, l: true, m: true, n: true, ng: true, p: true, r: true, s: true, sh: true, t: true, th: true, v: true, w: true, z: true, zh: true, dz: true, kw: true, ks: true, ts: true, kl: true, pl: true, bl: true, tr: true, kr: true, pr: true, gr: true, ing: true,

		b: true, ch: true, d: true, dh: true, f: true, g: true, hh: true, jh: true, k: true, l: true, m: true, n: true, ng: true, p: true, r: true, s: true, sh: true, t: true, th: true, v: true, w: true, z: true, zh: true, dz: true, kw: true, ks: true, ts: true, kl: true, pl: true, bl: true, tr: true, kr: true, pr: true, gr: true, ing: true, thr: true, sts: true, st: true, kt: true, fl: true,

		//64only
		//b: true, ch: true, d: true, dh: true, f: true, g: true, hh: true, jh: true, k: true, l: true, m: true, n: true, ng: true, p: true, r: true, s: true, sh: true, t: true, th: true, v: true, w: true, z: true, zh: true, dz: true, kw: true, ks: true, ts: true, kl: true, pl: true, bl: true, tr: true, kr: true, pr: true, gr: true, ing: true, fl: true,
	}
	return consonants
}

func consRemove(cons map[phoneme]bool, remove ...phoneme) map[phoneme]bool {
	ret := cons
	for _, c := range remove {
		ret[c] = false
	}
	return ret
}

func consonantsInGroup(phons []phoneme) []phoneme {
	inGroup := []phoneme{}
	for _, ph := range phons {
		if groups, ok := groupsByConsonant[ph]; ok {
			for _, g := range groups {
				phG, ok := consonantsByGroup[g]
				if !ok {
					continue
				}
				inGroup = append(inGroup, phG...)
			}
		} else {
			// Should never get here. We're missing an entry for the phoneme in
			// groupsByConsonant if we do, - or passed a vowel into r.phons!
		}
	}
	return inGroup
}

func pickConsonant(cons map[phoneme]bool) phoneme {
	for k, v := range cons {
		if v {
			return k
		}
	}
	return sil
}

type R_Cx struct {
	namedRule
}

func (g *jsgfGrammar) new_named_R(ruleGen func() (string, rule)) namedRule {
	name, r := ruleGen()
	// Check to see if we already have a rule with this name. If we do just
	// return
	if r, ok := g.rules[name]; ok {
		return r
	}
	new := namedRule{
		name,
		r,
	}
	g.rules[name] = new
	return new
}

func (gr *jsgfGrammar) new_R_Cx(exceptPhons ...phoneme) R_Cx {
	// Create a name for the rule
	name := "any_Cx"
	for _, ph := range exceptPhons {
		name += "_" + string(ph)
	}
	r := gr.new_R_named(name, func() rule {
		allCons := newConsonants()
		remaining := consRemove(allCons, exceptPhons...)

		remaining = consRemove(remaining, n, ng, ing) // removing n and ng so that when an "n" type sound (eg n, m, ng) is identified these effectively grouped to m .... should help with Russia(n) by Khurrum

		//remaining = consRemove(remaining, p, t, k, f, th, s, sh, st, jh, r, y, dz, ts, axn, axm, axl, kw, ks)
		//remaining = consRemove(remaining, p, t, k, f, th, s, sh, jh, r, y, dz, ts, axn, axl, kw, ks, uwl, sts)
		//remaining = consRemove(remaining, p, t, k, f, th, s, sh, jh, r, y, dz, kw, ks, ts, kl, pl, bl, tr, pr, "gr", kr)

		/*            // groups to reduce any_consonant by ..... the reason to do so is to help the different scan return a similar surprise
		p, b,
		t, d,
		k, g
		m, n, ng,
		f, v,
		th, dh,
		s, z,
		sh, zh,
		ch,
		jh,
		w, r, l, y
		hh
		*/

		// remaining := filterConsonants(exceptPhons)
		// Let's get the phonemes left as an array of phoneme rules
		phons := []rule{}
		for k, v := range remaining {
			if v {
				new := phonRule{
					k,
				}
				phons = append(phons, rule(new))
			}
		}
		return new_R_or(phons...)
	})
	return R_Cx{
		r,
	}
}

func (r R_Cx) generate() string {
	return r.namedRule.generate()
}

type R_Vx struct {
	namedRule
}

func (g *jsgfGrammar) new_R_Vx(exceptPhons ...phoneme) R_Vx {
	// Create a name for the rule
	name := "any_Vx"
	for _, ph := range exceptPhons {
		name += "_" + string(ph)
	}
	r := g.new_R_named(name, func() rule {
		vowels := map[phoneme]bool{
			//aa: true, ae: true, ah: true, ao: true, aw: true, ax: true, ay: true, eh: true, ehr: true, er: true, ey: true, ih: true, iy: true, oh: true, ow: true, oy: true, uw: true, uh: true, y: true, axl: true, axm: true, axn: true, uwn: true, uwl: true, uwm: true, axr: true,

			aa: true, ae: true, ah: true, ao: true, aw: true, ax: true, ay: true, eh: true, ehr: true, er: true, ey: true, ih: true, iy: true, oh: true, ow: true, oy: true, uw: true, uh: true, y: true, axl: true, axm: true, axn: true, uwn: true, uwl: true, uwm: true, axr: true, yuw: true, ihl: true, ehl: true,

			//64only
			//aa: true, ae: true, ah: true, ao: true, aw: true, ax: true, ay: true, eh: true, ehr: true, er: true, ey: true, ih: true, iy: true, oh: true, ow: true, oy: true, uw: true, uh: true, y: true, axl: true, axn: true, uwn: true, uwl: true, axr: true, ihl: true, ehl: true,
		}
		remaining := phonDiff(vowels, exceptPhons)
		phons := []rule{}
		for k, v := range remaining {
			if v {
				new := phonRule{
					k,
				}
				phons = append(phons, rule(new))
			}
		}
		return new_R_or(phons...)
	})
	return R_Vx{
		r,
	}
}

func (r R_Vx) generate() string {
	return r.namedRule.generate()
}

func phonDiff(phons map[phoneme]bool, diffs []phoneme) map[phoneme]bool {
	ret := phons
	for _, d := range diffs {
		ret[d] = false
	}
	return ret
}

type R_V_noSlide_X struct {
	namedRule
}

func (g *jsgfGrammar) new_R_noSlide_X(exceptPhons ...phoneme) R_V_noSlide_X {
	// Create a name for the rule
	name := "any_Vx"
	for _, ph := range exceptPhons {
		name += "_" + string(ph)
	}
	name += "_noSlide"
	r := g.new_R_named(name, func() rule {
		noSlideVowels := map[phoneme]bool{

			//aa: true, ae: true, ah: true, ao: true, ax: true, eh: true, ehr: true, er: true, ih: true, iy: true, oh: true, uw: true, uh: true, y: true, axl: true, axm: true, axn: true, uwn: true, uwl:  true, uwm: true, axr: true,

			// possibly this should have no compound phonemes (or ax)
			//aa: true, ae: true, ah: true, ao: true, ax: true, eh: true, er: true, ih: true, iy: true, oh: true, uw: true, uh: true, y: true,
			// remove er (23 Feb 22) as er is a slide/diphthong
			aa: true, ae: true, ah: true, ao: true, ax: true, eh: true, ih: true, iy: true, oh: true, uw: true, uh: true, y: true,

			// alternately, have to add the new compound phonemes
			//aa: true, ae: true, ah: true, ao: true, ax: true, eh: true, ehr: true, er: true, ih: true, iy: true, oh: true, uw: true, uh: true, y: true, axl: true, axm: true, axn: true, uwn: true, uwl:  true, uwm: true, axr: true, yuw: true, ihl: true, ehl: true,

		}
		remaining := phonDiff(noSlideVowels, exceptPhons)
		phons := []rule{}
		for k, v := range remaining {
			if v {
				new := phonRule{
					k,
				}
				phons = append(phons, rule(new))
			}
		}
		return new_R_or(phons...)
	})
	return R_V_noSlide_X{
		r,
	}
}

func (r R_V_noSlide_X) generate() string {
	return r.namedRule.generate()
}

type R_V_noSlide struct {
	namedRule
}

//                                      _
//  __ _ _ _ _  _   __ _______ __ _____| |
// / _` | ' \ || |  \ V / _ \ V  V / -_) |
// \__,_|_||_\_, |   \_/\___/\_/\_/\___|_|
//           |__/
// any_vowel

func (g *jsgfGrammar) new_R_V_noSlide() R_V_noSlide {
	name := "any_vowel_noSlide"
	r := g.new_R_named(name, func() rule {
		noSlideVowels := []phoneme{
			//aa, ae, ah, ax, eh, ehr, er, ih, iy, oh, uw, uh, y, axnoise, ernoise, axl, axm, axn, uwn, uwl, uwm,            // removing y from consonants and adding in here as a vowel - PE 25/1/2020
			// aa, ae, ah, ax, eh, ehr, er, ih, iy, oh, uw, uh, y, axl, axm, axn, uwn, uwl, uwm, // removing y from consonants and adding in here as a vow

			// removing ax, axm, axn, axl, uwl, uwm, uwn because they seem to be interferring (being surprises) when guards, especially on words ending with AX
			// PE 2nd Jan 2022
			// Not adding axr - 5 Jan 2022  (just because if the above have already been removed, it doesn't make sense to add it in)
			// Not including any compound phoneme, ax or any diphthong
			//aa, ae, ah, eh, er, ih, iy, oh, uw, uh, y,
			// removing er (23 Feb 22) as it's a slide/diphthong
			aa, ae, ah, eh, ih, iy, oh, uw, uh, y,
		}
		phons := []rule{}
		for _, v := range noSlideVowels {
			new := phonRule{
				v,
			}
			phons = append(phons, rule(new))
		}

		return new_R_or(phons...)
	})
	return R_V_noSlide{
		r,
	}
}

func (r R_V_noSlide) generate() string {
	return r.namedRule.generate()
}

type R_softFade struct {
	namedRule
}

func (g *jsgfGrammar) new_R_softFade() R_softFade {
	name := "soft_fade"
	r := g.new_R_named(name, func() rule {
		return new_R_and(new_R_ph(p), new_R_ph(f))
	})
	return R_softFade{
		r,
	}
}

func (r R_softFade) generate() string {
	return r.namedRule.generate()
}

type R_C struct {
	namedRule
}

func (gr *jsgfGrammar) new_R_C() R_C {
	name := "any_C"
	r := gr.new_R_named(name, func() rule {
		consonants := []phoneme{
			// b, ch, d, dh, f, g, hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh,
			//b, ch, d, dh, f, g, hh, jh, k, l, m, n, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, "gr", pr, kr, ing,

			b, ch, d, dh, f, g, hh, jh, k, l, m, n, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, "gr", pr, kr, ing, thr, sts, st, kt, fl,

			//64only
			b, ch, d, dh, f, g, hh, jh, k, l, m, n, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, "gr", pr, kr, ing, fl,

			/*            // groups to reduce any_consonant by ..... the reason to do so is to help the different scan return a similar surprise
			p, b,
			d, t,
			k, g
			m, n, ng,
			f, v,
			th, dh,
			s, z,
			sh, zh,
			ch,
			jh,
			w, r, l, y
			hh
			*/

			// remove p, k ,t as voiceless-polosives .... also helps resolve the d t s issues as "t" is gone
			// remove m & ng ... leaving n to coalesce around

			//b, ch, d, dh, f, g, hh, jh, l, n, s, sh, th, v, z, zh,

		}
		phons := []rule{}
		for _, c := range consonants {
			new := phonRule{
				c,
			}
			phons = append(phons, rule(new))
		}
		return new_R_or(phons...)
	})
	return R_C{
		r,
	}
}

func (r R_C) generate() string {
	return r.namedRule.generate()
}

func allIndexes(data []psPhonemeDatum, ph phoneme) []int {
	ret := []int{}
	for i, datum := range data {
		if datum.phoneme == ph {
			ret = append(ret, i)
		}
	}
	return ret
}

func getPhonInResultAt(res psPhonemeResults, start int) (phoneme, bool) {
	if start < 0 || len(res.data)-1 < start {
		return "", false
	}
	return res.data[start].phoneme, true
}

func inResultInRange(res psPhonemeResults, ph phoneme, start, end int) []int {
	if start < 0 || start+end > len(res.data)-1 {
		return []int{}
	}
	ret := []int{}
	for i := start; i <= start+end; i++ {
		if res.data[i].phoneme == ph {
			ret = append(ret, i)
		}
	}
	return ret
}

//===================================================================
//    ___                 _                        _
//   / _ \ _ __  ___ _ _ (_)_ _  __ _     _ _ _  _| |___
//  | (_) | '_ \/ -_) ' \| | ' \/ _` |   | '_| || | / -_)
//   \___/| .__/\___|_||_|_|_||_\__, |   |_|  \_,_|_\___|
//        |_|                   |___/
//
//===================================================================

type R_opening struct {
	guard phoneme
	rule
}

func (r R_opening) generate() string {
	return r.rule.generate()
}

func (r R_opening) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
	if len(res) < i+3 {
		// The number of phonemes after i should be at least 3 - eg. sil ay sil
		return []parseResult{}, parseFailed
	}
	// One result is a sil
	if res[i].phoneme != sil {
		// The first result should always be a sil
		return []parseResult{}, parseFailed
	}
	ret := []parseResult{}
	// The opening rule can always just be a sil
	parsRes := parseResult{
		i, i, false,
	}
	ret = append(ret, parsRes)
	// Either the guard will appear in the two phonemes immediately following sil
	// ()== phon[0]) or it won't. If it does, add it as a possible result
	if r.guard == res[i+1].phoneme || r.guard == res[i+2].phoneme {
		parsRes := parseResult{
			i, i + 2, false,
		}
		ret = append(ret, parsRes)

		if len(res) > i+5 {
			// The guard can also appear in the next two phonemes.
			if r.guard == res[i+3].phoneme || r.guard == res[i+4].phoneme {
				parsRes := parseResult{
					i, i + 4, false,
				}
				ret = append(ret, parsRes)
			}
		}
	}
	return ret, nil
}

func (gr *jsgfGrammar) openingRule(phons []phoneme) parsableRule {
	if len(phons) < 1 {
		return nullRule{}
	}
	targetPh := phons[0]
	r_pen, ok := gr.config.openingPenalty(targetPh)
	if !ok {
		str := fmt.Sprintln("preRule: no guard!")
		str += fmt.Sprintln("phons[0] =", phons[0])
		lastGuards := gr.config.getPens()
		str += fmt.Sprintln("lastGuards =", fmt.Sprint(lastGuards))
		str += fmt.Sprintln("grammar =", gr)
		_debug(str)
		log.Panic()

		return nullRule{}
	}
	var r_v, r_c rule
	exceptPhons := []phoneme{
		targetPh, r_pen.ph,
	}
	exceptPhons = append(exceptPhons, gr.config.getPens()...)
	// Configure the opening rule phoneme based on the first phoneme in the word
	/*
	  var r_openPh rule
	  r_openPh = nullRule{}
	  openPh:= map[phoneme]phoneme {}
	  if oPh, ok := openPh[phons[0]]; ok {
	    r_openPh = phonRule{
	      oPh,
	    }
	  }
	*/
	var retRule rule
	if isVowel(targetPh) {
		r_v = gr.new_R_noSlide_X(exceptPhons...)
		r_c = gr.new_R_Cx(exceptPhons...)
		retRule = new_R_and(new_R_ph(sil), new_R_opt(new_R_or(new_R_and(r_pen, r_v), new_R_and(r_v, r_pen))), new_R_opt(new_R_or(new_R_and(r_pen, r_c), new_R_and(r_c, r_pen))))
	} else {
		r_v = gr.new_R_V_noSlide()
		r_c = gr.new_R_Cx(exceptPhons...)
		retRule = new_R_and(new_R_ph(sil), new_R_opt(new_R_or(new_R_and(r_pen, r_c), new_R_and(r_c, r_pen))), new_R_opt(new_R_or(new_R_and(r_pen, r_v), new_R_and(r_v, r_pen))))
	}
	gr.config.addPen(r_pen.ph)

	oRule := R_opening{
		// phons[0],
		r_pen.ph,
		retRule,
	}
	// return retRule
	return oRule
}

type R_trap struct {
	trap phoneme
}

func (r R_trap) generate() string {
	return "sil " + string(r.trap)
}

type R_trappedOpening struct {
	rT R_trap
	rO R_opening
}

// Need to be a bit careful here. Just because we find sil r.rT.trap ...
// doesn't necessarily mean we have sil r.rT.trap sil ... The value of
// r.rT.trap could be te same as the first phoneme now that we are allowing
// other phonemes to represent the first phoneme. For a specific example,
// consider sil b where b is the trap and b is an alternate phoneme for p in
// the word P EH S T S...
func (r R_trappedOpening) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
	if i != 0 {
		return []parseResult{}, parseFailed
	}
	if len(res) < 2 {
		return r.rO.parse(res, i)
	}
	start := i
	if res[0].phoneme == sil && res[1].phoneme == r.rT.trap {
		start = i + 2
	}
	pRes_s, err := r.rO.parse(res, start)
	if err != nil {
		// Not necessarily a failure. Try this
		pRes_s, err = r.rO.parse(res, i)
		if err != nil {
			// NOW we have a parse failure...
			return []parseResult{}, parseFailed
		}
	}
	ret := []parseResult{}
	for _, pRes := range pRes_s {
		pRes.start = i
		ret = append(ret, pRes)
	}
	return ret, nil
}

//  _                  _
// (_)_ _  ___ ___ _ _| |_
// | | ' \(_-</ -_) '_|  _|
// |_|_||_/__/\___|_|  \__|
//

// This function inserts an extra P at the beginning of the grammar (even before s) to force the pocketsphinx engine into finding nothing
// (If removing then comment back in the function below)

//optional sill
func (r R_trappedOpening) generate() string {
	rTO := new_R_or(r.rO, new_R_and(r.rT, r.rO))
	return rTO.generate()
}

/*
//compulsory sil
func (r R_trappedOpening) generate() string {
  rTO := new_R_and(r.rT, r.rO)
  return rTO.generate()
}
*/

//     _ _    _
//  __(_) |  | |_ _ _ __ _ _ __
// (_-< | |  |  _| '_/ _` | '_ \
// /__/_|_|   \__|_| \__,_| .__/
//                        |_|
// sil_trap

func (g *jsgfGrammar) trappedOpeningRule(phons []phoneme) parsableRule {
	if len(phons) == 0 {
		return nullRule{}
	}

	// remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh) // remove everything except b & hh
	// remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, nd, ng, p, r, s, sh, st, t, th, v, w, z, zh, b, hh, dz, ts, axm, axn, axl, kw, ks) // remove everything
	//remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, b, hh, dz, ts, axn, axl, kw, ks, uwl, sts) // remove everything
	//remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, b, hh, dz, kw, ks, ts, kl, pl, bl, tr, pr, gr, kr, ing) // remove everything

	remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, b, hh, dz, kw, ks, ts, kl, pl, bl, tr, pr, gr, kr, ing, thr, sts, st, kt, fl) // remove everything

	//64only
	//remaining_pool := consRemove(newConsonants(), phons[0], ch, d, dh, f, "g", jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, b, hh, dz, kw, ks, ts, kl, pl, bl, tr, pr, gr, kr, ing, fl) // remove everything

	switch phons[0] {
	case k:
		remaining_pool[b] = true
	case p:
		remaining_pool[d] = true
	case b:
		remaining_pool[hh] = true
	case hh:
		remaining_pool[b] = true
	default:
		remaining_pool[p] = true
	}

	rT := R_trap{
		//pickConsonant(consRemove(newConsonants(), f, ch, phons[0], min_guard)),
		//pickConsonant(consRemove(newConsonants(), ch, phons[0], min_guard)),
		pickConsonant(remaining_pool),
	}

	g.config.setPen(rT.trap)

	rO := g.openingRule(phons)
	return R_trappedOpening{
		rT,
		rO.(R_opening),
	}
}

type R_phoneme struct {
	phon, guard phoneme
	otherPhons  []phoneme
	rule
}

func (r R_phoneme) generate() string {
	return r.rule.generate()
}

/*
func (r R_phoneme) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
  // First the easy bit...
  if r.phon == res[i].phoneme {
    parsRes := parseResult{
      i, i, true,
    }
    return []parseResult{
      parsRes,
    }, nil
  }
  ret := []parseResult{}
  indexes := allIndexes(res[i: min(i + 2, len(res))], r.guard)
  if len(indexes) == 0 {
    // This is an error. We found no phoneme and can't find a guard
    return []parseResult{}, parseFailed
  }
  // A single guard and phoneme is a perfectly valid result so add it to the results
  // to be returned
  parsRes := parseResult{
    i, i + 1, false,
  }
  ret = append(ret, parsRes)
  if i + 2 < len(res) - 1 && res[i + 2].phoneme == r.phon {
    // We found a phoneme so let's add it as a result and return
    parsRes = parseResult{
      i, i + 2, true,
    }
    ret = append(ret, parsRes)
    return ret, nil
  }
  // Search for another guard
  indexes = allIndexes(res[i + 2: min(i + 4, len(res))], r.guard)
  if len(indexes) == 1 {
    parsRes := parseResult{
      i, i + 3, false,
    }
    ret = append(ret, parsRes)
    if i + 4 < len(res) - 1 && res[i + 4].phoneme == r.phon {
      // We found a phoneme so let's add it as a result and return
      parsRes = parseResult{
        i, i + 4, true,
      }
      ret = append(ret, parsRes)
      return ret, nil
    }
    // Finally see if there are any more guards
    indexes = allIndexes(res[i + 4: min(i + 6, len(res))], r.guard)
    if len(indexes) == 1 {
      parsRes := parseResult{
        i, i + 5, false,
      }
      ret = append(ret, parsRes)
    }
  }
  return ret, nil
}
*/

func (r R_phoneme) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
	// First the easy bit...
	if r.phon == res[i].phoneme {
		parsRes := parseResult{
			i, i, true,
		}
		return []parseResult{
			parsRes,
		}, nil
	}
	ret := []parseResult{}
	if contains(r.otherPhons, res[i].phoneme) {
		parsRes := parseResult{
			i, i, false,
		}
		ret = append(ret, parsRes)
	}
	indexes := allIndexes(res[i:min(i+2, len(res))], r.guard)
	if len(indexes) == 0 {
		// This is an error. We found no phoneme, and can't find a guard
		return ret, nil
	}
	// A single guard and phoneme is a perfectly valid result so add it to the results
	// to be returned
	parsRes := parseResult{
		i, i + 1, false,
	}
	ret = append(ret, parsRes)
	if i+2 < len(res)-1 && res[i+2].phoneme == r.phon {
		// We found a phoneme so let's add it as a result and return
		parsRes := parseResult{
			i, i + 2, true,
		}
		ret = append(ret, parsRes)
		return ret, nil
	}
	if i+2 < len(res)-1 && contains(r.otherPhons, res[i+2].phoneme) {
		parsRes := parseResult{
			i, i + 2, true,
		}
		ret = append(ret, parsRes)
	}
	if i+2 < len(res)-1 {
		// Search for another guard
		indexes = allIndexes(res[i+2:min(i+4, len(res))], r.guard)
		if len(indexes) == 1 {
			parsRes := parseResult{
				i, i + 3, false,
			}
			ret = append(ret, parsRes)
			if i+4 < len(res)-1 && res[i+4].phoneme == r.phon {
				// We found a phoneme so let's add it as a result and return
				parsRes = parseResult{
					i, i + 4, true,
				}
				ret = append(ret, parsRes)
				return ret, nil
			}
			// Finally see if there are any more guards
			indexes = allIndexes(res[i+4:min(i+6, len(res))], r.guard)
			if len(indexes) == 1 {
				parsRes := parseResult{
					i, i + 5, false,
				}
				ret = append(ret, parsRes)
			}
		}
	}
	return ret, nil
}

type R_diphthongPhoneme struct {
	phon1, phon2, guard phoneme
	rule
}

func (r R_diphthongPhoneme) generate() string {
	return r.rule.generate()
}

func (r R_diphthongPhoneme) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
	// First the easy bit...
	// At the very least we expect to see the two phonemes that make up the
	// diphthong so if there aren't enough results return now
	if len(res)-1 < i+1 {
		return []parseResult{}, parseFailed
	}
	if r.phon1 == res[i].phoneme && r.phon2 == res[i+1].phoneme {
		parsRes := parseResult{
			i, i + 1, true,
		}
		return []parseResult{
			parsRes,
		}, nil
	}
	ret := []parseResult{}
	indexes := allIndexes(res[i:min(i+2, len(res))], r.guard)
	if len(indexes) == 0 {
		// This is an error. We found no phoneme and can't find a guard
		return []parseResult{}, parseFailed
	}
	// A single guard and phoneme is a perfectly valid result so add it to the results
	// to be returned
	parsRes := parseResult{
		i, i + 1, false,
	}
	ret = append(ret, parsRes)
	if i+3 < len(res)-1 && res[i+2].phoneme == r.phon1 && res[i+3].phoneme == r.phon2 {
		// We found a phoneme so let's add it as a result and return
		parsRes = parseResult{
			i, i + 3, true,
		}
		ret = append(ret, parsRes)
		return ret, nil
	}
	// Search for another guard
	indexes = allIndexes(res[i+2:min(i+4, len(res))], r.guard)
	if len(indexes) == 1 {
		parsRes := parseResult{
			i, i + 3, false,
		}
		ret = append(ret, parsRes)
		if i+5 < len(res)-1 && res[i+4].phoneme == r.phon1 && res[i+5].phoneme == r.phon2 {
			// We found a phoneme so let's add it as a result and return
			parsRes = parseResult{
				i, i + 5, true,
			}
			ret = append(ret, parsRes)
			return ret, nil
		}
		// Finally see if there are any more guards
		indexes = allIndexes(res[i+4:min(i+6, len(res))], r.guard)
		if len(indexes) == 1 {
			parsRes := parseResult{
				i, i + 5, false,
			}
			ret = append(ret, parsRes)
		}
	}
	return ret, nil
}

func (g *jsgfGrammar) preRule(phons []phoneme, currPh int, usePen phonRule) rule {
	if currPh < 0 || currPh > len(phons)-1 {
		return nullRule{}
	}

	exceptPhons := []phoneme{}
	if usePen != noPen {
		exceptPhons = append(exceptPhons, usePen.ph)
	}
	/*
	  if g.config.getLastPen() != sil {
	    exceptPhons = append(exceptPhons, g.config.getLastPen())
	  }
	*/
	exceptPhons = append(exceptPhons, g.config.getPens()...)
	exceptV_Phons := exceptPhons
	exceptC_Phons := exceptPhons

	currPhIsVowel := isVowel(phons[currPh])
	var r_x, r_v, r_c rule
	if currPhIsVowel {
		exceptV_Phons = append(exceptV_Phons, g.config.vowelHandler(phons[currPh])...)
		// Check for next phoneme
		if currPh+1 < len(phons) {
			if isVowel(phons[currPh+1]) {
				exceptV_Phons = append(exceptV_Phons, g.config.vowelHandler(phons[currPh+1])...)
			} else {
				exceptC_Phons = append(exceptC_Phons, phons[currPh+1])
			}
		}
		if currPh > 0 {
			if isVowel(phons[currPh-1]) {
				exceptV_Phons = append(exceptV_Phons, g.config.vowelHandler(phons[currPh-1])...)
			} else {
				exceptC_Phons = append(exceptC_Phons, phons[currPh-1])
			}
		}
		r_c = g.new_R_Cx(exceptC_Phons...)
		// r_v = g.new_R_Vx(exceptV_Phons...)
		r_v = g.new_R_noSlide_X(exceptV_Phons...)
	} else {
		exceptC_Phons = append(exceptC_Phons, phons[currPh])
		if currPh+1 < len(phons)-1 && !isVowel(phons[currPh+1]) {
			exceptC_Phons = append(exceptC_Phons, phons[currPh+1])
		}
		if currPh > 0 && !isVowel(phons[currPh-1]) && phons[currPh-1] != sil {
			exceptC_Phons = append(exceptC_Phons, phons[currPh-1])
		}
		r_x = g.new_R_Cx(exceptC_Phons...)
	}
	var retRule rule
	if usePen != noPen {
		if currPhIsVowel {
			vGuard := new_R_or(new_R_opt(new_R_and(r_v, usePen)), new_R_opt(new_R_and(usePen, r_v)))
			cGuard := new_R_or(new_R_opt(new_R_and(r_c, usePen)), new_R_opt(new_R_and(usePen, r_c)))
			retRule = new_R_and(vGuard, cGuard)
		} else {
			retRule = new_R_or(new_R_opt(new_R_and(r_x, usePen)), new_R_opt(new_R_and(usePen, r_x)))
		}
	} else {
		// retRule = new_R_or(new_R_opt(new_R_and(r_x)), new_R_opt(new_R_and(r_x)))
		// Get a load of _debug so we can figure out why we can't get hold of a
		// guard.
		// What'll be useful is:
		// currPh
		// previous and next phonemes
		// last gaurd used
		// grammar?
		str := fmt.Sprintln("preRule: no guard!")
		str += fmt.Sprintln("phons[currPh] =", phons[currPh])
		if currPh > 0 {
			str += fmt.Sprintln("phons[currPh - 1] =", phons[currPh-1])
		}
		if currPh < len(phons)-1 {
			str += fmt.Sprintln("phons[currPh + 1] =", phons[currPh+1])
		}
		lastGuards := g.config.getPens()
		str += fmt.Sprintln("lastGuards =", fmt.Sprint(lastGuards))
		str += fmt.Sprintln("grammar =", g)
		_debug(str)
		log.Panic()
	}

	return retRule
}

/*
func (g *jsgfGrammar) preRule(phons []phoneme, currPh int, usePen phonRule) rule {
  if currPh < 0 || currPh > len(phons) - 1 {
    return nullRule{}
  }
  currPhIsVowel := isVowel(phons[currPh])
  var exceptPhons = []phoneme{}
  if usePen != noPen {
    exceptPhons = append(exceptPhons, usePen.ph)
  }
  if g.config.getLastPen() != sil {
    exceptPhons = append(exceptPhons, g.config.getLastPen())
  }
  var r_x, r_v, r_c rule
  if currPhIsVowel {
    exceptPhons_v := append(exceptPhons, g.config.vowelHandler(phons[currPh])...)
    exceptPhons_c := exceptPhons
    if currPh > 0 && isVowel(phons[currPh - 1]) {
      // exceptPhons_v = append(exceptPhons, g.config.vowelHandler(phons[currPh])...)
      exceptPhons_v = append(exceptPhons_v, g.config.vowelHandler(phons[currPh - 1])...)
    } else {
      exceptPhons_c := append(exceptPhons, phons[currPh])
      if currPh > 0 && !isVowel(phons[currPh - 1]) && phons[currPh - 1] != sil {
        exceptPhons_c = append(exceptPhons_c, phons[currPh - 1])
      }
    }
    // Put rc = , rv = here
    r_c = g.new_R_Cx(exceptPhons_c...)
    r_v = g.new_R_Vx(exceptPhons_v...)
  } else {
    exceptPhons = append(exceptPhons, phons[currPh])
    if currPh > 0 && !isVowel(phons[currPh - 1]) && phons[currPh - 1] != sil {
      exceptPhons = append(exceptPhons, phons[currPh - 1])
    }
    r_x = g.new_R_Cx(exceptPhons...)
  }
  var retRule rule
  if usePen != noPen {
    if currPhIsVowel {
      vGuard := new_R_or(new_R_opt(new_R_and(r_v, usePen)), new_R_opt(new_R_and(usePen, r_v)))
      cGuard := new_R_or(new_R_opt(new_R_and(r_c, usePen)), new_R_opt(new_R_and(usePen, r_c)))
      retRule = new_R_and(vGuard, cGuard)
    } else {
      retRule = new_R_or(new_R_opt(new_R_and(r_x, usePen)), new_R_opt(new_R_and(usePen, r_x)))
    }
  } else {
    retRule = new_R_or(new_R_opt(new_R_and(r_x)), new_R_opt(new_R_and(r_x)))
    log.Panic()
  }

  return retRule
}
*/

func (g *jsgfGrammar) postRule(phons []phoneme, currPh int, usePen phonRule) rule {
	if currPh < 0 || currPh > len(phons)-1 {
		return nullRule{}
	}
	currPhIsVowel := isVowel(phons[currPh])
	var exceptPhons = []phoneme{}
	if usePen != noPen {
		exceptPhons = append(exceptPhons, usePen.ph)
	}
	/*
	  if g.config.getLastPen() != sil {
	    exceptPhons = append(exceptPhons, g.config.getLastPen())
	  }
	*/
	exceptPhons = append(exceptPhons, g.config.getPens()...)
	var r_x rule
	if currPhIsVowel {
		exceptPhons = append(exceptPhons, g.config.vowelHandler(phons[currPh])...)
		if currPh < len(phons)-1 && isVowel(phons[currPh+1]) {
			exceptPhons = append(exceptPhons, phons[currPh+1])
		}
		if currPh-1 > 0 && isVowel(phons[currPh-1]) {
			exceptPhons = append(exceptPhons, g.config.vowelHandler(phons[currPh-1])...)
		}
		// r_x = g.new_R_Vx(exceptPhons...)
		r_x = g.new_R_noSlide_X(exceptPhons...)
	} else {
		exceptPhons = append(exceptPhons, phons[currPh])
		if currPh < len(phons)-1 && !isVowel(phons[currPh+1]) && phons[currPh+1] != sil {
			exceptPhons = append(exceptPhons, phons[currPh+1])
		}
		if currPh-1 > 0 && !isVowel(phons[currPh-1]) {
			exceptPhons = append(exceptPhons, phons[currPh-1])
		}
		r_x = g.new_R_Cx(exceptPhons...)
	}
	var retRule rule
	if usePen != noPen {
		retRule = new_R_or(new_R_and(r_x, usePen), new_R_and(usePen, r_x))
	} else {
		retRule = new_R_or(new_R_and(r_x), new_R_and(r_x))
	}

	return retRule
}

func (g *jsgfGrammar) middleRule(phons []phoneme, currPh int) parsableRule {
	if currPh < 0 || currPh > len(phons)-1 {
		return nullRule{}
	}
	r_pen, ok := g.config.guardFor(phons, currPh)
	if !ok {
		// Generate an Error
		err := noGuardsError{}
		err.context = "middle rule..."
		err.this = phons[currPh]
		if currPh > 0 {
			err.previous = phons[currPh-1]
		}
		if currPh < len(phons)-1 {
			err.next = phons[currPh+1]
		}
		err.lastPenalties = g.config.getPens()
		err.grammar = *g
		_debug(err.Error())
		log.Panic()
	}
	r_ph := g.config.phonemeHandler(phons[currPh])
	retRule := new_R_and(g.preRule(phons, currPh, r_pen), new_R_or(r_ph, g.postRule(phons, currPh, r_pen)))
	// Now set last penalty here...
	g.config.setPen(r_pen.ph)
	var mRule parsableRule
	config := g.config
	if _, ok := config.(*jsgfDiphthong); ok {
		phs, ok := diphthongs[phons[currPh]]
		if ok {
			mRule = R_diphthongPhoneme{
				phs[0], phs[1],
				r_pen.ph,
				retRule,
			}
		} else {
			mRule = R_phoneme{
				phons[currPh],
				r_pen.ph,
				[]phoneme{},
				retRule,
			}
		}
	} else {
		mRule = R_phoneme{
			phons[currPh],
			r_pen.ph,
			[]phoneme{},
			retRule,
		}
	}
	return mRule
}

//       _ _                     _              _         _ _
//  __ _| | |_ ___ _ _ _ _  __ _| |_ ___     __| |_  __ _| | |___ _ _  __ _ ___
// / _` | |  _/ -_) '_| ' \/ _` |  _/ -_)   / _| ' \/ _` | | / -_) ' \/ _` / -_)
// \__,_|_|\__\___|_| |_||_\__,_|\__\___|   \__|_||_\__,_|_|_\___|_||_\__, \___|
//                                                                    |___/
// alternate, alternative, challenge
//
func (gr jsgfGrammar) otherPhons(phons []phoneme, currPh int) []phoneme {

	// ideas
	// 1) Pass in phoneme location so that 1st and last can be separated
	// 2) Work in any position .... for challenging vowels

	others := []phoneme{}

	if currPh < 0 || currPh > len(phons)-1 {
		return others
	}

	phon := phons[currPh]

	if currPh == 0 {
		switch phon {
		case k:
			others = append(others, d)
		//others = append(others, p, t)    //w
		//others = append(others, p, g, t, f)    //w
		//others = append(others, f, jh, ch)    //w
		//case b:
		//others = append(others, p, v, f)      // (r, z, hh)    .. tried y, aa, m, t, k    ... removed dh & th as they upset bars
		case f:
			others = append(others, s, p) // Note, can not do the "s" vs "f" the other way around because "f" will sit inside (is a shorter) "s"
		case w:
			others = append(others, v)
		case r:
			others = append(others, l, w)
		case s:
			others = append(others, z)
		case th:
			others = append(others, t, d, s) // was n prior to 9Apr22  PE
		case dh:
			others = append(others, t)
		default:
		}
		return others
	}

	/*
	  if currPh == len(phons) -1 {
	    switch phon {
	    case k:
	      //others = append(others, p, t, g, hh)
	      others = append(others, p, t)
	    case p:
	      others = append(others, t, k)
	      //others = append(others, m, jh)
	    case t:
	      others = append(others, p, k)
	      //others = append(others, m, jh)
	    case g:
	      //others = append(others, b, d, ch)
	      others = append(others, b, d, m)
	    case b:
	      others = append(others, g, d, ch)
	    case d:
	      //others = append(others, g, b, ch)
	      others = append(others, jh)
	    case ch:
	      others = append(others, g, b, d)
	    case uw:
	      others = append(others, jh, m)

	    default:
	    }
	    return others
	  }
	*/

	/*
	   switch phon {
	     case iy:
	       others = append(others, ih)
	     case ae:
	       others = append(others, eh)
	     default:
	     }
	*/

	return others
}

func (g *jsgfGrammar) endMiddleRule(phons []phoneme, currPh int) parsableRule {
	if currPh < 0 || currPh > len(phons)-1 {
		return nullRule{}
	}

	r_pen, ok := g.config.guardFor(phons, currPh)

	if !ok {
		// Generate an Error
		err := noGuardsError{}
		err.context = "middle rule..."
		err.this = phons[currPh]
		if currPh > 0 {
			err.previous = phons[currPh-1]
		}
		if currPh < len(phons)-1 {
			err.next = phons[currPh+1]
		}
		err.lastPenalties = g.config.getPens()
		err.grammar = *g
		_debug(err.Error())
		log.Panic()
	}
	r_ph := g.config.phonemeHandler(phons[currPh])
	rules := []rule{
		r_ph,
	}
	others := g.otherPhons(phons, currPh)

	//others := g.otherPhons(phons[currPh])
	for _, other := range others {
		otherRule := phonRule{
			other,
		}
		rules = append(rules, otherRule)
	}
	rules = append(rules, g.postRule(phons, currPh, r_pen))

	retRule := new_R_and(g.preRule(phons, currPh, r_pen), new_R_or(rules...))

	// Now set last penalty here...
	g.config.setPen(r_pen.ph)
	var mRule parsableRule
	config := g.config
	if _, ok := config.(*jsgfDiphthong); ok {
		phs, ok := diphthongs[phons[currPh]]
		if ok {
			mRule = R_diphthongPhoneme{
				phs[0], phs[1],
				r_pen.ph,
				retRule,
			}
		} else {
			mRule = R_phoneme{
				phons[currPh],
				r_pen.ph,
				others,
				retRule,
			}
		}
	} else {
		mRule = R_phoneme{
			phons[currPh],
			r_pen.ph,
			others,
			retRule,
		}
	}
	return mRule
}

//===================================================================
//    ___ _        _                        _
//   / __| |___ __(_)_ _  __ _     _ _ _  _| |___
//  | (__| / _ (_-< | ' \/ _` |   | '_| || | / -_)
//   \___|_\___/__/_|_||_\__, |   |_|  \_,_|_\___|
//                        |___/
//===================================================================

type R_closing struct {
	guard phoneme
	rule
}

func (r R_closing) generate() string {
	return r.rule.generate()
}

func (r R_closing) parse(res []psPhonemeDatum, i int) ([]parseResult, error) {
	if i < 0 || i > len(res) {
		return []parseResult{}, parseFailed
	}
	// If the final phoneme isn't a sil something isn't right so just return
	if res[len(res)-1].phoneme != "sil" {
		return []parseResult{}, parseFailed
	}
	// We might just have a sil
	if i == len(res)-1 {
		parsRes := parseResult{
			i, i, false,
		}
		return []parseResult{
			parsRes,
		}, nil
	}
	// Let's look for guards
	indexes := allIndexes(res[i:min(i+4, len(res))], r.guard)
	if len(res) == i+len(indexes)*2+1 {
		parsRes := parseResult{
			i, len(res) - 1, false,
		}
		return []parseResult{
			parsRes,
		}, nil
	} else {
		return []parseResult{}, parseFailed
	}
	/*
	  if len(indexes) == 2 {
	    if len(res) > i + 4 {
	      // We found 2 guards
	      parsRes := parseResult{
	        i, i + 4, false,
	      }
	      // There might also by a soft fade so let's check for that now.
	      if i + 5 < len(res) - 1 {
	        if res[i + 4].phoneme == p && res[i + 5].phoneme == f {
	          parsRes.end = i + 6
	        }
	      }
	      ret = append(ret, parsRes)
	      return ret, nil
	    } else {
	      return []parseResult{}, parseFailed
	    }
	  }
	  if len(indexes) == 1 {
	    if len(res) > i + 2 {
	      // We found 1 guard
	      parsRes := parseResult{
	        i, i + 2, false,
	      }
	      // There might also by a soft fade so let's check for that now.
	      if i + 3 < len(res) - 1 {
	        if res[i + 2].phoneme == p && res[i + 3].phoneme == f {
	          parsRes.end = i + 4
	        }
	      }
	      ret = append(ret, parsRes)
	      return ret, nil
	    } else {
	      return []parseResult{}, parseFailed
	    }
	  }
	  if len(indexes) == 0 {
	    // There might by a soft fade so let's check for that now.
	    if i + 1 < len(res) - 1 {
	      if res[i].phoneme == p && res[i + 1].phoneme == f {
	        parsRes := parseResult{
	          i, i + 2, false,
	        }
	        ret = append(ret, parsRes)
	        return ret, nil
	      }
	    }
	  }
	*/
	// It's possible we didn't find any guards in which case we should probably
	// return an error here because we already check for just a sil...
	// return ret, parseFailed
}

func (gr *jsgfGrammar) closingRule(phons []phoneme) parsableRule {
	r_pen, ok := gr.config.closingPenalty(phons[len(phons)-1])
	if !ok {
		// Generate an error
		err := noGuardsError{}
		err.context = "closing rule..."
		err.previous = phons[len(phons)-1]
		err.lastPenalties = gr.config.getPens()
		err.grammar = *gr
		_debug(err.Error())
		log.Panic()
	}
	r_any_v := gr.new_R_V_noSlide()
	//r_fade := gr.new_R_softFade()            // soft fade

	// Configure the closing rule phoneme based on the last phoneme in the word
	lastPh := phons[len(phons)-1]

	xPhons := gr.config.getPens()
	xPhons = append(xPhons, lastPh, r_pen.ph)
	r_Cx := gr.new_R_Cx(xPhons...)
	r_Vx := gr.new_R_Vx(xPhons...)
	var retRule rule
	if isVowel(lastPh) {
		retRule = new_R_and(new_R_or(new_R_opt(new_R_and(r_Cx, r_pen)), new_R_opt(new_R_and(r_pen, r_Cx))), new_R_or(new_R_opt(new_R_and(r_Vx, r_pen)), new_R_opt(new_R_and(r_pen, r_Vx))), new_R_ph(sil))
	} else {
		retRule = new_R_and(new_R_or(new_R_opt(new_R_and(r_any_v, r_pen)), new_R_opt(new_R_and(r_pen, r_any_v))), new_R_or(new_R_opt(new_R_and(r_Cx, r_pen)), new_R_opt(new_R_and(r_pen, r_Cx))), new_R_ph(sil))

	}
	gr.config.setPen(r_pen.ph)

	cRule := R_closing{
		r_pen.ph,
		retRule,
	}
	return cRule
}

type R_target struct {
	word  string
	rules []parsableRule
}

func (g *jsgfGrammar) new_R_target(word string, phons []phoneme) {
	rules := []parsableRule{}
	rules = append(rules, g.trappedOpeningRule(phons))

	for i := range phons {
		if i == 0 || i == len(phons)-1 {
			rules = append(rules, g.endMiddleRule(phons, i))

		} else {
			rules = append(rules, g.middleRule(phons, i))

		}
	}
	rules = append(rules, g.closingRule(phons))

	g.target = R_target{
		word,
		rules,
	}
}

func (t R_target) generate() string {
	str := "public <" + t.word + "> = ("
	strs := []string{}
	for _, r := range t.rules {
		strs = append(strs, r.generate())
	}
	str += strings.Join(strs, " ") + ")"
	return str
}

type jsgfConfig interface {
	// setLastPen(ph phoneme)
	setPen(ph phoneme)
	addPen(ph phoneme)
	getPens() []phoneme
	// getLastPen() phoneme
	// penalty([]phoneme, int) (phonRule, bool)
	openingPenalty(ph phoneme) (phonRule, bool)
	guardFor([]phoneme, int) (phonRule, bool)
	closingPenalty(ph phoneme) (phonRule, bool)
	phonemeHandler(phoneme) rule
	vowelHandler(phoneme) []phoneme
}

type jsgfStandard struct {
	lastPens []phoneme
}

func new_jsgfStandard() jsgfStandard {
	return jsgfStandard{
		[]phoneme{},
	}
}

var noPen = phonRule{
	sil, // indicating no phoneme penalty
}

/*
func (j jsgfStandard) getLastPen() phoneme {
  return j.lastPen
}

func (j *jsgfStandard) setLastPen(ph phoneme) {
  j.lastPen = ph
}
*/

func (j *jsgfStandard) setPen(ph phoneme) {
	j.lastPens = []phoneme{
		ph,
	}
}

func (j *jsgfStandard) addPen(ph phoneme) {
	j.lastPens = append(j.lastPens, ph)
}

func (j *jsgfStandard) getPens() []phoneme {
	if len(j.lastPens) == 0 {
		_debug("No last pens to return!")
	}
	return j.lastPens
}

func (j *jsgfStandard) openingPenalty(ph phoneme) (phonRule, bool) {
	// Don't want the following with a vowel: m, v?
	///penalties := consRemove(newConsonants(), m, v)
	// Or voiced soft phonemes, not sure what to include here so I;ve put a
	// couple of example phonemes in
	//penalties = consRemove(penalties, b, p)

	//penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, nd, ng, p, r, s, sh, tS, t, th, v, w, z, zh, dz, ts, axm, axl, axn, kw, ks)  // remove everything

	//penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing)  // remove everything

	penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing, thr, sts, st, kt, fl) // remove everything

	//64only
	//penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing, fl)  // remove everything
	//partner := []phoneme{}

	switch ph {

	case p:
		//partner = append(partner, b)
		penalties[b] = true
		penalties[hh] = true
		penalties[sh] = true
	case b:
		//partner = append(partner, p, dh, sh, zh, hh, ch, jh, m)     //.. p and anything with a "h" type sound
		penalties[p] = true
		penalties[hh] = true
		penalties[sh] = true

		penalties[d] = true
		penalties[l] = true
		penalties[f] = true
		penalties[v] = true
		penalties[dh] = true
	case f:
		penalties[th] = true
		penalties[s] = true
		penalties[v] = true
		penalties[z] = true
		penalties[dh] = true
	case hh:
		penalties[z] = true
		penalties[b] = true
		penalties[iy] = true
	case g:
		penalties[jh] = true
		penalties[d] = true
		//penalties[b] = true
		penalties[z] = true
		penalties[v] = true
		penalties[ch] = true
		penalties[k] = true
	case k:
		penalties[t] = true
		penalties[d] = true
		penalties[ch] = true
		penalties[g] = true
		penalties[p] = true
	case l:
		penalties[r] = true
		penalties[w] = true
		penalties[m] = true
		penalties[v] = true
		penalties[z] = true
	case m:
		//penalties[b] = true
		penalties[d] = true
		penalties[l] = true
		penalties[v] = true
		penalties[z] = true
	case r:
		//penalties[w] = true
		penalties[n] = true
		//penalties[l] = true
		penalties[hh] = true
		penalties[m] = true
	case s:
		// do not use hh - it can mess up the "ae" in eh-school
		penalties[z] = true
		penalties[t] = true
		penalties[sh] = true
		penalties[ch] = true
		penalties[th] = true
		//penalties[f] = true    // f is a fragment of S and will nearly always be found in S
	case t:
		penalties[d] = true
		penalties[f] = true
		penalties[s] = true
		penalties[ch] = true
		penalties[k] = true
		penalties[g] = true
	case v:
		penalties[z] = true
		penalties[m] = true
		penalties[n] = true
		penalties[l] = true
		//penalties[d] = true
		penalties[r] = true
	case w:
		penalties[z] = true
		penalties[m] = true
		penalties[n] = true
		penalties[l] = true
		penalties[d] = true
		penalties[r] = true
	case z:
		//partner = append(partner, s)
		penalties[s] = true
		penalties[v] = true
		penalties[b] = true
	case th:
		//partner = append(partner, f, t, z, dh, s)
		//penalties[f] = true    // f should not guard th because the system doesn't seem to be able to tell the difference
		penalties[t] = true
		penalties[z] = true
		penalties[dh] = true
		penalties[s] = true
		//case y    ... see vowel section

	case aa:
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	case ah:
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	case er:
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	case eh:
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	case ehr: // Added 8 March 2021 - simply copied eh rule to get strated.
		penalties[s] = true // was b, but always seems to be found.
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	case ow:
		/*
		   penalties[uw] = true
		   penalties[ao] = true
		   penalties[l] = true
		*/
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true

	case y:
		penalties[b] = true
		penalties[hh] = true
		penalties[r] = true
		penalties[l] = true
	default:
		penalties[z] = true
		//partner = append(partner, hh)
	}

	pen := pickConsonant(penalties)
	if ph != sil {
		return phonRule{
			pen,
		}, true
	} else {
		return phonRule{}, false
	}
}

var confusionMatrix = map[phoneme]map[phoneme]float64{

	p: {
		p: 86, t: 13.8, k: 18.2, f: 9.6, th: 8.8, s: 1.5, sh: 3, ch: 3.0, hh: 19.2, b: 21.3, d: 5.8, g: 3.1, v: 2.9, dh: 3.3, z: 3, zh: 0.5, jh: 0.5, y: 3, m: 1.4, n: 3, ng: 0.5, l: 1.3, r: 1, w: 1.3,
	},
	t: {
		p: 24.6, t: 95, k: 12.5, f: 7.7, th: 9.2, s: 4, sh: 2.7, ch: 5, hh: 11.3, b: 7.1, d: 20.4, g: 5, v: 2, dh: 5.4, z: 3.7, zh: 0.5, jh: 5, y: 1.5, m: 2.1, n: 3.3, ng: 2, l: 2.4, r: 0.8, w: 1, oy: 0.5, ao: 0.5, aa: 0.5,
	},
	k: {
		p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
	},
	f: {
		p: 24.6, t: 21.7, k: 13.9, f: 90, th: 11.3, s: 7.5, sh: 2.4, ch: 2.1, hh: 9.2, b: 15, d: 10.4, g: 5.5, v: 5.4, dh: 5, z: 2.4, zh: 1, jh: 2, y: 0.4, m: 1.7, n: 2, ng: 1.2, l: 0.8, r: 0.4, w: 0.8,
	},
	th: {
		p: 18.8, t: 24.6, k: 5.6, f: 39.3, th: 74, s: 12, sh: 3.7, ch: 2.4, hh: 7.1, b: 14.2, d: 7.9, g: 0.5, v: 4, dh: 9.2, z: 5.5, zh: 0.5, jh: 0.5, y: 2, m: 1.3, n: 2.9, ng: 0.5, l: 2.9, r: 0.5, w: 0.5,
	},
	s: {
		p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao: 1, oy: 0.1,
	},
	sh: {
		p: 5.9, t: 9.7, k: 6.7, f: 4.1, th: 6, s: 11.3, sh: 87, ch: 26, hh: 2.2, b: 0.5, d: 0.5, g: 0.5, v: 2.2, dh: 1.6, z: 15, zh: 10, jh: 16.1, y: 0.8, m: 1.7, n: 0.5, ng: 0.5, l: 0.7, r: 0.4, w: 0,
	},
	ch: {
		p: 7.7, t: 16.5, k: 5.5, f: 5.4, th: 3.7, s: 8.9, sh: 13, ch: 76, hh: 4, b: 4.6, d: 3.3, g: 3, v: 1.6, dh: 3.3, z: 1.4, zh: 5.4, jh: 34.2, y: 0.4, m: 1.5, n: 2.5, ng: 1.3, l: 0, r: 0, w: 0,
	},
	hh: {
		p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 2, sh: 2, ch: 2, hh: 99, b: 0.5, d: 1.3, g: 0.4, v: 4.6, dh: 2, z: 0.4, zh: 2, jh: 2, y: 2, m: 2, n: 2, ng: 2, l: 1.7, r: 2, w: 2,
	},
	b: {
		p: 8.1, t: 7.5, k: 6.7, f: 13.7, th: 7.4, s: 0.5, sh: 0.5, ch: 1.1, hh: 0.5, b: 91, d: 15, g: 14.6, v: 18.6, dh: 6.1, z: 2, zh: 0.5, jh: 0.5, y: 2.9, m: 7.6, n: 3.8, ng: 1.3, l: 5, r: 1.3, w: 5,
	},
	d: {
		p: 2.5, t: 16.3, k: 3, f: 3.1, th: 5.8, s: 2.7, sh: 1.5, ch: 3, hh: 8.8, b: 12.1, d: 92, g: 8, v: 5.5, dh: 12.9, z: 4.5, zh: 2.9, jh: 5.8, y: 6.3, m: 4.2, n: 12.5, ng: 2.5, l: 12.5, r: 2.1, w: 0.5,
	},
	g: {
		p: 3.3, t: 12.9, k: 9.2, f: 0.5, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 9.2, b: 10.1, d: 20.4, g: 100, v: 7.8, dh: 5, z: 0.5, zh: 0.5, jh: 9.5, y: 24.2, m: 2.5, n: 3.5, ng: 2, l: 3.6, r: 1.5, w: 1.7, oy: 0.3,
	},
	v: {
		p: 7.5, t: 12.9, k: 3.8, f: 12.1, th: 5.4, s: 2.6, sh: 2.1, ch: 0.6, hh: 6.7, b: 30, d: 15.8, g: 7, v: 91, dh: 13, z: 14.7, zh: 0.5, jh: 2.7, y: 3.0, m: 3.8, n: 5.1, ng: 2.7, l: 4.2, r: 4, w: 5.4, oy: 0.5,
	},
	dh: {
		p: 2.5, t: 11.3, k: 2.5, f: 4.2, th: 14.6, s: 3.3, sh: 1.7, ch: 0.5, hh: 2.1, b: 17.1, d: 29.2, g: 6.3, v: 31.4, dh: 79, z: 12.3, zh: 0.5, jh: 7.9, y: 1.3, m: 1.8, n: 6.2, ng: 2, l: 12.1, r: 2.1, w: 2.9,
	},
	z: {
		p: 1.6, t: 3.8, k: 0.5, f: 2.6, th: 10, s: 15, sh: 15, ch: 2.5, hh: 0.5, b: 7.9, d: 8.8, g: 1, v: 13.9, dh: 23.8, z: 98, zh: 8.3, jh: 4, y: 2.5, m: 2.6, n: 8.1, ng: 0.5, l: 2.6, r: 0.5, w: 3.8, oy: 0.5,
	},
	zh: {
		p: 1.6, t: 0.5, k: 1.4, f: 1.7, th: 2.9, s: 2.5, sh: 14.2, ch: 4.2, hh: 1.9, b: 0.5, d: 3.9, g: 22, v: 8.9, dh: 1, z: 13, zh: 45, jh: 45, y: 0, m: 1.3, n: 0.3, ng: 1.7, l: 0.8, r: 0.8, w: 10,
	},
	jh: {
		p: 5, t: 8.2, k: 5.5, f: 4.5, th: 4.8, s: 6.2, sh: 5.2, ch: 18.3, hh: 2.6, b: 0.5, d: 17.7, g: 23, v: 2, dh: 7.5, z: 7.8, zh: 17, jh: 59.8, y: 5.8, m: 1.7, n: 0.3, ng: 1.6, l: 3.1, r: 1.4, w: 0,
	},
	y: {
		p: 1.3, t: 0, k: 0.8, f: 0.9, th: 0.8, s: 1.4, sh: 3, ch: 0.9, hh: 2.1, b: 4.2, d: 2.5, g: 1.3, v: 0.7, dh: 1.3, z: 1.6, zh: 0, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 0, l: 8.5, r: 2.6, w: 2.4,
	},
	m: {
		p: 3.8, t: 9.2, k: 2.5, f: 2.9, th: 4.6, s: 1.7, sh: 0.8, ch: 0.7, hh: 3.9, b: 9.6, d: 5.8, g: 1.5, v: 7.3, dh: 2.1, z: 1.5, zh: 1.2, jh: 2.4, y: 0.9, m: 100, n: 20, ng: 10.8, l: 11.7, r: 5.4, w: 4.2, ing: 99,
	},
	n: {
		p: 1.6, t: 9.6, k: 0.6, f: 1, th: 1.6, s: 0.9, sh: 0.3, ch: 0.3, hh: 3.2, b: 1.7, d: 8.3, g: 1.9, v: 4.5, dh: 2.1, z: 3.8, zh: 0.3, jh: 0.3, y: 1.2, m: 24.1, n: 94, ng: 14.1, l: 17.5, r: 4, w: 1.8, ing: 99,
	},
	ng: {
		p: 0, t: 6.7, k: 0.4, f: 2.1, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 0, b: 2.5, d: 5.8, g: 6.3, v: 1.9, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0, m: 13.8, n: 35, ng: 60, l: 2.1, r: 1.7, w: 10, ing: 99,
	},
	l: {
		p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, uwl: 99, axl: 99,
	},
	r: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
	},
	w: {
		p: 1.7, t: 0, k: 0.4, f: 1.2, th: 0.4, s: 1.2, sh: 0, ch: 0.9, hh: 1.7, b: 5.8, d: 0.8, g: 0, v: 2.1, dh: 0.4, z: 3.9, zh: 1, jh: 3.2, y: 6.5, m: 5.8, n: 0.8, ng: 10, l: 2.5, r: 4.2, w: 99,
	},
	iy: {
		p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 1.4, sh: 0.5, ch: 1.5, hh: 0.5, b: 8.3, d: 2.5, g: 3, v: 4.6, dh: 1.7, z: 1.6, zh: 3, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 3, l: 8.5, r: 2.6, w: 3, ih: 99, aa: 0.5, ao: 0.5, aw: 0.5,
	},
	ih: {
		p: 26.3, t: 4.6, k: 12.1, f: 0.5, th: 5, s: 0.4, sh: 0.4, ch: 3, hh: 99, b: 8.3, d: 1.3, g: 2, v: 4.6, dh: 1.7, z: 0.4, zh: 2, jh: 0.5, y: 3, m: 2, n: 3, ng: 2, l: 1.7, r: 1.5, w: 2,
	},
	ey: {
		p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 1.4, sh: 0.8, ch: 0.9, hh: 99, b: 8.3, d: 2.5, g: 1.3, v: 4.6, dh: 1.7, z: 1.6, zh: 3, jh: 0.5, y: 2.2, m: 10, n: 4.2, ng: 2.6, l: 8.5, r: 2.6, w: 1, ao: 0.3, ehr: 99, eh: 99, er: 99,
	},
	eh: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 2.5, zh: 0.5, jh: 0.5, y: 0.5, m: 4.1, n: 4.1, ng: 1, l: 14.8, r: 96, w: 1, oy: 0.5, ax: 99, axr: 99, ehr: 99, ey: 99, er: 99,
	},
	ae: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 3.8, s: 0.8, sh: 0.4, ch: 0.9, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 2, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 0, l: 14.8, r: 96, w: 6.7, aa: 3, ao: 0.5,
	},
	aa: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 3, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.9, ng: 0.5, l: 14.8, r: 96, w: 6.7, ax: 99, axr: 99,
	},
	ah: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ax: 99, axr: 99,
	},
	ao: {
		p: 26.3, t: 0.5, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 1, l: 14.8, r: 96, w: 6.7, oy: 3, ih: 0.3, ax: 99, axl: 99,
	},
	ow: {
		p: 26.3, t: 7.9, k: 12.1, f: 1, th: 5, s: 0.5, sh: 1.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 10, jh: 3.2, y: 0.5, m: 5.8, n: 2.1, ng: 10, l: 14.8, r: 96, w: 99,
	},
	uh: {
		p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 0.4, sh: 0.4, ch: 0.8, hh: 99, b: 8.3, d: 1.3, g: 0.4, v: 4.6, dh: 1.7, z: 0.4, zh: 0, jh: 1, y: 0.8, m: 0, n: 0.8, ng: 0, l: 1.7, r: 0.8, w: 0.4, ax: 99, axr: 99, uwn: 99, uwl: 99,
	},
	uw: {
		p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axr: 99,
	},
	ay: {
		p: 0.5, t: 7.9, k: 12.1, f: 11.3, th: 3.8, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 3, m: 10, n: 4.2, ng: 10, l: 14.8, r: 96, w: 6.7,
	},
	oy: {
		p: 1.3, t: 0, k: 0.8, f: 0.9, th: 0.8, s: 1.4, sh: 0.8, ch: 0.9, hh: 2.1, b: 4.2, d: 2.5, g: 0.1, v: 0.7, dh: 1.3, z: 1.6, zh: 0, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 0, l: 8.5, r: 2.6, w: 2.4,
	},
	aw: {
		p: 2.5, t: 7.9, k: 1.7, f: 0.5, th: 0.5, s: 0.8, sh: 0.4, ch: 0.5, hh: 5.8, b: 0.5, d: 0.5, g: 4.9, v: 0.5, dh: 0.8, z: 0.5, zh: 0.4, jh: 0.5, y: 4.6, m: 3.6, n: 2.1, ng: 10, l: 14.8, r: 96, w: 6.7, oy: 0.5, ax: 99, axr: 99,
	},
	er: {
		p: 26.3, t: 7.9, k: 12.1, f: 1, th: 5, s: 0.5, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 1, dh: 1.7, z: 2.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 1, ax: 99, axr: 99,
	},
	ax: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ih: 0.3, aa: 99, ah: 99, oy: 3, axn: 99, axm: 99, axl: 99,
	},
	oh: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 3, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.9, ng: 0.5, l: 14.8, r: 96, w: 6.7, aa: 99, ah: 99, ax: 99, axr: 99,
	},
	ehr: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 2.5, zh: 0.5, jh: 0.5, y: 0.5, m: 4.1, n: 4.1, ng: 1, l: 14.8, r: 96, w: 1, oy: 0.5, ax: 99, axr: 99, er: 99, eh: 99, ey: 99,
	},
	axl: {
		p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99, axr: 99,
	},
	axm: {
		p: 3.8, t: 9.2, k: 2.5, f: 2.9, th: 4.6, s: 1.7, sh: 0.8, ch: 0.7, hh: 3.9, b: 9.6, d: 5.8, g: 1.5, v: 7.3, dh: 2.1, z: 1.5, zh: 1.2, jh: 2.4, y: 0.9, m: 100, n: 100, ng: 100, l: 11.7, r: 5.4, w: 4.2, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, axr: 99,
	},
	axn: {
		p: 1.6, t: 9.6, k: 0.6, f: 1, th: 1.6, s: 0.9, sh: 0.3, ch: 0.3, hh: 3.2, b: 1.7, d: 8.3, g: 1.9, v: 4.5, dh: 2.1, z: 3.8, zh: 0.3, jh: 0.3, y: 1.2, m: 100, n: 100, ng: 100, l: 17.5, r: 4, w: 1.8, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, axr: 99,
	},
	ks: {
		p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao: 1,
	},
	kw: {
		p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
	},
	dz: {
		p: 1.6, t: 3.8, k: 0.5, f: 2.6, th: 10, s: 15, sh: 15, ch: 2.5, hh: 0.5, b: 7.9, d: 8.8, g: 1, v: 13.9, dh: 23.8, z: 98, zh: 8.3, jh: 4, y: 2.5, m: 2.6, n: 8.1, ng: 0.5, l: 2.6, r: 0.5, w: 3.8, oy: 0.5,
	},
	//ts: {
	//	p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao:1,
	//},
	uwl: {
		p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99, axr: 99,
	},
	uwn: {
		p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
	},
	uwm: {
		p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
	},
	kl: {
		p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
	},
	bl: {
		p: 8.1, t: 7.5, k: 6.7, f: 13.7, th: 7.4, s: 0.5, sh: 0.5, ch: 1.1, hh: 0.5, b: 91, d: 15, g: 14.6, v: 18.6, dh: 6.1, z: 2, zh: 0.5, jh: 0.5, y: 2.9, m: 7.6, n: 3.8, ng: 1.3, l: 5, r: 1.3, w: 5,
	},
	pl: {
		p: 86, t: 13.8, k: 18.2, f: 9.6, th: 8.8, s: 1.5, sh: 3, ch: 3.0, hh: 19.2, b: 21.3, d: 5.8, g: 3.1, v: 2.9, dh: 3.3, z: 3, zh: 0.5, jh: 0.5, y: 3, m: 1.4, n: 3, ng: 0.5, l: 1.3, r: 1, w: 1.3,
	},
	pr: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
	},
	kr: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
	},
	tr: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
	},
	gr: {
		p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
	},
	ing: {
		p: 0, t: 6.7, k: 0.4, f: 2.1, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 0, b: 2.5, d: 5.8, g: 6.3, v: 1.9, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0, m: 13.8, n: 99, ng: 99, l: 2.1, r: 1.7, w: 10, ih: 99, axn: 99, axm: 99, uwn: 99, uwm: 99,
	},
	axr: {
		p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ih: 0.3, aa: 99, ah: 99, oy: 3, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
	},

	// Adding 14 Feb 2022 for the new compound phonemes
	thr: {
		p: 18.8, t: 24.6, k: 5.6, f: 39.3, th: 99, s: 12, sh: 2.5, ch: 2.5, hh: 7.1, b: 14.2, d: 7.9, g: 2.5, v: 10, dh: 10, z: 5.5, zh: 1.0, jh: 1.0, y: 2, m: 2.5, n: 2.5, ng: 0.5, l: 14.8, r: 99, w: 10.0,
	},
	sts: {
		p: 25, t: 99, k: 14.0, f: 20.0, th: 25.0, s: 99, sh: 15.0, ch: 7, hh: 4, b: 5, d: 20, g: 5, v: 2, dh: 5, z: 25, zh: 10, jh: 2.7, y: 1.5, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao: 1,
	},
	ts: {
		p: 25, t: 99, k: 14.0, f: 20.0, th: 25.0, s: 99, sh: 15.0, ch: 7, hh: 4, b: 5, d: 20, g: 5, v: 2, dh: 5, z: 25, zh: 10, jh: 2.7, y: 1.5, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao: 1,
	},
	st: {
		p: 25, t: 99, k: 14.0, f: 20.0, th: 25.0, s: 99, sh: 15.0, ch: 7, hh: 4, b: 5, d: 20, g: 5, v: 2, dh: 5, z: 25, zh: 10, jh: 2.7, y: 1.5, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao: 1,
	},
	kt: {
		p: 25, t: 99, k: 99, f: 7.7, th: 9.2, s: 1, sh: 1, ch: 20, hh: 11.3, b: 7.1, d: 20.4, g: 15, v: 0.5, dh: 5.4, z: 0.5, zh: 0.5, jh: 5, y: 1.5, m: 2.1, n: 3.3, ng: 3, l: 2.4, r: 0.8, w: 1, oy: 0.5, ao: 0.5, aa: 0.5, ih: 2,
	},
	fl: {
		p: 24.6, t: 21.7, k: 13.9, f: 99, th: 11.3, s: 7.5, sh: 2.4, ch: 1.5, hh: 9.2, b: 15, d: 10.4, g: 5.5, v: 5.4, dh: 5, z: 1.4, zh: 1, jh: 2, y: 2.0, m: 1.7, n: 2, ng: 1.2, l: 99, r: 15, w: 15, uwl: 99, axl: 99,
	},
	ihl: {
		p: 30, t: 8.3, k: 1.7, f: 5.8, th: 5, s: 1, sh: 0.5, ch: 1.0, hh: 30, b: 8.3, d: 1, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 1, jh: 1, y: 3, m: 10, n: 3, ng: 3.0, l: 99, r: 10, w: 10, uwl: 99, axl: 99, ih: 99,
	},
	yuw: {
		p: 1.3, t: 0.5, k: 0.8, f: 0.9, th: 0.8, s: 1.4, sh: 3, ch: 0.9, hh: 5, b: 5, d: 2.5, g: 3, v: 0.7, dh: 2.0, z: 0.5, zh: 0.5, jh: 4.6, y: 99, m: 10, n: 4.2, ng: 0, l: 10, r: 3, w: 10, uw: 99,
	},
	ehl: {
		p: 3.0, t: 3.0, k: 3.0, f: 8.8, th: 5.0, s: 2, sh: 0.5, ch: 1.0, hh: 99, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 4.0, z: 2.5, zh: 1.0, jh: 2, y: 1.0, m: 4, n: 1.1, ng: 1.0, l: 99, r: 4.2, w: 5.4, uwl: 99, axl: 99, oy: 0.5, ax: 99, axr: 99, eh: 99,
	},
}

// 64only
/*
p: {
	p: 86, t: 13.8, k: 18.2, f: 9.6, th: 8.8, s: 1.5, sh: 3, ch: 3.0, hh: 19.2, b: 21.3, d: 5.8, g: 3.1, v: 2.9, dh: 3.3, z: 3, zh: 0.5, jh: 0.5, y: 3, m: 1.4, n: 3, ng: 0.5, l: 1.3, r: 1, w: 1.3,
},
t: {
	p: 24.6, t: 95, k: 12.5, f: 7.7, th: 9.2, s: 4, sh: 2.7, ch: 5, hh: 11.3, b: 7.1, d: 20.4, g: 5, v: 2, dh: 5.4, z: 3.7, zh: 0.5, jh: 5, y: 1.5, m: 2.1, n: 3.3, ng: 2, l: 2.4, r: 0.8, w: 1, oy: 0.5, ao: 0.5, aa: 0.5,
},
k: {
	p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
},
f: {
	p: 24.6, t: 21.7, k: 13.9, f: 90, th: 11.3, s: 7.5, sh: 2.4, ch: 2.1, hh: 9.2, b: 15, d: 10.4, g: 5.5, v: 5.4, dh: 5, z: 2.4, zh: 1, jh: 2, y: 0.4, m: 1.7, n: 2, ng: 1.2, l: 0.8, r: 0.4, w: 0.8,
},
th: {
	p: 18.8, t: 24.6, k: 5.6, f: 39.3, th: 74, s: 12, sh: 3.7, ch: 2.4, hh: 7.1, b: 14.2, d: 7.9, g: 0.5, v: 4, dh: 9.2, z: 5.5, zh: 0.5, jh: 0.5, y: 2, m: 1.3, n: 2.9, ng: 0.5, l: 2.9, r: 0.5, w: 0.5,
},
s: {
	p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao:1, oy: 0.1,
},
sh: {
	p: 5.9, t: 9.7, k: 6.7, f: 4.1, th: 6, s: 11.3, sh: 87, ch: 26, hh: 2.2, b: 0.5, d: 0.5, g: 0.5, v: 2.2, dh: 1.6, z: 15, zh: 10, jh: 16.1, y: 0.8, m: 1.7, n: 0.5, ng: 0.5, l: 0.7, r: 0.4, w: 0,
},
ch: {
	p: 7.7, t: 16.5, k: 5.5, f: 5.4, th: 3.7, s: 8.9, sh: 13, ch: 76, hh: 4, b: 4.6, d: 3.3, g: 3, v: 1.6, dh: 3.3, z: 1.4, zh: 5.4, jh: 34.2, y: 0.4, m: 1.5, n: 2.5, ng: 1.3, l: 0, r: 0, w: 0,
},
hh: {
	p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 2, sh: 2, ch: 2, hh: 99, b: 0.5, d: 1.3, g: 0.4, v: 4.6, dh: 2, z: 0.4, zh: 2, jh: 2, y: 2, m: 2, n: 2, ng: 2, l: 1.7, r: 2, w: 2,
},
b: {
	p: 8.1, t: 7.5, k: 6.7, f: 13.7, th: 7.4, s: 0.5, sh: 0.5, ch: 1.1, hh: 0.5, b: 91, d: 15, g: 14.6, v: 18.6, dh: 6.1, z: 2, zh: 0.5, jh: 0.5, y: 2.9, m: 7.6, n: 3.8, ng: 1.3, l: 5, r: 1.3, w: 5,
},
d: {
	p: 2.5, t: 16.3, k: 3, f: 3.1, th: 5.8, s: 2.7, sh: 1.5, ch: 3, hh: 8.8, b: 12.1, d: 92, g: 8, v: 5.5, dh: 12.9, z: 4.5, zh: 2.9, jh: 5.8, y: 6.3, m: 4.2, n: 12.5, ng: 2.5, l: 12.5, r: 2.1, w: 0.5,
},
g: {
	p: 3.3, t: 12.9, k: 9.2, f: 0.5, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 9.2, b: 10.1, d: 20.4, g: 100, v: 7.8, dh: 5, z: 0.5, zh: 0.5, jh: 9.5, y: 24.2, m: 2.5, n: 3.5, ng: 2, l: 3.6, r: 1.5, w: 1.7, oy: 0.3,
},
v: {
	p: 7.5, t: 12.9, k: 3.8, f: 12.1, th: 5.4, s: 2.6, sh: 2.1, ch: 0.6, hh: 6.7, b: 30, d: 15.8, g: 7, v: 91, dh: 13, z: 14.7, zh: 0.5, jh: 2.7, y: 3.0, m: 3.8, n: 5.1, ng: 2.7, l: 4.2, r: 4, w: 5.4, oy:0.5,
},
dh: {
	p: 2.5, t: 11.3, k: 2.5, f: 4.2, th: 14.6, s: 3.3, sh: 1.7, ch: 0.5, hh: 2.1, b: 17.1, d: 29.2, g: 6.3, v: 31.4, dh: 79, z: 12.3, zh: 0.5, jh: 7.9, y: 1.3, m: 1.8, n: 6.2, ng: 2, l: 12.1, r: 2.1, w: 2.9,
},
z: {
	p: 1.6, t: 3.8, k: 0.5, f: 2.6, th: 10, s: 15, sh: 15, ch: 2.5, hh: 0.5, b: 7.9, d: 8.8, g: 1, v: 13.9, dh: 23.8, z: 98, zh: 8.3, jh: 4, y: 2.5, m: 2.6, n: 8.1, ng: 0.5, l: 2.6, r: 0.5, w: 3.8, oy:0.5,
},
zh: {
	p: 1.6, t: 0.5, k: 1.4, f: 1.7, th: 2.9, s: 2.5, sh: 14.2, ch: 4.2, hh: 1.9, b: 0.5, d: 3.9, g: 22, v: 8.9, dh: 1, z: 13, zh: 45, jh: 45, y: 0, m: 1.3, n: 0.3, ng: 1.7, l: 0.8, r: 0.8, w: 10,
},
jh: {
	p: 5, t: 8.2, k: 5.5, f: 4.5, th: 4.8, s: 6.2, sh: 5.2, ch: 18.3, hh: 2.6, b: 0.5, d: 17.7, g: 23, v: 2, dh: 7.5, z: 7.8, zh: 17, jh: 59.8, y: 5.8, m: 1.7, n: 0.3, ng: 1.6, l: 3.1, r: 1.4, w: 0,
},
y: {
	p: 1.3, t: 0, k: 0.8, f: 0.9, th: 0.8, s: 1.4, sh: 3, ch: 0.9, hh: 2.1, b: 4.2, d: 2.5, g: 1.3, v: 0.7, dh: 1.3, z: 1.6, zh: 0, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 0, l: 8.5, r: 2.6, w: 2.4,
},
m: {
	p: 3.8, t: 9.2, k: 2.5, f: 2.9, th: 4.6, s: 1.7, sh: 0.8, ch: 0.7, hh: 3.9, b: 9.6, d: 5.8, g: 1.5, v: 7.3, dh: 2.1, z: 1.5, zh: 1.2, jh: 2.4, y: 0.9, m: 100, n: 20, ng: 10.8, l: 11.7, r: 5.4, w: 4.2, ing: 99,
},
n: {
	p: 1.6, t: 9.6, k: 0.6, f: 1, th: 1.6, s: 0.9, sh: 0.3, ch: 0.3, hh: 3.2, b: 1.7, d: 8.3, g: 1.9, v: 4.5, dh: 2.1, z: 3.8, zh: 0.3, jh: 0.3, y: 1.2, m: 24.1, n: 94, ng: 14.1, l: 17.5, r: 4, w: 1.8, ing: 99,
},
ng: {
	p: 0, t: 6.7, k: 0.4, f: 2.1, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 0, b: 2.5, d: 5.8, g: 6.3, v: 1.9, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0, m: 13.8, n: 35, ng: 60, l: 2.1, r: 1.7, w: 10, ing: 99,
},
l: {
	p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, uwl: 99, axl: 99,
},
r: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
},
w: {
	p: 1.7, t: 0, k: 0.4, f: 1.2, th: 0.4, s: 1.2, sh: 0, ch: 0.9, hh: 1.7, b: 5.8, d: 0.8, g: 0, v: 2.1, dh: 0.4, z: 3.9, zh: 1, jh: 3.2, y: 6.5, m: 5.8, n: 0.8, ng: 10, l: 2.5, r: 4.2, w: 99,
},
iy: {
	p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 1.4, sh: 0.5, ch: 1.5, hh: 0.5, b: 8.3, d: 2.5, g: 3, v: 4.6, dh: 1.7, z: 1.6, zh: 3, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 3, l: 8.5, r: 2.6, w: 3, ih: 99, aa: 0.5, ao:0.5, aw:0.5,
},
ih: {
	p: 26.3, t: 4.6, k: 12.1, f: 0.5, th: 5, s: 0.4, sh: 0.4, ch: 3, hh: 99, b: 8.3, d: 1.3, g: 2, v: 4.6, dh: 1.7, z: 0.4, zh: 2, jh: 0.5, y: 3, m: 2, n: 3, ng: 2, l: 1.7, r: 1.5, w: 2,
},
ey: {
	p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 1.4, sh: 0.8, ch: 0.9, hh: 99, b: 8.3, d: 2.5, g: 1.3, v: 4.6, dh: 1.7, z: 1.6, zh: 3, jh: 0.5, y: 2.2, m: 10, n: 4.2, ng: 2.6, l: 8.5, r: 2.6, w: 1, ao: 0.3, ehr:99, eh:99, er:99,
},
eh: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 2.5, zh: 0.5, jh: 0.5, y: 0.5, m: 4.1, n: 4.1, ng: 1, l: 14.8, r: 96, w: 1, oy: 0.5, ax:99, axr: 99, ehr:99, ey:99, er:99,
},
ae: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 3.8, s: 0.8, sh: 0.4, ch: 0.9, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 2, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 0, l: 14.8, r: 96, w: 6.7, aa: 3, ao: 0.5,
},
aa: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 3, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.9, ng: 0.5, l: 14.8, r: 96, w: 6.7, ax:99, axr: 99,
},
ah: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ax:99, axr: 99,
},
ao: {
	p: 26.3, t: 0.5, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 1, l: 14.8, r: 96, w: 6.7, oy: 3, ih: 0.3, ax: 99, axl: 99,
},
ow: {
	 p: 26.3, t: 7.9, k: 12.1, f: 1, th: 5, s: 0.5, sh: 1.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 10, jh: 3.2, y: 0.5, m: 5.8, n: 2.1, ng: 10, l: 14.8, r: 96, w: 99,
},
uh: {
	p: 26.3, t: 4.6, k: 12.1, f: 11.3, th: 5, s: 0.4, sh: 0.4, ch: 0.8, hh: 99, b: 8.3, d: 1.3, g: 0.4, v: 4.6, dh: 1.7, z: 0.4, zh: 0, jh: 1, y: 0.8, m: 0, n: 0.8, ng: 0, l: 1.7, r: 0.8, w: 0.4, ax: 99, axr: 99, uwn: 99, uwl: 99,
},
uw: {
	p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axr: 99,
},
ay: {
	p: 0.5, t: 7.9, k: 12.1, f: 11.3, th: 3.8, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 3, m: 10, n: 4.2, ng: 10, l: 14.8, r: 96, w: 6.7,
},
oy: {
	p: 1.3, t: 0, k: 0.8, f: 0.9, th: 0.8, s: 1.4, sh: 0.8, ch: 0.9, hh: 2.1, b: 4.2, d: 2.5, g: 0.1, v: 0.7, dh: 1.3, z: 1.6, zh: 0, jh: 4.6, y: 70.6, m: 10, n: 4.2, ng: 0, l: 8.5, r: 2.6, w: 2.4,
},
aw: {
	p: 2.5, t: 7.9, k: 1.7, f: 0.5, th: 0.5, s: 0.8, sh: 0.4, ch: 0.5, hh: 5.8, b: 0.5, d: 0.5, g: 4.9, v: 0.5, dh: 0.8, z: 0.5, zh: 0.4, jh: 0.5, y: 4.6, m: 3.6, n: 2.1, ng: 10, l: 14.8, r: 96, w: 6.7, oy:0.5, ax: 99, axr: 99,
},
er: {
	p: 26.3, t: 7.9, k: 12.1, f: 1, th: 5, s: 0.5, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 1, dh: 1.7, z: 2.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 1, ax: 99, axr: 99,
},
ax: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ih: 0.3, aa: 99, ah: 99, oy: 3, axn: 99, axm: 99, axl: 99,
},
oh: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 3, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 2.9, ng: 0.5, l: 14.8, r: 96, w: 6.7, aa: 99, ah: 99, ax: 99, axr: 99,
},
ehr: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.5, sh: 0.5, ch: 0.5, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 1.7, z: 2.5, zh: 0.5, jh: 0.5, y: 0.5, m: 4.1, n: 4.1, ng: 1, l: 14.8, r: 96, w: 1, oy: 0.5, ax: 99, axr: 99, er:99, eh:99, ey:99,
},
axl: {
	p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99, axr: 99,
},
//axm: {
//	p: 3.8, t: 9.2, k: 2.5, f: 2.9, th: 4.6, s: 1.7, sh: 0.8, ch: 0.7, hh: 3.9, b: 9.6, d: 5.8, g: 1.5, v: 7.3, dh: 2.1, z: 1.5, zh: 1.2, jh: 2.4, y: 0.9, m: 100, n: 100, ng: 100, l: 11.7, r: 5.4, w: 4.2, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, axr: 99,
//},
axn: {
	p: 1.6, t: 9.6, k: 0.6, f: 1, th: 1.6, s: 0.9, sh: 0.3, ch: 0.3, hh: 3.2, b: 1.7, d: 8.3, g: 1.9, v: 4.5, dh: 2.1, z: 3.8, zh: 0.3, jh: 0.3, y: 1.2, m: 100, n: 100, ng: 100, l: 17.5, r: 4, w: 1.8, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, axr: 99,
},
ks: {
	p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao:1,
},
kw: {
	p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
},
dz: {
	p: 1.6, t: 3.8, k: 0.5, f: 2.6, th: 10, s: 15, sh: 15, ch: 2.5, hh: 0.5, b: 7.9, d: 8.8, g: 1, v: 13.9, dh: 23.8, z: 98, zh: 8.3, jh: 4, y: 2.5, m: 2.6, n: 8.1, ng: 0.5, l: 2.6, r: 0.5, w: 3.8, oy:0.5,
},
//ts: {
//	p: 6, t: 10.5, k: 5.2, f: 17.1, th: 24.6, s: 90, sh: 7, ch: 7, hh: 4, b: 0.5, d: 2.2, g: 0.5, v: 4.1, dh: 10, z: 25, zh: 7, jh: 2.7, y: 2, m: 0.8, n: 0.5, ng: 0.9, l: 1.4, r: 3, w: 0.5, ih: 0.3, aa: 0.5, ao:1,
//},
uwl: {
	p: 5.8, t: 8.3, k: 1.7, f: 5.8, th: 3.3, s: 2, sh: 0.5, ch: 1.2, hh: 2.1, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 2.5, jh: 4, y: 1.8, m: 10, n: 1.1, ng: 3.0, l: 99, r: 4.2, w: 5.4, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99, axr: 99,
},
uwn: {
	p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
},
//uwm: {
//	p: 1.7, t: 2, k: 2, f: 1.2, th: 2, s: 1.2, sh: 0, ch: 0.5, hh: 1.7, b: 5.8, d: 2, g: 3, v: 2.1, dh: 2, z: 0.5, zh: 0.5, jh: 0.5, y: 5.0, m: 5.8, n: 2, ng: 10, l: 2.5, r: 4.2, w: 99, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
//},
kl: {
	p: 25, t: 12.5, k: 98, f: 5.2, th: 5, s: 4, sh: 3.9, ch: 5.7, hh: 13.8, b: 4.2, d: 8, g: 12.5, v: 1.4, dh: 4.0, z: 0.4, zh: 0.5, jh: 6, y: 1.3, m: 1, n: 1.5, ng: 10, l: 1.6, r: 1.3, w: 1, ih: 2, oy: 0.5,
},
bl: {
	p: 8.1, t: 7.5, k: 6.7, f: 13.7, th: 7.4, s: 0.5, sh: 0.5, ch: 1.1, hh: 0.5, b: 91, d: 15, g: 14.6, v: 18.6, dh: 6.1, z: 2, zh: 0.5, jh: 0.5, y: 2.9, m: 7.6, n: 3.8, ng: 1.3, l: 5, r: 1.3, w: 5,
},
pl: {
	p: 86, t: 13.8, k: 18.2, f: 9.6, th: 8.8, s: 1.5, sh: 3, ch: 3.0, hh: 19.2, b: 21.3, d: 5.8, g: 3.1, v: 2.9, dh: 3.3, z: 3, zh: 0.5, jh: 0.5, y: 3, m: 1.4, n: 3, ng: 0.5, l: 1.3, r: 1, w: 1.3,
},
pr: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
},
kr: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
},
tr: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
},
gr: {
	p: 2.5, t: 7.9, k: 1.7, f: 2.1, th: 0.5, s: 0.9, sh: 1.5, ch: 1.5, hh: 5.8, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 0.8, z: 0.5, zh: 1, jh: 1.6, y: 0.5, m: 3.6, n: 2.1, ng: 0.5, l: 14.8, r: 96, w: 6.7,
},
ing: {
	p: 0, t: 6.7, k: 0.4, f: 2.1, th: 0.5, s: 0.5, sh: 0.5, ch: 0.5, hh: 0, b: 2.5, d: 5.8, g: 6.3, v: 1.9, dh: 1.7, z: 0.5, zh: 0.5, jh: 0.5, y: 0, m: 13.8, n: 99, ng: 99, l: 2.1, r: 1.7, w: 10, ih:99, axn: 99, axm: 99, uwn: 99, uwm: 99,
},
axr: {
	p: 26.3, t: 7.9, k: 12.1, f: 11.3, th: 5, s: 0.8, sh: 0.4, ch: 0.9, hh: 99, b: 14.6, d: 6.7, g: 4.9, v: 6.6, dh: 3, z: 0.5, zh: 0.4, jh: 0.5, y: 0.5, m: 3.6, n: 5, ng: 10, l: 14.8, r: 96, w: 6.7, ih: 0.3, aa: 99, ah: 99, oy: 3, ax: 99, axn: 99, axm: 99, axl: 99, uwn: 99, uwm: 99, uwl: 99,
},
fl: {
	p: 24.6, t: 21.7, k: 13.9, f: 99, th: 11.3, s: 7.5, sh: 2.4, ch: 1.5, hh: 9.2, b: 15, d: 10.4, g: 5.5, v: 5.4, dh: 5, z: 1.4, zh: 1, jh: 2, y: 2.0, m: 1.7, n: 2, ng: 1.2, l: 99, r: 15, w: 15, uwl: 99, axl: 99,
},
ihl: {
	p: 30, t: 8.3, k: 1.7, f: 5.8, th: 5, s: 1, sh: 0.5, ch: 1.0, hh: 30, b: 8.3, d: 1, g: 3.6, v: 8.7, dh: 3.3, z: 2.5, zh: 1, jh: 1, y: 3, m: 10, n: 3, ng: 3.0, l: 99, r: 10, w: 10, uwl: 99, axl: 99, ih: 99,
},
ehl: {
	p: 3.0, t: 3.0, k: 3.0, f: 8.8, th: 5.0, s: 2, sh: 0.5, ch: 1.0, hh: 99, b: 8.3, d: 2, g: 3.6, v: 8.7, dh: 4.0, z: 2.5, zh: 1.0, jh: 2, y: 1.0, m: 4, n: 1.1, ng: 1.0, l: 99, r: 4.2, w: 5.4, uwl: 99, axl: 99, oy: 0.5, ax:99, axr: 99, eh:99,
},

}
*/

// This is just one way of scoring the possible guards that are good to use
// between two adjacent expected phonemes. We could think up others
func biGuards(first, second map[phoneme]float64) map[phoneme]float64 {
	sum := make(map[phoneme]float64)
	diff := make(map[phoneme]float64)
	for k, v1 := range first {
		if v2, ok := second[k]; ok {
			sum[k] = v1 + v2
			diff[k] = math.Abs(v1 - v2)
		}
	}
	guards := make(map[phoneme]float64)
	for k, v1 := range sum {
		if v2, ok := diff[k]; ok {
			guards[k] = v1 + v2
		}
	}
	return guards
}

func (j *jsgfStandard) guardFor(phons []phoneme, currPh int) (phonRule, bool) {
	guards := make(map[phoneme]float64)

	if currPh > 0 {
		v1 := confusionMatrix[phons[currPh-1]]
		v2 := confusionMatrix[phons[currPh]]
		if v1 != nil && v2 != nil {
			guards = biGuards(v1, v2)
		}
	} else {
		guards = confusionMatrix[phons[currPh]]
	}
	if currPh < len(phons)-1 {
		v := confusionMatrix[phons[currPh+1]]
		if v != nil {
			guards = biGuards(guards, v)
		}
	}
	// Need to sort guards by value
	guardsByValue := make(map[float64][]phoneme)
	for k, v := range guards {
		if guardsByValue[v] == nil {
			guardsByValue[v] = []phoneme{}
		}
		guardsByValue[v] = append(guardsByValue[v], k)
	}
	var keys []float64
	for k := range guardsByValue {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// allow the same guard twice  -- PE  18 March 2020

	// Don't need to worry about removing the previous and next phoneme. These
	// will naturally appear later on in guardsByValue so should never be
	// selected. We do need to remove the previously used guards though
	//lastGuards := j.getPens()
	candidateGuards := make(map[phoneme]bool)
	str := ""
	for _, k := range keys {
		for _, ph := range guardsByValue[k] {
			str += fmt.Sprintln("Considering guard", ph, "with value", k, "for", phons[currPh], "...")
			//_debug("Considering guard", ph, "with value", k, "for", phons[currPh], "...")
			//if !contains(lastGuards, ph) {
			candidateGuards[ph] = true
			//}
		}
		str += fmt.Sprintln("candidateGuards =", candidateGuards)
		//_debug("candidateGuards =", candidateGuards)
		if len(candidateGuards) != 0 {
			// Randomly pick a guard - if there's more than one anyway...
			for k := range candidateGuards {
				str += fmt.Sprintln("Picking guard,", k)
				_debug(str)
				//_debug("Picking guard,", k)
				return phonRule{
					k,
				}, true
			}
		}
	}
	// If we get here we never got more than one guard. It should never
	// happen if the rest of the code is working okay...
	str += fmt.Sprintln("Failed to pick a guard!")
	_debug(str)
	return phonRule{}, false
}

func (j *jsgfStandard) closingPenalty(ph phoneme) (phonRule, bool) {
	// Don't want an f, p, etc or the phoneme

	// If the last target phoneme was a vowel, then the following were removed from guarding that last vowel
	// guards = consRemove(guards, ng, w, r, l, f, sh, p, b, th, s, t)
	// Removed from guarding vowels in general
	// guards = consRemove(guards, ng, b, p, k, hh, n, v, m, d, l, r, g, jh, dh)  - these have similar frequency

	//penalties := consRemove(newConsonants(), b, ph)

	//penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing)

	penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing, th, sts, st, kt, fl)

	//64only
	//penalties := consRemove(newConsonants(), b, ch, d, dh, f, "g", hh, jh, k, l, m, n, ng, p, r, s, sh, t, th, v, w, z, zh, dz, kw, ks, ts, kl, pl, bl, tr, kr, pr, gr, ing, th, fl)

	/*

	 if isVowel(ph) {
	 penalties = consRemove(penalties, hh, w, r, l, ng, p, f, dh, th, k, s, sh, z, zh, m, n, g, v)
	 } else {
	 penalties = consRemove(penalties, hh, w, r, l, ng, p, f, dh, th, k, s, sh, z, zh, v)
	 }
	*/

	//partner := []phoneme{}

	switch ph {
	case t:
		penalties[p] = true
		penalties[d] = true
		penalties[ch] = true
		penalties[g] = true
		penalties[k] = true
	case d:
		penalties[b] = true
		penalties[g] = true
		penalties[v] = true
		penalties[dh] = true
	case ch:
		penalties[p] = true
		penalties[hh] = true
		penalties[g] = true
		penalties[ch] = true
		penalties[sh] = true
	case g:
		penalties[d] = true
		penalties[hh] = true
		penalties[y] = true
	case k:
		penalties[p] = true
		penalties[g] = true
		penalties[t] = true
		penalties[hh] = true
	case l:
		penalties[r] = true
		penalties[w] = true
		penalties[f] = true
		penalties[b] = true
		penalties[d] = true
	case m:
		penalties[n] = true
		penalties[l] = true
		penalties[d] = true
		penalties[v] = true
		penalties[b] = true
	case n:
		penalties[m] = true
		penalties[ng] = true
		penalties[d] = true
		penalties[z] = true
		penalties[b] = true
	case p:
		penalties[k] = true
		penalties[t] = true
		penalties[b] = true
		penalties[v] = true
		penalties[hh] = true
		penalties[m] = true
	case r:
		penalties[l] = true
		penalties[w] = true
		penalties[v] = true
		penalties[hh] = true
	case s:
		penalties[z] = true
		penalties[t] = true
		penalties[sh] = true
		penalties[ch] = true
		penalties[th] = true
		penalties[f] = true
		penalties[hh] = true
	case v:
		penalties[z] = true
		penalties[m] = true
		penalties[n] = true
		penalties[l] = true
		penalties[d] = true
		penalties[r] = true
	case z:
		penalties[s] = true
		penalties[v] = true
		penalties[b] = true
	case eh:
		penalties[ng] = true
		penalties[n] = true
		penalties[r] = true
		penalties[l] = true
		penalties[dh] = true
		penalties[th] = true
	case ey:
		// have one of n or ng but not both because they are so similar likely to found similar noise
		penalties[n] = true
		penalties[r] = true
		penalties[w] = true
		penalties[z] = true
		penalties[hh] = true
		penalties[f] = true
	case uw:
		// have one of n or ng but not both because they are so similar likely to found similar noise
		penalties[n] = true
		penalties[r] = true
		penalties[w] = true
		penalties[d] = true
		penalties[hh] = true
		penalties[f] = true
	case er: // remove h sound
		penalties[n] = true
		penalties[r] = true
		penalties[w] = true
		penalties[d] = true
		penalties[l] = true
		penalties[b] = true
	case ah: // remove h sound
		penalties[n] = true
		penalties[r] = true
		penalties[w] = true
		penalties[d] = true
		penalties[l] = true
		penalties[b] = true
	case ao: // remove h & r sounds
		penalties[ch] = true
		penalties[n] = true
		penalties[r] = true
		penalties[w] = true
		penalties[d] = true
		penalties[hh] = true
	default:
		penalties[hh] = true
		penalties[m] = true
		penalties[ng] = true
		penalties[n] = true
		penalties[p] = true
		penalties[f] = true
	}

	/*
	   if ph != s {
	   	penalties[s] = true
	   }
	*/

	pen := pickConsonant(penalties)
	if ph != sil {
		return phonRule{
			pen,
		}, true
	} else {
		return phonRule{}, false
	}
}

func (j *jsgfStandard) phonemeHandler(ph phoneme) rule {
	if ph == k || ph == p {
		ph_r := new_R_ph(ph)
		return ph_r
		//return new_R_or(ph_r, new_R_and(ph_r, ph_r))
		//return new_R_or(ph_r, new_R_and(ph_r, ph_r), new_R_and(ph_r, new_R_ph(hh)))
		//return new_R_or(ph_r, new_R_and(ph_r, ph_r), new_R_and(ph_r, new_R_ph(hh)), new_R_ph(t), new_R_ph(g), new_R_ph(p))       // places a challenge of t, g & p versus k in black for example but creates problems in replied
	}
	return phonRule{
		ph,
	}
}

func (j *jsgfStandard) vowelHandler(ph phoneme) []phoneme {
	v_x := []phoneme{}
	if isVowel(ph) {
		v_x = append(v_x, ph)
	}
	if len(v_x) == 0 {
		// Something's gone wrong. Either vowelHandler has been called with ph not
		// a vowel or isVowel is causing a problem
		str := fmt.Sprintln("vowelHandler: ph =", ph)
		// Test isVowel and report anything unexpected
		vowels := []phoneme{
			//aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axm, axn, uwl, uwn, uwm, axr,

			aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axm, axn, uwl, uwn, uwm, axr, ihl, yuw, ehl,

			//64Only
			//aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axn, uwl, uwn, axr, ihl, ehl,
		}
		for _, vowel := range vowels {
			if !isVowel(vowel) {
				str += fmt.Sprintln("isVowel reporting", vowel, "not a vowel")
			}
		}
		_debug(str)
	}
	return v_x
}

type jsgfDiphthong struct {
	jsgfStandard
}

func new_jsgfDiphthong() jsgfDiphthong {
	return jsgfDiphthong{
		new_jsgfStandard(),
	}
}

/*
func (j jsgfDiphthong) getLastPen() phoneme {
  return j.jsgfStandard.getLastPen()
}

func (j *jsgfDiphthong) setLastPen(ph phoneme) {
  j.jsgfStandard.setLastPen(ph)
}
*/
func (j *jsgfDiphthong) getPens() []phoneme {
	return j.jsgfStandard.getPens()
}

func (j *jsgfDiphthong) setPen(ph phoneme) {
	j.jsgfStandard.setPen(ph)
}

func (j *jsgfDiphthong) addPen(ph phoneme) {
	j.jsgfStandard.addPen(ph)
}

/*
func (j *jsgfDiphthong) penalty(phons []phoneme, currPh int) (phonRule, bool) {
  return j.jsgfStandard.penalty(phons, currPh)
}
*/

func (j *jsgfDiphthong) guardFor(phons []phoneme, currPh int) (phonRule, bool) {
	return j.jsgfStandard.guardFor(phons, currPh)
}

func (j *jsgfDiphthong) phonemeHandler(ph phoneme) rule {
	if isDiphthong(ph) {
		diphs := diphthongs[ph]
		phonRules := []rule{}
		for _, ph := range diphs {
			new := phonRule{
				ph,
			}
			phonRules = append(phonRules, rule(new))
		}
		return new_R_and(phonRules...)
	}
	return j.jsgfStandard.phonemeHandler(ph)
}

func (j *jsgfDiphthong) vowelHandler(ph phoneme) []phoneme {
	v_x := []phoneme{}
	if isVowel(ph) {
		v_x = append(v_x, ph)
	}
	diphs, ok := diphthongs[ph]
	if ok {
		v_x = append(v_x, diphs[:]...)
	}
	if len(v_x) == 0 {
		// Something's gone wrong. Either vowelHandler has been called with ph not
		// a vowel or isVowel is causing a problem
		str := fmt.Sprintln("vowelHandler: ph =", ph)
		// Test isVowel and report anything unexpected
		vowels := []phoneme{
			//aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axm, axn, uwl, uwn, uwm, axr,

			aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axm, axn, uwl, uwn, uwm, axr, ihl, yuw, ehl,

			// 64only
			//aa, ae, ah, ao, aw, ax, ay, eh, ehr, er, ey, ih, iy, oh, ow, oy, uw, uh, y, axl, axn, uwl, uwn, axr, ihl, ehl,
		}
		for _, vowel := range vowels {
			if !isVowel(vowel) {
				str += fmt.Sprintln("isVowel reporting", vowel, "not a vowel")
			}
		}
		_debug(str)
	}
	return v_x
}
