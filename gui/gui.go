package gui

type Text struct {
	String string
}
func (*Text) iswidget() {
}

type Table struct {
	Rows [][]Widget
}
func (*Table) iswidget() {
}

type Button struct {
	Text
}

type Widget interface {
	iswidget()
}
