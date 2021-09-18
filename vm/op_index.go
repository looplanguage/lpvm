package vm

import (
	"fmt"
	"github.com/looplanguage/loop/models/object"
)

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASHMAP:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported on: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	array := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(array.Elements)) - 1

	if i > max {
		return vm.push(Null)
	}

	return vm.push(array.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	array := left.(*object.HashMap)
	i := index.(object.Hashable).Hash()

	if elem, ok := array.Pairs[i]; ok {
		return vm.push(elem.Value)
	}

	return vm.push(Null)
}
