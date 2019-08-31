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
		gwNetwork := chooseDefaultGatewayNetwork(node.Networks)
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
				if gwIp := net.ParseIP(iface.Config.Gateway[0]); gwIp != nil && gwNetwork == netname {
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

func getGatewayWeight(nic *types.NICInstance) int {
	cidrWeights := map[string]int{
		"192.168.0.0/16": 2,
		"172.16.0.0/12":  2,
		"10.0.0.0/8":     2,
		"fe80::/10":      1,
		"169.254.0.0/16": 1,
		"fc00::/7":       2,
	}

	var weight int
	if len(nic.Config.Gateway) == 0 {
		return 0
	}
	for _, gwCandidate := range nic.Config.Gateway {
		gwCandidateIP := net.ParseIP(gwCandidate)
		if gwCandidateIP == nil {
			continue
		}

		candidateWeight := 4
		for cidr, netWeight := range cidrWeights {
			_, candidateNet, _ := net.ParseCIDR(cidr)
			if candidateNet.Contains(gwCandidateIP) {
				candidateWeight = netWeight
			}
		}
		weight += candidateWeight
	}
	return weight
}

func chooseDefaultGatewayNetwork(networks map[string]*types.NICInstance) string {
	var gwNetwork string
	var gwWeight int

	for network, nic := range networks {
		if candidateWeight := getGatewayWeight(nic); candidateWeight > gwWeight {
			gwWeight = candidateWeight
			gwNetwork = network
		}
	}
	return gwNetwork
}
