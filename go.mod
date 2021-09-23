module github.com/looplanguage/lpvm

go 1.17

require (
	github.com/looplanguage/compiler v0.2.0
	github.com/looplanguage/loop v0.5.1
)

replace github.com/looplanguage/loop => ../loop
replace github.com/looplanguage/compiler => ../compiler