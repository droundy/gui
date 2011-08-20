package gui

import (
	"fmt"
	"strings"
	"strconv"
)

func leafLookup(w Widget, p string) Widget {
	if p == w.Name() {
		return w
	}
	return nil
}
func leafName(w interface{}) string {
	return fmt.Sprintf("%T", w)
}

type Text struct {
	String string
}
func (*Text) iswidget() {}
func (t *Text) Lookup(p string) Widget {
	return leafLookup(t, p)
}
func (t *Text) Name() string {
	return leafName(t)
}

type Table struct {
	Rows [][]Widget
}
func (*Table) iswidget() {}
func (t *Table) Name() string {
	return leafName(t)
}
func (t *Table) Lookup(p string) Widget {
	if p == t.Name() {
		return t
	}
	s := strings.SplitN(p, "/", 3)
	if len(s) != 3 {
		return nil
	}
	i, err := strconv.Atoi(s[0])
	if err != nil || i >= len(t.Rows) {
		return nil
	}
	r := t.Rows[i]
	j, err := strconv.Atoi(s[1])
	if err != nil || j >= len(r) {
		return nil
	}
	return r[j]
}

type Button struct {
	Text
}
func (b *Button) Lookup(p string) Widget {
	return leafLookup(b, p)
}
func (b *Button) Name() string {
	return "Button-" + b.Text.String
}

type Widget interface {
	iswidget()
	Lookup(string) Widget // nil indicates no such widget
	Name() string // this is a programmer-friendly name for the widget
}
