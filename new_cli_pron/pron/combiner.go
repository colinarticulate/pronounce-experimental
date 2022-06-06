package pron

import (
	"fmt"
	"log"
	"sort"
)

type parsableResult struct {
	psPhonemeResults
	R_target
}

type parseMap struct {
	mapping [][]parseResult
}

type psPhonemeDatumRef struct {
	scan, index int
}

type linkedPhonVerdict struct {
	phonVerdict
	links []psPhonemeDatumRef
}

type timeAlignedPhoneme struct {
	psPhonemeDatum
	links []psPhonemeDatumRef
}

type combiner struct {
	results []parsableResult
	phons   []phoneme
	parseMap
	ruleAlignedVerdict []linkedPhonVerdict
	timeAligner
	timeAlignedPhonemes []timeAlignedPhoneme
}

func newCombiner(phons []phoneme) combiner {
	return combiner{
		[]parsableResult{},
		phons,
		parseMap{},
		[]linkedPhonVerdict{},
		newTimeAligner(phons),
		[]timeAlignedPhoneme{},
	}
}

func (c *combiner) addResult(result psPhonemeResults, tRule R_target) {
	// A quick check to make sure hte ruke we're about to ad has the same length
	// as any rule already added
	if len(c.results) != 0 && len(c.results[0].rules) != len(tRule.rules) {
		return
	}
	c.results = append(c.results, parsableResult{
		result,
		tRule,
	})
	mapEntry := make([]parseResult, len(tRule.rules))
	c.parseMap.mapping = append(c.parseMap.mapping, mapEntry)
	// c.timeAligner.AddResult(result)
}

func (c *combiner) parse() {
	for i, result := range c.results {
		// Parse each set of results against the rules
		parsRes, err := parse(result.data, 0, result.rules)
		if err == nil {
			c.mapping[i] = parsRes
		} else {
			c.print()
			_debug("parse: err =", err)
		}
	}
}

func parse(results []psPhonemeDatum, atPosn int, rules []parsableRule) ([]parseResult, error) {
	if len(rules) == 0 {
		// Base case
		return []parseResult{}, nil
	}
	parseResults, err := rules[0].parse(results, atPosn)
	if err != nil {
		return []parseResult{}, parseFailed
	}
	for _, result := range parseResults {
		subresults, err := parse(results, result.end+1, rules[1:])
		if err == nil {
			return append([]parseResult{result}, subresults...), nil
		}
		// else we try the next parserResult...
	}
	return []parseResult{}, parseFailed
}

func mapLinkedVerdicts(linkedVerdicts []linkedPhonVerdict) []phonVerdict {
	retVerdicts := []phonVerdict{}
	for _, verdict := range linkedVerdicts {
		retVerdicts = append(retVerdicts, verdict.phonVerdict)
	}
	return retVerdicts
}

func (c combiner) whereInTime(phons []psPhonemeDatumRef) []int {
	intersects := func(l1s, l2s []psPhonemeDatumRef) bool {
		for _, l1 := range l1s {
			for _, l2 := range l2s {
				if l1 == l2 {
					return true
				}
			}
		}
		return false
	}
	where := []int{}
	for i, ph := range c.timeAlignedPhonemes {
		if intersects(phons, ph.links) {
			where = append(where, i)
		}
	}
	return where
}

func (c combiner) whereInTimeAfter(phons []psPhonemeDatumRef, time int) []int {
	intersects := func(l1s, l2s []psPhonemeDatumRef, time int) bool {
		for _, l1 := range l1s {
			for _, l2 := range l2s {
				l2Start := c.results[l2.scan].data[l2.index].start
				if l1 == l2 && l2Start >= time {
					return true
				}
			}
		}
		return false
	}
	where := []int{}
	for i, ph := range c.timeAlignedPhonemes {
		if intersects(phons, ph.links, time) {
			where = append(where, i)
		}
	}
	return where
}

func intersection(datumRefs1, datumRefs2 []psPhonemeDatumRef) []psPhonemeDatumRef {
	intersect := []psPhonemeDatumRef{}
	for _, ref1 := range datumRefs1 {
		for _, ref2 := range datumRefs2 {
			if ref1 == ref2 {
				intersect = append(intersect, ref1)
			}
		}
	}
	return intersect
}

type timedPhonVerdict struct {
	phonVerdict
	startsAt int
}

func (c combiner) insert(into []timedPhonVerdict, ph timeAlignedPhoneme) []timedPhonVerdict {
	newEntry := timedPhonVerdict{
		phonVerdict{
			ph.phoneme,
			missing,
		},
		ph.start,
	}
	for i, verdict := range into {
		if ph.start >= verdict.startsAt {
			if len(ph.links) == 3 {
				newEntry.goodBadEtc = possible
			} else if len(ph.links) > 3 {
				newEntry.goodBadEtc = good
			}
			return append(append(into[:i], newEntry), into[i:]...)
		}
	}
	return append(into, newEntry)
}

/*
func (c *combiner) integrate() []phonVerdict {
  ruleResult := c.ruleAlignedVerdict
  timeResult := c.timeAlignedPhonemes
  // log.Println("ruleResult =", c.ruleAlignedVerdict)
  // log.Println("timeResult =", c.timeAlignedPhonemes)

  if len(timeResult) == 0 {
    // There's nothing to do so just return the verdicts
    return mapLinkedVerdicts(ruleResult)
  }

  // Now combine the two different views
  verdicts := []phonVerdict{}
  iTime := 0
  for i := 0; i < len(ruleResult); i++ {
    if iTime >= len(timeResult) {
      // We've run out of time-aligned results so process what's left of the
      // rule-aligned verdicts, (one at a time)
      if ruleResult[i].goodBadEtc != surprise {
        verdict := ruleResult[i].phonVerdict
        verdict.goodBadEtc = missing

        verdicts = append(verdicts, verdict)
      } else {
        // If it's a surprise then the time aligned view reckons it's not
        // really there so do nothing
      }
      continue
    }
    if timeResult[iTime].phoneme == sil {
      // Skip over sil's in the timeResult
      iTime++
      if iTime >= len(timeResult) {
        verdicts = append(verdicts, ruleResult[i].phonVerdict)
        continue
      }
    }
    result := ruleResult[i]
    switch result.goodBadEtc {
    case good, possible:
      // If it's not present in the time-aligned view then it's probably not there, with some exceptions:
      // E.g. b's,
      // aa then iy - the iy can be short in this case, particular as part of a diphthong,
      // n's in the middle of a word
      if result.phon == timeResult[iTime].phoneme {
        // It probably is really there. - Yep, but we need to check how much
        // it's really there.
        intersectingLinks := intersection(result.links, timeResult[iTime].links)
        var goodBadEtc verdict
        switch len(intersectingLinks) {
        case 0, 1, 2:
          goodBadEtc = missing
        case 3:
          goodBadEtc = possible
        default:
          goodBadEtc = good
        }
        newVerdict := result.phonVerdict
        newVerdict.goodBadEtc = goodBadEtc
        verdicts = append(verdicts, newVerdict)
        iTime++
      } else {
     where := c.whereInTime(result.links)
        if len(where) == 0 {
          newVerdict := result.phonVerdict
          newVerdict.goodBadEtc = missing
          verdicts = append(verdicts, newVerdict)
        } else {
          intersectingLinks := intersection(result.links, timeResult[where[0]].links)
          if len(intersectingLinks) > 2 {
            for j := iTime; j < where[0]; j++ {
              if len(timeResult[j].links) > 2 {
                newSurprise := phonVerdict{
                  timeResult[j].phoneme, surprise,
                }
                verdicts = append(verdicts, newSurprise)
              }
              iTime++
            }
          }
          var goodBadEtc verdict
          switch len(intersectingLinks) {
          case 0, 1, 2:
            goodBadEtc = missing
          case 3:
            goodBadEtc = possible
          default:
            goodBadEtc = good
          }
          newVerdict := phonVerdict{
            result.phon, goodBadEtc,
          }
          verdicts = append(verdicts, newVerdict)
          iTime++
        }
      }
    case missing:
      // If it's missing in the rule-aligned view under what scenario might it actually be present in the time-aligned view?
      // It might be missing because it's attached to the wrong rule. For now assume that we have at least one phoneme instance in this rule.
      if result.phon == timeResult[iTime].phoneme {
        // It looks like the 'missing' phoneme is there
        numLinks := len(timeResult[iTime].links)
        if numLinks >= 3 {
          goodBadEtc := good
          if numLinks == 3 {
            goodBadEtc = possible
          }
          newVerdict := result.phonVerdict
          newVerdict.goodBadEtc = goodBadEtc
          verdicts = append(verdicts, newVerdict)
          iTime++
        } else {
          verdicts = append(verdicts, result.phonVerdict)
        }
      } else {
        // Perhaps it really is missing. I really should search forward here
        // to see if I can find it.
        verdicts = append(verdicts, result.phonVerdict)
      }
    case surprise:
      // This might not be time-aligned so not really there, if the phonemes are small it's probably noise.
      // Need to filter out guards here unless it's large and time-aligned.
      if result.phon == timeResult[iTime].phoneme{
        // Maybe this really is a surprise
        verdicts = append(verdicts, result.phonVerdict)
        iTime++
      } else {
        // Maybe it isn't. Or maybe it is and there are other time-aligned
        // surprises as well
        where := c.whereInTime(result.links)
        if len(where) == 0 {
          // There's something odd going on in this case. If we have a surprise
          // but it doesn't appear as a time-aligned surprise where's it gone?
          // It's probably not really there...
        } else {
          for j := iTime; j < where[0]; j++ {
            if len(timeResult[j].links) > 2 {
              newSurprise := phonVerdict{
                timeResult[j].phoneme, surprise,
              }
              verdicts = append(verdicts, newSurprise)
            }
            iTime++
          }
        }
      }
    }
  }
  // If there's anything left in the time-aligned results then these will be
  // surprises so add them in
  if iTime < len(timeResult) {
    for i := iTime; i < len(timeResult); i++ {
      verdict := phonVerdict{
        timeResult[i].phoneme,
        surprise,
      }
      verdicts = append(verdicts, verdict)
    }
  }
  return verdicts
}
*/

func (c *combiner) integrate() []phonVerdict {
	ruleResult := c.ruleAlignedVerdict
	timeResult := c.timeAlignedPhonemes
	if len(timeResult) == 0 {
		// There's nothing to do so just return the verdicts
		return mapLinkedVerdicts(ruleResult)
	}

	// Now combine the two different views
	verdicts := []phonVerdict{}
	iTime := 0
	timeStartFrom := 0
	hwmITime := iTime
	for i := 0; i < len(ruleResult); i++ {
		result := ruleResult[i]
		where := c.whereInTimeAfter(result.links, timeStartFrom)
		newITime := -1
		for _, j := range where {
			// We're not interested in time aligned phonemes that are before where
			// we think we are in the time aligned phonemes
			if j >= iTime {
				// Assume this is where the time aligned phonemes we're interested in
				// are
				newITime = j
				break
			}
		}
		_debug("i, iTime, newITime, timeStartFrom =", i, iTime, newITime, timeStartFrom)
		if len(where) == 0 || newITime < 0 {
			// We've not found any links in result.links in the time aligned
			// phonemes. How can that happen?
			// The phonemes could be filtered out by the time aligner, for
			// instance if they're too short or too long
			if result.goodBadEtc == surprise {
				// We've not found any surprising phonemes so move onto the next
				// rule aligned phoneme
				continue
			}
			// If the result is anything else then we should probably record this
			// phoneme as missing
			missingVerdict := phonVerdict{
				result.phon, missing,
			}
			verdicts = append(verdicts, missingVerdict)
		} else {
			timeLinks := []psPhonemeDatumRef{}
			for _, link := range timeResult[newITime].links {
				if c.results[link.scan].data[link.index].start >= timeStartFrom {
					timeLinks = append(timeLinks, link)
				}
			}
			intersectingLinks := intersection(result.links, timeLinks)
			if newITime > iTime {
				// We've found some phonemes in the time-aligned result that don't
				// appear in the rule-aligned result so add these as surprises
				for j := iTime; j < newITime; j++ {
					linksAfterCount := 0
					for _, link := range timeResult[j].links {
						if c.results[link.scan].data[link.index].start >= timeStartFrom {
							linksAfterCount++
						}
					}
					if linksAfterCount > 2 {
						newVerdict := phonVerdict{
							timeResult[j].phoneme, surprise,
						}
						verdicts = append(verdicts, newVerdict)
					}
				}
			}
			var goodBadEtc verdict
			switch len(intersectingLinks) {
			case 0, 1, 2:
				if result.goodBadEtc == surprise {
					// It looks like it's not really there
					continue
				}
				goodBadEtc = missing
			default:
				// We go with the rule-aligned verdict - which could be surprise,
				// possible or good.
				goodBadEtc = result.goodBadEtc
			}
			newVerdict := result.phonVerdict
			newVerdict.goodBadEtc = goodBadEtc
			verdicts = append(verdicts, newVerdict)
			iTime = newITime
			if iTime > hwmITime {
				hwmITime = iTime
			}
			timeStartFrom = c.results[intersectingLinks[0].scan].data[intersectingLinks[0].index].end
			for _, link := range intersectingLinks {
				end := c.results[link.scan].data[link.index].end
				if timeStartFrom > end {
					timeStartFrom = end
				}
			}
		}
	}
	// If there's anything left in the time-aligned results then these will be
	// surprises so add them in
	if hwmITime < len(timeResult)-1 {
		for i := hwmITime + 1; i < len(timeResult); i++ {
			verdict := phonVerdict{
				timeResult[i].phoneme,
				surprise,
			}
			verdicts = append(verdicts, verdict)
		}
	}
	return verdicts
}

// There's an implict assumption here that all rules are basically the same
// except for the guards - basically the first rule is an R_trapppedOpening
// rule, followed by the same number of R_phoneme rules and ending with a
// R_closing rule
func (c *combiner) ruleAlign() []linkedPhonVerdict {
	if len(c.results) == 0 {
		return []linkedPhonVerdict{}
	}
	verdicts := []linkedPhonVerdict{}
	i := 0
	for i < len(c.results[0].rules) {
		r := c.results[0].rules[i]
		if _, ok := r.(R_trappedOpening); ok {
			verdicts = append(verdicts, c.unexpectedOpening(i)...)
		}
		if _, ok := r.(R_phoneme); ok {
			verdicts = append(verdicts, c.unexpectedPhoneme(i)...)
		}
		if _, ok := r.(R_diphthongPhoneme); ok {
			verdicts = append(verdicts, c.unexpectedPhoneme(i)...)
		}
		if _, ok := r.(R_closing); ok {
			verdicts = append(verdicts, c.unexpectedClosing(i)...)
		}
		i++
	}
	c.ruleAlignedVerdict = verdicts
	return verdicts
}

func (c *combiner) unexpectedOpening(ruleIndex int) []linkedPhonVerdict {
	if len(c.results) == 0 {
		return []linkedPhonVerdict{}
	}
	unexpecteds := make(map[phoneme][]psPhonemeDatumRef)
	if ruleIndex < 0 || ruleIndex > len(c.results[0].rules)-1 {
		return []linkedPhonVerdict{}
	}
	if ruleIndex != 0 {
		return []linkedPhonVerdict{}
	}
	r, ok := c.results[0].rules[ruleIndex].(R_trappedOpening)
	if !ok {
		// Something's wrong so crash the program for now
		c.print()
		_debug("unexpectedOpening: type assertion failed")
		log.Panic()
	}
	for i, result := range c.mapping {
		// We may have no data - typically if pocketsphinx crashes as it
		// occasionally does
		if len(c.results[i].data) == 0 {
			continue
		}
		for j := result[ruleIndex].start; j <= result[ruleIndex].end; j++ {
			datum := c.results[i].data[j]
			ph := datum.phoneme
			if ph == sil {
				// Not interested in sil so throw it away
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			if ph == r.rT.trap {
				// Throw away the trap
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			if ph == r.rO.guard {
				if datum.end-datum.start <= 7 {
					// Throw away the guard (unless it's long) in any result that is
					// returned
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}
			if datum.end-datum.start <= 3 {
				// It's a short unexpected phoneme so probably not really there
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			ref := psPhonemeDatumRef{
				i, j,
			}
			if _, ok := unexpecteds[ph]; ok {
				unexpecteds[ph] = append(unexpecteds[ph], ref)
			} else {
				unexpecteds[ph] = []psPhonemeDatumRef{
					ref,
				}
			}
		}
	}
	retVerdict := []linkedPhonVerdict{}
	for ph := range unexpecteds {
		if len(unexpecteds[ph]) > 2 {
			// Return a phonVerdict
			lv := linkedPhonVerdict{
				phonVerdict{
					ph,
					surprise,
				},
				unexpecteds[ph],
			}
			retVerdict = append(retVerdict, lv)
		}
	}
	return retVerdict
}

func (c combiner) timeOrderAlignedPhons(alignedPhs map[phoneme][]psPhonemeDatumRef) []phoneme {
	type avgStart struct {
		ph    phoneme
		start int
	}
	starts := []avgStart{}
	for ph, refs := range alignedPhs {
		// Keep a tally of the average start time
		start := 0
		for _, ref := range refs {
			start += c.results[ref.scan].data[ref.index].start
		}
		if start != 0 {
			start = int(start / len(refs))
		}
		starts = append(starts, avgStart{
			ph, start,
		})
	}
	// Now return the phonemes sorted by average start time
	sort.Slice(starts, func(i int, j int) bool {
		return starts[i].start < starts[j].start
	})
	ret := []phoneme{}
	for _, start := range starts {
		ret = append(ret, start.ph)
	}
	return ret
}

// The code below is tightly bound to the grammar that was used to generate
// pocketsphinx results.
// The currect grammar for a (expected) phoneme rule is:
//
// [g + v] [g + c] (p | (g + a))
//
// where g = guard, v = vowel, c = consonant, p = expected phoneme, a = any (so
// vowel or consonant)
//
// In each of the pocketsphinx results I'm aligning phonemes found into one of
// four slots (indexes). So, v in slot 0, c in slot 1, p in slot 2 and a in
// slot 3.
// So in alignedPhs below you'll always see the expected phoneme (if found)
// in slot 2.
// There's a switch statement for unexpected phonemes and where they should go.
// Placement of these phonemes is dependent on the number of pocketsphinx
// phonemes that are mapped to this rule (see the condition for the switch
// - result[ruleIndex].end - result[ruleIndex].start + 1 which gives the
// number of phonemes)
// The cases considered are as follows:
// 1: The only way we can have one phoneme mapped to this rule is if we've
// found the phoneme, p, expected by this rule, and since we only ever get to
// the switch if we haven't found the expected phoneme this should never happen
// 2: If there are two phonemes these must be the g + a phonemes, so the a
// phoneme goes in slot 3
// 3: We've found either g + v or g + c, and the expected phoneme and the v
// phoneme goes in slot 0 , the c phoneme in slot 1
// 4: This case is more complex. There are a couple of possibilities here.
// Either we've found g + v, g + a or g + c, g + a, hence the check to see
// whether the current phoneme is a vowel when
// j - result[ruleIndex].start < 2
// 5: This is pretty straightforward. If there are 5 phonemes then we must have
// found g + v, g + c, p
// 6: This is also pretty straighforward. If there are 6 phonemes then we've
// found g + v, g + c, g + a
//
func (c *combiner) standardPhoneme(ruleIndex int) []linkedPhonVerdict {
	r, ok := c.results[0].rules[ruleIndex].(R_phoneme)
	if !ok {
		// Something's wrong so crash the program for now
		c.print()
		_debug("standardPhoneme: type assertion failed")
		log.Panic()
	}
	alignedPhs := make([]map[phoneme][]psPhonemeDatumRef, 4)
	for i := range alignedPhs {
		alignedPhs[i] = make(map[phoneme][]psPhonemeDatumRef)
		// Initialise alignedPhs with the expected phoneme regardless of what's in
		// the results - (if it doesn't appear at all we won't get an opportunity
		// later on to add it in)
		if i == 2 {
			alignedPhs[i][r.phon] = []psPhonemeDatumRef{}
		}
	}
	for i, result := range c.mapping {
		// We may have no data - typically if pocketsphinx crashes as it
		// occasionally does
		if len(c.results[i].data) == 0 {
			continue
		}
		r, ok := c.results[i].rules[ruleIndex].(R_phoneme)
		if !ok {
			// Something's wrong so crash the program for now
			c.print()
			_debug("standardPhoneme: type assertion failed")
			log.Panic()
		}
		for j := result[ruleIndex].start; j <= result[ruleIndex].end; j++ {
			datum := c.results[i].data[j]
			ph := datum.phoneme

			if ph == r.guard {
				if datum.end-datum.start < 5 {
					// Throw away the guard (unless it's long) in any result that is
					// returned
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			// Remove long phonemes

			//factor := float64(c.results[i].frate) / 100.0

			/*
			   if ph == b && (ruleIndex == 1) && (c.results[i].data[j].end - c.results[i].data[j].start > 18)   {
			     // Specifically, trying to get rid of any long b phoneme which is the
			     // first phoneme in a word.  Increased to 18 because of the "b" in black said by US female - which is really "ber" but pocketsphinx seems trained that way
			     datum.end = datum.start
			     c.results[i].data[j] = datum
			     continue
			   }
			*/

			if ph == v && (ruleIndex != len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start >= 20) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			//remove long f's except when in the last position, eg tressa "if"  3rd March 2021
			if ph == f && (ruleIndex != len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start > 30) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			//   s
			//=======

			// (I think, maybe wrong, that this is checking that an S in the last position is always greater than 5)
			if ph == s && (ruleIndex == len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start < 6) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == s && (ruleIndex == 1) && (c.results[i].data[j].end-c.results[i].data[j].start < 6) {
				// removing false S from the front of yet/set
				// The problem with killing of short s in the openning position is that words like "six" may start to fail
				// Instead if the S is going into a short vowel, then allow a short S, for example six & sit
				phonIndex := ruleIndex - 1
				if (phonIndex == 0) && (c.phons[phonIndex+1] == ih) {
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			/*
			   if ph == s && (c.results[i].data[j].end - c.results[i].data[j].start < 5) {
			   // Specifically, trying to get rid of any short s phoneme... This is
			   // useful in killing off false s's, for example next to a t as in vests
			     datum.end = datum.start
			     c.results[i].data[j] = datum
			     continue
			   }
			*/

			//   g
			//=======

			if ph == g && (c.results[i].data[j].end-c.results[i].data[j].start > 30) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			/*
				      if ph == g && (ruleIndex == 1) && (c.results[i].data[j].end - c.results[i].data[j].start <= 4) {

				        // to remove the false g at the front of yet/get
				        datum.end = datum.start
				        c.results[i].data[j] = datum
				        continue
				      }

					 if ph == g && (ruleIndex == len(c.results[i].rules) - 2) && (c.results[i].data[j].end - c.results[i].data[j].start <= 4) {
				        // to remove false g on ebb/egg
				        datum.end = datum.start
				        c.results[i].data[j] = datum
				        continue
				      }
			*/

			//   t
			//=======

			/*
			   if ph == t && (ruleIndex == len(c.results[i].rules) - 2) && (c.results[i].data[j].end - c.results[i].data[j].start <= 5) {
			     // Specifically, trying to get rid short false t in tomo-asked
			     datum.end = datum.start
			     c.results[i].data[j] = datum
			     continue
			   }
			*/

			if ph == t && (c.results[i].data[j].end-c.results[i].data[j].start <= 3) {
				phonIndex := ruleIndex - 1
				if (phonIndex > 0) && (phonIndex < len(c.phons)-1) && (c.phons[phonIndex-1] == s) && (c.phons[phonIndex+1] == s) {
					//
				} else if (phonIndex == len(c.phons)-1) && (c.phons[phonIndex-1] == s) {
					//
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			//   k
			//=======
			//(Careful of order)

			// Remove long voiceless plosives .... except when they are the very last phoneme .... because they run into silence
			if (ph == k || ph == p) && (ruleIndex != len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start > 20) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == k && (ruleIndex == len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start <= 5) {
				// Specifically, trying to get rid short false K's at the end of black, blah, blair etc
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			/*
			   if ph == k && (ruleIndex == 1) && (c.results[i].data[j].end - c.results[i].data[j].start <= 4) {
			     // Specifically, trying to get rid short false K's at beginning of time/climb

			     //
			     // concept - set the time information to 0, so that this rule effectively carries into the time-aligner (requires rule-aligned to be done first, which it currently is)
			     // c.results[i].data[j].end = c.results[i].data[j].start
			     //
			     datum.end = datum.start
			     c.results[i].data[j] = datum
			     continue
			   }
			*/

			/*
					  if ph == k && (c.results[i].data[j].end - c.results[i].data[j].start < 3) {
				        phonIndex := ruleIndex - 1
				        if phonIndex > 0 && phonIndex < len(c.phons) - 1 && c.phons[phonIndex - 1] == s && c.phons[phonIndex + 1] == w {     // careful of the K in squirrel S K W
				        } else if phonIndex > 0 && phonIndex < len(c.phons) - 1 && c.phons[phonIndex - 1] == s {     // careful of the K in school1_colin
				        } else {
				        datum.end = datum.start
				        c.results[i].data[j] = datum
				          continue
				        }
				      }
			*/

			//   hh
			//=======

			/*   // removed due to the long hh in tressa how 3rd March 2021
			       if ph == hh && (c.results[i].data[j].end - c.results[i].data[j].start > 21)  {
			      // Getting rid of false h in clay/hay
			       datum.end = datum.start
			       c.results[i].data[j] = datum
			     continue
			     }
			*/

			if ph == hh && (ruleIndex == 1) && (c.results[i].data[j].end-c.results[i].data[j].start < 4) {
				// Getting rid of false h in clay/hay
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == p && (ruleIndex == len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start < 5) {
				// removing fake p's in sit/sip
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == n && (ruleIndex == 1) && (c.results[i].data[j].end-c.results[i].data[j].start <= 3) {
				phonIndex := ruleIndex - 1
				if (phonIndex > 0) && (phonIndex < len(c.phons)-1) && (c.phons[phonIndex+1] == v) { // allow "nv" of any size as in conversation
				} else if (phonIndex > 0) && (phonIndex < len(c.phons)-1) && (c.phons[phonIndex-1] == ae) && (c.phons[phonIndex+1] == ih) { // allow very short n in animal
				} else if (phonIndex > 0) && (phonIndex < len(c.phons)-1) && (c.phons[phonIndex-1] == ae) && (c.phons[phonIndex+1] == iy) { // allow very short n in animal
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			if ph == "r" && (ruleIndex == 1) && (c.results[i].data[j].end-c.results[i].data[j].start < 3) {
				// Specifically, trying to get rid short false R's at beginning of replied - tomo
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == w && (ruleIndex == 1) && (c.results[i].data[j].end-c.results[i].data[j].start <= 4) {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == l && (c.results[i].data[j].end-c.results[i].data[j].start < 3) { // L's can be small, even though they shouldn't be, eg school1_colin
				phonIndex := ruleIndex - 1
				if phonIndex > 0 && phonIndex < len(c.phons)-1 && c.phons[phonIndex-1] == p {
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			if ph == l && (c.results[i].data[j].end-c.results[i].data[j].start > 37) { // to get rid of the long noise L in vater/water .... not sure how universal this is
				// ... increased to 35 because of long L in lonely1_philip
				// ... increased to 37 because of long L in animal2_deven
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			/*
			   	  if ph == "th" && (ruleIndex == 1) &&  (c.results[i].data[j].end - c.results[i].data[j].start <= 4){
			           // Specifically, trying to get rid fake th in thunder/tressa
			           datum.end = datum.start
			           c.results[i].data[j] = datum
			           continue
			         }


			   	  if ph == "dh" && (ruleIndex == 1) &&  (c.results[i].data[j].end - c.results[i].data[j].start <= 3){
			            // Specifically, trying to get rid fake th in thunder/tressa
			           datum.end = datum.start
			           c.results[i].data[j] = datum
			           continue
			         }

			         if ph == "dh" && (c.results[i].data[j].end - c.results[i].data[j].start < 3){
			            // to get rid fake dh in otter/other
			           datum.end = datum.start
			           c.results[i].data[j] = datum
			           continue
			         }
			*/

			//  Vowels
			//==========

			/*
			   	// Needs more investigation because "ih" can be very small especially if part of a diphthong or sat next to another vowel.
			   	  if ph == ih && (c.results[i].data[j].end - c.results[i].data[j].start <= 3)  {

			           phonIndex := ruleIndex - 1

			           if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == b) && (c.phons[phonIndex + 1] == k) {
			              //  keep  ... ie falls down to lenPh
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex + 1] == ng) {   // trying to catch "ing"s
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex + 1] == ng) {   // trying to catch "ity"
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex + 1] == ks) {   // ih into x
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == w) && (c.phons[phonIndex + 1] == g) {     // squiggle
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == w) && (c.phons[phonIndex + 1] == "r") {     // squirrel
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == v) && (c.phons[phonIndex + 1] == s) {     // conversation
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == n) && (c.phons[phonIndex + 1] == m) {     // animal
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == d) && (c.phons[phonIndex + 1] == s) {     // discoverable
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == l) && (c.phons[phonIndex + 1] == jh) {     // knowledge
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && (c.phons[phonIndex - 1] == t) && (c.phons[phonIndex + 1] == d) {
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && c.phons[phonIndex - 1] == er {
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && c.phons[phonIndex - 1] == ao {
			              //
			           } else if phonIndex > 0 && phonIndex < (len(c.phons) - 1) && c.phons[phonIndex - 1] == ah {
			              //
			           } else {
			           datum.end = datum.start
			           c.results[i].data[j] = datum
			             continue
			           }
			           //  to get rid short false ih in because_6t
			         }   // ... careful of the short ih in voyeuristic ... V OY ER IH S T IH K   // or AO, AH
			*/

			if ph == iy && (c.results[i].data[j].end-c.results[i].data[j].start <= 4) {
				//  to get rid short false iy in because_6t

				phonIndex := ruleIndex - 1

				if phonIndex > 0 && phonIndex < (len(c.phons)-1) && (c.phons[phonIndex-1] == b) && (c.phons[phonIndex+1] == k) {
					//  keep  ... ie falls down to lenPh
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == "r" { //if a semi-vowel is before iy
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == l {
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == w {

				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex+1] == "r" { //if a semi-vowel is after iy
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex+1] == l {
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex+1] == w {

				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == er {
					//  keep  ... ie falls down to lenPh
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == ao {
					//
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == ah {
					//
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex-1] == aa { // for the diphthong AY = aa iy
					//
				} else if (phonIndex > 0) && (phonIndex < (len(c.phons) - 1)) && c.phons[phonIndex-1] == n && c.phons[phonIndex+1] == m { // for the for iy in alternate spelling of animal
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			/*
			   if ph == eh && (c.results[i].data[j].end - c.results[i].data[j].start <= 3){
			     datum.end = datum.start
			     c.results[i].data[j] = datum
			     continue
			   }
			*/

			if ph == uw && (ruleIndex == len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start < 6) {
				// careful - short uw sound in eh-school1_hossein
				//  This is set for checking last position only at the moment
				phonIndex := ruleIndex - 1

				if phonIndex > 0 && phonIndex < (len(c.phons)-1) && (c.phons[phonIndex-1] == y) && (c.phons[phonIndex+1] == t) { // amputate
					//  keep  ... ie falls down to lenPh
				} else if (phonIndex > 0) && (phonIndex < (len(c.phons) - 1)) && c.phons[phonIndex-1] == y && c.phons[phonIndex+1] == f { // manufacture
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			if ph == ao && (c.results[i].data[j].end-c.results[i].data[j].start <= 5) {
				//   removing fake ao in cut/coat
				phonIndex := ruleIndex - 1
				if phonIndex > 0 && phonIndex < (len(c.phons)-1) && c.phons[phonIndex+1] == l { // ao getting swallowed up by following L as in walter
					// do nothing drop to bottom
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			if ph == ao && (c.results[i].data[j].end-c.results[i].data[j].start > 55) {
				//   removing fake ao in animal2_deven
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			if ph == ah && (ruleIndex == len(c.results[i].rules)-2) && (c.results[i].data[j].end-c.results[i].data[j].start <= 3) { //"ah" in the last position must be greater than 3

				phonIndex := ruleIndex - 1

				if phonIndex > 0 && phonIndex < (len(c.phons)-1) && (c.phons[phonIndex-1] == "r") && (c.phons[phonIndex+1] == f) {
					//  keep  ... ie falls down to lenPh
				} else if phonIndex > 0 && phonIndex < (len(c.phons)-1) && (c.phons[phonIndex-1] == "v") && (c.phons[phonIndex+1] == s) { // coversation
					//
				} else {
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}

			lenPh := c.results[i].data[j].end - c.results[i].data[j].start
			if ph == r.phon || contains(r.otherPhons, ph) {
				ref := psPhonemeDatumRef{
					i, j,
				}
				if _, ok := alignedPhs[2][ph]; ok {
					if lenPh != 0 {
						alignedPhs[2][ph] = append(alignedPhs[2][ph], ref)
					}
				} else {
					if lenPh != 0 {
						alignedPhs[2][ph] = []psPhonemeDatumRef{
							ref,
						}
					} else {
						alignedPhs[2][ph] = []psPhonemeDatumRef{}
					}
				}
				continue
			}

			// If we get here then this is an unexpected phoneme
			if lenPh <= 3 { // ************************************************** may need adjusting
				// It's a short unexpected phoneme so probably not really there, remove
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			// Work out where to store it
			var storeIndex = -1
			switch result[ruleIndex].end - result[ruleIndex].start + 1 {
			case 1:
				// Should never happen if this is an unexpected phoneme
				break
			case 2:
				storeIndex = 3
			case 3:
				if j-result[ruleIndex].start < 2 {
					storeIndex = 0
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 1
					break
				}
			case 4:
				if isVowel(ph) {

				}
				if j-result[ruleIndex].start < 2 {
					if isVowel(ph) {
						storeIndex = 0
					} else {
						storeIndex = 1
					}
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 3
					break
				}
			case 5:
				if j-result[ruleIndex].start < 2 {
					storeIndex = 0
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 1
					break
				}
			case 6:
				if j-result[ruleIndex].start < 2 {
					storeIndex = 0
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 1
					break
				}
				if j-result[ruleIndex].start < 7 {
					storeIndex = 3
					break
				}
			default:
				// Should never happen
				break
			}
			if storeIndex == -1 {
				// Something's not right with the code above so crash the program
				c.print()
				_debug("standardPhoneme: failed to set storeIndex")
				log.Panic()
			}
			ref := psPhonemeDatumRef{
				i, j,
			}
			if _, ok := alignedPhs[storeIndex][ph]; ok {
				alignedPhs[storeIndex][ph] = append(alignedPhs[storeIndex][ph], ref)
			} else {
				alignedPhs[storeIndex][ph] = []psPhonemeDatumRef{
					ref,
				}
			}
		}
	}

	vs := []linkedPhonVerdict{}
	for i, alignedPh := range alignedPhs {
		sortedPhs := c.timeOrderAlignedPhons(alignedPh)
		for _, ph := range sortedPhs {
			count := alignedPh[ph]
			// for ph, count := range alignedPh {
			var v verdict
			if i == 2 {
				// We've got a couple of possibilities here. it's the expected phoneme
				// or one of the other phonemes
				if ph == r.phon {
					if len(count) > 3 {
						v = good
					}
					if len(count) == 3 {
						v = possible
					}
					if len(count) < 3 {
						v = missing
					}
					lv := linkedPhonVerdict{
						phonVerdict{
							r.phon,
							v,
						},
						alignedPh[ph],
					}
					vs = append(vs, lv)
				} else {
					if len(count) > 2 {
						lv := linkedPhonVerdict{
							phonVerdict{
								ph,
								surprise,
							},
							alignedPh[ph],
						}
						vs = append(vs, lv)
					}
				}
			} else {
				// This is an unexpected phoneme
				if len(count) > 2 {
					lv := linkedPhonVerdict{
						phonVerdict{
							ph,
							surprise,
						},
						alignedPh[ph],
					}
					vs = append(vs, lv)
				}
			}
		}
	}

	/*
	  vs := []phonVerdict{}
	  unexpecteds := make(map[phoneme]int)
	  expectedCount := 0
	  for _, alignedPh := range alignedPhs {
	    for ph, count := range alignedPh {
	      if ph != r.phon {
	        if _, ok := unexpecteds[ph]; ok {
	          unexpecteds[ph] += count
	        } else {
	          unexpecteds[ph] = count
	        }
	      } else {
	        expectedCount = count
	      }
	    }
	  }

	  var v verdict
	  if expectedCount > 3 {
	    v = good
	  }
	  if expectedCount == 3 {
	    v = possible
	  }
	  if expectedCount < 3 {
	    v = missing
	  }
	  phV := phonVerdict{
	    r.phon,
	    v,
	  }
	  vs = append(vs, phV)

	  for ph, count := range unexpecteds {
	    if count > 2 {
	      phV := phonVerdict{
	        ph,
	        surprise,
	      }
	      vs = append(vs, phV)
	    }
	  }
	*/
	return vs
}

func makeDiphthong(ph1, ph2 phoneme) (phoneme, bool) {
	for k, v := range diphthongs {
		if v[0] == ph1 && v[1] == ph2 {
			return k, true
		}
	}
	// Something went wrong so just return any old phoneme and false. The caller
	// should not use the phoneme value anyway
	return aa, false
}

// The code below is very similar to that for the standardPhoneme above and is
// tightly bound to the grammar that was used to generate pocketsphinx results.
// The current grammar for a (expected) diphthong phoneme rule is:
//
// [g + v] [g + c] ((p1 p2) | (g + a))
//
// where g = guard, v = vowel, c = consonant, p1, p2 = expected phonemes that
// make up the diphthong vowel, a = any (so vowel or consonant)
// For details of how phonemes are aligned into slots see the detailed comment
// standardPhoneme above. The alignment is similar here - the only difference
// being that we have two (expected) phoneme slots that together make up the
// diphthong.
//
func (c *combiner) diphthongPhoneme(ruleIndex int) []linkedPhonVerdict {
	r, ok := c.results[0].rules[ruleIndex].(R_diphthongPhoneme)
	if !ok {
		// Something's wrong so crash the program for now
		c.print()
		_debug("diphthongPhoneme: type assertion failed")
		log.Panic()
	}
	alignedPhs := make([]map[phoneme][]psPhonemeDatumRef, 5)
	for i := range alignedPhs {
		alignedPhs[i] = make(map[phoneme][]psPhonemeDatumRef)
		// Initialise alignedPhs with the expected phoneme regardless of what's in
		// the results - (if it doesn't appear at all we won't get an opportunity
		// later on to add it in)
		if i == 2 {
			alignedPhs[i][r.phon1] = []psPhonemeDatumRef{}
		}
		if i == 3 {
			alignedPhs[i][r.phon2] = []psPhonemeDatumRef{}
		}
	}
	for i, result := range c.mapping {
		// We may have no data - typically if pocketsphinx crashes as it
		// occasionally does
		if len(c.results[i].data) == 0 {
			continue
		}
		r, ok := c.results[i].rules[ruleIndex].(R_diphthongPhoneme)
		if !ok {
			// Something's wrong so crash the program for now
			c.print()
			_debug("diphthongPhoneme: type assertion failed")
			log.Panic()
		}
		for j := result[ruleIndex].start; j <= result[ruleIndex].end; j++ {
			datum := c.results[i].data[j]
			ph := datum.phoneme
			if ph == r.guard {
				if datum.end-datum.start < 5 {
					// Throw away the guard (unless it's long) in any result that is
					// returned
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}
			lenPh := c.results[i].data[j].end - c.results[i].data[j].start
			if ph == r.phon1 {
				ref := psPhonemeDatumRef{
					i, j,
				}
				if _, ok := alignedPhs[2][ph]; ok {
					if lenPh != 0 {
						alignedPhs[2][ph] = append(alignedPhs[2][ph], ref)
					}
				} else {
					if lenPh != 0 {
						alignedPhs[2][ph] = []psPhonemeDatumRef{
							ref,
						}
					} else {
						alignedPhs[2][ph] = []psPhonemeDatumRef{}
					}
				}
				continue
			}
			if ph == r.phon2 {
				ref := psPhonemeDatumRef{
					i, j,
				}
				if _, ok := alignedPhs[3][ph]; ok {
					if lenPh != 0 {
						alignedPhs[3][ph] = append(alignedPhs[3][ph], ref)
					}
				} else {
					if lenPh != 0 {
						alignedPhs[3][ph] = []psPhonemeDatumRef{
							ref,
						}
					} else {
						alignedPhs[3][ph] = []psPhonemeDatumRef{}
					}
				}
				continue
			}

			// If we get here then this is an unexpected phoneme
			if lenPh <= 3 { // ************************************************** may need adjusting
				// It's a short unexpected phoneme so probably not really there, remove
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			// Work out where to store it
			var storeIndex = -1
			switch result[ruleIndex].end - result[ruleIndex].start + 1 {
			case 2:
				storeIndex = 4
			case 4:
				if j-result[ruleIndex].start < 2 {
					if isVowel(ph) {
						storeIndex = 0
					} else {
						storeIndex = 1
					}
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 4
					break
				}
			case 6:
				if j-result[ruleIndex].start < 2 {
					storeIndex = 0
					break
				}
				if j-result[ruleIndex].start < 4 {
					storeIndex = 1
					break
				}
				if j-result[ruleIndex].start < 6 {
					storeIndex = 4
					break
				}
			default:
				// Should never happen
				break
			}
			if storeIndex == -1 {
				// Something's not right with the code above so crash the program
				c.print()
				_debug("diphthongPhoneme: failed to set storeIndex")
				log.Panic()
			}
			ref := psPhonemeDatumRef{
				i, j,
			}
			if _, ok := alignedPhs[storeIndex][ph]; ok {
				alignedPhs[storeIndex][ph] = append(alignedPhs[storeIndex][ph], ref)
			} else {
				alignedPhs[storeIndex][ph] = []psPhonemeDatumRef{
					ref,
				}
			}
		}
	}

	vs := []linkedPhonVerdict{}
	for i, alignedPh := range alignedPhs {
		sortedPhs := c.timeOrderAlignedPhons(alignedPh)
		for _, ph := range sortedPhs {
			links := alignedPh[ph]
			// for ph, links := range alignedPh {
			var v verdict
			/*
			   if i == 2 {
			     // Then this is an expected phoneme. It's the first phoneme in a
			     // diphthong though so we need to loop around again to work out the
			     // count for the whole diphthong
			   } else if i == 3 {
			     // Then this is the second expected phoneme which is part of a
			     // diphthong
			     count = min(alignedPhs[2][r.phon1], count)
			     if count > 3 {
			       v = good
			     }
			     if count == 3 {
			       v = possible
			     }
			     if count < 3 {
			       v = missing
			     }
			     ph, ok := makeDiphthong(r.phon1, r.phon2)
			     if !ok {
			       // Something went wrong so crash the program for now
			       c.print()
			       _debug("diphthongPhoneme: failed to make diphthong")
			       log.Panic()
			     }
			     phV := phonVerdict{
			       ph,
			       v,
			     }
			     vs = append(vs, phV)

			   } else {
			*/
			count := len(links)
			if i == 2 || i == 3 {
				// Then this is an expected phoneme
				if count > 3 {
					v = good
				}
				if count == 3 {
					v = possible
				}
				if count < 3 {
					v = missing
				}
				phV := linkedPhonVerdict{
					phonVerdict{
						ph,
						v,
					},
					alignedPh[ph],
				}
				vs = append(vs, phV)
			} else {
				// This is an unexpected phoneme
				if count > 2 {
					lv := linkedPhonVerdict{
						phonVerdict{
							ph,
							surprise,
						},
						alignedPh[ph],
					}
					vs = append(vs, lv)
				}
			}
		}
	}
	/*
	  vs := []phonVerdict{}
	  unexpecteds := make(map[phoneme]int)
	  phon1ExpectedCount := 0
	  phon2ExpectedCount := 0
	  for _, alignedPh := range alignedPhs {
	    for ph, count := range alignedPh {
	      if ph == r.phon1 {
	        phon1ExpectedCount = count
	      } else if ph == r.phon2 {
	        phon2ExpectedCount = count
	      } else {
	        if _, ok := unexpecteds[ph]; ok {
	          unexpecteds[ph] += count
	        } else {
	          unexpecteds[ph] = count
	        }
	      }
	    }
	  }

	  expectedCount := min(phon1ExpectedCount, phon2ExpectedCount)
	  var v verdict
	  if expectedCount > 3 {
	    v = good
	  }
	  if expectedCount == 3 {
	    v = possible
	  }
	  if expectedCount < 3 {
	    v = missing
	  }
	  ph, ok := makeDiphthong(r.phon1, r.phon2)
	  if !ok {
	    // Something went wrong so crash the program for now
	    c.print()
	    _debug("diphthongPhoneme: failed to make diphthong")
	    log.Panic()
	  }
	  phV := phonVerdict{
	    ph,
	    v,
	  }
	  vs = append(vs, phV)

	  for ph, count := range unexpecteds {
	    if count > 2 {
	      phV := phonVerdict{
	        ph,
	        surprise,
	      }
	      vs = append(vs, phV)
	    }
	  }
	*/
	return vs
}

func (c combiner) unexpectedPhoneme(ruleIndex int) []linkedPhonVerdict {
	if len(c.results) == 0 {
		return []linkedPhonVerdict{}
	}
	if ruleIndex < 0 || ruleIndex > len(c.results[0].rules)-1 {
		return []linkedPhonVerdict{}
	}
	// The phonme could be a standard phoneme or a diphthong phoneme. Test for
	// each and treat separately
	_, ok := c.results[0].rules[ruleIndex].(R_phoneme)
	if ok {
		return c.standardPhoneme(ruleIndex)
	}
	_, ok = c.results[0].rules[ruleIndex].(R_diphthongPhoneme)
	if ok {
		return c.diphthongPhoneme(ruleIndex)
	}
	// Something's wrong so crash the program for now
	c.print()
	_debug("unexpectedPhoneme: type assertion failed")
	log.Panic()
	// Stupid compiler insists on a return here even though it's unreachable
	// code
	return []linkedPhonVerdict{}
}

func (c *combiner) unexpectedClosing(ruleIndex int) []linkedPhonVerdict {
	if len(c.results) == 0 {
		return []linkedPhonVerdict{}
	}
	unexpecteds := make(map[phoneme][]psPhonemeDatumRef)
	if ruleIndex < 0 || ruleIndex > len(c.results[0].rules)-1 {
		return []linkedPhonVerdict{}
	}
	if ruleIndex != len(c.results[0].rules)-1 {
		return []linkedPhonVerdict{}
	}
	r, ok := c.results[0].rules[ruleIndex].(R_closing)
	if !ok {
		// Something's wrong so crash the program for now
		c.print()
		_debug("unexpectedClosing: type assertion failed")
		log.Panic()
	}
	for i, result := range c.mapping {
		// We may have no data - typically if pocketsphinx crashes as it
		// occasionally does
		if len(c.results[i].data) == 0 {
			continue
		}
		for j := result[ruleIndex].start; j <= result[ruleIndex].end; j++ {
			datum := c.results[i].data[j]
			ph := datum.phoneme
			if ph == sil {
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			if ph == r.guard {
				if datum.end-datum.start < 7 {
					// Throw away the guard (unless it's long) in any result that is
					// returned
					datum.end = datum.start
					c.results[i].data[j] = datum
					continue
				}
			}
			if datum.end-datum.start <= 3 {
				// It's a short unexpected phoneme so probably not really there, remove
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}
			if datum.end-datum.start > 35 {
				// It's a long unexpected phoneme in the closing rule - probably noise, remove
				datum.end = datum.start
				c.results[i].data[j] = datum
				continue
			}

			ref := psPhonemeDatumRef{
				i, j,
			}
			if _, ok := unexpecteds[ph]; ok {
				unexpecteds[ph] = append(unexpecteds[ph], ref)
			} else {
				unexpecteds[ph] = []psPhonemeDatumRef{
					ref,
				}
			}
		}
	}
	retVerdict := []linkedPhonVerdict{}
	for ph := range unexpecteds {
		if len(unexpecteds[ph]) > 2 {
			// Return a phonVerdict
			lv := linkedPhonVerdict{
				phonVerdict{
					ph,
					surprise,
				},
				unexpecteds[ph],
			}
			retVerdict = append(retVerdict, lv)
		}
	}
	return retVerdict
}

func (c combiner) print() {
	if len(c.results) == 0 {
		_debug("Nothing to print")
		return
	}
	str := fmt.Sprintln("Mapped result =")
	for i, result := range c.mapping {
		str += fmt.Sprintln("result", i, "=")
		for _, parsedRes := range result {
			str += fmt.Sprintln("start:", parsedRes.start, "end:", parsedRes.end, "phonemeFound:", parsedRes.phonemeFound)
		}
	}
	_debug(str)
	str = fmt.Sprintln("Rules =")
	for i, parsRes := range c.results {
		str += fmt.Sprintln("rule", i, "=")
		for _, rule := range parsRes.rules {
			if r, ok := rule.(R_trappedOpening); ok {
				str += fmt.Sprintln("Trapped opening rule:")
				str += fmt.Sprintln("trap:", r.rT.trap)
				str += fmt.Sprintln("rule:", r.rT)
				str += fmt.Sprintln("Opening rule:")
				str += fmt.Sprintln("guard:", r.rO.guard)
				str += fmt.Sprintln("rule:", r.rO)
			}
			if r, ok := rule.(R_phoneme); ok {
				str += fmt.Sprintln("Phoneme rule:")
				str += fmt.Sprintln("guard:", r.guard, "phoneme:", r.phon)
				str += fmt.Sprintln("rule:", r.rule)
			}
			if r, ok := rule.(R_diphthongPhoneme); ok {
				str += fmt.Sprintln("Phoneme rule:")
				str += fmt.Sprintln("guard:", r.guard, "phonemes:", r.phon1, r.phon2)
				str += fmt.Sprintln("rule:", r.rule)
			}
			if r, ok := rule.(R_closing); ok {
				str += fmt.Sprintln("Closing rule:")
				str += fmt.Sprintln("guard:", r.guard)
				str += fmt.Sprintln("rule:", r.rule)
			}

		}
	}
	_debug(str)
}
