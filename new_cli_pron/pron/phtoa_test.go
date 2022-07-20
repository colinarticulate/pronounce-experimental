package pron

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/colinarticulate/dictionary"
)

func Test_wordListInApp(te *testing.T) {
	// GIVEN
	alphas := "animal"
	phons := []phoneme{
		ae, n, ih, m, ah, l,
	}

	// WHEN
	actual, _ := mapPhToA(phons, alphas)
	expected := []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				ih,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"m",
		},
		{
			[]phoneme{
				ah,
			},
			"a",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
	}

	// THEN

	// and GIVEN
	phons = []phoneme{
		ae, n, er, m, ah, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				er,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"m",
		},
		{
			[]phoneme{
				ah,
			},
			"a",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		ae, n, er, m, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				er,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"m",
		},
		{
			[]phoneme{
				l,
			},
			"al",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "because"
	phons = []phoneme{
		b, ih, k, ao, z,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				b,
			},
			"b",
		},
		{
			[]phoneme{
				ih,
			},
			"e",
		},
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				ao,
			},
			"au",
		},
		{
			[]phoneme{
				z,
			},
			"se",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		b, iy, k, ao, z,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				b,
			},
			"b",
		},
		{
			[]phoneme{
				iy,
			},
			"e",
		},
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				ao,
			},
			"au",
		},
		{
			[]phoneme{
				z,
			},
			"se",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "climbed"
	phons = []phoneme{
		k, l, ay, m, d,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				ay,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"mbe",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "caught"
	phons = []phoneme{
		k, aa, t,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				aa,
			},
			"au",
		},
		{
			[]phoneme{
				t,
			},
			"ght",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "coat"
	phons = []phoneme{
		k, ow, t,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				ow,
			},
			"oa",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "conversation"
	phons = []phoneme{
		k, aa, n, v, er, s, ey, sh, ao, n,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				aa,
			},
			"o",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				v,
			},
			"v",
		},
		{
			[]phoneme{
				er,
			},
			"er",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				ey,
			},
			"a",
		},
		{
			[]phoneme{
				sh,
			},
			"ti",
		},
		{
			[]phoneme{
				ao,
			},
			"o",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "greatest"
	phons = []phoneme{
		g, r, ey, t, ah, s, t,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				g,
			},
			"g",
		},
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ey,
			},
			"ea",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
		{
			[]phoneme{
				ah,
			},
			"e",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "hardly"
	phons = []phoneme{
		hh, aa, r, d, l, iy,
	}
	// WHEN
	actual, err := mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				hh,
			},
			"h",
		},
		{
			[]phoneme{
				aa, r,
			},
			"ar",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				iy,
			},
			"y",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		log.Println("hardly = ", actual, "err =", err)
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		hh, aa, d, l, iy,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				hh,
			},
			"h",
		},
		{
			[]phoneme{
				aa,
			},
			"ar",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				iy,
			},
			"y",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "last"
	phons = []phoneme{
		l, ae, s, t,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "last"
	phons = []phoneme{
		l, aa, s, t,
	}
	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				aa,
			},
			"a",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "other"
	phons = []phoneme{
		ah, dh, er,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ah,
			},
			"o",
		},
		{
			[]phoneme{
				dh,
			},
			"th",
		},
		{
			[]phoneme{
				er,
			},
			"er",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "replied"
	phons = []phoneme{
		r, ih, p, l, ay, d,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ih,
			},
			"e",
		},
		{
			[]phoneme{
				p,
			},
			"p",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				ay,
			},
			"ie",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	phons = []phoneme{
		r, iy, p, l, ay, d,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				iy,
			},
			"e",
		},
		{
			[]phoneme{
				p,
			},
			"p",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				ay,
			},
			"ie",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "russia"
	phons = []phoneme{
		r, ah, sh, ah,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ah,
			},
			"u",
		},
		{
			[]phoneme{
				sh,
			},
			"ss",
		},
		{
			[]phoneme{
				ah,
			},
			"ia",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "school"
	phons = []phoneme{
		s, k, uw, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				k,
			},
			"ch",
		},
		{
			[]phoneme{
				uw,
			},
			"oo",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "seat"
	phons = []phoneme{
		s, iy, t,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				iy,
			},
			"ea",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "sit"
	phons = []phoneme{
		s, ih, t,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				ih,
			},
			"i",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "squirrel"
	phons = []phoneme{
		s, k, w, er, ah, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				k, w,
			},
			"qu",
		},
		{
			[]phoneme{
				er,
			},
			"irr",
		},
		{
			[]phoneme{
				ah,
			},
			"e",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "time"
	phons = []phoneme{
		t, ay, m,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				t,
			},
			"t",
		},
		{
			[]phoneme{
				ay,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"me",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "water"
	phons = []phoneme{
		w, ao, t, er,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				w,
			},
			"w",
		},
		{
			[]phoneme{
				ao,
			},
			"a",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
		{
			[]phoneme{
				er,
			},
			"er",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "world"
	phons = []phoneme{
		w, er, l, d,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				w,
			},
			"w",
		},
		{
			[]phoneme{
				er,
			},
			"or",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_mapPhToA(te *testing.T) {
	// GIVEN
	alphas := "our"
	phons := []phoneme{
		aw, er,
	}

	// WHEN
	actual, _ := mapPhToA(phons, alphas)
	expected := []phonToAlphas{
		{
			[]phoneme{
				aw,
			},
			"ou",
		},
		{
			[]phoneme{
				er,
			},
			"r",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "vehicle"
	phons = []phoneme{
		v, iy, ih, k, ah, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				v,
			},
			"v",
		},
		{
			[]phoneme{
				iy,
			},
			"e",
		},
		{
			[]phoneme{
				ih,
			},
			"hi",
		},
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				ah, l,
			},
			"le",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "hardly"
	phons = []phoneme{
		hh, aa, d, l, iy,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				hh,
			},
			"h",
		},
		{
			[]phoneme{
				aa,
			},
			"ar",
		},
		{
			[]phoneme{
				d,
			},
			"d",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
		{
			[]phoneme{
				iy,
			},
			"y",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "exciting"
	phons = []phoneme{
		ih, k, s, ay, t, ih, ng,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ih,
			},
			"e",
		},
		{
			[]phoneme{
				k, s,
			},
			"xc",
		},
		{
			[]phoneme{
				ay,
			},
			"i",
		},
		{
			[]phoneme{
				t,
			},
			"t",
		},
		{
			[]phoneme{
				ih,
			},
			"i",
		},
		{
			[]phoneme{
				ng,
			},
			"ng",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "asked"
	phons = []phoneme{
		ae, s, k, t,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				k,
			},
			"ke",
		},
		{
			[]phoneme{
				t,
			},
			"d",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "treasuries"
	phons = []phoneme{
		t, r, eh, zh, er, iy, z,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				t,
			},
			"t",
		},
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				eh,
			},
			"ea",
		},
		{
			[]phoneme{
				zh,
			},
			"s",
		},
		{
			[]phoneme{
				er,
			},
			"ur",
		},
		{
			[]phoneme{
				iy,
			},
			"ie",
		},
		{
			[]phoneme{
				z,
			},
			"s",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "provisionally"
	phons = []phoneme{
		p, r, ah, v, ih, zh, ah, n, ah, l, iy,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				p,
			},
			"p",
		},
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ah,
			},
			"o",
		},
		{
			[]phoneme{
				v,
			},
			"v",
		},
		{
			[]phoneme{
				ih,
			},
			"i",
		},
		{
			[]phoneme{
				zh,
			},
			"s",
		},
		{
			[]phoneme{
				ah,
			},
			"io",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				ah,
			},
			"a",
		},
		{
			[]phoneme{
				l,
			},
			"ll",
		},
		{
			[]phoneme{
				iy,
			},
			"y",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}

	// and GIVEN
	alphas = "animal"
	phons = []phoneme{
		ae, n, er, m, ah, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				er,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"m",
		},
		{
			[]phoneme{
				ah,
			},
			"a",
		},
		{
			[]phoneme{
				l,
			},
			"l",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "animal"
	phons = []phoneme{
		ae, n, er, m, l,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				er,
			},
			"i",
		},
		{
			[]phoneme{
				m,
			},
			"m",
		},
		{
			[]phoneme{
				l,
			},
			"al",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "usurp"
	phons = []phoneme{
		y, uw, s, ah, p,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				y, uw,
			},
			"u",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				ah,
			},
			"ur",
		},
		{
			[]phoneme{
				p,
			},
			"p",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "conversation"
	phons = []phoneme{
		k, aa, n, v, s, ey, sh, n,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				k,
			},
			"c",
		},
		{
			[]phoneme{
				aa,
			},
			"o",
		},
		{
			[]phoneme{
				n,
			},
			"n",
		},
		{
			[]phoneme{
				v,
			},
			"ver",
		},
		{
			[]phoneme{
				s,
			},
			"s",
		},
		{
			[]phoneme{
				ey,
			},
			"a",
		},
		{
			[]phoneme{
				sh, n,
			},
			"tion",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		log.Println("conversation, actual =", actual)
		te.Error()
	}

	// and GIVEN
	alphas = "chair"
	phons = []phoneme{
		ch, eh,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ch,
			},
			"ch",
		},
		{
			[]phoneme{
				eh,
			},
			"air",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}

	// and GIVEN
	alphas = "russia"
	phons = []phoneme{
		r, ah, sh, ah,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ah,
			},
			"u",
		},
		{
			[]phoneme{
				sh,
			},
			"ssi",
		},
		{
			[]phoneme{
				ah,
			},
			"a",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		log.Println("russia, actual =", actual)
		te.Error()
	}

	// and GIVEN
	alphas = "fir"
	phons = []phoneme{
		f, er,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				f,
			},
			"f",
		},
		{
			[]phoneme{
				er,
			},
			"ir",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}

	// and GIVEN
	alphas = "rapport"
	phons = []phoneme{
		r, ae, p, ao, r,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				r,
			},
			"r",
		},
		{
			[]phoneme{
				ae,
			},
			"a",
		},
		{
			[]phoneme{
				p,
			},
			"pp",
		},
		{
			[]phoneme{
				ao, r,
			},
			"ort",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}

	// ...and GIVEN
	alphas = "aide-de-camp"
	phons = []phoneme{
		ey, d, d, ax, k, oh, m,
	}

	// WHEN
	actual, _ = mapPhToA(phons, alphas)
	expected = []phonToAlphas{
		{
			[]phoneme{
				ey,
			},
			"ai",
		},
		{
			[]phoneme{
				d,
			},
			"de",
		},
		{
			[]phoneme{
				d,
			},
			"-d",
		},
		{
			[]phoneme{
				ax,
			},
			"e",
		},
		{
			[]phoneme{
				k,
			},
			"-c",
		},
		{
			[]phoneme{
				oh,
			},
			"a",
		},
		{
			[]phoneme{
				m,
			},
			"mp",
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error(actual)
	}
}

func Test_trailingSilentAlphas(te *testing.T) {
	//GIVEN
	abc := "rapport"
	phons := []phoneme{
		r, ae, p, ao, r,
	}
	currMap := phonToAlphas{
		[]phoneme{
			ao, r,
		},
		"or",
	}

	// WHEN
	actual := trailingSilentAlphas(abc, len("rapp")-1, phons, len([]phoneme{r, ae, p})-1, currMap)
	expected := "t"

	// THEN
	if actual != expected {
		te.Error(actual)
	}

	// and GIVEN
	currMap = phonToAlphas{
		[]phoneme{
			ao, r,
		},
		"ort",
	}

	// WHEN
	actual = trailingSilentAlphas(abc, len("rapp")-1, phons, len([]phoneme{r, ae, p})-1, currMap)
	expected = ""

	// THEN
	if actual != expected {
		te.Error(actual)
	}
}

func uniqueWordsFromFile(filename string) []string {
	textBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}
	words := strings.Fields(string(textBytes))
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range words {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		words[j] = v
		j++
	}
	return words[:j]
}

func reallyPhonemes(word string, dict dictionary.Dictionary) ([]phoneme, error) {
	phs := []phoneme{}
	phStrs, err := dict.Phonemes(word)
	if err != nil {
		return phs, err
	}
	for _, phStr := range phStrs {
		phs = append(phs, phoneme(phStr))
	}
	return phs, nil
}

func convStrsToPhons(strs []string) []phoneme {
	phs := []phoneme{}
	for _, str := range strs {
		phs = append(phs, phoneme(str))
	}
	return phs
}

func checkMap(word string, mapping []phonToAlphas) bool {
	alphas := ""
	for _, m := range mapping {
		alphas += m.alphas
	}
	if word != alphas {
		return false
	}
	return true
}

func Test_bigmapPhToA(t *testing.T) {
	// GIVEN
	/*
	  words := uniqueWordsFromFile("../../test/alice.txt")
	  dict := dictionary.Create("../../test/cmudict-0.7b.txt")

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }

	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/goldilocks.txt")

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }

	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/words.txt")

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }

	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/dracula.txt")

	  log.Println("dracula...")
	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }
	*/

	/*
	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/windinthewillows.txt")
	  log.Println("len(words) in windinthewillows.tx =", len(words))

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }
	*/

	/*
	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/frankenstein.txt")
	  dict = dictionary.Create("../../test/frankensteindict.txt")

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }
	*/

	/*
	  // and GIVEN
	  words = uniqueWordsFromFile("../../test/top10000.txt")

	  for _, word := range words {
	    phons, err := reallyPhonemes(word, dict)
	    if err != nil {
	      log.Println("Info: Word,", word, " not in dictionary?")
	      continue
	    }
	    // WHEN
	    _, err = mapPhToA(phons, word)

	    // THEN
	    if err != nil {
	      t.Error(err)
	    }
	  }
	*/
	// GIVEN
	// dict := dictionary.Create("../../test/beepFiltered_er_axN_axL_uwL_StS.dict")
	dict := dictionary.Create("../../WorkingDictionary/wordDictionary/Out/sourceFiltered.dict")
	f, err := os.Create("test/mapPhToAFailures.txt")
	if err != nil {
		// Stop the test now
		log.Panic(err)
	}
	defer f.Close()

	for _, entry := range dict.AllEntries() {
		word := strings.ToLower(entry.Word())
		strs := entry.Phonemes()
		phons := convStrsToPhons(entry.Phonemes())
		// WHEN
		m, err := mapPhToA(phons, word)

		// THEN
		if !checkMap(word, m) {
			f.WriteString("Oops, bad map = " + word + ": " + fmt.Sprintf("%v", m) + "\n")
		}
		if err != nil {
			f.WriteString(word + ", " + strings.Join(strs, " ") + "\n")
			f.WriteString(err.Error() + "\n\n")
			// t.Error(err)
		}
	}
}
