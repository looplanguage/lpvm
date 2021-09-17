package vm

import (
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
				return nil
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue + rightValue

			vm.push(&object.Integer{Value: result})
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
		}
	}

	return nil
}
