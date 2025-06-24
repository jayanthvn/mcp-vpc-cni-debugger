
package collectors

import (
    "context"
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/jayanthvn/mcp-vpc-cni-debugger/pkg/models"
)

func GetENIFromIMDS() (*models.ENIInfo, error) {
    mac, err := getIMDS("network/interfaces/macs/")
    if err != nil {
        return nil, err
    }

    mac = string([]byte(mac)[:len(mac)-1]) // strip trailing /

    subnetId, _ := getIMDS(fmt.Sprintf("network/interfaces/macs/%s/subnet-id", mac))
    vpcId, _ := getIMDS(fmt.Sprintf("network/interfaces/macs/%s/vpc-id", mac))
    eniId, _ := getIMDS(fmt.Sprintf("network/interfaces/macs/%s/interface-id", mac))
    sgIds, _ := getIMDS(fmt.Sprintf("network/interfaces/macs/%s/security-group-ids", mac))

    return &models.ENIInfo{
        ENIID:  eniId,
        Device: "eth0",
        MAC:    mac,
        Subnet: subnetId,
        VPC:    vpcId,
        SGIDs:  []string{sgIds},
    }, nil
}

func getIMDS(path string) (string, error) {
    url := fmt.Sprintf("http://169.254.169.254/latest/meta-data/%s", path)
    req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
    if err != nil {
        return "", err
    }
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body), nil
}
