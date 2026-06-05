package main

import (
	"git.danicos.dev/daniel/curious-ape/pkg/config"
	"git.danicos.dev/daniel/curious-ape/pkg/deployment"
	"git.danicos.dev/daniel/curious-ape/pkg/deployment/secrets"
)

func main() {
	base := deployment.BaseStack()
	base.MarshalYaml(config.KUBERNETES_DEPLOYMENT)

	k3s := deployment.K3sStack()
	k3s.Add("secrets", secrets.ApeSecret)
	k3s.MarshalYamlFlat(config.KUBERNETES_DEPLOYMENT + "/overlays/k3s")
}
