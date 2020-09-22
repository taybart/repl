package repl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
)

type Repl struct {
	state    *terminal.State
	terminal *terminal.Terminal
	fd       int
	tty      *os.File
}

func New() (r *Repl, err error) {
	tty := os.Stdin
	r = &Repl{
		fd:  int(tty.Fd()),
		tty: tty,
	}
	r.MakeRaw()
	r.SetCursor("->")

	return
}

func (r *Repl) SetCursor(c string) {
	r.terminal = terminal.NewTerminal(r.tty, fmt.Sprintf("%s ", c))
}

func (r *Repl) MakeRaw() (err error) {
	r.state, err = terminal.MakeRaw(r.fd)
	return
}

func (r *Repl) Restore() error {
	return terminal.Restore(r.fd, r.state)
}

func (r *Repl) Do(eval func(string) bool) error {
	defer r.Restore()
	for {
		ln, err := r.terminal.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		err = r.Restore()
		if err != nil {
			return err
		}

		quit := eval(ln)
		if quit {
			return nil
		}

		err = r.MakeRaw()
		if err != nil {
			return err
		}

	}
}

func DoShell(cmd string) error {
	args := strings.Split(cmd, " ")
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
