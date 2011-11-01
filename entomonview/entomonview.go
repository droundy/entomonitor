package main

import (
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"έντομο"
	"github.com/droundy/gui"
	"github.com/droundy/gui/web"
	"github.com/droundy/goopt"
)

var port = goopt.Int([]string{"--port"}, 8080, "port on which to serve")

var bug = έντομο.Type("bug")
var todo = έντομο.Type("todo")

func main() {
	err := web.Serve(*port, Page)
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

func Header(page string, p chan<- string) gui.Widget {
	list := gui.Button("Bug list")
	newbug := gui.Button("Report new bug")
	about := gui.Button("About")
	go func() {
		for {
			select {
			case _ = <- list.Clicks():
				p <- "/"
			case _ = <- newbug.Clicks():
				p <- "/new"
			case _ = <- about.Clicks():
				p <- "/about"
			}
		}
	}()
	return gui.Row(list, newbug, about)
}

func Page() gui.Window {
	paths := make(chan string)
	x := gui.Column(Header("/", paths), BugList(paths))
	go func() {
		for {
			p := <- paths
			p = p[1:]
			fmt.Println("My path is actually", p)
			psplit := strings.SplitN(p, "-", 2)
			newp := x
			if len(psplit) != 2 {
				switch p {
				case "/", "":
					newp = gui.Column(Header(p, paths), BugList(paths))
				case "new":
					newp = gui.Column(Header(p, paths), NewBug(paths, έντομο.Type("bug")))
				default:
					if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".txt"); err==nil {
						newp = gui.Column(Header(p, paths), gui.Text(string(page)))
					} else if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".md"); err==nil {
						newp = gui.Column(Header(p, paths), gui.Text(string(page)))
					} else if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".html"); err==nil {
						newp = gui.Column(Header(p, paths),
							gui.Text("This html shouldn't be escaped"), gui.Text(string(page)))
					} else {
						newp = gui.Column(Header(p, paths), gui.Text("I don't understand: "+p))
					}
				}
			} else {
				bnum, err := strconv.Atoi(psplit[1])
				if err != nil || len(psplit[0]) == 0 {
					newp = gui.Column(Header(p, paths), BugList(paths))
				} else {
					newp = gui.Column(Header(p, paths), BugPage(paths, έντομο.Type(psplit[0]), bnum))
				}
			}
			x.Updater() <- newp
			x = newp
		}
	}()
	return gui.Window{"Title", "/", x}
}

type WhenToWrite bool

const (
	WriteNow   WhenToWrite = true
	WriteLater WhenToWrite = false
)

func AttributeChooser(b *έντομο.Bug, attr string, imm WhenToWrite) interface {
	gui.Widget
	gui.Changeable
} {
	opts := b.Type.AttributeOptions(attr)
	b.Initialize()
	b.Comments()
	var chooser interface {
		gui.Widget
		gui.Changeable
	}
	if len(opts) > 1 {
		val := 0
		for i,o := range opts {
			if o == b.Attributes[attr] {
				val = i
			}
		}
		chooser = gui.Menu(val, opts)
		fmt.Println("looking at", attr)
	} else {
		chooser = gui.EditText(b.Attributes[attr])
	}
	go func() {
		for {
			newvalue := <- chooser.Changes()
			if imm == WriteNow {
				b.WriteAttribute(attr, newvalue)
			} else {
				b.PendingChanges = append(b.PendingChanges, attr+":"+newvalue)
			}
		}
	}()
	return chooser
}

func BugList(p chan string) gui.Widget {
	bugs := []gui.Widget{}
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
		if err != nil || len(cs) == 0 {
			continue
		}
		// for _, c := range cs {
		// 	bugs = append(bugs, gui.Text(c.Author), gui.Text(c.Date), gui.Text(c.Text))
		// }
		bid := gui.Button(bugname)
		bstatus := AttributeChooser(b, "status", WriteNow)
		bdate := gui.Text(cs[0].Date)
		//btitle := AttributeChooser(b, "title")
		btitle := gui.Text(b.Attributes["title"])
		go func() {
			for {
				select {
				case _ = <- bid.Clicks():
					p <- "/" + bugname
				case _ = <- bdate.Clicks():
					p <- "/" + bugname
				}
			}
		}()
		bugtable = append(bugtable, []gui.Widget{bid, bstatus, bdate, btitle})
	}
	bugs = append(bugs, gui.Text(""), gui.Text(""), gui.Table(bugtable))
	return gui.Column(bugs...)
}

func BugPage(p chan string, btype έντομο.Type, bnum int) gui.Widget {
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
	bugs := []gui.Widget{}
	for attr := range b.Attributes {
		bugs = append(bugs, gui.Row(gui.Text(attr+":"), AttributeChooser(b, attr, WriteNow)))
	}
	for _, c := range cs {
		bugs = append(bugs, gui.Text(c.Author), gui.Text(c.Date), gui.Text(c.Text))
	}
	return gui.Column(bugs...)
}

func NewBug(p chan string, btype έντομο.Type) gui.Widget {
	attrs := []string{"title", "status"}
	fields := []gui.Widget{}
	b := btype.Create()
	for _, attr := range attrs {
		fields = append(fields, gui.Row(gui.Text(attr+":"), AttributeChooser(b, attr, WriteLater)))
	}
	maintext := gui.TextArea("")
	fields = append(fields, gui.Text("Comment:"))
	fields = append(fields, maintext)
	submit := gui.Button("Submit bug")
	go func() {
		mainstr := ""
		for {
			select {
			case mainstr = <- maintext.Changes():
				// Nothing to do...
			case _ = <- submit.Clicks():
				b.PendingChanges = append(b.PendingChanges, mainstr)
				b.FlushPending()
				p <- "/"
				return
			}
		}
	}()
	fields = append(fields, submit)
	return gui.Column(fields...)
}
