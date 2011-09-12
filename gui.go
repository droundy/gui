//target:github.com/droundy/gui
package gui

// This module will eventually provide the end-user interface for gui.
// The idea is that there is a separation between different possible
// ways to use the gui package.
//
// Developers who want to just write a portable gui application should
// just import this module (once it has been written) and one or more
// backend modules (currently github.com/gui/web is the only backend),
// which will take charge of the actual rendering.
//
// Developers who want to write a backend will need to import from
// github.com/gui/data, as well as this module (once it is written).
//
// Developers who want to develop new widgets will have the option of
// either importing from this module, or importing from gui/data.  If
// possible, they should only import from this module, as the precise
// data format is not considered stable.  It will periodically be
// necessary to make non-backwards-compatible changes to gui/data, so
// most users should instead (once this is written) import from gui
// directly and use this (stable, but non-existent) API.

import (
	"github.com/droundy/gui/data"
)

type Widget interface {
	ToRaw() data.Widget
}

type HandleClick func()
type Clickable interface {
	HandleClick()
}

type HandleChange func(old string)
type Changeable interface {
	HandleChange(old string)
}
