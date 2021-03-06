// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2014 Stanley Steel
package editor

import (
	"log"
	"time"

	"github.com/sesteel/go-view"
	"github.com/sesteel/go-view/color"
	"github.com/sesteel/go-view/geometry"
	"github.com/sesteel/go-view/tokenizer"
)

type CursorType int

const (
	BAR CursorType = iota
	BLINK
	BLOCK
	FADE
	OUTLINE
	UNDERLINE
)

// Cursor is used to store the position of the cursor via a Index.
type Cursor struct {
	Index
	Type     CursorType
	Color    *color.RGBA
	lastedit int
	wait     int
}

func (self Cursor) PreviousPos(lines []tokenizer.Line) Index {
	if len(lines) < self.Line-1 || self.Line < 0 {
		return Index(self.Index)
	}

	if self.Column == 0 {
		return Index{self.Line - 1, len(lines[self.Line-1].Characters)}
	} else {
		return Index{self.Line, self.Column - 1}
	}
}

func (self *Cursor) Draw(s *view.Surface, b *geometry.Bounds, e *Editor) {
	offset := int(e.vscroll.Offset())
	surfaces := e.lineSurfaces
	now := time.Now().Nanosecond()

	// the cursor is not visible
	if offset > self.Line || b == nil {
		return
	}

	yoff := 0.0

	for i, surface := range surfaces {
		if i == self.Line {
			break
		}

		if surface == nil {
			log.Println("Cursor cannot draw on nil surface... ", self)
		}

		if i >= offset {
			yoff += float64(surface.Height())
		}
	}

	y := b.Y + yoff

	blink := func(max float64) {
		if now < 450000000 {
			self.Color.A = 0
		} else {
			self.Color.A = max
		}
	}

	prepare := func() {
		s.SetSourceRGBA(*self.Color)
		s.SetLineWidth(1)
	}

	// fade
	switch self.Type {

	case BAR:
		self.Color.A = 1
		prepare()
		s.MoveTo(b.X, y-b.Height)
		s.LineTo(b.X, y+b.Height/2)
		s.Stroke()

	case BLINK:
		blink(1)
		prepare()
		s.MoveTo(b.X, y-b.Height)
		s.LineTo(b.X, y+b.Height/2)
		s.Stroke()

	case BLOCK:
		blink(0.5)
		prepare()
		s.RoundedRectangle(b.X+ALIGN, y-b.Height+ALIGN, b.Width, b.Height+b.Height/2, 1, 1, 1, 1)
		s.Fill()
		s.RoundedRectangle(b.X+ALIGN, y-b.Height+ALIGN, b.Width, b.Height+b.Height/2, 1, 1, 1, 1)
		s.Stroke()

	case FADE:
		if now < 300000000 {
			self.Color.A += 0.05
		} else if now < 600000000 {
			self.Color.A = 1
		} else if self.Color.A > 0 {
			self.Color.A -= 0.05
		}
		prepare()
		s.MoveTo(b.X, y-b.Height)
		s.LineTo(b.X, y+b.Height/2)
		s.Stroke()

	case OUTLINE:
		blink(0.7)
		prepare()
		s.RoundedRectangle(b.X, y-b.Height+ALIGN, b.Width, b.Height+b.Height/2, 2, 2, 2, 2)
		s.Stroke()

	case UNDERLINE:
		blink(1)
		prepare()
		s.MoveTo(b.X, y+ALIGN+b.Height/2)
		s.LineTo(b.X+b.Width, y+ALIGN+b.Height/2)
		s.Stroke()
	}
	s.Flush()
}
