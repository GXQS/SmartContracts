package vm

import "github.com/GXQS/SmartContracts/pkg/state"

type Reverter struct {
	snapshots []state.SnapshotID
	db        state.Database
}

func NewReverter(db state.Database) *Reverter { return &Reverter{db: db} }

func (r *Reverter) Begin() {
	r.snapshots = append(r.snapshots, r.db.Snapshot())
}

func (r *Reverter) Revert() {
	if len(r.snapshots) == 0 {
		return
	}
	id := r.snapshots[len(r.snapshots)-1]
	r.snapshots = r.snapshots[:len(r.snapshots)-1]
	r.db.RevertToSnapshot(id)
}

func (r *Reverter) Commit() {
	if len(r.snapshots) > 0 {
		r.snapshots = r.snapshots[:len(r.snapshots)-1]
	}
}
