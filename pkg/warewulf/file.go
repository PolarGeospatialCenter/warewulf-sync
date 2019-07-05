package warewulf

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/go-test/deep"
)

func NewFileFromWWObject(obj map[string]interface{}) *File {
	f := &File{}

	if nameSt, ok := obj["NAME"].(string); ok {
		f.Name = nameSt
	}

	if pathSt, ok := obj["PATH"].(string); ok {
		f.Destination = pathSt
	}

	if originSt, ok := obj["ORIGIN"].(string); ok {
		f.Source = originSt
	}

	if modeSt, ok := obj["MODE"].(string); ok {
		mode, _ := strconv.Atoi(modeSt)
		f.Mode = mode
	}

	if uidSt, ok := obj["UID"].(string); ok {
		uid, _ := strconv.Atoi(uidSt)
		f.Owner = uid
	}

	if gidSt, ok := obj["GID"].(string); ok {
		gid, _ := strconv.Atoi(gidSt)
		f.Group = gid
	}

	return f
}

func (f *File) NewCmd() [][]string {
	return [][]string{
		[]string{"wwsh", "file", "new", f.Name},
	}
}

func (f *File) UpdateCmd() [][]string {
	cmds := [][]string{}
	cmd := []string{
		"wwsh", "file", "set", f.Name,
		"-u", f.UIDString(),
		"-g", f.GIDString(),
	}

	if f.Mode != 0 {
		cmd = append(cmd, "-m", strconv.FormatInt(int64(f.Mode), 8))
	}

	if f.Destination != "" {
		cmd = append(cmd, "-D", f.Destination)
	}

	if f.Source != "" {
		filePath := f.Source
		// add command to download file if it's source is a URL
		if u, err := url.Parse(f.Source); err == nil && u.IsAbs() {
			tmpFile := path.Join("/tmp", fmt.Sprintf("%s.tmp", f.Name))
			cmds = append(cmds, []string{"curl", f.Source, "-o", tmpFile})
			filePath = tmpFile
		}
		cmd = append(cmd, "-o", filePath)
	}
	cmds = append(cmds, cmd)

	cmds = append(cmds, []string{
		"wwsh", "file", "sync", f.Name,
	})

	return cmds

}

func (f *File) DeleteCmd() [][]string {
	return [][]string{
		[]string{"wwsh", "file", "delete", f.Name},
	}
}

func (f *File) IdString() string {
	return f.Name
}

func (f *File) ResolveRelativePaths(base string) {
	if u, err := url.Parse(f.Source); err == nil && u.IsAbs() {
		return
	}

	if !path.IsAbs(f.Source) && f.Source != "" {
		f.Source = path.Join(base, f.Source)
	}
}

func (f *File) UIDString() string {
	return strconv.Itoa(f.Owner)
}

func (f *File) GIDString() string {
	return strconv.Itoa(f.Group)
}

func (f *File) Equals(other interface{}) (bool, string) {
	otherFile, ok := other.(*File)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}
	if diff := deep.Equal(f, otherFile); len(diff) > 0 {
		return false, strings.Join(diff, " ")
	}
	return true, ""
}
