package view

import (
	"ui/tokenizer"
	"ui/view/color"
)

func (self *Surface) drawTextToken(tkn *tokenizer.Token, bounds Bounds, style Style) {
	if tkn.Selected {
		self.SetSourceRGBA(color.Selection)
		self.Rectangle(bounds.X, bounds.Y, bounds.Width, bounds.Height)
		self.Fill()
	}
	
	self.SelectFontFace(style.FontName(), style.FontSlant(), style.FontWeight())
	self.SetFontSize(style.FontSize())
	self.SetSourceRGBA(style.Foreground())
	self.MoveTo(bounds.X, bounds.Y)
	self.ShowText(tkn.Value)
}

var EXTENTS map[string]*TextExtents = make(map[string]*TextExtents)

// DrawWrappedPlainText uses the Style's foreground color to draw plain
// ascii formatted text within the bounds.  It stops rendering at last visible
// line, but continues to calculate total height.
// returns height 
func (self *Surface) DrawWrappedPlainText(tokens []*tokenizer.Token, bounds Bounds, offset ScrollOffset, style Style) (linesDrawn, lines, height float64) {
	bounds.X += style.PaddingLeft()
	bounds.Y += style.PaddingTop()
	bounds.Width -= (style.PaddingLeft() + style.PaddingRight())
	bounds.Height -= (style.PaddingTop() + style.PaddingBottom())
	
	var lineHeight, x, y float64 = 0, 0, 0
	
	self.SelectFontFace(style.FontName(), style.FontSlant(), style.FontWeight())
	self.SetFontSize(style.FontSize())
	self.SetSourceRGBA(style.Foreground())
	
	spaceExtents := self.TextExtents("M")
	y += spaceExtents.Height
	spaceExtents.Width *= 1
	spaceExtents.Height *= 1.75
	 	
	self.SetFontOptions(defaultFontOptions)
	selected := make([]*tokenizer.Token, 0)
	
	for i := 0; i < len(tokens); i++ {
		t   := tokens[i]
		e   := EXTENTS[t.Value]
		if e == nil {
			e = self.TextExtents(t.Value)
			EXTENTS[t.Value] = e
		}
		
		if t.Type == tokenizer.NEWLINE {
			x = 0
			if ScrollOffset(lines) >= offset {
				y += spaceExtents.Height
			}
			lines++
			continue
		}	
		
		if t.Type == tokenizer.TAB {
			x += spaceExtents.Width * float64(style.TabWidth())
			continue
		}
	
		if e.Height > lineHeight {
			lineHeight = e.Height
		}
		
		if x + e.Width > bounds.Width {
			
			if ScrollOffset(lines) >= offset {
				y += spaceExtents.Height
			}
			x = 0
			lines++
		}
		
		var b Bounds
		b.X = bounds.X + x
		b.Y = bounds.Y + y
		b.Width  = e.Xadvance - b.X
		b.Height = e.Yadvance - b.Y
		
		if t.Selected {
			selected = append(selected, t)
		}
		
		if y < bounds.Height && ScrollOffset(lines) >= offset {
			linesDrawn = float64(ScrollOffset(lines) - offset)
			self.drawTextToken(t, b, style)
		} 
		x += e.Xadvance
		y += e.Yadvance 
	}
	
	return linesDrawn, lines, y + style.PaddingBottom()
}

func (self *Surface) DrawVerticalOverflow(boundsHeight, height, percent float64, s Style) {
	if height < 5 {
		height = 10
	}
	switch s.OverflowY() {
	case STYLE_OVERFLOW_Y_SCROLL:
		self.SetSourceRGBA(color.HexRGBA(0x00000025))
		self.RoundedRectangle(float64(self.GetWidth()) - 11, boundsHeight * percent, 10, height, 3, 3, 3, 3)
		self.Fill()
	case STYLE_OVERFLOW_Y_SHADE:
	}
}

func (self *Surface) DrawHorizontalOverflow(boundsWidth float64, s Style) {
	ratio := float64(self.GetWidth()) / boundsWidth
	self.SetSourceRGBA(color.HexRGBA(0x00000025))
	self.RoundedRectangle(float64(self.GetHeight()) - 10, 50, 10, float64(self.GetWidth()) * ratio, 2, 2, 2, 2)
	self.Fill()
}