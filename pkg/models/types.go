
package models

type VpcCniPodNetwork struct {
    PodName      string                 `json:"podName"`
    Namespace    string                 `json:"namespace"`
    PodIP        string                 `json:"podIP"`
    PodPID       string                 `json:"podPID"`
    ENI          *ENIInfo       `json:"eni,omitempty"`
    IPAM         *IPAMEntry     `json:"ipam,omitempty"`
    RouteRules   []string               `json:"routeRules"`
    IPTablesSNAT []string               `json:"iptablesSNAT"`
    PodNetwork   map[string]interface{} `json:"podNetwork"`
    HostRouting  map[string]interface{} `json:"hostRouting"`
    Anomalies    []string               `json:"anomalies"`
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
