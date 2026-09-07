#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -e -o pipefail

# Set the annotation `containerImage` in the ClusterServiceVersion to the
# operator image URL found in the deployment spec.
# This ensures OperatorHub displays the correct image reference in the UI.
#
# Requires: bin/yq (install via 'make yq')

HACK_DIR=$(dirname "$0")
CSV_FILE="${HACK_DIR}/../bundle/manifests/vault-kms-plugin-openshift-provider.clusterserviceversion.yaml"
YQ="${HACK_DIR}/../bin/yq"

if [ ! -f "${YQ}" ]; then
  echo "yq not found at ${YQ} — run 'make yq' first"
  exit 1
fi

IMAGE=$(cat "${CSV_FILE}" | "${YQ}" \
  '.spec.install.spec.deployments.[] | select(.name == "vault-kms-controller-manager") | .spec.template.spec.containers.[] | select(.name == "manager") | .image')

cat "${CSV_FILE}" | "${YQ}" ".metadata.annotations.containerImage |= (\"${IMAGE}\")" > /tmp/csv.yaml
mv /tmp/csv.yaml "${CSV_FILE}"

echo "Set containerImage annotation to: ${IMAGE}"
