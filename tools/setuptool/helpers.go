package main

import (
	"runtime"
	"strconv"
)

func classedPrintln(class string, msg string) {
	println(class + msg)
}

func errPrintln(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	classedPrintln("ERROR", ": FILE: "+filename+", LINE: "+strconv.Itoa(line)+", "+msg)
}

func infoPrintln(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	classedPrintln("INFO", ": FILE: "+filename+", LINE: "+strconv.Itoa(line)+", "+msg)
}
