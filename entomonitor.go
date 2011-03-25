//target:entomonitor
package main

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"entomon"
	"github.com/droundy/goopt"
)

var action = goopt.Alternatives([]string{"-A", "--action"},
	[]string{"help", "new-issue"}, "select the action to be performed")

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
	if *action == "help" {
		fmt.Println(goopt.Usage())
		os.Exit(0)
	}
	switch *action {
	case "new-issue":
		fmt.Print("What is your name? ")
		inp, e := bufio.NewReaderSize(os.Stdin, 1)
		dieOn(e)
		name, e := inp.ReadString('\n')
		if len(name) < 1 {
			name = "\n"
		}
		fmt.Print("What is the problem? ")
		bugtext, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println("Done here")
		dieOn(entomon.NewIssue(name[:len(name)-1], string(bugtext)))
	default:
		fmt.Println("I should do", *action)
		os.Exit(1)
	}
}
