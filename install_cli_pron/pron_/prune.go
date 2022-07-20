package pron

type link struct {
	from, to int
}

type linkSet struct {
	links map[link]bool
}

func sgn(i int) int {
	if i < 0 {
		return -1
	}
	if i == 0 {
		return 0
	}
	return 1
}

type linkWithConfidence struct {
	link
	confidence int
}

type linkWithConfidenceSet map[linkWithConfidence]bool

func createLinkWithConfidenceSet(actual candidateData, expected []phoneme) linkWithConfidenceSet {
	debug("createLinkWithConfidenceSet->")
	ret := make(linkWithConfidenceSet)
	for i, phE := range expected {
		for j, phA := range actual {
			// Coding in the er_ending rule here!
			//
			if i == len(expected)-1 && phE == er && phA.phoneme == ah {
				debug("Creating link between", phE, "and", phA)
				link := linkWithConfidence{
					link{
						i,
						j,
					},
					phA.confidence,
				}
				ret[link] = true
				continue
			}
			if phE == phA.phoneme {
				debug("Creating link between", phE, "and", phA)
				link := linkWithConfidence{
					link{
						i,
						j,
					},
					phA.confidence,
				}
				ret[link] = true
				continue
			}
		}
	}
	debug("createLinkWithConfidenceSet->")
	return ret
}

func pruneWithConfidence(links linkWithConfidenceSet, pa pruneAssist) linkWithConfidenceSet {
	remaining := linkSet{
		make(map[link]bool),
	}
	removed := linkSet{
		make(map[link]bool),
	}
	max := 0
	for key := range links {
		if key.confidence > max {
			max = key.confidence
		}
	}
	linksIn := linkSet{
		make(map[link]bool),
	}
	for link := range links {
		linksIn.links[link.link] = true
	}
	for i := max; i > 0; i-- {
		links_i := linkSet{
			make(map[link]bool),
		}
		for link := range links {
			if link.confidence == i && !removed.contains(link.link) {
				links_i.insert(link.link)
			}
		}
		remaining_i := prune(links_i, pa)
		// A bit of housekeeping...
		remaining.formUnion(remaining_i)
		removed.formUnion(links_i.subtracting(remaining_i))

		// Remove all links with a lower confidence
		for rem := range remaining.links {
			others := linkSet{
				make(map[link]bool),
			}
			for other := range links {
				if other.confidence < i {
					others.links[other.link] = true
				}
			}
			conflicting := conflicts(rem, others)
			removed.formUnion(conflicting)
		}
	}
	// Finally, return all links in remaininf as a linkWithConfidenceSet
	ret := make(linkWithConfidenceSet)
	for link := range links {
		if remaining.contains(link.link) {
			ret[link] = true
		}
	}
	return ret
}

func createLinkSet(actual []psPhonemeDatum, expected []phoneme) linkSet {
	debug("createLinkSet->")
	ret := linkSet{
		make(map[link]bool),
	}
	for i, phE := range expected {
		for j, phA := range actual {
			// Coding in the er_ending rule here!
			//
			if i == len(expected)-1 && phE == er && phA.phoneme == ah {
				debug("Creating link between", phE, "and", phA)
				ret.insert(link{
					i,
					j,
				})
				continue
			}
			if phE == phA.phoneme {
				debug("Creating link between", phE, "and", phA)
				ret.insert(link{
					i,
					j,
				})
				continue
			}
		}
	}
	debug("Links created =", ret)
	debug("createLinkSet->")
	return ret
}

func (l linkSet) subtracting(m linkSet) linkSet {
	ret := linkSet{
		make(map[link]bool),
	}
	for k := range l.links {
		if _, ok := m.links[k]; !ok {
			ret.insert(k)
		}
	}
	return ret
}

func (l *linkSet) formUnion(with linkSet) {
	for k := range with.links {
		l.insert(k)
	}
}

func (l *linkSet) insert(k link) {
	l.links[k] = true
}

func (l *linkSet) remove(k link) {
	delete(l.links, k)
}

func (l linkSet) contains(k link) bool {
	_, ok := l.links[k]
	return ok
}

func (l linkSet) linkWithFrom(from int) (link, bool) {
	for k := range l.links {
		if k.from == from {
			return k, true
		}
	}
	return link{}, false
}

func fromConflicts(l link, others linkSet) linkSet {
	ret := linkSet{
		make(map[link]bool),
	}
	for k := range others.links {
		if l == k {
			// A link cannot conflict with itself
			//
			continue
		}
		if l.from == k.from {
			ret.insert(k)
		}
	}
	return ret
}

func toConflicts(l link, others linkSet) linkSet {
	ret := linkSet{
		make(map[link]bool),
	}
	for k := range others.links {
		if l == k {
			// A link cannot conflict with itself
			//
			continue
		}
		if l.to == k.to {
			ret.insert(k)
		}
	}
	return ret
}

func crossoverConflicts(l link, others linkSet) linkSet {
	ret := linkSet{
		make(map[link]bool),
	}
	for k := range others.links {
		if l == k {
			continue
		}
		if sgn(l.from-k.from)+sgn(l.to-k.to) == 0 {
			ret.insert(k)
		}
	}
	return ret
}

func conflicts(l link, others linkSet) linkSet {
	ret := fromConflicts(l, others)
	ret.formUnion(toConflicts(l, others))
	ret.formUnion(crossoverConflicts(l, others))
	return ret
}

func costToKeep(l linkSet) int {
	return len(l.links)
}

func prune(links linkSet, pa pruneAssist) linkSet {
	debug("prune:->")
	remaining := linkSet{
		make(map[link]bool),
	}
	removed := linkSet{
		make(map[link]bool),
	}

	for k := range links.links {
		debug("Processing link k =", k)
		if removed.contains(k) {
			continue
		}
		confLinkSet := conflicts(k, links).subtracting(removed)
		if len(confLinkSet.links) == 0 {
			// There are no conflicting links so keep this one and move onto the
			// next link.
			//
			remaining.insert(k)
			continue
		}

		debug("Conflicting link set =", confLinkSet)
		thisCost := costToKeep(confLinkSet)
		debug("Cost to keep link,", k, "=", thisCost)

		type linksWithCost struct {
			cost int
			linkSet
		}
		minCostLinks := linksWithCost{
			thisCost,
			linkSet{
				make(map[link]bool),
			},
		}
		for k := range confLinkSet.links {
			linksConf_k := conflicts(k, links).subtracting(removed)
			cost := costToKeep(linksConf_k)
			debug("Cost to keep link,", k, "=", cost)
			if cost < minCostLinks.cost {
				links := make(map[link]bool)
				links[k] = true
				minCostLinks = linksWithCost{
					cost,
					linkSet{
						links,
					},
				}
				continue
			}
			if cost == minCostLinks.cost {
				minCostLinks.links[k] = true
			}
		}
		minCost := minCostLinks.cost
		if thisCost < minCost {
			debug("Keeping", k)
			remaining.insert(k)
			debug("Removing", confLinkSet)
			removed.formUnion(confLinkSet)
			continue
		}
		if thisCost > minCost {
			debug("Removing", k)
			removed.insert(k)
			continue
		}
		if thisCost == minCost {
			l := pa.resolve(k, minCostLinks.linkSet)
			if l == k {
				debug("Keeping", k)
				remaining.insert(k)
				debug("Removing", confLinkSet)
				removed.formUnion(confLinkSet)
				continue
			}
		}
		debug("Removing", k)
		removed.insert(k)
	}
	debug("remaining =", remaining)
	return remaining
}
