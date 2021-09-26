package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/lpvm/flags"
	"github.com/looplanguage/lpvm/repl"
	"github.com/looplanguage/lpvm/vm"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flags.Parse()

	if flags.File == "" {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	compiler.RegisterGobTypes()
	bts, err := ioutil.ReadFile(flags.File)

	if err != nil {
		log.Fatalln(err)
	}

	var constantBytes bytes.Buffer
	constantBytes.Write(bts)

	dec := gob.NewDecoder(&constantBytes)
	var Bytecode compiler.Bytecode
	err = dec.Decode(&Bytecode)

	if err != nil {
		log.Fatal(err)
	}

	machine := vm.Create(&Bytecode)
	err = machine.Run(nil)

	if err != nil {
		log.Fatal(err)
	}

	if machine.LastPoppedStackElem() != nil {
		fmt.Println(machine.LastPoppedStackElem().Inspect())
	}
}
