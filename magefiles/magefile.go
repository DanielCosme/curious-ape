//go:build mage

package main

// NOTE: Mage https://github.com/magefile/mage
import (
	"fmt"
	// "os"
	// "os/exec"
	"github.com/magefile/mage/sh"
	// Mage other packages
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	// "github.com/magefile/mage/mage"
	// "github.com/magefile/mage/parse"
	// "github.com/magefile/mage/target"

	"github.com/danielcosme/curious-ape/pkg/root"
	"github.com/danielcosme/curious-ape/pkg/target"
)

const tmpDir = "./tmp"

var Env = map[string]string{
	root.ENVIRONMENT: "dev",
}

var Default = Run

var Aliases = map[string]any{
	"v": Version,
	"r": Run,
	"t": Test,
}

var r target.Runner
var binOutput string

func init() {
	binOutput = fmt.Sprintf("%s/%s", tmpDir, root.APP)
	r = target.NewRunner(Env, nil)
}

func Run() error {
	mg.Deps(Build)
	return r.RunV("run", target.New(binOutput))
}

// Builds Binary
func Build() error {
	c := target.New("go")
	versionFlag := fmt.Sprintf("-X main.version=%s-dev", root.VERSION)
	c.Args("build", "-ldflags", versionFlag, "-o="+binOutput, "./cmd/web")
	return r.RunV("build", c)
}

// Install development environment tools
func Tools() {
	ts := []target.Target{
		target.NewA("go", "install", "github.com/air-verse/air@latesta"),
		target.NewA("go", "install", "-tags", "'sqlite3'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"),
		target.NewA("go", "get", "-tool", "github.com/rakyll/gotest@latest"),
		target.NewA("go", "get", "-tool", "honnef.co/go/tools/cmd/staticcheck@latest"),
		target.NewA("go", "get", "-tool", "github.com/stephenafamo/bob/gen/bobgen-sqlite@v0.42.0"),
		target.NewA("go", "get", "-tool", "github.com/magefile/mage@latest"),
	}
	runSteps("tools", ts)
}

func Test() error {
	return r.RunV("test", target.NewA("go", "tool", "gotest", "./..."))
}

func Tag() error {
	return r.RunV("tag", target.NewA("git", "tag", root.VERSION))
}

func Version() error {
	return sh.RunV("echo", root.VERSION)
}

func runSteps(target string, ts []target.Target) error {
	for _, t := range ts {
		err := r.RunV(target, t)
		assert(err)
	}
	return nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

// A custom install step if you need your bin someplace other than go/bin
// func Install() error {
// 	mg.Deps(Build)
// 	fmt.Println("Installing...")
// 	return os.Rename("./MyApp", "/usr/bin/MyApp")
// }
//
// // Manage your deps, or running package managers.
// func InstallDeps() error {
// 	fmt.Println("Installing Deps...")
// 	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
// 	return cmd.Run()
// }
//
// // Clean up after yourself
// func Clean() {
// 	fmt.Println("Cleaning...")
// 	os.RemoveAll("MyApp")
// }

/*
type Build mg.Namespace

// Builds the site using hugo.
func (Build) Site() error {
  return nil
}

// Builds the pdf docs.
func (Build) Docs() {}

$ mage build:site
*/
