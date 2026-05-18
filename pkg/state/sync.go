package state

func DeterministicReplay(target *Trie, updates [][2][]byte) {
	for _, pair := range updates {
		target.Update(pair[0], pair[1])
	}
}
