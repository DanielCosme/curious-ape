package main

import (
	"strings"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/target"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Generate kubenetes manifests from Go
func Build_kube() error {
	c := target.NewA("go", "run", "./cmd/kubernetes/main.go")
	return r.RunV("build kubernetes deployment", c)
}

// Encrypts all secrets.
func Encrypt() error {
	enc := target.NewA("./scripts/enc_dec.fish", "enc")
	return r.RunV("encryp secrets", enc)
}

// Decrypts all secrets.
func Decrypt() error {
	dec := target.NewA("./scripts/enc_dec.fish", "dec")
	return r.RunV("decrypt secrets", dec)
}

// Encrypts SOPS Secrets for Kubernetes GITOPS (Flux)
func Enc_sops() error {
	mg.SerialDeps(Encrypt, Decrypt)
	// SECRETS_ENC_PATH
	// AGE_KEY_NO_PQ
	// SECRETS_FOLDER
	dec := target.NewA("./scripts/encrypt_sops.sh")
	return r.RunV("Encrypt SOPS", dec)
}

func Flux_reconcile() error {
	return r.RunV("Reconcile GitOps deployment", target.NewA("flux", "--kubeconfig", "/home/daniel/.kube/charlie", "reconcile", "source", "git", "flux-system"))
}

func Version_image() error {
	s := strings.TrimPrefix(config.VERSION, "v")
	return sh.RunV("echo", s)
}

func Registry() error {
	return sh.RunV("echo", config.REGISTRY)
}

// Create tag and push to origin
func Tag() error {
	ts := []target.Target{
		target.NewA("git", "checkout", "main"),
		target.NewA("git", "pull"),
		target.NewA("git", "tag", config.VERSION),
		target.NewA("git", "push", "--tags"),
	}
	return runSteps("tag and push", ts)
}
