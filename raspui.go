// raspui project raspui.go
package raspui

import (
	"fmt"
	"image"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	standbyOff          = 0
	standbyActivating   = 1
	standbyOn           = 2
	standbyDeactivating = 3
)

var xoff, yoff int
var xsize, ysize float64
var elements []drawable
var standby int
var touchticker *time.Ticker
var gpioActionChan chan func()
var invalidChan chan drawable
var lastx, lasty int
var touchelement touchable = nil
var bright = 255

func CreateRaspUI() error {
	//load calibration
	arr, err := ioutil.ReadFile("calibration.conf")
	if err != nil {
		return err
	}

	s := strings.Split(string(arr), ",")
	xoff, err = strconv.Atoi(s[0])
	if err != nil {
		return err
	}

	yoff, err = strconv.Atoi(s[1])
	if err != nil {
		return err
	}

	xsize, err = strconv.ParseFloat(s[2], 64)
	if err != nil {
		return err
	}

	ysize, err = strconv.ParseFloat(s[3], 64)
	if err != nil {
		return err
	}

	standby = standbyOff

	initDisplay()
	clearDisplay(WHITE)

	elements = make([]drawable, 0, 10)

	touchticker = time.NewTicker(time.Millisecond * 25)
	gpioActionChan = make(chan func(), 100)
	invalidChan = make(chan drawable, 100)
	go uiLoop()

	//btn := CreateButton(20, 20, 100, 30, "test", func(x, y int) { fmt.Println(x, y) })
	//AddElement(btn)

	return nil
}

func Standby() {
	standby = standbyActivating
	invokeLater(func() { setBacklightPwmValue(0) })
}

func CloseRaspUI() {
	closeGpio()
}

func SetDisplayBacklight(brightness int) {
	bright = brightness
	invokeLater(func() { setBacklightPwmValue(brightness) })
}

func AddElement(elmnt drawable) {
	elements = append(elements, elmnt)
	invokeLater(elmnt.draw)
}

func RemoveElement(elmnt drawable) {
	for i, v := range elements {
		if elmnt == v {
			r := elmnt.getRect()
			invokeLater(func() { drawFilledBox(r.x, r.y, r.x+r.width, r.y+r.height, WHITE) }) //remove from ui
			elements = append(elements[:i], elements[i+1:]...)
		}
	}
}

func convertToPixel(p *image.Point) (x, y int) {
	x = int(float64(p.X-xoff) * xsize)
	y = int(float64(p.Y-yoff) * ysize)
	if x < 0 && y < 0 {
		return lastx, lasty
	}
	lastx, lasty = x, y
	return x, y
}

func invokeLater(function func()) {
	gpioActionChan <- function
}

func declareInvalid(element drawable) {
	invalidChan <- element
}

/*func touchDetection() {
	touchstate := getTouchState()
	if touchstate != NOTOUCH {
		p := getTouchXY()
		x, y := convertToPixel(p)
		//fmt.Println("touch ", x, y)
		delivered := false
		for _, v := range elements {
			if t, ok := v.(touchable); ok {
				//all touchables on the display -> find the touched one
				r := v.getRect()
				if r.x <= x && r.y <= y && r.x+r.width >= x && r.y+r.height >= y {
					//current element is touched
					if t != touchelement {
						//new element is pressed
						if touchelement != nil {
							touchelement.stoptouch()
						}
						touchelement = t
					}
					//send touchcoordinates
					touchelement.touch(p.X, p.Y)
					delivered = true
					break
				}
			}
		}
		if delivered == false && touchelement != nil {
			touchelement.stoptouch()
			touchelement = nil
		}
	} else {
		if touchelement != nil {
			touchelement.stoptouch()
			touchelement = nil
		}
	}
}*/

func touchDetection() {
	touchstate := getTouchState()
	if standby == standbyOff {
		if touchstate != NOTOUCH {
			p := getTouchXY()
			x, y := convertToPixel(p)
			if touchelement != nil {
				if touchelement.touch(x, y) == false {
					touchelement.stoptouch()
					touchelement = nil
				}
			}

			if touchelement == nil {
				for _, v := range elements {
					if t, ok := v.(touchable); ok {
						//all touchables on the display -> find the touched one
						if t.istouchable(x, y) {
							if t.touch(x, y) == true {
								touchelement = t
							}
							break
						}
					}
				}
			}
		} else {
			if touchelement != nil {
				touchelement.stoptouch()
				touchelement = nil
			}
		}
	} else if standby == standbyActivating {
		if touchstate == NOTOUCH {
			standby = standbyOn
			if touchelement != nil {
				touchelement.stoptouch()
				touchelement = nil
			}
		}
	} else if standby == standbyOn {
		if touchstate != NOTOUCH {
			standby = standbyDeactivating
			invokeLater(func() { setBacklightPwmValue(bright) })
		}
	} else if standby == standbyDeactivating {
		if touchstate == NOTOUCH {
			standby = standbyOff
		}
	}
}

func uiLoop() {
	for {
		select {
		case f := <-gpioActionChan:
			f()
		case i := <-invalidChan:
			for _, v := range elements {
				if i == v {
					i.draw()
					break
				}
			}
		case <-touchticker.C:
			touchDetection()
		}
	}
}

func Calibrate() {
	initDisplay()
	clearDisplay(WHITE)

	//werte ermitteln
	c1, c2, tx1, tx2, ty1, ty2 := 0, 0, 0, 0, 0, 0
	drawFilledBox(18, 18, 22, 22, BLUE)
	for getTouchState() == NOTOUCH {
		time.Sleep(time.Millisecond * 100)
	}
	for getTouchState() != NOTOUCH {
		p := getTouchXY()
		c1++
		tx1 += p.X
		ty1 += p.Y
		time.Sleep(time.Millisecond * 50)
	}
	clearDisplay(WHITE)
	drawFilledBox(298, 218, 302, 222, BLUE)
	for getTouchState() == NOTOUCH {
		time.Sleep(time.Millisecond * 100)
	}
	for getTouchState() != NOTOUCH {
		p := getTouchXY()
		c2++
		tx2 += p.X
		ty2 += p.Y
		time.Sleep(time.Millisecond * 50)
	}
	clearDisplay(WHITE)

	//durchschnitt ausrechnen
	tx1 = tx1 / c1
	tx2 = tx2 / c2
	ty1 = ty1 / c1
	ty2 = ty2 / c2

	//touchgroesse in pixel ausrechnen
	tx := 280.0 / float64(tx2-tx1)
	ty := 200.0 / float64(ty2-ty1)

	//offset ausrechnen
	offx := tx1 - int(20/tx)
	offy := ty1 - int(20/ty)

	//ausgeben
	fmt.Println("Touchsize X, Y", tx, ty, "Offset X, Y", offx, offy)
	arr := []byte(strconv.Itoa(offx) + "," + strconv.Itoa(offy) +
		"," + strconv.FormatFloat(tx, 'f', 2, 64) +
		"," + strconv.FormatFloat(ty, 'f', 2, 64))
	ioutil.WriteFile("calibration.conf", arr, 0644)
}
