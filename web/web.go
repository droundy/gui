package web

import (
	"gui/data"
	"io"
	"fmt"
	"os"
	"http"
	"html"
	"path"
)

func WidgetToHtml(parent string, widget data.Widget) (out string) {
	mypath := path.Join(parent, widget.Name())
	switch widget := widget.(type) {
	case *data.Text:
		return html.EscapeString(widget.String)
	case *data.EditText:
		myname := widget.Text.String
		return `<input type="text" onchange="say('` + mypath +
			`',  'onchange:'+this.value)" value="` + html.EscapeString(myname) + `" />`
	case *data.TextArea:
		myname := widget.Text.String
		return `<textarea cols="80" rows="5" onchange="say('` + mypath +
			`',  'onchange:'+this.value)">` + html.EscapeString(myname) + `</textarea>`
	case *data.Table:
		out = "<table>\n"
		for i,r := range widget.Rows {
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
	case *data.Button:
		myname := widget.Text.String
		return `<input type="submit" onclick="say('` + mypath +
			`',  'onclick')" value="` + html.EscapeString(myname) + `" />`
	case *data.Window:
		return WidgetToHtml(parent, widget.Widget)
	default:
		panic(fmt.Sprintf("Unhandled data.Widget type! %T", widget))
	}
	return
}

func Serve(port int, newWidget func() *data.Window) os.Error {
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
		widget := newWidget()

		n := "job-" + (<- uniqueids)
		oldtitle := widget.Title
		oldpath := widget.Path
		io.WriteString(w, skeletonpage(oldtitle, n, req))

		cc := commChannel{ n, make(chan []byte), make(chan *http.Request), make(chan event) }
		go func() {
			//fmt.Printf("widget is:\n%#v\n", widget)
			//fmt.Println("Html is:", WidgetToHtml("", widget))
			cc.pages <- []byte(`settitle ` + oldtitle)
			cc.pages <- []byte(`setpath ` + oldpath)
			cc.pages <- []byte(WidgetToHtml("", widget))
			for {
				e := <- cc.events
				fmt.Printf("Event is %#v\n", e)
				//fmt.Printf("Corresponding widget is %#v\n", widget.Lookup(e.widget))
				newWidget, refresh := widget.Handle(data.Event{e.widget, e.event})
				//fmt.Printf("New widget is %#v\n", newWidget)
				if newWidget, ok := newWidget.(*data.Window); ok {
					widget = newWidget
					refresh = true
				}
				if widget.Title != oldtitle {
					oldtitle = widget.Title
					cc.pages <- []byte(`settitle ` + oldtitle)
				}
				if widget.Path != oldpath {
					oldpath = widget.Path
					cc.pages <- []byte(`setpath ` + oldpath)
				}
				if refresh {
					cc.pages <- []byte(WidgetToHtml("", widget))
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
             history.pushState('', response.substr(8), response.substr(8));
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
