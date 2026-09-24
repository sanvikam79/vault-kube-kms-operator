// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package controller

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	kmsv1alpha1 "github.com/hashicorp/vault-kube-kms-operator/api/v1alpha1"
	"github.com/hashicorp/vault-kube-kms-operator/internal/version"
)

var _ = Describe("VaultKMSConfig Controller", func() {
	const (
		timeout          = 10 * time.Second
		interval         = 250 * time.Millisecond
		testConfigName   = "test-config"
		testVaultAddress = "https://vault.example.com:8200"
		testVaultKeyPath = "transit/keys/my-key"
		testSecretName   = "vault-approle-creds"
	)

	AfterEach(func() {
		config := &kmsv1alpha1.VaultKMSConfig{}
		err := k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, config)
		if err == nil {
			Expect(k8sClient.Delete(ctx, config)).To(Succeed())
		}
	})

	Context("when a VaultKMSConfig is created", func() {
		It("should populate the status with the plugin image", func() {
			config := &kmsv1alpha1.VaultKMSConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testConfigName,
				},
				Spec: kmsv1alpha1.VaultKMSConfigSpec{
					VaultAddress: testVaultAddress,
					VaultKeyPath: testVaultKeyPath,
					Authentication: kmsv1alpha1.VaultAuthentication{
						Type: kmsv1alpha1.VaultAuthenticationTypeAppRole,
						AppRole: kmsv1alpha1.VaultAppRoleAuthentication{
							Secret: kmsv1alpha1.VaultSecretReference{
								Name: testSecretName,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, config)).To(Succeed())

			Eventually(func(g Gomega) {
				var fetched kmsv1alpha1.VaultKMSConfig
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, &fetched)).To(Succeed())
				g.Expect(fetched.Status.KMSPluginImage).To(Equal(version.PluginImage))
			}, timeout, interval).Should(Succeed())
		})

		It("should set the plugin image regardless of optional spec fields", func() {
			config := &kmsv1alpha1.VaultKMSConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testConfigName,
				},
				Spec: kmsv1alpha1.VaultKMSConfigSpec{
					VaultAddress:       testVaultAddress,
					VaultKeyPath:       testVaultKeyPath,
					VaultNamespace:     "admin/team-a",
					VaultAuthNamespace: "admin/auth",
					TLS: kmsv1alpha1.VaultTLSConfig{
						CABundle: kmsv1alpha1.VaultConfigMapReference{
							Name: "vault-ca-bundle",
						},
						ServerName: "vault.internal.example.com",
					},
					Authentication: kmsv1alpha1.VaultAuthentication{
						Type: kmsv1alpha1.VaultAuthenticationTypeAppRole,
						AppRole: kmsv1alpha1.VaultAppRoleAuthentication{
							Secret: kmsv1alpha1.VaultSecretReference{
								Name: testSecretName,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, config)).To(Succeed())

			Eventually(func(g Gomega) {
				var fetched kmsv1alpha1.VaultKMSConfig
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, &fetched)).To(Succeed())
				g.Expect(fetched.Status.KMSPluginImage).To(Equal(version.PluginImage))
			}, timeout, interval).Should(Succeed())
		})

		It("should preserve the plugin image after spec changes", func() {
			config := &kmsv1alpha1.VaultKMSConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testConfigName,
				},
				Spec: kmsv1alpha1.VaultKMSConfigSpec{
					VaultAddress: testVaultAddress,
					VaultKeyPath: testVaultKeyPath,
					Authentication: kmsv1alpha1.VaultAuthentication{
						Type: kmsv1alpha1.VaultAuthenticationTypeAppRole,
						AppRole: kmsv1alpha1.VaultAppRoleAuthentication{
							Secret: kmsv1alpha1.VaultSecretReference{
								Name: testSecretName,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, config)).To(Succeed())

			Eventually(func(g Gomega) {
				var fetched kmsv1alpha1.VaultKMSConfig
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, &fetched)).To(Succeed())
				g.Expect(fetched.Status.KMSPluginImage).To(Equal(version.PluginImage))
			}, timeout, interval).Should(Succeed())

			By("updating the spec")
			var current kmsv1alpha1.VaultKMSConfig
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, &current)).To(Succeed())
			current.Spec.VaultAddress = "https://vault-new.example.com:8200"
			current.Spec.VaultKeyPath = "transit/keys/new-key"
			Expect(k8sClient.Update(ctx, &current)).To(Succeed())

			Eventually(func(g Gomega) {
				var fetched kmsv1alpha1.VaultKMSConfig
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testConfigName}, &fetched)).To(Succeed())
				g.Expect(fetched.Status.KMSPluginImage).To(Equal(version.PluginImage))
			}, timeout, interval).Should(Succeed())
		})
	})
})
