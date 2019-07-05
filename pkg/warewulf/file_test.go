package warewulf

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalFileObject(t *testing.T) {
	testData := `[{"SIZE":"13","FORMAT":"data","NAME":"test","UID":"0","_TIMESTAMP":1537207319,"_ID":7,"GID":"0","PATH":"/test/output.txt","ORIGIN":"/test","CHECKSUM":"da40f7ebcf60d9491d47a681d0537d8e","_TYPE":"file","FILETYPE":"32768","MODE":"360"}]`
	obj := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(testData), &obj)
	if err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}
	f := NewFileFromWWObject(obj[0])
	if f.Name != "test" {
		t.Errorf("File name doesn't match expected value: %s", f.Name)
	}

}

func TestFileUpdateCmdLocal(t *testing.T) {
	f := File{Name: "test.local", Source: "file/test.local"}
	cmds := f.UpdateCmd()

	if len(cmds) != 2 {
		t.Errorf("Wrong number of commands returned: %v", cmds)
	}
}

func TestFileUpdateCmdRemote(t *testing.T) {
	f := File{Name: "test.remote", Source: "https://foo.test/file/test.remote"}
	cmds := f.UpdateCmd()

	t.Logf("Commands: %v", cmds)
	if len(cmds) != 3 {
		t.Errorf("Wrong number of commands returned: %v", cmds)
	}

	if cmds[0][0] != "curl" {
		t.Errorf("First command should be curl, got: %v", cmds[0])
	}
}
