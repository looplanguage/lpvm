package repl

import (
	"bufio"
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/loop/lexer"
	"github.com/looplanguage/loop/parser"
	"github.com/looplanguage/lpvm/vm"
	"io"
)

// TODO: For testing, remove in eventual build & replace with it's own executable.

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	i := 0

	for {
		i++
		io.WriteString(out, fmt.Sprintf("%d", i))
		io.WriteString(out, " > ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.Create(line)
		p := parser.Create(l)

		program := p.Parse()
		if len(p.Errors) != 0 {
			for _, e := range p.Errors {
				fmt.Println(e)
			}
			continue
		}

		comp := compiler.Create()
		err := comp.Compile(program)

		if err != nil {
			fmt.Fprintf(out, "Compilation failed. \n%s\n", err)
			continue
		}

		machine := vm.Create(comp.Bytecode())

		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "vm failed running bytecode with: \n%s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
