package main

import (
	"http"
	"fmt"
	"strings"
	"έντομο"
	"github.com/droundy/gui"
	"github.com/droundy/goopt"
)

var port = goopt.Int([]string{"--port"}, 8080, "port on which to serve")

var bug = έντομο.Type("bug")
var todo = έντομο.Type("todo")

func main() {
	http.HandleFunc("/style.css", styleServer)

	//gui.HandleSeparate("/bug/", BugPage)
	err := gui.RunSeparate(*port, Page)
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

type PageType struct {
	gui.Widget
	gui.PathHandler
}

func Page() gui.Widget {
	var x = new(PageType)
	x.Widget = BugList()
	x.OnPath(func() gui.Refresh {
		p := x.GetPath()
		if len(p) > 5 && p[:5] == "/bug-" {
			x.Widget = BugPage(p)
		}
		fmt.Println("My path is actually", p)
		return gui.NeedsRefresh
	})
	return x
}

func BugList() gui.Widget {
	bugs := []gui.Widget{gui.Text("This will be a bug browser, some day!")}
	bl, err := bug.List()
	if err != nil {
		panic("bug.List: " + err.String())
	}
	bugtable := [][]gui.Widget{{
		gui.Text("id"), gui.Text("status"), gui.Text("date"), gui.Text("bug"),
	}}
	for bnum, b := range bl {
		b.Comments() // to get attributes
		bugname := fmt.Sprint(bug, "-", bnum)
		// bugs = append(bugs, gui.Text(""), gui.Text(bugname))
		cs, err := b.Comments()
		if err != nil {
			continue
		}
		// for _, c := range cs {
		// 	bugs = append(bugs, gui.Text(c.Author), gui.Text(c.Date), gui.Text(c.Text))
		// }
		lines := strings.Split(cs[0].Text, "\n", 2)
		status, _ := b.Attributes["status"]
		bugtable = append(bugtable, []gui.Widget{
			gui.Button(bugname), gui.Text(status), gui.Text(cs[0].Date), gui.Text(lines[0])})
	}
	bugs = append(bugs, gui.Empty(), gui.Empty(), gui.Table(bugtable...))
	return gui.Column(bugs...)
}

func BugPage(b string) gui.Widget {
	bugs := []gui.Widget{gui.Text("This will show a particular bug, some day! " + b)}
	return gui.Column(bugs...)
}
