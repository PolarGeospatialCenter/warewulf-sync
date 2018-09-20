package warewulf

import "testing"

func TestParseJsonDumpEmpty(t *testing.T) {
	output := ""
	objs, err := ParseJsonDump([]byte(output))
	if err != nil {
		t.Errorf("Got unexpected error parsing empty result: %v", err)
	}

	if objs == nil {
		t.Errorf("Got nil result")
	}

	if len(objs) != 0 {
		t.Errorf("Expecting zero length result: got %d", len(objs))
	}
}
