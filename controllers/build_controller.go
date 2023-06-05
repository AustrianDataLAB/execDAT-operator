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
	pointer "k8s.io/utils/pointer"
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

	var resourceName string = build.Name
	var resourceNamespace string = build.Namespace

	scriptTemplates := []string{"./templates/init.sh.tmpl"}
	templateData := TemplateData{
		BaseImage: build.Spec.BaseImage,
		// GitRepo:   build.Spec.SourceCode.URL,
		// GitBranch: build.Spec.SourceCode.Branch,
		// BuildCmd:  build.Spec.SourceCode.BuildCMD,
	}
	init_sh, err := GenerateScript(scriptTemplates, templateData)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error generating script: %v", err)
	}

	podSpec := kcore.PodSpec{}
	podSpec.RestartPolicy = kcore.RestartPolicyNever
	podSpec.Containers = []kcore.Container{
		{
			Name:    "buildah",
			Image:   "ghcr.io/austriandatalab/execdat-operator-buildah:main",
			Command: []string{"/bin/bash", "-c", "--"},
			Args:    []string{"trap : TERM INT; echo \"$INIT_SH\" | bash"},
			Env: []kcore.EnvVar{
				{Name: "INIT_SH", Value: init_sh},
				{Name: "BASE_IMAGE", Value: build.Spec.BaseImage},
				// {Name: "IMAGE_NAME", Value: build.Spec.ImageName},
				// {Name: "IMAGE_TAG", Value: build.Spec.ImageTag},
				// {Name: "IMAGE_REGISTRY", Value: build.Spec.ImageRegistry},
				// {Name: "IMAGE_REGISTRY_USER", Value: build.Spec.ImageRegistryUser},
				// {Name: "IMAGE_REGISTRY_PASSWORD", Value: build.Spec.ImageRegistryPassword},
				// {Name: "IMAGE_REGISTRY_INSECURE", Value: build.Spec.ImageRegistryInsecure},
				// {Name: "IMAGE_REGISTRY_VERIFY_TLS", Value: build.Spec.ImageRegistryVerifyTLS},
				{Name: "ENTRYPOINT", Value: build.Spec.SourceCode.Entrypoint},
				{Name: "GIT_REPO", Value: build.Spec.SourceCode.URL},
				{Name: "GIT_BRANCH", Value: build.Spec.SourceCode.Branch},
				{Name: "BUILD_CMD", Value: build.Spec.SourceCode.BuildCMD},
			},
			SecurityContext: &kcore.SecurityContext{
				Capabilities: &kcore.Capabilities{
					Add: []kcore.Capability{
						"SETUID",
						"SETGID",
					},
				},
			},
		},
	}
	podSpec.SecurityContext = &kcore.PodSecurityContext{
		RunAsUser:  pointer.Int64(1000),
		RunAsGroup: pointer.Int64(1000),
	}

	job := &kbatch.Job{}
	job.GenerateName = resourceName + "-"
	job.ObjectMeta.Namespace = resourceNamespace
	job.Spec.TTLSecondsAfterFinished = pointer.Int32(60)
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
