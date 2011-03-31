//target:entomonitor
package main

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"έντομο"
	"github.com/droundy/goopt"
)

var action = goopt.Alternatives([]string{"-A", "--action"},
	[]string{"help", "new-issue", "comment"}, "select the action to be performed")
var message = goopt.String([]string{"-m", "--message"}, "", "short message")
var bugid = goopt.String([]string{"-b", "--bug"}, "", "bug ID")

func dieOn(err os.Error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var bug = έντομο.Type("bug")

func main() {
	goopt.Parse(func() []string { return nil })
	pname, err := έντομο.ProjectName()
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
		_, err := bug.New(*message)
		dieOn(err)
	case "comment":
		if *bugid == "" {
			fmt.Print("Which bug? ")
			inp, e := bufio.NewReaderSize(os.Stdin, 1)
			dieOn(e)
			*bugid, err = inp.ReadString('\n')
			if len(*bugid) > 0 {
				*bugid = (*bugid)[:len(*bugid)-1]
			}
		}
		b, err := έντομο.LookupBug(*bugid)
		dieOn(err)
		fmt.Print("What do you want to say? ")
		bugtext, _ := ioutil.ReadAll(os.Stdin)
		b.AddComment(string(bugtext))
	default:
		fmt.Println("I should do", *action)
		os.Exit(1)
	}
}
