package main

import (
	"github.com/looplanguage/lpvm/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
