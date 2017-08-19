package raspui

import _ "fmt"

type Tabpane struct {
	tabs         []*Tab
	currentindex int
	touchelement touchable
}

type Tab struct {
	elements []drawable
	name     string
	tabpane  *Tabpane
}

func CreateTabpane() *Tabpane {
	tp := &Tabpane{}
	tp.tabs = make([]*Tab, 0, 5)
	tp.currentindex = -1
	tp.touchelement = nil
	return tp
}

func (tp *Tabpane) AddTab(tab *Tab) {
	tab.tabpane = tp
	tp.tabs = append(tp.tabs, tab)
	declareInvalid(tp)
}

func (tp *Tabpane) RemoveTab(tab *Tab) {
	for i, v := range tp.tabs {
		if tab == v {
			tab.tabpane = nil
			tp.tabs = append(tp.tabs[:i], tp.tabs[i+1:]...)
			declareInvalid(tp)
			return
		}
	}
}

func (tp *Tabpane) doDraw() {
	if tp.currentindex == -1 && len(tp.tabs) > 0 {
		tp.currentindex = 0
		tp.tabs[0].show()
	}
	x := 1
	tx := 0
	for i, v := range tp.tabs {
		w := len(v.name)*8 + 12
		drawBox(x, 1, x+w, 38, DARKBLUE)
		if i != tp.currentindex {
			drawFilledBox(x+1, 2, x+w-1, 37, BLUE)
			drawText(x+7, 12, 0, v.name, BLACK, BLUE)
		} else {
			tx = x
		}
		x += w + 3
	}
	drawLine(0, 38, 320, 38, DARKBLUE)
	drawLine(0, 39, 320, 39, LIGHTBLUE)
	if tp.currentindex != -1 {
		t := tp.tabs[tp.currentindex]
		w := len(t.name)*8 + 12
		drawFilledBox(tx+1, 2, tx+w-1, 38, LIGHTBLUE)
		drawText(tx+7, 12, 0, t.name, BLACK, LIGHTBLUE)
	}
}

func (tp *Tabpane) getRect() rect {
	return rect{0, 0, 320, 40, true}
}

func (tp *Tabpane) setRect(r rect) {
	//always full size
}

func (tp *Tabpane) setDrawable(draw bool) {
	//always drawable
}

func (tp *Tabpane) invalidateTab(t *Tab) {
	if tp.currentindex != -1 {
		if tp.tabs[tp.currentindex] == t {
			t.hide()
			t.show()
		}
	}
}

func (tp *Tabpane) istouchable(x, y int) bool {
	if y <= 39 {
		tl := 1
		for _, v := range tp.tabs {
			w := len(v.name)*8 + 12
			tl += w + 3
		}
		if x <= tl {
			return true
		}
	}
	return false
}

func (tp *Tabpane) touch(x, y int) bool {
	tabx := 1
	for i, v := range tp.tabs {
		w := len(v.name)*8 + 12
		if tabx+1 <= x && tabx+w-1 >= x {
			//tab v is touched
			if tp.currentindex != i {
				tp.tabs[tp.currentindex].hide()
				tp.currentindex = i
				declareInvalid(tp)
				v.show()
			}
			break
		}
		tabx += w + 3
	}
	return true
}

func (tp *Tabpane) stoptouch() {
}

func CreateTab(name string) *Tab {
	t := &Tab{}
	t.name = name
	t.tabpane = nil
	return t
}

func (t *Tab) AddElement(elmnt drawable) {
	r := elmnt.getRect()
	r.y += 40
	elmnt.setRect(r)
	t.elements = append(t.elements, elmnt)
	if t.tabpane != nil {
		t.tabpane.invalidateTab(t)
	}
}

func (t *Tab) RemoveElement(elmnt drawable) {
	for i, v := range t.elements {
		if elmnt == v {
			t.elements = append(t.elements[:i], t.elements[i+1:]...)
			if t.tabpane != nil {
				t.tabpane.invalidateTab(t)
			}
			return
		}
	}
}

func (t *Tab) show() {
	for _, v := range t.elements {
		AddElement(v)
	}
}

func (t *Tab) hide() {
	for _, v := range t.elements {
		RemoveElement(v)
	}
}
