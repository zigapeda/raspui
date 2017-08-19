package raspui

type rect struct {
	x      int
	y      int
	width  int
	height int
	draw   bool
}

type drawable interface {
	doDraw()
	getRect() rect
	setRect(rect)
	setDrawable(bool)
}

type touchable interface {
	istouchable(int, int) bool
	touch(int, int) bool
	stoptouch()
}
