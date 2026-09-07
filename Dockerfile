# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Build the manager binary
FROM golang:1.26 AS builder
ARG TARGETOS
ARG TARGETARCH

# LD_FLAGS is injected by the CRT release pipeline to bake version.PluginImage
# and other build metadata into the binary at compile time.
# e.g. -X github.com/hashicorp/vault-kms-plugin-openshift-provider/internal/version.PluginImage=docker.io/hashicorp/vault-kube-kms@sha256:...
ARG LD_FLAGS

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source
COPY . .

# Build — LD_FLAGS injected via --build-arg in release pipeline
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags "$LD_FLAGS" -a -o manager ./cmd/main.go

# ubi-minimal intermediate stage — used only to refresh CA certificates.
# ubi-micro has no package manager so cannot run microdnf itself;
# we copy the updated CA bundle from here into the final image.
# This matches the VSO Dockerfile pattern exactly.
FROM registry.access.redhat.com/ubi10/ubi-minimal:10.2 AS build-ubi
RUN microdnf --refresh --assumeyes upgrade ca-certificates

# Final runtime image — ubi-micro is the smallest Red Hat UBI variant.
# It contains no package manager, minimising the CVE surface area for
# Red Hat Partner Connect preflight certification.
# ubi-micro is FIPS-validated and accepted in the Red Hat Catalog.
FROM registry.access.redhat.com/ubi10/ubi-micro:10.2 AS release-ubi

ENV BIN_NAME=vault-kms-plugin-openshift-provider
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
ARG PRODUCT_NAME=$BIN_NAME
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

# OCI + Red Hat certification labels — required by preflight and HC IPS-002.
LABEL name="Vault KMS Plugin OpenShift Provider" \
      maintainer="Team Vault <vault@hashicorp.com>" \
      vendor="HashiCorp" \
      version=$PRODUCT_VERSION \
      release=$PRODUCT_VERSION \
      revision=$PRODUCT_REVISION \
      org.opencontainers.image.licenses="BUSL-1.1" \
      summary="OLM operator that publishes the Vault KMS plugin image reference for OpenShift encryption controllers." \
      description="OLM operator that publishes the Vault KMS plugin image reference for OpenShift encryption controllers."

WORKDIR /
COPY --from=builder /workspace/manager .

# Copy license for Red Hat certification.
COPY LICENSE /licenses/copyright.txt
# Copy license to conform to HC IPS-002
COPY LICENSE /usr/share/doc/$PRODUCT_NAME/LICENSE.txt

# Copy refreshed CA bundle from ubi-minimal intermediate stage.
# ubi-micro ships no microdnf so cannot update certs itself.
COPY --from=build-ubi /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/

# Run as non-root
USER 65532:65532

ENTRYPOINT ["/manager"]
