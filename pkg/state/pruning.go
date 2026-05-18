package state

func PruneOldSnapshots(db *MemoryDB, keep int) {
	if keep < 0 {
		keep = 0
	}
	if len(db.history) <= keep {
		return
	}
	db.history = db.history[len(db.history)-keep:]
}
