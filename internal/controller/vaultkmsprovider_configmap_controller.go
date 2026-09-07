// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package controller

import (
	"context"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/hashicorp/vault-kms-plugin-openshift-provider/internal/version"
)

const (
	ConfigMapName = "ibm-kms-vault-plugin-provider"

	// LabelKMSPluginImage is required by Red Hat's cluster-kube-apiserver-operator
	// to locate the ConfigMap that provides the plugin image reference.
	LabelKMSPluginImage = "config.openshift.io/kms-plugin-image"
)

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch

type VaultKMSProviderConfigMapReconciler struct {
	client.Client
	Namespace string
}

func (r *VaultKMSProviderConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if req.Name != ConfigMapName || req.Namespace != r.Namespace {
		return ctrl.Result{}, nil
	}

	desired := r.desiredConfigMap()

	var existing corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: ConfigMapName, Namespace: r.Namespace}, &existing)
	if errors.IsNotFound(err) {
		logger.Info("creating configmap", "name", ConfigMapName)
		return ctrl.Result{}, r.Create(ctx, desired)
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	// Only update when data or labels actually differ — avoids noisy no-op updates.
	if reflect.DeepEqual(existing.Data, desired.Data) &&
		existing.Labels[LabelKMSPluginImage] == desired.Labels[LabelKMSPluginImage] {
		return ctrl.Result{}, nil
	}

	existing.Data = desired.Data
	existing.Labels = desired.Labels
	logger.Info("updating configmap", "name", ConfigMapName)
	return ctrl.Result{}, r.Update(ctx, &existing)
}

func (r *VaultKMSProviderConfigMapReconciler) desiredConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: r.Namespace,
			// Required label so Red Hat's cluster-kube-apiserver-operator
			// can discover this ConfigMap via a label selector.
			Labels: map[string]string{
				LabelKMSPluginImage: "true",
			},
		},
		Data: map[string]string{
			// version.PluginImage is set at compile time via -ldflags in release builds.
			// In local / Kind e2e builds it uses the default defined in internal/version/version.go.
			"image": version.PluginImage,
		},
	}
}

func (r *VaultKMSProviderConfigMapReconciler) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	desired := r.desiredConfigMap()

	var existing corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: ConfigMapName, Namespace: r.Namespace}, &existing)
	if errors.IsNotFound(err) {
		logger.Info("creating initial configmap", "name", ConfigMapName)
		return r.Create(ctx, desired)
	}
	if err != nil {
		return err
	}

	// Fix wrong content immediately on startup, before the informer cache warms up.
	// Without this, a ConfigMap tampered while the operator was dead would remain
	// wrong until the informer re-List completes and triggers Reconcile (10+ seconds).
	if !reflect.DeepEqual(existing.Data, desired.Data) ||
		existing.Labels[LabelKMSPluginImage] != desired.Labels[LabelKMSPluginImage] {
		existing.Data = desired.Data
		existing.Labels = desired.Labels
		logger.Info("correcting configmap on startup", "name", ConfigMapName)
		return r.Update(ctx, &existing)
	}
	return nil
}

func (r *VaultKMSProviderConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.Add(r); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}
