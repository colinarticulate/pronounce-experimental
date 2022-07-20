package pron

import (
	"testing"
)

func Test_new_testGrammarFile(t *testing.T) {
	// GIVEN
	g1 := new_jsgfStandard()

	// WHEN
	c := new_testGrammarFile(&g1)
	actual := c.config

	// THEN
	if _, ok := actual.(*jsgfStandard); !ok {
		t.Error()
	}

	// and GIVEN
	g2 := new_jsgfDiphthong()

	// WHEN
	c = new_testGrammarFile(&g2)
	actual = c.config

	// THEN
	if _, ok := actual.(*jsgfDiphthong); !ok {
		t.Error()
	}
}

type wrapper struct {
	j jsgfGrammar
}

func new(config jsgfConfig) wrapper {
	return wrapper{
		new_testGrammarFile(config),
	}
}

func Test_wrapper(t *testing.T) {
	// GIVEN
	g1 := new_jsgfStandard()
	w := new(&g1)

	// WHEN
	actual := w.j.config

	// THEN
	if _, ok := actual.(*jsgfStandard); !ok {
		t.Error()
	}
}
