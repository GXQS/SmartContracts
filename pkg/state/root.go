package state

func RootOf(db Database) []byte {
	return db.Root()
}
