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

	//kapps "k8s.io/api/apps/v1"
	kbatch "k8s.io/api/batch/v1"
	kcore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	taskv1alpha1 "github.com/AustrianDataLAB/execDAT-operator/api/v1alpha1"
)

// BuildReconciler reconciles a Build object
type BuildReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func create_init_sh(spec taskv1alpha1.BuildSpec, IMAGE_NAME string) string {
	nl := " \n"

	// Declare variable
	init_sh := fmt.Sprintf(`
		#!/bin/bash
		export BUILD="$(buildah from %s)"
		export ENTRYPOINT="%s"
		export IMAGE_NAME="%s"
		export USER="%s"
		export PASS="%s"
		export REGISTRY_URL="%s"
	`,
		spec.BaseImage,
		spec.SourceCode.Entrypoint,
		IMAGE_NAME,
		"test",
		"testPW",
		"localhost:5000")

	init_sh += nl

	// loop over all dependencies and fill init.sh with the commands to install them
	for i := 0; i < len(spec.SourceCode.Dependencies.OS); i++ {
		if i == 0 {
			init_sh += "buildah run $BUILD apt update"
			init_sh += nl
		}
		init_sh += "buildah run $BUILD apt install " + spec.SourceCode.Dependencies.OS[i].Name
		if spec.SourceCode.Dependencies.OS[i].Version != "" {
			init_sh += "=" + spec.SourceCode.Dependencies.OS[i].Version
		}
		init_sh += nl
	}

	init_sh += `
		buildah config --entrypoint "$ENTRYPOINT" $BUILD
		buildah commit $BUILD $IMAGE_NAME
		buildah push --tls-verify=false --creds $USER:$PASS $IMAGE_NAME docker://$REGISTRY_URL/$USER/$IMAGE_NAME
	`

	return init_sh
}

//+kubebuilder:rbac:groups=task.execd.at,resources=builds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=task.execd.at,resources=builds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=task.execd.at,resources=builds/finalizers,verbs=update

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

//+kubebuilder:rbac:resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:resources=configmaps/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Build object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BuildReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	log := logger.WithValues("taskv1alpha1.Build", req.NamespacedName)

	// TODO(user): your logic here

	build := &taskv1alpha1.Build{}
	if err := r.Get(ctx, req.NamespacedName, build); err != nil {
		log.Error(err, "unable to fetch build")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var resourceName string = build.Name
	var resourceNamespace string = build.Namespace

	cm := &kcore.ConfigMap{}
	cm.Name = resourceName
	cm.Namespace = resourceNamespace
	cm.Data = map[string]string{
		"init.sh": create_init_sh(build.Spec, resourceName+"-"+string(build.UID)),
	}

	if err := controllerutil.SetControllerReference(build, cm, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundConfigMap := &kcore.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("Creating CM", "ConfigMap", cm.Name)
		err = r.Create(ctx, cm)
	} else if err == nil {
		log.V(1).Info("ConfigMap already created", "ConfiMap", cm.Name)
	}

	podSpec := kcore.PodSpec{}
	podSpec.RestartPolicy = kcore.RestartPolicyNever
	podSpec.Containers = []kcore.Container{
		{
			Name:    "buildah",
			Image:   "ghcr.io/austriandatalab/execdat-operator-buildah:main",
			Command: []string{"sh"},
			Args:    []string{"/mnt/init.sh"},
			Env: []kcore.EnvVar{
				{Name: "BASE_IMAGE", Value: build.Spec.BaseImage},
				{Name: "GIT_REPO", Value: build.Spec.SourceCode.URL},
				{Name: "GIT_BRANCH", Value: build.Spec.SourceCode.Branch},
				{Name: "BUILD_CMD", Value: build.Spec.SourceCode.BuildCMD},
			},
		},
	}

	job := &kbatch.Job{}
	job.GenerateName = resourceName + "-"
	job.ObjectMeta.Namespace = resourceNamespace
	var ttl int32 = 60
	job.Spec.TTLSecondsAfterFinished = &ttl
	job.Spec.Template.Spec = podSpec

	if err := controllerutil.SetControllerReference(build, job, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundJob := &kbatch.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: resourceNamespace}, foundJob)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("Creating Job", "job", job.Name)
		err = r.Create(ctx, job)
	} else if err == nil {
		log.V(1).Info("Job already created", "job", job.Name)
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuildReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.Build{}).
		Complete(r)
}
