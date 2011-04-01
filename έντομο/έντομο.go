//target:έντομο
package έντομο

import (
	"fmt"
	"os"
	"exec"
	"time"
	"bytes"
	"strings"
	"sort"
	"path/filepath"
	"io/ioutil"
	"github.com/droundy/goopt"
)

const ndate = 17

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
	err = WriteComment(".entomon/"+string(t)+"/"+b.Id, text)
	return
}

func (t Type) List() (out []Bug, err os.Error) {
	d, err := os.Open(".entomon/"+string(t), os.O_RDONLY, 0)
	defer d.Close()
	if err != nil {
		return out, nil
	}
	ns, err := d.Readdirnames(-1)
	if err != nil {
		return
	}
	sort.SortStrings(ns)
	for _, n := range ns {
		if len(n) > ndate+2 {
			out = append(out, Bug{n, t, nil})
		}
	}
	return
}

func LookupBug(b string) (out Bug, err os.Error) {
	xs := strings.Split(b, "-", 2)
	t := Type(xs[0])
	if len(xs) < 2 {
		err = os.NewError("Should have a '-' in " + b)
		return
	}
	num := 0
	_, err = fmt.Sscan(xs[1], &num)
	if err != nil {
		return
	}
	bs, err := t.List()
	if err != nil {
		return
	}
	if num >= len(bs) || num < 0 {
		err = os.NewError(fmt.Sprint("Num out of range in ", b, " num is ", num, " len is ", len(bs), " type is ", t))
		return
	}
	return bs[num], err
}

type Bug struct {
	Id string
	Type
	Attributes map[string]string
}

func (b *Bug) String() string {
	return fmt.Sprint(b.Type, "/", b.Id)
}

type Comment struct {
	Author string
	Date   string
	Text   string
}

func (b *Bug) stripAttributes(c Comment) Comment {
	if b.Attributes == nil {
		b.Attributes = make(map[string]string)
		d, err := os.Open(".entomon/"+string(b.Type)+"/defaults", os.O_RDONLY, 0)
		as := []string{}
		if err == nil {
			as, _ = d.Readdirnames(-1)
			// No point checking error, since we'll want to assume no
			// defaults in that case, which is what we woulld have anyhow.
			d.Close()
		}
		for _, a := range as {
			x, err := ioutil.ReadFile(".entomon/" + string(b.Type) + "/defaults/" + a)
			if err == nil {
				xx := strings.Split(string(x), "\n", 2)
				b.Attributes[a] = xx[0]
			}
		}
	}
	firstl := []string{}
	t := c.Text
	for {
		lines := strings.Split(t, "\n", 2)
		if len(lines) < 2 {
			break
		}
		att := strings.Split(lines[0], ": ", 2)
		if len(att) != 2 {
			firstl = append(firstl, lines[0])
			t = lines[1]
			continue
		}
		ch := strings.Split(att[1], " -> ", 2)
		if len(att[0]) == 0 {
			break // This isn't an Attribute: line
		} else {
			if len(ch) == 2 {
				// It is a change thing like Foo: baz -> bar
				old, ok := b.Attributes[att[0]]
				if ok && old == ch[0] {
					b.Attributes[att[0]] = ch[1]
				} else {
					fmt.Println("ch bad", ch)
				}
			} else {
				b.Attributes[att[0]] = att[1]
			}
		}
		t = lines[1]
	}
	firstl = append(firstl, t)
	c.Text = strings.Join(firstl, "\n")
	return c
}

func (b *Bug) AddComment(t string) os.Error {
	return WriteComment(fmt.Sprint(".entomon/", b.Type, "/", b.Id), t)
}

func (b *Bug) Comments() (out []Comment, err os.Error) {
	d, err := os.Open(".entomon/"+string(b.Type)+"/"+b.Id, os.O_RDONLY, 0)
	defer d.Close()
	if err != nil {
		return out, nil
	}
	ns, err := d.Readdirnames(-1)
	d.Close()
	if err != nil {
		return
	}
	sort.SortStrings(ns)
	for _, n := range ns {
		dateauthor := strings.Split(n, "--", 2)
		if len(dateauthor) == 2 {
			c := Comment{dateauthor[1], dateauthor[0], ""}
			x, err := ioutil.ReadFile(".entomon/" + b.String() + "/" + n)
			if err != nil {
				return out, err
			}
			c.Text = string(x)
			c = b.stripAttributes(c)
			out = append(out, c)
		}
	}
	return
}
