package deployment

import (
	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/go-kube/pkg/kube"

	"danicos.dev/daniel/go-kube/pkg/stack"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"
)

var Secret = struct {
	Name          string
	ConfigKey     string
	LitestreamKey string
}{
	Name:          "curious-ape-secret",
	ConfigKey:     "config.json",
	LitestreamKey: "litestream.yml",
}

var Namespace = kube.Namespace(config.KUBERNETES_NAME)

var meta kube.Metadata
var SRV core.Service
var PVC core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata(config.KUBERNETES_NAME, Namespace)
	PVC = meta.PVC()
	SRV = meta.Service(config.KUBERNETES_PORT)
}

func BaseStack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":  Namespace,
			"service":    SRV,
			"pvc":        PVC,
			"deployment": Deployment(),
		},
	)
	return kz.Stack("base")
}

func K3sStack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"ingress": Ingress(),
		},
	)
	return kz.Stack("k3s")
}

func SecretsStack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"secrets": ApeSecret,
		},
	)
	return kz.Stack("config")
}

func Deployment() apps.Deployment {
	dataVolume := kube.NewVolumeFrom(kube.VolumeSourcePVC, "data", PVC.Name)
	configVolume := kube.NewVolumeFromSecret("config", Secret.Name, []core.KeyToPath{{
		Key:  Secret.ConfigKey,
		Path: Secret.ConfigKey,
	}})
	litestreamVolume := kube.NewVolumeFromSecret("litestream-vol", Secret.Name, []core.KeyToPath{{
		Key:  Secret.LitestreamKey,
		Path: Secret.LitestreamKey,
	}})
	podSpec := core.PodSpec{
		InitContainers: []core.Container{{
			Name:  "restore-litestream",
			Image: config.LITESTREAM_IMAGE,
			Command: []string{
				"litestream",
				"restore",
				"-if-db-not-exists",
				"-if-replica-exists",
				"/db-data/ape.db",
			},
			VolumeMounts: []core.VolumeMount{
				{
					Name:      dataVolume.Name,
					MountPath: "/db-data",
				},
				{
					Name:      litestreamVolume.Name,
					MountPath: "/etc/litestream.yml",
					SubPath:   Secret.LitestreamKey,
				},
			},
		}},
		Containers: []core.Container{
			{
				Name:  config.KUBERNETES_NAME,
				Image: config.KUBERNETES_IMAGE,
				Ports: []core.ContainerPort{{ContainerPort: int32(config.KUBERNETES_PORT)}},
				Env:   []core.EnvVar{{Name: "APE_ENVIRONMENT", Value: "prod"}},
				VolumeMounts: []core.VolumeMount{
					{
						Name:      configVolume.Name,
						MountPath: "/app/config.json",
						SubPath:   Secret.ConfigKey,
					},
					{
						Name:      dataVolume.Name,
						MountPath: "/app/db-data",
					},
				},
			},
			{
				Name:  "replicate-litestream",
				Image: config.LITESTREAM_IMAGE,
				Command: []string{
					"litestream",
					"replicate",
				},
				VolumeMounts: []core.VolumeMount{
					{
						Name:      dataVolume.Name,
						MountPath: "/db-data",
					},
					{
						Name:      litestreamVolume.Name,
						MountPath: "/etc/litestream.yml",
						SubPath:   Secret.LitestreamKey,
					},
				},
			},
		},
		Volumes: []core.Volume{
			dataVolume,
			configVolume,
			litestreamVolume,
		},
	}
	return kube.NewDeployment(meta, podSpec)
}

func Ingress() net.Ingress {
	rules := []kube.IngressRule{
		{
			Host:        config.KUBERNETES_HOST,
			ServiceName: SRV.Name,
			PortNumber:  config.KUBERNETES_PORT,
		},
	}
	return kube.Ingress(Namespace.Name, rules, true)
}
