package utils

import (
    "runtime/debug"
    "log"
)

func PrintRecover(e interface{}) interface{} {
    if e != nil {
        log.Println("[!] catch exception:", e)
        debug.PrintStack()
    }
    return e
}
