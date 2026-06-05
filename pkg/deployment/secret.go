package deployment

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	core "k8s.io/api/core/v1"
)

var ApeSecret core.Secret

func init() {
	s := Secret
	meta := kube.NewMetadata(s.Name, Namespace)
	ApeSecret = kube.SecretFromFile(s.ConfigKey, "./deployment/secrets/config.json", meta)
	ApeSecret.Data[s.LitestreamKey] = kube.ReadFileBytes("deployment/secrets/litestream.yaml")
}
