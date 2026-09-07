#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -e -o pipefail

# Append the minimum supported OpenShift version annotation to the OLM bundle
# metadata/annotations.yaml file.
#
# Red Hat's certification pipeline reads this annotation to gate the operator
# against OCP versions that pre-date the KMSEncryptionProvider feature gate.
#
# The default below (v4.18) reflects the first OCP release expected to ship
# the cluster-kube-apiserver-operator changes that read the plugin ConfigMap.
# Override by setting OPENSHIFT_MINIMUM_VERSION=vX.Y before calling this script.

ANNOTATIONS="bundle/metadata/annotations.yaml"

OPENSHIFT_MINIMUM_VERSION="${OPENSHIFT_MINIMUM_VERSION:-v4.18}"

# Idempotent — do not append twice.
if grep -q "com.redhat.openshift.versions" "${ANNOTATIONS}" 2>/dev/null; then
  echo "com.redhat.openshift.versions already set in ${ANNOTATIONS}"
  exit 0
fi

{
  echo ""
  echo "  # Minimum OpenShift version — set by hack/set_openshift_minimum_version.sh"
  echo "  com.redhat.openshift.versions: \"${OPENSHIFT_MINIMUM_VERSION}\""
} >> "${ANNOTATIONS}"

echo "Set com.redhat.openshift.versions = ${OPENSHIFT_MINIMUM_VERSION} in ${ANNOTATIONS}"
