package cmd

import (
	"testing"

	"github.com/PolarGeospatialCenter/warewulf-sync/pkg/warewulf"
)

func TestMakeSyncableMap(t *testing.T) {
	db1 := warewulf.NewDB()
	db1.Files["test"] = &warewulf.File{}

	files := MakeSyncableMap(db1.Files)

	if len(files) != len(db1.Files) {
		t.Errorf("Wrong number of elements returned")
	}
	for _, f := range files {
		t.Log(f)
	}
}
