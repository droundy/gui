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
func (w *Text) Handle(event Event) Widget {
	return nil
}

type EditText struct {
	Text
	HandleChanged
}
func (w *EditText) Handle(event Event) Widget {
	if event.Widget != w.Name() {
		return nil
	}
	switch strings.SplitN(event.Event, ":", 2)[0] {
	case "onchange":
		old := w.Text.String
		w.Text.String = strings.SplitN(event.Event, ":", 2)[1]
		if w.HandleChanged != nil {
			return w.HandleChanged(old)
		}
	}
	return nil
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
	if i,j,rest,ok := t.lookInside(p); ok {
		return t.Rows[i][j].Lookup(rest)
	}
	return nil
}
func (t *Table) lookInside(p string) (i int, j int, rest string, ok bool) {
	s := strings.SplitN(p, "/", 3)
	if len(s) != 3 {
		return
	}
	i, err := strconv.Atoi(s[0])
	if err != nil || i >= len(t.Rows) {
		return
	}
	r := t.Rows[i]
	j, err = strconv.Atoi(s[1])
	if err != nil || j >= len(r) {
		return
	}
	return i, j, s[2], true
}
func (w *Table) Handle(event Event) Widget {
	if i,j,rest,ok := w.lookInside(event.Widget); ok {
		event.Widget = rest
		newij := w.Rows[i][j].Handle(event)
		if newij != nil {
			w.Rows[i][j] = newij
		}
		return w
	}
	return nil
}

type Button struct {
	Text
	HandleClick
}
func (b *Button) Lookup(p string) Widget {
	return leafLookup(b, p)
}
func (b *Button) Name() string {
	return "Button-" + b.Text.String
}
func (w *Button) Handle(event Event) Widget {
	if event.Widget != w.Name() {
		return nil
	}
	switch event.Event {
	case "onclick":
		if w.HandleClick != nil {
			return w.HandleClick()
		} else {
			fmt.Println("This button doesn't do anything")
		}
	}
	return nil
}

type Widget interface {
	iswidget()
	Lookup(string) Widget // nil indicates no such widget
	Name() string // this is a programmer-friendly name for the widget
	// We have no static checking of which events are handled. A nil
	// return value from Handle means that nothing was changed and we
	// don't need to redraw the widget.
	Handle(event Event) Widget
}

type Event struct {
	Widget string
	Event string
}

type HandleClick func() Widget
type HandleChanged func(old string) Widget
