package target

import "fmt"
import "github.com/magefile/mage/sh"
import "github.com/fatih/color"

type Target struct {
	Bin    string
	args   []string
	Silent bool
	Msg    string
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

func (t Target) SetMsg(m string) Target {
	t.Msg = m
	return t
}

type Runner struct {
	env   map[string]string
	info  *color.Color
	error *color.Color
}

func NewRunner(env map[string]string, c *color.Color) (r Runner) {
	r.info = c
	if c == nil {
		r.info = color.New(color.FgGreen)
	}
	r.error = color.New(color.FgRed)
	r.env = env
	return r
}

func (r *Runner) RunV(targetName string, t Target) error {
	return r.run(true, targetName, t)
}

func (r *Runner) Run(targetName string, t Target) error {
	return r.run(false, targetName, t)
}

func (r *Runner) run(verbose bool, targetName string, t Target) error {
	tar := "Target: "
	td := fmt.Sprintf("<%s> ", targetName)
	if len(t.args) == 0 {
		r.info.Println(tar, td, t.Bin)
	} else {
		r.info.Println(tar, td, t.Bin, t.args)
	}
	if t.Msg != "" {
		r.error.Println("  ", t.Msg)
	}
	if verbose {
		return sh.RunWithV(r.env, t.Bin, t.args...)
	}
	return sh.RunWith(r.env, t.Bin, t.args...)
}
