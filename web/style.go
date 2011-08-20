package web

import (
	"http"
	"io"
)

func styleServer(c http.ResponseWriter, req *http.Request) {
	c.Header().Set("Content-Type", "text/css")
	io.WriteString(c, Style)
}

var Style string = `
html {
    margin: 0;
    padding: 0;
}

body {
    margin: 0;
    padding: 0;
    background: #ffffff;
    font-family: arial,helvetica,"sans serif";
    font-size: 12pt;
    background: white;
}
h1 {
font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 16pt;
}
h2 { font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 14pt;
}
td tr.odd {
  background-color: #bbbbff;
}
td dr.even {
  background-color: #ffffff;
}
p {
font-family: arial,helvetica,"sans serif";
font-size:12pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 12pt;
}
a {
  color: #555599;
}
`
