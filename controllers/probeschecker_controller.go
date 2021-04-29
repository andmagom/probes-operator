/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"

	"context"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/example/probes-operator/api/v1alpha1"
)

// ProbesCheckerReconciler reconciles a ProbesChecker object
type ProbesCheckerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.vmware.com,resources=probescheckers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.vmware.com,resources=probescheckers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.vmware.com,resources=probescheckers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProbesChecker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ProbesCheckerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("probeschecker", req.NamespacedName)

	// your logic here
	probeschecker := &cachev1alpha1.ProbesChecker{}
	err := r.Get(ctx, req.NamespacedName, probeschecker)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("ProbesChecker resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ProbesChecker")
		return ctrl.Result{}, err
	}

	pods := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
	}
	err = r.List(ctx, pods, opts...)
	fmt.Printf("Pods Size: %d\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("PodName: %s :: ContainersSize: %d\n", pod.Name, len(pod.Spec.Containers))
		for _, container := range pod.Spec.Containers {
			fmt.Printf("ContainerName: %s :: ProbesEnabled: %v\n", container.Name, container.LivenessProbe)
		}
	}

	// return ctrl.Result{}, nil		->  return with error
	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProbesCheckerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.ProbesChecker{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
