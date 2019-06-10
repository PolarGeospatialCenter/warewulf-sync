package warewulf

import (
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v2"
)

func LoadYaml(yamlPath string) (*DB, error) {
	log.Printf("Loading yaml from %s", yamlPath)
	db := NewDB()

	files := make([]*File, 0)
	err := LoadYamlFile(path.Join(yamlPath, "files.yml"), &files)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		f.ResolveRelativePaths(yamlPath)
		db.Files[f.IdString()] = f
	}

	nodes := make([]*Node, 0)
	err = LoadYamlFile(path.Join(yamlPath, "nodes.yml"), &nodes)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		db.Nodes[n.IdString()] = n
	}

	vnfs := make([]*Vnfs, 0)
	err = LoadYamlFile(path.Join(yamlPath, "vnfs.yml"), &vnfs)
	if err != nil {
		return nil, err
	}

	for _, v := range vnfs {
		db.Vnfs[v.IdString()] = v
	}

	bootstrap := make([]*Bootstrap, 0)
	err = LoadYamlFile(path.Join(yamlPath, "bootstraps.yml"), &bootstrap)
	if err != nil {
		return nil, err
	}

	for _, b := range bootstrap {
		db.Bootstraps[b.IdString()] = b
	}

	roles := RoleList{}
	err = LoadYamlFile(path.Join(yamlPath, "roles.yml"), &roles)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		db.Roles[role.IdString()] = role
	}
	return db, nil
}

func LoadYamlFile(filePath string, result interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, result)
	return err
}
