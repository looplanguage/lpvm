package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/loop/models/object"
	"log"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants []object.Object
	variables []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames     []*Frame
	frameIndex int
}

func Create(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)

	frames[0] = mainFrame

	return &VM{
		constants:  bytecode.Constants,
		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    make([]object.Object, GlobalsSize),
		frames:     frames,
		frameIndex: 1,
		variables:  make([]object.Object, GlobalsSize),
	}
}

func CreateWithStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := Create(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) callFunction(numArgs int) error {
	switch fn := vm.stack[vm.sp-1-numArgs].(type) {
	case *object.Closure:
		return vm.callUserClosure(fn, numArgs)
	case *object.BuiltinFunction:
		return vm.callBuiltinFunction(fn, numArgs)
	}

	return fmt.Errorf("attempt to call non-function. got=%q", vm.stack[vm.sp-1].Type())
}

func (vm *VM) callBuiltinFunction(fn *object.BuiltinFunction, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := fn.Function(args)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		return vm.push(result)
	}

	return vm.push(Null)
}

type MemoizedFunction struct {
	Id     int
	Args   []object.Object
	Result object.Object
}

type MemoizedKey struct {
	Id   int
	Args string
}

// Cache needs to be on function ID & arguments, **NOT** just its ID!
var MemoizedFunctions = make(map[MemoizedKey]*MemoizedFunction, 5000)
var GetFunctionResult *MemoizedFunction

func (vm *VM) callUserClosure(cl *object.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments. expected=%d. got=%d", cl.Fn.NumParameters, numArgs)
	}

	var args = []object.Object{}

	arg := 0

	for arg < numArgs {
		args = append(args, vm.stack[vm.sp-numArgs+arg])
		arg++
	}

	concatArgs := ""

	for _, a := range args {
		concatArgs += a.Type() + a.Inspect()
	}

	var MF *MemoizedFunction
	key := MemoizedKey{
		Id:   cl.Fn.Id,
		Args: concatArgs,
	}

	MF = MemoizedFunctions[key]

	if MF != nil && MF.Result != nil {
		vm.sp = vm.sp - numArgs

		for i := 0; i < numArgs; i++ {
			vm.pop()
		}

		vm.push(MF.Result)
	} else {
		/*val := &MemoizedFunction{
			Id:   cl.Fn.Id,
			Args: args,
		}

		MemoizedFunctions[key] = val

		GetFunctionResult = val*/

		frame := NewFrame(cl, vm.sp-numArgs)
		vm.pushFrame(frame)

		vm.sp = frame.basePointer + cl.Fn.NumLocals
	}

	return nil
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.frameIndex] = f
	vm.frameIndex++
}

func (vm *VM) pushClosure(constIndex int, numFree int) error {
	constant := vm.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		log.Fatalf("type is not function. got=%+v\n", constant)
	}

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &object.Closure{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) popFrame(returnValue object.Object) *Frame {
	if GetFunctionResult != nil {
		GetFunctionResult.Result = returnValue
		GetFunctionResult = nil
	}

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
