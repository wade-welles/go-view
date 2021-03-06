// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2014 Stanley Steel
package editor

import (
	"log"

	"github.com/sesteel/go-view"
	"github.com/sesteel/go-view/color"
	"github.com/sesteel/go-view/event"
	"github.com/sesteel/go-view/tokenizer"
)

// Selection represents a character range where text has been
// selected; commonly used for copy and cutting operations.
type Selection struct {
	Range
}

func (self Selection) IndexInSelection(i Index) bool {
	self.Normalize()
	if i.Line > self.Start.Line && i.Line < self.End.Line {
		return true
	} else if i.Line == self.Start.Line && i.Column >= self.Start.Column {
		return true
	} else if i.Line == self.End.Line && i.Column <= self.End.Column {
		return true
	}
	return false
}

func (self Selection) drawCharBG(s *view.Surface, lines tokenizer.Line, i Index, x, y, w, h float64) {
	w += 2
	h += 1
	s.Save()
	s.SetSourceRGBA(color.Blue1)
	s.Rectangle(x, y, w, h)
	s.Fill()
	self.Normalize()

	// s.SetSourceRGBA(color.Blue2)

	// top := func() {
	// 	s.MoveTo(x, y)
	// 	s.LineTo(x+w, y)
	// }

	// left := func() {
	// 	s.MoveTo(x, y)
	// 	s.LineTo(x, y+h)
	// }

	// bottom := func() {
	// 	s.MoveTo(x, y+h)
	// 	s.LineTo(x+w, y+h)
	// }

	// right := func() {
	// 	s.MoveTo(x+w, y)
	// 	s.LineTo(x+w, y+h)
	// }

	// if i.Column == 0 {
	// 	left()
	// }

	// if i.Line == self.Start.Line {
	// 	top()
	// 	if i.Column == self.Start.Column {
	// 		left()
	// 	}
	// } else {
	// 	// fmt.Println(i)
	// 	a := self.IndexInSelection(Index{i.Line - 1, i.Column})
	// 	b := len(lines[i.Line]) > 0
	// 	c := true
	// 	if i.Line > 0 {
	// 		c = len(lines[i.Line-1]) <= i.Column
	// 	}
	// 	if !a || b && c {
	// 		top()
	// 	}
	// }

	// if i.Line == self.End.Line {
	// 	bottom()
	// 	if i.Column == self.End.Column {
	// 		right()
	// 	}
	// } else {
	// 	if !self.IndexInSelection(Index{i.Line + 1, i.Column}) || len(lines[i.Line+1]) <= i.Column {
	// 		bottom()
	// 	}
	// }

	// if i.Column == len(lines[i.Line])-1 {
	// 	right()
	// }

	// s.Stroke()
	s.Restore()
}

func (self *Editor) handleScrollMapOffset(ev event.Mouse) event.Mouse {
	if self.DrawScrollMap {
		log.Println("smw:", self.scrollMap.Width)
		ev.X -= self.scrollMap.Width
	}
	return ev
}

func (self *Editor) addTextSelectionBehavior() {

	self.AddMouseButtonPressHandler(func(ev event.Mouse) {
		switch ev.Button {
		case event.MOUSE_BUTTON_LEFT:
			// ev := self.handleScrollMapOffset(ev)
			self.MoveCursor(float64(ev.X), float64(ev.Y))
			kb := event.LastKeyboardState()
			if kb.CtrlOnly() {
				self.Selection = &Selection{Range{Index{-1, -1}, Index{-1, -1}}}
			} else {
				self.Selections = make([]*Selection, 0)
			}
			self.Redraw()
		}
	})

	self.AddMouseButtonReleaseHandler(func(ev event.Mouse) {
		switch ev.Button {
		case event.MOUSE_BUTTON_LEFT:
			if self.Selection.Start.Line > -1 {
				// ev := self.handleScrollMapOffset(ev)
				idx := self.FindClosestIndex(ev.X, ev.Y)
				l, _ := idx.Line, idx.Column
				if l >= 0 {
					self.Selection = &Selection{Range{Index{-1, -1}, Index{-1, -1}}}
				}
			}
			self.Redraw()
		}
	})

	self.AddMousePositionHandler(func(ev event.Mouse) {
		if ev.LeftPressed {
			// ev := self.handleScrollMapOffset(ev)
			idx := self.FindClosestIndex(ev.X, ev.Y)
			if idx.Line >= 0 && idx.Column >= 0 && self.Selection.Start.Line == -1 {
				self.Selection.Start = idx
				self.Selection.End.Column = idx.Column + 1
				self.Selections = append(self.Selections, self.Selection)
			} else if idx.Line >= 0 && idx.Column >= 0 && self.Selection.Start.Line > -1 {
				self.Selection.End = idx
				self.Selection.End.Column = idx.Column - 1
				self.MoveCursor(float64(ev.X), float64(ev.Y))
			}
			self.Redraw()
		}
	})
}
