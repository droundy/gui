Experimental gui library
========================

I always build with gb.  I am not sure if I've made this thing
goinstall-friendly.  If you run into trouble with goinstall, please
let me know!

This library is still very experimental.  It is a replacement of the
former gui library that I (David Roundy) worked on.

The gui package exports almost everything you should need to create a
gui program, with the single exception being the lack of a function to
make your gui actually run.  That is provided by a "back end" package.
Currently there is just one back end, *gui/web*, which provides a web
interface.  I hope that eventually, further back ends will be
provided, possibly as third-party packages.

There aren't many widgets yet.  Widget-creation comes in one of two
varieties.  A composite widget can be created in a package that only
imports from gui itself.  However, if you want to create a new
primitive widget type (e.g. a color map or file selector), you'll need
to create a new type in gui/data, and add support for it in gui/web
(and any other back ends that are developed).

The API is based around chans and interfaces, and hopefully is a nice
go-like API.
