# kubectl-toolbox-plugin 🧰

A powerful kubectl plugin that provides essential toolbox utilities for Kubernetes cluster management and diagnostics.

## Overview

`kubectl-toolbox-plugin` is a comprehensive CLI tool designed to simplify Kubernetes operations and provide advanced cluster diagnostics capabilities. Built as a kubectl plugin, it seamlessly integrates with your existing Kubernetes workflow while offering enhanced functionality for cluster management.

## Features

- 🔍 **Cluster Health Checks**: Comprehensive cluster diagnostics and health monitoring
- 📦 **Pod Management**: Advanced pod operations and status monitoring
- 🔐 **External Secrets Integration**: External secrets operator compatibility checks
- 🎯 **Version Validation**: Kubernetes version compatibility verification
- ⚙️ **Configurable Environment**: Flexible configuration through environment variables
- 🚀 **Fast & Lightweight**: Minimal resource footprint with maximum performance

## Installation

### Prerequisites

- Go 1.23+
- kubectl installed and configured
- Access to a Kubernetes cluster

## 📦 Installation

```bash
git clone https://github.com/your-org/kubectl-toolbox-plugin.git
cd kubectl-toolbox-plugin
go build -o kubectl-toolbox
sudo mv kubectl-toolbox /usr/local/bin/
```

### Using the Build Script

```bash
chmod +x build.sh
./build.sh
```

## 🛠️ Usage

Once installed, you can use the plugin directly with kubectl:

```bash
kubectl toolbox --help
```

### Available Commands

#### Health Checks
```bash
# Run comprehensive cluster checks
kubectl toolbox check all

# Check pod status
kubectl toolbox check pods

# Verify cluster version compatibility
kubectl toolbox check version

# Validate external secrets configuration
kubectl toolbox check extsecret
```

#### Configuration

The tool supports various configuration options through environment variables:

- `TB_NAMESPACE`: Default namespace for operations
- `TB_KUBECONFIG`: Path to kubeconfig file
- `TB_KUBECONTEXT`: Kubernetes context to use
- `TB_LOGLEVEL`: Logging level (debug, info, warn, error)

### Command Line Options

```bash
  -n, --namespace string     namespace scope for this request
      --kubeconfig string    path to the kubeconfig file  
      --kube-context string  name of the kubeconfig context to use
  -l, --loglevel string      Log level (default "info")
```

## Configuration

### Environment Variables

The plugin reads configuration from the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `TB_NAMESPACE` | Default Kubernetes namespace | `default` |
| `TB_KUBECONFIG` | Path to kubeconfig file | `~/.kube/config` |
| `TB_KUBECONTEXT` | Kubernetes context name | Current context |
| `TB_LOGLEVEL` | Log level | `info` |

### Example Configuration

```bash
export TB_NAMESPACE=production
export TB_LOGLEVEL=debug
export TB_KUBECONTEXT=my-cluster-context

kubectl toolbox check all
```

## Container Support

The project includes container support with the provided `Containerfile`:

```bash
# Build container image
podman build -t kubectl-toolbox-plugin .

# Run in container
podman run --rm -v ~/.kube:/tmp/.kube kubectl-toolbox-plugin check all
podman run --rm -v ~/.kube:/tmp/.kube kubectl-toolbox-plugin check all --kubeconfig=/tmp/.kube/cluster.yml
```

## 🤝 Contributing
Contributions are welcome!  
Please open an issue or submit a PR to discuss improvements or new features.