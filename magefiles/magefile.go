package main

// NOTE: Mage https://github.com/magefile/mage
import (
	/*
		Mage other packages
		"github.com/magefile/mage/mage"
		"github.com/magefile/mage/parse"
		"github.com/magefile/mage/target"
	*/

	"fmt"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"

	"git.danicos.dev/daniel/curious-ape/pkg/config"
	"git.danicos.dev/daniel/curious-ape/pkg/target"
)

const tmpDir = config.TMP_DIR

var Default = Run
var Aliases = map[string]any{
	"v": Version,
	"r": Run,
	"t": Test,
}

var r target.Runner
var binOutput string
var dbLocation string
var prodHost string

func init() {
	binOutput = fmt.Sprintf("%s/%s", tmpDir, config.APP_NAME)
	dbLocation = binOutput + ".db"

	Env := map[string]string{
		config.ENVIRONMENT: "dev",
		"PROD_OUTPUT":      fmt.Sprintf("%s/%s", config.DEPLOYMENT_DIR, config.APP_NAME),
		"DEV_OUTPUT":       binOutput,
		"SECRETS_PATH":     config.DEPLOYMENT_DIR + "/secrets",
		"ENC_SECRETS_PATH": config.DEPLOYMENT_DIR + "/enc",
		"KUBE_SECRETS":     config.KUBERNETES_SECRETS,
		"KUBE_ENC_SECRETS": config.KUBERNETES_ENC_SECRETS,
	}
	r = target.NewRunner(Env, nil)
}

// "/overlays/secrets/secret, kustomization.yaml"
func Run() error {
	mg.Deps(Build)
	return r.RunV("run", target.New(binOutput))
}

// Builds Development Binary
func Build() error {
	c := target.New("./scripts/build.fish")
	return r.RunV("build", c)
}

// Builds production static Binary
func Build_prod() error {
	c := target.NewA("./scripts/build.fish", "prod")
	return r.RunV("build", c)
}

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

// Install development environment tools
func Tools() {
	ts := []target.Target{
		// target.NewA("go", "install", "github.com/air-verse/air@latest"),
		target.NewA("go", "install", "-tags", "'sqlite3'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"),
		target.NewA("go", "get", "-tool", "github.com/rakyll/gotest@latest"),
		target.NewA("go", "get", "-tool", "honnef.co/go/tools/cmd/staticcheck@latest"),
		target.NewA("go", "get", "-tool", "github.com/stephenafamo/bob/gen/bobgen-sqlite@v0.45.0"),
		target.NewA("go", "get", "-tool", "github.com/magefile/mage@v1.17.2"),
	}
	runSteps("tools", ts)
}

func Audit() error {
	ts := []target.Target{
		target.NewA("go", "mod", "tidy"),
		target.NewA("go", "mod", "verify"),
		target.NewA("go", "fmt", "./..."),
		target.NewA("go", "vet", "./..."),
		target.NewA("go", "tool", "staticcheck", "-checks='inherit,-ST1001'", "./cmd...", "./pkg..."),
	}
	return runSteps("audit", ts)
}

func Ci() {
	mg.SerialDeps(Test, Audit)
}

func Test() error {
	return r.RunV("test", target.NewA("go", "tool", "gotest", "./..."))
}

func Tag() error {
	return r.RunV("tag", target.NewA("git", "tag", config.VERSION))
}

func Version() error {
	return sh.RunV("echo", config.VERSION)
}

func Version_image() error {
	s := strings.TrimPrefix(config.VERSION, "v")
	return sh.RunV("echo", s)
}

func Registry() error {
	return sh.RunV("echo", config.REGISTRY)
}

func runSteps(target string, ts []target.Target) error {
	var err error
	for _, t := range ts {
		if t.Silent {
			err = r.Run(target, t)
		} else {
			err = r.RunV(target, t)
		}
		assert(err)
	}
	return nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
