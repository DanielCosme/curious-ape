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

	"github.com/danielcosme/curious-ape/pkg/config"
	"github.com/danielcosme/curious-ape/pkg/target"
)

const tmpDir = "./tmp"

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
	}
	r = target.NewRunner(Env, nil)

	prodHost = fmt.Sprintf("%s@%s", config.PROD_ADMIN, config.PROD_HOST)
}

func Run() error {
	mg.Deps(Build)
	return r.RunV("run", target.New(binOutput))
}

// Builds Binary
func Build() error {
	c := target.New("./scripts/build.fish")
	return r.RunV("build", c)
}

// Builds production static Binary
func Build_prod() error {
	c := target.NewA("./scripts/build.fish", "prod")
	return r.RunV("build", c)
}

func Deploy() error {
	mg.SerialDeps(Install)

	enc := target.NewA("ssh", prodHost, "sudo", "systemctl", "restart", "curious-ape")
	return r.RunV("deploy", enc)
}

func Install() error {
	mg.SerialDeps(Build_prod, Decrypt)

	installDir := tmpDir + "/deployment"
	ts := []target.Target{
		target.NewA("mkdir", "-p", installDir),
		target.NewA("mv", config.DEPLOYMENT_DIR+"/ape", installDir+"/ape"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/curious-ape.service", installDir+"/curious-ape.service"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/secrets/config.json", installDir+"/config.json"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/litestream/litestream", installDir+"/litestream"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/litestream/etc/litestream.service", installDir+"/litestream.service"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/secrets/litestream.yaml", installDir+"/litestream.yaml"),
		target.NewA("cp", "./scripts/install.fish", installDir+"/install.fish"),
		target.NewA("cp", config.DEPLOYMENT_DIR+"/envfile", installDir+"/envfile"),
		target.NewA("rm", "-r", config.DEPLOYMENT_DIR+"/secrets"),

		target.NewA("rsync", "--compress", "--recursive", installDir, prodHost+":/tmp"),
		target.NewA("ssh", prodHost, "/tmp/deployment/install.fish"),
		target.NewA("rm", "-r", installDir),
	}
	return runSteps("install server", ts)
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

// Install development environment tools
func Tools() {
	ts := []target.Target{
		// target.NewA("go", "install", "github.com/air-verse/air@latest"),
		target.NewA("go", "install", "-tags", "'sqlite3'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"),
		target.NewA("go", "get", "-tool", "github.com/rakyll/gotest@latest"),
		target.NewA("go", "get", "-tool", "honnef.co/go/tools/cmd/staticcheck@latest"),
		target.NewA("go", "get", "-tool", "github.com/stephenafamo/bob/gen/bobgen-sqlite@v0.42.0"),
		target.NewA("go", "get", "-tool", "github.com/magefile/mage@latest"),
	}
	runSteps("tools", ts)
}

func Logs_Prod() error {
	t := target.NewA("ssh", prodHost, "sudo", "journalctl", "--lines", "40", "-fu", "curious-ape.service")
	return r.RunV("logs prod", t)
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

// Push Pushes repository to GitHub and creates new Version
func Push() error {
	mg.SerialDeps(Ci)

	diff := target.NewA("git", "diff", "--exit-code").SetMsg("working tree cannot be dirty")
	diff.Silent = true

	branch, err := sh.Output("git", "branch", "--show-current")
	assert(err)
	ts := []target.Target{
		diff,
		target.NewA("test", branch, "=", "main").SetMsg("branch has to be main"),
	}
	err = runSteps("push", ts)
	assert(err)

	err = Tag()
	assert(err)

	ts = []target.Target{
		// NOTE: not pushing to origin (GitHub)
		target.NewA("git", "push", "apex", "--follow-tags"),
	}
	return runSteps("push", ts)
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
