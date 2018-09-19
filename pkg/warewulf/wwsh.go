package warewulf

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type wwObjId int

func IdFromString(id string) wwObjId {
	objId, _ := strconv.Atoi(id)
	return wwObjId(objId)
}

type wwObj map[string]interface{}

func (o wwObj) ID() wwObjId {
	return wwObjId(o["_ID"].(float64))
}

func ObjJsonDump(args ...string) ([]wwObj, error) {
	cmd := []string{"wwsh", "object", "jsondump"}
	cmd = append(cmd, args...)
	c := exec.Command(cmd[0], cmd[1:]...)

	output, err := c.Output()
	if err, ok := err.(*exec.ExitError); !ok && err != nil {
		return nil, err
	}

	result := make([]wwObj, 0)
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal output of wwsh object jsondump: %v", err)
	}
	return result, nil
}

type idNameMap map[wwObjId]string

type wwDBLoader struct {
	VnfsMap      idNameMap
	BootstrapMap idNameMap
	FileMap      idNameMap
}

func LoadWwshDB() (*DB, error) {
	loader := &wwDBLoader{
		VnfsMap:      idNameMap{},
		BootstrapMap: idNameMap{},
		FileMap:      idNameMap{},
	}

	db := NewDB()

	files, err := loader.getFiles()
	if err != nil {
		return nil, err
	}
	db.Files = files

	bootstraps, err := loader.getBootstraps()
	if err != nil {
		return nil, err
	}
	db.Bootstraps = bootstraps

	vnfs, err := loader.getVnfs()
	if err != nil {
		return nil, err
	}
	db.Vnfs = vnfs

	nodes, err := loader.getNodes()
	if err != nil {
		return nil, err
	}
	db.Nodes = nodes

	return db, nil
}

func (l *wwDBLoader) getNodes() (map[string]*Node, error) {
	objs, err := ObjJsonDump("-t", "node")
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]*Node, len(objs))
	for _, obj := range objs {
		n := NewNodeFromWWObject(obj, l.FileMap, l.BootstrapMap, l.VnfsMap)
		nodes[n.IdString()] = n
	}
	return nodes, nil
}

func (l *wwDBLoader) getFiles() (map[string]*File, error) {
	objs, err := ObjJsonDump("-t", "file")
	if err != nil {
		return nil, err
	}

	files := make(map[string]*File, len(objs))
	for _, obj := range objs {
		f := NewFileFromWWObject(obj)
		l.FileMap[obj.ID()] = f.IdString()
		files[f.IdString()] = f
	}
	return files, nil
}

func (l *wwDBLoader) getBootstraps() (map[string]*Bootstrap, error) {
	objs, err := ObjJsonDump("-t", "bootstrap")
	if err != nil {
		return nil, err
	}

	bootstraps := make(map[string]*Bootstrap, len(objs))
	for _, obj := range objs {
		b := NewBootstrapFromWWObject(obj)
		l.BootstrapMap[obj.ID()] = b.IdString()
		bootstraps[b.IdString()] = b
	}
	return bootstraps, nil
}

func (l *wwDBLoader) getVnfs() (map[string]*Vnfs, error) {
	objs, err := ObjJsonDump("-t", "vnfs")
	if err != nil {
		return nil, err
	}

	vnfs := make(map[string]*Vnfs, len(objs))
	for _, obj := range objs {
		v := NewVnfsFromWWObject(obj)
		l.VnfsMap[obj.ID()] = v.IdString()
		vnfs[v.IdString()] = v
	}
	return vnfs, nil
}
