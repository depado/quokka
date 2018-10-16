package utils

import (
	"fmt"

	"github.com/fatih/color"
)

// Green is a simple green foreground color
var Green = color.New(color.FgGreen)

// OkPrefix is the prefix that should prefix output when everything is ok
var OkPrefix = Green.Sprint("»")

// ErrPrefix is the prefix that should output when an error occured
var ErrPrefix = color.New(color.FgRed).Sprint("»")

// OkPrintln prints with a green prefix
func OkPrintln(opts ...interface{}) {
	fmt.Println(append([]interface{}{OkPrefix}, opts...)...)
}

// ErrPrintln prints with a red prefix
func ErrPrintln(opts ...interface{}) {
	fmt.Println(append([]interface{}{ErrPrefix}, opts...)...)
}
