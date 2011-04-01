package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"exec"
)

func Edit(start_text string) (out string, err os.Error) {
	f, err := ioutil.TempFile("", "ento-")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(f, start_text)
	if err != nil {
		return
	}
	fname := f.Name()
	defer os.Remove(fname)
	f.Close()
	editor, err := exec.LookPath("emacs")
	if err != nil {
		return
	}
	pid, err := exec.Run(editor, []string{editor, fname}, nil, "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		return
	}
	ws, err := pid.Wait(0)
	if err != nil {
		return
	}
	if ws.ExitStatus() != 0 {
		err = os.NewError(fmt.Sprintf("%s %s exited with '%v'", editor, fname, ws.ExitStatus()))
	}
	bs, err := ioutil.ReadFile(fname)
	return string(bs), err
}
