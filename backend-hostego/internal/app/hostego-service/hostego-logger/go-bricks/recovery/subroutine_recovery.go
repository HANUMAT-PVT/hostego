package Recovery

import (
	"fmt"
	"log"
)

func GoRoutineRecovery(fn func()) {
	go func() {
		defer RecoverFromPanic(nil)
		fn()
	}()
}

func GoRoutineCustomRecovery(fn func(), recoveryFunc func(err error)) {
	go func() {
		defer RecoverFromPanic(recoveryFunc)
		fn()
	}()
}

func RecoverFromPanic(callback func(err error)) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		newStack := stack(0)
		log.Println("%s\n", string(newStack))
		if callback != nil {
			callback(err)
		}
	}
}
