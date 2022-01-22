package pron

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_R_and(t *testing.T) {
	// GIVEN
	r1 := phonRule{
		aa,
	}
	r := R_and{
		[]rule{
			r1,
		},
	}

	// WHEN
	actual := r.generate()

	// THEN
	if actual != "aa" {
		t.Error()
	}

	// and GIVEN
	r2 := phonRule{
		ae,
	}
	r = R_and{
		[]rule{
			r1, r2,
		},
	}

	// WHEN
	actual = r.generate()

	// THEN
	if actual != "(aa ae)" {
		t.Error()
	}
}

type testConfig struct {
	lastPens []phoneme
}

func new_testConfig() testConfig {
	return testConfig{
		[]phoneme{},
	}
}

func (j *testConfig) setPen(ph phoneme) {
	j.lastPens = []phoneme{
		ph,
	}
}

func (j *testConfig) addPen(ph phoneme) {
	j.lastPens = append(j.lastPens, ph)
}

func (j testConfig) getPens() []phoneme {
	return j.lastPens
}

func (j *testConfig) openingPenalty(ph phoneme) (phonRule, bool) {
	return phonRule{
		m,
	}, true
}

func (j *testConfig) guardFor(phons []phoneme, currPh int) (phonRule, bool) {
	if len(j.lastPens) == 0 || j.lastPens[0] == n {
		return phonRule{
			m,
		}, true
	}
	return phonRule{
		n,
	}, true
}

func (j *testConfig) closingPenalty(ph phoneme) (phonRule, bool) {
	if len(j.lastPens) == 0 || j.lastPens[0] == n {
		return phonRule{
			m,
		}, true
	}
	return phonRule{
		n,
	}, true
}

func (j testConfig) phonemeHandler(ph phoneme) rule {
	return phonRule{
		ph,
	}
}

func (j testConfig) vowelHandler(ph phoneme) []phoneme {
	return []phoneme{
		ph,
	}
}

func Test_openingRule(t *testing.T) {
	// GIVEN
	config := testConfig{}
	testGrammar := jsgfGrammar{
		"#JSGF V1.0",
		"grammar template",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{
			"any_vowel_noSlide": {
				"any_vowel_noSlide",
				nullRule{},
			},
		},
		&config,
	}
	// WHEN
	// with starting vowel...
	phons := []phoneme{
		aa,
	}
	r := testGrammar.openingRule(phons)
	actual := r.generate()

	// THEN
	if actual != "(sil [((n <any_Vx_aa_noSlide>) | (<any_Vx_aa_noSlide> n))] [((n <any_C>) | (<any_C> n))])" {
		t.Error()
	}

	// and GIVEN
	// Reset config
	config = testConfig{}
	testGrammar.config = &config

	// WHEN
	// with starting consonant...
	phons = []phoneme{
		b,
	}
	r = testGrammar.openingRule(phons)
	actual = r.generate()

	// THEN
	if actual != "(sil [((n <any_Cx_b>) | (<any_Cx_b> n))] [((n <any_vowel_noSlide>) | (<any_vowel_noSlide> n))])" {
		t.Error()
	}
}

func Test_trappedOpening(t *testing.T) {
	// GIVEN
	config := testConfig{}
	_ = jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}

	// WHEN
	r := R_trappedOpening{
		R_trap{
			aa,
		},
		R_opening{
			b,
			nullRule{},
		},
	}
	_ = r.generate()
	// expected := "sil aa | "

	// THEN
}

func Test_new_R_Cx(t *testing.T) {
	// GIVEN
	config := testConfig{}
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons := []phoneme{
		m, zh,
	}

	// WHEN
	_ = testGrammar.new_R_Cx(phons...)

	// THEN
}

func Test_preRule(t *testing.T) {
	// GIVEN
	config := new_testConfig()
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons := []phoneme{
		aa, b,
	}

	// WHEN
	currPh := 0
	pen, _ := config.guardFor(phons, currPh)
	ru := testGrammar.preRule(phons, currPh, pen)
	actual := ru.generate()

	// THEN
	if actual != "(([(<any_Vx_n_aa> n)] | [(n <any_Vx_n_aa>)]) ([(<any_Cx_n_b> n)] | [(n <any_Cx_n_b>)]))" {
		fmt.Println("!", actual)
		t.Error()
	}

	// and WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	ru = testGrammar.preRule(phons, currPh, pen)
	actual = ru.generate()

	// THEN
	if actual != "([(<any_Cx_n_b> n)] | [(n <any_Cx_n_b>)])" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		aa, ae,
	}

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	ru = testGrammar.preRule(phons, currPh, pen)
	actual = ru.generate()

	// THEN
	if actual != "(([(<any_Vx_n_ae_aa> n)] | [(n <any_Vx_n_ae_aa>)]) ([(<any_Cx_n> n)] | [(n <any_Cx_n>)]))" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		b, d,
	}

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	ru = testGrammar.preRule(phons, currPh, pen)
	actual = ru.generate()

	// THEN
	if actual != "([(<any_Cx_n_d_b> n)] | [(n <any_Cx_n_d_b>)])" {
		t.Error()
	}

	// and finally GIVEN
	config = testConfig{
		[]phoneme{
			n,
		},
	}
	testGrammar.config = &config

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	ru = testGrammar.preRule(phons, currPh, pen)
	actual = ru.generate()

	// THEN
	if actual != "([(<any_Cx_m_n_d_b> m)] | [(m <any_Cx_m_n_d_b>)])" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		r, iy, p, l, ay, d,
	}

	// WHEN
	currPh = 1
	config.setPen(d)
	ru = testGrammar.preRule(phons, currPh, phonRule{d})
	actual = ru.generate()

	// THEN
	fmt.Println("!*", actual)
	fmt.Println("ru =", ru)
}

func TestPostRule(t *testing.T) {
	// GIVEN
	config := new_testConfig()
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons := []phoneme{
		sil, aa, b, sil,
	}

	// WHEN
	currPh := 1
	pen, _ := config.guardFor(phons, currPh)
	r := testGrammar.postRule(phons, currPh, pen)
	actual := r.generate()

	// THEN
	if actual != "((<any_Vx_n_aa> n) | (n <any_Vx_n_aa>))" {
		fmt.Println("!!", actual)
		fmt.Println("r =", r)
		t.Error()
	}

	// and WHEN
	currPh = 2
	pen, _ = config.guardFor(phons, currPh)
	r = testGrammar.postRule(phons, currPh, pen)
	actual = r.generate()

	// THEN
	if actual != "((<any_Cx_n_b> n) | (n <any_Cx_n_b>))" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		sil, aa, ae, sil,
	}

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	r = testGrammar.postRule(phons, currPh, pen)
	actual = r.generate()

	// THEN
	if actual != "((<any_Vx_n_aa_ae> n) | (n <any_Vx_n_aa_ae>))" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		sil, b, d, sil,
	}

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	r = testGrammar.postRule(phons, currPh, pen)
	actual = r.generate()

	// THEN
	if actual != "((<any_Cx_n_b_d> n) | (n <any_Cx_n_b_d>))" {
		t.Error()
	}

	// and finally GIVEN
	config = testConfig{
		[]phoneme{
			n,
		},
	}
	testGrammar.config = &config

	// WHEN
	currPh = 1
	pen, _ = config.guardFor(phons, currPh)
	r = testGrammar.postRule(phons, currPh, pen)
	actual = r.generate()

	// THEN
	if actual != "((<any_Cx_m_n_b_d> m) | (m <any_Cx_m_n_b_d>))" {
		t.Error()
	}
}

func Test_middleRule(t *testing.T) {
	// GIVEN
	config := new_testConfig()
	testGrammar := jsgfGrammar{
		"#JSGF V1.0",
		"grammar template",
		R_target{
			"",
			[]parsableRule{
				R_opening{
					g,
					phonRule{
						sil,
					},
				},
				R_phoneme{
					aa,
					g,
					[]phoneme{},
					phonRule{
						aa,
					},
				},
				R_phoneme{
					b,
					g,
					[]phoneme{},
					phonRule{
						b,
					},
				},
				R_phoneme{
					d,
					g,
					[]phoneme{},
					phonRule{
						d,
					},
				},
				R_closing{
					g,
					phonRule{
						sil,
					},
				},
			},
		},
		map[string]namedRule{
			"any_Vx_aa": {
				"any_Vx_aa",
				nullRule{},
			},
			"any_Cx_b": {
				"any_Cx_b",
				nullRule{},
			},
			"any_Cx_b_d": {
				"any_Cx_b_d",
				nullRule{},
			},
		},
		&config,
	}
	phons := []phoneme{
		sil, aa, b, d, sil,
	}

	// WHEN
	currPh := 2
	ru := testGrammar.middleRule(phons, currPh)
	fmt.Println("ru =", ru)
	actual := ru.generate()

	// THEN
	if actual != "(([(<any_Cx_n_b_d> n)] | [(n <any_Cx_n_b_d>)]) (b | ((<any_Cx_n_b_d> n) | (n <any_Cx_n_b_d>))))" {
		t.Error()
	}

	// and finally GIVEN
	config = testConfig{
		[]phoneme{
			n,
		},
	}
	testGrammar.config = &config

	// WHEN
	currPh = 2
	ru = testGrammar.middleRule(phons, currPh)
	actual = ru.generate()

	// THEN
	if actual != "(([(<any_Cx_m_n_b_d> m)] | [(m <any_Cx_m_n_b_d>)]) (b | ((<any_Cx_m_n_b_d> m) | (m <any_Cx_m_n_b_d>))))" {
		t.Error()
	}

	// and WHEN
	currPh = 3
	ru = testGrammar.middleRule(phons, currPh)
	actual = ru.generate()

	// THEN
	if actual != "(([(<any_Cx_n_m_d_b> n)] | [(n <any_Cx_n_m_d_b>)]) (d | ((<any_Cx_n_m_d_b> n) | (n <any_Cx_n_m_d_b>))))" {
		t.Error()
	}

	// and GIVEN
	config = testConfig{
		[]phoneme{
			g,
		},
	}
	testGrammar = jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons = []phoneme{
		r, ih, p, l, ay, d,
	}

	// WHEN
	currPh = 1
	ru = testGrammar.middleRule(phons, currPh)
	actual = ru.generate()
	expected := "((([(<any_Vx_n_g_ih_noSlide> n)] | [(n <any_Vx_n_g_ih_noSlide>)]) ([(<any_Cx_n_g_p_r> n)] | [(n <any_Cx_n_g_p_r>)])) (ih | ((<any_Vx_n_g_ih_noSlide> n) | (n <any_Vx_n_g_ih_noSlide>))))"

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error(actual)
	}

	// and GIVEN
	stdConfig := jsgfStandard{
		[]phoneme{
			g,
		},
	}
	testGrammar.config = &stdConfig

	currPh = 1
	ru = testGrammar.middleRule(phons, currPh)
	actual = ru.generate()
	expected = "((([(<any_Vx_n_g_ih_noSlide> n)] | [(n <any_Vx_n_g_ih_noSlide>)]) ([(<any_Cx_n_g_p_r> n)] | [(n <any_Cx_n_g_p_r>)])) (ih | ((<any_Vx_n_g_ih_noSlide> n) | (n <any_Vx_n_g_ih_noSlide>))))"

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error(actual)
	}
}

// A test to make sure a named parsableRule is added to the grammar
func Test_softFade(t *testing.T) {
	// GIVEN
	config := testConfig{}
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}

	// WHEN
	_ = testGrammar.new_R_softFade()

	// THEN
	if _, ok := testGrammar.rules["soft_fade"]; !ok {
		t.Error()
	}
}

func Test_closingRule(t *testing.T) {
	// GIVEN
	config := testConfig{}
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons := []phoneme{
		aa,
	}

	// WHEN
	r := testGrammar.closingRule(phons)
	actual := r.generate()

	// THEN
	if actual != "(([(<any_Cx_n> n)] | [(n <any_Cx_n>)]) ([(<any_vowel_noSlide> n)] | [(n <any_vowel_noSlide>)]) [<soft_fade>] sil)" {
		t.Error()
	}

	// and GIVEN
	phons = []phoneme{
		b,
	}
	// Reset config
	config = testConfig{}
	testGrammar.config = &config

	// WHEN
	r = testGrammar.closingRule(phons)
	actual = r.generate()

	// THEN
	if actual != "(([(<any_vowel_noSlide> n)] | [(n <any_vowel_noSlide>)]) ([(<any_Cx_b> n)] | [(n <any_Cx_b>)]) [<soft_fade>] sil)" {
		t.Error()
	}
}

func Test_new_R_target(t *testing.T) {
	// GIVEN
	config := new_testConfig()
	testGrammar := jsgfGrammar{
		"",
		"",
		R_target{
			"",
			[]parsableRule{},
		},
		map[string]namedRule{},
		&config,
	}
	phons := []phoneme{
		aa,
	}

	// WHEN
	testGrammar.new_R_target("aah", phons)
	actual := testGrammar.target.generate()

	// THEN
	// Note the penalty flips between m and n because of the way testConfig.penalty is defined
	if actual != "public <aah> = ((sil [((n <any_Vx_aa_noSlide>) | (<any_Vx_aa_noSlide> n))] [((n <any_C>) | (<any_C> n))]) ((([(<any_Vx_m_n_aa> m)] | [(m <any_Vx_m_n_aa>)]) ([(<any_Cx_m_n> m)] | [(m <any_Cx_m_n>)])) (aa | ((<any_Vx_m_n_aa> m) | (m <any_Vx_m_n_aa>)))) (([(<any_Cx_n> n)] | [(n <any_Cx_n>)]) ([(<any_vowel_noSlide> n)] | [(n <any_vowel_noSlide>)]) [<soft_fade>] sil))" {

		t.Error()
	}
}

// func Test_SaveToDisk(t *testing.T) {
// 	// GIVEN
// 	config := testConfig{}
// 	word := "last"
// 	dict := dictionary.New("../../test/cmudict-0.7b.txt")
// 	entry := dict.Lookup(word)
// 	phonStrs := entry[0].Phonemes()
// 	phons := []phoneme{}
// 	for _, phStr := range phonStrs {
// 		phons = append(phons, phoneme(phStr))
// 	}
// 	grammar := jsgfGrammar{
// 		"#JSGF V1.0",
// 		"grammar test",
// 		R_target{
// 			"",
// 			[]parsableRule{},
// 		},
// 		map[string]namedRule{},
// 		&config,
// 	}

// 	// WHEN
// 	grammar.new_R_target(word, phons)

// 	// THEN
// 	grammar.SaveToDisk("savedGrammar.txt")
// }

func Test_phonemeHandler(t *testing.T) {
	// GIVEN
	config := jsgfStandard{}

	// and WHEN
	// Check this still works for phonemes other than k and p
	r := config.phonemeHandler(b)
	actual := r.generate()

	// THEN
	if actual != "b" {
		t.Error()
	}
}

func Test_parse(te *testing.T) {
	// GIVEN
	r := R_phoneme{
		t,
		g,
		[]phoneme{},
		R_and{
			[]rule{},
		},
	}
	results := psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0,
				10,
			},
			{
				t,
				11,
				20,
			},
			{
				eh,
				21,
				30,
			},
			{
				s,
				31,
				40,
			},
			{
				t,
				41,
				50,
			},
			{
				sil,
				51,
				60,
			},
		},
	}

	// WHEN
	actual, _ := r.parse(results.data, 4)
	expected := []parseResult{
		{
			4, 4, true,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	results = psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0,
				10,
			},
			{
				z,
				11,
				15,
			},
			{
				g,
				16,
				20,
			},
			{
				eh,
				21,
				30,
			},
			{
				s,
				31,
				40,
			},
			{
				t,
				41,
				50,
			},
			{
				sil,
				51,
				60,
			},
		},
	}

	// WHEN
	actual, err := r.parse(results.data, 1)
	expected = []parseResult{
		{
			1, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	results = psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0,
				10,
			},
			{
				z,
				11,
				14,
			},
			{
				g,
				15,
				18,
			},
			{
				t,
				19,
				20,
			},
			{
				eh,
				21,
				30,
			},
			{
				s,
				31,
				40,
			},
			{
				t,
				41,
				50,
			},
			{
				sil,
				51,
				60,
			},
		},
	}

	// WHEN
	actual, _ = r.parse(results.data, 1)
	expected = []parseResult{
		{
			1, 2, false,
		},
		{
			1, 3, true,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	results = psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0,
				10,
			},
			{
				z,
				11,
				15,
			},
			{
				g,
				16,
				20,
			},
			{
				eh,
				21,
				30,
			},
			{
				g,
				31,
				40,
			},
			{
				sil,
				41,
				50,
			},
		},
	}

	// WHEN
	actual, _ = r.parse(results.data, 1)
	expected = []parseResult{
		{
			1, 2, false,
		},
		{
			1, 4, false,
		},
	}
	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	results = psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0,
				10,
			},
			{
				z,
				11,
				15,
			},
			{
				g,
				16,
				20,
			},
			{
				eh,
				21,
				30,
			},
			{
				g,
				31,
				40,
			},
			{
				d,
				41,
				50,
			},
			{
				g,
				51,
				60,
			},
			{
				sil,
				71,
				80,
			},
		},
	}

	// WHEN
	actual, _ = r.parse(results.data, 1)
	expected = []parseResult{
		{
			1, 2, false,
		},
		{
			1, 4, false,
		},
		{
			1, 6, false,
		},
	}
	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rO := R_opening{
		v,
		R_and{
			[]rule{},
		},
	}
	data := []psPhonemeDatum{
		{
			sil,
			0,
			7,
		},
		{
			w,
			8,
			13,
		},
		{
			v,
			14,
			23,
		},
		{
			ow,
			24,
			42,
		},
		{
			dh,
			43,
			52,
		},
		{
			z,
			53,
			59,
		},
		{
			t,
			60,
			62,
		},
		{
			er,
			63,
			78,
		},
		{
			sil,
			79,
			115,
		},
	}

	// WHEN
	actual, _ = rO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
		{
			0, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rO = R_opening{
		b,
		R_and{
			[]rule{},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			47,
		},
		{
			b,
			48,
			51,
		},
		{
			eh,
			52,
			62,
		},
		{
			s,
			63,
			85,
		},
		{
			k,
			86,
			91,
		},
		{
			uw,
			92,
			103,
		},
		{
			l,
			104,
			123,
		},
		{
			sil,
			124,
			138,
		},
	}

	// WHEN
	actual, _ = rO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
		{
			0, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rO = R_opening{
		p,
		R_and{
			[]rule{},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			32,
		},
		{
			l,
			33,
			39,
		},
		{
			p,
			40,
			127,
		},
		{
			p,
			128,
			186,
		},
		{
			uw,
			187,
			204,
		},
		{
			w,
			205,
			207,
		},
		{
			s,
			208,
			228,
		},
		{
			aa,
			229,
			258,
		},
		{
			ng,
			259,
			274,
		},
		{
			n,
			275,
			300,
		},
		{
			l,
			301,
			303,
		},
		{
			iy,
			304,
			324,
		},
		{
			sil,
			325,
			430,
		},
	}

	// WHEN
	actual, _ = rO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
		{
			0, 2, false,
		},
		{
			0, 4, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rC := R_closing{
		sh,
		R_and{
			[]rule{},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			15,
		},
		{
			ch,
			16,
			18,
		},
		{
			th,
			19,
			72,
		},
		{
			iy,
			73,
			110,
		},
		{
			t,
			111,
			129,
		},
		{
			sil,
			156,
			186,
		},
	}

	// WHEN
	actual, _ = rC.parse(data, 5)
	expected = []parseResult{
		{
			5, 5, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// And some other tests around soft fade and the closing rule

	// and WHEN
	data = []psPhonemeDatum{
		{
			sil,
			156,
			186,
		},
	}

	// GIVEN
	actual, _ = rC.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	data = []psPhonemeDatum{
		{
			p,
			73,
			110,
		},
		{
			sh,
			111,
			129,
		},
		{
			sil,
			156,
			186,
		},
	}

	// WHEN
	actual, _ = rC.parse(data, 0)
	expected = []parseResult{
		{
			0, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	data = []psPhonemeDatum{
		{
			p,
			130,
			132,
		},
		{
			sh,
			133,
			155,
		},
		{
			sil,
			156,
			186,
		},
	}

	// WHEN
	actual, _ = rC.parse(data, 0)
	expected = []parseResult{
		{
			0, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	data = []psPhonemeDatum{
		{
			p,
			73,
			110,
		},
		{
			sh,
			111,
			129,
		},
		{
			sh,
			130,
			132,
		},
		{
			f,
			133,
			155,
		},
		{
			sil,
			156,
			186,
		},
	}

	// WHEN
	actual, _ = rC.parse(data, 0)
	expected = []parseResult{
		{
			0, 4, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rC = R_closing{
		dh,
		R_and{
			[]rule{},
		},
	}
	data = []psPhonemeDatum{
		{
			dh,
			52,
			63,
		},
		{
			sil,
			75,
			89,
		},
	}

	// WHEN
	actual, err = rC.parse(data, 0)
	expected = []parseResult{}

	// THEN
	if err == nil {
		te.Error()
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rP := R_phoneme{
		ah, d,
		[]phoneme{},
		R_and{
			[]rule{},
		},
	}
	data = []psPhonemeDatum{
		{
			d,
			39,
			41,
		},
		{
			y,
			42,
			51,
		},
		{
			ah,
			52,
			63,
		},
		{
			dh,
			52,
			63,
		},
		{
			sil,
			75,
			89,
		},
	}

	// WHEN
	actual, _ = rP.parse(data, 0)
	expected = []parseResult{
		{
			0, 1, false,
		},
		{
			0, 2, true,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rTO := R_trappedOpening{
		R_trap{
			zh, // A guess becasue the console log does't make clear what the trap is
		},
		R_opening{
			ch,
			R_and{
				[]rule{},
			},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			14,
		},
		{
			g,
			15,
			20,
		},
		{
			w,
			21,
			33,
		},
		{
			aa,
			34,
			57,
		},
		{
			z,
			58,
			60,
		},
		{
			g,
			61,
			63,
		},
	}

	// WHEN
	actual, _ = rTO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rTO = R_trappedOpening{
		R_trap{
			k,
		},
		R_opening{
			b,
			R_and{
				[]rule{},
			},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			18,
		},
		{
			k,
			19,
			22,
		},
		{
			sil,
			23,
			59,
		},
		{
			ch,
			60,
			62,
		},
		{
			dh,
			63,
			80,
		},
		{
			ih,
			81,
			97,
		},
		{
			v,
			98,
			117,
		},
		{
			sil,
			118,
			136,
		},
	}

	// WHEN
	actual, _ = rTO.parse(data, 0)
	expected = []parseResult{
		{
			0, 2, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rTO = R_trappedOpening{
		R_trap{
			n,
		},
		R_opening{
			b,
			R_and{
				[]rule{},
			},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0,
			4,
		},
		{
			ae,
			5,
			1,
		},
		{
			n,
			12,
			14,
		},
		{
			d,
			15,
			17,
		},
		{
			iy,
			18,
			21,
		},
		{
			m,
			22,
			31,
		},
		{
			ah,
			32,
			38,
		},
		{
			n,
			39,
			41,
		},
		{
			l,
			42,
			51,
		},
		{
			sil,
			52,
			59,
		},
	}

	// WHEN
	actual, _ = rTO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		fmt.Println(actual)
		te.Error()
	}

	// and GIVEN
	rTO = R_trappedOpening{
		R_trap{
			b,
		},
		R_opening{
			n, // Don't really care what this is
			R_and{
				[]rule{},
			},
		},
	}
	data = []psPhonemeDatum{
		{
			sil,
			0, 0,
		},
		{
			b,
			0, 0,
		},
		{
			aa, // Don't care what this is
			0, 0,
		},
	}

	// WHEN
	actual, _ = rTO.parse(data, 0)
	expected = []parseResult{
		{
			0, 0, false,
		},
	}

	// THEN
	fmt.Println(actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	data = []psPhonemeDatum{
		{
			sil,
			0, 0,
		},
		{
			b,
			0, 0,
		},
		{
			sil,
			0, 0,
		},
		{
			aa,
			0, 0,
		},
		{
			sil,
			0, 0,
		},
	}

	// WHEN
	actual, _ = rTO.parse(data, 0)
	expected = []parseResult{
		{
			0, 2, false,
		},
	}

	// THEN
	fmt.Println(actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_guardFor(te *testing.T) {
	// GIVEN
	jsgf := jsgfStandard{
		[]phoneme{
			ng,
		},
	}
	phons := []phoneme{
		r, ih, z, ay, d,
	}

	// WHEN
	_, ok := jsgf.guardFor(phons, 3)

	// THEN
	if !ok {
		te.Error()
	}

	// and GIVEN
	jsgfDiph := jsgfDiphthong{
		jsgfStandard{
			[]phoneme{
				v, f,
			},
		},
	}
	phons = []phoneme{
		t, ey, l,
	}

	// WHEN
	actual, ok := jsgfDiph.guardFor(phons, 0)

	// THEN
	if !ok {
		fmt.Println(actual)
		te.Error()
	}

	// and GIVEN
	jsgf = jsgfStandard{
		[]phoneme{
			zh,
		},
	}
	phons = []phoneme{
		k, ao, t,
	}

	// WHEN
	actual, ok = jsgf.guardFor(phons, 2)

	// THEN
	fmt.Println("caught actual =", actual)
	if !ok {
		te.Error()
	}
}

func Test_clayMay(te *testing.T) {
	// GIVEN
	results := psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0, 16,
			},
			{
				m,
				17, 17,
			},
			{
				zh,
				20, 20,
			},
			{
				w,
				23, 40,
			},
			{
				ey,
				41, 71,
			},
			{
				m,
				72, 83,
			},
			{
				jh,
				84, 88,
			},
			{
				sil,
				89, 137,
			},
		},
	}

	// and GIVEN
	rTO := R_trappedOpening{
		R_trap{
			g,
		},
		R_opening{
			s,
			R_and{
				[]rule{},
			},
		},
	}

	// WHEN
	actual, _ := rTO.parse(results.data, 0)
	expected := []parseResult{
		{
			0, 0, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	r := R_phoneme{
		m,
		dh,
		[]phoneme{},
		R_and{
			[]rule{},
		},
	}

	// WHEN
	actual, _ = r.parse(results.data, 1)
	expected = []parseResult{
		{
			1, 1, true,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	r = R_phoneme{
		ey,
		zh,
		[]phoneme{},
		R_and{
			[]rule{},
		},
	}

	// WHEN
	actual, _ = r.parse(results.data, 2)
	expected = []parseResult{
		{
			2, 3, false,
		},
		{
			2, 4, true,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	rC := R_closing{
		jh,
		R_and{
			[]rule{},
		},
	}

	// WHEN
	actual, _ = rC.parse(results.data, 5)
	expected = []parseResult{
		{
			5, 7, false,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and WHEN
	actual, err := rC.parse(results.data, 4)

	// THEN
	if err == nil {
		// There should be an error as the 4 is derived from a previous parse of
		// the last phoneme (see the {2, 3, false} result of the previous test)
		te.Error()
	}

	// and GIVEN
	results = psPhonemeResults{
		120,
		[]psPhonemeDatum{
			{
				sil,
				0, 8,
			},
			{
				s,
				9, 13,
			},
			{
				ch,
				14, 14,
			},
			{
				dh,
				17, 20,
			},
			{
				ih,
				21, 21,
			},
			{
				z,
				23, 35,
			},
			{
				w,
				36, 36,
			},
			{
				ih,
				38, 62,
			},
			{
				v,
				63, 63,
			},
			{
				z,
				65, 75,
			},
			{
				v,
				76, 76,
			},
			{
				sil,
				78, 104,
			},
		},
	}
	rC = R_closing{
		v,
		R_and{
			[]rule{},
		},
	}

	// WHEN
	actual, err = rC.parse(results.data, 7)
	expected = []parseResult{
		{
			7, 11, false,
		},
	}

	// THEN
	if err != nil {
		te.Error()
	} else {
		if !reflect.DeepEqual(actual, expected) {
			te.Error()
		}
	}
}
