package collectors

import (
    "context"
    "errors"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "fmt"
    "crypto/sha1"
    "encoding/hex"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    "github.com/jayanthvn/mcp-vpc-cni-debugger/pkg/models"
)

// GenerateVpcCniPodReport gathers VPC CNI-related network diagnostics for a given pod
func GenerateVpcCniPodReport(podName, namespace string) (*models.VpcCniPodNetwork, error) {
    // Step 1: Access Kubernetes API
    cfg, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        return nil, err
    }

    // Step 2: Get Pod Info
    pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    podIP := pod.Status.PodIP
    if podIP == "" {
        return nil, errors.New("pod IP not assigned yet")
    }

    // Step 3: Get container PID using crictl
    containerID := strings.TrimPrefix(pod.Status.ContainerStatuses[0].ContainerID, "containerd://")
    fmt.Printf("Using containerID: %s\n", containerID)
    out, err := exec.Command("crictl", "--runtime-endpoint", "unix:///run/containerd/containerd.sock", "inspect", "--output", "go-template", "--template", "{{.info.pid}}", containerID).CombinedOutput()

    if err != nil {
        return nil, errors.New("crictl inspect failed: " + string(out))
    }

    pid := strings.TrimSpace(string(out))
    if _, err := strconv.Atoi(pid); err != nil {
        return nil, errors.New("invalid PID returned")
    }

    // Step 4: Collect network diagnostics (IPAM, ENI, routes, SNAT)
    ipam, _ := GetIPAMEntry(podIP)
    eni, _ := GetENIFromIMDS()
    routes, _ := GetRoutingInfo()
    snat, _ := GetSNATRules()

    // Step 5: Check static ARP and compare MAC
    hasARP, podMAC, err := CheckStaticARP(pid)
    arpStatus := fmt.Sprintf("Static ARP 169.254.1.1 present: %v, MAC: %s", hasARP, podMAC)
    if err != nil {
	    arpStatus = fmt.Sprintf("Static ARP check failed: %v", err)
    }

    // Get host veth MAC
    hostMAC, err := GetMACForHostVeth(namespace, podName)
    if err != nil {
	    arpStatus += fmt.Sprintf("; Host veth MAC check failed: %v", err)
    } else if hasARP && podMAC != hostMAC {
	    arpStatus += fmt.Sprintf("; MISMATCH with host-side veth MAC: %s", hostMAC)
    } else if hasARP {
	    arpStatus += "; MAC matches host-side veth"
    }

    // Step 6: Return structured response
    return &models.VpcCniPodNetwork{
        PodName:      podName,
        Namespace:    namespace,
        PodIP:        podIP,
        PodPID:       pid,
        ENI:          eni,
        IPAM:         ipam,
        RouteRules:   flattenRoutes(routes),
        IPTablesSNAT: snat,
        Anomalies:    []string{arpStatus},
    }, nil
}

func flattenRoutes(routes []models.LinuxRoute) []string {
    flat := []string{}
    for _, r := range routes {
        entry := fmt.Sprintf("rule: %s, route: %s", r.Rule, r.Route)
        flat = append(flat, entry)
    }
    return flat
}

// CheckStaticARP verifies if 169.254.1.1 is a static ARP entry inside the pod's network namespace.
func CheckStaticARP(pid string) (found bool, mac string, err error) {
    arpPath := fmt.Sprintf("/proc/%s/net/arp", pid)
    data, err := os.ReadFile(arpPath)
    if err != nil {
        return false, "", fmt.Errorf("failed to read ARP table: %w", err)
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines[1:] { // skip header
        fields := strings.Fields(line)
        if len(fields) >= 6 && fields[0] == "169.254.1.1" {
            if fields[2] == "0x6" { // flags indicating static/permanent
                return true, fields[3], nil
            }
            return false, "", nil // entry exists but not static
        }
    }

    return false, "", nil // not found
}

// GeneratePodHostVethName generates the name for Pod's host-side veth device.
func GeneratePodHostVethName(prefix string, podNamespace string, podName string) string {
    suffix := GeneratePodHostVethNameSuffix(podNamespace, podName)
    return fmt.Sprintf("%s%s", prefix, suffix)
}

// GeneratePodHostVethNameSuffix generates the name suffix for Pod's host-side veth.
func GeneratePodHostVethNameSuffix(podNamespace string, podName string) string {
    h := sha1.New()
    h.Write([]byte(fmt.Sprintf("%s.%s", podNamespace, podName)))
    return hex.EncodeToString(h.Sum(nil))[:11]
}

// GetMACForHostVeth uses the generated host veth name to read its MAC from sysfs
func GetMACForHostVeth(podNamespace, podName string) (string, error) {
    eniName := GeneratePodHostVethName("eni", podNamespace, podName)
    macPath := fmt.Sprintf("/sys/class/net/%s/address", eniName)
    macBytes, err := os.ReadFile(macPath)
    if err != nil {
        return "", fmt.Errorf("failed to read MAC from %s: %w", macPath, err)
    }

    return strings.TrimSpace(string(macBytes)), nil
}
