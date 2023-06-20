package v1alpha1

import (
	"fmt"

	lib "github.com/AustrianDataLAB/execDAT-operator/lib"
	kcore "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func (build *Build) SetPodSpec(podSpec *kcore.PodSpec) error {

	scriptTemplates := []string{"./templates/init.sh.tmpl"}
	templateData := lib.TemplateData{
		BaseImage: build.Spec.BaseImage,
		// GitRepo:   build.Spec.SourceCode.URL,
		// GitBranch: build.Spec.SourceCode.Branch,
		// BuildCmd:  build.Spec.SourceCode.BuildCMD,
	}
	init_sh, err := lib.GenerateScript(scriptTemplates, templateData)
	if err != nil {
		return fmt.Errorf("error generating script: %v", err)
	}

	podSpec.RestartPolicy = kcore.RestartPolicyNever

	podSpec.Volumes = []kcore.Volume{
		{
			Name: "custom-uidmap",
			VolumeSource: kcore.VolumeSource{
				ConfigMap: &kcore.ConfigMapVolumeSource{
					LocalObjectReference: kcore.LocalObjectReference{Name: "custom-uidmap"},
				},
			}},
	}
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
				{Name: "STORAGE_DRIVER", Value: "vfs"},
				{Name: "BUILDAH_FORMAT", Value: "docker"},
				{Name: "BUILDAH_ISOLATION", Value: "chroot"},
			},
			SecurityContext: &kcore.SecurityContext{
				Capabilities: &kcore.Capabilities{
					Add: []kcore.Capability{
						"SETUID",
						"SETGID",
					},
				},
			},
			VolumeMounts: []kcore.VolumeMount{
				{Name: "custom-uidmap", MountPath: "/etc/subuid", SubPath: "subuid"},
				{Name: "custom-uidmap", MountPath: "/etc/subgid", SubPath: "subgid"},
			},
		},
	}
	podSpec.SecurityContext = &kcore.PodSecurityContext{
		RunAsUser:  pointer.Int64(1000),
		RunAsGroup: pointer.Int64(1000),
	}
	return nil
}
