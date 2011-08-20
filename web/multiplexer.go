package web

import (
	"http"
)

type commChannel struct {
	string
	pages chan []byte // this carries the html sent to the client
	requests chan *http.Request // this carries posts sent by the client

	// the events is the chan carrying events to the user
	// goroutine, which will normally respond by sending a page through
	// the "pages" chan.
	events chan event
}

type pageRequest struct {
	string
	ch chan []byte // send the page here please!
}

type event struct {
	string
	widget string
	event string
}

var newconnection chan<- commChannel
var closeconnection chan<- string
var pagerequests chan<- pageRequest
var incomingevents chan<- event

func init() {
	connmap := make(map[string]commChannel)

	// first let's create the "new connection" creator
	nc := make(chan commChannel)
	newconnection = nc
	// now let's create the "closer" channel
	cl := make(chan string)
	closeconnection = cl
	// finally, we create the chan for distributing pages to the worthy
	pr := make(chan pageRequest)
	pagerequests = pr

	ie := make(chan event)
	incomingevents = ie

	go func() {
		for {
			select {
			case cc := <- nc:
				connmap[cc.string] = cc
			case toclose := <- cl:
				connmap[toclose] = connmap[toclose], false
			case req := <- pr:
				if cc,ok := connmap[req.string]; ok {
					go func() {
						// when ready, send a page from the source to the sink
						req.ch <- (<- cc.pages)
					}()
				}
			case e := <- ie:
				if cc,ok := connmap[e.string]; ok {
					go func() {
						// send the event to the appropriate channel...
						cc.events <- e
					}()
				}
			}
		}
	}()
}
