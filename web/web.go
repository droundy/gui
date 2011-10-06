//target:github.com/droundy/gui/exp/web
package web

import (
	"github.com/droundy/gui/exp/data"
	"github.com/droundy/gui/exp/gui"
	"io"
	"fmt"
	"os"
	"http"
	"html"
	"path"
	"strings"
	"strconv"
)

func widgetName(w data.Widget) (out string) {
	out = fmt.Sprintf("%t", w)
	switch w := w.(type) {
	}
	return
}

func lookupWidget(p string, w data.Widget) data.Widget {
	if p == widgetName(w) {
		return w
	}
	switch w := w.(type) {
	case *data.Text, *data.EditText, *data.Button:
		return w
	case *data.Column:
		s := strings.SplitN(p, "/", 2)
		if len(s) != 2 {
			panic("Weird bug:  not a nice name")
		}
		i, err := strconv.Atoi(s[0])
		if err != nil || i >= len(*w) {
			panic("Weird bug:  not a good row")
		}
		return lookupWidget(s[1], (*w)[i])
	case *data.Table:
		s := strings.SplitN(p, "/", 3)
		if len(s) != 3 {
			panic("Weird bug:  not a nice name")
		}
		i, err := strconv.Atoi(s[0])
		if err != nil || i >= len(*w) {
			panic("Weird bug:  not a good row")
		}
		r := (*w)[i]
		j, err := strconv.Atoi(s[1])
		if err != nil || j >= len(r) {
			panic(fmt.Sprint("Weird bug:  not a good column: ", s[1]))
		}
		return lookupWidget(s[2], r[j])
	default:
		panic(fmt.Sprintf("Unknown lookupWidget type %#v\n", w))
	}
	return nil
}

func WidgetToHtml(parent string, widget data.Widget) (out string) {
	mypath := path.Join(parent, widgetName(widget))
	switch widget := widget.(type) {
	case *data.Text:
		return html.EscapeString(widget.String)
	case *data.EditText:
		myname := widget.String
		return `<input type="text" onchange="say('` + mypath +
			`',  'onchange:'+this.value)" value="` + html.EscapeString(myname) + `" />`
	case *data.TextArea:
	 	myname := widget.String
	 	return `<textarea cols="80" rows="5" onchange="say('` + mypath +
	 		`',  'onchange:'+this.value)">` + html.EscapeString(myname) + `</textarea>`
	case *data.Table:
		out = "<table>\n"
		for i,r := range *widget {
			class := "even" // I define classes for even and odd rows
			switch {        // so you can alternate colors if you like.
			case i == 0:    // I also define "even first" as a possible header
				class = "even first"
			case i & 1 == 1:
				class = "odd"
			}
			out += `  <tr class="`+ class + `">` + "\n"
			for j,w := range r {
				whtml := WidgetToHtml(path.Join(parent, fmt.Sprint(i, "/", j)), w)
				out += "    <td>" + whtml + "</td>\n"
			}
			out += "  </tr>\n"
		}
		out += "</table>\n"
	case *data.Column:
		out = ""
		for i,w := range *widget {
			class := "even" // I define classes for even and odd rows
			switch {        // so you can alternate colors if you like.
			case i == 0:    // I also define "even first" as a possible header
				class = "even first"
			case i & 1 == 1:
				class = "odd"
			}
			whtml := WidgetToHtml(path.Join(parent, fmt.Sprint(i)), w)
			out += `<p class="` + class + `">` + whtml + "</p>\n"
		}
	case *data.Button:
		myname := widget.String
		return `<input type="submit" onclick="say('` + mypath +
			`',  'onclick')" value="` + html.EscapeString(myname) + `" />`
	case *data.Menu:
		myname := widget.Options[widget.Value]
		out = `<select onchange="say('` + mypath +
			`',  'onchange:'+this.value)" value="` + html.EscapeString(myname) + `">`
		for i,v := range widget.Options {
			if i == widget.Value {
				out += "\n<option value=\"" + v + `" selected='selected'>` + v + "</option>"
			} else {
				out += "\n<option value=\"" + v + `">` + v + "</option>"
			}
		}
		out += "\n</select>"
	// case *data.Window:
	// 	return WidgetToHtml(parent, widget.Widget)
	default:
		panic(fmt.Sprintf("Unhandled data.Widget type in WidgetToHtml! %T", widget))
	}
	return
}

func Serve(port int, newWidget func() gui.Window) os.Error {
	// We have a style sheet called style.css
	http.HandleFunc("/style.css", styleServer)
	http.HandleFunc("/jsupdate", func(w http.ResponseWriter, req *http.Request) {
		//fmt.Println("query = ", req.URL.RawQuery)

		// Wait for the response...
		ch := make(chan []byte)
		pagerequests <- pageRequest{ req.URL.RawQuery, ch }
		w.Write( <-ch )
	})
	// Events come via the URL "/say"
	http.HandleFunc("/say", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		//fmt.Printf("Got a message: %v with url %#v and RawURL %v\n",
		//	req.URL.RawQuery, req.URL, req.RawURL)
		//fmt.Printf("Form is %#v\n", req.Form)
		if u, ok := req.Form["user"]; ok {
			if w, ok := req.Form["widget"]; ok {
				if e, ok := req.Form["event"]; ok {
					if len(u) == 1 && len(w) == 1 && len(e) == 1 {
						incomingevents <- event{ u[0], w[0], e[0] }
					}
				}
			}
		}
	})
	// And we listen on the root for our gui program
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// This is the generator of pages
		window := newWidget()

		n := "job-" + (<- uniqueids)
		oldtitle := window.Title
		oldpath := window.Path
		io.WriteString(w, skeletonpage(oldtitle, n, req))

		cc := commChannel{ n, make(chan []byte), make(chan *http.Request), make(chan event) }
		go func() {
			cc.pages <- []byte(`settitle ` + oldtitle)
			cc.pages <- []byte(`setpath ` + oldpath)
			cc.pages <- []byte(WidgetToHtml("", window.Contents.Raw()))
			for {
				select {
				case e := <- cc.events:
					fmt.Printf("Event is %#v\n", e)
					eventfor := lookupWidget(e.widget, window.Contents.Raw())
					switch {
					case e.event[:7] == "onclick":
						if eventfor, ok := eventfor.(gui.Clickable); ok {
							eventfor.Clicks() <- struct{}{}
						} else {
							fmt.Printf("Got unexpected click for %#v\n", eventfor)
						}
					case e.event[:8] == "onchange":
						if eventfor, ok := eventfor.(gui.Changeable); ok {
							eventfor.Changes() <- e.event[9:]
						} else {
							fmt.Printf("Got unexpected change for %#v\n", eventfor)
						}
					default:
						fmt.Printf("Got weird event: %#v\n", e.event)
					}
				case window.Contents = <- window.Contents.Updater():
					cc.pages <- []byte(WidgetToHtml("", window.Contents.Raw()))					
				}
			}
		}()

		newconnection <- cc
	})
	return http.ListenAndServe(fmt.Sprint(":", port), nil)
}

var uniqueids <-chan string

func init() {
	n := make(chan string)
	uniqueids = n
	go func() {
		for i:=0; true; i++ {
			n <- fmt.Sprint(i)
		}
	}()
}

func skeletonpage(title, query string, req *http.Request) string {
	return `<!DOCTYPE HTML>
<html>
  <head>
    <title>` + title + `</title>
    <link href="/style.css" rel="stylesheet" type="text/css" />
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
  </head>
  
  <body>
    <div id="everything">
      Everything goes here.
    </div>
  </body>
  <script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js"></script>
  <script type="text/javascript">
    function say(who, what) {
      $.post("say", { user: "`+query+`", widget: who, event: what })
    };
    var client = new function() {
      var _poll = function() {
        $.get('/jsupdate?`+ query +`', function(response) {
           var everything = document.getElementById("everything")
           if (everything == null) {
             return
           }
           if (response.substr(0,8) == 'setpath ') {
             if (history.pushState) { // workaround for older browsers
               history.pushState('', response.substr(8), response.substr(8));
             }
           } else if (response.substr(0,9) == 'settitle ') {
             document.title = response.substr(9)
           } else if (response.length < 10) {
             return // looks like server has exited
           } else {
             everything.innerHTML=response;
           }
           _poll();
        });
      }
      
      _poll();
    }
  </script>
</html>
`
}
