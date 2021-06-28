package main

import (
	"gorm.io/gorm"
)

type Estimation struct {
	gorm.Model
	ID             uint    `json:"id"`
	ClientName     string  `json:"client_name"`
	EstimationName string  `json:"estimation_name"`
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
						Price:     0,
					},
				},
			},
		},
	}
}

type RedirectResponse struct {
	URL string `json:"url"`
}
