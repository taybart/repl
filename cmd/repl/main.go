package main

import (
	"fmt"
	"os"

	"github.com/taybart/repl"
)

func main() {
	r, err := repl.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.SetCursor("»")

	r.Do(func(cmd string) bool {
		return cmd == "quit" || repl.DoShell(cmd) != nil
	})
}
