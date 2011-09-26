package main

import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/droundy/gui/data"
	"github.com/droundy/gui/web"
	"time"
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

func NewWidget() *data.Window {
	window := data.Window{ "Class survey", "", nil }
	
	team := &data.Menu {
		Options: []string{
			"", "archimedes", "boltzmann", "curie", "doppler", "euler", "feynman", "galileo",
			"hamilton", "ising", "joule", "kelvin", "lagrange", "maxwell", "newton", "onsager", "planck",
		},
	}
	teamrow := &data.Table{{ &data.Text{String: "Team:"}, team }}

	namebox := &data.EditText{}
	namerow := &data.Table {{ &data.Text{String: "Name:"}, namebox }}
	namebox.HandleChanged = func(old string) (modified data.Widget, refresh bool) {
		window.Title = `Survey of ` + namebox.Text.String
		return
	}
	partnerbox := &data.EditText{}
	partnerrow := &data.Table{{ &data.Text{String: "Partner:"}, partnerbox }}
	dotoday := &data.TextArea{}
	learntoday := &data.TextArea{}
	workwell := &data.TextArea{}

	button := &data.Button{Text: data.Text{String: "Submit"}}

	widget := &data.Table{
		{ teamrow },
		{ namerow },
		{ partnerrow },
		{ &data.Text{String: "What did you do today?"} },
		{ dotoday },
		{ &data.Text{String: "What is one thing you learned today?"} },
		{ learntoday },
		{ &data.Text{String: "What is one thing that didn't work well today?"} },
		{ workwell },
		{ button },
	}
	window.Widget = widget
	button.HandleClick = func() (modified data.Widget, refresh bool) {
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
		}
		f, err := os.Create(filepath.Join(dir, namebox.Text.String))
		if err != nil {
			fmt.Println("ERROR CREATING FILE", filepath.Join(dir, namebox.Text.String), "!")
		}
		defer f.Close()
		_,err = fmt.Fprintf(f, "\\daily{%s}{%s}{%s}{%s}{\n%s}{\n%s}{\n%s}\n",
			t.Format("3:04PM"),
			namebox.Text.String, partnerbox.Text.String, team.String(),
			CleanLatex(dotoday.Text.String),
			CleanLatex(learntoday.Text.String),
			CleanLatex(workwell.Text.String))
		if err == nil {
			surveyfile.WriteString(t.Format("\\input{2006-01-02/15.04.05/" +
				namebox.Text.String +"}\n"))
		} else {
			fmt.Println("I ran into a bug!", err)
		}

		*widget = [][]data.Widget {
			{ &data.Text{ String: "Thank you, " + namebox.Text.String + "!" } },
		}
		return nil, true
	}
	return &window
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
