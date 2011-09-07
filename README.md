Experimental gui library
========================

Build with gb.  I haven't yet made this thing goinstall-friendly.

This library is still very experimental.  It is a replacement of the
former gui library that I (David Roundy) worked on.  The gui package
itself---which is intended to be the front end used by gui
writers---hasn't even been started.  So far, I've created gui/data,
which is the backend-agnostic data module, and a gui/web, which
provides a web server back end.  There isn't much there, but I have
written a functioning example using it.

Eventually, the idea is that there can be multiple back ends, so you
can write a program that pretty easily can be compiled either with GTK
or as a windows native program or whatever.  But there is no timeline
for completion.
