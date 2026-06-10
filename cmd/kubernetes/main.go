package main

import (
	"io/fs"
	"os"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/deployment"
)

func main() {
	base := deployment.BaseStack()
	base.MarshalYaml(config.KUBERNETES_DEPLOYMENT)

	k3s := deployment.K3sStack()
	k3s.MarshalYaml(config.KUBERNETES_DEPLOYMENT + "/overlays")

	os.MkdirAll(config.KUBERNETES_SECRETS, fs.ModeDir|0755)
	secrets := deployment.SecretsStack()
	secrets.MarshalYamlFlat(config.KUBERNETES_SECRETS)
}
