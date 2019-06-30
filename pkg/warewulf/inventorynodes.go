package warewulf

import (
	"fmt"
	"log"
	"net"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type InventoryNodeGetter interface {
	GetAll() ([]*types.InventoryNode, error)
}

func (db *DB) LoadNodesFromInventory(inv InventoryNodeGetter, system string) error {

	nodes, err := inv.GetAll()
	if err != nil {
		return err
	}

	wwnodes := make([]*Node, 0)
	for _, node := range nodes {
		if node.System.ID() != system {
			continue
		}
		wwnode := &Node{}
		wwnode.Name = node.Hostname
		wwnode.RoleName = node.Role
		if role, ok := db.Roles[node.Role]; ok {
			wwnode.Role = role
		} else {
			log.Printf("unable to find role '%s' in warewulf configuration repo: skipping node", node.Role)
			continue
		}
		wwnode.IPxeUrl = node.Environment.IPXEUrl

		if node.Metadata != nil {
			if console, ok := node.Metadata.GetString("serial_console"); ok {
				wwnode.Console = console
			}
		}

		if wwmaster, ok := node.Environment.Metadata.GetString("wwmaster"); ok {
			wwnode.Master = wwmaster
		}

		wwnode.PostNetDown = true
		wwnode.Role = db.Roles[node.Role]
		wwnode.Interfaces = make(NetDevList, 0, len(node.Networks))
		for netname, iface := range node.Networks {
			netDev := &NetDev{}

			netDev.Interface = fmt.Sprintf("%s0", netname)

			if len(iface.Config.IP) > 0 {
				ip, mask, err := net.ParseCIDR(iface.Config.IP[0])
				if err != nil {
					return err
				}
				netDev.Ip = ip.String()
				if subnetMask := net.IP(mask.Mask).To4(); subnetMask != nil {
					netDev.Netmask = subnetMask.String()
				}
				if gwIp := net.ParseIP(iface.Config.Gateway[0]); gwIp != nil {
					netDev.Gateway = gwIp.String()
				}
			}

			if len(iface.Interface.NICs) > 0 {
				netDev.HwAddr = iface.Interface.NICs[0].String()
			}

			netDev.MTU = fmt.Sprintf("%d", iface.Network.MTU)

			wwnode.Interfaces = append(wwnode.Interfaces, netDev)
		}
		wwnode.LastModified = &node.LastUpdated
		wwnodes = append(wwnodes, wwnode)
	}
	db.Nodes = make(map[string]*Node, len(wwnodes))
	for _, n := range wwnodes {
		db.Nodes[n.IdString()] = n
	}
	return nil

}
