package main

import (
	"fmt"
	"math"
)

type WorldMap struct {
	World
	Nations   []*Nation
	Data      map[*Coordinate]*CoordinateInfo
	rawData   []*CoordinateInfo
	coodMap   map[int]map[int]*CoordinateInfo
	Character Coordinate
}

type Coordinate struct {
	// KISS, mah dudes
	X int
	Y int
}

type CoordinateInfo struct {
	GeograficInfo
	PoliticalInfo
	Coordinate
	BasicNode
}

type GeograficInfo struct {
	Z     float64
	Biome int
	*Resources
}

var BRUSH_SIZE = 1

func (world *WorldMap) GetSurroundingAltitude(x, y, lvl int, meta map[Coordinate]float64) map[Coordinate]float64 {
	newMeta := meta
	if x > world.Width-1 || y > world.Height-1 || x < 0 || y < 0 {
		fmt.Println("Forbidden Coordinate access.")
		return nil
	}
	if lvl >= (BRUSH_SIZE - 1) {
		return map[Coordinate]float64{
			{x - 1, y + 1}: world.GetAltitude(x-1, y+1),
			{x - 1, y}:     world.GetAltitude(x-1, y),
			{x - 1, y - 1}: world.GetAltitude(x-1, y-1),
			{x + 1, y + 1}: world.GetAltitude(x+1, y+1),
			{x + 1, y}:     world.GetAltitude(x+1, y),
			{x + 1, y - 1}: world.GetAltitude(x+1, y-1),
			{x, y - 1}:     world.GetAltitude(x, y-1),
			{x, y + 1}:     world.GetAltitude(x, y+1),
		}
	}
	newMeta = world.GetSurroundingAltitude(x, y, lvl+1, newMeta)
	for i := -(BRUSH_SIZE - lvl); i <= BRUSH_SIZE; i++ {
		for k := -(BRUSH_SIZE - lvl); k <= BRUSH_SIZE; k++ {
			if i == (BRUSH_SIZE-lvl) || k == (BRUSH_SIZE-lvl) || k == -(BRUSH_SIZE-lvl) || i == -(BRUSH_SIZE-lvl) {
				newMeta[Coordinate{k, i}] = world.GetAltitude(k, i)
			}
		}
	}
	return newMeta
}
func (world *WorldMap) GetAverageSurroundingAltitude(x, y int) float64 {
	a := world.GetSurroundingAltitude(x, y, 0, nil)
	if a == nil {
		fmt.Println("Forbidden Coordinate access.")
		return 0
	}
	total := float64(0)
	for i := range a {
		total += a[i]
	}

	return total / float64(len(a))
}

func (world *WorldMap) GetAreaSameAltitudeSize(threshold float64, coord Coordinate, meta map[Coordinate]bool) (int, map[Coordinate]bool) {
	if meta == nil {
		meta = map[Coordinate]bool{}
	}
	meta[coord] = true
	if int(world.GetAltitude(coord.X, coord.Y)) == 0 {
		return 0, meta
	}

	altitudes := world.GetSurroundingAltitude(coord.X, coord.Y, BRUSH_SIZE-1, nil)
	area := 0
	for i := range altitudes {
		if !meta[i] {
			if math.Abs(world.GetAltitude(coord.X, coord.Y)-altitudes[i]) <= threshold {
				area++
			}
			other, newMeta := world.GetAreaSameAltitudeSize(threshold, i, meta)
			area += other
			meta = newMeta
		}
	}
	return area, meta
}

func (world *WorldMap) GetTotalAreaOfIsland(coord Coordinate, meta map[Coordinate]bool) (int, map[Coordinate]bool) {
	if meta == nil {
		meta = map[Coordinate]bool{}
	}
	meta[coord] = true
	if int(world.GetAltitude(coord.X, coord.Y)) == 0 {
		return 0, meta
	}

	altitudes := world.GetSurroundingAltitude(coord.X, coord.Y, BRUSH_SIZE-1, nil)
	area := 0
	for i := range altitudes {
		if !meta[i] {
			if altitudes[i] >= 0.5 {
				area++
			}
			other, newMeta := world.GetTotalAreaOfIsland(i, meta)
			area += other
			meta = newMeta
		}
	}
	return area, meta
}

func (world *WorldMap) GetNumberSizeOfSurroundingArea(coord Coordinate) int {
	altitudes := world.GetSurroundingAltitude(coord.X, coord.Y, 0, nil)
	area := 0
	for i := range altitudes {
		if altitudes[i] >= 0.5 {
			area++
		}
	}
	return area
}

func (world *WorldMap) GetAltitude(x, y int) float64 {
	coord := world.GetCoordinateInfo(x, y)
	if x > world.Width-1 || y > world.Height-1 || x < 0 || y < 0 || coord == nil {
		if coord == nil {
			// fmt.Println("SUPER Forbidden Coordinate access.")
		}
		return 0
	}
	alt := coord.Z
	return alt
}

func (world *WorldMap) GetCoordinateInfo(x, y int) *CoordinateInfo {
	if world.coodMap == nil {
		world.coodMap = map[int]map[int]*CoordinateInfo{}
		for i := range world.Data {
			if world.coodMap[i.X] == nil {
				world.coodMap[i.X] = map[int]*CoordinateInfo{}
			}
			world.coodMap[i.X][i.Y] = world.Data[i]
		}
	}

	return world.coodMap[x][y]
}

func (world *WorldMap) GetSurroundingStatistics(coordinate *Coordinate, radius float64) *CoordinateInfo {
	stats := &CoordinateInfo{
		GeograficInfo: GeograficInfo{
			Resources: &Resources{0, 0, 0, 0, 0, 0},
		},
	}
	for k := range world.Data {
		// ta uma merda, eu sei
		dist := CalculeDistanceFromCoordinate(coordinate.X, coordinate.Y, *k)
		if dist <= radius && world.Data[k].Resources != nil {
			stats.Resources.Add(*world.Data[k].Resources)
		}
	}
	return stats
}

func CalculeDistanceFromCoordinate(x, y int, c Coordinate) float64 {
	return math.Pow(float64(c.X-x), 2) + math.Pow(float64(c.Y-y), 2)
}
