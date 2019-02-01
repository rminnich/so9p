package so9p

import (
	"log"
)

var debugPrint = true

func debugPrintf(fmt string, a ...interface{}) {
	if debugPrint {
		log.Printf(fmt, a...)
	}
}

var verbose = func(string, ...interface{}) {}
