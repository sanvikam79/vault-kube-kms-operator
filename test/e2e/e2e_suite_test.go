//go:build e2e
// +build e2e

package e2e

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var namespace string

func init() {
	flag.StringVar(&namespace, "e2e.namespace", "openshift-kms-plugin-provider", "namespace where the operator is installed")
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	_, _ = fmt.Fprintf(GinkgoWriter, "Starting vault-kms-plugin-openshift-provider e2e test suite\n")
	RunSpecs(t, "e2e suite")
}
