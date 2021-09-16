package vm

import "github.com/looplanguage/compiler/code"

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
		}
	}

	return nil
}
