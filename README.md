
# MCP VPC CNI Debugger

A lightweight MCP server to inspect and debug AWS VPC CNI state per node. Designed for use in EKS clusters and integration with Amazon Q CLI.

---

## ✨ Features

- Inspects pod networking state from:
  - AWS VPC CNI IPAM cache
  - ENI metadata via IMDS
  - Host-level routes (`ip rule`, `ip route`)
  - IPTables SNAT configuration
- Exposes a REST API
- Deployable as a DaemonSet
- Integrates with Amazon Q CLI as a tool

---

## 🚀 Getting Started

### 1. Build & Push Docker Image

```bash
make docker-build
make docker-push
```

### 2. Deploy to EKS

```bash
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/daemonset.yaml
```

Optionally expose:
```bash
kubectl apply -f deploy/service.yaml
```

### 3. Test API

```bash
kubectl port-forward -n kube-system pod/<mcp-pod> 8080:8080
curl http://localhost:8080/mcp/network/pod/<namespace>/<pod-name>
```

---

## 🤖 Amazon Q CLI Integration

1. Install Q CLI:

```bash
pip install amazon-q-cli
q configure
```

2. Create tool config:

```yaml
name: mcp-debugger
description: "Query pod network context"
schemaVersion: 1
type: http
configuration:
  method: GET
  url: http://localhost:8080/mcp/network/pod/{{namespace}}/{{podName}}
  inputParameters:
    - name: namespace
      required: true
      type: string
    - name: podName
      required: true
      type: string
```

3. Register tool:

```bash
q tools add --file tool-config.yaml
```

4. Use Q CLI:

```bash
q "Why can’t pod nginx in namespace web connect to the internet?"
```

---

## 🛠️ Project Layout

```
cmd/                # Entry point (main.go)
pkg/
  collectors/       # IPAM, ENI, routes, iptables
  models/           # MCP schema structs
deploy/             # Kubernetes manifests
Dockerfile
Makefile
README.md
```

---

## 📄 License

MIT or Apache 2.0 (your choice)
