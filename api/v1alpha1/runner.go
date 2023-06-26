package v1alpha1

import (
	kcore "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

type RunPodSpecData struct {
	INIT_SH        string
	ImageName      string
	ImageTag       string
	InputDataPath  string
	OutputDataPath string
}

// SetPodSpec sets the pod spec for the run
func (run *Run) SetPodSpec(podSpec *kcore.PodSpec, runPodSpecData RunPodSpecData) error {

	podSpec.RestartPolicy = kcore.RestartPolicyNever

	podSpec.Volumes = []kcore.Volume{
		{
			Name: "input",
			VolumeSource: kcore.VolumeSource{
				EmptyDir: &kcore.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "output",
			VolumeSource: kcore.VolumeSource{
				EmptyDir: &kcore.EmptyDirVolumeSource{},
			},
		},
	}

	podSpec.Containers = []kcore.Container{
		{
			Name:            "buildah",
			Image:           "harbor.caas-0013.dev.austrianopencloudcommunity.org/execdev/" + runPodSpecData.ImageName + ":" + runPodSpecData.ImageTag,
			ImagePullPolicy: kcore.PullIfNotPresent,
			Command:         []string{"/bin/bash", "-c", "--"},
			Args:            []string{"trap : TERM INT; echo \"$INIT_SH\" | bash"},
			TTY:             true,
			Stdin:           true,
			Env: []kcore.EnvVar{
				{Name: "INIT_SH", Value: runPodSpecData.INIT_SH},
				{Name: "MINIO_ENDPOINT", Value: "http://minio.single-minio.svc.cluster.local:9000"},
				{Name: "MINIO_ACCESS_KEY", Value: "cache-user-1"},
				{Name: "MINIO_SECRET_KEY", Value: "CACHE_USER_PASS_XYZ_123"},
				{Name: "BUCKET_NAME", Value: "cache-bucket-1"},
				{Name: "HOME", Value: "/tmp"},
			},
			SecurityContext: &kcore.SecurityContext{
				RunAsUser:  pointer.Int64(1000),
				RunAsGroup: pointer.Int64(1000),
			},
			VolumeMounts: []kcore.VolumeMount{
				{Name: "input", MountPath: runPodSpecData.InputDataPath},
				{Name: "output", MountPath: runPodSpecData.OutputDataPath},
			},
		},
	}
	podSpec.SecurityContext = &kcore.PodSecurityContext{
		RunAsUser:  pointer.Int64(1000),
		RunAsGroup: pointer.Int64(1000),
	}
	return nil
}
