package raspui

type Slider struct {
	rect
	pressed bool
	minx    int
	maxx    int
	pos     int
	lastpos int
	offset  int
	step    float64
	change  func(int)
	release func(int)
}

func CreateSlider(x, y, width, height, minvalue, maxvalue int) *Slider {
	sld := &Slider{}
	sld.x = x
	sld.y = y
	sld.width = width
	sld.height = height
	sld.minx = x + height/2
	sld.maxx = x + width - height/2
	sld.pos = x + height/2
	sld.lastpos = sld.pos
	difv := maxvalue - minvalue
	difx := sld.maxx - sld.minx
	sld.offset = minvalue
	sld.step = float64(difv) / float64(difx)
	sld.change = nil
	sld.release = nil
	return sld
	//sld ist bla
}

func (s *Slider) draw() {
	ly := s.y + (s.height / 2)
	drawFilledBox(s.x, s.y, s.x+s.width, s.y+s.height, WHITE)
	drawLine(s.minx, ly, s.maxx, ly, GRAY)
	drawCircle(s.pos, ly, s.height/2, DARKBLUE)
	if s.pressed == true {
		drawFilledCircle(s.pos, ly, s.height/2-1, LIGHTBLUE)
	} else {
		drawFilledCircle(s.pos, ly, s.height/2-1, BLUE)
	}
}

func (s *Slider) getRect() rect {
	return s.rect
}

func (s *Slider) setRect(r rect) {
	s.rect = r
}

func (s *Slider) istouchable(x, y int) bool {
	return s.x <= x && s.y <= y && s.x+s.width >= x && s.y+s.height >= y
}

func (s *Slider) touch(x, y int) bool {
	if s.pressed == false {
		s.pressed = true
		invokeLater(s.draw)
	}
	if x < s.minx {
		s.pos = s.minx
	} else if x > s.maxx {
		s.pos = s.maxx
	} else {
		s.pos = x
	}
	if s.lastpos+3 < s.pos || s.lastpos-3 > s.pos {
		if s.change != nil {
			s.change(s.GetValue())
		}
		invokeLater(s.draw)
		s.lastpos = s.pos
	}
	return true
}

func (s *Slider) stoptouch() {
	s.pressed = false
	s.lastpos = s.pos
	if s.release != nil {
		s.release(s.GetValue())
	}
	invokeLater(s.draw)
}

func (s *Slider) SetChangeFunc(change func(int)) {
	s.change = change
}

func (s *Slider) SetReleaseFunc(release func(int)) {
	s.release = release
}

func (s *Slider) SetValue(value int) {
	s.pos = int(float64(value-s.offset)/s.step) + s.minx
	s.lastpos = s.pos
	declareInvalid(s)
}

func (s *Slider) GetValue() int {
	return int(float64(s.pos-s.minx)*s.step) + s.offset
}
