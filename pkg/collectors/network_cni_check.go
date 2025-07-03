package collectors

import (
    "context"
    "os/exec"
    "strconv"
    "strings"

    "github.com/pkg/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// GetPodIPAndPID returns the IP and PID of a running pod
func GetPodIPAndPID(clientset *kubernetes.Clientset, namespace, name string) (string, string, error) {
    pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
        return "", "", errors.Wrap(err, "failed to get pod")
    }

    podIP := pod.Status.PodIP
    if podIP == "" {
        return "", "", errors.New("pod IP not assigned yet")
    }

    containerID := strings.TrimPrefix(pod.Status.ContainerStatuses[0].ContainerID, "containerd://")
    out, err := exec.Command("crictl", "inspect", "--output", "go-template", "--template", "{{.info.pid}}", containerID).CombinedOutput()
    if err != nil {
        return "", "", errors.Wrap(err, "crictl inspect failed: "+string(out))
    }

    pid := strings.TrimSpace(string(out))
    if _, err := strconv.Atoi(pid); err != nil {
        return "", "", errors.New("invalid PID returned")
    }

    return podIP, pid, nil
}

