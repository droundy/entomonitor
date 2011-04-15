package main

import (
	"http"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
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

func Header(page string, p gui.PathHandler) gui.Widget {
	elems := []gui.Widget{}
	list := gui.Button("Bug list")
	list.OnClick(func() gui.Refresh {
		p.SetPath("/")
		return gui.NeedsRefresh
	})
	elems = append(elems, list)
	newbug := gui.Button("Report new bug")
	newbug.OnClick(func() gui.Refresh {
		p.SetPath("/new")
		return gui.NeedsRefresh
	})
	about := gui.Button("About")
	about.OnClick(func() gui.Refresh {
		p.SetPath("/about")
		return gui.NeedsRefresh
	})
	elems = append(elems, newbug, about)
	return gui.Row(elems...)
}

func Page() gui.Widget {
	x := gui.MakePathHandler(nil)
	x.SetWidget(gui.Paragraphs(Header("/", x), BugList(x)))
	x.OnPath(func() gui.Refresh {
		p := x.GetPath()[1:]
		fmt.Println("My path is actually", p)
		psplit := strings.Split(p, "-", 2)
		if len(psplit) != 2 {
			switch p {
			case "/", "":
				x.SetWidget(gui.Paragraphs(Header(p, x), BugList(x)))
			case "new":
				x.SetWidget(gui.Paragraphs(Header(p, x), NewBug(x, έντομο.Type("bug"))))
			default:
				if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".txt"); err == nil {
					x.SetWidget(gui.Paragraphs(Header(p, x), gui.Text(string(page))))
				} else if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".md"); err == nil {
					x.SetWidget(gui.Paragraphs(Header(p, x), gui.Text(string(page))))
				} else if page, err := ioutil.ReadFile(".entomon/Static/" + p + ".html"); err == nil {
					x.SetWidget(gui.Paragraphs(Header(p, x),
						gui.Text("This html shouldn't be escaped"), gui.Text(string(page))))
				} else {
					x.SetWidget(gui.Paragraphs(Header(p, x), gui.Text("I don't understand: "+p)))
				}
			}
			return gui.NeedsRefresh
		}
		bnum, err := strconv.Atoi(psplit[1])
		if err != nil || len(psplit[0]) == 0 {
			x.SetWidget(gui.Paragraphs(Header(p, x), BugList(x)))
			return gui.NeedsRefresh
		}
		x.SetWidget(gui.Paragraphs(Header(p, x), BugPage(x, έντομο.Type(psplit[0]), bnum)))
		return gui.NeedsRefresh
	})
	return x
}

type WhenToWrite bool

const (
	WriteNow   WhenToWrite = true
	WriteLater WhenToWrite = false
)

func AttributeChooser(b *έντομο.Bug, attr string, imm WhenToWrite) interface {
	gui.Widget
	gui.String
	gui.Changeable
} {
	opts := b.Type.AttributeOptions(attr)
	b.Initialize()
	b.Comments()
	if len(opts) > 1 {
		menu := gui.Menu(opts...)
		fmt.Println("looking at", attr)
		menu.SetString(b.Attributes[attr])
		menu.OnChange(func() gui.Refresh {
			if imm == WriteNow {
				b.WriteAttribute(attr, menu.GetString())
			} else {
				b.PendingChanges = append(b.PendingChanges, attr+":"+menu.GetString())
			}
			return gui.NeedsRefresh
		})
		return menu
	}
	edit := gui.EditText(b.Attributes[attr])
	edit.OnChange(func() gui.Refresh {
		if imm == WriteNow {
			b.WriteAttribute(attr, edit.GetString())
		} else {
			b.PendingChanges = append(b.PendingChanges, attr+":"+edit.GetString())
		}
		return gui.NeedsRefresh
	})
	return edit
}

func BugList(p gui.PathHandler) gui.Widget {
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
		setpath := func() gui.Refresh { return p.SetPath("/" + bugname) }
		bid := gui.Button(bugname)
		bid.OnClick(setpath)
		bstatus := AttributeChooser(b, "status", WriteNow)
		bdate := gui.Text(cs[0].Date)
		bdate.OnClick(setpath)
		//btitle := AttributeChooser(b, "title")
		btitle := gui.Text(b.Attributes["title"])
		bugtable = append(bugtable, []gui.Widget{bid, bstatus, bdate, btitle})
	}
	bugs = append(bugs, gui.Empty(), gui.Empty(), gui.Table(bugtable...))
	return gui.Column(bugs...)
}

func BugPage(p gui.PathHandler, btype έντομο.Type, bnum int) gui.Widget {
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

func NewBug(p gui.PathHandler, btype έντομο.Type) gui.Widget {
	attrs := []string{"title", "status"}
	fields := []gui.Widget{}
	b := btype.Create()
	for _, attr := range attrs {
		fields = append(fields, gui.Row(gui.Text(attr+":"), AttributeChooser(b, attr, WriteLater)))
	}
	maintext := gui.EditText("")
	fields = append(fields, gui.Row(gui.Text("Comment:"), maintext))
	submit := gui.Button("Submit bug")
	submit.OnClick(func() gui.Refresh {
		b.PendingChanges = append(b.PendingChanges, maintext.GetString())
		b.FlushPending()
		p.SetPath("/")
		return gui.NeedsRefresh
	})
	fields = append(fields, submit)
	return gui.Column(fields...)
}
