package main

import (
	"gorm.io/gorm"
)

type Estimation struct {
	gorm.Model
	ID             uint    `json:"id"`
	ClientName     string  `json:"client_name"`
	EstimationName string  `json:"estimation_name"`
	TemplateID     uint    `json:"template_id"`
	Groups         []Group `json:"groups"`
	SubTotal       int     `json:"sub_total"`
	Tax            int     `json:"tax"`
	Total          int     `json:"total"`
}

type Group struct {
	gorm.Model
	ID           uint   `json:"id"`
	EstimationID uint   `json:"estimation_id"`
	Name         string `json:"name"`
	Items        []Item `json:"items"`
	Order        int    `json:"order"`
}

type Item struct {
	gorm.Model
	ID        uint   `json:"id"`
	GroupID   uint   `json:"group_id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Unit      string `json:"unit"`
	UnitPrice int    `json:"unit_price"`
	Price     int    `json:"price"`
	Order     int    `json:"order"`
}

func (e Estimation) Init() {

	g := Group{}
	g.Init()

	e = Estimation{
		ClientName:     "",
		EstimationName: "",
		Groups: []Group{
			g,
		},
	}
}

func (g Group) Init() {

	i := Item{}
	i.Init()

	g = Group{
		Name: "",
		Items: []Item{
			i,
		},
	}
}

func (e *Estimation) SeparateGroup(gidx int, iidx int) {

	nName := e.Groups[gidx].Name
	items := e.Groups[gidx].Items[iidx:]
	nItems := make([]Item, len(items))
	copy(nItems, items)

	ng := Group{
		Name:  nName,
		Items: nItems,
	}
	e.Groups[gidx].Items = e.Groups[gidx].Items[:iidx]
	nGroups := []Group{}
	befG := e.Groups[:gidx+1]
	afG := e.Groups[gidx+1:]

	nGroups = append(nGroups, befG...)
	nGroups = append(nGroups, ng)
	nGroups = append(nGroups, afG...)

	e.Groups = nGroups
}

func (i Item) Init() {
	i = Item{
		Name:      "",
		Amount:    0,
		Unit:      "",
		UnitPrice: 0,
		Price:     0,
	}
}

type RedirectResponse struct {
	URL string `json:"url"`
}
