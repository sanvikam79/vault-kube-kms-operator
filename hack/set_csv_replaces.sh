#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -e -o pipefail

# Set the `spec.replaces` field in the ClusterServiceVersion (CSV) to the
# previous git tag, so that the OLM upgrade graph is built correctly and
# previous versions show up in OperatorHub.
#
# https://olm.operatorframework.io/docs/concepts/olm-architecture/operator-catalog/creating-an-update-graph/
#
# Override the resolved previous tag by setting PREVIOUS_VERSION=vX.Y.Z before
# calling this script (used by the CRT release pipeline when the git history is
# not fully available).

CSV="bundle/manifests/vault-kube-kms-operator.clusterserviceversion.yaml"

# Already set — nothing to do (idempotent re-run safety).
CHECK="$(grep 'replaces:' "${CSV}" || true)"
if [ -n "${CHECK}" ]; then
  echo "replaces already set in ${CSV}:"
  echo "${CHECK}"
  exit 0
fi

if [ -z "${PREVIOUS_VERSION}" ]; then
  PREVIOUS_GIT_TAG="$(git describe --abbrev=0 --tags "$(git rev-list --tags --skip=1 --max-count=1)" 2>/dev/null || true)"
  PREVIOUS_VERSION="${PREVIOUS_GIT_TAG}"
fi

if [ -z "${PREVIOUS_VERSION}" ]; then
  echo "No previous git tag found — skipping spec.replaces (first release, no upgrade graph needed)"
  exit 0
fi

echo "  replaces: vault-kube-kms-operator.${PREVIOUS_VERSION}" >> "${CSV}"
echo "Set spec.replaces = vault-kube-kms-operator.${PREVIOUS_VERSION} in ${CSV}"
