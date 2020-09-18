package main

import (
	"github.com/taybart/repl"
)

func main() {
	repl.SetCursor("»")
	repl.Do(func(cmd string) bool {
		return cmd == "quit" || repl.DoShell(cmd) != nil
	})
}
