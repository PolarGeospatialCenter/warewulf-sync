package warewulf

import "fmt"

type Node struct {
	Name        string     `yaml:"hostname"`
	Interfaces  NetDevList `yaml:"interfaces"`
	RoleName    string     `yaml:"role"`
	Console     string     `yaml:"console"`
	PxeLoader   string     `yaml:"pxeloader"`
	IPxeUrl     string     `yaml:ipxeurl`
	PostNetDown bool       `yaml:"postnetdown"`
	Role        *Role      `yaml:"-"`
}

type NetDevList []*NetDev

type NetDev struct {
	Interface string `yaml:"interface"`
	Ip        string `yaml:"ip"`
	HwAddr    string `yaml:"mac"`
	MTU       string `yaml:"mtu"`
	Netmask   string
	Gateway   string
}

type File struct {
	Name        string
	Source      string
	Destination string
	Mode        int
	Owner       int
	Group       int
}

type Vnfs struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
}

type Bootstrap struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
}

type Role struct {
	Name          string
	VnfsName      string   `yaml:"vnfs"`
	BootstrapName string   `yaml:"bootstrap"`
	FileNames     []string `yaml:"files"`
	Groups        []string
}

type RoleList []*Role

type DB struct {
	Nodes      map[string]*Node
	Files      map[string]*File
	Vnfs       map[string]*Vnfs
	Bootstraps map[string]*Bootstrap
	Roles      map[string]*Role
}

func NewDB() *DB {
	db := &DB{}
	db.Files = make(map[string]*File)
	db.Nodes = make(map[string]*Node)
	db.Vnfs = make(map[string]*Vnfs)
	db.Bootstraps = make(map[string]*Bootstrap)
	db.Roles = make(map[string]*Role)
	return db
}

func (db *DB) String() string {
	return fmt.Sprintf("%d nodes, %d files, %d vnfs, %d bootstraps", len(db.Nodes), len(db.Files), len(db.Vnfs), len(db.Bootstraps))
}
