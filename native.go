// raspui project raspui.go
package raspui

/*
#cgo LDFLAGS: -lbcm2835
#include <bcm2835.h>
#include "RAIO8870.h"
#include "tft.h"

uint8_t getTouchState() {
	RAIO_gettouch();
	if(my_touch.state == pressed || my_touch.state == down) {
		return 1;
	}
	return 0;
}

uint32_t getTouchX() {
	return my_touch.touch_x;
}

uint32_t getTouchY() {
	return my_touch.touch_y;
}
*/
import "C"

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	"image"
	"unsafe"
)

const (
	RED = int(C.COLOR_RED)
	BLUE = int(C.COLOR_BLUE)
	GREEN = int(C.COLOR_GREEN)
	BLACK = int(C.COLOR_BLACK)
	WHITE = int(C.COLOR_WHITE)
	CYAN = int(C.COLOR_CYAN)
	YELLOW = int(C.COLOR_YELLOW)
	MAGENTA = int(C.COLOR_MAGENTA)
	DARKGREEN = int(C.COLOR_DARK_GREEN)
	DARKBLUE  = 0x01
	LIGHTBLUE = 0xB7
	GRAY      = 0x92
)

const (
	PRESSED = 1
	NOTOUCH = 3
)

func getTouchState() int {
	if int(C.getTouchState()) == 1 {
		return PRESSED
	}
	return NOTOUCH
}

func getTouchXY() *image.Point {
	x := int(C.getTouchX())
	y := int(C.getTouchY())
	return &image.Point{x, y}
}

func clearDisplay(color int) {
	C.Text_Background_Color (C.uint8_t(color))
	C.RAIO_clear_screen()	
}

func drawLine(x1, y1, x2, y2, color int) {
	C.Text_Foreground_Color (C.uint8_t(color))
	C.Set_Geometric_Coordinate(C.uint16_t(x1), C.uint16_t(y1), C.uint16_t(x2), C.uint16_t(y2))
	C.RAIO_StartDrawing(C.LINE)
}

func drawBox(x1, y1, x2, y2, color int) {
	C.Text_Foreground_Color (C.uint8_t(color))
	C.Set_Geometric_Coordinate(C.uint16_t(x1), C.uint16_t(y1), C.uint16_t(x2), C.uint16_t(y2))
	C.RAIO_StartDrawing(C.SQUARE_NONFILL)
}

func drawFilledBox(x1, y1, x2, y2, color int) {
	C.Text_Foreground_Color (C.uint8_t(color))
	C.Set_Geometric_Coordinate(C.uint16_t(x1), C.uint16_t(y1), C.uint16_t(x2), C.uint16_t(y2))
	C.RAIO_StartDrawing(C.SQUARE_FILL)
}

func drawCircle(x, y, rad, color int) {
	C.Text_Foreground_Color (C.uint8_t(color))
	C.Set_Geometric_Coordinate_circle(C.uint16_t(x), C.uint16_t(y), C.uint8_t(rad))
	C.RAIO_StartDrawing(C.CIRCLE_NONFILL)
}

func drawFilledCircle(x, y, rad, color int) {
	C.Text_Foreground_Color (C.uint8_t(color))
	C.Set_Geometric_Coordinate_circle(C.uint16_t(x), C.uint16_t(y), C.uint8_t(rad))
	C.RAIO_StartDrawing(C.CIRCLE_FILL)
}

func drawText(x, y, size int, text string, fg, bg int) error {
	//RAIO_print_text( (DISPLAY_WIDTH/4)-45, (DISPLAY_HEIGHT/4)*3-20, "Button 3 is", COLOR_GREEN, COLOR_BLACK );

	buf := new(bytes.Buffer)
	w, err := charset.NewWriter("latin1", buf)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, text)
	w.Close()

	b := []byte(buf.String())
	b = append(b, 0)
	ctext := (*C.uchar)(unsafe.Pointer(&b[0]))
	C.RAIO_SetFontSizeFactor(C.uint8_t(size))
	C.RAIO_print_text(C.uint16_t(x), C.uint16_t(y), ctext, C.uint8_t(bg), C.uint8_t(fg))
	return nil
}

func initDisplay() {
	C.bcm2835_init()
	C.TFT_init_board()
	C.TFT_hard_reset()
	C.RAIO_init()
	C.RAIO_SetBacklightPWMValue(255)
}

func closeGpio() {
	C.TFT_hard_reset()
	C.RAIO_SetBacklightPWMValue(0)
	C.bcm2835_close()
}

func setBacklightPwmValue(brightness int) {
	C.RAIO_SetBacklightPWMValue(C.uint8_t(brightness))
}

func SetGpioOutput(pin int) {
	fmt.Println("set pin", pin, "to output")
}

func SetGpioInput(pin, mode int) {
	fmt.Println("set pin", pin, "to input")
}

func SetGpio(pin, level int) {
	fmt.Println("set pin", pin, "to", level)
}
