package warewulf

import (
	"fmt"
)

func (r *Role) IdString() string {
	return r.Name
}

func (r *Role) Equals(other interface{}) (bool, string) {
	otherRole, ok := other.(*Role)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if r.BootstrapName != otherRole.BootstrapName ||
		r.VnfsName != otherRole.VnfsName {
		return false, "Role Vnfs, or Bootstrap mismatch"
	}

	if equal, reason := StringListEquals(r.FileNames, otherRole.FileNames); !equal {
		return false, fmt.Sprintf("List of files doesn't match: %s", reason)
	}

	if equal, reason := StringListEquals(r.Groups, otherRole.Groups); !equal {
		return false, fmt.Sprintf("List of groups doesn't match: %s", reason)
	}

	return true, ""
}

func StringListEquals(l1, l2 []string) (bool, string) {
	if len(l1) != len(l2) {
		return false, fmt.Sprintf("length of lists doesn't match: %d != %d", len(l1), len(l2))
	}

	l1Map := make(map[string]struct{}, 0)
	for _, v := range l1 {
		l1Map[v] = struct{}{}
	}
	for _, v := range l2 {
		if _, ok := l1Map[v]; !ok {
			return false, fmt.Sprintf("%s not found", v)
		}
	}

	return true, ""
}

func (r RoleList) RequiredObjectsOnly(db *DB) (*DB, error) {
	resultDB := NewDB()

	for _, role := range r {
		vnfs, ok := db.Vnfs[role.VnfsName]
		if !ok {
			return nil, fmt.Errorf("vnfs not found")
		}
		resultDB.Vnfs[role.VnfsName] = vnfs

		bootstrap, ok := db.Bootstraps[role.BootstrapName]
		if !ok {
			return nil, fmt.Errorf("bootstrap not found")
		}

		resultDB.Bootstraps[role.BootstrapName] = bootstrap

		for _, fName := range role.FileNames {
			file, ok := db.Files[fName]
			if !ok {
				return nil, fmt.Errorf("File not found: %s", fName)
			}
			resultDB.Files[fName] = file
		}
		for _, n := range db.Nodes {
			if n.RoleName == role.Name {
				n.Role = role
				resultDB.Nodes[n.IdString()] = n
			}
		}
	}

	return resultDB, nil
}
