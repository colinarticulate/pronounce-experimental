package pron

import (
	"sort"
)

// So if qr is a quantumResult then qr[aa] = 4 says that for a given timeslot (quantum
// of time) phoneme aa 'covers'the time slot 4 times
type quantumResult map[phoneme][]psPhonemeDatumRef

// quantisedResult[3][aa] = 4 means the phoneme aa covers time slot 3, 4 times
type timeAligner struct {
	results         []psPhonemeResults
	phons           []phoneme
	quantisedResult []quantumResult
}

func newTimeAligner(phons []phoneme) timeAligner {
	return timeAligner{
		[]psPhonemeResults{},
		phons,
		[]quantumResult{},
	}
}

// There must be a better way to do this but for now... This takes a result
// and creates a timeslot entry for each phoneme for each timeslot covered by
// the phoneme. Each entry contains an array of references to the phoneme
// instances that cover the timeslot
func (ta *timeAligner) AddResult(result psPhonemeResults) {
	// Add the result to the set of results
	ta.results = append(ta.results, result)
	i := len(ta.results) - 1
	for j, phonInst := range result.data {
		if phonInst.end > len(ta.quantisedResult)-1 {
			// Grow the number of timeslots if we need to
			for k := len(ta.quantisedResult); k <= phonInst.end; k++ {
				ta.quantisedResult = append(ta.quantisedResult, make(quantumResult))
			}
		}
		// Add entries for each timeslot covered by this phoneme
		for k := phonInst.start; k <= phonInst.end; k++ {
			if phonInst.end-phonInst.start == 0 {
				// It's not really there so there's nothing to add for this phoneme.
				// Just continue on to the next timeslot
				continue
			}
			ref := psPhonemeDatumRef{
				i, j,
			}
			if _, ok := ta.quantisedResult[k][phonInst.phoneme]; ok {
				ta.quantisedResult[k][phonInst.phoneme] = append(ta.quantisedResult[k][phonInst.phoneme], ref)
			} else {
				ta.quantisedResult[k][phonInst.phoneme] = []psPhonemeDatumRef{
					ref,
				}
			}
		}
	}
}

func union(refs1, refs2 []psPhonemeDatumRef) []psPhonemeDatumRef {
	contains := func(x []psPhonemeDatumRef, y psPhonemeDatumRef) bool {
		for _, ref := range x {
			if ref == y {
				return true
			}
		}
		return false
	}
	refs := refs1
	for _, ref := range refs2 {
		if !contains(refs, ref) {
			refs = append(refs, ref)
		}
	}
	return refs
}

func (ta timeAligner) filterRef(refs []psPhonemeDatumRef) []psPhonemeDatumRef {
	filtered := []psPhonemeDatumRef{}
	for _, ref := range refs {
		datum := ta.results[ref.scan].data[ref.index]
		if datum.end-datum.start > 0 {
			filtered = append(filtered, ref)
		}
	}
	return filtered
}

func (ta timeAligner) timeAlign() []timeAlignedPhoneme {
	// Report phonemes with a count > 2
	qRes := []timeAlignedPhoneme{}
	if len(ta.quantisedResult) == 0 {
		return []timeAlignedPhoneme{}
	}
	currentRes := make(map[phoneme]timeAlignedPhoneme)
	for phon := range ta.quantisedResult[0] {
		// Only interested in phonemes that 'cover' the timeslot 3 or more times.
		if len(ta.quantisedResult[0][phon]) > 2 {
			currentRes[phon] = timeAlignedPhoneme{
				psPhonemeDatum{
					phon,
					0, 0,
				},
				ta.filterRef(ta.quantisedResult[0][phon]),
			}
		}
	}
	// We're going to loop over the timeslots checking what phonemes cover the timeslot
	// to arrive at a time-aligned view of phonemes
	for i, result := range ta.quantisedResult {
		// If there are any phonemes in the current result that aren't in this timeslot then
		// append them to the time-aligned result and delete them from the current result
		for phon := range currentRes {
			if _, ok := result[phon]; !ok {
				qRes = append(qRes, currentRes[phon])
				delete(currentRes, phon)
			}
		}
		// Now loop over the phonemes that cover this timeslot
		for phon := range result {
			if _, ok := currentRes[phon]; ok {
				if len(result[phon]) > 2 {
					// Update the end time for this phoneme
					temp := currentRes[phon]
					temp.end = i
					// Also need to add in any additonal refs not already in temp.links
					temp.links = ta.filterRef(union(temp.links, result[phon]))
					currentRes[phon] = temp
				} else {
					// This phoneme appears less than 3 times in this timeslot so it's time to add
					// it to the time-aligned result and remove it from the current result...
					qRes = append(qRes, currentRes[phon])
					delete(currentRes, phon)
				}
			} else {
				// This is a new phoneme so add it to the current result.
				if len(result[phon]) > 2 {
					currentRes[phon] = timeAlignedPhoneme{
						psPhonemeDatum{
							phon,
							i, i,
						},
						ta.filterRef(result[phon]),
					}
				}
			}
		}
	}
	// Add whatever's left in the current result to the time-aligned result and sort by
	// start time
	for phon := range currentRes {
		qRes = append(qRes, currentRes[phon])
	}
	sort.Slice(qRes, func(i int, j int) bool {
		return qRes[i].start < qRes[j].start
	})
	sharedLinks := func(lPh1, lPh2 timeAlignedPhoneme) bool {
		if lPh1.phoneme != lPh2.phoneme {
			return false
		}
		for _, l1 := range lPh1.links {
			for _, l2 := range lPh2.links {
				if l1 == l2 {
					return true
				}
			}
		}
		return false
	}
	for i := 0; i < len(qRes)-1; i++ {
		for j := i + 1; j < len(qRes); j++ {
			// If there's a shared link and the gap between the time aligned
			// phonemes isn't 'too big' then join the phonemes together
			if sharedLinks(qRes[i], qRes[j]) && qRes[j].start-qRes[i].end < 3 {
				qRes[i].end = qRes[j].end
				qRes[i].links = union(qRes[i].links, qRes[j].links)
				qRes = append(qRes[:j], qRes[j+1:]...)
			}
		}
	}
	return modifiedResults(qRes)
}

func modifiedResults(results []timeAlignedPhoneme) []timeAlignedPhoneme {
	modRes := []timeAlignedPhoneme{}
	for i, res := range results {
		switch res.phoneme {

		case s:
			if i == 0 { // when S is the first letter
				if res.end-res.start > 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 {
				if res.end-res.start >= 4 { // in combiner s => 7, that is, each individual length but be >=7, here the combined lengths must be => 5
					modRes = append(modRes, res)
				}
			} else if i == len(results)-1 { // when s is the last letter
				if res.end-res.start > 5 {
					modRes = append(modRes, res)
				}
			}

		case t:
			if res.end-res.start >= 4 {
				modRes = append(modRes, res)
			} else if (i > 0) && (i < len(results)-1) && (results[i-1].phoneme == s) && (results[i+1].phoneme == s) {
				if res.end-res.start > 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == s && results[i+1].phoneme == r {
				if res.end-res.start > 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == s && results[i+1].phoneme == k {
				if res.end-res.start > 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == n && results[i+1].phoneme == r {
				if res.end-res.start > 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i <= len(results)-1 && results[i-1].phoneme == s { //  short t in voyeuristic ..... may need more work, seems too loose
				if res.end-res.start > 2 {
					modRes = append(modRes, res)
				}
			}

			/*
			   case th:
			     if res.end - res.start >= 4 {          //
			       modRes = append(modRes, res)
			     }

			   case dh:
			     if res.end - res.start >= 4 {          //
			       modRes = append(modRes, res)
			     }
			*/

		case l:
			if res.end-res.start >= 3 { // remove fake L in otter/walter
				modRes = append(modRes, res)
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == p { // too allow short L's after P
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == iy && results[i+1].phoneme == iy { // too allow short L's between iy & iy, as in swahili   ....... L's get swallowed up when between two vowels
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == ih && results[i+1].phoneme == ih { // too allow short L's between ih & ih, as in billings
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == uw && results[i+1].phoneme == ey { // too allow short L's between uw & ey, as in inoculates
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == ah && results[i+1].phoneme == ey { // too allow short L's between ah & ey, as in inoculates
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == uh && results[i+1].phoneme == ey { // too allow short L's between ah & ey, as in inoculates
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == ao && results[i+1].phoneme == ah { // too allow short L's between ao & ah, as PALAEONTOLOGY P EY L Y N T AO L AH JH IY
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == iy && results[i+1].phoneme == ey { // too allow short L's between ao & ah, as relations
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == iy && results[i+1].phoneme == eh { // too allow short L's between ao & ah, as relations (because ey is diphthong eh iy)
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if (i > 0 && i < len(results)-1) && results[i-1].phoneme == ah && results[i+1].phoneme == iy { // as in "aly"
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			}

		case w:
			if res.end-res.start >= 3 { // 5-- Really want this to be just for first position. (Remove fake w in water1_hossein)
				// 5 is too large ... killing of w in waving1_rishka  ... individual scans could be set to <= 4 in combiner
				modRes = append(modRes, res)
			} else {
				if i > 0 && i < len(results)-1 && results[i-1].phoneme == k {
					if res.end-res.start >= 3 {
						modRes = append(modRes, res)
					}
				}
			}

		/*   // removed due to the long hh in tressa how 3rd March 2021
		case hh:
		  if (res.end - res.start >= 4) && (res.end - res.start <= 22) {        // remove fake hh in clay/hay
		    modRes = append(modRes, res)
		    } else if (i > 0) && (i < len(results) - 1) && (results[i + 1].phoneme == uw) {   // short hh in who1_rishka
		      if res.end - res.start >= 3 {
		        modRes = append(modRes, res)
		      }
		  }
		*/

		/*
		   case b:
		     if (res.end - res.start > 1) && (res.end - res.start <= 22) {        // 18 required to remove fake b from because6T, but ebb has b=18
		       modRes = append(modRes, res)
		     } else if (i > 0) && (i < len(results) - 1) && (results[i - 1].phoneme == m) && (results[i + 1].phoneme == l) {   // as in crumble
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     }

		   case v:
		     if (res.end - res.start <= 25) {        // remove fake v from my1_aa_hossein/may
		       modRes = append(modRes, res)
		     }
		*/

		//if k is in the last position it must be greater than 4 -- kill off fake k's on the end of blat/black   ***** update when available *****
		// watch also sk in school
		case k:
			if res.end-res.start > 1 { // remove fake k in asked1_tomo   ..... kills some openning K's
				modRes = append(modRes, res)
			} else if (i > 0) && (i < len(results)-1) && (results[i-1].phoneme == s) && (results[i+1].phoneme == w) {
				if res.end-res.start >= 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i+1].phoneme == s { // This is X,   k into s,  k becomes suppressed
			} else if i == len(results) {
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			}

		case g:
			if (res.end-res.start > 1) && (res.end-res.start <= 30) { // greater than 1 to stop g g in grass1_paul
				modRes = append(modRes, res)
			}
		/*
		   case g:
		    if res.end - res.start >= 4 {          // remove the fake g in yet/get  ..... this really needs to be limited to the first position, 3 at other times
		      modRes = append(modRes, res)
		    }
		*/

		// Vowels
		// ======

		/*
		   case ih:
		     if res.end - res.start >= 3 {
		       modRes = append(modRes, res)
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == v && results[i + 1].phoneme == t {   // the short ih in tomo connectivity
		         if res.end - res.start >= 2 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == v && results[i + 1].phoneme == s {   // the short ih in conversation
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i + 1].phoneme == t {   // trying to catch "ity"
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i + 1].phoneme == ng {   // trying to catch "ing"
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i + 1].phoneme == ks {   // into X, like six
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && (results[i - 1].phoneme == er)  {
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && (results[i - 1].phoneme == ao)  {
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && (results[i - 1].phoneme == ah)  {
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == r && results[i + 1].phoneme == p {   // the short ih retried/replied
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == r && results[i + 1].phoneme == s {   // the short ih voyeuristic (maybe combine with above)
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == w && results[i + 1].phoneme == g {   // the short ih squiggle
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == w && results[i + 1].phoneme == r {   // the short ih squirrel
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == n && results[i + 1].phoneme == m {   // the short ih animal
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == d && results[i + 1].phoneme == s {   // the short ih discoverable
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == l && results[i + 1].phoneme == jh {   // the short ih knowledge
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     } else if i > 0 && i < len(results) - 1 && results[i - 1].phoneme == t && results[i + 1].phoneme == d {
		         if res.end - res.start >= 1 {
		           modRes = append(modRes, res)
		         }
		     }
		*/

		case iy:
			if res.end-res.start >= 4 {
				modRes = append(modRes, res)
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == r && results[i+1].phoneme == p { // the short ih retried/replied
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == b && results[i+1].phoneme == k { // the short iy in because
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == n && results[i+1].phoneme == m { // the short iy in alternate spelling for animal
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == r && results[i+1].phoneme == l { // the short iy in relations
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == sh && results[i+1].phoneme == l { // the short iy in shelia
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == aa { // short when in diphthong AY --> aa iy
				if res.end-res.start >= 1 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == ao { // short when in diphthong OY --> ao iy
				if res.end-res.start >= 1 {
					modRes = append(modRes, res)
				}
			}

			/*
			   case eh:
			     if res.end - res.start >= 3 {
			       modRes = append(modRes, res)
			     }
			*/

		case uw:
			if res.end-res.start >= 5 { // careful - short uw sound in eh-school1_hossein
				modRes = append(modRes, res)
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == k && results[i+1].phoneme == l { // the short uw in  eh-school1_hossein
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == y && results[i+1].phoneme == t { // the short uw in amputate
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == y && results[i+1].phoneme == f { // the short uw in manufacture
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			}

		// When in the last position >= 6 .... gets rid of fake uw in clay/clay         ** update when available ***********
		// However, when inside a word it can be shorter, such as in eh-school1_hossein

		case ao:
			if (res.end-res.start >= 4) && (res.end-res.start <= 49) { // The larger ao is to remove huge fake in animal2_deven
				// had to extend from 45 to 47 to cover long ao in caught1_rishka
				modRes = append(modRes, res)
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == w && results[i+1].phoneme == l { // the short iy in because
				if res.end-res.start >= 3 {
					modRes = append(modRes, res)
				}
			} else if i == len(results) {  // allowing longer "ao" sounds at the end of a word like ..your..
				if (res.end-res.start >= 4) && (res.end-res.start <= 60){
					modRes = append(modRes, res)
				}
			}

		case aa:
			if (res.end-res.start >= 3) && (res.end-res.start < 60) { // needed to extend 45 to 60 to cover the "aa" in because1_rishka
				modRes = append(modRes, res)
			}

		case ah:
			if res.end-res.start >= 3 {
				modRes = append(modRes, res)
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == y { // the short ah in amputate (vowel following vowel)
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == er { // the short ah in squirrel
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i-1].phoneme == r { // the short ah sometimes seen in russia (the first ah)
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			} else if i > 0 && i < len(results)-1 && results[i+1].phoneme == n { // the short ah at the start of under2_tressa
				if res.end-res.start >= 2 {
					modRes = append(modRes, res)
				}
			}

		default:
			modRes = append(modRes, res)
		}
	}
	return modRes
}
