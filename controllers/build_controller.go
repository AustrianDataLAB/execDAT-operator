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
	"fmt"
	"strings"
	"time"

	//kapps "k8s.io/api/apps/v1"

	kbatch "k8s.io/api/batch/v1"
	kcore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	pointer "k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	taskv1alpha1 "github.com/AustrianDataLAB/execDAT-operator/api/v1alpha1"
	lib "github.com/AustrianDataLAB/execDAT-operator/lib"
)

// BuildReconciler reconciles a Build object
type BuildReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=task.execd.at,resources=builds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=task.execd.at,resources=builds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=task.execd.at,resources=builds/finalizers,verbs=update

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// +kubebuilder:rbac:resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:resources=configmaps/status,verbs=get
func (r *BuildReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	log := logger.WithValues("taskv1alpha1.Build", req.NamespacedName)

	build := &taskv1alpha1.Build{}
	if err := r.Get(ctx, req.NamespacedName, build); err != nil {
		log.V(1).Info("unable to fetch build", "build", req.NamespacedName)
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if build.Status.CurrentPhase != nil && *build.Status.CurrentPhase == taskv1alpha1.CurrentPhaseBuildComplete {
		log.V(1).Info("Build already completed", "build", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	var resourceName string = build.Name
	var resourceNamespace string = build.Namespace

	scriptTemplates := []string{"./templates/init_build.sh.tmpl"}
	init_sh, err := lib.CreateTemplate(scriptTemplates, build.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	dockerTemplates := []string{"./templates/Dockerfile.tmpl"}
	dockerfile, err := lib.CreateTemplate(dockerTemplates, build.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	newBuildPodSpecData := taskv1alpha1.BuildPodSpecData{
		INIT_SH:    init_sh,
		Dockerfile: dockerfile,
		ImageName:  build.ObjectMeta.Labels["runRef"],
		ImageTag:   strings.Split(build.ObjectMeta.Name, build.ObjectMeta.GenerateName)[1],
	}

	podSpec := &kcore.PodSpec{}
	if err := build.SetPodSpec(podSpec, newBuildPodSpecData); err != nil {
		return ctrl.Result{}, err
	}

	job := &kbatch.Job{}
	job.Name = resourceName
	job.ObjectMeta.Namespace = resourceNamespace
	job.Spec.TTLSecondsAfterFinished = pointer.Int32(60)
	job.Spec.Template.Spec = *podSpec

	if err := controllerutil.SetControllerReference(build, job, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundJob := &kbatch.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, foundJob)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("Creating Job", "job", job.Name)
		//TODO create the Job in a separate namespace
		err = r.Create(ctx, job)
	} else if err == nil {

		log.V(1).Info("Job already created", "job", job.Name)
		if build.Status.CurrentPhase == nil {
			build.Status.CurrentPhase = taskv1alpha1.CurrentPhaseBuilding.Ptr()
			if err := r.Status().Update(ctx, build); err != nil {
				return ctrl.Result{RequeueAfter: time.Second * 10}, fmt.Errorf("failed to update build status: %w", err)
			}
		}

		// Check if the Job has completed
		jobComplete, err := CheckJobCompletion(job, r.Client)
		if err != nil {
			log.Error(err, "Failed to check Job completion")
			return ctrl.Result{}, err
		}
		if jobComplete {
			log.V(1).Info("Job has completed", "job", job.Name)
			build.Status.CurrentPhase = taskv1alpha1.CurrentPhaseBuildComplete.Ptr()
			if err := r.Status().Update(ctx, build); err != nil {
				return ctrl.Result{RequeueAfter: time.Second * 10}, fmt.Errorf("failed to update build status: %w", err)
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuildReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.Build{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
