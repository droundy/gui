package main

import (
	"os"
	"fmt"
	"gui"
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

func NewWidget() gui.Widget {
	namebox := &gui.EditText{}
	namerow := &gui.Table {
		[][]gui.Widget{
			{ &gui.Text{"Name:"}, namebox },
		},
	}
	partnerbox := &gui.EditText{}
	partnerrow := &gui.Table {
		[][]gui.Widget{
			{ &gui.Text{"Partner:"}, partnerbox },
		},
	}
	dotoday := &gui.TextArea{}
	learntoday := &gui.TextArea{}

	button := &gui.Button{Text: gui.Text{"Submit"}}

	widget := &gui.Table{
		[][]gui.Widget{
			{ namerow },
			{ partnerrow },
			{ &gui.Text{"What did you do today?"} },
			{ dotoday },
			{ &gui.Text{"What did you learn today?"} },
			{ learntoday },
			{ button },
		},
	}
	button.HandleClick = func() (modified gui.Widget, refresh bool) {
		fmt.Println("Name:", namebox.Text.String)
		fmt.Println("Partner:", partnerbox.Text.String)
		fmt.Println("Did >>>>>>")
		fmt.Println(dotoday.Text.String)
		fmt.Println("Learned >>>>>>")
		fmt.Println(learntoday.Text.String)

		widget.Rows = [][]gui.Widget {
			{ &gui.Text{ "Thank you, " + namebox.Text.String + "!" } },
		}
		return nil, true
	}
	return widget
}
