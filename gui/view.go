package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type enterEntry struct {
	widget.Entry
	OnSubmit  func(string)
	innerText string
	KeyDownCB func(*fyne.KeyEvent)
}

func (e *enterEntry) onEnter() {
	e.OnSubmit(e.Entry.Text)
}

func (e *enterEntry) MinSize() fyne.Size {
	return fyne.NewSize(200, e.Entry.MinSize().Height)
}

func (e *enterEntry) KeyDown(key *fyne.KeyEvent) {
	if e.Entry.Focused() {
		switch key.Name {
		case fyne.KeyReturn:
			e.onEnter()
		default:
			e.Entry.KeyDown(key)
		}
	}
}

func (e *enterEntry) OnChanged(text string) {
	fmt.Println(text)
	e.innerText = text
}

func NewEnterEntry(onSubmit func(string), startingTxt string) *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Entry.OnChanged = entry.OnChanged
	entry.OnSubmit = onSubmit
	entry.Entry.SetText(startingTxt)
	return entry
}

func NewHeader(app fyne.App) *widget.Box {
	return widget.NewHBox(
		widget.NewLabel("World Visualizer"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	)
}

type WorldImgHolder struct {
	widget.BaseWidget
	imgGen func(x, y, w, h int) color.Color
}

type WorldRenderer struct {
	w         *WorldImgHolder
	container *fyne.Container
}

func (w *WorldImgHolder) MinSize() fyne.Size {
	return fyne.NewSize(400, 400)
}

func NewWorldVisualizer(imageFunction func(x, y, w, h int) color.Color) *WorldImgHolder {
	w := &WorldImgHolder{}
	w.ExtendBaseWidget(w)
	w.imgGen = imageFunction
	return w
}

func (w *WorldImgHolder) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	renderer := &WorldRenderer{}
	renderer.container = fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(w.MinSize().Width, w.MinSize().Height)), canvas.NewRasterWithPixels(w.imgGen))
	renderer.w = w
	return renderer
}

func (w *WorldRenderer) MinSize() fyne.Size {
	return fyne.Size{400, 400}
}

func (w *WorldRenderer) Refresh() {
	canvas.Refresh(w.w)
}

func (w *WorldRenderer) ApplyTheme() {
}

func (w *WorldRenderer) BackgroundColor() color.Color {
	return color.White
}

func (w *WorldRenderer) Objects() []fyne.CanvasObject {
	return w.container.Objects
}

func (w *WorldRenderer) Destroy() {
}

func (w *WorldRenderer) Layout(size fyne.Size) {
}
