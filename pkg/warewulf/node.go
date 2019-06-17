package warewulf

import (
	"fmt"
	"strings"

	"github.com/go-test/deep"
)

func NewNodeFromWWObject(obj map[string]interface{}, fileIndex, bootstrapIndex, vnfsIndex idNameMap) *Node {
	n := &Node{}

	if nameSt, ok := obj["NODENAME"].(string); ok {
		n.Name = nameSt
	}

	n.Interfaces = make([]*NetDev, 0)
	if netDevWrapper, ok := obj["NETDEVS"].(map[string]interface{}); ok {
		if netDevs, ok := netDevWrapper["ARRAY"].([]interface{}); ok {
			for _, devI := range netDevs {
				dev, ok := devI.(map[string]interface{})
				if !ok {
					continue
				}
				d := &NetDev{
					Interface: dev["NAME"].(string),
				}

				if v, ok := dev["IPADDR"].(string); ok {
					d.Ip = v
				}
				if v, ok := dev["NETMASK"].(string); ok {
					d.Netmask = v
				}
				if v, ok := dev["GATEWAY"].(string); ok {
					d.Gateway = v
				}
				if v, ok := dev["HWADDR"].(string); ok {
					d.HwAddr = v
				}
				n.Interfaces = append(n.Interfaces, d)
			}
		}
	}

	if console, ok := obj["CONSOLE"].(string); ok {
		n.Console = console
	}

	if postNetDown, ok := obj["POSTNETDOWN"].(float64); ok {
		n.PostNetDown = postNetDown == 1
	}

	role := &Role{}
	attachRole := false

	if groupsI, ok := obj["GROUPS"]; ok {
		if groups, ok := groupsI.([]interface{}); ok {
			role.Groups = make([]string, 0, len(groups))
			for _, g := range groups {
				role.Groups = append(role.Groups, g.(string))
			}
			attachRole = true
		}
	}

	if fileIds, ok := obj["FILEIDS"].([]interface{}); ok {
		role.FileNames = make([]string, 0, len(fileIds))
		for _, fileId := range fileIds {
			role.FileNames = append(role.FileNames, fileIndex[IdFromString(fileId.(string))])
		}
		attachRole = true
	}

	if bootstrapId, ok := obj["BOOTSTRAPID"].(string); ok {
		role.BootstrapName = bootstrapIndex[IdFromString(bootstrapId)]
		attachRole = true
	}

	if vnfsId, ok := obj["VNFSID"].(string); ok {
		role.VnfsName = vnfsIndex[IdFromString(vnfsId)]
		attachRole = true
	}

	if master, ok := obj["MASTER"].(string); ok {
		n.Master = master
	}

	if attachRole {
		n.Role = role
	}

	return n
}

func (n *Node) NewCmd() [][]string {
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{"wwsh", "node", "new", n.Name, "--nodhcp"})
	return cmds
}

func (n *Node) UpdateCmd() [][]string {
	cmds := make([][]string, 0)
	cmds = append(cmds, n.DeleteCmd()...)
	cmds = append(cmds, n.NewCmd()...)
	cmd := []string{"wwsh", "node", "set", n.Name}
	if n.Role != nil {
		if len(n.Role.Groups) > 0 {
			cmd = append(cmd, "-g", strings.Join(n.Role.Groups, ","))
		}
	}
	cmds = append(cmds, cmd)

	if n.Console != "" {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--console", n.Console,
		})
	}

	if n.PostNetDown {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--postnetdown", "1",
		})
	}

	if n.PxeLoader != "" {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--pxeloader", n.PxeLoader,
		})
	}

	if n.Master != "" {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--master", n.Master,
		})
	}

	if n.IPxeUrl != "" {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--ipxeurl", n.IPxeUrl,
		})
	}

	for _, dev := range n.Interfaces {
		cmd := []string{
			"wwsh", "node", "set", n.Name, "--nodhcp",
			"--netdev", dev.Interface,
		}
		if dev.Ip != "" {
			cmd = append(cmd, "--ipaddr", dev.Ip)
		}
		if dev.HwAddr != "" {
			cmd = append(cmd, "--hwaddr", dev.HwAddr)
		}
		if dev.Netmask != "" {
			cmd = append(cmd, "--netmask", dev.Netmask)
		}
		if dev.Gateway != "" {
			cmd = append(cmd, "--gateway", dev.Gateway)
		}
		if dev.MTU != "" {
			cmd = append(cmd, "--mtu", dev.MTU)
		}
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, []string{
		"wwsh", "provision", "set", n.Name,
		"--bootstrap", n.Role.BootstrapName,
		"--vnfs", n.Role.VnfsName,
	})

	if len(n.Role.FileNames) > 0 {
		cmds = append(cmds, []string{
			"wwsh", "provision", "set", n.Name,
			"--files", strings.Join(n.Role.FileNames, ","),
		})
	}

	return cmds
}

func (n *Node) IdString() string {
	return n.Name
}

func (n *Node) DeleteCmd() [][]string {
	cmds := make([][]string, 0)
	cmds = append(cmds, []string{"wwsh", "node", "delete", n.Name})
	return cmds
}

func (n *Node) Equals(other interface{}) (bool, string) {
	otherNode, ok := other.(*Node)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if n.Name != otherNode.Name ||
		n.Console != otherNode.Console {
		return false, "Node Name or Console mismatch"
	}

	if n.PostNetDown != otherNode.PostNetDown {
		return false, "PostNetDown mismatch"
	}

	if n.Role == nil && otherNode.Role != nil ||
		n.Role != nil && otherNode.Role == nil {
		return false, "One role is nil, the other is not"
	}

	if n.Role != nil {
		if equals, reason := n.Role.Equals(otherNode.Role); !equals {
			return equals, reason
		}
	}

	return true, " "
}

func (l NetDevList) Equals(other interface{}) (bool, string) {
	otherList, ok := other.([]*NetDev)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if len(l) != len(otherList) {
		return false, "Element count mismatch"
	}

	otherMap := make(map[string]*NetDev, len(l))
	for _, d := range otherList {
		otherMap[d.Interface] = d
	}

	for _, d := range l {
		otherD, ok := otherMap[d.Interface]
		if !ok {
			return false, "device not found in other"
		}

		if equal, reason := d.Equals(otherD); !equal {
			return false, reason
		}
	}

	return true, ""
}

func (n *NetDev) Equals(other interface{}) (bool, string) {
	otherNetDev, ok := other.(*NetDev)
	if !ok {
		return false, fmt.Sprintf("Wrong type or non-existent %T", other)
	}

	if diff := deep.Equal(n, otherNetDev); len(diff) > 0 {
		return false, strings.Join(diff, " ")
	}
	return true, ""
}
