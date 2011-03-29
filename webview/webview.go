package main

import (
	"http"
	"fmt"
	"github.com/droundy/gui"
	"github.com/droundy/goopt"
)

var port = goopt.Int([]string{"--port"}, 8080, "port on which to serve")

// FIXME: There is no way to make the following return something that
// is both a gui.Bool and a gui.HasText.
type labelledcheckbox struct {
	gui.Widget
	gui.String
	gui.Changeable
	gui.Bool
}
func LabelledCheckbox(l string) interface { gui.Widget; gui.String; gui.Changeable; gui.Bool } {
	cb := gui.Checkbox()
	label := gui.Text(l)
	table := gui.Row(cb, label)
	label.OnClick(func() gui.Refresh {
		cb.Toggle()
		return cb.HandleChange()
	})
	out := labelledcheckbox{ table, label, cb, cb }
	return &out
}

type radiobuttons struct {
	gui.Widget
	gui.String
	gui.Changeable
}
	
func RadioButtons(vs... string) interface{ gui.Widget; gui.String; gui.Changeable } {
	var bs []interface{ gui.Changeable; gui.Bool; gui.String }
	var ws []gui.Widget
	for _,v := range vs {
		b := gui.RadioButton(v)
		bs = append(bs, b)
		ws = append(ws, b)
	}
	col := gui.Column(ws...)
	grp := gui.RadioGroup(bs...)
	return &radiobuttons{ col, grp, grp }
}


func main() {
	http.HandleFunc("/style.css", styleServer)

	err := gui.Run(*port,
		gui.Column(
		gui.Text("This will be a bug browser, some day!"),
		))
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

func styleServer(c http.ResponseWriter, req *http.Request) {
	c.SetHeader("Content-Type", "text/css")
	fmt.Fprint(c, `
html {
    margin: 0;
    padding: 0;
}

body {
    margin: 0;
    padding: 0;
    background: #ffffff;
    font-family: arial,helvetica,"sans serif";
    font-size: 12pt;
}
h1 {
font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 16pt;
}
h2 { font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 14pt;
}
p {
font-family: arial,helvetica,"sans serif";
font-size:12pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 12pt;
}
a {
  color: #555599;
}
`)
}
