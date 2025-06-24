
package collectors

import (
    "fmt"
    "os/exec"
    "strings"
    "mcp-vpc-cni-debugger/pkg/models"
)

func GetRoutingInfo() ([]models.LinuxRoute, error) {
    var routes []models.LinuxRoute

    ruleOut, err := exec.Command("ip", "rule").Output()
    if err != nil {
        return nil, err
    }
    for _, line := range strings.Split(string(ruleOut), "\n") {
        if strings.TrimSpace(line) == "" {
            continue
        }
        routeOut, err := exec.Command("ip", "route", "show", "table", "main").Output()
        if err != nil {
            continue
        }
        routes = append(routes, models.LinuxRoute{
            Rule:  line,
            Route: strings.TrimSpace(string(routeOut)),
        })
    }

    return routes, nil
}
