package pron

import (
	"fmt"
	"reflect"
	"testing"
)

/*
func Test_integrate(te *testing.T) {
  // GIVEN
  phons := []phoneme{
    eh, m, p, er, er, z,
  }
  c := newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        eh, good,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
      },
    },
    {
      phonVerdict{
        m, good,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
      },
    },
    {
      phonVerdict{
        p, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
      },
    },
    {
      phonVerdict{
        er, good,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
      },
    },
    {
      phonVerdict{
        er, good,
      },
      []psPhonemeDatumRef{
        {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5},
      },
    },
    {
      phonVerdict{
        z, good,
      },
      []psPhonemeDatumRef{
        {0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        eh,
        19, 33,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {4, 1}, {2, 1}, {3, 1},
      },
    },
    {
      psPhonemeDatum{
        m,
        34, 45,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
      },
    },
    {
      psPhonemeDatum{
        p,
        46, 56,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
      },
    },
    {
      psPhonemeDatum{
        er,
        58, 89,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 5}, {2, 5}, {1, 5}, {3, 5}, {0, 5},
      },
    },
    {
      psPhonemeDatum{
        z,
        90, 111,
      },
      []psPhonemeDatumRef{
        {0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6},
      },
    },
  }

  // WHEN
  actual := c.integrate()
  expected := []phonVerdict{
    {
      eh, good,
    },
    {
      m, good,
    },
    {
      p, good,
    },
    {
      er, good,
    },
    {
      er, good,
    },
    {
      z, good,
    },
  }

  if !reflect.DeepEqual(actual, expected) {
    fmt.Println("emperors actual =", actual)
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    r, iy, p, l, ay, d,
  }

  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        r, possible,
      },
      []psPhonemeDatumRef{
        {
          0, 1,
        },
        {
          2, 1,
        },
        {
          3, 1,
        },
      },
    },
    {
      phonVerdict{
        iy, good,
      },
      []psPhonemeDatumRef{
        {
          0, 2,
        },
        {
          1, 2,
        },
        {
          2, 2,
        },
        {
          3, 2,
        },
        {
          4, 2,
        },
      },
    },
    {
      phonVerdict{
        p, good,
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
      phonVerdict{
        l, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        aa, missing,
      },
      []psPhonemeDatumRef{
        {
        },
      },
    },
    {
      phonVerdict{
        iy, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        ae, surprise,
      },
      []psPhonemeDatumRef{
        {
          0, 6,
        },
        {
          1, 5,
        },
        {
          2, 5,
        },
        {
          3, 5,
        },
        {
          4, 6,
        },
      },
    },
    {
      phonVerdict{
        d, missing,
      },
      []psPhonemeDatumRef{
        {
          1, 7,
        },
        {
          2, 7,
        },
      },
    },
    {
      phonVerdict{
        k, surprise,
      },
      []psPhonemeDatumRef{
        {
          0, 9,
        },
        {
          3, 8,
        },
        {
          4, 8,
        },
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        sil,
        0, 3,
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
        r,
        3, 6,
      },
      []psPhonemeDatumRef{
        {
          0, 1,
        },
        {
          2, 1,
        },
        {
          3, 1,
        },
        {
          1, 1,
        },
        {
          4, 1,
        },
      },
    },
    {
      psPhonemeDatum{
        iy,
        7, 17,
      },
      []psPhonemeDatumRef{
        {
          0, 2,
        },
        {
          1, 2,
        },
        {
          2, 2,
        },
        {
          3, 2,
        },
        {
          4, 2,
        },
      },
    },
    {
      psPhonemeDatum{
        p,
        18, 30,
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
          4, 3,
        },
        {
          3, 3,
        },
      },
    },
    {
      psPhonemeDatum{
        l,
        31, 32,
      },
      []psPhonemeDatumRef{
        {
          1, 4,
        },
        {
          2, 4,
        },
        {
          3, 4,
        },
        {
          4, 4,
        },
      },
    },
    {
      psPhonemeDatum{
        ae,
        33, 54,
      },
      []psPhonemeDatumRef{
        {
          0, 6,
        },
        {
          1, 5,
        },
        {
          2, 5,
        },
        {
          4, 6,
        },
        {
          3, 5,
        },
      },
    },
    {
      psPhonemeDatum{
        ng,
        57, 61,
      },
      []psPhonemeDatumRef{
        {
          0, 8,
        },
        {
          3, 7,
        },
        {
          4, 7,
        },
      },
    },
    {
      psPhonemeDatum{
        k,
        63, 80,
      },
      []psPhonemeDatumRef{
        {
          0, 9,
        },
        {
          3, 8,
        },
        {
          4, 8,
        },
      },
    },
    {
      psPhonemeDatum{
        sil,
        83, 112,
      },
      []psPhonemeDatumRef{
        {
          0, 10,
        },
        {
          3, 9,
        },
        {
          4, 9,
        },
        {
          1, 8,
        },
        {
          2, 8,
        },
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      r, possible,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    hh, ey,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        hh, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        l, surprise,
      },
      []psPhonemeDatumRef{
        {2, 4}, {3, 2}, {4, 2},
      },
    },
    {
      phonVerdict{
        ey, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        sil,
        0, 15,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        w,
        25, 32,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        l,
        33, 39,
      },
      []psPhonemeDatumRef{
        {2, 4}, {3, 2}, {4, 2},
      },
    },
    {
      psPhonemeDatum{
        ey,
        41, 75,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        sil,
        80, 141,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      hh, missing,
    },
    {
      w, surprise,
    },
    {
      l, surprise,
    },
    {
      ey, good,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    k, ih, t,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        k, missing,
      },
      []psPhonemeDatumRef{
        {},
      } ,
    },
    {
      phonVerdict{
        ih, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {3, 4}, {4, 3},
      },
    },
    {
      phonVerdict{
        t, possible,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        s,
        8, 21,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        ih,
        23, 24,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 2}, {3, 4}, {4, 3},
      },
    },
    {
      psPhonemeDatum{
        t,
        36, 48,
      },
      []psPhonemeDatumRef{
        {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      k, missing,
    },
    {
      s, surprise,
    },
    {
      ih, good,
    },
    {
      t, possible,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    p, ih, t,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        p, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        ih, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 4}, {3, 3}, {4, 4},
      },
    },
    {
      phonVerdict{
        t, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        s,
        9, 21,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        ih,
        22, 35,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {3, 3}, {2, 2}, {4, 4}, {2, 4},
      },
    },
    {
      psPhonemeDatum{
        t,
        36, 53,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      p, missing,
    },
    {
      s, surprise,
    },
    {
      ih, good,
    },
    {
      t, good,
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
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        k, good,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {2, 1}, {2, 1}, {4, 1},
      } ,
    },
    {
      phonVerdict{
        ao, missing,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 2}, {3, 2},
      } ,
    },
    {
      phonVerdict{
        ch, surprise,
      },
      []psPhonemeDatumRef{
        {0, 3}, {2, 3}, {3, 3},
      } ,
    },
    {
      phonVerdict{
        ah, surprise,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 2}, {3, 2},
      } ,
    },
    {
      phonVerdict{
        t, missing,
      },
      []psPhonemeDatumRef{
        {1, 3}, {4, 3},
      } ,
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        k,
        17, 26,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {}, // {0, 1}, {1, 1}, {3, 1}, {2, 1}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        ah,
        27, 48,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 2}, {3, 2},
      },
    },
    {
      psPhonemeDatum{
        ch,

        49, 59,
      },
      []psPhonemeDatumRef{
      {0, 3}, {2, 3}, {3, 3},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      k, good,
    },
    {
      ao, missing,
    },
    {
      ah, surprise,
    },
    {
      ch, surprise,
    },
    {
      t, missing,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    ah, dh, er,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        ah, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      } ,
    },
    {
      phonVerdict{
        dh, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      } ,
    },
    {
      phonVerdict{
        er, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      } ,
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        ao,
        19, 37,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        er,
        54, 69,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      ah, missing,
    },
    {
      dh, missing,
    },
    {
      ao, surprise,
    },
    {
      er, good,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    w, ao, t, er,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        w, missing,
      },
      []psPhonemeDatumRef{
        {1, 1},
      },
    },
    {
      phonVerdict{
        ao, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 2}, {2, 3}, {3, 3},
      },
    },
    {
      phonVerdict{
        t, good,
      },
      []psPhonemeDatumRef{
        {1, 3}, {2, 6}, {3,4}, {4, 6},
      },
    },
    {
      phonVerdict{
        er, good,
      },
      []psPhonemeDatumRef{
        {0, 5}, {1, 4}, {3, 5}, {4, 7},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        ao,
        30, 38,
      },
      []psPhonemeDatumRef{
        {1, 2}, {2, 3}, {3, 3},
      },
    },
    {
      psPhonemeDatum{
        l,
        41, 53,
      },
      []psPhonemeDatumRef{
        {0, 1}, {2, 5}, {4, 5},
      },
    },
    {
      psPhonemeDatum{
        t,
        54, 64,
      },
      []psPhonemeDatumRef{
        {1, 3}, {2, 6}, {3, 4}, {4, 6},
      },
    },
    {
      psPhonemeDatum{
        er,
        66, 84,
      },
      []psPhonemeDatumRef{
        {1, 4}, {3, 5}, {4, 7}, {0, 5},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      w, missing,
    },
    {
      ao, good,
    },
    {
      t, good,
    },
    {
      er, good,
    },
  }

  // THEN
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    r, ih, p, l, ay, d,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        r, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        ih, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        p, good,
      },
      []psPhonemeDatumRef{
        {0, 6}, {2, 4}, {3, 6}, {4, 6},
      },
    },
    {
      phonVerdict{
        l, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        ay, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        d, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        r,
        31, 49,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        p,
        68, 74,
      },
      []psPhonemeDatumRef{
        {0, 6}, {3, 6}, {4, 6},
      },
    },
    {
      psPhonemeDatum{
        p,
        76, 77,
      },
      []psPhonemeDatumRef{
        {0, 6}, {2, 4}, {3, 6},
      },
    },
    {
      psPhonemeDatum{
        l,
        79, 86,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        ay,
        87, 115,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        d,
        116, 126,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      r, good,
    },
    {
      ih, good,
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
  fmt.Println(actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    ah, dh, ah,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        ah, possible,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 1}, {4, 3},
      },
    },
    {
      phonVerdict{
        dh, missing,
      },
      []psPhonemeDatumRef{
        {1, 2},
      },
    },
    {
      phonVerdict{
        ah, possible,
      },
      []psPhonemeDatumRef{
        {1, 3}, {2, 5,}, {3, 5},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        ao,
        30, 43,
      },
      []psPhonemeDatumRef{
        {0, 1}, {2, 1}, {3, 1}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        z,
        53, 58,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 4}, {3, 4},
      },
    },
    {
      psPhonemeDatum{
        ah,
        60, 76,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {3, 5}, {4, 3}, {2, 5},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      ao, surprise,
    },
    {
      z, surprise,
    },
    {
      ah, good,
    },
    {
      dh, missing,
    },
    {
      ah, missing,
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    w, ao, t, er,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        w, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        v, surprise,
      },
      []psPhonemeDatumRef{
        {2, 1}, {3, 1}, {4,1},
      },
    },
    {
      phonVerdict{
        ao, missing,
      },
      []psPhonemeDatumRef{
        {2, 2}, {3, 2},
      },
    },
    {
      phonVerdict{
        t, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 5},
      },
    },
    {
      phonVerdict{
        er, good,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {3, 4}, {2, 4}, {4, 6},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        v,
        32, 37,
      },
      []psPhonemeDatumRef{
        {2, 1}, {3, 1}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        l,
        43, 57,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {4, 3},
      },
    },
    {
      psPhonemeDatum{
        t,
        65, 74,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 5},
      },
    },
    {
      psPhonemeDatum{
        er,
        75, 94,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {3, 4}, {2, 4}, {4, 6},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      w, missing,
    },
    {
      v, surprise,
    },
    {
      ao, missing,
    },
    {
      l, surprise,
    },
    {
      t, good,
    },
    {
      er, good,
    },
  }

  // THEN
  fmt.Println("actual =", actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    ae, n, ah, m, ah, l,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        ae, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        n, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        ah, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        m, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        ah, good,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      phonVerdict{
        l, possible,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme {
    {
      psPhonemeDatum{
        ae,
        30, 40,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        n,
        41, 42,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        ah,
        43, 64,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        m,
        52, 55,
      },
      []psPhonemeDatumRef{
        {}, {}, {}, {},
      },
    },
    {
      psPhonemeDatum{
        l,
        65, 69,
      },
      []psPhonemeDatumRef{
        {}, {}, {},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      ae, good,
    },
    {
      n, good,
    },
    {
      ah, good,
    },
    {
      m, good,
    },
    {
      ah, good,
    },
    {
      l, possible,
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    fmt.Println("actual =", actual)
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    hh, ey,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        hh, missing,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        w, surprise,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2},
      },
    },
    {
      phonVerdict{
        g, surprise,
      },
      []psPhonemeDatumRef{
        {1, 1}, {3, 2}, {4, 1},
      },
    },
    {
      phonVerdict{
        ey, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 5}, {4, 5},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        g,
        26, 30,
      },
      []psPhonemeDatumRef{
        {1, 1}, {3, 2}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        w,
        36, 51,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2}, {3, 4},
      },
    },
    {
      psPhonemeDatum{
        ey,
        53, 88,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 5}, {4, 5},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      hh, missing,
    },
    {
      g, surprise,
    },
    {
      w, surprise,
    },
    {
      g, surprise,
    },
    {
      ey, missing,
    },
  }

  // THEN
  fmt.Println("hay actual =", actual)
  if !reflect.DeepEqual(actual, expected) {
    fmt.Println("hay actual =", actual)
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    ae, n, ah, m, ah, l,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        ae, good,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
      },
    },
    {
      phonVerdict{
        n, possible,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 2}, {4, 2},
      },
    },
    {
      phonVerdict{
        ah, good,
      },
      []psPhonemeDatumRef{
      },
    },
    {
      phonVerdict{
        m, good,
      },
      []psPhonemeDatumRef{
        {1, 4}, {2, 4}, {3, 4}, {4, 4},
      },
    },
    {
      phonVerdict{
        ah, good,
      },
      []psPhonemeDatumRef{
        {1, 5}, {2, 5}, {3, 5}, {4, 5},
      },
    },
    {
      phonVerdict{
        l, good,
      },
      []psPhonemeDatumRef{
        {0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        ae,
        27, 43,
      },
      []psPhonemeDatumRef{
        {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        n,
        46, 49,
      },
      []psPhonemeDatumRef{
        {0, 2}, {2, 2}, {4, 2},
      },
    },
    {
      psPhonemeDatum{
        m,
        52, 54,
      },
      []psPhonemeDatumRef{
        {1, 4}, {3, 4}, {4, 4},
      },
    },
    {
      psPhonemeDatum{
        ah,
        55, 70,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 5}, {2, 3}, {3, 5}, {4, 5}, {2, 5},
      },
    },
    {
      psPhonemeDatum{
        l,
        71, 78,
      },
      []psPhonemeDatumRef{
        {1, 6}, {2, 6}, {3, 6}, {4, 6}, {0, 6},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      ae, good,
    },
    {
      n, possible,
    },
    {
      ah, missing,
    },
    {
      m, possible,
    },
    {
      ah, good,
    },
    {
      l, good,
    },
  }

  // THEN
  if !reflect.DeepEqual(actual, expected) {
    fmt.Println("animal actual =", actual)
    te.Error()
  }

  // and GIVEN
  phons = []phoneme{
    s, k, uw, l,
  }
  c = newCombiner(phons)
  c.ruleAlignedVerdict = []linkedPhonVerdict{
    {
      phonVerdict{
        ae, surprise,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 1},
      },
    },
    {
      phonVerdict{
        ch, surprise,
      },
      []psPhonemeDatumRef{
        {1, 1}, {2, 1}, {3, 1},
      },
    },
    {
      phonVerdict{
        s, good,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
      },
    },
    {
      phonVerdict{
        k, good,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
      },
    },
    {
      phonVerdict{
        uw, good,
      },
      []psPhonemeDatumRef{
        {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5},
      },
    },
    {
      phonVerdict{
        l, good,
      },
      []psPhonemeDatumRef{
        {0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6},
      },
    },
  }
  c.timeAlignedPhonemes = []timeAlignedPhoneme{
    {
      psPhonemeDatum{
        ch,
        18, 23,
      },
      []psPhonemeDatumRef{
        {1, 1}, {2, 1}, {3, 1},
      },
    },
    {
      psPhonemeDatum{
        ae,
        24, 40,
      },
      []psPhonemeDatumRef{
        {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 1},
      },
    },
    {
      psPhonemeDatum{
        s,
        41, 66,
      },
      []psPhonemeDatumRef{
        {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
      },
    },
    {
      psPhonemeDatum{
        k,
        67, 71,
      },
      []psPhonemeDatumRef{
        {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
      },
    },
    {
      psPhonemeDatum{
        uw,
        72, 84,
      },
      []psPhonemeDatumRef{
        {0, 5}, {2, 5}, {3, 5}, {1, 5}, {4, 5},
      },
    },
    {
      psPhonemeDatum{
        l,
        85, 109,
      },
      []psPhonemeDatumRef{
        {0, 6}, {1, 6}, {2, 6}, {4, 6}, {3, 6},
      },
    },
  }

  // WHEN
  actual = c.integrate()
  expected = []phonVerdict{
    {
      ch, surprise,
    },
    {
      ae, surprise,
    },
    {
      s, good,
    },
    {
      k, good,
    },
    {
      uw, good,
    },
    {
      l, good,
    },
  }

  // THEN
  fmt.Println("school actual =", actual)
  if !reflect.DeepEqual(actual, expected) {
    te.Error()
  }
}
*/

func Test_integrate(te *testing.T) {
	// GIVEN
	phons := []phoneme{
		r, ah, sh, ah,
	}
	c := newCombiner(phons)
	c.results = []parsableResult{
		{
			psPhonemeResults{
				72,
				[]psPhonemeDatum{
					{sil, 0, 0}, {r, 63, 68}, {ah, 69, 69}, {zh, 74, 89}, {y, 90, 97}, {ah, 99, 121}, {sil, 122, 122},
				},
			},
			R_target{},
		},
		{
			psPhonemeResults{
				80,
				[]psPhonemeDatum{
					{sil, 0, 0}, {r, 64, 69}, {ah, 70, 70}, {zh, 74, 89}, {y, 90, 98}, {ah, 99, 110}, {n, 111, 120}, {b, 121, 131}, {sil, 133, 133},
				},
			},
			R_target{},
		},
		{
			psPhonemeResults{
				91,
				[]psPhonemeDatum{
					{sil, 0, 0}, {r, 64, 68}, {ah, 69, 75}, {sh, 76, 90}, {ah, 91, 121}, {sil, 122, 122},
				},
			},
			R_target{},
		},
		{
			psPhonemeResults{
				105,
				[]psPhonemeDatum{
					{sil, 0, 0}, {r, 64, 69}, {ah, 70, 75}, {sh, 76, 90}, {y, 90, 99}, {er, 100, 132}, {sil, 133, 133},
				},
			},
			R_target{},
		},
		{
			psPhonemeResults{
				125,
				[]psPhonemeDatum{
					{sil, 0, 0}, {r, 64, 68}, {ah, 69, 74}, {sh, 74, 90}, {ah, 91, 110}, {m, 111, 121}, {l, 122, 122}, {sil, 124, 124},
				},
			},
			R_target{},
		},
	}
	c.ruleAlignedVerdict = []linkedPhonVerdict{
		{
			phonVerdict{
				r, good,
			},
			[]psPhonemeDatumRef{
				{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
			},
		},
		{
			phonVerdict{
				ah, possible,
			},
			[]psPhonemeDatumRef{
				{2, 2}, {3, 2}, {4, 2},
			},
		},
		{
			phonVerdict{
				sh, possible,
			},
			[]psPhonemeDatumRef{
				{2, 3}, {3, 3}, {4, 3},
			},
		},
		{
			phonVerdict{
				ah, good,
			},
			[]psPhonemeDatumRef{
				{0, 5}, {1, 5}, {2, 4}, {4, 4},
			},
		},
	}
	c.timeAlignedPhonemes = []timeAlignedPhoneme{
		{
			psPhonemeDatum{
				r,
				64, 68,
			},
			[]psPhonemeDatumRef{
				{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
			},
		},
		{
			psPhonemeDatum{
				ah,
				70, 74,
			},
			[]psPhonemeDatumRef{
				{2, 2}, {3, 2}, {4, 2},
			},
		},
		{
			psPhonemeDatum{
				sh,
				76, 90,
			},
			[]psPhonemeDatumRef{
				{2, 3}, {3, 3}, {4, 3},
			},
		},
		{
			psPhonemeDatum{
				y,
				90, 97,
			},
			[]psPhonemeDatumRef{
				{0, 4}, {1, 4}, {3, 4},
			},
		},
		{
			psPhonemeDatum{
				ah,
				99, 110,
			},
			[]psPhonemeDatumRef{
				{0, 5}, {1, 5}, {2, 4}, {4, 4},
			},
		},
	}

	// WHEN
	actual := c.integrate()
	expected := []phonVerdict{
		{
			r, good,
		},
		{
			ah, possible,
		},
		{
			sh, possible,
		},
		{
			y, surprise,
		},
		{
			ah, good,
		},
	}

	// THEN
	fmt.Println("russia actual =", actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_insert(te *testing.T) {
	// GIVEN
	// Don't care what the phonemes are for these tests
	c := newCombiner([]phoneme{})
	verdicts := []timedPhonVerdict{
		{
			phonVerdict{
				t, good,
			},
			10,
		},
		{
			phonVerdict{
				eh, good,
			},
			20,
		},
		{
			phonVerdict{
				s, good,
			},
			30,
		},
		{
			phonVerdict{
				t, good,
			},
			40,
		},
	}

	// WHEN
	ph := timeAlignedPhoneme{
		psPhonemeDatum{
			aa,
			5, 10,
		},
		[]psPhonemeDatumRef{
			{}, {}, {}, {},
		},
	}
	actual := c.insert(verdicts, ph)
	expected := []timedPhonVerdict{
		{
			phonVerdict{
				aa, good,
			},
			5,
		},
		{
			phonVerdict{
				t, good,
			},
			10,
		},
		{
			phonVerdict{
				eh, good,
			},
			20,
		},
		{
			phonVerdict{
				s, good,
			},
			30,
		},
		{
			phonVerdict{
				t, good,
			},
			40,
		},
	}

	// THEN
	fmt.Println(actual)
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}

func Test_ruleAlign(te *testing.T) {
	// GIVEN
	phons := []phoneme{
		k, ao, t,
	}
	c := newCombiner(phons)
	c.results = []parsableResult{
		{
			psPhonemeResults{
				90,
				[]psPhonemeDatum{
					{
						sil, 0, 16,
					},
					{
						k, 17, 24,
					},
					{
						ah, 26, 48,
					},
					{
						ch, 49, 60,
					},
					{
						t, 61, 61,
					},
					{
						sil, 64, 106,
					},
				},
			},
			R_target{
				"cot",
				[]parsableRule{
					R_trappedOpening{
						R_trap{
							f,
						},
						R_opening{
							z,
							R_and{
								[]rule{},
							},
						},
					},
					R_phoneme{
						k, m,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						ao, ch,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						t, w,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_closing{
						n,
						R_and{
							[]rule{},
						},
					},
				},
			},
		},
		{
			psPhonemeResults{
				105,
				[]psPhonemeDatum{
					{
						sil, 0, 16,
					},
					{
						k, 17, 26,
					},
					{
						ao, 27, 46,
					},
					{
						t, 47, 63,
					},
					{
						sil, 64, 103,
					},
				},
			},
			R_target{
				"cot",
				[]parsableRule{
					R_trappedOpening{
						R_trap{
							sh,
						},
						R_opening{
							s,
							R_and{
								[]rule{},
							},
						},
					},
					R_phoneme{
						k, m,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						ao, sh,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						t, w,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_closing{
						b,
						R_and{
							[]rule{},
						},
					},
				},
			},
		},
		{
			psPhonemeResults{
				120,
				[]psPhonemeDatum{
					{
						sil, 0, 17,
					},
					{
						k, 18, 26,
					},
					{
						ah, 27, 48,
					},
					{
						ch, 49, 59,
					},
					{
						t, 60, 60,
					},
					{
						sil, 63, 103,
					},
				},
			},
			R_target{
				"cot",
				[]parsableRule{
					R_trappedOpening{
						R_trap{
							s,
						},
						R_opening{
							ch,
							R_and{
								[]rule{},
							},
						},
					},
					R_phoneme{
						k, th,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						ao, ch,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						t, w,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_closing{
						g,
						R_and{
							[]rule{},
						},
					},
				},
			},
		},
		{
			psPhonemeResults{
				143,
				[]psPhonemeDatum{
					{
						sil, 0, 16,
					},
					{
						k, 17, 25,
					},
					{
						ah, 26, 48,
					},
					{
						ch, 49, 59,
					},
					{
						t, 59, 59,
					},
					{
						sil, 62, 105,
					},
				},
			},
			R_target{
				"cot",
				[]parsableRule{
					R_trappedOpening{
						R_trap{
							w,
						},
						R_opening{
							ch,
							R_and{
								[]rule{},
							},
						},
					},
					R_phoneme{
						k, s,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						ao, ch,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						t, w,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_closing{
						b,
						R_and{
							[]rule{},
						},
					},
				},
			},
		},
		{
			psPhonemeResults{
				190,
				[]psPhonemeDatum{
					{
						sil, 0, 18,
					},
					{
						k, 19, 27,
					},
					{
						ao, 28, 39,
					},
					{
						t, 40, 61,
					},
					{
						sil, 61, 105,
					},
				},
			},
			R_target{
				"cot",
				[]parsableRule{
					R_trappedOpening{
						R_trap{
							w,
						},
						R_opening{
							d,
							R_and{
								[]rule{},
							},
						},
					},
					R_phoneme{
						k, n,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						ao, sh,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_phoneme{
						t, w,
						[]phoneme{},
						R_and{
							[]rule{},
						},
					},
					R_closing{
						g,
						R_and{
							[]rule{},
						},
					},
				},
			},
		},
	}
	c.parseMap = parseMap{
		[][]parseResult{
			{
				{0, 0, false}, {1, 1, true}, {2, 3, false}, {4, 4, true}, {5, 5, false},
			},
			{
				{0, 0, false}, {1, 1, true}, {2, 2, true}, {3, 3, true}, {4, 4, false},
			},
			{
				{0, 0, false}, {1, 1, true}, {2, 3, false}, {4, 4, true}, {5, 5, false},
			},
			{
				{0, 0, false}, {1, 1, true}, {2, 3, false}, {4, 4, true}, {5, 5, false},
			},
			{
				{0, 0, false}, {1, 1, true}, {2, 2, true}, {3, 3, true}, {4, 4, false},
			},
		},
	}

	// WHEN
	actual := c.ruleAlign()
	expected := []linkedPhonVerdict{
		{
			phonVerdict{
				k, good,
			},
			[]psPhonemeDatumRef{
				{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
			},
		},
		{
			phonVerdict{
				ao, missing,
			},
			[]psPhonemeDatumRef{
				{1, 2}, {4, 2},
			},
		},
		{
			phonVerdict{
				ah, surprise,
			},
			[]psPhonemeDatumRef{
				{0, 2}, {2, 2}, {3, 2},
			},
		},
		{
			phonVerdict{
				ch, surprise,
			},
			[]psPhonemeDatumRef{
				{0, 3}, {2, 3}, {3, 3},
			},
		},
		{
			phonVerdict{
				t, missing,
			},
			[]psPhonemeDatumRef{
				{1, 3}, {4, 3},
			},
		},
	}

	// THEN
	if !reflect.DeepEqual(actual, expected) {
		te.Error()
	}
}
