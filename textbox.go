package raspui

type Textbox struct {
	rect
	text      string
	oldlength int
}

func CreateTextbox(x, y, width, height int, text string) *Textbox {
	txt := &Textbox{}
	txt.x = x
	txt.y = y
	txt.width = width
	txt.height = height
	txt.text = text
	txt.oldlength = 0
	return txt
}

func (t *Textbox) draw() {
	//text with size 0 has a height of 16 pixel
	ty := int((t.height-16)/2 + t.y)
	drawFilledBox(t.x, ty, t.x+8*t.oldlength, ty+16, WHITE)
	if len(t.text) * 8 > t.width {
		drawText(t.x, ty, 0, t.text[0:t.width/8], BLACK, WHITE)
	} else {
		drawText(t.x, ty, 0, t.text, BLACK, WHITE)
	}
}

func (t *Textbox) getRect() rect {
	return t.rect
}

func (t *Textbox) setRect(r rect) {
	t.rect = r
}

func (t *Textbox) SetText(text string) {
	t.oldlength = len(t.text)
	t.text = text
	declareInvalid(t)
}

func (t *Textbox) GetText() string {
	return t.text
}
