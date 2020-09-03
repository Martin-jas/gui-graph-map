package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func FromCoordinateInfoToWorldMap(coords []CoordinateInfo, wn World) *WorldMap {
	world := WorldMap{wn, nil, map[*Coordinate]*CoordinateInfo{}, nil, nil, Coordinate{0, 0}}
	if len(world.rawData) == 0 {
		world.rawData = []*CoordinateInfo{}
	}
	for i := range coords {
		world.rawData = append(world.rawData, &coords[i])
		world.Data[&(coords[i].Coordinate)] = &coords[i]
	}
	return &world
}

func TryToGenerateZBasedOnPlacas(coords []CoordinateInfo, wn World) *WorldMap {
	world := FromCoordinateInfoToWorldMap(coords, wn)
	num := 20
	for i := 0; i < num; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		x := rand.Intn(wn.Width)
		rand.Seed(int64(time.Now().Nanosecond()))
		y := rand.Intn(wn.Height)
		point := world.GetCoordinateInfo(x, y)
		point.Z = 1

	}
	return world
}

func GenerateWorldBasedOnCoordinates(coords []CoordinateInfo, wn World) *WorldMap {
	// return TryToGenerateZBasedOnPlacas(coords, wn)
	world := FromCoordinateInfoToWorldMap(coords, wn)
	s := time.Now()
	for p := 0; p < world.NumIterations; p++ {
		for k := range world.Data {
			rand.Seed(int64(time.Now().Nanosecond()))
			step := rand.Intn(3) - 1
			point := world.Data[k]
			avarageZ := world.GetAverageSurroundingAltitude(k.X, k.Y)
			if avarageZ <= 0 && step < 0 {
				step = 0
			}
			point.Z = math.Max(avarageZ+(float64(step)*world.Smoothness), 0) // chaotic way
		}
	}
	fmt.Println(fmt.Sprintf("\n[Z-Lvl] = Time taken - %v", time.Since(s)))
	// fmt.Println("After initial generation")
	// PrintCoordinates(world)
	s = time.Now()
	// Remove too small islands
	islands := []map[Coordinate]bool{}
	for k := range world.Data {
		flag := false
		for island := range islands {
			if islands[island][*k] {
				flag = true
			}
		}
		if flag {
			continue
		}
		islandSize, islandPoints := world.GetTotalAreaOfIsland(*k, nil)
		islands = append(islands, islandPoints) // TODO: check how to save this info
		if islandSize < world.MinimumIslandSize && islandSize != 0 {
			for t := range islandPoints {
				c := world.GetCoordinateInfo(t.X, t.Y)
				if c != nil {
					c.Z = 0
				} else {
					fmt.Println(fmt.Sprintf("Nil coordinate at[%v, %v]", t.X, t.Y))
				}
			}
		}
	}
	fmt.Println(fmt.Sprintf("\n [Cleaning-Islands] = Time taken - %v", time.Since(s)))
	s = time.Now()
	// fmt.Println("After cleaning islands")
	// PrintCoordinates(world)
	// Remove weird "listras"
	for k := range world.Data {
		area := world.GetNumberSizeOfSurroundingArea(*k)
		if area < 3 && area != 0 {
			world.Data[k].Z = 0
		}
		// Remove close water spac
		if area >= 6 && world.GetAltitude(k.X, k.Y) == 0 {
			world.Data[k].Z = 1
		}
	}
	// fmt.Println("After cleaning water in earth")
	// PrintCoordinates(world)
	fmt.Println(fmt.Sprintf("\n [Cleaning-No islands] = Time taken - %v", time.Since(s)))
	return world
}

func GenerateResources(world *WorldMap, height_threshold float64) *WorldMap {
	s := time.Now()
	// lets base for now the resource distribution on altitude only
	var lastArea map[Coordinate]bool
	for k := range world.Data {
		if world.Data[k].Resources == nil {
			world.Data[k].Resources = &Resources{}
		}
		if !lastArea[*k] {
			lastArea = map[Coordinate]bool{}
		}
		// mountains have rawProducts
		notSea := world.GetAltitude(k.X, k.Y) > 0.5
		diff := world.GetAverageSurroundingAltitude(k.X, k.Y) - world.GetAltitude(k.X, k.Y) - 0.3
		enoughHeightDifference := diff > 0
		if notSea && enoughHeightDifference {
			world.Data[k].RawProduct += int(math.Abs(diff) * 1000)
		}
		// hills have food
		if math.Abs(diff+0.3) < height_threshold {
			var area int
			area, lastArea = world.GetAreaSameAltitudeSize(height_threshold, *k, lastArea)
			world.Data[k].Food += int(area * 10)
		}
	}
	fmt.Println(fmt.Sprintf("\n [Resources] = Time taken - %v", time.Since(s)))
	return world
}

func GenerateInitialNations(world *WorldMap) *WorldMap {
	s := time.Now()
	suitable := map[*Coordinate]*CoordinateInfo{}
	for k := range world.Data {
		statis := world.GetSurroundingStatistics(k, world.InitCivRadius)
		if statis.Resources.RawProduct > 0 && statis.Resources.Food > 0 {
			// fmt.Printf("Suitable place: (%v, %v). RawProduct: %v, Food: %v.\n", k.X, k.Y, statis.Resources.RawProduct, statis.Resources.Food)
			statis.BasicNode = world.Data[k].BasicNode
			statis.Coordinate = world.Data[k].Coordinate
			suitable[k] = statis
		}
	}
	// With the suitable places here we can create civs
	ns := []*Nation{}
	for i := 0; i < world.NumCivs; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		temp := Nation{
			Name:       fmt.Sprintf("[Civ %v]", i),
			Population: rand.Intn(16)*1000 + 1000,
			Index:      i,
			BasicNode:  BasicNode{"Nation", nil},
		}
		ns = append(ns, &temp)
	}
	for n := range ns {
		fmt.Printf("%v - Population: %v.\n", ns[n].Name, ns[n].Population)
		cityPopulations := []int{}
		maxPop := 0
		for leftPopulation := ns[n].Population; leftPopulation > 0; {
			cityPop := leftPopulation
			rand.Seed(int64(time.Now().Nanosecond()))
			cityPop = rand.Intn(leftPopulation)
			if cityPop < 1000 || leftPopulation < 1000 {
				cityPop = leftPopulation
			}
			cityPopulations = append(cityPopulations, cityPop)
			if cityPop > maxPop {
				maxPop = cityPop
			}
			leftPopulation -= cityPop
		}

		cities := []City{}
		for t := range cityPopulations {
			for k := range suitable {
				if suitable[k].Food >= cityPopulations[t] && suitable[k].owner == nil && distFromOtherCity(suitable, suitable[k].X, suitable[k].Y) > float64(6-(50-world.Height)/10) {
					suitable[k].owner = ns[n]
					cit := City{
						Name:       ns[n].Name + " city",
						owner:      ns[n],
						area:       suitable[k],
						Population: cityPopulations[t],
						IsCapital:  maxPop == cityPopulations[t],
						BasicNode: BasicNode{
							"City",
							nil,
						},
					}
					cities = append(cities, cit)
					cityType := "City"
					if cit.IsCapital {
						cityType = "Capital city"
					}
					fmt.Printf("    %v at: (%v, %v). Population: %v.\n", cityType, k.X, k.Y, cit.Population)
					world.Data[k].owner = ns[n]
					break
				}
			}
		}
		ns[n].cities = &cities
		world.Nations = append(world.Nations, ns[n])
	}
	fmt.Println(fmt.Sprintf("\n [Civilizations] = Time taken - %v", time.Since(s)))
	return world
}

func distFromOtherCity(cities map[*Coordinate]*CoordinateInfo, x, y int) float64 {
	min := 9999.0
	for i := range cities {
		city := cities[i]
		if city.owner != nil {
			dist := float64((city.X-x)*(city.X-x) + (city.Y-y)*(city.Y-y))
			if dist > 0 && dist < min {
				min = math.Sqrt(float64(dist))
			}
		}
	}
	return min
}
