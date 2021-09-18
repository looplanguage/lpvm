package vm

import (
	"fmt"
	"github.com/looplanguage/loop/models/object"
)

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{
			Key:   key,
			Value: value,
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("incorrect key type: %s", key.Type())
		}

		hashedPairs[hashKey.Hash()] = pair
	}

	return &object.HashMap{Pairs: hashedPairs}, nil
}
