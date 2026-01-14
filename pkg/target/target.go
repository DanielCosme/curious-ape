package target

import "fmt"
import "github.com/magefile/mage/sh"
import "github.com/fatih/color"

type Target struct {
	Bin  string
	args []string
}

func New(bin string) (c Target) {
	c.Bin = bin
	return NewA(bin)
}

func NewA(bin string, args ...string) (c Target) {
	c.Bin = bin
	c.args = args
	return c
}

func (t *Target) Args(args ...string) {
	t.args = append(t.args, args...)
}

type Runner struct {
	env   map[string]string
	color *color.Color
}

func NewRunner(env map[string]string, c *color.Color) (r Runner) {
	r.color = c
	if c == nil {
		r.color = color.New(color.FgGreen)
	}
	r.env = env
	return r
}

func (r *Runner) RunV(targetName string, t Target) error {
	tar := "Target: "
	td := fmt.Sprintf("<%s> ", targetName)
	r.color.Println(tar, td, t.Bin, t.args)
	return sh.RunWithV(r.env, t.Bin, t.args...)
}
