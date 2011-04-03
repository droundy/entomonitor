package main

import (
	"http"
	"fmt"
	"strings"
	"strconv"
	"έντομο"
	"github.com/droundy/gui"
	"github.com/droundy/goopt"
)

var port = goopt.Int([]string{"--port"}, 8080, "port on which to serve")

var bug = έντομο.Type("bug")
var todo = έντομο.Type("todo")

func main() {
	http.HandleFunc("/style.css", styleServer)

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
	x.Widget = BugList(x)
	x.OnPath(func() gui.Refresh {
		p := x.GetPath()[1:]
		fmt.Println("My path is actually", p)
		psplit := strings.Split(p, "-", 2)
		if len(psplit) != 2 {
			fmt.Println("Not a particular bug")
			x.Widget = BugList(x)
			return gui.NeedsRefresh
		}
		bnum, err := strconv.Atoi(psplit[1])
		if err != nil || len(psplit[0]) == 0 {
			fmt.Println("Not a particular bug ii")
			x.Widget = BugList(x)
			return gui.NeedsRefresh
		}
		x.Widget = BugPage(x, έντομο.Type(psplit[0]), bnum)
		return gui.NeedsRefresh
	})
	return x
}

func BugList(p gui.HasPath) gui.Widget {
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
		setpath := func() gui.Refresh { p.SetPath("/" + bugname); return p.HandlePath() }
		bid := gui.Button(bugname)
		bid.OnClick(setpath)
		bstatus := gui.Text(status)
		bstatus.OnClick(setpath)
		bdate := gui.Text(cs[0].Date)
		bdate.OnClick(setpath)
		btitle := gui.Text(lines[0])
		btitle.OnClick(setpath)
		bugtable = append(bugtable, []gui.Widget{bid, bstatus, bdate, btitle})
	}
	bugs = append(bugs, gui.Empty(), gui.Empty(), gui.Table(bugtable...))
	return gui.Column(bugs...)
}

func BugPage(p gui.HasPath, btype έντομο.Type, bnum int) gui.Widget {
	bl, err := btype.List()
	if err != nil {
		return gui.Text("Error: " + err.String())
	}
	if bnum >= len(bl) {
		return gui.Text(fmt.Sprint("Error: no such ", btype, " as number ", bnum))
	}
	b := bl[bnum]
	cs, err := b.Comments()
	if err != nil {
		return gui.Text("Error: " + err.String())
	}
	listall := gui.Button("List all bugs")
	listall.OnClick(func() gui.Refresh {
		p.SetPath("/")
		return p.HandlePath()
	})
	bugs := []gui.Widget{
		listall,
	}
	for _, c := range cs {
		bugs = append(bugs, gui.Text(c.Author), gui.Text(c.Date), gui.Text(c.Text))
	}
	return gui.Column(bugs...)
}
