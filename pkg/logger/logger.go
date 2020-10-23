package logger

import (
	"fmt"

	"github.com/fatih/color"
)

// Logger empty struct
type Logger struct {
}

// NewLogger returns Logger struct
func NewLogger() *Logger {
	return &Logger{}
}

// Info method for information stdout
func (l *Logger) Info(msg string, args ...interface{}) {
	if msg == "" {
		fmt.Println("")
		return
	}

	c := color.New(color.FgHiCyan)
	c.Println(fmt.Sprintf(msg, args...))
}

// Error method for error stdout
func (l *Logger) Error(err error) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("%#v", err))
}

// Instructions method for instruction stdout
func (l *Logger) Instructions(msg string, args ...interface{}) {
	white := color.New(color.FgHiWhite)
	white.Println("")
	white.Println(fmt.Sprintf(msg, args...))
}
