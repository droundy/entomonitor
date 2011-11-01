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

	if err = exec.Command("emacs", fname).Run(); err != nil {
		return
	}
	bs, err := ioutil.ReadFile(fname)
	return string(bs), err
}
