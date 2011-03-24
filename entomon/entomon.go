//target:entomon
package entomon

import (
	"os"
	"bytes"
	"path/filepath"
	"io/ioutil"
)

func isEntomonHere() bool {
	fi, err := os.Stat(".entomon")
	return err == nil && fi.IsDirectory()
}

func findEntomon() os.Error {
	origd, err := os.Getwd()
	if err != nil {
		// If we can't read our working directory, let's just fail!
		return err
	}
	oldd := origd
	for !isEntomonHere() {
		err = os.Chdir("..")
		if err != nil {
			// If we can't cd .. then we'll just use the original directory.
			goto giveup
		}
		newd, err := os.Getwd()
		if err != nil || newd == oldd {
			// Either something weird happened or we're at the root
			// directory.  In either case, we'll just go with the original
			// directory.
			goto giveup
		}
		oldd = newd
	}
	return nil
giveup:
	// If nothing else works we'll just use the original directory.
	err = os.Chdir(origd)
	if err != nil {
		return err
	}
	return os.MkdirAll(".entomon", 0777)
}

func ProjectName() (name string, err os.Error) {
	err = findEntomon()
	x, err := ioutil.ReadFile(".entomon/ProjectName")
	if err == nil {
		lns := bytes.Split(x, []byte{'\n'}, 2)
		return string(lns[0]), nil
	}
	origd, err := os.Getwd()
	return filepath.Base(origd), err
}
