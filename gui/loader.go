package gui

import (
	"fmt"

	"fyne.io/fyne/widget"
)

type Loader struct {
	*widget.Box
}

func NewLoader(log *LogBox, cbSave func(string), cbGenerateOnly func(string)) *Loader {
	worldNameEntry := NewEnterEntry(cbSave, "testera")
	a := widget.NewHBox(widget.NewButton("Load/Create", func() {
		if worldNameEntry.innerText != "" {
			log.Log(fmt.Sprintf("Loading [%v] world...", worldNameEntry.innerText))
			if cbSave != nil {
				cbSave(worldNameEntry.innerText)
			}
		}
	}),
		widget.NewButton("Generate", func() {
			if worldNameEntry.innerText != "" {
				log.Log(fmt.Sprintf("Generating [%v] world...", worldNameEntry.innerText))
				if cbGenerateOnly != nil {
					cbGenerateOnly(worldNameEntry.innerText)
				}
			}
		}),
		worldNameEntry)
	return &Loader{
		a,
	}
}
