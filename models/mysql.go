package models

import "github.com/lee0720/nuwa/pkg/essentials"

type MarketData struct {
	essentials.BasicModel

	EntityID   string  `gorm:"VARCHAR(100);index;not null"`
	EntityName string  `gorm:"VARCHAR(100);not null"`
	PE         float64 `gorm:"FLOAT(4,4)"`
	PB         float64 `gorm:"FLOAT(4,4)"`
	PS         float64 `gorm:"FLOAT(4,4)"`
	DateOn     string  `gorm:"VARCHAR(8);not null"`
}
