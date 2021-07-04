package main

import (
	"testing"
)

func TestSeparateTest(t *testing.T) {

	est := Estimation{
		Groups: []Group{
			{
				Name: "name1",
				Items: []Item{
					{
						Name: "Name1",
					},
					{
						Name: "Name2",
					},
					{
						Name: "Name3",
					},
					{
						Name: "Name4",
					},
					{
						Name: "Name5",
					},
					{
						Name: "Name6",
					},
				},
			},
		},
	}

	est.SeparateGroup(0, 3)

}
