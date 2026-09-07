// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package version

// PluginImage is the digest-pinned vault-kube-kms container image reference
// written into the ibm-kms-vault-plugin-provider ConfigMap.
//
// In release builds this variable is overridden at compile time via -ldflags:
//
//	-X github.com/hashicorp/vault-kms-plugin-openshift-provider/internal/version.PluginImage=<digest>
//
// The default below is used for local development and Kind e2e testing only.
var PluginImage = "docker.io/hashicorp/vault-kube-kms:0.1.0-beta"
