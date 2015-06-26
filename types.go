package raspui

type rect struct {
	x      int
	y      int
	width  int
	height int
}

type drawable interface {
	draw()
	getRect() rect
	setRect(rect)
}

type touchable interface {
	istouchable(int, int) bool
	touch(int, int) bool
	stoptouch()
}
