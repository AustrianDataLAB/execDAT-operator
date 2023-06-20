package v1alpha1

import (
	kcore "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

type PodSpecData struct {
	INIT_SH    string
	Dockerfile string
	ImageName  string
}

// SetPodSpec sets the pod spec for the build
func (build *Build) SetPodSpec(podSpec *kcore.PodSpec, podSpecData PodSpecData) error {

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
			Image:   "quay.io/buildah/stable",
			Command: []string{"/bin/bash", "-c", "--"},
			Args:    []string{"trap : TERM INT; echo \"$INIT_SH\" | bash"},
			Env: []kcore.EnvVar{
				{Name: "INIT_SH", Value: podSpecData.INIT_SH},
				{Name: "DOCKERFILE", Value: podSpecData.Dockerfile},
				{Name: "BASE_IMAGE", Value: build.Spec.BaseImage},
				{Name: "IMAGE_NAME", Value: podSpecData.ImageName},
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
				{Name: "REGISTRY_USER", ValueFrom: &kcore.EnvVarSource{
					SecretKeyRef: &kcore.SecretKeySelector{
						LocalObjectReference: kcore.LocalObjectReference{
							Name: "buildah-registry-credentials",
						},
						Key: "REGISTRY_USER",
					},
				}},
				{Name: "REGISTRY_KEY", ValueFrom: &kcore.EnvVarSource{
					SecretKeyRef: &kcore.SecretKeySelector{
						LocalObjectReference: kcore.LocalObjectReference{
							Name: "buildah-registry-credentials",
						},
						Key: "REGISTRY_KEY",
					},
				}},
				{Name: "REGISTRY_BASE", ValueFrom: &kcore.EnvVarSource{
					SecretKeyRef: &kcore.SecretKeySelector{
						LocalObjectReference: kcore.LocalObjectReference{
							Name: "buildah-registry-credentials",
						},
						Key: "REGISTRY_BASE",
					},
				}},
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
