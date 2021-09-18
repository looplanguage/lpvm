package vm

import (
	"fmt"
	"github.com/looplanguage/loop/models/object"
	"strconv"
)

func (vm *VM) OpAdd() error {
	right := vm.pop()
	left := vm.pop()

	switch lValue := left.(type) {
	case *object.String:
		if rValue, ok := right.(*object.String); ok {
			vm.push(&object.String{Value: lValue.Value + rValue.Value})
			return nil
		} else if rValue, ok := right.(*object.Integer); ok {
			vm.push(&object.String{Value: lValue.Value + strconv.FormatInt(rValue.Value, 10)})
			return nil
		}
	case *object.Integer:
		if rValue, ok := right.(*object.String); ok {
			vm.push(&object.String{Value: strconv.FormatInt(lValue.Value, 10) + rValue.Value})
			return nil
		} else if rValue, ok := right.(*object.Integer); ok {
			vm.push(&object.Integer{Value: rValue.Value + lValue.Value})
			return nil
		}
	}

	return fmt.Errorf("unknown operation exception. got=%q. got=%q", left.Type(), right.Type())
}
