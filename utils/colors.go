package utils

import (
	"fmt"
	"os"

	"github.com/depado/gorich"
)

var Debug bool

func OkPrintln(args ...any) {
	gorich.Println("[green]»[/]", fmt.Sprint(args...))
}

func OkPrintf(format string, args ...any) {
	gorich.Printf("[green]»[/] "+format+"\n", args...)
}

func ErrPrintln(args ...any) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
}

func ErrPrintf(format string, args ...any) {
	gorich.Printf("[red]»[/] "+format+"\n", args...)
}

func FatalPrintln(args ...any) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
	os.Exit(1)
}

func ExitPrintln(args ...any) {
	gorich.Println("[red]»[/]", fmt.Sprint(args...))
	os.Exit(0)
}

func DebugPrintf(format string, args ...any) {
	if Debug {
		gorich.Printf("[dim cyan]»[/] "+format+"\n", args...)
	}
}
