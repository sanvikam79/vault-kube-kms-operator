# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

schema = 1
artifacts {
  container = [
    "vault-kms-plugin-openshift-provider_release-ubi_linux_amd64_${version}_${commit_sha}.docker.redhat.tar",
    "vault-kms-plugin-openshift-provider_release-ubi_linux_arm64_${version}_${commit_sha}.docker.redhat.tar",
    "vault-kms-plugin-openshift-provider_release-ubi_linux_amd64_${version}_${commit_sha}.docker.tar",
    "vault-kms-plugin-openshift-provider_release-ubi_linux_arm64_${version}_${commit_sha}.docker.tar",
  ]
}
