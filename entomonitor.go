package main

import (
	"fmt"
	"github.com/droundy/goopt"
)

var all = goopt.Flag([]string{"-a", "--all"}, []string{"--interactive"}, "hello", "goodbye")

func main() {
	goopt.Parse(func() []string { return nil })
	fmt.Println("Args are", goopt.Args)
}
