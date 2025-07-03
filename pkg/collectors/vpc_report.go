package collectors

import (
    "context"
    "errors"
    "os/exec"
    "strconv"
    "strings"
    "fmt"

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

    // Step 5: Return structured response
    return &models.VpcCniPodNetwork{
        PodName:      podName,
        Namespace:    namespace,
        PodIP:        podIP,
        PodPID:       pid,
        ENI:          eni,
        IPAM:         ipam,
        RouteRules:   flattenRoutes(routes),
        IPTablesSNAT: snat,
        Anomalies:    []string{}, // placeholder for further validation logic
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
