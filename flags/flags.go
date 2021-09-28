package flags

import (
	"flag"
	"strings"
)

var Optimizations = map[string]bool{}
var File string

func Parse() {
	optimizationsFlag := flag.String("o", "", "Specify which VM optimizations you'd like to activate. Seperated by a comma")

	flag.Parse()

	File = flag.Arg(0)

	optimizations := strings.Split(*optimizationsFlag, ",")

	for _, optimization := range optimizations {
		Optimizations[optimization] = true
	}
}

func OptimizationEnabled(optimization string) bool {
	if opt, ok := Optimizations[optimization]; ok {
		return opt
	}

	return false
}
