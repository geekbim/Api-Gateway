package logger

import (
	"fmt"
	"log"
)

func Info(message string) {
	log.Println(fmt.Sprintf(" INFO: %v", message))
}

func Error(err error) {
	log.Println(fmt.Sprintf(" ERROR: %v", err.Error()))
}
