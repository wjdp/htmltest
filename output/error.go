package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// CheckErrorGeneric aborts if error present with its message.
// Use for when error is not abnormal.
func CheckErrorGeneric(err error) {
	if err != nil {
		AbortWith(err.Error())
	}
}

// CheckErrorPanic panics when error is present with its message.
// Use when the presence of the error is not expected under the great majority
// of circumstances.
func CheckErrorPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// AbortWith prints the provided message and ends execution of the program
// with code 1.
func AbortWith(a ...interface{}) {
	color.Set(color.FgRed)
	fmt.Println(a...)
	color.Unset()
	os.Exit(1)
}
