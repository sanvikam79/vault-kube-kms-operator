# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

schema = 1
artifacts {
  container = [
    "vault-kube-kms-operator_release-ubi_linux_amd64_${version}_${commit_sha}.docker.redhat.tar",
    "vault-kube-kms-operator_release-ubi_linux_arm64_${version}_${commit_sha}.docker.redhat.tar",
    "vault-kube-kms-operator_release-ubi_linux_amd64_${version}_${commit_sha}.docker.tar",
    "vault-kube-kms-operator_release-ubi_linux_arm64_${version}_${commit_sha}.docker.tar",
  ]
}
