package gui

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type InfoBox struct {
	widget.BaseWidget
	info   []string
	lines  []fyne.CanvasObject
	window fyne.Window
}

type InfoBoxRenderer struct {
	container   *fyne.Container
	linesHolder *widget.Box
	infoBox     *InfoBox
}

func (l *InfoBox) ShowInfo(e string, payload interface{}) {
	s := payload.([]string)
	l.ExtendBaseWidget(l)
	l.info = s
	l.lines = []fyne.CanvasObject{}
	for inf := range s {
		t := canvas.NewText(s[inf], color.White)
		t.Alignment = fyne.TextAlignLeading
		l.lines = append(l.lines, t)
	}
	scrollCont := widget.Renderer(l).(*InfoBoxRenderer).container.Objects[0].(*widget.ScrollContainer)
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

func (l *InfoBox) MinSize() fyne.Size {
	// Change to the size of the window
	return fyne.Size{400, 200}
}

func (l *InfoBox) Refresh() {
}

func (l *InfoBoxRenderer) MinSize() fyne.Size {
	return fyne.Size{400, 500}
}

func (l *InfoBoxRenderer) MaxSize() fyne.Size {
	return fyne.Size{400, 500}
}
func NewInfoBox(w fyne.Window) *InfoBox {
	entry := &InfoBox{}
	entry.lines = []fyne.CanvasObject{}
	entry.window = w
	entry.ExtendBaseWidget(entry)
	return entry
}

func (l *InfoBox) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	renderer := &InfoBoxRenderer{infoBox: l}
	box := widget.NewVBox(l.lines...)
	cont := widget.NewVScrollContainer(box)
	renderer.container = fyne.NewContainerWithLayout(layout.NewHBoxLayout(), cont)
	renderer.linesHolder = box
	return renderer
}

func (l *InfoBoxRenderer) Refresh() {
	canvas.Refresh(l.container)
}

func (l *InfoBoxRenderer) ApplyTheme() {
}

func (l *InfoBoxRenderer) BackgroundColor() color.Color {
	return color.White
}

func (l *InfoBoxRenderer) Objects() []fyne.CanvasObject {
	return l.container.Objects
}

func (l *InfoBoxRenderer) Destroy() {
}

func (l *InfoBoxRenderer) Layout(size fyne.Size) {
	headerHeight := l.infoBox.MinSize().Height
	gridSize := size.Add(fyne.NewSize(0, headerHeight))
	scrollSize := gridSize.Min(l.MaxSize())
	l.container.Resize(scrollSize)
	l.container.Objects[0].Resize(scrollSize)
	l.infoBox.Resize(scrollSize)
}
