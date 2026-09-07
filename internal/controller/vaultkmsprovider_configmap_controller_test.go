// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/hashicorp/vault-kms-plugin-openshift-provider/internal/version"
)

const (
	testNamespace = "openshift-kms-plugin-provider"
	timeout       = 10 * time.Second
	interval      = 250 * time.Millisecond
)

var _ = Describe("VaultKMSProviderConfigMapReconciler", Ordered, func() {
	ctx := context.Background()

	BeforeAll(func() {
		// Create the namespace that the reconciler targets.
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testNamespace}}
		err := k8sClient.Create(ctx, ns)
		if err != nil {
			// Namespace may already exist from a previous test run in the same envtest session.
			Expect(err.Error()).To(ContainSubstring("already exists"))
		}
	})

	namespacedName := types.NamespacedName{Name: ConfigMapName, Namespace: testNamespace}

	It("should create the ConfigMap with the correct image value", func() {
		Eventually(func(g Gomega) {
			cm := &corev1.ConfigMap{}
			g.Expect(k8sClient.Get(ctx, namespacedName, cm)).To(Succeed())
			g.Expect(cm.Data["image"]).To(Equal(version.PluginImage))
		}, timeout, interval).Should(Succeed())
	})

	It("should create the ConfigMap with the required Red Hat label", func() {
		Eventually(func(g Gomega) {
			cm := &corev1.ConfigMap{}
			g.Expect(k8sClient.Get(ctx, namespacedName, cm)).To(Succeed())
			g.Expect(cm.Labels[LabelKMSPluginImage]).To(Equal("true"))
		}, timeout, interval).Should(Succeed())
	})

	It("should restore the ConfigMap after deletion", func() {
		// Delete the ConfigMap.
		cm := &corev1.ConfigMap{}
		Expect(k8sClient.Get(ctx, namespacedName, cm)).To(Succeed())
		Expect(k8sClient.Delete(ctx, cm)).To(Succeed())

		// Reconciler should recreate it.
		Eventually(func(g Gomega) {
			restored := &corev1.ConfigMap{}
			g.Expect(k8sClient.Get(ctx, namespacedName, restored)).To(Succeed())
			g.Expect(restored.Data["image"]).To(Equal(version.PluginImage))
		}, timeout, interval).Should(Succeed())
	})

	It("should restore tampered ConfigMap data", func() {
		// Patch the data to a wrong value.
		cm := &corev1.ConfigMap{}
		Expect(k8sClient.Get(ctx, namespacedName, cm)).To(Succeed())
		patch := cm.DeepCopy()
		patch.Data = map[string]string{"image": "tampered-value"}
		Expect(k8sClient.Update(ctx, patch)).To(Succeed())

		// Reconciler should correct it.
		Eventually(func(g Gomega) {
			updated := &corev1.ConfigMap{}
			g.Expect(k8sClient.Get(ctx, namespacedName, updated)).To(Succeed())
			g.Expect(updated.Data["image"]).To(Equal(version.PluginImage))
		}, timeout, interval).Should(Succeed())
	})

	It("should restore the required label after it is removed", func() {
		// Remove the label.
		cm := &corev1.ConfigMap{}
		Expect(k8sClient.Get(ctx, namespacedName, cm)).To(Succeed())
		patch := cm.DeepCopy()
		delete(patch.Labels, LabelKMSPluginImage)
		Expect(k8sClient.Update(ctx, patch)).To(Succeed())

		// Reconciler should restore it.
		Eventually(func(g Gomega) {
			restored := &corev1.ConfigMap{}
			g.Expect(k8sClient.Get(ctx, namespacedName, restored)).To(Succeed())
			g.Expect(restored.Labels[LabelKMSPluginImage]).To(Equal("true"))
		}, timeout, interval).Should(Succeed())
	})
})
