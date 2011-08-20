package web

import (
	"gui"
	"io"
	"fmt"
	"os"
	"http"
	"html"
	"path"
)

func WidgetToHtml(parent string, widget gui.Widget) (out string) {
	switch widget := widget.(type) {
	case *gui.Text:
		return html.EscapeString(widget.String)
	case *gui.Table:
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
				whtml := WidgetToHtml(fmt.Sprint(i, "/", j), w)
				out += "    <td>" + whtml + "</td>\n"
			}
			out += "  </tr>\n"
		}
		out += "</table>\n"
	case *gui.Button:
		myname := widget.Text.String
		mypath := path.Join(parent, myname)
		return `<input type="submit" onclick="say('` + mypath +
			`',  'onclick')" value="` + html.EscapeString(myname) + `" />`
	default:
		panic(fmt.Sprintf("Unhandled gui.Widget type! %T", widget))
	}
	return
}

func Serve(port int, widget gui.Widget) os.Error {
	// We have a style sheet called style.css
	http.HandleFunc("/style.css", styleServer)
	http.HandleFunc("/jsupdate", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("query = ", req.URL.RawQuery)

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
		n := "job-" + (<- uniqueids)
		io.WriteString(w, skeletonpage(n, req))

		cc := commChannel{ n, make(chan []byte), make(chan *http.Request), make(chan event) }
		go func() {
			// This is the generator of pages
			fmt.Println("Html is:", WidgetToHtml("", widget))
			cc.pages <- []byte(WidgetToHtml("", widget))
			// FIXME I should handle events next!
			for {
				e := <- cc.events
				fmt.Printf("Event is %#v\n", e)
				cc.pages <- []byte(WidgetToHtml("", widget))
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

func skeletonpage(query string, req *http.Request) string {
	return `<!DOCTYPE HTML>
<html>
  <head>
    <title>Long polling test</title>
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
           //var received_msg = evt.data;
           //alert("Message is received: " + received_msg);
           everything.innerHTML=response;
           //$('textarea').text(response);
          _poll();
        });
      }
      
      $.get('/jsstatus', function(response) {
        $('textarea').text(response);
        _poll();
      });
    }
  </script>
</html>
`
}
