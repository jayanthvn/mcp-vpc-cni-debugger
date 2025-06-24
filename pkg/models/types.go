
package models

type VpcCniPodNetwork struct {
    PodName      string         `json:"podName"`
    Namespace    string         `json:"namespace"`
    PodIP        string         `json:"podIP"`
    ENI          *ENIInfo       `json:"eni,omitempty"`
    IPAM         *IPAMEntry     `json:"ipam,omitempty"`
    RouteRules   []LinuxRoute   `json:"routeRules,omitempty"`
    IPTablesSNAT []string       `json:"iptablesSnat,omitempty"`
    NodeInfo     *NodeMetadata  `json:"nodeInfo,omitempty"`
    Anomalies    []string       `json:"anomalies,omitempty"`
}

type ENIInfo struct {
    ENIID     string   `json:"eniId"`
    Device    string   `json:"device"`
    MAC       string   `json:"mac"`
    Subnet    string   `json:"subnet"`
    VPC       string   `json:"vpc"`
    SGIDs     []string `json:"sgIds"`
}

type IPAMEntry struct {
    IP        string `json:"ip"`
    Subnet    string `json:"subnet"`
    Allocated bool   `json:"allocated"`
}

type LinuxRoute struct {
    Rule  string `json:"rule"`
    Route string `json:"route"`
}

type NodeMetadata struct {
    InstanceID    string `json:"instanceId"`
    Hostname      string `json:"hostname"`
    KernelVersion string `json:"kernelVersion"`
}
