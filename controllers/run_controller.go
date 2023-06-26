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
	// Define labels to filter Jobs
	labels := map[string]string{
		"runRef":        run.Name,
		"runGeneration": fmt.Sprint(run.ObjectMeta.Generation),
	}
	build.ObjectMeta.Labels = labels
	build.Spec = run.Spec.Build

	if err := controllerutil.SetControllerReference(run, build, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundBuildList := &taskv1alpha1.BuildList{}
	err := r.List(ctx,
		foundBuildList,
		client.InNamespace(resourceNamespace),
		client.MatchingLabels(labels),
	)
	if err == nil && len(foundBuildList.Items) <= 0 {
		log.V(1).Info("Creating Build", "build", build.Name)
		if err := r.Create(ctx, build); err != nil {
			log.Error(err, "Failed to create new Build", "build", build.Name)
			return ctrl.Result{}, err
		}
		run.Status.CurrentPhase = taskv1alpha1.CurrentPhaseBuilding.Ptr()
		if err := r.Status().Update(ctx, run); err != nil {
			log.Error(err, "Failed to update Run status")
			return ctrl.Result{}, err
		}
	} else if len(foundBuildList.Items) > 0 {

		log.V(1).Info("Build already created", "build", build.Name)

		currentBuild := foundBuildList.Items[0]

		if currentBuild.Status.CurrentPhase == nil {
			return ctrl.Result{}, nil
		}

		switch *currentBuild.Status.CurrentPhase {
		case taskv1alpha1.CurrentPhaseBuildComplete:
			log.V(1).Info("Build completed", "run", build.Name)
			//TODO: create RUN
			scriptTemplates := []string{"./templates/init_run.sh.tmpl"}
			init_sh, err := lib.CreateTemplate(scriptTemplates, run.Spec)
			if err != nil {
				return ctrl.Result{}, err
			}

			newRunPodSpecData := taskv1alpha1.RunPodSpecData{
				INIT_SH:        init_sh,
				ImageName:      currentBuild.ObjectMeta.Labels["runRef"],
				ImageTag:       strings.Split(currentBuild.ObjectMeta.Name, currentBuild.ObjectMeta.GenerateName)[1],
				InputDataPath:  run.Spec.InputData.DataPath,
				OutputDataPath: run.Spec.OutputData.DataPath,
			}

			podSpec := &kcore.PodSpec{}
			if err := run.SetPodSpec(podSpec, newRunPodSpecData); err != nil {
				return ctrl.Result{}, err
			}

			job := &kbatch.Job{}
			job.Name = resourceName
			job.ObjectMeta.Namespace = resourceNamespace
			job.Spec.TTLSecondsAfterFinished = pointer.Int32(60)
			job.Spec.Template.Spec = *podSpec

			if err := controllerutil.SetControllerReference(run, job, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			foundJob := &kbatch.Job{}
			err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, foundJob)
			if err != nil && errors.IsNotFound(err) && *run.Status.CurrentPhase != taskv1alpha1.CurrentPhaseRunCompleted {

				log.V(1).Info("Creating Run Job", "job", job.Name)
				//TODO create the Job in a separate namespace
				if err := r.Create(ctx, job); err != nil {
					log.Error(err, "Failed to create new Run Job", "job", job.Name)
					return ctrl.Result{}, err
				}

				run.Status.CurrentPhase = taskv1alpha1.CurrentPhaseRunning.Ptr()
				if err := r.Status().Update(ctx, run); err != nil {
					log.Error(err, "Failed to update Run status")
					return ctrl.Result{}, err
				}

			} else if err == nil {
				log.V(1).Info("Run Job already created", "job", job.Name)

				// Check if the Job has completed
				jobComplete, err := CheckJobCompletion(job, r.Client)
				if err != nil {
					log.Error(err, "Failed to check Job completion")
					return ctrl.Result{}, err
				}
				if jobComplete {
					log.V(1).Info("Job has completed", "job", job.Name)
					run.Status.CurrentPhase = taskv1alpha1.CurrentPhaseRunCompleted.Ptr()
					if err := r.Status().Update(ctx, run); err != nil {
						return ctrl.Result{RequeueAfter: time.Second * 10}, fmt.Errorf("failed to update run status: %w", err)
					}
					return ctrl.Result{}, nil
				}
				return ctrl.Result{RequeueAfter: time.Second * 10}, nil
			}

		case taskv1alpha1.CurrentPhaseFailed:
			log.V(1).Info("Build failed", "build", build.Name)
			run.Status.CurrentPhase = taskv1alpha1.CurrentPhaseFailed.Ptr()
			if err := r.Status().Update(ctx, run); err != nil {
				log.Error(err, "Failed to update Run status")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.Run{}).
		Owns(&taskv1alpha1.Build{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
