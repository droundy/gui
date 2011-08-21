package main

import (
	"os"
	"fmt"
	"gui"
	"web"
)

func main() {
	button := &gui.Button{Text: gui.Text{"Hello world"}}
	button.HandleClick = func() gui.Widget {
		if button.Text.String == "Hello world" {
			button.Text.String = "Goodbye world"
		} else {
			button.Text.String = "Hello world"
		}
		return button
	}
	widget := &gui.Table{
		[][]gui.Widget{
			{ &gui.Text{"Hello world"} },
			{ button },
			{ &gui.Text{"Goodbye world"}, &gui.EditText{Text: gui.Text{"And the end"}} },
		},
	}
	err := web.Serve(12345, widget)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	} else {
		fmt.Println("Exited successfully!")
	}
}
