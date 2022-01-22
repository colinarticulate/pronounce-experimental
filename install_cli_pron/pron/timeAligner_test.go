package pron

import (
	"log"
	"reflect"
	"testing"
)

func Test_AddResult(t *testing.T) {
	// GIVEN
	phons := []phoneme{
		dh, iy,
	}
	ta := newTimeAligner(phons)

	// WHEN
	result := psPhonemeResults{
		90,
		[]psPhonemeDatum{
			{
				sil,
				0, 22,
			},
			{
				z,
				23, 44,
			},
			{
				ah,
				46, 77,
			},
			{
				dh,
				78, 83,
			},
			{
				sil,
				88, 128,
			},
		},
	}
	ta.AddResult(result)

	// THEN
	if len(ta.quantisedResult) != 129 {
		t.Error()
	}
	links, ok := ta.quantisedResult[50][ah]
	if !ok {
		t.Error()
	}
	if len(links) != 1 {
		t.Error()
	}
	links, ok = ta.quantisedResult[77][ah]
	if !ok {
		t.Error()
	}
	if len(links) != 1 {
		t.Error()
	}
	links, ok = ta.quantisedResult[78][dh]
	if !ok {
		t.Error()
	}
	if len(links) != 1 {
		t.Error()
	}

	// and WHEN
	result = psPhonemeResults{
		120,
		[]psPhonemeDatum{
			{
				sil,
				0, 23,
			},
			{
				s,
				23, 44,
			},
			{
				ah,
				45, 65,
			},
			{
				dh,
				66, 68,
			},
			{
				sil,
				72, 96,
			},
		},
	}
	ta.AddResult(result)

	// THEN
	if len(ta.quantisedResult) != 129 {
		t.Error()
	}
	links, ok = ta.quantisedResult[50][ah]
	if !ok {
		t.Error()
	}
	if len(links) != 2 {
		t.Error()
	}
	expected := quantumResult{
		s: []psPhonemeDatumRef{
			{
				1, 1,
			},
		},
		z: []psPhonemeDatumRef{
			{
				0, 1,
			},
		},
	}
	if !reflect.DeepEqual(ta.quantisedResult[25], expected) {
		t.Error()
	}
	expected = quantumResult{
		ah: []psPhonemeDatumRef{
			{
				1, 2,
			},
		},
	}
	if !reflect.DeepEqual(ta.quantisedResult[45], expected) {
		t.Error()
	}
	expected = quantumResult{
		ah: []psPhonemeDatumRef{
			{
				0, 2,
			},
			{
				1, 2,
			},
		},
	}
	if !reflect.DeepEqual(ta.quantisedResult[46], expected) {
		t.Error()
	}
}

func Test_timeAlign(te *testing.T) {
	// GIVEN
	phons := []phoneme{
		dh, iy,
	}
	ta := newTimeAligner(phons)
	results := []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					sil,
					0, 22,
				},
				{
					z,
					23, 44,
				},
				{
					ah,
					46, 77,
				},
				{
					dh,
					78, 83,
				},
				{
					sil,
					88, 128,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					sil,
					12, 23,
				},
				{
					dh,
					24, 50,
				},
				{
					iy,
					50, 67,
				},
				{
					sil,
					68, 105,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					sil,
					0, 23,
				},
				{
					s,
					23, 44,
				},
				{
					ah,
					45, 65,
				},
				{
					dh,
					66, 68,
				},
				{
					sil,
					72, 96,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					sil,
					0, 21,
				},
				{
					z,
					24, 44,
				},
				{
					ah,
					45, 75,
				},
				{
					sil,
					76, 129,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					sil,
					0, 21,
				},
				{
					m,
					22, 24,
				},
				{
					z,
					24, 45,
				},
				{
					ah,
					45, 64,
				},
				{
					sil,
					66, 79,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}

	// WHEN
	actual := ta.timeAlign()
	expected := []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				sil,
				0, 22,
			},
			[]psPhonemeDatumRef{
				{
					0, 0,
				},
				{
					2, 0,
				},
				{
					3, 0,
				},
				{
					4, 0,
				},
				{
					1, 0,
				},
			},
		},
		{
			psPhonemeDatum{
				z,
				24, 44,
			},
			[]psPhonemeDatumRef{
				{
					0, 1,
				},
				{
					3, 1,
				},
				{
					4, 2,
				},
			},
		},
		{
			psPhonemeDatum{
				ah,
				45, 65,
			},
			[]psPhonemeDatumRef{
				{
					2, 2,
				},
				{
					3, 2,
				},
				{
					4, 3,
				},
				{
					0, 2,
				},
			},
		},
		{
			psPhonemeDatum{
				sil,
				72, 105,
			},
			[]psPhonemeDatumRef{
				{
					1, 3,
				},
				{
					2, 4,
				},
				{
					4, 4,
				},
				{
					3, 3,
				},
				{
					0, 4,
				},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		k, ao, t,
	}
	ta = newTimeAligner(phons)
	results = []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					sil,
					0, 17,
				},
				{
					k,
					18, 26,
				},
				{
					ao,
					27, 47,
				},
				{
					t,
					48, 68,
				},
				{
					sil,
					69, 122,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					sil,
					0, 18,
				},
				{
					k,
					19, 26,
				},
				{
					ao,
					27, 48,
				},
				{
					t,
					49, 68,
				},
				{
					sil,
					69, 108,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					sil,
					0, 18,
				},
				{
					k,
					19, 27,
				},
				{
					ao,
					28, 48,
				},
				{
					t,
					48, 68,
				},
				{
					sil,
					68, 102,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					sil,
					0, 18,
				},
				{
					k,
					19, 27,
				},
				{
					ao,
					27, 48,
				},
				{
					t,
					49, 65,
				},
				{
					sil,
					66, 120,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					sil,
					18, 22,
				},
				{
					k,
					22, 29,
				},
				{
					ao,
					30, 51,
				},
				{
					t,
					51, 68,
				},
				{
					sil,
					68, 87,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}

	// WHEN
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				sil,
				0, 18,
			},
			[]psPhonemeDatumRef{
				{
					0, 0,
				},
				{
					1, 0,
				},
				{
					2, 0,
				},
				{
					3, 0,
				},
				{
					4, 0,
				},
			},
		},
		{
			psPhonemeDatum{
				k,
				19, 27,
			},
			[]psPhonemeDatumRef{
				{
					0, 1,
				},
				{
					1, 1,
				},
				{
					2, 1,
				},
				{
					3, 1,
				},
				{
					4, 1,
				},
			},
		},
		{
			psPhonemeDatum{
				ao,
				27, 48,
			},
			[]psPhonemeDatumRef{
				{
					0, 2,
				},
				{
					1, 2,
				},
				{
					3, 2,
				},
				{
					2, 2,
				},
				{
					4, 2,
				},
			},
		},
		{
			psPhonemeDatum{
				t,
				49, 68,
			},
			[]psPhonemeDatumRef{
				{
					0, 3,
				},
				{
					1, 3,
				},
				{
					2, 3,
				},
				{
					3, 3,
				},
				{
					4, 3,
				},
			},
		},
		{
			psPhonemeDatum{
				sil,
				68, 108,
			},
			[]psPhonemeDatumRef{
				{
					2, 4,
				},
				{
					3, 4,
				},
				{
					4, 4,
				},
				{
					0, 4,
				},
				{
					1, 4,
				},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		ae,
	}
	ta = newTimeAligner(phons)
	results = []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					ae,
					31, 49,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					ae,
					35, 60,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					ae,
					40, 68,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					ae,
					48, 87,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					ae,
					65, 109,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}

	// WHEN
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				ae,
				40, 60,
			},
			[]psPhonemeDatumRef{
				{
					0, 0,
				},
				{
					1, 0,
				},
				{
					2, 0,
				},
				{
					3, 0,
				},
			},
		},
		{
			psPhonemeDatum{
				ae,
				65, 68,
			},
			[]psPhonemeDatumRef{
				{
					2, 0,
				},
				{
					3, 0,
				},
				{
					4, 0,
				},
			},
		},
	}

	// THEN
	log.Println(actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		k, l, ay, m,
	}
	results = []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					k,
					14, 26,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					k,
					15, 27,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					k,
					17, 28,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					k,
					15, 25,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					k,
					19, 19,
				},
			},
		},
	}
	ta = newTimeAligner(phons)

	// WHEN
	for _, result := range results {
		ta.AddResult(result)
	}
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				k,
				15, 26,
			},
			[]psPhonemeDatumRef{
				{0, 0}, {1, 0}, {3, 0}, {2, 0}, {4, 0},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual[0].psPhonemeDatum, expected[0].psPhonemeDatum) {
		te.Error()
	}
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		ao, dh, ah,
	}
	results = []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					ah,
					59, 71,
				},
				{
					ah,
					76, 79,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					ah,
					30, 40,
				},
				{
					ah,
					60, 76,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					ah,
					62, 81,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					ah,
					59, 83,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					ah,
					61, 74,
				},
				{
					ah,
					76, 79,
				},
			},
		},
	}
	ta = newTimeAligner(phons)

	// WHEN
	for _, result := range results {
		ta.AddResult(result)
	}
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				ah,
				60, 79,
			},
			[]psPhonemeDatumRef{
				{0, 0}, {1, 1}, {3, 0}, {4, 0}, {2, 0}, {0, 1}, {4, 1},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		w, ao, t, ah,
	}
	ta = newTimeAligner(phons)

	// WHEN
	results = []psPhonemeResults{
		{
			105,
			[]psPhonemeDatum{
				{
					ao,
					59, 68,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					ao,
					61, 68,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					ao,
					59, 67,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				ao,
				61, 67,
			},
			[]psPhonemeDatumRef{
				{0, 0}, {1, 0}, {2, 0},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		w, ao, t, ah,
	}
	ta = newTimeAligner(phons)

	// WHEN
	results = []psPhonemeResults{
		{
			90,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					hh,
					4, 4,
				},
				{
					s,
					49, 49,
				},
				{
					w,
					52, 58,
				},
				{
					ah,
					59, 68,
				},
				{
					dh,
					69, 69,
				},
				{
					t,
					72, 78,
				},
				{
					ah,
					79, 87,
				},
				{
					d,
					88, 97,
				},
				{
					g,
					98, 159,
				},
				{
					sil,
					160, 160,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					w,
					53, 58,
				},
				{
					ao,
					59, 68,
				},
				{
					t,
					69, 79,
				},
				{
					ah,
					80, 89,
				},
				{
					sil,
					90, 90,
				},
			},
		},
		{
			120,
			[]psPhonemeDatum{
				{
					sil,
					11, 11,
				},
				{
					w,
					55, 60,
				},
				{
					ao,
					61, 68,
				},
				{
					t,
					69, 81,
				},
				{
					ah,
					82, 91,
				},
				{
					sil,
					92, 92,
				},
			},
		},
		{
			143,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					w,
					53, 53,
				},
				{
					aa,
					57, 67,
				},
				{
					dh,
					68, 68,
				},
				{
					t,
					71, 81,
				},
				{
					ah,
					82, 87,
				},
				{
					sil,
					88, 88,
				},
			},
		},
		{
			190,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					hh,
					7, 47,
				},
				{
					sil,
					48, 48,
				},
				{
					w,
					56, 56,
				},
				{
					ao,
					59, 67,
				},
				{
					t,
					68, 79,
				},
				{
					ah,
					80, 87,
				},
				{
					t,
					87, 101,
				},
				{
					iy,
					101, 127,
				},
				{
					sil,
					128, 128,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}
	actual = ta.timeAlign()
	expected = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				ao,
				61, 67,
			},
			[]psPhonemeDatumRef{
				{1, 2}, {2, 2}, {4, 4},
			},
		},
	}

	// THEN
	log.Println(actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		r, ah, sh, ah,
	}
	ta = newTimeAligner(phons)
	results = []psPhonemeResults{
		{
			72,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					b,
					14, 14,
				},
				{
					ih,
					18, 32,
				},
				{
					l,
					33, 33,
				},
				{
					ih,
					38, 50,
				},
				{
					ng,
					51, 78,
				},
				{
					sil,
					79, 79,
				},
			},
		},
		{
			80,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					b,
					14, 18,
				},
				{
					ih,
					19, 33,
				},
				{
					l,
					34, 34,
				},
				{
					ih,
					38, 51,
				},
				{
					ng,
					53, 80,
				},
				{
					sil,
					81, 81,
				},
			},
		},
		{
			91,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					b,
					14, 18,
				},
				{
					ih,
					19, 33,
				},
				{
					l,
					34, 37,
				},
				{
					ih,
					38, 52,
				},
				{
					ng,
					53, 80,
				},
				{
					sil,
					81, 81,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					b,
					15, 18,
				},
				{
					ih,
					19, 34,
				},
				{
					l,
					35, 38,
				},
				{
					ih,
					39, 53,
				},
				{
					ng,
					54, 81,
				},
				{
					sil,
					82, 82,
				},
			},
		},
		{
			125,
			[]psPhonemeDatum{
				{
					sil,
					3, 3,
				},
				{
					b,
					14, 19,
				},
				{
					ih,
					20, 38,
				},
				{
					l,
					39, 39,
				},
				{
					ih,
					42, 53,
				},
				{
					ng,
					54, 84,
				},
				{
					sil,
					85, 85,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}

	// WHEN
	actual = ta.timeAlign()

	// THEN
	log.Println("billing, actual =", actual)

	// and GIVEN
	phons = []phoneme{
		r, ah, sh, ah,
	}
	ta = newTimeAligner(phons)
	results = []psPhonemeResults{
		{
			72,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					r,
					24, 43,
				},
				{
					ah,
					44, 49,
				},
				{
					ch,
					50, 63,
				},
				{
					y,
					64, 75,
				},
				{
					ah,
					76, 89,
				},
				{
					sil,
					90, 90,
				},
			},
		},
		{
			80,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					r,
					24, 44,
				},
				{
					ah,
					45, 45,
				},
				{
					sh,
					49, 66,
				},
				{
					ah,
					68, 89,
				},
				{
					sil,
					90, 90,
				},
			},
		},
		{
			91,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					r,
					20, 20,
				},
				{
					l,
					26, 26,
				},
				{
					r,
					30, 43,
				},
				{
					ah,
					44, 47,
				},
				{
					sh,
					48, 65,
				},
				{
					ah,
					66, 74,
				},
				{
					ae,
					75, 82,
				},
				{
					n,
					84, 103,
				},
				{
					sil,
					104, 104,
				},
			},
		},
		{
			105,
			[]psPhonemeDatum{
				{
					sil,
					0, 0,
				},
				{
					v,
					20, 20,
				},
				{
					n,
					24, 24,
				},
				{
					r,
					27, 44,
				},
				{
					ah,
					45, 48,
				},
				{
					sh,
					49, 66,
				},
				{
					ah,
					67, 89,
				},
				{
					sil,
					90, 90,
				},
			},
		},
		{
			125,
			[]psPhonemeDatum{
				{
					sil,
					3, 3,
				},
				{
					b,
					22, 22,
				},
				{
					l,
					27, 31,
				},
				{
					r,
					32, 45,
				},
				{
					ah,
					46, 48,
				},
				{
					sh,
					49, 65,
				},
				{
					dh,
					66, 66,
				},
				{
					ih,
					68, 83,
				},
				{
					m,
					84, 102,
				},
				{
					w,
					103, 103,
				},
				{
					sil,
					106, 106,
				},
			},
		},
	}
	for _, result := range results {
		ta.AddResult(result)
	}

	// WHEN
	actual = ta.timeAlign()

	// THEN
	log.Println("russia, actual =", actual)

}
