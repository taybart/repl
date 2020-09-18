package repl

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	Stdout io.Writer
	Stdin  io.Reader
)

var config = struct {
	Cursor string
}{
	Cursor: "->",
}

func SetCursor(c string) {
	config.Cursor = c
}

func Do(eval func(string) bool) error {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		err = fmt.Errorf("can't open /dev/tty: %w", err)
		return err
	}
	fd := int(tty.Fd())
	termState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, termState)

	n := terminal.NewTerminal(os.Stdin, fmt.Sprintf("%s ", config.Cursor))
	n.SetSize(int(^uint(0)>>1), 0)
	Stdout = os.Stdout
	for {
		ln, err := n.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		terminal.Restore(fd, termState)
		if err != nil {
			return err
		}

		quit := eval(ln)
		if quit {
			return nil
		}

		termState, err = terminal.MakeRaw(fd)
		if err != nil {
			return err
		}
	}
}

func newTerm() {

}
