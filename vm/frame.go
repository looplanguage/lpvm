package vm

import (
	"github.com/looplanguage/compiler/code"
	"github.com/looplanguage/loop/models/object"
)

type Frame struct {
	closure     *object.Closure
	ip          int
	basePointer int
}

func NewFrame(fn *object.Closure, basePointer int) *Frame {
	return &Frame{closure: fn, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.closure.Fn.Instructions
}
