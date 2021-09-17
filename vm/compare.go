package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/code"
	"github.com/looplanguage/loop/models/object"
)

func (vm *VM) compareOperator(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER || right.Type() == object.INTEGER {
		return vm.compareInteger(op, left, right)
	}

	switch op {
	case code.OpEquals:
		return vm.push(getBoolean(left == right))
	case code.OpNotEquals:
		return vm.push(getBoolean(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) compareInteger(op code.OpCode, left object.Object, right object.Object) error {
	leftObj, ok := left.(*object.Integer)
	if !ok {
		return fmt.Errorf("left comparison is not of type integer. got=%q", left.Type())
	}

	rightObj, ok := right.(*object.Integer)
	if !ok {
		return fmt.Errorf("right comparison is not of type integer. got=%q", left.Type())
	}

	leftValue := leftObj.Value
	rightValue := rightObj.Value

	switch op {
	case code.OpEquals:
		return vm.push(getBoolean(leftValue == rightValue))
	case code.OpNotEquals:
		return vm.push(getBoolean(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(getBoolean(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func getBoolean(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}
