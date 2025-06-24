
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/jayanthvn/mcp-vpc-cni-debugger/pkg/collectors"
)

func main() {
    r := gin.Default()

    r.GET("/mcp/network/pod/:namespace/:name", func(c *gin.Context) {
        ns := c.Param("namespace")
        name := c.Param("name")

        report, err := collectors.GenerateVpcCniPodReport(name, ns)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, report)
    })

    r.Run(":8080")
}
