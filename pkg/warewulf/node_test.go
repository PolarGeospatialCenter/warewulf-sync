package warewulf

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalNodeObject(t *testing.T) {
	testData := `[{"NAME":["test123"],"_TIMESTAMP":1537232705,"_ID":1,"_HWADDR":["00:08:b2:9a:34:f8"],"NODENAME":"test123","_IPADDR":["172.16.0.1"],"NETDEVS":{"ARRAY":[{"NETMASK":"255.255.255.0","GATEWAY":"172.16.0.254","NAME":"eth0","IPADDR":"172.16.0.1","HWADDR":"00:08:b2:9a:34:f8"}]},"ARCH":"x86_64","_TYPE":"node","GROUPS":["foo","bar"],"POSTNETDOWN":1}]`
	obj := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(testData), &obj)
	if err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}
	n := NewNodeFromWWObject(obj[0], idNameMap{}, idNameMap{}, idNameMap{})
	if n.Name != "test123" {
		t.Errorf("Node name doesn't match expected value: %s", n.Name)
	}

	if len(n.Interfaces) != 1 {
		t.Errorf("Wrong number of interfaces on node")
	}

	if n.Role == nil {
		t.Fatalf("Role set to nil")
	}

	if len(n.Role.Groups) != 2 {
		t.Errorf("Wrong number of groups returned")
	}

	if !n.PostNetDown {
		t.Errorf("PostNetDown not set")
	}

}
