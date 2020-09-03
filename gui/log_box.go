package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var first = 0

type LogBox struct {
	widget.BaseWidget
	logs   []string
	lines  []fyne.CanvasObject
	window fyne.Window
}

type LogBoxRenderer struct {
	container   *fyne.Container
	linesHolder *widget.Box
	logBox      *LogBox
}

func (l *LogBox) Log(s string) {
	fmt.Println(s)
	l.ExtendBaseWidget(l)
	l.logs = append(l.logs, s)
	t := canvas.NewText(s, color.White)
	t.Alignment = fyne.TextAlignLeading
	l.lines = append(l.lines, t)
	scrollCont := widget.Renderer(l).(*LogBoxRenderer).container.Objects[0].(*widget.ScrollContainer)
	externalCont := scrollCont.Content.(*widget.Box)
	externalCont.Children = l.lines
	scrollCont.Scrolled(&fyne.ScrollEvent{
		DeltaY: -600,
	})
	widget.Refresh(externalCont)
	scrollCont.Scrolled(&fyne.ScrollEvent{
		DeltaY: -600,
	})
}

func (l *LogBox) MinSize() fyne.Size {
	// Change to the size of the window
	return fyne.Size{400, 500}
}

func (l *LogBox) Refresh() {
}

func (l *LogBoxRenderer) MinSize() fyne.Size {
	return fyne.Size{400, 500}
}

func (l *LogBoxRenderer) MaxSize() fyne.Size {
	return fyne.Size{400, 500}
}
func NewLogBox(w fyne.Window) *LogBox {
	entry := &LogBox{}
	entry.lines = []fyne.CanvasObject{}
	entry.window = w
	entry.ExtendBaseWidget(entry)
	return entry
}

func (l *LogBox) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	renderer := &LogBoxRenderer{logBox: l}
	box := widget.NewVBox(l.lines...)
	cont := widget.NewVScrollContainer(box)
	renderer.container = fyne.NewContainerWithLayout(layout.NewHBoxLayout(), cont)
	renderer.linesHolder = box
	return renderer
}

func (l *LogBoxRenderer) Refresh() {
	canvas.Refresh(l.container)
}

func (l *LogBoxRenderer) ApplyTheme() {
}

func (l *LogBoxRenderer) BackgroundColor() color.Color {
	return color.White
}

func (l *LogBoxRenderer) Objects() []fyne.CanvasObject {
	return l.container.Objects
}

func (l *LogBoxRenderer) Destroy() {
}

func (l *LogBoxRenderer) Layout(size fyne.Size) {
	headerHeight := l.logBox.MinSize().Height
	gridSize := size.Add(fyne.NewSize(0, headerHeight))
	scrollSize := gridSize.Min(l.MaxSize())
	l.container.Resize(scrollSize)
	l.container.Objects[0].Resize(scrollSize)
	l.logBox.Resize(scrollSize)
}
