package pexception

import "fmt"

func RecoverFromPanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}
