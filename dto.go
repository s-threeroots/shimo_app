package main

import (
	"gorm.io/gorm"
)

type Estimation struct {
	gorm.Model
	ClientName     string
	EstimationName string
	Groups         []Group
	SubTotal       int
	Tax            int
	Total          int
}

type Group struct {
	gorm.Model
	EstimationID int
	Name         string
	Items        []Item
}

type Item struct {
	gorm.Model
	GroupID   int
	Name      string
	Amount    int
	Unit      string
	UnitPrice int
}

func (e Estimation) Init() {

	e = Estimation{
		ClientName:     "",
		EstimationName: "",
		Groups: []Group{
			{
				Name: "",
				Items: []Item{
					{
						Name:      "",
						Amount:    0,
						Unit:      "",
						UnitPrice: 0,
					},
				},
			},
		},
	}
}
