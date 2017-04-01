package utils

import (
    "runtime"
    "log"
)

func PrintRecover(e interface {}) interface {} {
	if e != nil {
		log.Println("[!] recover:", e)
		for skip := 1; ; skip++ {
			pc, file, line, isOK := runtime.Caller(skip)
			if !isOK {
				break
			}

			f := runtime.FuncForPC(pc)
			log.Printf("    %s:%d %s()\n", file, line, f.Name())
		}
	}
	return e
}
