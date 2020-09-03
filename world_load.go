package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"image/color"

	"fyne.io/fyne"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var LoadedMap *WorldMap
var LocalBroker Broker
var InitialConcepts = []Concept{
	{BasicNode{"Concept", nil}, "Concept"},
	{BasicNode{"Concept", nil}, "City"},
	{BasicNode{"Concept", nil}, "Nation"},
	{BasicNode{"Concept", nil}, "Location"},
}

type CollectionContainer struct {
	Nodes               driver.Collection
	Graph               driver.Graph
	FromEdgesCollection driver.Graph
	db                  driver.Database
}

type Concept struct {
	BasicNode
	Name string
}

type World struct {
	BasicNode
	Name              string
	Width             int
	Height            int
	Smoothness        float64
	MinimumIslandSize int
	HeightThreshold   float64
	InitCivRadius     float64
	NumCivs           int
	NumIterations     int
}

func (n *World) SaveData() {
	n.BasicNode.SaveData(n)
}

var container = CollectionContainer{}

func p(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		// Handle error
	}
}

type PoliticalInfo struct {
	owner  *Nation
	solved bool
}

type Nation struct {
	BasicNode
	Name       string
	Population int
	cities     *[]City
	Index      int
}

type City struct {
	Name       string
	owner      *Nation
	area       *CoordinateInfo
	IsCapital  bool
	Population int
	BasicNode
}

func (n *Nation) SaveData() {
	n.BasicNode.SaveData(n)
	n.ConnectToConcept("Nation")
}

func (n *City) SaveData() {
	n.BasicNode.SaveData(n)
	n.ConnectToConcept("City")
}

func (n *Concept) SaveData() {
	n.BasicNode.SaveData(n)
	n.ConnectToConcept("Concept")
}

func (n *CoordinateInfo) SaveData() {
	n.BasicNode.SaveData(n)
}

func (n *CoordinateInfo) GetOwner() *Nation {
	s := time.Now()
	if n.solved || n.documentMeta == nil {
		return n.owner
	} else {
		cursor, err := container.db.Query(nil, SOLVE_FIELDS, map[string]interface{}{"id": n.documentMeta.ID})
		p(err)
		defer cursor.Close()
		count := 0
		for {
			count++
			var doc struct {
				Nation
				Relation string
			}
			meta, err := cursor.ReadDocument(nil, &doc)
			if driver.IsNoMoreDocuments(err) {
				n.solved = true
				break
			} else if err != nil {
				p(err)
			}
			if doc.Relation == "OWNS" {
				n.solved = true
				n.owner = &doc.Nation
				n.owner.documentMeta = &meta
			}
			fmt.Println(fmt.Sprintf("\n [Access owner - Using query] = Time taken - %v", time.Since(s)))
			return n.owner
		}
	}
	return nil
}

type Resources struct {
	Food    int
	Wealth  int
	Product int

	RawProduct int

	// Instaled Productions
	FProduction int
	PProduction int
}

func (r *Resources) Add(n Resources) *Resources {
	r.FProduction += n.FProduction
	r.Food += n.Food
	r.Wealth += n.Wealth
	r.Product += n.Product
	r.RawProduct += n.RawProduct
	r.PProduction += n.PProduction

	return r
}

var params = map[string]float64{
	"WORLD_WIDTH":       50,
	"WORLD_HEIGHT":      50,
	"SMOOTHNESS":        2.1,
	"TOO_SMALL_ISLANDS": 2,
	"HEIGHT_THRESHOLD":  0.2,
	"INIT_CIV_RADIUS":   2.0,
	"NUM_CIV":           5,
}

func CharacterListener(e string, k interface{}) {
	key := k.(*fyne.KeyEvent)
	if LoadedMap == nil {
		return
	}
	switch key.Name {
	case fyne.KeyDown:
		LoadedMap.Character.Y = (LoadedMap.Character.Y + 1) % LoadedMap.Height
	case fyne.KeyUp:
		LoadedMap.Character.Y = (LoadedMap.Height + LoadedMap.Character.Y - 1) % LoadedMap.Height
	case fyne.KeyLeft:
		LoadedMap.Character.X = (LoadedMap.Width + LoadedMap.Character.X - 1) % LoadedMap.Width
	case fyne.KeyRight:
		LoadedMap.Character.X = (LoadedMap.Character.X + 1) % LoadedMap.Width
	}
	amd, _ := json.MarshalIndent(LoadedMap.GetCoordinateInfo(LoadedMap.Character.X, LoadedMap.Character.Y), "", "\t")
	lines := strings.Split(string(amd), "\n\t")
	LocalBroker.EmitEvent("ShowInfo", lines) // Make it less ugly
}

func InitWorldLoader(log func(string), broker Broker) (func(string, bool), func(x, y, w, h int) color.Color, func(string, string), func(string) float64) {
	LocalBroker = broker
	LocalBroker.SubscribeToEvent("KeyDownEvent", CharacterListener)
	return func(worldName string, save bool) {
			if len(worldName) < 4 {
				log(fmt.Sprintf("World name [%v] is too small, loading fail...", worldName))
				return
			}
			// Get the world
			WORLD_NAME := worldName
			world := World{
				BasicNode{"World", nil},
				worldName,
				int(params["WORLD_WIDTH"]),
				int(params["WORLD_HEIGHT"]),
				params["SMOOTHNESS"],
				int(params["TOO_SMALL_ISLANDS"]),
				params["HEIGHT_THRESHOLD"],
				params["INIT_CIV_RADIUS"],
				int(params["NUM_CIV"]),
				5,
			}
			dbSetup(WORLD_NAME) // TODO: Add mode to generate and save
			wMeta := GetNodeByName(worldName, nil)
			worldFound := wMeta != nil
			var coordinates []CoordinateInfo
			var nations []*Nation
			if LoadedMap != nil && worldName == LoadedMap.Name && save {
				world = LoadedMap.World
				LoadedMap.SaveData()
				wMeta = LoadedMap.documentMeta
				world.documentMeta = wMeta
				for i := range LoadedMap.rawData {
					coordinates = append(coordinates, *LoadedMap.rawData[i])
				}
				UpdateCoordinates(LoadedMap.rawData, world.documentMeta)
				SaveNationsAndCities(LoadedMap)
				log(fmt.Sprintf("[%v] saved!", WORLD_NAME))
			} else if !worldFound || !save {
				log(fmt.Sprintf("[%v] will be created!", WORLD_NAME))
				if save {
					world.SaveData()
				}
				log("Generating coordinates...")
				c := GenerateWorldCoordinates(int(params["WORLD_WIDTH"]), int(params["WORLD_HEIGHT"]))
				log("Generating coordinates...Done!")
				coordinates = SaveCoordinates(c, world.BasicNode, save)
				log("Generating z levels...")
				worldMap := GenerateWorldBasedOnCoordinates(coordinates, world)
				log("Generating z levels...Done!")
				log("Generating resources...")
				worldMap = GenerateResources(worldMap, world.HeightThreshold)
				if save {
					UpdateCoordinates(worldMap.rawData, world.documentMeta)
				}
				log("Generating resources... Done")
				log("Generating nations...")
				worldMap = GenerateInitialNations(worldMap)
				log("Generating nations...Done!")
				if save {
					SaveNationsAndCities(worldMap)
				}
			} else {
				msg := fmt.Sprintf("[%v] Found!\n", WORLD_NAME)
				log(msg)
				fmt.Println(msg)
				GetNodeByName(WORLD_NAME, &world)
				coordinates = []CoordinateInfo{}
				nations := []*Nation{}
				s := time.Now()
				cursor, err := container.db.Query(nil, GET_COORDINATES, map[string]interface{}{"name": world.Name})
				p(err)
				defer cursor.Close()
				count := 0
				for {
					count++
					var doc struct {
						Vertex CoordinateInfo
						Owner  []Nation
					}
					meta, err := cursor.ReadDocument(nil, &doc)
					if driver.IsNoMoreDocuments(err) {
						break
					} else if err != nil {
						p(err)
					}
					doc.Vertex.documentMeta = &meta
					if len(doc.Owner) > 0 {
						doc.Vertex.owner = &doc.Owner[0]
						nations = append(nations, doc.Vertex.owner)
					}
					coordinates = append(coordinates, doc.Vertex)
				}
				if len(coordinates) == 0 {
					p(fmt.Errorf("NO COORDINATES FOUND"))
				}
				fmt.Println(fmt.Sprintf("\n [Load graph] = Time taken - %v", time.Since(s)))
			}
			wMap := FromCoordinateInfoToWorldMap(coordinates, world)
			wMap.Nations = nations
			PrintCoordinates(wMap, log)
			LoadedMap = wMap
			log("Finished loading " + worldName)
		}, func(x, y, w, h int) color.Color {
			if LoadedMap == nil {
				return checkerPatternFunc(x, y, w, h)
			}
			return worldPixel(x, y, w, h, LoadedMap)
		}, func(parmName string, value string) {
			f, _ := strconv.ParseFloat(value, 64)
			params[parmName] = f
		}, func(parmName string) float64 {
			return params[parmName]
		}
}

func SaveCoordinates(coordinates []Coordinate, wMeta BasicNode, save bool) []CoordinateInfo {
	finalCoordinates := []CoordinateInfo{}
	coordinatesToWorldEdges := ECollection("FROM")
	for c := range coordinates {
		cc := CoordinateInfo{
			Coordinate: coordinates[c],
			GeograficInfo: GeograficInfo{
				Z:         0,
				Biome:     0,
				Resources: &Resources{},
			},
		}
		if save {
			cMeta, err := container.Nodes.CreateDocument(nil, cc)
			cc.documentMeta = &cMeta
			p(err)
			_, err = coordinatesToWorldEdges.CreateDocument(nil, driver.EdgeDocument{
				From: cMeta.ID,
				To:   wMeta.documentMeta.ID,
			})
			p(err)
		}
		finalCoordinates = append(finalCoordinates, cc)
	}
	return finalCoordinates
}

func UpdateCoordinates(coordinates []*CoordinateInfo, wMeta *driver.DocumentMeta) []CoordinateInfo {
	a := []CoordinateInfo{}
	for i := range coordinates {
		coordinates[i].SaveData()
		a = append(a, *coordinates[i])
	}
	return a
}
func SaveNationsAndCities(w *WorldMap) {

	for n := range w.Nations {
		w.Nations[n].SaveData()
		if w.Nations[n].cities != nil {
			savedCities := SaveCities(w.Nations[n], *w.Nations[n].cities)
			w.Nations[n].cities = &savedCities
		}
	}
}
func SaveCities(nation *Nation, cities []City) []City {
	cityLocationEdge := ECollection("IS_AT")
	ownsCollection := ECollection("OWNS")
	for k := range cities {
		cities[k].SaveData()
		cities[k].ConnectToNode(cityLocationEdge, cities[k].area.documentMeta.ID)
		nation.ConnectToNode(ownsCollection, cities[k].documentMeta.ID)
	}
	return cities
}

//
func dbSetup(worldName string) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	p(err)
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.JWTAuthentication("root", ""),
	})
	p(err)
	fmt.Println(c)
	var db driver.Database
	// Open "map" database
	does, err := c.DatabaseExists(nil, "worlds")
	p(err)
	if !does {
		db, err = c.CreateDatabase(nil, "worlds", &driver.CreateDatabaseOptions{})
		p(err)
	} else {
		db, err = c.Database(nil, "worlds")
		p(err)
	}

	container.db = db
	gExists, err := container.db.GraphExists(nil, "worldGraphs")
	p(err)

	nodes := createCollectionIfNeeded(container.db, "nodes", false)
	container.Nodes = nodes
	createCollectionIfNeeded(container.db, "IS_A", true)
	createCollectionIfNeeded(container.db, "IS_AT", true)
	createCollectionIfNeeded(container.db, "OWNS", true)
	createCollectionIfNeeded(container.db, "FROM", true)
	if !gExists {
		g, err := container.db.CreateGraph(nil, "worldGraphs", &driver.CreateGraphOptions{
			EdgeDefinitions: []driver.EdgeDefinition{
				{
					Collection: "IS_A",
					From:       []string{"nodes"},
					To:         []string{"nodes"},
				},
				{
					Collection: "IS_AT",
					From:       []string{"nodes"},
					To:         []string{"nodes"},
				},
				{
					Collection: "OWNS",
					From:       []string{"nodes"},
					To:         []string{"nodes"},
				},
				{
					Collection: "FROM",
					From:       []string{"nodes"},
					To:         []string{"nodes"},
				},
			},
		})
		p(err)
		container.Graph = g
	} else {
		g, err := container.db.Graph(nil, "worldGraphs")
		p(err)
		container.Graph = g

	}

	// createCollectionIfNeeded(container.MapGraph)
	initializeConceptsIfNeeded(nodes)
}

func createCollectionIfNeeded(db driver.Database, name string, edge bool) driver.Collection {
	var col driver.Collection
	docType := driver.CollectionTypeDocument
	if edge {
		docType = driver.CollectionTypeEdge
	}
	// Open collection
	cdoes, err := db.CollectionExists(nil, name)
	p(err)
	if !cdoes {
		col, err = db.CreateCollection(nil, name, &driver.CreateCollectionOptions{Type: docType})
		p(err)
	} else {
		col, err = db.Collection(nil, name)
		p(err)
	}
	return col
}

func initializeConceptsIfNeeded(nd driver.Collection) {
	for i := range InitialConcepts {
		var meta driver.DocumentMeta
		conc := Concept{}
		cursor, err := nd.Database().Query(nil, FIND_NODE_BY_NAME, map[string]interface{}{"name": InitialConcepts[i].Name})
		p(err)
		defer cursor.Close()
		if !cursor.HasMore() {
			InitialConcepts[i].SaveData()
			meta = *InitialConcepts[i].documentMeta
		} else {
			meta, err = cursor.ReadDocument(nil, &conc)
			p(err)
		}
		InitialConcepts[i].documentMeta = &meta
	}
}

func ECollection(name string) driver.Collection {
	e, err := container.db.Collection(nil, name)
	p(err)
	return e
}

func ConnectToConcept(conceptName string, id driver.DocumentID) {
	toBeEdge := ECollection("IS_A")
	var cId driver.DocumentID
	for i := range InitialConcepts {
		if InitialConcepts[i].Name == conceptName {
			cId = InitialConcepts[i].BasicNode.documentMeta.ID
			break
		}
	}
	if cId == "" {
		p(fmt.Errorf("Concept not found"))
	}
	toBeEdge.CreateDocument(nil, driver.EdgeDocument{
		From: id,
		To:   cId,
	})
}

func GetNodeByName(name string, result interface{}) *driver.DocumentMeta {
	r := result
	if r == nil {
		r = &map[string]interface{}{}
	}
	cursor, err := container.db.Query(nil, FIND_NODE_BY_NAME, map[string]interface{}{"name": name})
	if driver.IsNotFound(err) {
		return nil
	}
	defer cursor.Close()
	if cursor.HasMore() {
		meta, err := cursor.ReadDocument(nil, r)
		p(err)
		return &meta
	}
	return nil
}

//TODO: Time passing
//TODO: Make concepts connect to each other
// ----
//TODO: Time Passing
