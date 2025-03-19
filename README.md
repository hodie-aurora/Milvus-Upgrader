
# Milvus-Upgrader

**Milvus-Upgrader** is a Go-based tool for upgrading Milvus instances (version 2.2.3 and above) in Kubernetes clusters. It streamlines and automates the Milvus upgrade process, supporting both minor and major version upgrades.

---

## Features

- **Supports Operator-Deployed Milvus Clusters**: Currently focuses on upgrading Milvus clusters deployed via the Milvus Operator.
- **Minor Version Upgrades**: Automatically updates the image version and triggers redeployment.
- **Major Version Upgrades**: Feature not yet implemented, planned for future releases (TBD).
- **Command-Line Interface**: Built with Cobra for a simple and intuitive CLI experience.

---

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/hodie-aurora/milvus-upgrader.git
   cd milvus-upgrader
   ```
2. **Build the Program**:

   ```bash
   go build -o milvus-upgrade
   ```
3. **(Optional) Dependency Check**:

   - Ensure Go version 1.16 or higher is installed.
   - Run `go mod tidy` to fetch dependencies.

---

## Usage

Run the tool with the following command:

```bash
./milvus-upgrade upgrade --instance <instance-name> --namespace <namespace> --source-version <current-version> --target-version <target-version>
```

### Parameters

- `-i, --instance`: Name of the Milvus instance (required).
- `-n, --namespace`: Kubernetes namespace (default: `default`).
- `-s, --source-version`: Current Milvus version (required, e.g., `v2.2.3`).
- `-t, --target-version`: Target Milvus version (required, e.g., `v2.2.5`).
- `-f, --force`: Force the upgrade without confirmation.
- `-k, --skip-checks`: Skip pre-upgrade version checks.
- `--kubeconfig`: Path to the kubeconfig file (defaults to `$HOME/.kube/config`).

### Examples

- **Upgrade a Milvus Instance**:

  ```bash
  ./milvus-upgrade upgrade --instance my-release --namespace default --source-version v2.2.3 --target-version v2.2.5
  ```
- **Force Upgrade and Skip Checks**:

  ```bash
  ./milvus-upgrade upgrade -i my-release -n default -s v2.2.3 -t v2.2.5 -f -k
  ```

---

## Notes

- **Prerequisites**:
  - A Kubernetes cluster with the Milvus Operator installed.
  - The target Milvus instance must be deployed and running.
- **Current Limitations**:
  - Major version upgrades (e.g., v2.x to v3.x) are not yet implemented and will return a “not supported” message.
  - Helm-deployed Milvus clusters are not supported.
- **Verify the Upgrade**:
  - After upgrading, check the result with:
    ```bash
    kubectl get milvuses my-release -n default -o jsonpath='{.spec.components.image}'
    ```

---

## Debugging and Logging

- The tool logs Kubernetes client initialization and upgrade steps during execution.
- If issues arise (e.g., “server could not find the requested resource”), verify:
  1. `kubectl get crd milvuses.milvus.io` exists.
  2. `kubectl get milvuses <instance-name> -n <namespace>` returns the instance details.

---

## Contributing

Contributions are welcome! Please submit issues or PRs, ensuring code adheres to Go conventions. Run the following before submitting:

```bash
go fmt ./...
go test ./...
```
