//target:github.com/droundy/gui/exp/data
package data

type Widget interface {
	iswidget()
}

type Button struct {
	String string
	ClickChan chan struct{}
}
func (w *Button) iswidget() {}
func (w *Button) Clicks() chan struct{} {
	return w.ClickChan
}
func (w *Button) Raw() Widget {
	return w
}

type Text struct {
	String string
	ClickChan chan struct{}
}
func (w *Text) iswidget() {}
func (w *Text) Clicks() chan struct{} {
	return w.ClickChan
}
func (w *Text) Raw() Widget {
	return w
}

type EditText struct {
	String string
	ClickChan chan struct{}
	ChangeChan chan string
}
func (w *EditText) iswidget() {}
func (w *EditText) Changes() chan string {
	return w.ChangeChan
}
func (w *EditText) Clicks() chan struct{} {
	return w.ClickChan
}
func (w *EditText) Raw() Widget {
	return w
}

type TextArea struct {
	String string
	ClickChan chan struct{}
	ChangeChan chan string
}
func (w *TextArea) iswidget() {}
func (w *TextArea) Changes() chan string {
	return w.ChangeChan
}
func (w *TextArea) Clicks() chan struct{} {
	return w.ClickChan
}
func (w *TextArea) Raw() Widget {
	return w
}

type Menu struct {
	Value int
	Options []string
	ChangeChan chan string
}
func (w *Menu) iswidget() {}
func (w *Menu) Changes() chan string {
	return w.ChangeChan
}
func (w *Menu) Raw() Widget {
	return w
}


type Column []Widget
func (w *Column) iswidget() {}

type Table [][]Widget
func (w *Table) iswidget() {}

type Window struct {
	Title string
	Path string
	Contents Widget
}
