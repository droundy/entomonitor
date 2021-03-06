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
  args := []string{"var", "GIT_AUTHOR_IDENT"}
  o, err := exec.Command("git", args...).CombinedOutput()
  if err != nil {
    return err.String()
  }
	lines := bytes.SplitN(o, []byte{'\n'}, 2)
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
		lns := bytes.SplitN(x, []byte{'\n'}, 2)
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
	comment, err := os.OpenFile(fname, os.O_CREATE+os.O_WRONLY+os.O_EXCL, 0777)
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

func (t Type) Create() *Bug {
	var b Bug
	b.Type = t
	b.Id = createName()
	b.Attributes = make(map[string]string)
	b.PendingChanges = make(chan string, 100)
	return &b
}

func (t Type) New(text string) (b *Bug, err os.Error) {
	b = new(Bug)
	b.Type = t
	b.Id = createName()
	b.Attributes = make(map[string]string)
	b.PendingChanges = make(chan string, 100)
	err = WriteComment(".entomon/"+string(t)+"/"+b.Id, text)
	return
}

func (t Type) List() (out []*Bug, err os.Error) {
	d, err := os.Open(".entomon/" + string(t))
	defer d.Close()
	if err != nil {
		return out, nil
	}
	ns, err := d.Readdirnames(-1)
	if err != nil {
		return
	}
	sort.Strings(ns)
	for _, n := range ns {
		if len(n) > ndate+2 {
			out = append(out, &Bug{n, t, nil, nil})
		}
	}
	return
}

func (t Type) AttributeOptions(attr string) []string {
	data, err := ioutil.ReadFile(".entomon/" + string(t) + "/options/" + attr)
	if err != nil {
		return nil // FIXME: should I verify ENOEXIST error?
	}
	as := []string{}
	for _, a := range strings.Split(string(data), "\n") {
		if len(a) > 0 {
			as = append(as, a)
		}
	}
	return as
}

func (t Type) ListAttributes() []string {
	data, err := ioutil.ReadFile(".entomon/" + string(t) + "/attributes")
	if err != nil {
		return nil // FIXME: should I verify ENOEXIST error?
	}
	as := []string{}
	for _, a := range strings.Split(string(data), "\n") {
		if len(a) > 0 {
			as = append(as, a)
		}
	}
	return as
}

func LookupBug(b string) (out *Bug, err os.Error) {
	xs := strings.SplitN(b, "-", 2)
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
	Attributes     map[string]string
	PendingChanges chan string
}

func (b *Bug) String() string {
	return fmt.Sprint(b.Type, "/", b.Id)
}

type Comment struct {
	Author string
	Date   string
	Text   string
}

func (b *Bug) Initialize() {
	if b.Attributes == nil || len(b.Attributes) == 0 {
		b.Attributes = make(map[string]string)
		d, err := os.Open(".entomon/" + string(b.Type) + "/defaults")
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
				//fmt.Println("Setting default value", string(x), "for attribute", a)
				xx := strings.SplitN(string(x), "\n", 2)
				b.Attributes[a] = xx[0]
			}
		}
	}
}

func (b *Bug) stripAttributes(c Comment) Comment {
	b.Initialize()
	t := c.Text
	for {
		lines := strings.SplitN(t, "\n", 2)
		if len(lines) < 2 {
			break
		}
		att := strings.SplitN(lines[0], ": ", 2)
		if len(att) != 2 || len(att[0]) == 0 {
			// We've passed by all the Attribute: lines
			break
		}
		ch := strings.SplitN(att[1], " -> ", 2)
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
		t = lines[1]
	}
	c.Text = t
	return c
}

func (b *Bug) AddComment(t string) os.Error {
	return WriteComment(fmt.Sprint(".entomon/", b.Type, "/", b.Id), t)
}

func (b *Bug) ScheduleChange(s string) {
	b.PendingChanges <- s
}

func (b *Bug) FlushPending() os.Error {
	pendingchanges := []string{}
	notdone := true
	for notdone {
		select {
		case ch := <- b.PendingChanges:
			pendingchanges = append(pendingchanges, ch)
		default:
			notdone = false
		}
	}
	c := strings.Join(pendingchanges, "\n")
	return b.AddComment(c)
}

func (b *Bug) Comments() (out []Comment, err os.Error) {
	d, err := os.Open(".entomon/" + string(b.Type) + "/" + b.Id)
	defer d.Close()
	if err != nil {
		return out, nil
	}
	ns, err := d.Readdirnames(-1)
	d.Close()
	if err != nil {
		return
	}
	sort.Strings(ns)
	for _, n := range ns {
		dateauthor := strings.SplitN(n, "--", 2)
		// Now let's put the date into the current timezone...
		date, err := time.Parse(time.RFC3339, dateauthor[0])
		if err == nil {
			dateauthor[0] = time.SecondsToLocalTime(date.Seconds()).Format(
				"January 2, 2006 3:04PM")
		}
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

func (b *Bug) WriteAttribute(attr, value string) os.Error {
	b.Attributes[attr] = value
	return b.AddComment(attr + ": " + value)
}

func (b *Bug) ScheduleAttribute(attr, value string) {
	b.Attributes[attr] = value
	b.ScheduleChange(attr + ": " + value)
}
