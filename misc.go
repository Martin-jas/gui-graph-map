package main

import (
	"encoding/json"
	"fmt"
	fcolor "image/color"
	"math"

	"fyne.io/fyne/theme"
	"github.com/fatih/color"
)

func GenerateWorldCoordinates(w, h int) []Coordinate {
	fmt.Printf("Generating world with '%v' width and '%v' height.\n", w, h)
	a := []Coordinate{}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			a = append(a, Coordinate{j, i})
		}
	}
	return a
}

func PrintCoordinates(world *WorldMap, log func(string)) {
	w := world.Width
	h := world.Height
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("_")
	}
	fmt.Println(" ")
	for i := 0; i < h; i++ {
		fmt.Printf("| ")
		for j := 0; j < w; j++ {
			z := int(math.Round(world.GetAltitude(j, i)))
			var c *color.Color
			if z == 0 {
				c = color.New(color.FgBlue)
			} else if z == -1 {
				c = color.New(color.FgRed)
				z = 0
			} else {
				c = color.New(color.FgGreen)
			}
			c.Printf("%v", z)

			if j < w-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println(" |")
	}
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("‾")
	}
	fmt.Println(" ")
}

func PrintCoordinatesWithRichness(world *WorldMap) {
	w := world.Width
	h := world.Height
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("_")
	}
	fmt.Println(" ")
	for i := 0; i < h; i++ {
		fmt.Printf("| ")
		for j := 0; j < w; j++ {
			tile := world.GetCoordinateInfo(j, i)
			z := int(math.Round(world.GetAltitude(j, i)))
			isFoodT := tile != nil && tile.Resources != nil && tile.Food > tile.RawProduct
			isProductT := tile != nil && tile.Resources != nil && tile.RawProduct > 0
			var c *color.Color
			if z == 0 {
				c = color.New(color.FgBlue)
			} else {
				c = color.New(color.FgGreen)
			}
			if isFoodT && !isProductT {
				c = color.New(color.FgHiGreen)
				c = c.Add(color.Bold)
			} else if isProductT {
				c = c.Add(color.BgCyan)
			}
			c.Printf("%v", z)

			if j < w-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println(" |")
	}
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("‾")
	}
	fmt.Println(" ")
}

func PrintCoordinatesWithNations(world *WorldMap) {
	w := world.Width
	h := world.Height
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("_")
	}
	fmt.Println(" ")
	for i := 0; i < h; i++ {
		fmt.Printf("| ")
		for j := 0; j < w; j++ {
			z := int(math.Round(world.GetAltitude(j, i)))
			var c *color.Color
			if z == 0 {
				c = color.New(color.FgBlue)
			} else {
				c = color.New(color.FgGreen)
			}
			if world.GetCoordinateInfo(j, i).owner != nil {
				c = color.New(color.Attribute(world.GetCoordinateInfo(j, i).owner.Index) + 93)
				c.Printf("V")
			} else {
				c.Printf("%v", z)
			}

			if j < w-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println(" |")
	}
	fmt.Printf(" ")
	for b := 0; b < 2*w+1; b++ {
		fmt.Printf("‾")
	}
	fmt.Println(" ")
}

func prettyPrint(i interface{}) {
	amd, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(amd))
}

func checkerPatternFunc(x, y, w, h int) fcolor.Color {
	x /= (w / 10)
	y /= (w / 10)

	if x%2 == y%2 {
		return theme.BackgroundColor()
	}

	return theme.ButtonColor()
}

func worldPixel(x, y, w, h int, world *WorldMap) fcolor.Color {
	wResolution := float64(world.Width) / float64(w)
	hResolution := float64(world.Height) / float64(h)

	xAdjusted := float64(x) * wResolution
	yAdjusted := float64(y) * hResolution

	coord := world.GetCoordinateInfo(int(xAdjusted), int(yAdjusted))
	if coord == nil {
		return fcolor.Black
	}
	if world.Character.X == coord.X && world.Character.Y == coord.Y {
		return CivColors[5] // Draw "Character"
	}
	if coord.Z <= 0.8 {
		return fcolor.RGBA{35, 99, 163, 255}
	}
	if coord.GetOwner() == nil {
		alt := uint8(int(coord.Z-1) * 10)
		return fcolor.RGBA{102 - alt, 153 - alt, 102 - alt, 255}
	}
	return CivColors[coord.GetOwner().Index]
}

var CivColors = []fcolor.RGBA{
	//red
	{
		250, 15, 15, 255,
	},
	//green
	{
		45, 255, 45, 255,
	},
	//purple
	{
		150, 45, 255, 255,
	},
	//pink
	{
		255, 45, 209, 255,
	},
	//white
	{
		255, 255, 255, 255,
	},
	//black
	{
		0, 0, 0, 255,
	},
}
