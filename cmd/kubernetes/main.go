package main

import (
	"git.danicos.dev/daniel/curious-ape/pkg/config"
	"git.danicos.dev/daniel/curious-ape/pkg/deployment"
)

func main() {
	s := deployment.BaseStack()
	s.MarshalYaml(config.KUBERNETES_DEPLOYMENT)
}
