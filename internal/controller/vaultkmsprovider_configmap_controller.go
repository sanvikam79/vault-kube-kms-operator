package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ConfigMapName = "ibm-kms-vault-plugin-provider"
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

	existing.Data = desired.Data
	logger.Info("updating configmap", "name", ConfigMapName)
	return ctrl.Result{}, r.Update(ctx, &existing)
}

func (r *VaultKMSProviderConfigMapReconciler) desiredConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: r.Namespace,
		},
		Data: map[string]string{
			"image": "quay.io/kevinrizza/test-vault-plugin-image:latest",
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
	return err
}

func (r *VaultKMSProviderConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.Add(r); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}
