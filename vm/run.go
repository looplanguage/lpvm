package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/code"
	"github.com/looplanguage/loop/models/object"
)

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.OpCode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			err := vm.OpAdd()
			if err != nil {
				return err
			}
		case code.OpMultiply:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue * rightValue

			vm.push(&object.Integer{Value: result})
		case code.OpDivide:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue / rightValue

			vm.push(&object.Integer{Value: result})
		case code.OpSubtract:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue - rightValue

			vm.push(&object.Integer{Value: result})
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			vm.push(True)
		case code.OpFalse:
			vm.push(False)
		case code.OpEquals, code.OpNotEquals, code.OpGreaterThan:
			err := vm.compareOperator(op)
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpIfNotTrue:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			obj, ok := condition.(*object.Boolean)
			if !ok {
				return fmt.Errorf("condition is not boolean. got=%q", condition.Type())
			}

			if !obj.Value {
				ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
