package utils

import (
	"fmt"
	"os"

	"github.com/depado/gorich"
)

var Debug bool

func OkPrintln(args ...interface{}) {
	gorich.Println("[green]»[/]", fmt.Sprint(args...))
}

func OkPrintf(format string, args ...interface{}) {
	gorich.Printf("[green]»[/] "+format+"\n", args...)
}

func ErrPrintln(args ...interface{}) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
}

func ErrPrintf(format string, args ...interface{}) {
	gorich.Printf("[red]»[/] "+format+"\n", args...)
}

func FatalPrintln(args ...interface{}) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
	os.Exit(1)
}

func FatalPrintf(format string, args ...interface{}) {
	gorich.Printf("[red]»[/] "+format+"\n", args...)
	os.Exit(1)
}

func ExitPrintln(args ...interface{}) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
	os.Exit(0)
}

func ExitPrintf(format string, args ...interface{}) {
	gorich.Printf("[red]»[/] "+format+"\n", args...)
	os.Exit(0)
}

func DebugPrintln(args ...interface{}) {
	if Debug {
		gorich.Println("[dim cyan]»[/]", fmt.Sprint(args...))
	}
}

func DebugPrintf(format string, args ...interface{}) {
	if Debug {
		gorich.Printf("[dim cyan]»[/] "+format+"\n", args...)
	}
}
