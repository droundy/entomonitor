//target:entomonitor
package main

import (
	"fmt"
	"os"
	"entomon"
	"github.com/droundy/goopt"
)

var all = goopt.Flag([]string{"-a", "--all"}, []string{"--interactive"}, "hello", "goodbye")

func dieOn(err os.Error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	goopt.Parse(func() []string { return nil })
	pname, err := entomon.ProjectName()
	dieOn(err)
	fmt.Println("Project name is", pname)
}
