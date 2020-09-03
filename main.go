package main

import (
	"fmt"
	gui "terminalgame/gui"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()
	w := app.NewWindow("Hello")
	log := gui.NewLogBox(w)
	info := gui.NewInfoBox(w)
	broker := Broker{}
	broker.SubscribeToEvent("ShowInfo", info.ShowInfo)
	loadWorld, imageFunc, setParams, getParams := InitWorldLoader(log.Log, broker)
	worldVisualizer := gui.NewWorldVisualizer(imageFunc)
	widthSize := gui.NewEnterEntry(func(s string) { setParams("WORLD_WIDTH", s) }, fmt.Sprintf("%v", getParams("WORLD_WIDTH")))
	heightSize := gui.NewEnterEntry(func(s string) { setParams("WORLD_HEIGHT", s) }, fmt.Sprintf("%v", getParams("WORLD_HEIGHT")))
	keyDownListener := func(f *fyne.KeyEvent) { broker.EmitEvent("KeyDownEvent", f); worldVisualizer.Refresh() }
	w.Canvas().(desktop.Canvas).SetOnKeyDown(keyDownListener)
	w.SetContent(widget.NewHBox(
		widget.NewVBox(
			widthSize,
			heightSize,
			info,
		),
		widget.NewVBox(
			gui.NewHeader(app),
			widget.NewGroup(""),
			gui.NewLoader(log, func(s string) { loadWorld(s, true); worldVisualizer.Refresh() }, func(s string) { loadWorld(s, false); worldVisualizer.Refresh() }),
			worldVisualizer,
		),
		log,
	))

	w.ShowAndRun()
}

func addFuncToKeyCallBack(f1, f2 func(*fyne.KeyEvent)) func(*fyne.KeyEvent) {
	return func(a *fyne.KeyEvent) {
		f1(a)
		f2(a)
	}
}

// TODO: Add Game time clocks
// TODO: Add Button "Advance Step" and "Play"
// TODO: Add food production, stock and consume
// TODO: Add Money
// TODO: Add trade ( buy from someone's stock )
// TODo: add product production, stock and consume
