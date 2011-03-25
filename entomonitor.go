//target:entomonitor
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"entomon"
	"github.com/droundy/goopt"
)

var action = goopt.Alternatives([]string{"-A", "--action"},
	[]string{"help", "new-issue"}, "select the action to be performed")
var message = goopt.String([]string{"-m", "--message"}, "", "short message")

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
		if *message == "" {
			fmt.Print("What is the problem? ")
			bugtext, _ := ioutil.ReadAll(os.Stdin)
			fmt.Println("Done here")
			*message = string(bugtext)
		}
		dieOn(entomon.NewIssue(*message))
	default:
		fmt.Println("I should do", *action)
		os.Exit(1)
	}
}
