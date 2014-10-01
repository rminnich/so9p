package so9p

import (
	"log"
)

var DebugPrint = true

func DebugPrintf(fmt string, a ...interface{}) {
	if DebugPrint {
		log.Printf(fmt, a)
	}
}
