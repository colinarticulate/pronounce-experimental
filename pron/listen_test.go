package pron

import (
	"fmt"
	"reflect"
	"testing"
)

func deepEqual(p, r []candidatePhoneme) bool {
	contains := func(x candidatePhoneme, y []candidatePhoneme) bool {
		for _, y1 := range y {
			if x == y1 {
				return true
			}
		}
		return false
	}
	if len(p) != len(r) {
		return false
	}
	for _, p1 := range p {
		if !contains(p1, r) {
			return false
		}
	}
	return true
}

func Test_powerset(te *testing.T) {
	// GIVEN
	ps := candidateData{
		{
			psPhonemeDatum{
				t,
				3,
				6,
			},
			2,
		},
	}

	// WHEN
	actual := ps.powerset()
	expected := []candidatePhoneme{
		{
			psPhonemeDatum{
				t,
				3,
				4,
			},
			2,
		},
		{
			psPhonemeDatum{
				t,
				4,
				5,
			},
			2,
		},
		{
			psPhonemeDatum{
				t,
				5,
				6,
			},
			2,
		},
		{
			psPhonemeDatum{
				t,
				3,
				5,
			},
			2,
		},
		{
			psPhonemeDatum{
				t,
				4,
				6,
			},
			2,
		},
		{
			psPhonemeDatum{
				t,
				3,
				6,
			},
			2,
		},
	}

	// THEN
	if !deepEqual(actual, expected) {
		te.Error()
	}
}

func Test_merge(t *testing.T) {
	// GIVEN
	results := []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					ae,
					7,
					16,
				},
				{
					n,
					17,
					20,
				},
				{
					ih,
					21,
					36,
				},
				{
					m,
					37,
					48,
				},
				{
					f,
					49,
					54,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					ae,
					7,
					16,
				},
				{
					n,
					17,
					21,
				},
				{
					ih,
					22,
					27,
				},
				{
					m,
					28,
					30,
				},
				{
					l,
					31,
					50,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					ae,
					9,
					17,
				},
				{
					n,
					18,
					21,
				},
				{
					ih,
					22,
					37,
				},
				{
					m,
					38,
					46,
				},
				{
					v,
					47,
					50,
				},
				{
					th,
					51,
					58,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					ae,
					8,
					16,
				},
				{
					n,
					17,
					20,
				},
				{
					ih,
					21,
					26,
				},
				{
					m,
					27,
					29,
				},
				{
					l,
					30,
					49,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					ae,
					10,
					17,
				},
				{
					n,
					17,
					22,
				},
				{
					ih,
					22,
					27,
				},
				{
					m,
					27,
					50,
				},
				{
					th,
					51,
					57,
				},
			},
		},
	}

	// WHEN
	actual := merge(results)
	expected := []candidatePhoneme{}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func Test_trimInitialPlosives(t *testing.T) {
	// GIVEN
	results := []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					b,
					20,
					35,
				},
				{
					ih,
					36,
					100,
				},
				{
					m,
					101,
					120,
				},
			},
		},
	}

	// WHEN
	actual := trimInitialPlosives(results)
	expected := []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					ih,
					36,
					100,
				},
				{
					m,
					101,
					120,
				},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func Test_getVerdict(t *testing.T) {
	// GIVEN
	phToAbcEntry := phonToAlphas{
		[]phoneme{
			aa, b, k,
		},
		"abc",
	}
	verdicts := []phonVerdict{
		{
			aa,
			good,
		},
		{
			b,
			possible,
		},
		{
			k,
			good,
		},
	}

	// WHEN
	actual := getVerdict(phToAbcEntry, verdicts)
	expected := []LettersVerdict{
		{
			"abc",
			[]phoneme{
				aa, b, k,
			},
			possible,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}

	// and GIVEN
	verdicts = []phonVerdict{
		{
			aa,
			good,
		},
		{
			b,
			possible,
		},
		{
			er,
			surprise,
		},
		{
			k,
			good,
		},
	}

	// WHEN
	actual = getVerdict(phToAbcEntry, verdicts)
	expected = []LettersVerdict{
		{
			"abc",
			[]phoneme{
				aa, b, k,
			},
			missing,
		},
		{
			"er",
			[]phoneme{
				er,
			},
			surprise,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func Test_lookFor(te *testing.T) {
	// GIVEN
	result := variantResult{
		0,
		[]phonVerdict{
			{
				hh, missing,
			},
			{
				iy, surprise,
			},
			{
				ih, surprise,
			},
			{
				iy, good,
			},
			{
				ah, missing,
			},
			{
				er, surprise,
			},
		},
		[]phoneme{
			hh, iy, ah,
		},
	}

	// WHEN
	actual := lookFor(iy, result.verdict)
	expected := []int{
		1, 3,
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		fmt.Println("here actual =", actual)
		te.Error()
	}

	// and GIVEN
	// Check that lookFor stops when it finds a good or possible phoneme before
	// it finds the phoneme it's looking for
	result = variantResult{
		0,
		[]phonVerdict{
			{
				hh, good,
			},
			{
				iy, surprise,
			},
			{
				ih, surprise,
			},
			{
				iy, good,
			},
			{
				ah, missing,
			},
			{
				er, surprise,
			},
		},
		[]phoneme{
			hh, iy, ah,
		},
	}

	// WHEN
	actual = lookFor(iy, result.verdict)

	// THEN
	if len(actual) != 0 {
		te.Error()
	}
}

/*
func Test_trySpelling(te *testing.T) {
  // GIVEN
  result := variantResult{
    0,
    []phonVerdict{
      {
        hh, good,
      },
      {
        iy, surprise,
      },
      {
        ih, surprise,
      },
      {
        iy, good,
      },
      {
        ah, missing,
      },
      {
        er, surprise,
      },
    },
    []phoneme{
      hh, iy, ah,
    },
  }
  spellings := [][]phoneme{
    []phoneme{
      hh, iy, ah,
    },
    []phoneme{
      hh, iy, er,
    },
  }

  // WHEN
  _, found := trySpelling(result, 0, spellings[0])

  // THEN
  if found {
    te.Error()
  }

  // and WHEN
  actual, found := trySpelling(result, 0, spellings[1])
  expected := variantResult{
    0,
    []phonVerdict{
      {
        hh, good,
      },
      {
        iy, surprise,
      },
      {
        ih, surprise,
      },
      {
        iy, good,
      },
      {
        er, possible,
      },
    },
    []phoneme{
      hh, iy, er,
    },
  }

  // THEN
  if !found {
    te.Error()
  } else {
    if !reflect.DeepEqual(actual, expected) {
      te.Error()
    }
  }

  // and GIVEN
  result = variantResult{
    0,
    []phonVerdict{
      {
        m, missing,
      },
      {
        g, surprise,
      },
      {
        w, surprise,
      },
      {
        ey, good,
      },
    },
    []phoneme{
      m, ey,
    },
  }

  // WHEN
  actual, found = trySpelling(result, 0, result.phons)
  expected = variantResult{
    0,
    []phonVerdict{
      {
        m, missing,
      },
      {
        g, surprise,
      },
      {
        w, surprise,
      },
      {
        ey, good,
      },
    },
    []phoneme{
      m, ey,
    },
  }

  // THEN
  if found {
    te.Error()
  }
}
*/

func Test_trySpelling2(te *testing.T) {
	// GIVEN
	result := []phonVerdict{
		{
			z, surprise,
		},
		{
			ah, surprise,
		},
		{
			dh, possible,
		},
		{
			iy, missing,
		},
	}
	spelling := []phoneme{
		dh, ah,
	}

	// WHEN
	actual, found := trySpelling(result, spelling)

	// THEN
	if found {
		fmt.Println("the actual =", actual)
		te.Error()
	}

	// and GIVEN
	result = []phonVerdict{
		{
			hh, good,
		},
		{
			iy, surprise,
		},
		{
			ih, surprise,
		},
		{
			iy, good,
		},
		{
			ah, missing,
		},
		{
			er, surprise,
		},
	}
	spelling = []phoneme{
		hh, iy, er,
	}

	// WHEN
	actual, found = trySpelling(result, spelling)
	expected := []phonVerdict{
		{
			hh, good,
		},
		{
			iy, surprise,
		},
		{
			ih, surprise,
		},
		{
			iy, good,
		},
		{
			er, possible,
		},
	}

	// THEN
	if !found {
		te.Error()
	} else {
		if !reflect.DeepEqual(actual, expected) {
			te.Error()
		}
	}

	// and GIVEN
	result = []phonVerdict{
		{
			z, surprise,
		},
		{
			dh, missing,
		},
		{
			ah, good,
		},
	}
	spelling = []phoneme{
		dh, ah,
	}

	// WHEN
	actual, found = trySpelling(result, spelling)

	// THEN
	if found {
		te.Error()
	}

	// GIVEN
	result = []phonVerdict{
		{w, good}, {ao, good}, {n, good}, {d, good}, {ah, good}, {r, good}, {ah, missing}, {r, possible},
	}
	spelling = []phoneme{
		w, ao, n, d, ah, r,
	}

	// WHEN
	fmt.Println("trySpelling for wanderer...")
	actual, found = trySpelling(result, spelling)

	// THEN
	fmt.Println("wanderer, actual =", actual)
	if found {
		te.Error()
	}
}

func Test_searchForBest(te *testing.T) {
	// GIVEN
	// A real example
	result := variantResult{
		0,
		[]phonVerdict{
			{
				hh, good,
			},
			{
				er, missing,
			},
			{
				ah, surprise,
			},
		},
		[]phoneme{
			hh, er,
		},
	}
	spellings := [][]phoneme{
		{
			hh, er,
		},
		{
			hh, ah,
		},
	}

	// WHEN
	results := []variantResult{
		result,
	}
	actual, found := searchForBest(results, spellings)
	expected := variantResult{
		0,
		[]phonVerdict{
			{
				hh, good,
			},
			{
				ah, possible,
			},
		},
		[]phoneme{
			hh, ah,
		},
	}

	// THEN
	if !found {
		te.Error()
	} else {
		if !reflect.DeepEqual(actual, expected) {
			fmt.Println("her, actual =", actual)
			te.Error()
		}
	}

	// and GIVEN
	result = variantResult{
		0,
		[]phonVerdict{
			{
				m, missing,
			},
			{
				g, surprise,
			},
			{
				w, surprise,
			},
			{
				ey, good,
			},
		},
		[]phoneme{
			m, ey,
		},
	}

	// WHEN
	results = []variantResult{
		result,
	}
	spellings = [][]phoneme{
		result.phons,
	}
	actual, found = searchForBest(results, spellings)

	// THEN
	if found {
		te.Error()
	}

	// and GIVEN
	result = variantResult{
		0,
		[]phonVerdict{
			{
				p, missing,
			},
			{
				ae, good,
			},
			{
				s, good,
			},
			{
				t, good,
			},
		},
		[]phoneme{
			p, ae, s, t,
		},
	}

	// WHEN
	results = []variantResult{
		result,
	}
	spellings = [][]phoneme{
		{
			p, ae, s, t,
		},
		{
			p, aa, s, t,
		},
	}
	actual, found = searchForBest(results, spellings)

	// THEN
	if found {
		te.Error()
	}

	// and GIVEN
	results = []variantResult{
		{
			-10,
			[]phonVerdict{
				{
					er, missing,
				},
				{
					eh, missing,
				},
				{
					s, missing,
				},
				{
					t, missing,
				},
				{
					s, missing,
				},
				{
					ae, surprise,
				},
				{
					er, surprise,
				},
				{
					eh, surprise,
				},
				{
					s, surprise,
				},
				{
					t, surprise,
				},
			},
			[]phoneme{
				er, eh, s, t, s,
			},
		},
		{
			2,
			[]phonVerdict{
				{
					aa, missing,
				},
				{
					r, good,
				},
				{
					eh, good,
				},
				{
					s, good,
				},
				{
					t, good,
				},
				{
					s, missing,
				},
			},
			[]phoneme{
				aa, r, eh, s, t, s,
			},
		},
		{
			1,
			[]phonVerdict{
				{
					ah, missing,
				},
				{
					ae, surprise,
				},
				{
					r, good,
				},
				{
					eh, good,
				},
				{
					s, good,
				},
				{
					t, good,
				},
				{
					s, missing,
				},
			},
			[]phoneme{
				ah, r, eh, s, t, s,
			},
		},
	}

	// WHEN
	spellings = [][]phoneme{
		{
			er, eh, s, t, s,
		},
		{
			aa, r, eh, s, t, s,
		},
		{
			ah, r, eh, s, t, s,
		},
	}
	actual, found = searchForBest(results, spellings)

	// THEN
	if found {
		te.Error()
	}

	// and GIVEN
	results = []variantResult{
		{
			-2,
			[]phonVerdict{
				{
					z, surprise,
				},
				{
					ah, surprise,
				},
				{
					dh, possible,
				},
				{
					iy, missing,
				},
			},
			[]phoneme{
				dh, iy,
			},
		},
		{
			-4,
			[]phonVerdict{
				{
					z, surprise,
				},
				{
					ah, surprise,
				},
				{
					dh, missing,
				},
				{
					er, missing,
				},
			},
			[]phoneme{
				dh, er,
			},
		},
		{
			-1,
			[]phonVerdict{
				{
					z, surprise,
				},
				{
					dh, missing,
				},
				{
					ah, good,
				},
			},
			[]phoneme{
				dh, ah,
			},
		},
	}

	// WHEN
	spellings = [][]phoneme{
		{
			dh, iy,
		},
		{
			dh, er,
		},
		{
			dh, ah,
		},
	}
	actual, found = searchForBest(results, spellings)

	// THEN
	if found {
		te.Error()
	}

	// and GIVEN
	results = []variantResult{
		{
			0,
			[]phonVerdict{
				{
					w, good,
				},
				{
					ao, good,
				},
				{
					z, missing,
				},
				{
					zh, surprise,
				},
			},
			[]phoneme{
				w, ao, z,
			},
		},
		{
			-1,
			[]phonVerdict{
				{
					w, good,
				},
				{
					aa, missing,
				},
				{
					ao, surprise,
				},
				{
					g, surprise,
				},
				{
					z, good,
				},
			},
			[]phoneme{
				w, aa, z,
			},
		},
		{
			0,
			[]phonVerdict{
				{
					w, good,
				},
				{
					ah, missing,
				},
				{
					ao, surprise,
				},
				{
					z, possible,
				},
			},
			[]phoneme{
				w, ah, z,
			},
		},
	}

	// WHEN
	spellings = [][]phoneme{
		{
			w, ao, z,
		},
		{
			w, aa, z,
		},
		{
			w, ah, z,
		},
	}
	actual, found = searchForBest(results, spellings)
	expected = variantResult{
		0,
		[]phonVerdict{
			{
				w, good,
			},
			{
				ao, possible,
			},
			{
				z, possible,
			},
		},
		[]phoneme{
			w, ah, z,
		},
	}

	// THEN
	if !found {
		te.Error()
	} else {
		if !reflect.DeepEqual(actual, expected) {
			fmt.Println("was actual =", actual)
			te.Error()
		}
	}

	// and GIVEN
	results = []variantResult{
		{
			-1,
			[]phonVerdict{
				{
					ch, good,
				},
				{
					eh, missing,
				},
				{
					ae, surprise,
				},
			},
			[]phoneme{
				ch, eh,
			},
		},
		{
			0,
			[]phonVerdict{
				{
					ch, good,
				},
				{
					eh, possible,
				},
				{
					er, missing,
				},
				{
					ae, surprise,
				},
			},
			[]phoneme{
				ch, eh, er,
			},
		},
		{
			0,
			[]phonVerdict{
				{
					ch, possible,
				},
				{
					eh, missing,
				},
				{
					ae, surprise,
				},
				{
					r, good,
				},
			},
			[]phoneme{
				ch, eh, r,
			},
		},
	}
	spellings = [][]phoneme{
		{
			ch, eh,
		},
		{
			ch, eh, er,
		},
		{
			ch, eh, r,
		},
	}

	// WHEN
	actual, found = searchForBest(results, spellings)
	expected = variantResult{
		0,
		[]phonVerdict{
			{
				ch, good,
			},
			{
				eh, possible,
			},
			{
				ae, surprise,
			},
		},
		[]phoneme{
			ch, eh, er,
		},
	}

	// THEN
	if !found {
		te.Error()
	} else {
		if !reflect.DeepEqual(actual, expected) {
			te.Error()
		}
	}

	// and GIVEN
	results = []variantResult{
		{
			5,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {ah, good}, {r, good}, {aa, surprise},
			},
			[]phoneme{
				w, ao, n, d, ah, r,
			},
		},
		{
			3,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {r, possible}, {er, missing}, {aa, surprise},
			},
			[]phoneme{
				w, ao, n, d, r, er,
			},
		},
		{
			3,
			[]phonVerdict{
				{w, good}, {aa, good}, {n, good}, {d, good}, {er, good}, {er, missing}, {aa, surprise},
			},
			[]phoneme{
				w, ao, n, d, er, er,
			},
		},
		{
			5,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {er, good}, {r, good}, {ah, surprise},
			},
			[]phoneme{
				w, aa, n, d, er, r,
			},
		},
		{
			3,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {er, good}, {r, possible}, {er, missing}, {eh, surprise}, {s, surprise},
			},
			[]phoneme{
				w, ao, n, d, er, r, er,
			},
		},
		{
			5,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {ah, good}, {r, possible}, {ah, missing},
			},
			[]phoneme{
				w, ao, n, d, ah, r, ah,
			},
		},
		{
			4,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {ah, good}, {r, good}, {er, missing}, {aa, surprise},
			},
			[]phoneme{
				w, ao, n, d, ah, r, er,
			},
		},
		{
			6,
			[]phonVerdict{
				{w, good}, {ao, good}, {n, good}, {d, good}, {ah, good}, {r, good}, {ah, missing}, {r, possible},
			},
			[]phoneme{
				w, ao, n, d, ah, r, ah, r,
			},
		},
	}
	spellings = [][]phoneme{
		{
			w, ao, n, d, ah, r,
		},
		{
			w, ao, n, d, r, er,
		},
		{
			w, ao, n, d, er, er,
		},
		{
			w, aa, n, d, er, r,
		},
		{
			w, ao, n, d, er, r, er,
		},
		{
			w, ao, n, d, ah, r, ah,
		},
		{
			w, ao, n, d, ah, r, er,
		},
		{
			w, ao, n, d, ah, r, ah, r,
		},
	}

	// WHEN
	actual, found = searchForBest(results, spellings)

	// THEN
	fmt.Println("wanderer, actual =", actual)
}

func Test_betterResult(t *testing.T) {
	// GIVEN
	results := []variantResult{
		{
			7,
			[]phonVerdict{
				{
					hh,
					good,
				},
				{
					aa,
					good,
				},
				{
					d,
					missing,
				},
				{
					l,
					good,
				},
				{
					iy,
					good,
				},
			},
			[]phoneme{
				hh, aa, d, l, iy,
			},
		},
		{
			8,
			[]phonVerdict{
				{
					hh,
					good,
				},
				{
					aa,
					good,
				},
				{
					r,
					missing,
				},
				{
					d,
					possible,
				},
				{
					l,
					good,
				},
				{
					iy,
					good,
				},
			},
			[]phoneme{
				hh, aa, r, d, l, iy,
			},
		},
	}
	best := results[1]

	// WHEN
	isBest := couldBeBetter(best)
	spellings := [][]phoneme{
		{
			hh, aa, d, l, iy,
		},
	}
	better, found := searchForBetter(results[1], spellings)

	expected := variantResult{
		8,
		[]phonVerdict{
			{
				hh,
				good,
			},
			{
				aa,
				good,
			},
			{
				d,
				possible,
			},
			{
				l,
				good,
			},
			{
				iy,
				good,
			},
		},
		[]phoneme{
			hh, aa, r, d, l, iy,
		},
	}

	// THEN
	if isBest != true {
		t.Error()
	}
	if found != true {
		t.Error()
	} else {
		if !reflect.DeepEqual(better, expected) {
			fmt.Println("hardly, better =", better)
			t.Error()
		}
	}

	// and GIVEN
	results = []variantResult{
		{
			7,
			[]phonVerdict{
				{
					hh,
					good,
				},
				{
					aa,
					good,
				},
				{
					d,
					missing,
				},
				{
					l,
					good,
				},
				{
					iy,
					good,
				},
			},
			[]phoneme{
				hh, aa, d, l, iy,
			},
		},
		{
			10,
			[]phonVerdict{
				{
					hh,
					good,
				},
				{
					aa,
					good,
				},
				{
					r,
					possible,
				},
				{
					d,
					possible,
				},
				{
					l,
					good,
				},
				{
					iy,
					good,
				},
			},
			[]phoneme{
				hh, aa, r, d, l, iy,
			},
		},
	}
	best = results[1]

	// WHEN
	isBest = couldBeBetter(best)

	// THEN
	if isBest != false {
		t.Error()
	}
}

func Test_searchForBetter(te *testing.T) {
	// GIVEN
	result := variantResult{
		0,
		[]phonVerdict{
			{ch, good}, {eh, possible}, {er, missing}, {ae, surprise},
		},
		[]phoneme{
			ch, eh, er,
		},
	}
	spellings := [][]phoneme{
		{
			ch, eh,
		},
		{
			ch, eh, er,
		},
		{
			ch, eh, r,
		},
	}

	// WHEN
	actual, ok := searchForBetter(result, spellings)
	expected := variantResult{
		0,
		[]phonVerdict{
			{ch, good}, {eh, possible}, {ae, surprise},
		},
		[]phoneme{
			ch, eh,
		},
	}

	// THEN
	if ok != true {
		te.Error()
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_trySpelling(te *testing.T) {
	// GIVEN
	result := []phonVerdict{
		{ch, good}, {eh, possible}, {er, missing}, {ae, surprise},
	}
	spelling := []phoneme{
		ch, eh,
	}

	// WHEN
	actual, ok := trySpelling(result, spelling)
	expected := []phonVerdict{
		{ch, good}, {eh, possible}, {ae, surprise},
	}

	// THEN
	if ok != true {
		te.Error()
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

/*
func Test_searchForBetter(te *testing.T) {
  // GIVEN
  results := []variantResult{
    {
      0,
      []phonVerdict{
        {
          hh, good,
        },
        {
          aa, good,
        },
        {
          d, missing,
        },
        {
          l, good,
        },
        {
          iy, good,
        },
      },
      []phoneme{
        hh, aa, d, l, iy,
      },
    },
    {
      0,
      []phonVerdict{
        {
          hh, good,
        },
        {
          aa, good,
        },
        {
          r, missing,
        },
        {
          d, possible,
        },
        {
          l, good,
        },
        {
          iy, good,
        },
      },
      []phoneme{
        hh, aa, r, d, l, iy,
      },
    },
  }

  // WHEN
  actual, ok := searchForBetter(results)
  expected := variantResult{
    0,
    []phonVerdict{
      {
        hh, good,
      },
      {
        aa, good,
      },
      {
        d, possible,
      },
      {
        l, good,
      },
      {
        iy, good,
      },
    },
    []phoneme{
      hh, aa, d, l, iy,
    },
  }

  // THEN
  if ok != true {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  results = []variantResult{
    {
      0,
      []phonVerdict{
        {
          w, good,
        },
        {
          ao, good,
        },
        {
          z, missing,
        },
      },
      []phoneme{
        w, ao, z,
      },
    },
    {
      0,
      []phonVerdict{
        {
          w, good,
        },
        {
          aa, missing,
        },
        {
          ao, surprise,
        },
        {
          z, missing,
        },
      },
      []phoneme{
        w, aa, z,
      },
    },
    {
      0,
      []phonVerdict{
        {
          w, good,
        },
        {
          ah, missing,
        },
        {
          ao, surprise,
        },
        {
          z, possible,
        },
      },
      []phoneme{
        w, ah, z,
      },
    },
  }

  // WHEN
  actual, ok = searchForBetter(results)
  expected = variantResult{
    0,
    []phonVerdict{
      {
        w, good,
      },
      {
        ao, good,
      },
      {
        z, possible,
      },
    },
    []phoneme{
      w, ao, z,
      },
  }

  // THEN
  if ok != true {
    te.Error()
  }
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and finally test to check that when we only have one verdict and it
  // could be better we return false
  // GIVEN
  results = []variantResult{
    {
      0,
      []phonVerdict{
        {
          w, good,
        },
        {
          ao, good,
        },
        {
          z, missing,
        },
      },
      []phoneme{
        w, ao, z,
      },
    },
  }

  // WHEN
  actual, ok = searchForBetter(results)

  // THEN
  if ok != false {
    te.Error()
  }

  // and GIVEN
  results = []variantResult{
    {
      0,
      []phonVerdict{
        {
          k, good,
        },
        {
          ao, good,
        },
        {
          t, missing,
        },
        {
          ch, surprise,
        },
      },
      []phoneme{
        k, ao, t,
      },
    },
    {
      0,
      []phonVerdict{
        {
          k, good,
        },
        {
          ao, good,
        },
        {
          r, missing,
        },
        {
          t, possible,
        },
      },
      []phoneme{
        k, ao, r, t,
      },
    },
  }

  // WHEN
  actual, ok = searchForBetter(results)
  expected = variantResult{
    0,
    []phonVerdict{
      {
        k, good,
      },
      {
        ao, good,
      },
      {
        t, possible,
      },
    },
    []phoneme{
      k, ao, t,
    },
  }

  // THEN
  if !ok {
    te.Error()
  } else {
    if !reflect.DeepEqual(actual, expected) {
      te.Error()
    }
  }

  // and GIVEN
  results = []variantResult{
    {
      3,
      []phonVerdict{
        {
          hh, good,
        },
        {
          aa, missing,
        },
        {
          ao, surprise,
        },
        {
          d, missing,
        },
        {
          l, good,
        },
        {
          iy, good,
        },
      },
      []phoneme{
        hh, aa, d, l, iy,
      },
    },
    {
      7,
      []phonVerdict{
        {
          hh, good,
        },
        {
          aa, possible,
        },
        {
          r, missing,
        },
        {
          d, possible,
        },
        {
          l, good,
        },
        {
          iy, good,
        },
      },
      []phoneme{
        hh, aa, r, d, l, iy,
      },
    },
  }

  // WHEN
  actual, ok = searchForBetter(results)
  expected = variantResult{
    3,
    []phonVerdict{
      {
        hh, good,
      },
      {
        aa, possible,
      },
      {
        d, possible,
      },
      {
        l, good,
      },
      {
        iy, good,
      },
    },
    []phoneme{
      hh, aa, d, l, iy,
    },
  }

  // THEN
  if !ok {
    te.Error()
  } else {
    if !reflect.DeepEqual(actual, expected) {
      te.Error()
    }
  }
}
*/

func Test_parsePsData(te *testing.T) {
	// GIVEN
	logfile := "test/psOutfile1.txt"

	// WHEN
	actual := parsePsData(logfile)
	expected := []psPhonemeDatum{
		{
			sil,
			0, 10,
		},
		{
			t,
			11, 20,
		},
		{
			eh,
			21, 30,
		},
		{
			s,
			31, 40,
		},
		{
			t,
			41, 50,
		},
		{
			sil,
			51, 60,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	logfile = "test/psOutfile2.txt"

	// WHEN
	actual = parsePsData(logfile)
	expected = []psPhonemeDatum{
		{
			sil,
			0, 11,
		},
		{
			t,
			12, 21,
		},
		{
			eh,
			22, 31,
		},
		{
			s,
			32, 41,
		},
		{
			t,
			42, 52,
		},
		{
			sil,
			52, 61,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

/*
func Test_updateVerdict(te *testing.T) {
  // GIVEN
   rVerd := []phonVerdict{
     {
       ae, good,
     },
     {
       n, missing,
     },
     {
       ih, good,
     },
     {
       m, good,
     },
     {
       l, good,
     },
   }
   taPhons := []psPhonemeDatum{
     {
       ae,
       19, 29,
     },
     {
       ih,
       32, 37,
     },
     {
       l,
       57, 70,
     },
   }

  // WHEN
  actual := updateVerdict(rVerd, taPhons)
  expected := []phonVerdict{
    {
      ae, good,
    },
    {
      n, missing,
    },
    {
      ih, good,
    },
    {
      m, missing,
    },
    {
      l, good,
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  rVerd = []phonVerdict{
    {
      r, good,
    },
    {
      iy, missing,
    },
    {
      aa, surprise,
    },
    {
      p, good,
    },
    {
      l, good,
    },
    {
      ay, good,
    },
    {
      d, good,
    },
  }
  taPhons = []psPhonemeDatum{
    {
      r,
      25, 42,
    },
    {
      aa,
      44, 54,
    },
    {
      p,
      63, 70,
    },
    {
      l,
      73, 78,
    },
    {
      aa,
      79, 101,
    },
    {
      iy,
      101, 109,
    },
    {
      d,
      110, 118,
    },
  }

  // WHEN
  conf := newPsConfig{
    []psPhonemeSettings{},
    "replied",
    [][]phoneme{
      {
        r, iy, p, l, ay, d,
      },
    },
    R_target{},
    neighbourRules,
    "",
  }
  result := psPhonemeResults{
    100,
    taPhons,
  }
  actual = updateVerdict(rVerd, conf.normalise(result, true).data)
  expected = []phonVerdict{
    {
      r, good,
    },
    {
      iy, missing,
    },
    {
      aa, surprise,
    },
    {
      p, good,
    },
    {
      l, good,
    },
    {
      ay, good,
    },
    {
      d, good,
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  rVerd = []phonVerdict{
    {
      w, good,
    },
    {
      ao, good,
    },
    {
      l, good,
    },
    {
      t, good,
    },
    {
      ah, good,
    },
  }
  taPhons = []psPhonemeDatum{
    {
      w,
      22, 42,
    },
    {
      ao,
      42, 51,
    },
    {
      l,
      52, 60,
    },
    {
      t,
      61, 66,
    },
    {
      z,
      70, 77,
    },
    {
      ah,
      78, 92,
    },
  }

  // WHEN
  actual = updateVerdict(rVerd, taPhons)
  expected = []phonVerdict{
    {
      w, good,
    },
    {
      ao, good,
    },
    {
      l, good,
    },
    {
      t, good,
    },
    {
      z, surprise,
    },
    {
      ah, good,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
}
*/

/*
func Test_normalise(te *testing.T) {
  // GIVEN
  conf := newPsConfig{
    []psPhonemeSettings{},
    "tried",
    [][]phoneme{
      {
        t, r, ay, d,
      },
    },
    R_target{},
    neighbourRules,
    "",
  }
  result := psPhonemeResults{
    100,
    []psPhonemeDatum{
      {
        t,
        0, 10,
      },
      {
        r,
        11, 20,
      },
      {
        aa,
        21, 25,
      },
      {
        iy,
        26, 30,
      },
      {
        d,
        31, 40,
      },
    },
  }

  // WHEN
  actual := conf.normalise(result, true)
  expected := psPhonemeResults{
    100,
    []psPhonemeDatum{
      {
        t,
        0, 10,
      },
      {
        r,
        11, 20,
      },
      {
        ay,
        21, 30,
      },
      {
        d,
        31, 40,
      },
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN

  // WHEN
  actual = conf.normalise(result, false)
  expected = psPhonemeResults{
    100,
    []psPhonemeDatum{
      {
        t,
        0, 10,
      },
      {
        r,
        11, 20,
      },
      {
        aa,
        21, 25,
      },
      {
        iy,
        26, 30,
      },
      {
        d,
        31, 40,
      },
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
}
*/

func Test_verdictWithDiphthongs(te *testing.T) {
	// A real example
	// GIVEN
	phons := []phoneme{
		r, iy, p, l, ay, d,
	}
	verdicts := []phonVerdict{
		{
			r, possible,
		},
		{
			iy, good,
		},
		{
			p, good,
		},
		{
			l, missing,
		},
		{
			aa, missing,
		},
		{
			iy, missing,
		},
		{
			ae, surprise,
		},
		{
			d, missing,
		},
		{
			k, surprise,
		},
	}

	// WHEN
	actual := verdictWithDiphthongs(phons, verdicts)
	expected := []phonVerdict{
		{
			r, possible,
		},
		{
			iy, good,
		},
		{
			p, good,
		},
		{
			l, missing,
		},
		{
			ay, missing,
		},
		{
			ae, surprise,
		},
		{
			d, missing,
		},
		{
			k, surprise,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		ay,
	}
	verdicts = []phonVerdict{
		{
			aa, good,
		},
		{
			iy, good,
		},
	}

	// WHEN
	actual = verdictWithDiphthongs(phons, verdicts)
	expected = []phonVerdict{
		{
			ay, good,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		m, ey,
	}
	verdicts = []phonVerdict{
		{
			m, good,
		},
		{
			eh, missing,
		},
		{
			aa, surprise,
		},
		{
			iy, missing,
		},
	}

	// WHEN
	actual = verdictWithDiphthongs(phons, verdicts)
	expected = []phonVerdict{
		{
			m, good,
		},
		{
			ey, missing,
		},
		{
			aa, surprise,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_mapLinkedVerdicts(te *testing.T) {
	// GIVEN
	lv := []linkedPhonVerdict{
		{
			phonVerdict{
				aa, good,
			},
			[]psPhonemeDatumRef{},
		},
		{
			phonVerdict{
				iy, good,
			},
			[]psPhonemeDatumRef{},
		},
	}

	// WHEN
	actual := mapLinkedVerdicts(lv)
	expected := []phonVerdict{
		{
			aa, good,
		},
		{
			iy, good,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_timeNormalise(te *testing.T) {
	// GIVEN
	results := psPhonemeResults{
		100,
		[]psPhonemeDatum{
			{
				sil, 0, 4,
			},
			{
				k, 5, 7,
			},
			{
				ae, 8, 15,
			},
			{
				ae, 16, 19,
			},
			{
				t, 20, 24,
			},
			{
				sil, 25, 29,
			},
			{
				sil, 30, 34,
			},
		},
	}

	// WHEN
	actual := timeNormalise(results)
	expected := psPhonemeResults{
		100,
		[]psPhonemeDatum{
			{
				sil, 0, 4,
			},
			{
				k, 5, 5,
			},
			{
				ae, 8, 15,
			},
			{
				ae, 16, 19,
			},
			{
				t, 20, 24,
			},
			{
				sil, 25, 34,
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_publish(te *testing.T) {
	// GIVEN
	phons := []phoneme{
		ae, n, ih, m, l,
	}
	verdicts := []phonVerdict{
		{
			ae, good,
		},
		{
			n, good,
		},
		{
			ih, good,
		},
		{
			m, good,
		},
		{
			l, good,
		},
	}

	// WHEN
	actual, err := publish("animal", phons, verdicts)
	expected := []LettersVerdict{
		{
			"a", []phoneme{ae}, good,
		},
		{
			"n", []phoneme{n}, good,
		},
		{
			"i", []phoneme{ih}, good,
		},
		{
			"m", []phoneme{m}, good,
		},
		{
			"al", []phoneme{l}, good,
		},
	}

	// THEN
	if err != nil {
		te.Error(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}

	// ...and GIVEN
	phons = []phoneme{
		k, aa, n, t,
	}
	verdicts = []phonVerdict{
		{
			k, good,
		},
		{
			aa, good,
		},
		{
			n, good,
		},
		{
			t, good,
		},
	}

	// WHEN
	actual, err = publish("can't", phons, verdicts)
	expected = []LettersVerdict{
		{
			"c", []phoneme{k}, good,
		},
		{
			"a", []phoneme{aa}, good,
		},
		{
			"n", []phoneme{n}, good,
		},
		{
			"'t", []phoneme{t}, good,
		},
	}

	// THEN
	if err != nil {
		te.Error(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}
}

func TestTrimDupSurprises(te *testing.T) {
	// GIVEN
	v := []phonVerdict{
		{
			t, good,
		},
		{
			eh, good,
		},
		{
			s, good,
		},
		{
			t, good,
		},
		{
			t, surprise,
		},
	}

	// WHEN
	actual := trimDupSurprises(v)
	expected := []phonVerdict{
		{
			t, good,
		},
		{
			eh, good,
		},
		{
			s, good,
		},
		{
			t, good,
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}

	// ...and GIVEN
	v = []phonVerdict{
		{
			t, good,
		},
		{
			eh, good,
		},
		{
			t, surprise,
		},
		{
			s, good,
		},
		{
			t, good,
		},
	}

	// WHEN
	actual = trimDupSurprises(v)

	// THEN
	if !reflect.DeepEqual(actual, v) {
		te.Error(actual)
	}
}
