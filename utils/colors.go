package utils

import (
	"fmt"
	"log"
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
	log.Println(append([]interface{}{OkPrefix}, opts...)...)
}

// OkSprintln will return the string with a green prefix
func OkSprintln(opts ...interface{}) string {
	return fmt.Sprintln(append([]interface{}{OkPrefix}, opts...)...)
}

// ErrPrintln prints with a red prefix
func ErrPrintln(opts ...interface{}) {
	log.Println(append([]interface{}{ErrPrefix}, opts...)...)
}

// ErrSprintln will return the string with a red prefix
func ErrSprintln(opts ...interface{}) string {
	return fmt.Sprintln(append([]interface{}{ErrPrefix}, opts...)...)
}

// FatalPrintln prints out information with a red prefix and exits the program
func FatalPrintln(opts ...interface{}) {
	log.Println(append([]interface{}{ErrPrefix}, opts...)...)
	os.Exit(1)
}
