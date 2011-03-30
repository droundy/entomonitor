//target:έντομο
package έντομο

import (
	"fmt"
	"os"
	"exec"
	"rand"
	crand "crypto/rand"
	"time"
	"bytes"
	"strings"
	"path/filepath"
	"io/ioutil"
	"github.com/droundy/goopt"
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

func getDefaultAuthor() string {
	args := []string{"git", "var", "GIT_AUTHOR_IDENT"}
	git, err := exec.LookPath("git")
	if err != nil {
		return err.String()
	}
	pid, err := exec.Run(git, args, nil, "", exec.PassThrough, exec.Pipe, exec.PassThrough)
	if err != nil {
		return err.String()
	}
	o, err := ioutil.ReadAll(pid.Stdout)
	if err != nil {
		return err.String()
	}
	_, err = pid.Wait(0) // could have been os.WRUSAGE
	if err != nil {
		return err.String()
	}
	const ndate = 17
	lines := bytes.Split(o, []byte{'\n'}, 2)
	if len(lines[0]) > ndate {
		lines[0] = lines[0][:len(lines[0])-ndate]
	}
	return string(lines[0])
}

var Author = goopt.String([]string{"--author"}, getDefaultAuthor(), "author of this change")

func createName() string {
	*Author = strings.Replace(strings.Replace(strings.Replace(*Author, "\n", " ", -1), "/", "-", -1), "\\", "-", -1)
	return time.SecondsToUTC(time.Seconds()).Format(time.RFC3339) + "--" + *Author
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

func WriteComment(dname, text string) os.Error {
	id := createName()
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
	_, err = fmt.Fprintln(comment, text)
	return err
}

// A Type represents a class of items.  Typical Types would be "bug",
// "issue" or "todo".  A Type would not normally be useful for
// distinguishing between "wishlist", "feature-request" or "real"
// bugs, which should be managed by an Attribute, so you can easily
// change a bug from one to another (since bug reporters commonly
// misclassify bugs!).
type Type string

// έντομο

func (t Type) New(text string) (b Bug, err os.Error) {
	b.Type = t
	b.Id = createName()
	b.Attributes = make(map[string]string)
	err = WriteComment(".entomon/" + string(t) + "/"+b.Id, text)
	return
}

type Bug struct {
	Id string
	Type
	Attributes map[string]string
}

func (b *Bug) String() string {
	return fmt.Sprint(b.Type, "-", b.Id)
}
