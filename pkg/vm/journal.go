package vm

type JournalEntry struct {
	Key      string
	Previous []byte
}

type Journal struct {
	entries []JournalEntry
}

func (j *Journal) Append(e JournalEntry) { j.entries = append(j.entries, e) }
func (j *Journal) Entries() []JournalEntry {
	out := make([]JournalEntry, len(j.entries))
	copy(out, j.entries)
	return out
}
func (j *Journal) Reset() { j.entries = j.entries[:0] }
