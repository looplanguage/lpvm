package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/loop/models/object"
	"net/http"
)

type HttpHandler struct {
	vm      *VM
	Closure *object.Closure
}

func (httpHandler *HttpHandler) handleHttpRequest(w http.ResponseWriter, req *http.Request) {
	newVM := CreateWithStore(&compiler.Bytecode{
		Constants:    httpHandler.vm.constants,
		Instructions: httpHandler.Closure.Fn.Instructions,
	}, httpHandler.vm.globals, httpHandler.vm.variables)

	newVM.push(&object.BuiltinFunction{Function: func(args []object.Object) object.Object {
		if len(args) != 1 {
			return &object.Error{Message: "expected function 'writeHttp'. got=NULL"}
		}
		if args[0].Type() != "STRING" {
			fmt.Fprintf(w, "wrong type to 'writeHttp'. expected=%q. got=%q", "STRING", args[0].Type())

			return &object.Null{}
		}

		write := args[0].(*object.String).Value

		fmt.Fprintf(w, write)

		return &object.Null{}
	}})

	err := newVM.Run(nil)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func (vm *VM) replaceBuiltinFunctions() {
	var builtins = []object.Builtin{
		{
			Name: "httpServer",
			Builtin: &object.BuiltinFunction{Function: func(args []object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: "(not enough arguments) First argument is not a function"}
				}

				handlers, ok := args[0].(*object.HashMap)

				if !ok {
					return &object.Error{Message: "First argument is not a function"}
				}

				for _, handler := range handlers.Pairs {
					if handler.Key.Type() != "STRING" {
						return &object.Error{Message: fmt.Sprintf("handler requires to be a string. got=%q", handler.Key.Type())}
					}
					if handler.Value.Type() != "CLOSURE" {
						return &object.Error{Message: fmt.Sprintf("handler requires to be a function. got=%q", handler.Value.Type())}
					}

					pattern := handler.Key.(*object.String)
					closure := handler.Value.(*object.Closure)

					handler := HttpHandler{
						vm:      vm,
						Closure: closure,
					}

					http.HandleFunc(pattern.Value, handler.handleHttpRequest)
				}

				http.ListenAndServe(":8090", nil)
				return &object.String{Value: "Server Stopped"}
			}},
		},
	}

	// Replace
	for _, rb := range builtins {
		for k, b := range object.Builtins {
			if b.Name == rb.Name {
				object.Builtins[k].Builtin = rb.Builtin
			}
		}
	}
}
