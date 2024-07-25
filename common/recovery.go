package common

import "log"

func Recovery() {
	if r := recover(); r != nil {
		log.Println("Recovered from panic: ", r)
	}
}
