package debug

import "fmt"

var DEBUG = true

func LogPrintln(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func LogPrint(a ...any) {
	if DEBUG {
		fmt.Print(a...)
	}
}

func LogPrintf(format string, a ...any) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}
