package main

import (
	"os"
	"fmt"
	"gui/data"
	"web"
)

func main() {

	err := web.Serve(12345, NewWidget)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	} else {
		fmt.Println("Exited successfully!")
	}
}

func NewWidget() *data.Window {
	window := data.Window{ "Class survey", "", nil }
	
	namebox := &data.EditText{}
	namerow := &data.Table {
		[][]data.Widget{
			{ &data.Text{"Name:"}, namebox },
		},
	}
	namebox.HandleChanged = func(old string) (modified data.Widget, refresh bool) {
		window.Title = `Survey of ` + namebox.Text.String
		return
	}
	partnerbox := &data.EditText{}
	partnerrow := &data.Table {
		[][]data.Widget{
			{ &data.Text{"Partner:"}, partnerbox },
		},
	}
	dotoday := &data.TextArea{}
	learntoday := &data.TextArea{}
	workwell := &data.TextArea{}

	button := &data.Button{Text: data.Text{"Submit"}}

	widget := &data.Table{
		[][]data.Widget{
			{ namerow },
			{ partnerrow },
			{ &data.Text{"What did you do today?"} },
			{ dotoday },
			{ &data.Text{"What is one thing you learned today?"} },
			{ learntoday },
			{ &data.Text{"What is one thing that didn't work well today?"} },
			{ workwell },
			{ button },
		},
	}
	window.Widget = widget
	button.HandleClick = func() (modified data.Widget, refresh bool) {
		fmt.Println("Name:", namebox.Text.String)
		fmt.Println("Partner:", partnerbox.Text.String)
		fmt.Println("Did >>>>>>")
		fmt.Println(dotoday.Text.String)
		fmt.Println("Learned >>>>>>")
		fmt.Println(learntoday.Text.String)

		widget.Rows = [][]data.Widget {
			{ &data.Text{ "Thank you, " + namebox.Text.String + "!" } },
		}
		return nil, true
	}
	return &window
}
