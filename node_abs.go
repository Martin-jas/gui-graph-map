package main

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
)

type NodeData interface {
	SaveData()
	ConnectToConcept(cName string)
	ConnectToNode(id driver.DocumentID)
	CreateComplete()
	SolveFields()
}

type BasicNode struct {
	NodeType     string
	documentMeta *driver.DocumentMeta
}

func (b *BasicNode) SaveData(data interface{}) {
	if b.documentMeta != nil {
		_, err := container.Nodes.UpdateDocument(nil, b.documentMeta.Key, data)
		p(err)
	} else {
		nMeta, err := container.Nodes.CreateDocument(nil, data)
		p(err)
		b.documentMeta = &nMeta
	}
}

func (b *BasicNode) ConnectToConcept(cName string) {
	if b.documentMeta == nil {
		p(fmt.Errorf("Tried to connect to concept [%v] without saving self [%v].", cName, b.NodeType))
	}
	toBeEdge := ECollection("IS_A")
	var cId driver.DocumentID
	for i := range InitialConcepts {
		if InitialConcepts[i].Name == cName {
			cId = InitialConcepts[i].BasicNode.documentMeta.ID
			break
		}
	}
	if cId == "" {
		c := Concept{BasicNode{"Concept", nil}, cName}
		c.SaveData()
		cId = c.documentMeta.ID
	}
	_, err := toBeEdge.CreateDocument(nil, driver.EdgeDocument{
		From: b.documentMeta.ID,
		To:   cId,
	})
	p(err)
}

func (b *BasicNode) ConnectToNode(edgeCollection driver.Collection, id driver.DocumentID) {
	_, err := edgeCollection.CreateDocument(nil, driver.EdgeDocument{
		From: b.documentMeta.ID,
		To:   id,
	})
	p(err)
}

func (b *BasicNode) SolveFields() {
	p(fmt.Errorf("tried to solve with basic node!-%v", b.NodeType))
}
