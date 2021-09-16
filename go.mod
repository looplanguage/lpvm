module github.com/looplanguage/lpvm

go 1.17

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
)

require (
	github.com/fatih/color v1.12.0
	github.com/looplanguage/compiler v0.1.0
	github.com/looplanguage/loop v0.1.0
)

replace github.com/looplanguage/compiler => ../compiler
