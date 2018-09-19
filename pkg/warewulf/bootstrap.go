package warewulf

import (
	"fmt"
	"path"
)

func NewBootstrapFromWWObject(obj map[string]interface{}) *Bootstrap {
	b := &Bootstrap{}

	if nameSt, ok := obj["NAME"].(string); ok {
		b.Name = nameSt
	}

	return b
}

func (b *Bootstrap) IdString() string {
	return b.Name
}

func (b *Bootstrap) NewCmd() [][]string {
	tmpFile := path.Join("/tmp", fmt.Sprintf("%s.wwbs", b.Name))
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{
		"curl", b.Source,
		"-o", tmpFile,
	})
	cmds = append(cmds, []string{
		"wwsh", "bootstrap", "import", tmpFile,
		"--name", b.Name,
	})
	return cmds
}

func (b *Bootstrap) UpdateCmd() [][]string {
	cmds := make([][]string, 0)
	return cmds
}

func (b *Bootstrap) DeleteCmd() [][]string {
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{
		"wwsh", "bootstrap", "delete", b.Name,
	})
	return cmds
}

func (b *Bootstrap) Equals(other interface{}) (bool, string) {
	otherBootstrap, ok := other.(*Bootstrap)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if b.Name != otherBootstrap.Name {
		return false, "Bootstrap name mismatch"
	}

	return true, ""
}
