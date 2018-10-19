package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Green is a simple green foreground color
var Green = color.New(color.FgGreen)

// OkPrefix is the prefix that should prefix output when everything is ok
var OkPrefix = Green.Sprint("»")

// ErrPrefix is the prefix that should output when an error occurred
var ErrPrefix = color.New(color.FgRed).Sprint("»")

// OkPrintln prints with a green prefix
func OkPrintln(opts ...interface{}) {
	fmt.Println(append([]interface{}{OkPrefix}, opts...)...)
}

// ErrPrintln prints with a red prefix
func ErrPrintln(opts ...interface{}) {
	fmt.Println(append([]interface{}{ErrPrefix}, opts...)...)
}

// FatalPrintln prints out information with a red prefix and exits the program
func FatalPrintln(opts ...interface{}) {
	fmt.Println(append([]interface{}{ErrPrefix}, opts...)...)
	os.Exit(1)
}
