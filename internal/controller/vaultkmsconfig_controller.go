// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/equality"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kmsv1alpha1 "github.com/hashicorp/vault-kube-kms-operator/api/v1alpha1"
	"github.com/hashicorp/vault-kube-kms-operator/internal/version"
)

// +kubebuilder:rbac:groups=kms.openshift.io,resources=vaultkmsconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kms.openshift.io,resources=vaultkmsconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kms.openshift.io,resources=vaultkmsconfigs/finalizers,verbs=update

type VaultKMSConfigReconciler struct {
	client.Client
}

func (r *VaultKMSConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var config kmsv1alpha1.VaultKMSConfig
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	desired := kmsv1alpha1.VaultKMSConfigStatus{
		// version.PluginImage is set at compile time via -ldflags in release builds.
		// In local / Kind e2e builds it uses the default defined in internal/version/version.go.
		KMSPluginImage: version.PluginImage,
	}

	if !equality.Semantic.DeepEqual(config.Status, desired) {
		logger.Info("updating VaultKMSConfig status")
		config.Status = desired
		if err := r.Status().Update(ctx, &config); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *VaultKMSConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kmsv1alpha1.VaultKMSConfig{}).
		Complete(r)
}
