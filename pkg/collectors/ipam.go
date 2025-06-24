
package collectors

import (
    "encoding/json"
    "errors"
    "os"
    "github.com/jayanthvn/mcp-vpc-cni-debugger/pkg/models"
)

func GetIPAMEntry(ip string) (*models.IPAMEntry, error) {
    data, err := os.ReadFile("/var/run/aws-node/ipam.json")
    if err != nil {
        return nil, err
    }
    var ipamState map[string]interface{}
    if err := json.Unmarshal(data, &ipamState); err != nil {
        return nil, err
    }

    pool, ok := ipamState["addresses"].(map[string]interface{})
    if !ok {
        return nil, errors.New("unexpected format in IPAM file")
    }

    _, found := pool[ip]
    return &models.IPAMEntry{
        IP: ip,
        Allocated: found,
        Subnet: "",
    }, nil
}
