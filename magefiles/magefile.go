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

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/target"
)

const tmpDir = config.TMP_DIR

var Default = Live.All
var Aliases = map[string]any{
	"v": Version,
	"r": Run,
	"t": Test,
}

var r target.Runner
var devOutput string
var prodOutput string
var dbLocation string

func init() {
	devOutput = fmt.Sprintf("%s/%s", tmpDir, config.APP_NAME)
	prodOutput = fmt.Sprintf("bin/%s", config.APP_NAME)
	dbLocation = devOutput + ".db"

	Env := map[string]string{
		config.ENVIRONMENT: "dev",
		"PROD_OUTPUT":      prodOutput,
		"DEV_OUTPUT":       devOutput,
		"SECRETS_PATH":     config.DEPLOYMENT_DIR + "/secrets",
		"ENC_SECRETS_PATH": config.DEPLOYMENT_DIR + "/enc",
		"KUBE_SECRETS":     config.KUBERNETES_SECRETS,
		"KUBE_ENC_SECRETS": config.KUBERNETES_ENC_SECRETS,
	}
	r = target.NewRunner(Env, nil)
}

// Build and run server
func Run() error {
	mg.Deps(Build)
	return r.RunV("run", target.New(prodOutput))
}

// Builds production static Binary
func Build() error {
	mg.Deps(Build_Templ)
	c := target.New("./scripts/build.sh")
	return r.RunV("build", c)
}

func Build_Templ() error {
	return r.RunV("build templ", target.NewA("go", "tool", "templ", "generate"))
}

// Install development environment tools
func Tools() {
	ts := []target.Target{
		target.NewA("go", "get", "-tool", "github.com/air-verse/air@latest"),
		target.NewA("go", "install", "-tags", "'sqlite3'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"),
		target.NewA("go", "get", "-tool", "github.com/rakyll/gotest@latest"),
		target.NewA("go", "get", "-tool", "honnef.co/go/tools/cmd/staticcheck@latest"),
		target.NewA("go", "get", "-tool", "github.com/stephenafamo/bob/gen/bobgen-sqlite@v0.45.0"),
		target.NewA("go", "get", "-tool", "github.com/magefile/mage@v1.17.2"),
		target.NewA("go", "get", "-tool", "github.com/a-h/templ/cmd/templ"),
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

// Download Javascript Libraries
func Download() error {
	return r.RunV("download JS libraries", target.NewA("go", "run", "./cmd/downloader/main.go"))
}

// Run test and audit tasks
func Ci() {
	mg.SerialDeps(Test, Audit)
}

func Test() error {
	return r.RunV("test", target.NewA("go", "tool", "gotest", "./..."))
}

func Version() error {
	return sh.RunV("echo", config.VERSION)
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
