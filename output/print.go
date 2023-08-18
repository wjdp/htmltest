package output

import (
	"fmt"

	"github.com/fatih/color"
)

func Warn(a ...interface{}) {
	color.Set(color.FgYellow)
	fmt.Println(a...)
	color.Unset()
}

func Debug(a ...interface{}) {
	fmt.Println(a...)
}
