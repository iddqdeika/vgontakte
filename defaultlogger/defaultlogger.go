package defaultlogger

import (
	"fmt"
	"time"
	"vgontakte/vgontakte"
)

var Logger vgontakte.Logger = &logger{}

type logger struct {
}

func (l *logger) Info(args ...interface{}) {
	fmt.Println(time.Now().String(), ": ", args)
}

func (l *logger) Infof(format string, args ...interface{}) {
	fmt.Printf(time.Now().String()+": "+format, args)
}

func (l *logger) Error(args ...interface{}) {
	fmt.Println(time.Now().String(), ": ", args)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	fmt.Printf(time.Now().String()+": "+format, args)
}
