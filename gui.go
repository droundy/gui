//target:github.com/droundy/gui
package gui

import "github.com/droundy/gui/data"

// This class defines a high-level API for creating gui programs.  The
// core of this API is the Widget interface.  A number of functions
// are provided which create Widgets, but you can also create your own
// Widgets if you'd like.

// A Widget must support just two methods.  The first is a "Raw"
// method, which provides a representation for the widget in terms of
// concrete data structures, which is obvious.
//
// The second is an accessor to a chan that can be used to replace
// this widget.  For simple widgets, this is a passive accessor for a
// chan that is never used by the widget itself, but is instead
// written to by anyone who wishes to replace this widget, and read
// from by the owner of this widget---which is either its containing
// widget, or the back-end itself, if this happens to be a top-level
// widget.  So only "container" widgets (that can hold other widgets)
// need to do anything tricky with their Updaters.
type Widget interface {
	Raw() data.Widget
	// You can replace this widget by writing to this chan.  Its owner
	// had by golly better be prepared to read from this chan.
	Updater() chan Widget
}

type Updateable chan Widget

func (w *Updateable) Updater() chan Widget {
	return chan Widget(*w)
}

type Clickable interface {
	Clicks() chan struct{}
}
type Changeable interface {
	Changes() chan string
}

type button struct {
	data.Button
	Updateable
}

func Button(name string) interface {
	Widget
	Clickable
} {
	return &button{
		data.Button{name, make(chan struct{})},
		make(Updateable),
	}
}

type text struct {
	data.Text
	Updateable
}

func Text(t string) interface {
	Widget
	Clickable
} {
	return &text{
		data.Text{t, make(chan struct{})},
		make(Updateable),
	}
}

type edittext struct {
	data.EditText
	Updateable
}

func EditText(t string) interface {
	Widget
	Clickable
	Changeable
} {
	return &edittext{
		data.EditText{t, make(chan struct{}), make(chan string)},
		make(Updateable),
	}
}

type textarea struct {
	data.TextArea
	Updateable
}

func TextArea(t string) interface {
	Widget
	Clickable
	Changeable
} {
	return &textarea{
		data.TextArea{t, make(chan struct{}), make(chan string)},
		make(Updateable),
	}
}

type menu struct {
	data.Menu
	Updateable
}
func Menu(value int, options ...string) interface { Widget; Changeable } {
	if value < 0 || value > len(options) {
		panic("value out of range")
	}
	return &menu{
		data.Menu{value, options, make(chan string)},
		make(Updateable),
	}
}

type column struct {
	elems []Widget
	Updateable
}

func (w *column) Raw() data.Widget {
	dw := make(data.Column, len(w.elems))
	for i, sw := range w.elems {
		dw[i] = sw.Raw()
	}
	return &dw
}
func Column(es ...Widget) interface{ Widget } {
	setme := make(Updateable);
	replacements := make(chan struct { int; Widget })
	for i,w := range es {
		thisi := i
		thisw := w
		go func() {
			for {
				x := <-thisw.Updater()
				replacements <- struct {
					int
					Widget
				}{thisi, x}
			}
		}()
	}
	go func() {
		for {
			r := <-replacements
			wnew := make([]Widget, len(es))
			copy(wnew, es)
			es = wnew
			wnew[r.int] = r.Widget
			setme <- &column{es, setme}
		}
	}()
	return &column{es, setme}
}

func Row(elems ...Widget) interface {
	Widget
} {
	return Table([][]Widget{elems})
}

type table struct {
	elems [][]Widget
	Updateable
}

func (w *table) Raw() data.Widget {
	dw := make(data.Table, len(w.elems))
	for i, r := range w.elems {
		dw[i] = make([]data.Widget, len(r))
		for j, entry := range r {
			dw[i][j] = entry.Raw()
		}
	}
	return &dw
}
func Table(es [][]Widget) interface {
	Widget
} {
	setme := make(Updateable)
	replacements := make(chan struct {
		i, j int
		Widget
	})
	for i, r := range es {
		for j, w := range r {
			thisw := w
			thisi := i
			thisj := j
			go func() {
				for {
					x := <-thisw.Updater()
					replacements <- struct {
						i, j int
						Widget
					}{thisi, thisj, x}
				}
			}()
		}
	}
	go func() {
		for {
			r := <-replacements
			wnew := make([][]Widget, len(es))
			copy(wnew, es)
			es = wnew // so next change will preserve this one
			wnew[r.i][r.j] = r.Widget
			setme <- &table{es, setme}
		}
	}()
	return &table{es, setme}
}

type Window struct {
	Title    string
	Path     string
	Contents Widget
}
