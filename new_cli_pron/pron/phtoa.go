package pron

import (
	"fmt"
	"strings"
)

type phonToAlphas struct {
	phons  []phoneme
	alphas string
}

func (p phonToAlphas) equal(q phonToAlphas) bool {
	if p.alphas == q.alphas {
		if len(p.phons) != len(q.phons) {
			return false
		}
		for i := 0; i < len(p.phons); i++ {
			if p.phons[i] != q.phons[i] {
				return false
			}
		}
		return true
	}
	return false
}

type phToAError struct {
	phon phoneme
	word string
}

func (p phToAError) Error() string {
	return fmt.Sprintf("Phoneme translation failure, translating phoneme %s in word, %s", string(p.phon), p.word)
}

func stringAt(s string, start, length int, ts ...string) bool {
	if start < 0 || start+length > len(s)+1 {
		return false
	}
	subS := s[start:]
	for _, t := range ts {
		if len(t) != length {
			continue
		}
		if strings.HasPrefix(subS, t) {
			return true
		}
	}
	return false
}

/*
func getStringAt(s string, start, length int, ts ...string) (string, bool) {
  if start < 0 || start + length > len(s) + 1 {
    return "", false
  }
  subS := s[start:]
  for _, t := range ts {
    if len(t) != length {
      continue
    }
    if strings.HasPrefix(subS, t) {
      return t, true
    }
  }
  return "", false
}
*/

// Returns the first string in ts that appears as a prefix of s[start:] - so
// the order in which strings are passed to this function is important
func getStringAt(s string, start, length int, ts ...string) (string, bool) {
	if start < 0 {
		return "", false
	}
	subS := s[start:]
	for _, t := range ts {
		if strings.HasPrefix(subS, t) {
			return t, true
		}
	}
	return "", false
}

func nextPhIsConsonant(ps []phoneme, start int) bool {
	if start < 0 || start+1 >= len(ps) {
		return false
	}
	return isConsonant(ps[start+1])
}

func phsAt(ps []phoneme, start, length int, qs ...[]phoneme) ([]phoneme, bool) {
	hasPrefix := func(p, q []phoneme) bool {
		for i, q1 := range q {
			if p[i] != q1 {
				return false
			}
		}
		return true
	}
	if start < 0 || start+length > len(ps) {
		return []phoneme{}, false
	}
	subPh := ps[start:]
	for _, q := range qs {
		if len(q) != length {
			continue
		}
		if hasPrefix(subPh, q) {
			return q, true
		}
	}
	return []phoneme{}, false
}

func isConsonant(ph phoneme) bool {
	consonants := []phoneme{
		b, bl, ch, d, dh, dz, f, g, gr, hh, jh, k, kl, kr, ks, kw, l, m, n, ng, p, pl, pr, r, s, sh, t, th, tr, ts, v, w, z, zh,
	}
	for _, c := range consonants {
		if ph == c {
			return true
		}
	}
	return false
}

func isVowel(ph phoneme) bool {
	vowels := []phoneme{
		aa, ae, ah, ao, aw, ay, eh, er, ey, ih, iy, ow, oy, uw, uh, y,
	}
	for _, v := range vowels {
		if ph == v {
			return true
		}
	}
	return false
}

func trailingSilentAlphas(abc string, currAbc int, phons []phoneme, currPh int, currMap phonToAlphas) string {
	unmappedAlphas := ""
	if currAbc+len(currMap.alphas) < len(abc) {
		unmappedAlphas = abc[currAbc+len(currMap.alphas):]
	}
	if (currPh+len(currMap.phons) == len(phons)) &&
		(len(unmappedAlphas) != 0) {
		return unmappedAlphas
	}
	return ""
}

func isSilentE(abc string, currAbc int, phons []phoneme, currPh int, currMap phonToAlphas) bool {
	nextPhIndex := currPh + len(currMap.phons)
	if nextPhIndex > len(phons)-1 {
		// There is no next phoneme so check to see if the next character is an e
		if stringAt(abc, currAbc+len(currMap.alphas), 1, "e") {
			return true
		}
		return false
	}
	nextPh := phons[nextPhIndex]

	// In most cases, if the next phoneme is not a consonant then this is likely
	// not a silent e, but there are exceptions...
	//
	// ...words like
	// G R EY V Y AA D
	if stringAt(abc, currAbc+len(currMap.alphas), 2, "ey") && nextPh == y {
		return true
	}
	// T AY M AW T
	if stringAt(abc, currAbc+len(currMap.alphas), 2, "eo") && nextPh == aw {
		return true
	}
	// P AY N AE P L
	if stringAt(abc, currAbc+len(currMap.alphas), 2, "ea") && nextPh == ae {
		return true
	}
	// H OW M OW N ER
	if stringAt(abc, currAbc+len(currMap.alphas), 2, "eo") && nextPh == ow {
		return true
	}
	// K L OW S AH P
	if stringAt(abc, currAbc+len(currMap.alphas), 2, "eu") && nextPh == ah {
		return true
	}

	// Should really treat ed as a special case here

	if !isConsonant(nextPh) {
		return false
	}
	// If we get this far then we have two adjacent consonant phonemes.
	// Check to see if the next letter is an e {
	//

	// BUT...
	// This is pretty ugly but having two adjacent phoneme consonants and a next letter
	// e isn't enough. Take G AH V N for instance (and there are plenty of other words
	// like it)... So we test for that explictly here
	if stringAt(abc, currAbc+len(currMap.alphas), 3, "ern") && nextPh == n {
		return false
	}
	// Now we check if the next lexical character is e
	if stringAt(abc, currAbc+len(currMap.alphas), 1, "e") {
		return true
	}
	return false
}

func (m phonToAlphas) mapB(phons []phoneme, currPh int, alphas string, currAbc int) phonToAlphas {
	new := m
	currAbc += len(new.alphas)
	currPh += len(new.phons)

	// Treat february as a special case
	if stringAt(alphas, currAbc, 6, "bruary") {
		new = phonToAlphas{
			append(m.phons, []phoneme{
				b,
			}...),
			m.alphas + "br",
		}
		return new
	}
	if stringAt(alphas, currAbc, 2, "bb") {
		if _, ok := phsAt(phons, currPh, 2, []phoneme{b, b}); !ok {
			// As in stuBBorn,...
			new = phonToAlphas{
				append(m.phons, []phoneme{
					b,
				}...),
				m.alphas + "bb",
			}
			return new	
		}
	}
	if stringAt(alphas, currAbc, 2, "pb") {
		// There's a silent p here
		// As in cuPBoard,...
		new = phonToAlphas{
			append(m.phons, []phoneme{
				b,
			}...),
			m.alphas + "pb",
		}
		return new
	}
	if stringAt(alphas, currAbc, 1, "b") {
		new = phonToAlphas{
			append(m.phons, []phoneme{
				b,
			}...),
			m.alphas + "b",
		}
		return new
	}
	return new
}

func (m phonToAlphas) mapL(phons []phoneme, currPh int, alphas string, currAbc int) phonToAlphas {
	new := m
	currAbc += len(new.alphas)
	currPh += len(new.phons)

	if s, ok := getStringAt(alphas, currAbc, 0, "lel"); ok {
		if _, ok := phsAt(phons, currPh, 2, []phoneme{l, l}, []phoneme{l, ey}); !ok {
			// As in candLELight,...
			// But not as in candLelight, ukeLele,...
			new = phonToAlphas{
				[]phoneme{
					l,
				},
				s,
			}
			return new
		}
	}
	if s, ok := getStringAt(alphas, currAbc, 0, "hl", "ll"); ok {
		if _, ok := phsAt(phons, currPh, 2, []phoneme{l, l}); !ok {
			// As in daHLia, yeLLow,...
			// But not as in goaLLess, we...
			new = phonToAlphas{
				[]phoneme{
					l,
				},
				s,
			}
			return new
		}
	}
	// Check for the silent h in delhi
	if s, ok := getStringAt(alphas, currAbc, 0, "lh"); ok {
		if _, ok := phsAt(phons, currPh, 2, []phoneme{l, hh}); !ok {
			// The h isn't sounded
			// As in deLHi,...
			new = phonToAlphas{
				[]phoneme{
					l,
				},
				s,
			}
			return new
		}
	}
	if stringAt(alphas, currAbc, 1, "l") {
		new = phonToAlphas{
			[]phoneme{
				l,
			},
			"l",
		}
		return new
	}
	return new
}

func mapPhToA(phons []phoneme, alphas string) ([]phonToAlphas, error) {
	fail := func(p phoneme, a string) ([]phonToAlphas, error) {
		ret := []phonToAlphas{}
		err := phToAError{
			p,
			a,
		}
		return ret, err
	}
	ret := []phonToAlphas{}
	// current := 0
	currPh := 0
	currAbc := 0
	punctuationSkipped := ""
	for currPh < len(phons) {
		// for _, phon := range phons {
		var new phonToAlphas
		phon := phons[currPh]
		// Check we haven't run out of letters and return if we have
		if currAbc > len(alphas)-1 {
			return fail(phons[currPh], alphas)
		}
		if (alphas[currAbc] == '\'' && phon != ih) || alphas[currAbc] == '-' {
			// ' does not typically get expressed as a phoneme so move the current
			// character on. If the current phoneme is ih though then it does, for
			// instance in james's which has phonemes jh ey m z ih z
			// - has no effect on pronunciation
			//
			punctuationSkipped = string(alphas[currAbc])
			currAbc++
			continue
		}
		switch phon {
		case aa:
			if s, ok := getStringAt(alphas, currAbc, 0, "aar", "arr", "ear", "har", "er", "aa", "ar", "or"); ok {
				// As in AARdvark, stARRed, hEARt, philHARmonic, clERk, aardvARk, tomORrow,...
				// Now check whether we have a US r phoneme following
				if _, ok := phsAt(phons, currPh, 2, []phoneme{aa, r}); ok {
					s = s[:len(s)-1]
				}
				new = phonToAlphas{
					[]phoneme{
						aa,
					},
					s,
				}
				break
			}

			// Note that aa can sound like oh in cot and ah as in barn
			if s, ok := getStringAt(alphas, currAbc, 2, "au", "ha", "ho", "ou"); ok {
				// As in cAUght, gymkHAna, HOnest, cOUgh,...
				new = phonToAlphas{
					[]phoneme{
						aa,
					},
					s,
				}
				break
			}
			// Some words have an -ah ending with the 'h' silent
			if stringAt(alphas, currAbc, 2, "ah") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{aa, hh}); !ok {
					// The 'h' isn't sounded
					// As in shAH,...
					new = phonToAlphas{
						[]phoneme{
							aa,
						},
						"ah",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "o"); ok {
				// As in , gEnre, ,...
				new = phonToAlphas{
					[]phoneme{
						aa,
					},
					s,
				}
				break
			}
		case ae:
			if s, ok := getStringAt(alphas, currAbc, 2, "ai", "au", "ei"); ok {
				// As in plAIts, drAUght, revEIlle,,...
				new = phonToAlphas{
					[]phoneme{
						ae,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "i"); ok {
				// As in bAt, rEbate, merIngue,..
				new = phonToAlphas{
					[]phoneme{
						ae,
					},
					s,
				}
				break
			}
		case ah:
			if alphas == "hiccough" {
				new = phonToAlphas{
					[]phoneme{
						ah, p,
					},
					"ough",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "wer") {
				// As in ansWER,...
				s := "wer"
				// But only if the 'r' isn't sounded
				// As in ansWErable,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ah, r}); ok {
					s = "we"
				}
				new = phonToAlphas{
					[]phoneme{
						ah,
					},
					s,
				}
				break
			}
			// if p, ok := phsAt(phons, currPh, 2, []phoneme{ah, l}, []phoneme{ah, m}, []phoneme{ah, n}); ok {
			// 	// Trying to spot an inserted schwa as in bottLe, theisM, wasN't,...
			// 	if stringAt(alphas, currAbc, 2, "le") {
			// 		// Silent e?
			// 		new = phonToAlphas{
			// 			p,
			// 			"le",
			// 		}
			// 		break
			// 	}
			// 	if s, ok := getStringAt(alphas, currAbc, 1, "l", "m", "n"); ok {
			// 		new = phonToAlphas{
			// 			p,
			// 			s,
			// 		}
			// 		break
			// 	}
			// }
			// if stringAt(alphas, currAbc, 3, "iou") {
			// 	// As in suspicIOUs,...
			// 	new = phonToAlphas{
			// 		[]phoneme{
			// 			ah,
			// 		},
			// 		"iou",
			// 	}
			// 	break
			// }
			if s, ok := getStringAt(alphas, currAbc, 3, "our", "ure"); ok {
				// As in flavOUR, futURE...
				// But not if the r is sounded, for instance in armOURy, usUREr...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ah, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ah,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ur") {
				// As in dUration,...
				// If the r is sounded, we should only grab the 'u'
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ah, r}); ok {
					new = phonToAlphas{
						[]phoneme{
							ah,
						},
						"u",
					}
					break
				}
			}
			// Handling this separately as this is a mix of a silent e followed by a
			// vowel so I may handle this more generally at some point
			if stringAt(alphas, currAbc, 2, "ea") {
				// As in likEAble,...
				new = phonToAlphas{
					[]phoneme{
						ah,
					},
					"ea",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "au", "ei", "ia", "ie", "io", "oo", "ou", "re", "ua", "ui", "ur", "yr"); ok {
				// As in becAUse, forEIgn, catERpillar, RussIA, conscIEnce, percussIOn, blOOd, rOUgh, theatRE, usUAlly, biscUIt, treasURe, zephYR,...
				new = phonToAlphas{
					[]phoneme{
						ah,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ar", "er", "or"); ok {
				// As in wizARd, pERcussion tractOR...
				// Need to be careful here. I don't want to consume the r if there's
				// an r phoneme in the phonetic spelling, for instance as in
				// d(d)o(ao)c(k)u(y uh)m(m)e(eh)n(n)t(t)a(ah)r(r)y(iy)
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ah, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ah,
						},
						s,
					}
					break
				}
			}
			// Some words have an -ah ending with the 'h' silent
			if stringAt(alphas, currAbc, 2, "ah") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ah, hh}); !ok {
					// The 'h' isn't sounded
					// As in purdAH,...
					new = phonToAlphas{
						[]phoneme{
							ah,
						},
						"ah",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "i", "o", "u", "y"); ok {
				// As in ..., propYlene,...
				new = phonToAlphas{
					[]phoneme{
						ah,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "r") {
				// As in houR,...
				new = phonToAlphas{
					[]phoneme{
						ah,
					},
					"r",
				}
				break
			}
		case ao:
			if stringAt(alphas, currAbc, 3, "hon") {
				// Words like HOnest start with a silent 'h'
				// As in HOnour,... and they don't necessarily start with hon-, for
				// instance disHOnourable,...
				new = phonToAlphas{
					[]phoneme{
						ao,
					},
					"ho",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ort") {
				// As in the borrowed French word, rappORT,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, ch}, []phoneme{ao, dh}, []phoneme{ao, sh}, []phoneme{ao, t}, []phoneme{ao, tr}, []phoneme{ao, th}); !ok {
					// But not as in fORTunate, nORTHern, abORTion, repORT, pORTrait, nORTH,...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						"ort",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "orp"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, p}, []phoneme{ao, f}); !ok {
					// The p is not sounded so swallow it here
					// As in the French cORPs,...
					// But not as in cORpulent, amORphous,...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						s,
					}
					break
				}
			}
			// Borrowed French words ending -eur
			if stringAt(alphas, currAbc, 3, "eur") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, r}); !ok {
					// As in sabotEUR,...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						"eur",
					}
					break
				}
			}
			// The test for a phonetic 'r' following the 'ao' doesn't work for 
			// words like storeroom and forerunner so pulling these out as a special
			// case
			if (strings.HasPrefix(alphas, "storeroom") || strings.HasPrefix(alphas, "forerunner")) && stringAt(alphas, currAbc, 3, "ore") {
				new = phonToAlphas{
					[]phoneme{
						ao,
					},
					// The phonetic 'r' that follows 'ao' is for the lexical 'r' in room
					"ore",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "aor", "aur", "oar", "oor", "orr", "our", "wor", "or"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, r}); !ok {
					// As in extrAORdinary, dinosAUR, bOARd, flOOR, abhORRed, yOURs, sWORd, fORtnight...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						s,
					}
					break
				} else {
					// As in AUric, hOAry, mOOrish, pOUring, stOry,...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						// The 'r' is sounded so don't grab it here
						s[:len(s)-1],
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 4, "awer") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, ax}, []phoneme{ao, axr}); !ok {
					// The 'er' is not sounded
					// As in drAWER,...
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						"awer",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "au", "aw", "oa", "ou"); ok {
				// As in applAUd, clAWs, brOAd, thOUght,...
				new = phonToAlphas{
					[]phoneme{
						ao,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "arr", "ar"); ok {
				// As in wARRed, wAR,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ao, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ao,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "o", "a"); ok {
				// As in fOr, wAter...
				new = phonToAlphas{
					[]phoneme{
						ao,
					},
					s,
				}
				break
			}
		case aw:
			if stringAt(alphas, currAbc, 4, "ough") {
				// As in bOUGH,...
				new = phonToAlphas{
					[]phoneme{
						aw,
					},
					"ough",
				}
				break
			}
			if stringAt(alphas, currAbc, 4, "hour") {
				new = phonToAlphas{
					[]phoneme{
						aw,
					},
					"hou",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "our"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{aw, axn}); ok {
					// As in sOURNess,...
					new = phonToAlphas{
						[]phoneme{
							aw,
						},
						s,
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 3, []phoneme{aw, ax, r}); ok {
					// As in sOURest,...
					// Leave the 'r' as it's sounded
					new = phonToAlphas{
						p[:len(p)-1],
						s[:len(s)-1],
					}
					break
				}

			}
			if stringAt(alphas, currAbc, 3, "our") {
				new = phonToAlphas{
					[]phoneme{
						aw,
					},
					"ou",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ow") {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{aw, w}); ok {
					// The lexical w is sounded as in bOWing,...
					new = phonToAlphas{
						p,
						"ow",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ao", "au", "ho", "ou", "ow"); ok {
				// As in mAO, sAUerkrAUt, HOur, sOUnd, gOWn,...
				new = phonToAlphas{
					[]phoneme{
						aw,
					},
					s,
				}
				break
			}
		case ay:
			if stringAt(alphas, currAbc, 4, "eigh") {
				// As in hEIGHt,...
				new = phonToAlphas{
					[]phoneme{
						ay,
					},
					"eigh",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ai") {
				// As in nAIve,...
				// So check for this as two separate phonemes
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, iy}); ok {
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"a",
					}
					break
				}
			}
			if str, ok := getStringAt(alphas, currAbc, 0, "ais", "is"); ok {
				// As in AISle, ISland,...
				// But NOT as in ISolate, demonISe...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, s}, []phoneme{ay, z}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						str,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "aye", "eye", "igh", "ai", "ay", "uy"); ok {
				// As in AYE, EYE, hIGH, thAI, bUY, paraguAY...
				new = phonToAlphas{
					[]phoneme{
						ay,
					},
					s,
				}
				break
			}
			// Catching -ihi- words.
			if stringAt(alphas, currAbc, 3, "ihi") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, hh}); !ok {
					// The h is not sounded so swallow it here
					// As in nIHilism,...
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"ih",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ir") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, axr}); ok {
					// This sounds like ire so leave the re to be mapped to the axr phoneme
					// AS in wIre,...
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"i",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 3, []phoneme{ay, ax, r}); ok {
					// This sounds like 'f-ire', 're' so process the 'i' now and
					// leave the 'ax r' to be processed separately
					new = phonToAlphas{
						[]phoneme{
							ay, ax,
						},
						"i",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, er}); ok {
					// This sounds like 'f-ire', 're' so process the 'i' now and
					// leave the 'er' to be processed separately
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"i",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ir") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, r}); ok {
					// The r is sounded so don't swallow the r here
					// As in Ironic,...
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"i",
					}
				} else {
					// The r isn't sounded so swallow it now
					// As in IRon,...
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"ir",
					}
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ie") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, ah}, []phoneme{ay, eh}, []phoneme{ay, ih}, []phoneme{ay, iy}); ok {
					// Looks like the 'e' is sounded so don't swallow it here
					// As in dIEt, quIEscent, socIEtal, quIEtus...
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						"i",
					}
					break
				}
			}
			if str, ok := getStringAt(alphas, currAbc, 2, "ia", "ei", "ie"); ok {
				// As in dIAl, EIther, frIEd,...
				// But NOT as in trIAngle, quIEt, lIAr, strIAtion,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ay, ae}, []phoneme{ay, ax}, []phoneme{ay, axr}, []phoneme{ay, axl}, []phoneme{ay, ey}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ay,
						},
						str,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ae"); ok {
				// As in mAEstro,...
				new = phonToAlphas{
					[]phoneme{
						ay,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "i", "u", "y"); ok {
				// As in tIme, flUtist, wrYly,...
				new = phonToAlphas{
					[]phoneme{
						ay,
					},
					s,
				}
				break
			}
		case b:
			new = new.mapB(
				phons, currPh, alphas, currAbc,
			)
			break
		case bl:
			// Catch deleted schwas. There are several patterns...
			if s, ok := getStringAt(alphas, currAbc, 0, "bell", "boll", "bal", "bel", "bol"); ok {
				// As in laBELLed, gamBOLLing, suBALtern, laBEL, gamBOL,...
				new = phonToAlphas{
					[]phoneme{
						bl,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "bbl", "bl"); ok {
				new = phonToAlphas{
					[]phoneme{
						bl,
					},
					s,
				}
				break
			}
		case ch:
			if stringAt(alphas, currAbc, 3, "tch") {
				new = phonToAlphas{
					// As in stiTCH,...
					[]phoneme{
						ch,
					},
					"tch",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "cc", "ch", "cz", "tt", "c"); ok {
				// As in cappuCCino, CHeese, CZech, aTTune, Cello,...
				new = phonToAlphas{
					[]phoneme{
						ch,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "t") {
				// As in signature, overture, etc
				new = phonToAlphas{
					[]phoneme{
						ch,
					},
					"t",
				}
				break
			}
		case d:
			// A training pronunciation to get out of the way first
			if stringAt(alphas, currAbc, 1, "t") {
				// As in lofT,...
				new = phonToAlphas{
					[]phoneme{
						d,
					},
					"t",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "dd", "ld"); ok {
				// As in riDDle, and the silent l in wouLD,...
				p := []phoneme{d}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{d, d}); ok {
					// As in miDDay,... There's a double phonetic 'd' though, so don't grab
					//both lexical 'd's here
					s = "d"
				}
				new = phonToAlphas{
					p,
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ed") {
				// As in wanderED,...
				// Although there is a kind of silent e here it doesn't get picked up
				// by isSilentE in this case because the preceding phoneme is er and so
				// not a consonant
				new = phonToAlphas{
					[]phoneme{
						d,
					},
					"ed",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "d") {
				new = phonToAlphas{
					[]phoneme{
						d,
					},
					"d",
				}
				break
			}
		case dh:
			if stringAt(alphas, currAbc, 2, "th") {
				new = phonToAlphas{
					[]phoneme{
						dh,
					},
					"th",
				}
				break
			}
		case eh:
			if stringAt(alphas, currAbc, 1, "x") {
				// As in Xmas,...
				if p, ok := phsAt(phons, currPh, 3, []phoneme{eh, k, s}); ok {
					new = phonToAlphas{
						p,
						"x",
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 2, []phoneme{eh, ks}); ok {
					new = phonToAlphas{
						p,
						"x",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 3, "hei") {
				// As in HEIr,...
				new = phonToAlphas{
					[]phoneme{
						eh,
					},
					"hei",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "are") {
				// As in awARE,...
				// But not if the r is sounded as in aRena,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{eh, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							eh,
						},
						"are",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "air", "are", "ear"); ok {
				// As in chAIR, fARE, EARthenware,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{eh, r}, []phoneme{eh, er}); !ok {
					new = phonToAlphas{
						[]phoneme{
							eh,
						},
						s,
					}
					break
				} // else we handle the case with no r phoneme below
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ae", "ai", "ay", "ea", "ei", "eo", "ie"); ok {
				// As in AErobic, sAId, sAYs, endEAvour, thEIr, lEOpard, frIEnd,...
				new = phonToAlphas{
					[]phoneme{
						eh,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ar") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{eh, r}); !ok {
					// As in scARce,...
					new = phonToAlphas{
						[]phoneme{
							eh,
						},
						"ar",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "i", "u"); ok {
				// As in contrAry, rEd, squIrrel, bUried,...
				new = phonToAlphas{
					[]phoneme{
						eh,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "r") {
				// As in houR,...
				new = phonToAlphas{
					[]phoneme{
						eh,
					},
					"r",
				}
				break
			}
		case er:
			// Pulling this out as a special case for now because I can't think of
			// any other words this occurs in.
			if stringAt(alphas, currAbc-1, 4, "iron") {
				// As in iROn,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"ro",
				}
				break
			}
			// Another special case. olo as in colonel sounds as er.
			if stringAt(alphas, currAbc, 3, "olo") {
				// As in cOLOnel,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"olo",
				}
				break
			}
			// Lots of herb- words with a silent 'h' phonetic variant
			if strings.HasPrefix(alphas, "herb") && currAbc == 0 {
				// I'm being deliberately specific here and only looking for words
				// that start herb-
				// As in HERbal,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"her",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ure") {
				// As in treasURE, futURE,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"ure",
				}
				break
			}
			// Some er sounds can have a sounded r when combined with for instance -ing.
			if s, ok := getStringAt(alphas, currAbc, 3, "err", "eur", "irr", "urr", "er"); ok {
				// As in refERRed, entreprenEURial, stIRRed, pURRed, transfERable,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{er, r}); ok {
					// As in refErring, etc...
					// In this case don't swallow the r now, save it for the r phoneme
					s = s[:1]
				}
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "arr", "ear", "err", "eur", "irr", "our", "urr", "wer"); ok {
				// As in ARRay, hEARd, transfERRed, sabotEUR, stIRRed, yOURself, pURRed, ansWERed,...
				// What the hell! yourself is y er s eh l f but your is y ao r - why
				// the diffference?
				// That has to be a US English thing...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ah", "ar", "er", "eu", "ia", "ir", "or", "re", "ur", "yr"); ok {
				// As in hookAH,..., massEUse,..., theatRE, fUR, zephYR,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "a", "i") {
				// As in umbrellA, anImal...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"i",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "r") {
				// As in houR,...
				new = phonToAlphas{
					[]phoneme{
						er,
					},
					"r",
				}
				break
			}
		case ey:
			if stringAt(alphas, currAbc, 4, "eigh") {
				// As in wEIGH...
				new = phonToAlphas{
					[]phoneme{
						ey,
					},
					"eigh",
				}
				break
			}
			// Capture borrowed French word -er endings
			if stringAt(alphas, currAbc, 2, "er") {
				// As in dossiER,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ey, r}); !ok {
					// But not as in paYRoll,...
					new = phonToAlphas{
						[]phoneme{
							ey,
						},
						"er",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ai") {
				// As in pAId
				s := "ai"
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ey, ih}); ok {
					// The 'ai' is split across two phonetic vowels as in algebrAist,...
					s = "a"
				}
				new = phonToAlphas{
					[]phoneme{
						ey,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ea") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ey, aa}, []phoneme{ey, ax}, []phoneme{ey, axr}); ok {
					// As in rEAl (the unit of currency), eritrEA, eritrEA,...
					// Only swallow the e, the a will be picked up later
					new = phonToAlphas{
						[]phoneme{
							ey,
						},
						"e",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ei") {
				// As in cunEIform,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ey, ih}); ok {
					// The ei is split across two phonetic vowels so only grab
					// the first one
					new = phonToAlphas{
						[]phoneme{
							ey,
						},
						"e",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "au", "ay", "ea", "ee", "ei", "ey"); ok {
				// As in gAUge, wAY, grEAt, , soirEE, EIght, gourmET, thEY,...
				new = phonToAlphas{
					[]phoneme{
						ey,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "et") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ey, t}); !ok {
					// Swallow the t here, as in cabarET,...
					new = phonToAlphas{
						[]phoneme{
							ey,
						},
						"et",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "a") {
				new = phonToAlphas{
					[]phoneme{
						ey,
					},
					"a",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "e") {
				// As in borrowed French words like touchE,...
				new = phonToAlphas{
					[]phoneme{
						ey,
					},
					"e",
				}
				break
			}
		case f:
			if stringAt(alphas, currAbc, 2, "gh") {
				// As in couGH,...
				new = phonToAlphas{
					[]phoneme{
						f,
					},
					"gh",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "pph", "ph"); ok {
				// As in saPPhire, PHosPHorus,...
				new = phonToAlphas{
					[]phoneme{
						f,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ff") {
				new = phonToAlphas{
					[]phoneme{
						f,
					},
					"ff",
				}
				break
			}
			// Catch a possible silent l
			if stringAt(alphas, currAbc, 2, "lf") {
				// As in caLf,...
				new = phonToAlphas{
					[]phoneme{
						f,
					},
					"lf",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ft") {
				// A 't' following 'f' can be silent. Trying to catch this here because
				// it's harder to catch on the next phoneme
				if _, ok := phsAt(phons, currPh, 2, []phoneme{f, t}, []phoneme{f, d}, []phoneme{f, th}); !ok {
					// As in soFTen, but not as in liFT, loFT*, fiFTh ...
					// *This is a training pronunciation
					new = phonToAlphas{
						[]phoneme{
							f,
						},
						"ft",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "f") {
				new = phonToAlphas{
					[]phoneme{
						f,
					},
					"f",
				}
				break
			}
		case g:
			// A very special case to start with
			if alphas == "blackguard" && stringAt(alphas, currAbc, 2, "ck") {
				// There's a silent ck, and also a silent u after the g
				new = phonToAlphas{
					[]phoneme{
						g,
					},
					"ckgu",
				}
				break
			}
			// Trying to trap the silent h that can follow ex-
			if stringAt(alphas, currAbc, 2, "xh") {
				// As in eXHaust,...
				_, ok_gz := phsAt(phons, currPh, 2, []phoneme{g, z})
				_, ok_gzh := phsAt(phons, currPh, 3, []phoneme{g, z, hh})
				if ok_gz && !ok_gzh {
					// Okay, the h is silent so include it now
					new = phonToAlphas{
						[]phoneme{
							g, z,
						},
						"xh",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "x") {
				// As in eXactly but the g phoneme must be followed by a z...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{g, z}, []phoneme{g, zh}); ok {
					new = phonToAlphas{
						[]phoneme{
							g, z,
						},
						"x",
					}
				}
				break
			}
			// Handling silent u which typically follows a g
			if stringAt(alphas, currAbc, 3, "gue") {
				// As in dialoGUE, and many other words...
				if currPh == len(phons)-1 {
					new = phonToAlphas{
						[]phoneme{
							g,
						},
						"gue",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{g, ah}, []phoneme{g, er}, []phoneme{g, eh}, []phoneme{g, y}, []phoneme{g, ih}); !ok {
					// So NOT as in beleaGUEred, beleaGUEred, GUEss, arGUE, vaGUEst,...
					// If we get here then the ue is not sounded
					new = phonToAlphas{
						[]phoneme{
							g,
						},
						"gue",
					}
					break
				} else if _, ok := phsAt(phons, currPh, 2, []phoneme{g, y}); !ok {
					// So NOT as in arGUE,...
					// In other cases the u is silent and the vowel phoneme can be
					// processed next time round the loop - I think anyway...
					new = phonToAlphas{
						[]phoneme{
							g,
						},
						"gu",
					}
					break
				}
			}
			// There's also a silent u in other non-gue words
			if stringAt(alphas, currAbc, 2, "gu") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{g, aa}, []phoneme{g, ae}, []phoneme{g, eh}, []phoneme{g, ay}); ok {
					// As in GUard, GUarantee, GUarantee, GUide,...
					new = phonToAlphas{
						[]phoneme{
							g,
						},
						"gu",
					}
					break
				}
			}
			// A silent h can also follow a g
			if stringAt(alphas, currAbc, 2, "gh") {
				// As in GHost,...
				// But check the h is silent first
				if _, ok := phsAt(phons, currPh, 2, []phoneme{g, hh}); !ok {
					new = phonToAlphas{
						[]phoneme{
							g,
						},
						"gh",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "gg") {
				new = phonToAlphas{
					[]phoneme{
						g,
					},
					"gg",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "g") {
				new = phonToAlphas{
					[]phoneme{
						g,
					},
					"g",
				}
				break
			}
		case gr:
			if s, ok := getStringAt(alphas, currAbc, 0, "gar", "ggr", "gr"); ok {
				// As in marGARet, aGGRegate, GRipe,...
				// TDOD: Don't like adding margaret here as it's a proper name and the only
				// example of this mapping I can find in the dictionary. Should margaret be
				// removed from the dictionary?
				new = phonToAlphas{
					[]phoneme{
						gr,
					},
					s,
				}
				break
			}
		case hh:
			if stringAt(alphas, currAbc, 2, "wh") {
				// As in WHo, WHole,...
				new = phonToAlphas{
					[]phoneme{
						hh,
					},
					"wh",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "h") {
				new = phonToAlphas{
					[]phoneme{
						hh,
					},
					"h",
				}
				break
			}
		case ih:
			// Deal with training 'ih first
			if alphas == "'ih" {
				new = phonToAlphas{
					[]phoneme{
						ih,
					},
					"'ih",
				}
				break
			}
			// Catch borrowed French words ending -ier
			if stringAt(alphas, currAbc, 3, "ier") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, ey}, []phoneme{ih, ehr}); ok {
					// As in atelIer, concIerge,...
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						"i",
					}
					break
				}
			}
			// Some more French stuff, this time -ez- in rendezvous
			if stringAt(alphas, currAbc, 2, "ez") {
				// Check the lexical 'z' isn't sounded though
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, z}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						"ez",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ei") {
				// As in nuclEI,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, ay}); ok {
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						"e",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 3, "hea") {
				new = phonToAlphas{
					[]phoneme{
						ih,
					},
					"hea",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ea", "ia", "ie"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, aa}, []phoneme{ih, ax}, []phoneme{ih, ey}, []phoneme{ih, ae}, []phoneme{ih, eh}, []phoneme{ih, ih}, []phoneme{ih, iy}); ok {
					// As in cavIAr, folIAte, nausEAte, avIAte, asIAtic, fIEsta, foolIAge, medIEval...
					// Don't swallow the a here leave it for the next phoneme
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						s[:1],
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oeo", "eo"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, oh}, []phoneme{ih, ow}); ok {
					// As in homOEopathy, gEology, apothEosis...
					// The 'o' is sounded so don't swallow the 'o' here (as we do below in thEOry)
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						s[:len(s) - 1],
					}
					break
				}
			}
			if _, ok := getStringAt(alphas, currAbc, 2, "ey"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, y}); ok {
					// The y is sounded so don't swallow it here
					// as in bEyond
					new = phonToAlphas{
						[]phoneme{
							eh,
						},
						"e",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ae", "ai", "ay", "ea", "ee", "ei", "eo", "ey", "ia", "ie", "ui"); ok {
				// As in archAEology, portrAIt, mondAY, EAr (x! TODO: This example is wrong!), bEEn, wherEIn, thEOry, convERsation, donkEYs, carrIAge, sIEve, bUIlding,...
				new = phonToAlphas{
					[]phoneme{
						ih,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "er"); ok {
				// As in ERupt,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ih,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "hi"); ok {
				// It looks like we have a silent h
				// As in exHIbit,...
				new = phonToAlphas{
					[]phoneme{
						ih,
					},
					s,
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "i", "o", "u", "y"); ok {
				// As in encourAgIng, Erupt, bIn ,wOmen, bUsily, abYss,...
				new = phonToAlphas{
					[]phoneme{
						ih,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "'s", "s"); ok {
				// As in James'S, bridgeS,... and many other plural nouns
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ih, z}); ok {
					new = phonToAlphas{
						[]phoneme{
							ih, z,
						},
						s,
					}
					break
				}
			}
		case ing:
			p := []phoneme{
				ing,
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "eng", "ing"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ing, g}); ok {
					// As in ENGland, lINGer,...
					// We should save the g for later
					new = phonToAlphas{
						p,
						s[:len(s) - 1],
					}
				} else {
					// As in waxING, and just about any -ing word...
					new = phonToAlphas{
						p,
						s,
					}
				}
				break
			}			
			if s, ok := getStringAt(alphas, currAbc, 0, "eing", "uing", "eng"); ok {
				// As in agEING, catalogUING, ENGland,...
				new = phonToAlphas{
					p,
					s,
				}
				break
			}			
			if s, ok := getStringAt(alphas, currAbc, 3, "inc", "ink", "inq", "inx", "ync", "ynx"); ok {
				// As in zINc, twINking, delINquent, sphINx, sYNchronous, pharYNx,...
				new = phonToAlphas{
					[]phoneme{
						ing,
					},
					s[:len(s) - 1],
				}
				break
			}			
		case iy:
			if stringAt(alphas, currAbc-1, 4, "peop") {
				// Treating pEOple as a special case here. There are other occurrences
				// of the letters eo which break otherwise, like stero
				new = phonToAlphas{
					[]phoneme{
						iy,
					},
					"eo",
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "aeo"); ok {
				// As in palAEontology...
				new = phonToAlphas{
					[]phoneme{
						iy,
					},
					"ae",
				}
				break
			}
			if strings.HasPrefix(alphas, "here") {
				// I'm being quite specific here. I don't want to break words like etHEREal.
				// As in words like HEREafter,...
				if p, ok := phsAt(phons, currPh, 3, []phoneme{iy, ax, r}); ok {
					new = phonToAlphas{
						p,
						"ere",
					}
					break
				}
			}
			// Catch -eying words early, else they get caught up in the patterns that follow
			if stringAt(alphas, currAbc, 5, "eying") {
				// As in curtsEYing,...
				new = phonToAlphas{
					[]phoneme{
						iy,
					},
					"ey",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ear", "eer"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{iy, axm}, []phoneme{iy, axn}); ok {
					// As in EARMark, shEERNess,...
					new = phonToAlphas{
						[]phoneme{
							iy,
						},
						s,
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 3, []phoneme{iy, ax, r}); ok {
					// As in EARRING,...
					new = phonToAlphas{
						[]phoneme{
							iy, ax,
						},
						"ea",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{iy, ax}); ok {
					// As in EAR,...
					new = phonToAlphas{
						[]phoneme{
							iy, ax,
						},
						"ear",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "eh") {
				// As in vEHicle,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{iy, hh}); !ok {
					// But not as in vEHicular,... - where the phonetic 'h' is sounded
					new = phonToAlphas{
						[]phoneme{
							iy,
						},
						"eh",
					}
					break
				}
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "ei"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{iy, ing}); ok {
					new = phonToAlphas{
						[]phoneme{
							iy,
						},
						"e",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ae", "ay", "ea", "ee", "ei", "ey", "ie", "oe"); ok {
				// As in palAEontology, sundAY, mEAt, lEEt, EIther, whiskEYs, pastrIEs, OEstrus,...
				// But not as in words like rEARM, rEAct, twentIEth, wobblIER, clEARLy, sturdIEr, recrEAte, rEIntegrate,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{iy, aa}, []phoneme{iy, ae}, []phoneme{iy, ax}, []phoneme{iy, axr}, []phoneme{iy, axl}, []phoneme{iy, eh}, []phoneme{iy, er}, []phoneme{iy, ey}, []phoneme{iy, ih}); !ok {
						new = phonToAlphas{
						[]phoneme{
							iy,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ae", "ea", "ee", "e"); ok {
				// There seems to be a pattern that 'e' followed by a sounded 'r' can
				// be sounded as iy ah
				// As in chimAEra, chEEr, bactEria,...
				if _, ok := phsAt(phons, currPh, 3, []phoneme{iy, ax, r}); ok {
					new = phonToAlphas{
						[]phoneme{
							iy, ax,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ae", "oe"); ok {
				// As in pAEan, diarrhOEa,...
				// if _, ok := phsAt(phons, currPh, 3, []phoneme{iy, ax, r}); ok {
				new = phonToAlphas{
					[]phoneme{
						iy, ax,
					},
					s,
				}
				break
				// }
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "e", "i", "y"); ok {
				// As in shE, shylY,...
				new = phonToAlphas{
					[]phoneme{
						iy,
					},
					s,
				}
				break
			}
		case jh:
			if s, ok := getStringAt(alphas, currAbc, 2, "ch", "dg", "di", "dj", "gg"); ok {
				// As in sandwiCH, heDGing, solDIer, aDJust, suGGest,...
				new = phonToAlphas{
					[]phoneme{
						jh,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "g", "j"); ok {
				// As in ginGer,... We don't want to swallow up the e in er
				new = phonToAlphas{
					[]phoneme{
						jh,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "de", "du") {
				// As in granDEur, eDUcate... But leave the u to be picked up by the following
				// phoneme
				new = phonToAlphas{
					[]phoneme{
						jh,
					},
					"d",
				}
				break
			}
		case k:
			if alphas == "ok" {
				new = phonToAlphas{
					[]phoneme{
						k, ey,
					},
					"k",
				}
				break
			}
			// The k phoneme can optionally appear in -ngth words according to Google
			// and they do appear in the CMU dictionary
			if stringAt(alphas, currAbc-2, 4, "ngth") {
				if _, ok := phsAt(phons, currPh-1, 3, []phoneme{ng, k, th}); ok {
					// As in lengTH,...
					new = phonToAlphas{
						[]phoneme{
							k, th,
						},
						"th",
					}
					break
				}
			}
			// Detecting words ending que...
			if stringAt(alphas, currAbc, 3, "que") && currPh == len(phons)-1 {
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					"que",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "kh"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, hh}); !ok {
					// It looks like the h is silent so swallow it now.
					new = phonToAlphas{
						[]phoneme{
							k,
						},
						s,
					}
					break
				}
			}
			// Catch a possible silent l
			if stringAt(alphas, currAbc, 2, "lk") {
				// As in waLK, foLK,...
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					"lk",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ct") {
				// But make sure the lexical 't' isn't sounded - and 't' isn't always sound as t!
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, ch}, []phoneme{k, sh}, []phoneme{k, t}, []phoneme{k, tr}); !ok {
					// As in aCTuary, traCTion, striCTly, buTTRess,...
					new = phonToAlphas{
						[]phoneme{
							k,
						},
						"ct",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "cqu", "qu"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, w}); ok {
					// As aCQUire, QUoth and all sorts of words containing qu...
					new = phonToAlphas{
						[]phoneme{
							k, w,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "cqu"); ok {
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					s,
				}
				break
			}		
			// Check for q on it's own. These will be in borrowed words or foreign placenames
			if stringAt(alphas, currAbc, 1, "q") {
				// As in Qatar,...
				// new = phonToAlphas{
				// 	[]phoneme{
				// 		k,
				// 	},
				// 	"q",
				// }
				// break
			}
			if stringAt(alphas, currAbc, 4, "xion") {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{k, sh, n}); ok {
					new = phonToAlphas{
						[]phoneme{
							k, sh, n,
						},
						"xion",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ckc", "kc"); ok {
				// As in blaCKCurrant,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, k}, []phoneme{k, ch}, []phoneme{k, kl}, []phoneme{k, kr}); !ok {
					// But not as in blaCKcurrant, baCKchat, saCKcloth, coCKcrow...
					new = phonToAlphas{
						[]phoneme{
							k,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "cch", "ch", "ck", "qu"); ok {
				// As in saCCHarine, CHoir (or loCH approximately), chiCKen, QUay,...
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "kk"); ok {
				p := []phoneme{k}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, k}); ok {
					// As in treKKing,...
					// There's a double phonetic 'k' though, so don't grab both
					// lexical 'k's here
					s = "k"
				}
				new = phonToAlphas{
					p,
					s,
				}
				break
			}
			// Trapping words like exchequer, else the 'c' gets swallowed by
			// processing of 'xc' further on
			if stringAt(alphas, currAbc, 3, "xch") {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{k, s, ch}); ok {
					// As in eXchange
					new = phonToAlphas{
						[]phoneme{
							k, s,
						},
						"x",
					}
					break
				}
			}
			if str, ok := getStringAt(alphas, currAbc, 2, "cc", "xc"); ok {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{k, s, ch}, []phoneme{k, s, k}, []phoneme{k, s, kl}, []phoneme{k, s, kr}); ok {
					// As in eXCHange, eXCoriate, eXClude,, eXCRete...
					// Save the (second) k (or kl, kr) for later...
					new = phonToAlphas{
						[]phoneme{
							k, s,
						},
						str[:1],
					}
					break
				}
				// We need to be careful here and check the phonemes for k, s
				if _, ok := phsAt(phons, currPh, 2, []phoneme{k, s}); ok {
					// As in suCCess, eXCited,...
					new = phonToAlphas{
						[]phoneme{
							k, s,
						},
						str,
					}
					break
				}
				// As in suCCour,...
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					str,
				}
				break
			}
			if str, ok := getStringAt(alphas, currAbc, 2, "xs"); ok {
				// As in coXSwain,...
				new = phonToAlphas{
					[]phoneme{
						k, s,
					},
					str,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "c", "k", "q"); ok {
				// Note the 'q' here
				// As in ...,Qatar,...
				new = phonToAlphas{
					[]phoneme{
						k,
					},
					s,
				}
				break
			}
			// Trying to trap the silent h that can follow ex-
			if stringAt(alphas, currAbc, 2, "xh") {
				// As in eXHibition,...
				_, ok_gz := phsAt(phons, currPh, 2, []phoneme{k, s})
				_, ok_gzh := phsAt(phons, currPh, 3, []phoneme{k, s, hh})
				if ok_gz && !ok_gzh {
					// Okay, the h is silent so include it now
					new = phonToAlphas{
						[]phoneme{
							k, s,
						},
						"xh",
					}
					break
				}
			}
			//The 't' following x is sometimes not pronounced so catch it here
			if _, ok := getStringAt(alphas, currAbc, 0, "xtb"); ok {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{k, s, b}); ok {
					// The 't' is silent so swallow it now
					// As in teXTbook,... (I think this and its plural are the only examples)
					new = phonToAlphas{
						[]phoneme{
							k, s,
						},
						"xt",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "x") {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{k, s}, []phoneme{k, sh}); ok {
				// As in eXit, anXious,...
				new = phonToAlphas{
						p,
						"x",
					}
				}
			}
		case kl:
			if s, ok := getStringAt(alphas, currAbc, 0, "call", "ckcl", "quel", "ccl", "chl", "ckl", "col", "cul", "kel", "khl", "ctl", "cl", "kl"); ok {
				// As in dramatiCALLy, saCKCLoth, uniQUELy, aCCLaim, CHLorine, tiCKLer, choCOLate, faCULty, liKELy, neKHLudoff (seriously?
				// -it's a training word), striCTLy, CLock, KLingon,...
				new = phonToAlphas{
					[]phoneme{
						kl,
					},
					s,
				}
				break
			}
		case kr:
			if s, ok := getStringAt(alphas, currAbc, 0, "ckcr", "cker", "ccr", "chr", "ckr", "cr", "kr"); ok {
				// As in coCKCRow, ..., aCCRete, laCHRymose, coCKRoach, maCKERel, CRisis, sauerKRaut,...
				new = phonToAlphas{
					[]phoneme{
						kr,
					},
					s,
				}
				break
			}
		case l:
			new = new.mapL(
				phons, currPh, alphas, currAbc,
			)
			if !new.equal(phonToAlphas{}) {
				break
			}
			// Check for borrowed Italian words like seraglio
			if s, ok := getStringAt(alphas, currAbc, 0, "gl"); ok {
				// As in intaGLio,...
				new = phonToAlphas{
					[]phoneme{
						l,
					},
					s,
				}
				break
			}
			// Check for deleted schwa in -ally, -ully word endings
			if s, ok := getStringAt(alphas, currAbc, 0, "all", "ial", "ill", "oll", "ull"); ok {
				// As in principALLy, specIAL, pencILLed, gambOLLing, wonderfULLy,...
				new = phonToAlphas{
					[]phoneme{
						l,
					},
					s,
				}
				break
			}
			// Check for a deleted schwa, for instance in a(ae)n(n)i(ih)m(m)al(l).
			// This should be dealt with more generically (ie.e not just for the l
			// phoneme but for now this will catch some of the more common examples
			if s, ok := getStringAt(alphas, currAbc, 2, "al", "el", "il", "ol", "ul"); ok {
				// As in animAL, squirrEL, councIL, cathOLic, facULty,...
				new = phonToAlphas{
					[]phoneme{
						l,
					},
					s,
				}
				break
			}
		case m:
			// A couple of abbreviations whcih can't easily be mapped phoneme
			// by phoneme
			if alphas == "mr" || alphas == "mrs" {
				new = phonToAlphas{
					phons,
					alphas,
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 3, "med", "mes"); ok {
				// Don't want to swallow these characters unless they're at the end of
				// a word. Think of MEDical, tiMEShift,...
				if currAbc+3 == len(alphas) {
					// As in tiMED, liMES,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"me",
					}
					break
				}
			}
			// Treating dumbbell as special case. The first b should really be
			// attached to the m phoneme
			if alphas == "dumbbell" {
				new = phonToAlphas{
					[]phoneme{
						m,
					},
					"mb",
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 0 , "mem"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, m}); ok {
					// As in hoMEmade,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"me",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 3, []phoneme{m, ax, m}, []phoneme{m, ih, m}, []phoneme{m, eh, m}); !ok {
					// Not as in MEMber, MEMento, imMEMorial,...
					// But as in hoMEMade,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mem",
					}
					break
				}
			}
			if _, ok := getStringAt(alphas, currAbc, 0 , "mbm"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, m}); ok {
					// As in entoMBment,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mb",
					}
					break
				} else {
					// As in entoMBMent,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mbm",
					}
					break
				}
			}
			// Catch a possible mbl, before we test for mb
			if _, ok := getStringAt(alphas, currAbc, 0, "mboll", "mbl"); ok {
				// As in gamBOLLed, tuMBLer,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, bl}); ok {
					// As in raMBLer,... but leave the bl for later
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"m",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "mb") {
				// As in cliMber,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, b}); !ok {
					// Not as in claMBer,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mb",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "mn") {
				// As in hyMN,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, n}); !ok {
					// Not as in hyMNal,...
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mn",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "mm") {
				// As in diMMer,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, m}); !ok {
					// As in rooMMate,...
					// There's a double phonetic 'm' though, so don't grab both
					// lexical 'm's here
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						"mm",
					}
					break
				}
			}
			// Catch a possible silent g, or l
			if s, ok := getStringAt(alphas, currAbc, 2, "gm", "lm"); ok {
				// As in diaphraGM, caLM,...
				new = phonToAlphas{
					[]phoneme{
						m,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "mp"); ok {
				// As in redeMPtion,...
				// Check for a silent p though. A lexical p can also be sounded as f (as in eMPhatic)!
				if _, ok := phsAt(phons, currPh, 2, []phoneme{m, f}, []phoneme{m, p}, []phoneme{m, pl}, []phoneme{m, pr}); !ok {
					// The lexical 'p' is silent so grab it now.
					new = phonToAlphas{
						[]phoneme{
							m,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "m") {
				new = phonToAlphas{
					[]phoneme{
						m,
					},
					"m",
				}
				break
			}
			// In the word drachm, the ch is silent
			if stringAt(alphas, currAbc, 3, "chm") {
				new = phonToAlphas{
					[]phoneme{
						m,
					},
					"chm",
				}
			}
		case n:
			if s, ok := getStringAt(alphas, currAbc, 2, "ln", "mn"); ok {
				// As in lincoLN, mnemonic, ...
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "nn") {
				// As in naNNy,...
				s := "nn"
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, n}); ok {
					// But not as in uNnatural,...
					// There's a double-n phoneme so save the other one for the other
					// lexical 'n'
					s = "n"
				}
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "nkn"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, n}, []phoneme{n, k}); ok {
					// As in uNknown, baNknote,...
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						"n",
					}
				} else {
					// Asin uNKNown,...
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						s,
					}
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "nen"); ok {
				// This could be represented by one phoneme as in oNENess but
				// we need to be careful
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, n}); ok {
					// As in oNeness,...
					// But not as in oNENess,...
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						"n",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 3, []phoneme{n, ax, n}, []phoneme{n, eh, n}, []phoneme{n, ih, n}); ok {
					// Not as in oppoNENt, uNENding, lINEN,...
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						"n",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, axn}); !ok {
					// Not as in oppoNENt,...
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						s,
					}
					break
				}
			}
			// Catch the borrowed 'ñ' (n-yah)
			if s, ok := getStringAt(alphas, currAbc, 2, "gn", "nh"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, y}); ok {
					// As in siGNor, piraNHa,...
					new = phonToAlphas{
						[]phoneme{
							n, y,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "gnn"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, n}); !ok {
					// As in foreiGNNess
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						s,
					}
					break
				}
			}

			if s, ok := getStringAt(alphas, currAbc, 2, "dn", "gn", "kn", "mp", "pn"); ok {
				// As in weDNesday, siGN, KNee, coMPtroller, PNeumatic,...
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "nd") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, jh}); ok {
					// As in graNdeur,...
					// Leave the d to be processed later
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						"n",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{n, d}, []phoneme{n, dz}); !ok {
					// Looks like the d is not sounded as in laNDs...
					// So swallow it now
					new = phonToAlphas{
						[]phoneme{
							n,
						},
						"nd",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 4, "wain") {
				// As in coxswain,...
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					"wain",
				}
				break
			}
			// Check for deleted schwa
			if s, ok := getStringAt(alphas, currAbc, 0, "ain", "an", "ern", "ian", "ion", "ten"); ok {
				// As in certAIN, tartAN, govERNment, alsatIAN, fashIONable, sofTEN...
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "en", "in", "on"); ok {
				// As in christEN, basIN, arsON,...
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "n") {
				new = phonToAlphas{
					[]phoneme{
						n,
					},
					"n",
				}
				break
			}
		case ng:
			if stringAt(alphas, currAbc, 4, "ngue") {
				// As in toNGUE,...
				new = phonToAlphas{
					[]phoneme{
						ng,
					},
					"ngue",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "nc") {
				// As in the place name, Altrincham, which Wikipedia confirms is pronounced with a
				// phonetic 'ng'
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ng, k}, []phoneme{ng, kr}); !ok {
					// But not as in acupuNCture, paNCReas...
					new = phonToAlphas{
						[]phoneme{
							ng,
						},
						"ch",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "nd") {
				// As in haNDkerchief,...
				new = phonToAlphas{
					[]phoneme{
						ng,
					},
					"nd",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ng") {
				// If next phoneme is gr don't swallow the g here
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ng, gr}); ok {
					// As in coNgregate,...
					new = phonToAlphas{
						[]phoneme{
							ng,
						},
						"n",
					}
					break
				}
				// If next phoneme is g then both phonemes map to the letters ng
				// As in aNGuish,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ng, g}); ok {
					new = phonToAlphas{
						[]phoneme{
							ng,
						},
						"n",
					}
				} else {
					// As in wiNG,...
					new = phonToAlphas{
						[]phoneme{
							ng,
						},
						"ng",
					}
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "n") {
				// As in think...
				new = phonToAlphas{
					[]phoneme{
						ng,
					},
					"n",
				}
				break
			}
		case ow:
			if s, ok := getStringAt(alphas, currAbc, 0, "ough", "aoh"); ok {
				// As in pharAOH, furlOUGH,... at least with a US accent
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					s,
				}
				break
			}
			if stringAt(alphas, 0, 6, "brooch") {
				// I think this is about the only word in the English language with oo
				// represented by the phoneme ow. In other oo words the oo is sounded
				// as two separate vowels, as in cooperate, microorganism,...
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					"oo",
				}
				break
			}
			// Borrowed French words ending -eau, -eaus, -eaux. Do we really
			// pluralise borrowed French words with an s?
			if s, ok := getStringAt(alphas, currAbc, 0, "eaus", "eaux", "eau"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ow, z}); ok {
					// Don't swallow the s, or x as it's sounded
					s = "eau"
				}
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "au", "oa", "oe", "ou"); ok {
				// As in chAUvanist, cOAt, tOE, sOUlful,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ow, aa}, []phoneme{ow, ae}, []phoneme{ow, ao}, []phoneme{ow, ax}, []phoneme{ow, axr}, []phoneme{ow, eh}, []phoneme{ow, ey}, []phoneme{ow, ih}, []phoneme{ow, iy}); !ok {
					// But not as in kOAla, radiOActive, cOAuthor, bOA, bOA, whosOEver, OAsis, borrOWIng, micrOElectronics...
					new = phonToAlphas{
						[]phoneme{
							ow,
						},
						s,
					}
					break
				}
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "oww"); ok {
				// As in glOWWorm,...
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					"ow",
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "ow"); ok {
				// As in micrOwave,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ow, w}); ok {
					// The lexical 'w' is sounded so leave it for the phonetic 'w'
					new = phonToAlphas{
						[]phoneme{
							ow,
						},
						"o",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "eo", "ew", "ow"); ok {
				// As in yEOman, bOW,...
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "oh") {
				// As in OHm,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ow, hh}); !ok {
					// But not as in bOHemia,...
					new = phonToAlphas{
						[]phoneme{
							ow,
						},
						"oh",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ot") {
					// As in depot,...
					if _, ok := phsAt(phons, currPh, 2, []phoneme{ow, t}, []phoneme{ow, dh}, []phoneme{ow, sh}, []phoneme{ow, th}, []phoneme{ow, tr}); !ok {
					// But not as in rOTe, bOTH, pOTion, clOTHE, synchrOTRon,...
					// So it looks like the t is silent here
					new = phonToAlphas{
						[]phoneme{
							ow,
						},
						"ot",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "o") {
				// As in Over,...
				new = phonToAlphas{
					[]phoneme{
						ow,
					},
					"o",
				}
				break
			}
		case oy:
			if stringAt(alphas, currAbc, 3, "uoy") {
				new = phonToAlphas{
					[]phoneme{
						oy,
					},
					"uoy",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "oi", "oy"); ok {
				// As in chOIce, tOY...
				new = phonToAlphas{
					[]phoneme{
						oy,
					},
					s,
				}
				break
			}
		case p:
			if stringAt(alphas, currAbc, 2, "pp") {
				new = phonToAlphas{
					[]phoneme{
						p,
					},
					"pp",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "bp") {
				// As in subpoena,...
				new = phonToAlphas{
					[]phoneme{
						p,
					},
					"bp",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "pt") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{p, ch}, []phoneme{p, sh}, []phoneme{p, t}, []phoneme{p, tr}, []phoneme{p, th}); !ok {
					// It's not caPTure, descriPTIon, comPTRoller, uPTurn, upTHrust,...
					// So it looks like there's a silent t in words like bankruptcy
					new = phonToAlphas{
						[]phoneme{
							p,
						},
						"pt",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "p") {
				new = phonToAlphas{
					[]phoneme{
						p,
					},
					"p",
				}
				break
			}
		case pl:
			if s, ok := getStringAt(alphas, currAbc, 0, "pall", "p-l", "pel", "ppl", "pul", "pl"); ok {
				// As in princiPALLy, hooP-La, hoPELess, suPPLy, sePULchre, couPLet,...
				new = phonToAlphas{
					[]phoneme{
						pl,
					},
					s,
				}
				break
			}
		case pr:
			if s, ok := getStringAt(alphas, currAbc, 0, "par", "per", "pir", "por", "ppr", "pr"); ok {
				// as in comPARably, temPERature, asPIRin, contemPORary, aPPRaise, PRove,...
				new = phonToAlphas{
					[]phoneme{
						pr,
					},
					s,
				}
				break
			}
		case r:

			if s, ok := getStringAt(alphas, currAbc, 2, "aur", "our", "rrh", "ar", "ir", "or", "rh", "rr", "ur", "wr"); ok {
				// As in restAURant, labOURer, ciRRHosis, solitARy, aspIRin, conservatORy, RHythm, eRRor, natURal, WRite,...
				// But NOT as in speaRHead,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{r, hh}); !ok {
					new = phonToAlphas{
						[]phoneme{
							r,
						},
						s,
					}
					break
				}
			}
			// Is there a better way of handling these 'special' cases?
			//
			// Handle forehead(s) as a special case
			if strings.HasPrefix(alphas, "forehead") {
				new = phonToAlphas{
					[]phoneme{
						r,
					},
					"re",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "r") {
				new = phonToAlphas{
					[]phoneme{
						r,
					},
					"r",
				}
				break
			}
		case s:
			// Treating this a special case. There are spellings of conversation in
			// the dictionary which drop the er phoneme
			if stringAt(alphas, currAbc, 3, "ers") {
				new = phonToAlphas{
					[]phoneme{
						s,
					},
					"ers",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ces") {
				// We have to be careful here. Thre are lots of options
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, ax}, []phoneme{s, eh}, []phoneme{s, ih}, []phoneme{s, iy}, []phoneme{s, sh}, []phoneme{s, s}); !ok {
					// So, not as in neCEssary, anCEstry, spiCEs, faeCEs, spaCESHip, iCESkate,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"ces",
					}
					break
				}
			}
			// Focussing on postscript(s) which can be rendered P OH S K R IH P T (S)
			if stringAt(alphas, currAbc, 5, "stscr") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, k}, []phoneme{s, kr}); ok {
					// As in poSTScript,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"sts",
					}
					break
				}
			}
			// breaststroke is a right mare to parse! The -ststr- can be represented phonetically as
			// S T S T R, S S T R, S T S tR, S tS T R, S S tR, S T R, S tS tR, S tR.
			if _, ok := getStringAt(alphas, currAbc, 0, "sts"); ok {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{s, t, r}); ok {
					// As in breaSTStroke,...
					// By looking at three phonemes we know tha the t being sounded
					// is the second t.
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"sts",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, tr}); ok {
					// As in breaSTStroke,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"sts",
					}
					break
				}
			}
			// if str, ok := getStringAt(alphas, currAbc, 0, "sts"); ok {
			// 	if _, ok := phsAt(phons, currPh, 2, []phoneme{s, t}, []phoneme{s, s}); !ok {
			// 		// The t doesn't appear to be sounded,
			// 		new = phonToAlphas{
			// 			[]phoneme{
			// 				s,
			// 			},
			// 			str,
			// 		}

			// 	}
			// }
			if stringAt(alphas, currAbc, 3, "sth") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, th}, []phoneme{s, t}); !ok {
					// There is a silent th so swallow it now
					// As in iSTHmus,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"sth",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "st") {
				// Look out for a silent 't'. It typically follows an 's'.
				// if _, ok := phsAt(phons, currPh, 2, []phoneme{s, ch}, []phoneme{s, t}, []phoneme{s, th}, []phoneme{s, tr}, []phoneme{s, ts}); !ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, ax}, []phoneme{s, axl}, []phoneme{s, axn}, []phoneme{s, k}, []phoneme{s, l}, []phoneme{s, m}, []phoneme{s, n}, []phoneme{s, p}, []phoneme{s, s}); ok {
					// As in caSTle, caSTle, muSTn't, waiSTcoat, neSTle, adjuSTment, cheSTnut, poSTpone, breaSTstroke,...
						// There is a silent t
					// As in thiSTle,... but not as in reStless,..
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"st",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 4, "sten") {
				// As in faSTEN,... The t here is typically silent so include it
				// here
				// But be careful, we don't want to swallow the t in sTencil...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, t}); !ok {
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"st",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "tz") {
				// As in walTZ,...
				new = phonToAlphas{
					[]phoneme{
						s,
					},
					"tz",
				}
				break
			}
			if stringAt(alphas, currAbc-1, 3, "nds") {
				if _, ok := phsAt(phons, currPh-1, 2, []phoneme{n, s}); ok {
					// We've found a silent (or suppressed) 'd' as in wiNDSwept, wiNDSor...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"ds",
					}
					break
				}
			}
			if stringAt(alphas, currAbc-1, 3, "cts") {
				if _, ok := phsAt(phons, currPh-1, 2, []phoneme{k, s}); ok {
					// We've found a silent (or suppressed) 't' as in reflecTs,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"ts",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ic") {
				// As in the posh pronunciation of medICine,...
				new = phonToAlphas{
					[]phoneme{
						s,
					},
					"ic",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "sc") {
				// As in SCene,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, ch}, []phoneme{s, k}, []phoneme{s, kr}, []phoneme{s, kl}); !ok {
					// But not as in miSChief, SChool, SCrum, diSClose,...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"sc",
					}
				} else {
					// As in SCHism, aSCRibe, diSCLaimer...
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"s",
					}
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ps") {
				// As in PSychic,...
				new = phonToAlphas{
					[]phoneme{
						s,
					},
					"ps",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ss") {
				// As in paSS,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{s, s}, []phoneme{s, sh}); !ok {
					// But not as in miSSpell, miSSHapen,... and many other mis- words (in which we leave the second
					// lexical 's' for the second phonetic 's')
					new = phonToAlphas{
						[]phoneme{
							s,
						},
						"ss",
					}
					break
				}
			}
			if st, ok := getStringAt(alphas, currAbc, 1, "c", "s", "z"); ok {
				// As in truCe, whiSt, glitZy,...
				new = phonToAlphas{
					[]phoneme{
						s,
					},
					st,
				}
				break
			}
		case sh:
			if s, ok := getStringAt(alphas, currAbc, 4, "cian", "sion", "tion"); ok {
				// As in electriCIAN, penSION, naTION,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{sh, n}); ok {
					new = phonToAlphas{
						[]phoneme{
							sh, n,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "cesh", "sesh"); ok {
				// As in apprentiCESHip, horSESHoe,...
				new = phonToAlphas{
					[]phoneme{
						sh,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 5, "ssion"); ok {
				// As in percuSSION,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{sh, n}); ok {
					new = phonToAlphas{
						[]phoneme{
							sh, n,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "sch", "ssi"); ok {
				// As in SCHnapps, discuSSIon,...
				new = phonToAlphas{
					[]phoneme{
						sh,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ci", "ti"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{sh, ih}, []phoneme{sh, iy}, []phoneme{sh, y}); ok {
					// Only swallow the c or t in this case because the following i is
					// sounded separately.
					// As in appreCiate, negoTiable, assoCiative,...
					new = phonToAlphas{
						[]phoneme{
							sh,
						},
						s[:1],
					}
					break
				}
			}
			// Catch threshold, where the h is pronounced...
			// But be careful. There are words like fishhok where the h is pronounced
			// and in these words we do want to swallow the lexical sh
			_, shh := getStringAt(alphas, currAbc, 0, "shh")
			if _, ok := getStringAt(alphas, currAbc, 0, "sh"); ok && !shh {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{sh, hh}); ok {
					new = phonToAlphas{
						[]phoneme{
							sh,
						},
						"s",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ch", "ci", "sc", "sh", "ss", "ti"); ok {
				// As in CHagrin, vivaCIous, conSCience, SHed, seSSion, raTIon,...
				new = phonToAlphas{
					[]phoneme{
						sh,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "c", "s", "x"); ok {
				// As in oCeanography, Sugar, anXious...
				new = phonToAlphas{
					[]phoneme{
						sh,
					},
					s,
				}
				break
			}
		case t:
			// A training pronunciation to get out of the way first
			if stringAt(alphas, currAbc, 1, "t") {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{t, ch}); ok {
					new = phonToAlphas{
						p,
						"t",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ts") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{t, s}, []phoneme{t, sh}); !ok {
					// A special case, the lexical 's' isn't sounded, as in TSetse,...
					new = phonToAlphas{
						[]phoneme{
							t,
						},
						"ts",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ed") {
				// As in askED,...
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					"ed",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ort") {
				// As in comfORTable,...
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					"ort",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "th") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{t, th}); ok {
					// A special case, as in eighTH,...
					// The phonetic t is not represented in the lexical spelling. I'm not
					// aware of any other examples.
					new = phonToAlphas{
						[]phoneme{
							t, th,
						},
						"th",
					}
					break
				}
				// Some words sound 'th' as t like THyme
				if _, ok := phsAt(phons, currPh, 2, []phoneme{t, hh}); !ok {
					// As in discoTHeque, but not as in poTHole,...
					new = phonToAlphas{
						[]phoneme{
							t,
						},
						"th",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 3, "ght") {
				// As in wriGHT,...
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					"ght",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "cht") {
				// As in yachT...
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					"cht",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "tt") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{t, t}, []phoneme{t, tr}); ok {
					// As in posTTraumatic, posTTRaumatic,...
					// Only swallow one t if the next phoneme is also a t (or is a tr)
					new = phonToAlphas{
						[]phoneme{
							t,
						},
						"t",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{t, th}); ok {
					// Leave the second t for the 'th' phoneme, as in cuTthroat,...
					new = phonToAlphas{
						[]phoneme{
							t,
						},
						"t",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "bt", "ct", "pt", "tt"); ok {
				// As in suBTle, indiCT, emPTy, boTTle,...
				// But only
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					s,
				}
				break
			}
			if st, ok := getStringAt(alphas, currAbc, 0, "zz", "z"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{t, s}); ok {
					// As in schmalZ, piZZa...
					new = phonToAlphas{
						p,
						st,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "d", "t"); ok {
				// As in learneD, lookeD,...
				new = phonToAlphas{
					[]phoneme{
						t,
					},
					s,
				}
				break
			}
		case th:
			if s, ok := getStringAt(alphas, currAbc, 0, "dth", "fth", "tth", "th"); ok {
				// As in thousanDTH, fiFTH, maTTHew, wealTH,...
				new = phonToAlphas{
					[]phoneme{
						th,
					},
					s,
				}
				break
			}
		case tr:
			if s, ok := getStringAt(alphas, currAbc, 0, "taur", "tar", "ter", "tor", "ttr", "ptr", "tr"); ok {
				// As in resTAURant, planeTARy, cemeTERy, hisTORy, aTTRibute, temPTRess, sTRike,...
				new = phonToAlphas{
					[]phoneme{
						tr,
					},
					s,
				}
				break
			}
		case ts:
			new = phonToAlphas{
				[]phoneme{
					ts,
				},
				"ts",
			}
		case uh:
			if s, ok := getStringAt(alphas, currAbc, 3, "our"); ok {
				// Catch words ending in -our (the axr phoneme is only used at
				//  the end of a word)
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uh, axr}); ok {
					// As in velOUR,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				// As in gOURmand,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uh, axm}); ok {
					new = phonToAlphas{
						p[:len(p)-1],
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "oor"); ok {
				// As in mOORland,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uh, axl}); ok {
					new = phonToAlphas{
						p[:len(p)-1],
						s,
					}
					break
				}
			} // Catch an uh phoneme transitioning to r
			// if s, ok := getStringAt(alphas, currAbc, 3, "ewer", "oor", "our"); ok {
			if s, ok := getStringAt(alphas, currAbc, 3, "eur", "ewer", "oor", "our"); ok {
				if p, ok := phsAt(phons, currPh, 3, []phoneme{uh, ax, r}); ok {
					// As in plEUral, brEWEry, mOOrish, tOUring,...
					new = phonToAlphas{
						p[:len(p)-1],
						s[:len(s)-1],
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uh, axn}); ok {
					// As in pOORness, tOURniquet,...
					new = phonToAlphas{
						[]phoneme{
							uh,
						},
						s,
					}
					break
				}
				// As in brEWER, mOOR, tOUR,...
				new = phonToAlphas{
					[]phoneme{
						uh, ax,
					},
					s,
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "ure"); ok {
				// Check for a sounded phonetic 'r'
				if p, ok := phsAt(phons, currPh, 3, []phoneme{uh, ax, r}); ok {
					// As in assUredly,... - so don't swallow the 'r' here
					new = phonToAlphas{
						p[:len(p)-1],
						"u",
					}
					break
				}
				// Check for 'axm'
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uh, axm}); ok {
					// As in allUREMent,... -
					new = phonToAlphas{
						// Leave the axm for later
						p[:len(p)-1],
						s,
					}
					break
				}
				// As in allURE,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uh, ax}); ok {
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ar") {
				// As in onwARd,...
				new = phonToAlphas{
					[]phoneme{
						uh,
					},
					"ar",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "oul") {
				// As in wOULd,...
				new = phonToAlphas{
					[]phoneme{
						uh,
					},
					"oul",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ou") {
				// As in shOUld,...
				new = phonToAlphas{
					[]phoneme{
						uh,
					},
					"ou",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "oo") {
				// As in wOOl,...
				new = phonToAlphas{
					[]phoneme{
						uh,
					},
					"oo",
				}
				break
			}
			// Transitioning from u to r is often rendered phonetically as uh ax r
			// so trying to capture that here
			if stringAt(alphas, currAbc, 1, "u") {
				// As in allUring,...
				if p, ok := phsAt(phons, currPh, 3, []phoneme{uh, ax, r}); ok {
					new = phonToAlphas{
						// If the r is sounded process it separately
						p[:len(p)-1],
						"u",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "o", "u"); ok {
				// As in wOlf, fUll,...
				new = phonToAlphas{
					[]phoneme{
						uh,
					},
					s,
				}
				break
			}
		case uw:
			if stringAt(alphas, currAbc, 4, "ough") {
				// As in thrOUGH,...
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					"ough",
				}
				break
			}
			// Some borrowed French words to deal with first
			if stringAt(alphas, currAbc, 3, "hou") {
				// As in silHOUette,...
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					"hou",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oup"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, p}); !ok {
					// The p is not sounded so swallow it here
					// As in the French cOUP,...
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 3, "oeu") {
				// As in manOEUvre,...
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					"oeu",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ieu") {
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					"ieu",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ous") {
				// As in rendezvOUS,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, s}, []phoneme{uw, z}); !ok {
					// The lexical 's' isn't sounded, so swallow it now
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"ous",
					}
					break
				}
			}
			if _, ok := getStringAt(alphas, currAbc, 4, "uest", "uism"); ok {
				// As in trUest, trUism,...
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					"u",
				}
				// Don't let this drop through to the next if because "ui" as in jUIce
				// will swallow the "i"
				break
			}
			if stringAt(alphas, currAbc, 2, "oo") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, oh}); ok {
					// As in zOOlogy,...
					// Only swallow the first 'o'. The second 'o' belongs to the phonetic oh
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"o",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ue") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, ax}, []phoneme{uw, axl}, []phoneme{uw, ih}); ok {
					// As in grUelling, grUelling, sUEt,...
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"u",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ew") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, w}); ok {
					// As in sEwerage,...,...
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"e",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ew", "eu", "oe", "oo", "ou", "ue", "wo"); ok {
				// As in nEW, slEUth, shOE, fOOl, sOUp, blUE, tWOsome...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, eh}); !ok {
					// But not as in whOEver,...
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						s,
					}
					break
				}
			}
			// Handle fluid separately else it screws up other uw words
			if stringAt(alphas, currAbc, 2, "ui") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uw, ing}); ok {
					// Leave the lexical 'i' to be processed as part of the ing later
					// As in constrUing, and many others...
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"u",
					}
					break
				}
				// As in jUIce,...
				// But NOT as in flUId,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{uw, ih}, []phoneme{uw, ah}); !ok {
					new = phonToAlphas{
						[]phoneme{
							uw,
						},
						"ui",
					}
					break
				} else {
					new = phonToAlphas{
						p,
						"ui",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "o", "u"); ok {
				// As in whO, tUne...
				new = phonToAlphas{
					[]phoneme{
						uw,
					},
					s,
				}
				break
			}
		case uwm:
			if s, ok := getStringAt(alphas, currAbc, 0, "ombm", "oomm"); ok {
				// We have pronunciation variants with uwm and uwm, m phonemes
				p := []phoneme{
					uwm,
		 		}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{uwm, m}); ok {
					// Leave the phonetic m for procesing later
					// As in entOMBment, rOOMmate...
					new = phonToAlphas{
						p,
						"omb",
					}
				} else {
					// As in entOMBMent, rOOMMate...
					new = phonToAlphas{
						p,
						s,
					}
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ulme", "ewm", "eum", "oom", "omb", "oum", "om", "um"); ok {
				// As in levenshULME, crEWMan, rhEUMatism, rOOM, khartOUM, wOMB, whOMsoever, costUMe,...
				new = phonToAlphas{
					[]phoneme{
						uwm,
					},
					s,
				}
				break
			}
		case uwn:
			if s, ok := getStringAt(alphas, currAbc, 0, "uen", "ewn", "oon", "oun", "on", "un"); ok {
				// As in blUENess, strEWN, pontOON, wOUNded, cantON*, fortUNe,...
				// *This pronunciation looks suspicious. TODO: Check to see if
				// this should be removed from the dictionary
				new = phonToAlphas{
					[]phoneme{
						uwn,
					},
					s,
				}
				break
			}
		case v:
			// pocketsphinx training means it sometimes swallows a vowel sound
			// following the v which leads to spellings like conversation where the
			// ...vers... is repesented phonetically as ... v s ...
			if stringAt(alphas, currAbc, 3, "ver") {
				// A very specific test for now for the case conversation
				if _, ok := phsAt(phons, currPh, 2, []phoneme{v, s}); ok {
					new = phonToAlphas{
						[]phoneme{
							v,
						},
						"ver",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ph", "vv"); ok {
				// As in nePHew, saVVy,...
				new = phonToAlphas{
					[]phoneme{
						v,
					},
					s,
				}
				break
			}
			// Catch a possible silent l
			if stringAt(alphas, currAbc, 2, "lv") {
				// As in haLVe,...
				new = phonToAlphas{
					[]phoneme{
						v,
					},
					"lv",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "f", "v"); ok {
				new = phonToAlphas{
					[]phoneme{
						v,
					},
					s,
				}
				break
			}
		case w:
			if stringAt(alphas, currAbc, 2, "wh") {
				new = phonToAlphas{
					[]phoneme{
						w,
					},
					"wh",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "w", "u"); ok {
				// As in Wall, langUid...
				new = phonToAlphas{
					[]phoneme{
						w,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oirm"); ok {
				if p, ok := phsAt(phons, currPh, 3, []phoneme{w, ay, axm}); ok {
					// As in chOIRMaster,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oir"); ok {
				// Check for the borrowed French oir
				if p, ok := phsAt(phons, currPh, 2, []phoneme{w, aa}); ok {
					if _, ok := phsAt(phons, currPh, 3, []phoneme{w, aa, r}); ok {
						// We have an r phoneme, it's the last phoneme in the word
						// As in boudOIR,...
						if len(phons) == currPh + 3 {
							p = append(p, r)
						} else {
							// It isn't the last phoneme so don't swallow the lexical r here
							// As in sOIree,...
							s = s[:len(s) -1]
						}
					}
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 3, []phoneme{w, ay, ax}, []phoneme{w, ay, axr}, []phoneme{w, ay, er}); ok {
					// This is probably choir, check to see if there's an r to swallow
					if _, ok := phsAt(phons, currPh, 4, []phoneme{w, ay, ax, r}); ok {
						p = append(p, r)
					}
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oi"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{w, aa}); ok {
					// As in cOIffure,...
					new = phonToAlphas{
						[]phoneme{
							w, aa,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 1, "o") {
				// As in vOyeur, One,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{w, aa}, []phoneme{w, ah}); ok {
					new = phonToAlphas{
						p,
						"o",
					}
				}
				break
			}
		case y:
			if stringAt(alphas, currAbc, 3, "aeo") {
				// As in palAEOntology,...
				new = phonToAlphas{
					[]phoneme{
						y,
					},
					"aeo",
				}
				break
			}
			// An early check for the word ewe.
			if alphas == "ewe" {
				// Note we can't check for the substring ewe else we'll break handling of words
				// like newest
				new = phonToAlphas{
					[]phoneme{
						y, uw,
					},
					"ewe",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "ean") {
				// As in cetacEAN,...
				new = phonToAlphas{
					[]phoneme{
						y,
					},
					"e",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ule"); ok {
				if p, ok := phsAt(phons, currPh, 3, []phoneme{y, uwl, ih}); ok {
					// The e is probably sounded so don't swallow it here
					// As in amULet,...
					new = phonToAlphas{
						p[:len(p) - 1],
						s[:len(s) - 1],
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 6, "you'll", "oule", "yule", "uel", "ule", "eul", "ul"); ok {
				// As in YOU'LL, YULE, valUELess, capsULE, EULogy, reticULar...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, uwl}); ok {
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				// As in the training pronunciation of inarticULate,...
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, axl}); ok {
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "eum", "um"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, uwm}); ok {
					// As in pnEUMonia, hUMan,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "eun", "ewn", "ugn", "un"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, uwn}); ok {
					// As in EUNuch, hEWN, impUGN, mUNicipal,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			// The many sounds of 'ure'...
			if s, ok := getStringAt(alphas, currAbc, 3, "eur", "ure", "ur"); ok {
				// As in EURope, tenUREd, cURious, 
				// Catch a possible sounded phonetic 'r'
				if _, ok := phsAt(phons, currPh, 3, []phoneme{y, ax, r}, []phoneme{y, uh, r}); ok {
					// As in failURE,...
					new = phonToAlphas{
						[]phoneme{
							y, ax,
						},
						"u",
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 4, []phoneme{y, uh, ax, r}, []phoneme{y, uw, ax, r}, []phoneme{y, uh, ah, r}, []phoneme{y, uw, ah, r}); ok {
					new = phonToAlphas{
						// Leave the r to be processed separately
						p[:len(p)-1],
						"u",
					}
					break
				}
				// Okay, we've covered the cases with a phonetic 'r'
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, ax}, []phoneme{y, axr}, []phoneme{y, er}); ok {
					// As in failURE, failURE,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 3, []phoneme{y, uh, axm}, []phoneme{y, uw, axm}); ok {
					// As in procUREment,...
					new = phonToAlphas{
						p[:len(p)-1],
						s,
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 3, []phoneme{y, uh, ax}, []phoneme{y, uh, axr}, []phoneme{y, uh, er}, []phoneme{y, uw, ax}, []phoneme{y, uw, er}); ok {
					// As in cURE, liquEUR, pURE,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ut"); ok {
				if _, ok:= phsAt(phons, currPh, 3, []phoneme{y, uw, sh}, []phoneme{y, uw, ch}, []phoneme{y, uw, tr}, []phoneme{y, uw, t}); !ok {
					// Not as in restitUtion, fUture, nUtrition, tUtor,...
					if _, ok := phsAt(phons, currPh, 2, []phoneme{y, uw}); ok {
						// As in debUT,...
						new = phonToAlphas{
							[]phoneme{
								y, uw,
							},
							s,
						}
						break
					}
				}			
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "eau", "eu"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, uw}); ok {
					// As in bEAUty, tEUtonic,...
					new = phonToAlphas{
						p,
						s,
					}
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "eo"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{y, ax}); ok {
					// The 'o' appears to be represented phonetically so just take
					// the 'e' now and leave the 'o' for future processing
					// As in metEorological,...
					new = phonToAlphas{
						[]phoneme{
							y,
						},
						"e",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ui"); ok {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{y, uw, ing}); ok {
					// The i is sounded so leave it for later
					// As in continUING,...
					new = phonToAlphas{
						[]phoneme{
							y, uw,
						},
						"i",
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{y, uw}); ok {
					if _, ok := phsAt(phons, currPh, 3, []phoneme{y, uw, ih}); !ok {
						// The lexical 'i' isn't sounded separately so swallow it all
						new = phonToAlphas{
							[]phoneme{
								y, uw,
							},
							s,
						}
					} else {
						// As in acUity
						new = phonToAlphas{
							[]phoneme{
								y, uw,
							},
							"u",
						}
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ew"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, uw}, []phoneme{y, uh}); ok {
					// As in stEWard,...
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{y, uwl}); ok {
					// As in nEWLy,...
					new = phonToAlphas{
						// Leave the l for uwL
						[]phoneme{
							y,
						},
						s,
					}
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "i", "y"); ok {
				// As in millIon, onIon, You,...
				new = phonToAlphas{
					[]phoneme{
						y,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "j") {
				// As in the German Junker,...
				new = phonToAlphas{
					[]phoneme{
						y,
					},
					"j",
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "u") {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{y, ax}, []phoneme{y, ah}, []phoneme{y, uh}, []phoneme{y, uw}); ok {
					// As in inarticUlate*, articUlate, tUreen, fUture,... At least according to the CMU dictionary
					// *This is found in the training pronunciation of inarticULate,...
					new = phonToAlphas{
						p,
						"u",
					}
				}
				break
			}
		case z:
			if _, ok := getStringAt(alphas, currAbc, 0, "ds"); ok {
				// As in bonDS,...
				new = phonToAlphas{
					[]phoneme{
						z,
					},
					"s",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "spb") {
				// As in raSPBerry,...
				// Check for a silent 'p'
				if _, ok := phsAt(phons, currPh, 2, []phoneme{z, b}); ok {
					// Looks like the 'p' is silent
					new = phonToAlphas{
						[]phoneme{
							z,
						},
						"sp",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "cz"); ok {
				// As in CZar,...
				new = phonToAlphas{
					[]phoneme{
						z,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "se", "ze"); ok {
				_, ok := phsAt(phons, currPh, 2, []phoneme{z, d})
				if ok || currPh == len(phons)-1 {
					// As in raiSE, raiSEd, raZE, raZEd,...
					new = phonToAlphas{
						[]phoneme{
							z,
						},
						s,
					}
					break
				}
			}
			if str, ok := getStringAt(alphas, currAbc, 2, "ss", "zz"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{z, s}); ok {
					// As in newSStand,...
					new = phonToAlphas{
						[]phoneme{
							z,
						},
						// Leave the second lexical 's for the phonetic 's'
						"s",
					}
					break
				}
				// As in sciSSors, fiZZed,...
				new = phonToAlphas{
					[]phoneme{
						z,
					},
					str,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "s", "x", "z"); ok {
				// As in many (but not all) plural word, for instance dogS, Xylem, zoo,...
				new = phonToAlphas{
					[]phoneme{
						z,
					},
					s,
				}
				break
			}
		case zh:
			if s, ok := getStringAt(alphas, currAbc, 4, "sion", "tion"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{zh, n}); ok {
					// As in fuSION, equaTION,...
					new = phonToAlphas{
						[]phoneme{
							zh, n,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "j", "s", "t", "z", "g"); ok {
				// As in beiJing, uSual, equaTion, aZure, beiGe,...
				new = phonToAlphas{
					[]phoneme{
						zh,
					},
					s,
				}
				break
			}
		case ax:
			if stringAt(alphas, currAbc, 2, "ve") {
				// As in should'VE,...
				new = phonToAlphas{
					[]phoneme{
						ax, v,
					},
					"ve",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "wer") {
				// As in ansWER,...
				s := "wer"
				// But only if the 'r' isn't sounded
				// As in ansWErable,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); ok {
					s = "we"
				}
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "wai") {
				// As in coxsWAIn,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					"wai",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "anc") {
				if _, ok := phsAt(phons, currPh, 3, []phoneme{ax, n, k}, []phoneme{ax, n, s}); !ok {
					// As in blANCmange,..
					// But not as in melANCholy, vagrANCy, and many other words...
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						"anc",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "har", "her"); ok {
				// As in philHARmonic, shepHERd,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "ha") {
				// As in FulHAm and many other place names
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					"ha",
				}
				break
			}
			if stringAt(alphas, currAbc, 3, "lel") {
				// As in candLELight,...
				if _, ok := phsAt(phons, currPh, 3, []phoneme{ax, l, l}); !ok {
					// The second lexical l isn't represented phonetically so swallow it here
					new = phonToAlphas{
						[]phoneme{
							ax, l,
						},
						"lel",
					}
					break
				}
			}
			if p, ok := phsAt(phons, currPh, 2, []phoneme{ax, l}, []phoneme{ax, m}, []phoneme{ax, n}); ok {
				// Trying to spot an inserted schwa as in bottLe, theisM, wasN't,...
				if stringAt(alphas, currAbc, 2, "le") {
					// Silent e?
					new = phonToAlphas{
						p,
						"le",
					}
					break
				}
				if s, ok := getStringAt(alphas, currAbc, 1, "l", "m", "n"); ok {
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "eau", "eou", "iou"); ok {
				// As in burEAUcrat, outragEOUs, suspicIOUs,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "eur", "oar", "our", "ure"); ok {
				// As in amatEUR, cupbOARd, flavOUR, futURE...
				// But not if the r is sounded, for instance in armOURy, usUREr...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "ur") {
				// As in aubURn,...
				// If the r is sounded, we should only grab the 'u'
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); ok {
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						"u",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 4, "ough") {
				// As in thorOUGH,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					"ough",
				}
				break
			}
			// Handling this separately as this is a mix of a silent e followed by a
			// vowel so I may handle this more generally at some point
			if stringAt(alphas, currAbc, 2, "ea") {
				// As in likEAble,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					"ea",
				}
				break
			}
			if stringAt(alphas, currAbc, 2, "yr") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); ok {
					// The lexical 'r' is sounded so leave it for later
					// As in labYRinth,...
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						"y",
					}
				} else {
					// As in zephYR,...
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						"yr",
					}
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ei", "ia", "ie", "io", "iu", "oi", "ou", "ua", "ui"); ok {
				// As in forEIgn, russIA, conscIEnce, percussIOn, tortOIse, belgIUm, luxuriOUs, usUAlly, biscUIt,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
			// Treat 're' separately.
			if stringAt(alphas, currAbc, 2, "re") {
				// As in theatRE,...
				p := []phoneme{ax}
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); ok {
					// But as in acREage...
					p = []phoneme{ax, r}
				}
				new = phonToAlphas{
					p,
					"re",
				}
				break
			}
			if stringAt(alphas, currAbc, 4, "erwr") {
				// As in undERWRite,...
				// Don't swallow the r phoneme here! It belong with the
				// lexical wr
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					"er",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ar", "er", "ir", "or", "ur"); ok {
				// As in wizARd, pERcussion weIRd, tractOR, pURveyor,...
				// Need to be careful here. I don't want to consume the r if there's
				// an r phoneme in the phonetic spelling, for instance as in
				// d(d)o(ao)c(k)u(y uh)m(m)e(eh)n(n)t(t)a(ax)r(r)y(iy)
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						s,
					}
					break
				}
			}
			// Some words have an -ah ending with the 'h' silent
			if stringAt(alphas, currAbc, 2, "ah") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ax, hh}); !ok {
					// The 'h' isn't sounded
					// As in purdAH,...
					new = phonToAlphas{
						[]phoneme{
							ax,
						},
						"ah",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e", "i", "o", "u", "y"); ok {
				// As in - pretty much anything you can think of and in particular, as in pYjamas,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "rr", "r"); ok {
				if p, ok := phsAt(phons, currPh, 2, []phoneme{ax, r}); ok {
					new = phonToAlphas{
						p,
						s,
					}
					break
				}
				// As in houR,...
				new = phonToAlphas{
					[]phoneme{
						ax,
					},
					s,
				}
				break
			}
		case axr:
			if s, ok := getStringAt(alphas, currAbc, 0, "ough", "eur", "our", "ure", "wer", "ar", "er", "ia", "ir", "or", "re", "ur", "yr", "a", "e", "o", "r"); ok {
				// As in thorOUGH, amatEUR, favOUR, sutURE, ansWER, liAR, harriER, nostalgIA, fIR, tailOR, metRE, lemUR, zephYR, troikA, timbrE, ontO*, cobbleR,...
				// *onto is a bit suspect for a standalone word pronunciation. In the case of cobbler the test for silent e might
				// swallow the lexical e so that all we're left with is the lexical r
				new = phonToAlphas{
					[]phoneme{
						axr,
					},
					s,
				}
				break
			}
		case oh:
			if s, ok := getStringAt(alphas, currAbc, 3, "eau", "au", "ho", "oh", "ou"); ok {
				// As in burEAUcracy, sAUsage, HOnest, jOHn, cOUgh,...
				new = phonToAlphas{
					[]phoneme{
						oh,
					},
					s,
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 1, "ow", "a", "e", "o"); ok {
				// As in knOWledge, whAt, Encore, cOttage,...
				new = phonToAlphas{
					[]phoneme{
						oh,
					},
					s,
				}
				break
			}
		case ehr:
			if stringAt(alphas, currAbc, 4, "heir") {
				// As in HEIR, but
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ehr, r}); !ok {
					// As in HEIR,...
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						"heir",
					}
				} else {
					// Don't swallow the r up, as in HEIress,...
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						"hei",
					}
				}
			}
			if strings.HasPrefix(alphas, "there") || strings.HasPrefix(alphas, "where") {
				// I'm being quite specific here. I don't want to break words like etHEREal.
				// As in words like THEREafter, WHEREof...
				if p, ok := phsAt(phons, currPh, 3, []phoneme{ehr, r, eh}); ok {
					// We don't want to swallow the second lexical e in words like whERever,...
					new = phonToAlphas{
						p[:len(p)-1],
						"er",
					}
					break
				}
				if p, ok := phsAt(phons, currPh, 2, []phoneme{ehr, r}); ok {
					new = phonToAlphas{
						p,
						"ere",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 4, "ayor"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ehr, r}); ok {
					// As in mAYOral,...
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						s[:len(s)-1],
					}
				} else {
					// As in mAYOR,...
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						s,
					}
				}
				break
			}
			if _, ok := getStringAt(alphas, currAbc, 0, "ar"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ehr,r }); ok {
					// The r is sounded so don't grab it here
					// As in phAraoh,...
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						"r",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "aire", "ayer", "air", "are", "ear", "eir", "ere", "ar", "er"); ok {
				// As in millionAIRE, prAYER, eclAIR, bewARE, wEAR, thEIR, whERE, scARce, concERto,...
				// But only if the r isn't sounded
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ehr, r}); !ok {
					new = phonToAlphas{
						[]phoneme{
							ehr,
						},
						s,
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ae", "ai", "ea"); ok {
				// As in AErial, AIring, wEAr (for alternative W EHR R phonetic spelling)...
				new = phonToAlphas{
					[]phoneme{
						ehr,
					},
					s,
				}
				break
			}

			if s, ok := getStringAt(alphas, currAbc, 1, "a", "e"); ok {
				// As in aquArium, whEreas,...
				new = phonToAlphas{
					[]phoneme{
						ehr,
					},
					s,
				}
				break
			}
		case axl:
			if stringAt(alphas, currAbc, 3, "lel") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{axl, l}); !ok {
					// As in candLELight,...
					new = phonToAlphas{
						[]phoneme{
							axl,
						},
						"lel",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "le") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{axl, eh}); ok {
					// The lexical 'e' is sounded, so leave it for later
					// As in coaLEsce,...
					new = phonToAlphas{
						[]phoneme{
							axl,
						},
						"l",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "rl") {
				// As in douRLy,... and many others in which the r isn't sounded
				new = phonToAlphas{
					[]phoneme{
						axl,
					},
					"rl",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "ou'll", "oughl", "iall", "ourl", "uall", "urel", "wale", "ael", "all", "arl", "aul", "ell", "erl", "ial", "ill", "oll", "orl", "rel", "ull", "al", "el", "il", "le", "ol", "ul", "yl"); ok {
				// As in yOU'LL, thorOUGHLy, partIALLy, odOURLess, usUALLy, leisURELy, gunWALE, michAEL, nationALLy, regulARLy, epAULettes, modELLing, eastERLy, specIAL, councILLor, pOLLute, fORLorn, seveRELy, awfULLy, typicAL, caramEL, civIL, articLE, viOLate, awfUL, sibYL,...
				new = phonToAlphas{
					[]phoneme{
						axl,
					},
					s,
				}
				break
			}
			// In some words axL is capturing a transition from a preceding vowel to l.
			// The vowel is still represented phonetically and has been mapped to
			// the (lexical) vowel so that all that remains is to map the phonetic
			// axL to the lexical l.
			if stringAt(alphas, currAbc, 1, "l") {
				// As in
				new = phonToAlphas{
					[]phoneme{
						axl,
					},
					"l",
				}
				break
			}
		case axm:
			// Catch a possible double phonetic m first
			if stringAt(alphas, currAbc, 3, "omm") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{axm, m}); ok {
					// As in bottOMmost,...
					new = phonToAlphas{
						[]phoneme{
							axm,
						},
						// Leave the second m until later
						"om",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 3, "emn", "umn"); ok {
				// As in solEMNly, colUMN,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{axm, n}); !ok {
					// But not as in solEMNity, colUMNist,...
					new = phonToAlphas{
						[]phoneme{
							axm,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 3, "ham") {
				// As in many place names like wrexHAM, ...
				new = phonToAlphas{
					[]phoneme{
						axm,
					},
					"ham",
				}
				break
			}
			if s, ok := getStringAt(alphas, currAbc, 2, "ancm", "urem", "amm", "arm", "erm", "iam", "irm", "ium", "olm", "omm", "orm", "umm", "urm", "am", "em", "im", "om", "um"); ok {
				// As in bLANCMange, measUREMent, grAMMatical, philhARMonic, vERMillion, parlIAMent, affIRMation, belgIUM, malcOLM, cOMMercial, infORMation, consUMMation, sURMise, fAMiliar, acadEMy, anIMal (not sure this should
				// be a schwa though), incOMe, vacuUM...
				new = phonToAlphas{
					[]phoneme{
						axm,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "m") {
				// As in prisM,...
				new = phonToAlphas{
					[]phoneme{
						axm,
					},
					"m",
				}
				break
			}
		case axn:
			if s, ok := getStringAt(alphas, currAbc, 0, "eignn", "enn", "ann", "onn", "ornn"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{axn, n}); ok {
					// As in forEIGNness, unevENness, humaNness, commONness, stubbORNness,...
					// The second lexical 'n' is sounded so dont swallow it here
					new = phonToAlphas{
						[]phoneme{
							axn,
						},
						s[:len(s) - 1],
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "oughn", "eignn", "eign", "erin", "ionn", "ornn", "ourn", "wain", "ain", "ann", "ean", "enn", "ien", "eon", "ern", "ian", "ign", "ion", "oen", "oln", "omp", "onn",
				"orn", "ren", "urn", "an", "en", "in", "on", "un", "n"); ok {
				// As in thorOUGHNess, forEIGNNess, forEIGN, vetERINary, legIONNaire, stubbORNNess, sojOURN, coxsWAIN, certAIN, ANNexe, pagEANt, unevENNess, conscIENce, burgEON, hibERNate, clinicIAN, ensIGN, equatION, rOENtgen, lincOLN, cOMPtroller, cONNection,
				// holbORN, meagRENess, aubURN, laymAN, conferENce, origINal, sextON, volUNteer, shouldN't,...
				new = phonToAlphas{
					[]phoneme{
						axn,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "n") {
				// This is a syllabic consonant, as in ?
				new = phonToAlphas{
					[]phoneme{
						axn,
					},
					"n",
				}
			}
		case ks:
			if s, ok := getStringAt(alphas, currAbc, 2, "xe"); ok {
				// Catch things like aXE, but not eXercise...
				if nextPhIsConsonant(phons, currPh) {
					new = phonToAlphas{
						[]phoneme{
							ks,
						},
						s,
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "xc") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ks, ch}, []phoneme{ks, k}, []phoneme{ks, kl}, []phoneme{ks, kr}); !ok {
					// As in eXCHange, eXCise but not as in eXCoriate, eXCLude, eXCRete, ...
					new = phonToAlphas{
						[]phoneme{
							ks,
						},
						"xc",
					}
					break
				}
			}
			if stringAt(alphas, currAbc, 2, "xs") {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ks, s}); !ok {
					// The phonetic 's' isn't sounded so grab it now
					// As in coXSwain,...
					new = phonToAlphas{
						[]phoneme{
							ks,
						},
						"xs",
					}
					break
				}
			}
			// Trying to trap the silent h that can follow ex-
			if stringAt(alphas, currAbc, 2, "xh") {
				// As in eXHibition,...
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ks, hh}); !ok {
					// Okay, the h is silent so include it now
					new = phonToAlphas{
						[]phoneme{
							ks,
						},
						"xh",
					}
					break
				}
			}
			//The 't' following x is sometimes not pronounced so catch it here
			if _, ok := getStringAt(alphas, currAbc, 0, "xtb"); ok {
				if _, ok := phsAt(phons, currPh, 2, []phoneme{ks, b}); ok {
					// The 't' is silent so swallow it now
					// As in teXTbook,... (I think this and its plural are the only examples)
					new = phonToAlphas{
						[]phoneme{
							ks,
						},
						"xt",
					}
					break
				}
			}
			if s, ok := getStringAt(alphas, currAbc, 0, "chs", "cks", "kes", "cc", "cs", "cz", "ks"); ok {
				// As in daCHShund, triCKSter, spoKESperson, aCCede, froliCSome, eczema, pranKSter,...
				new = phonToAlphas{
					[]phoneme{
						ks,
					},
					s,
				}
				break
			}
			// As in eXit,...
			new = phonToAlphas{
				[]phoneme{
					ks,
				},
				"x",
			}
		case kw:
			if s, ok := getStringAt(alphas, currAbc, 2, "cqu", "qu"); ok {
				// As in aCQUaint, QUiet,...
				new = phonToAlphas{
					[]phoneme{
						kw,
					},
					s,
				}
				break
			}
		case dz:
			if s, ok := getStringAt(alphas, currAbc, 2, "des", "d's", "ds"); ok {
				// As in blonDES, world's, trenDS,...
				new = phonToAlphas{
					[]phoneme{
						dz,
					},
					s,
				}
				break
			}
		case uwl:
			if s, ok := getStringAt(alphas, currAbc, 3, "o'll", "oel", "ool", "oul", "uel", "ule", "ul"); ok {
				// As in whO'LL, shOELace, schOOL, ampOULe, clUELess, rULE, unrULy...
				new = phonToAlphas{
					[]phoneme{
						uwl,
					},
					s,
				}
				break
			}
			if stringAt(alphas, currAbc, 1, "l") {
				new = phonToAlphas{
					[]phoneme{
						uwl,
					},
					"l",
				}
				break
			}
		default:
			break
		}
		if new.equal(phonToAlphas{}) {
			return fail(phon, alphas)
		}
		if isSilentE(alphas, currAbc, phons, currPh, new) {
			new.alphas += "e"
		}
		if punctuationSkipped != "" {
			new.alphas = punctuationSkipped + new.alphas
			currAbc -= len(punctuationSkipped)
			punctuationSkipped = ""
			// currAbc -= len(new.alphas)
		}
		// Catch any trailing characters not yet mapped to phonemes.
		// TODO: This is a bit crude but will do for now.
		new.alphas += trailingSilentAlphas(alphas, currAbc, phons, currPh, new)
		currPh += len(new.phons)
		currAbc += len(new.alphas)
		ret = append(ret, new)
	}
	if currAbc != len(alphas) || currPh != len(phons) {
		// This is also a failure. Since we loop on phonemes it's most likely
		// that currAbc != len(alphas)
		return fail(phons[len(phons) - 1], alphas)
	}
	return ret, nil
}
