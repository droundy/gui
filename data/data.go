//target:github.com/droundy/gui/data
package data

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

type Menu struct {
	Value int
	Options []string
	HandleChanged
}
func (*Menu) iswidget() {}
func (w *Menu) Lookup(p string) Widget {
	return leafLookup(w, p)
}
func (w *Menu) Name() string {
	return leafName(w)
}
func (w *Menu) Handle(event Event) (modified Widget, refresh bool) {
	//fmt.Printf("Got handle of %#v in %#v\n", event, w)
	if event.Widget != w.Name() {
		fmt.Println("It isn't me:", w.Name(), event.Widget)
		return
	}
	switch strings.SplitN(event.Event, ":", 2)[0] {
	case "onchange":
		old := w.Options[w.Value]
		newv := strings.SplitN(event.Event, ":", 2)[1]
		if newv != old {
			for i := range w.Options {
				if newv == w.Options[i] {
					w.Value = i
					if w.HandleChanged != nil {
						return w.HandleChanged(old)
					}
					return
				}
			}
			fmt.Println("New value doesn't make sense in Menu.Handle onchange")
		}
	}
	return
}
func (w *Menu) SetString(newv string) {
	if newv != w.Options[w.Value] {
		for i := range w.Options {
			if newv == w.Options[i] {
				w.Value = i
				return
			}
		}
		fmt.Println("New value doesn't make sense in Menu.SetString")
	}
}
func (w *Menu) String() string {
	return w.Options[w.Value]
}

type Text struct {
	String string
	HandleClick
}
func (*Text) iswidget() {}
func (t *Text) Lookup(p string) Widget {
	return leafLookup(t, p)
}
func (t *Text) Name() string {
	return leafName(t)
}
func (w *Text) Handle(event Event) (Widget, bool) {
	fmt.Println("Got Handle even in *gui.Text")
	return nil, false
}

type EditText struct {
	Text
	HandleChanged
}
func (t *EditText) Lookup(p string) Widget {
	return leafLookup(t, p)
}
func (t *EditText) Name() string {
	return leafName(t)
}
func (w *EditText) Handle(event Event) (modified Widget, refresh bool) {
	//fmt.Printf("Got handle of %#v in %#v\n", event, w)
	if event.Widget != w.Name() {
		fmt.Println("It isn't me:", w.Name(), event.Widget)
		return
	}
	switch strings.SplitN(event.Event, ":", 2)[0] {
	case "onchange":
		old := w.Text.String
		w.Text.String = strings.SplitN(event.Event, ":", 2)[1]
		if w.HandleChanged != nil {
			return w.HandleChanged(old)
		}
	}
	return
}

type TextArea struct {
	EditText
}

type Column []Widget
func (*Column) iswidget() {}
func (t *Column) Name() string {
	return leafName(t)
}
func (t *Column) Lookup(p string) Widget {
	if p == t.Name() {
		return t
	}
	if i,rest,ok := t.lookInside(p); ok {
		return (*t)[i].Lookup(rest)
	}
	return nil
}
func (t *Column) lookInside(p string) (i int, rest string, ok bool) {
	s := strings.SplitN(p, "/", 2)
	if len(s) != 2 {
		fmt.Println("Weird bug:  not a nice name")
		return
	}
	i, err := strconv.Atoi(s[0])
	if err != nil || i >= len(*t) {
		fmt.Println("Weird bug:  not a good row")
		return
	}
	return i, s[1], true
}
func (w *Column) Handle(event Event) (modified Widget, refresh bool) {
	if i,rest,ok := w.lookInside(event.Widget); ok {
		event.Widget = rest
		//fmt.Printf("Passing off %#v to %#v\n", event, w[i])
		newi, refresh := (*w)[i].Handle(event)
		if newi != nil {
			(*w)[i] = newi
			fmt.Println("Something changed.")
			return w, refresh
		}
		return nil, refresh
	}
	return
}

type Table [][]Widget
func (*Table) iswidget() {}
func (t *Table) Name() string {
	return leafName(t)
}
func (t *Table) Lookup(p string) Widget {
	if p == t.Name() {
		return t
	}
	if i,j,rest,ok := t.lookInside(p); ok {
		return (*t)[i][j].Lookup(rest)
	}
	return nil
}
func (t *Table) lookInside(p string) (i int, j int, rest string, ok bool) {
	s := strings.SplitN(p, "/", 3)
	if len(s) != 3 {
		fmt.Println("Weird bug:  not a nice name")
		return
	}
	i, err := strconv.Atoi(s[0])
	if err != nil || i >= len(*t) {
		fmt.Println("Weird bug:  not a good row")
		return
	}
	r := (*t)[i]
	j, err = strconv.Atoi(s[1])
	if err != nil || j >= len(r) {
		fmt.Println("Weird bug:  not a good column: ", s[1])
		return
	}
	return i, j, s[2], true
}
func (w *Table) Handle(event Event) (modified Widget, refresh bool) {
	if i,j,rest,ok := w.lookInside(event.Widget); ok {
		event.Widget = rest
		//fmt.Printf("Passing off %#v to %#v\n", event, w[i][j])
		newij, refresh := (*w)[i][j].Handle(event)
		if newij != nil {
			(*w)[i][j] = newij
			fmt.Println("Something changed.")
			return w, refresh
		}
		return nil, refresh
	}
	return
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
func (w *Button) Handle(event Event) (modified Widget, refresh bool) {
	if event.Widget != w.Name() {
		return
	}
	switch event.Event {
	case "onclick":
		if w.HandleClick != nil {
			return w.HandleClick()
		} else {
			fmt.Println("This button doesn't do anything")
		}
	}
	return
}

type Widget interface {
	iswidget()
	Lookup(string) Widget // nil indicates no such widget
	Name() string // this is a programmer-friendly name for the widget
	// We have no static checking of which events are handled. A nil
	// return value from Handle means that this widget wasn't changed.
	// The bool argument indicates whether we need to redraw everything,
	// e.g. if Handle changed some other widget.
	Handle(event Event) (modified Widget, refresh bool)
}

type Event struct {
	Widget string
	Event string
}

type HandleClick func() (modified Widget, refresh bool)
type HandleChanged func(old string) (modified Widget, refresh bool)
