package warewulf

import (
	"fmt"
	"path"
)

func NewVnfsFromWWObject(obj map[string]interface{}) *Vnfs {
	v := &Vnfs{}

	if nameSt, ok := obj["NAME"].(string); ok {
		v.Name = nameSt
	}

	return v
}

func (v *Vnfs) IdString() string {
	return v.Name
}

func (v *Vnfs) NewCmd() [][]string {
	tmpFile := path.Join("/tmp", fmt.Sprintf("%s.vnfs", v.Name))
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{
		"curl", v.Source,
		"-o", tmpFile,
	})
	cmds = append(cmds, []string{
		"wwsh", "vnfs", "import", tmpFile,
		"--name", v.Name,
	})
	return cmds
}

func (v *Vnfs) UpdateCmd() [][]string {
	cmds := make([][]string, 0)
	return cmds
}

func (v *Vnfs) DeleteCmd() [][]string {
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{
		"wwsh", "vnfs", "delete", v.Name,
	})
	return cmds
}

func (v *Vnfs) Equals(other interface{}) (bool, string) {
	otherVnfs, ok := other.(*Vnfs)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if v.Name != otherVnfs.Name {
		return false, "Vnfs name mismatch"
	}

	return true, ""
}
