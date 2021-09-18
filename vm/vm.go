package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/loop/models/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames     []*Frame
	frameIndex int
}

func Create(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)

	frames[0] = mainFrame

	return &VM{
		constants:  bytecode.Constants,
		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    make([]object.Object, GlobalsSize),
		frames:     frames,
		frameIndex: 1,
	}
}

func CreateWithStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := Create(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.frameIndex] = f
	vm.frameIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.frameIndex--

	return vm.frames[vm.frameIndex]
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
