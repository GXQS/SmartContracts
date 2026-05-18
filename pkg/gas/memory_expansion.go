package gas

func (s Schedule) MemoryExpansionCost(before, after uint64) uint64 {
	if after <= before {
		return 0
	}
	return quadCost(after) - quadCost(before)
}

func quadCost(words uint64) uint64 {
	return (words*words)/512 + 3*words
}
