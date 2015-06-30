raspui
======

Here you see raspui, a simple library for drawing a GUI on the C-Berry-Touch display for the Raspberry PI.

![raspuiTextboxes](http://rlc2.zigapeda.de/textbox.jpg)

![raspuiSlider](http://rlc2.zigapeda.de/slider.jpg)

It's designed to draw a lightweight GUI to the display without the need of any framebuffer drivers or x installations. Currently there are four elements which can be used:
 * Textbox
 * Button
 * Slider
 * Tabpane

Additionally there are some features:
 * Standby mode (display turns off and waits for touch to turn back on)
 * Backlight brightness is adjustable

Planned features are:
 * Controll the GPIO pins on the display adapter board with the library
 * Spinner to adjust values precisely

For building the framework you need:
 * Working Go environment
 * GCC
 * libbcm2835
 * C-Berry-Touch sourcecode files RAIO8870.h, RAIO8870.c, tft.h, tft.c from admatec which are delivered with the display.

Example

```go
package main

import (
        "github.com/zigapeda/raspui"
        "fmt"
        "strconv"
)

var bright int = 255
var light *raspui.Textbox
var slid *raspui.Slider

func setlight(v int) {
        bright = v
        light.SetText(strconv.Itoa(bright))
        raspui.SetDisplayBacklight(bright)
}

func main() {
        err := raspui.CreateRaspUI()
        defer raspui.CloseRaspUI()
        if err != nil {
                fmt.Println(err)
                return
        }

        exit := make(chan int)

        tabpane := raspui.CreateTabpane()
        tab1 := raspui.CreateTab("Stats")
        tab1.AddElement(raspui.CreateTextbox(5, 2, 150, 20, "Channel 1: 120 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 27, 150, 20, "Channel 2:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 52, 150, 20, "Channel 3:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 77, 150, 20, "Channel 4:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 102, 150, 20, "Channel 5:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 127, 150, 20, "Channel 6:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 152, 150, 20, "Channel 7:  90 °C"))
        tab1.AddElement(raspui.CreateTextbox(5, 177, 150, 20, "Channel 8:  90 °C"))
        tab2 := raspui.CreateTab("Settings")
        slid = raspui.CreateSlider(130, 50, 170, 20, 2, 255)
        slid.SetChangeFunc(setlight)
        light = raspui.CreateTextbox(130, 20, 100, 20, "255")
        tab2.AddElement(slid)
        tab2.AddElement(light)
        tabpane.AddTab(tab1)
        tabpane.AddTab(tab2)
        raspui.AddElement(tabpane)


        raspui.AddElement(raspui.CreateButton(278, 1, 40, 35, " X ", func() {exit <- 1}))
        raspui.AddElement(raspui.CreateButton(236, 1, 40, 35, " _ ", func() {raspui.Standby()}))
        raspui.AddElement(raspui.CreateButton(194, 1, 40, 35, " T ", func() {light.SetText("Test Pressed")}))

        <- exit
}
```
