package main

import (
	"log"
	"runtime"
)

func FatalHandleError(err error) {
	if err != nil {
		// 1 in Caller() will actually log where the error happened.
		// 0 in Caller will not log where the error happened.
		pc, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[%s:%d]: %v", runtime.FuncForPC(pc).Name(), filename, line, err)
	}
}

func PrintHandleError(err error) {
	if err != nil {
		// 1 in Caller() will actually log where the error happened.
		// 0 in Caller will not log where the error happened.
		pc, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d]: %v\n", runtime.FuncForPC(pc).Name(), filename, line, err)
	}
}
