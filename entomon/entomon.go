//target:entomon
package entomon

import (
	"fmt"
	"os"
	"rand"
	crand "crypto/rand"
	"time"
	"bytes"
	"strings"
	"path/filepath"
	"io/ioutil"
)

func init() {
	// Here we seed the random number generator.  I take the lazy
	// expensive route and seed it with a cryptographically strong seed.
	// It'd be better to use the time and maybe something else local
	// like hostname and/or user name.  But that'd be more coding.
	b := make([]byte, 8)
	crand.Read(b) // ignore output, since there's nothing we can do here
	var seed int64
	for _, x := range b {
		seed = int64(x) + 256*seed
	}
	rand.Seed(seed)
}

var (
	vowels = []string{"a", "a", "ae", "au",
		"e", "e", "e", "ea", "ee",
		"i", "io", "iu", "ie", "iou",
		"o", "oo", "ou",
		"u", "uou"}
	consonants = []string{"b", "c", "ch", "d", "f", "g", "h", "j", "k", "l", "m", "n",
		"p", "qu",
		"r", "rd", "rf", "rm", "rn", "rk", "rl", "rj", "rp", "rs", "rt", "rsh", "rth", "rv",
		"s", "sh", "st", "sk", "sch",
		"t", "th", "tch",
		"v", "w", "x", "y", "z", "",
	}
	postconsonants = []string{"r", "l", "", "", "", "", "", "", "-", "-"}
)

func randString() string {
	out := ""
	for len(out) < 30 {
		out = out + consonants[rand.Intn(len(consonants))] +
			postconsonants[rand.Intn(len(postconsonants))] +
			vowels[rand.Intn(len(vowels))]
	}
	return string(out[:30])
}

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

func WriteComment(dname, author, text string) os.Error {
	id := randString()
	err := os.MkdirAll(dname, 0777)
	if err != nil {
		return err
	}
	fname := dname + "/" + id
	comment, err := os.Open(fname, os.O_CREAT+os.O_WRONLY+os.O_EXCL, 0777)
	if err != nil {
		return err
	}
	defer comment.Close()
	now := time.Seconds()
	_, err = fmt.Fprintln(comment, time.SecondsToUTC(now).Format(time.RFC3339))
	if err != nil {
		goto cleanup
	}
	_, err = fmt.Fprintln(comment, author)
	if err != nil {
		goto cleanup
	}
	_, err = fmt.Fprintln(comment, text)
	return err
cleanup:
	comment.Close()
	os.Remove(fname) // We should try to clean up...
	return err
}

func NewIssue(author, text string) os.Error {
	if strings.IndexRune(author, '\n') != -1 {
		return os.NewError("NewIssue: Invalid character '\\n' in author")
	}
	id := randString()
	err := WriteComment(".entomon/issue/"+id, author, text)
	if err != nil {
		return err
	}
	fmt.Println("Created issue", id)
	return nil
}
