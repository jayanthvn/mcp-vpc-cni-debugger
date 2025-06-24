
package collectors

import (
    "os/exec"
    "strings"
)

func GetSNATRules() ([]string, error) {
    var snatRules []string
    out, err := exec.Command("iptables", "-t", "nat", "-S").Output()
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        if strings.Contains(line, "AWS-SNAT-CHAIN") {
            snatRules = append(snatRules, line)
        }
    }
    return snatRules, nil
}
