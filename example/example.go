package main

import (
	"os"
	"fmt"
	"gui"
	"web"
)

func main() {
	widget := &gui.Table{
		[][]gui.Widget{
			{ &gui.Text{"Hello world"} },
			{ &gui.Button{Text: gui.Text{"Hello world"}} },
			{ &gui.Text{"Goodbye world"}, &gui.Text{"And the end"} },
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
