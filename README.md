# vault-kms-plugin-openshift-provider

An OLM operator that registers the [HashiCorp Vault KMS Plugin](https://developer.hashicorp.com/vault/docs/platform/k8s/kms) as the encryption-at-rest provider for OpenShift clusters.

## Overview

Red Hat's OpenShift encryption controllers (`cluster-kube-apiserver-operator`) discover the KMS plugin image to use by reading a well-known ConfigMap in the `openshift-kms-plugin-provider` namespace. This operator's sole responsibility is to create and continuously enforce that ConfigMap.

```
ConfigMap
  namespace: openshift-kms-plugin-provider
  name:      ibm-kms-vault-plugin-provider
  labels:
    config.openshift.io/kms-plugin-image: "true"
  data:
    image: docker.io/hashicorp/vault-kube-kms@sha256:<digest>
```

If the ConfigMap is deleted or its `image` value or label is modified, the reconciler detects the drift and restores the desired state immediately.

## How it works

1. **Install via OLM** — the operator is installed from OperatorHub into the `openshift-kms-plugin-provider` namespace (created by CVO).
2. **On startup** — the `Start()` runnable creates the ConfigMap if it does not exist.
3. **On change** — the controller watches all ConfigMaps in its namespace; any change to `ibm-kms-vault-plugin-provider` triggers reconciliation and the ConfigMap is restored.
4. **Red Hat reads it** — `cluster-kube-apiserver-operator` watches the ConfigMap via the `config.openshift.io/kms-plugin-image` label selector and uses the `image` value to deploy the KMS plugin DaemonSet.

## Plugin image

The image reference is baked into the operator binary at build time via `-ldflags`:

```makefile
LD_FLAGS := -X github.com/hashicorp/vault-kms-plugin-openshift-provider/internal/version.PluginImage=<digest-pinned-image>
```

This ensures every operator release is tied to a specific, immutable plugin image digest.

## Development

### Prerequisites

- Go 1.26+
- [operator-sdk](https://sdk.operatorframework.io/) v1.x
- [Kind](https://kind.sigs.k8s.io/) (for e2e tests)
- Docker (for building images)

### Running unit tests

```bash
make test
```

### Running e2e tests (requires Kind)

```bash
make test-e2e
```

### Building the operator image

```bash
make docker-build IMG=<your-image>
```

### Generating bundle manifests

```bash
make bundle IMG=<your-image>
```

## License

Copyright (c) HashiCorp, Inc.  
Licensed under the [Business Source License 1.1](LICENSE).
