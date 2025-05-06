package logger

import (
	"fmt"
	"runtime"
	"strings"
)

func fileAndFuncInfo(skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)

	funcName := runtime.FuncForPC(pc).Name()

	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
		dot := strings.LastIndex(funcName, ".")
		if dot >= 0 {
			funcName = funcName[dot+1:]
		}
	}
	return funcName, fmt.Sprintf("%s:%d", file, line)
}
