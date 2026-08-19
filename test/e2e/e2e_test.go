//go:build e2e
// +build e2e

package e2e

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kevinrizza/vault-kms-plugin-openshift-provider/test/utils"
)

var _ = Describe("Vault KMS Plugin OpenShift Provider", Ordered, func() {
	SetDefaultEventuallyTimeout(2 * time.Minute)
	SetDefaultEventuallyPollingInterval(time.Second)

	Context("OLM installation", func() {
		It("should have a CSV in Succeeded phase", func() {
			cmd := exec.Command("kubectl", "get", "csv",
				"-n", namespace,
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal("Succeeded"))
		})

		It("should have the controller manager pod running", func() {
			verifyPodRunning := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "pods",
					"-l", "control-plane=controller-manager",
					"-n", namespace,
					"-o", "jsonpath={.items[0].status.phase}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Running"))
			}
			Eventually(verifyPodRunning).Should(Succeed())
		})
	})

	Context("ConfigMap reconciliation", func() {
		It("should create the ConfigMap with the correct data", func() {
			verifyConfigMap := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "configmap",
					"ibm-kms-vault-plugin-provider",
					"-n", namespace,
					"-o", "jsonpath={.data.image}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("quay.io/kevinrizza/test-vault-plugin-image:latest"))
			}
			Eventually(verifyConfigMap).Should(Succeed())
		})

		It("should restore the ConfigMap after deletion", func() {
			By("deleting the ConfigMap")
			cmd := exec.Command("kubectl", "delete", "configmap",
				"ibm-kms-vault-plugin-provider",
				"-n", namespace)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("verifying the ConfigMap is recreated")
			verifyRecreated := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "configmap",
					"ibm-kms-vault-plugin-provider",
					"-n", namespace,
					"-o", "jsonpath={.data.image}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("quay.io/kevinrizza/test-vault-plugin-image:latest"))
			}
			Eventually(verifyRecreated).Should(Succeed())
		})

		It("should restore the ConfigMap data after modification", func() {
			By("modifying the ConfigMap data")
			cmd := exec.Command("kubectl", "patch", "configmap",
				"ibm-kms-vault-plugin-provider",
				"-n", namespace,
				"--type", "merge",
				"-p", `{"data":{"image":"tampered-value"}}`)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("verifying the ConfigMap data is restored")
			verifyRestored := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "configmap",
					"ibm-kms-vault-plugin-provider",
					"-n", namespace,
					"-o", "jsonpath={.data.image}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("quay.io/kevinrizza/test-vault-plugin-image:latest"))
			}
			Eventually(verifyRestored).Should(Succeed())
		})
	})
})
