package gas

func StorageAccessCost(warm bool) uint64 {
	if warm {
		return 20
	}
	return 100
}
