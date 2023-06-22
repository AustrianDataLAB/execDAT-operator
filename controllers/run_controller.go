/*
Copyright 2023 Thomas Weber.

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
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	taskv1alpha1 "github.com/AustrianDataLAB/execDAT-operator/api/v1alpha1"
)

// RunReconciler reconciles a Run object
type RunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=task.execd.at,resources=runs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=task.execd.at,resources=runs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=task.execd.at,resources=runs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Run object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	log := logger.WithValues("taskv1alpha1.Run", req.NamespacedName)

	run := &taskv1alpha1.Run{}
	if err := r.Get(ctx, req.NamespacedName, run); err != nil {
		log.V(1).Info("unable to fetch run", "run", req.NamespacedName)
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.V(1).Info("reconciling run", "run", run)
	var resourceName string = run.Name
	var resourceNamespace string = run.Namespace

	build := &taskv1alpha1.Build{}
	build.GenerateName = resourceName + "-"
	build.ObjectMeta.Namespace = resourceNamespace
	build.Spec = run.Spec.Build

	if err := controllerutil.SetControllerReference(run, build, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, build); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.Run{}).
		Complete(r)
}
