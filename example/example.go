package main

import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/droundy/gui"
	"github.com/droundy/gui/web"
	"time"
	"strings"
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

var surveyfile *os.File

func init() {
	sf, err := os.OpenFile("survey.tex", os.O_WRONLY + os.O_APPEND + os.O_CREATE, 0666)
	surveyfile = sf
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error appending to survey.tex:", err)
		os.Exit(1)
	}
}

func NewWidget() gui.Window {
	window := gui.Window{ "Class survey", "", nil }
	
	teamname := ""
	team := gui.Menu(0,
		[]string{
		"", "archimedes", "boltzmann", "curie", "doppler", "euler", "feynman", "galileo",
		"hamilton", "ising", "joule", "kelvin", "lagrange", "maxwell", "newton", "onsager", "planck",
	})
	teamrow := gui.Table([][]gui.Widget{{ gui.Text("Team:"), team }})

	name := ""
	namebox := gui.EditText(name)
	namerow := gui.Row(gui.Text("Name:"), namebox)
	// namebox.HandleChanged = func(old string) (modified gui.Widget, refresh bool) {
	// 	window.Title = `Survey of ` + namebox.Text.String
	// 	return
	// }
	partner := ""
	partnerbox := gui.EditText(partner)
	partnerrow := gui.Row(gui.Text("Partner:"), partnerbox)
	donetext := ""
	dotoday := gui.TextArea(donetext)
	learnedtoday := ""
	learntoday := gui.TextArea(learnedtoday)
	problems := ""
	workwell := gui.TextArea(problems)

	button := gui.Button("Submit")

	widget := gui.Table([][]gui.Widget{
		{ teamrow },
		{ namerow },
		{ partnerrow },
		{ gui.Text("What did you do today?") },
		{ dotoday },
		{ gui.Text("What is one thing you learned today?") },
		{ learntoday },
		{ gui.Text("What is one thing that didn't work well today?") },
		{ workwell },
		{ button },
	})
	window.Contents = widget
	go func() {
		for {
			select {
			case teamname = <- team.Changes():
				fmt.Println("Team name changed to", teamname)
			case name = <- namebox.Changes():
				//fmt.Println("Name changed to", name)
			case partner = <- partnerbox.Changes():
				//fmt.Println("Partner changed to", partner)
			case donetext = <- dotoday.Changes():
				fmt.Println("Done text is", donetext)
			case learnedtoday = <- learntoday.Changes():
				fmt.Println("Learned today is", learnedtoday)
			case problems = <- workwell.Changes():
				fmt.Println("Problems is", problems)
			case _ = <- button.Clicks():
				t := time.LocalTime()
				// First let's see if today has already been created
				if _,err := os.Stat(t.Format("2006-01-02")); err != nil {
					surveyfile.WriteString(t.Format("\\thisday{Monday}{2006-01-02}\n\n"))
				} else {
					fmt.Println("Day already exists.")
				}
				
				dir := t.Format("2006-01-02/15.04.05")
				err := os.MkdirAll(dir, 0777)
				if err != nil {
					fmt.Println("ERROR CREATING DIRECTORY", dir, "!")
					return
				}
				f, err := os.Create(filepath.Join(dir, name))
				if err != nil {
					fmt.Println("ERROR CREATING FILE", filepath.Join(dir, name), "!", err)
					return
				}
				defer f.Close()
				_,err = fmt.Fprintf(f, "\\daily{%s}{%s}{%s}{%s}{\n%s\n}{\n%s\n}{\n%s\n}\n",
					t.Format("3:04PM"),
					name, partner, teamname,
					IndentWrapText("  ", CleanLatex(donetext)),
					IndentWrapText("  ", CleanLatex(learnedtoday)),
					IndentWrapText("  ", CleanLatex(problems)))
				if err == nil {
					surveyfile.WriteString(t.Format("\\input{2006-01-02/15.04.05/" +
						name +"}\n"))
				} else {
					fmt.Println("I ran into a bug!", err)
					return
				}
				window.Contents.Updater() <- gui.Text("Thank you, " + name + "!")
			}
		}
	}()
	return window
}

func CleanLatex(input string) (out string) {
	aminmath := false
	outints := []int{}
	for _,c := range input {
		switch c {
		case '$':
			aminmath = !aminmath
		}
		if !aminmath {
			switch c {
			case '_', '^', '\\':
				outints = append(outints, '\\')
			}
		}
		outints = append(outints, c)
	}
	out = string(outints)
	//out = strings.Replace(out, "_", "\\_")
	return out

	//underscore := regexp.MustCompile(`_`)
	//return underscore.ReplaceAllString(input, `\_`)
}

func IndentWrapText(indent, input string) string {
	out := []string{}
	nextline := indent
	words := strings.Split(input, " ")
	for _,w := range words {
		if len(nextline) + 1 + len(w) < 80 {
			nextline += " " + w
		} else {
			out = append(out, nextline)
			nextline = indent + w
		}
	}
	out = append(out, nextline)
	return strings.Join(out, "\n")
}
