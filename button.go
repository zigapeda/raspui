package raspui

type Button struct {
	rect
	text     string
	pressed  bool
	callback func()
}

func CreateButton(x, y, width, height int, text string, callback func()) *Button {
	btn := &Button{}
	btn.x = x
	btn.y = y
	btn.width = width
	btn.height = height
	btn.text = text
	btn.pressed = false
	btn.callback = callback
	return btn
}

func (b *Button) draw() {
	//TODO override olt text if new text is shorter
	drawFilledBox(b.x, b.y, b.x+b.width, b.y+b.height, DARKBLUE)
	//text with size 0 has a height of 16 pixel
	ty := int((b.height-16)/2 + b.y)
	if b.pressed == false {
		drawFilledBox(b.x+1, b.y+1, b.x+b.width-1, b.y+b.height-1, BLUE)
		drawText(b.x+10, ty, 0, b.text, BLACK, BLUE)
	} else {
		drawFilledBox(b.x+1, b.y+1, b.x+b.width-1, b.y+b.height-1, LIGHTBLUE)
		drawText(b.x+10, ty, 0, b.text, BLACK, LIGHTBLUE)
	}
}

func (b *Button) getRect() rect {
	return b.rect
}

func (b *Button) setRect(r rect) {
	b.rect = r
}

func (b *Button) istouchable(x, y int) bool {
	return b.x <= x && b.y <= y && b.x+b.width >= x && b.y+b.height >= y
}

func (b *Button) touch(x, y int) bool {
	if b.pressed == false {
		b.pressed = true
		invokeLater(b.draw)
		b.callback()
	}
	return true
}

func (b *Button) stoptouch() {
	b.pressed = false
	invokeLater(b.draw)
}

func (b *Button) SetText(text string) {
	b.text = text
	invokeLater(b.draw)
}

func (b *Button) GetText() string {
	return b.text
}
