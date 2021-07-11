package models

import (
	"time"

	gorm "gorm.io/gorm"
)


type Day struct {
	gorm.Model
	CreationDay time.Time `json:"creationday"` //redundant with auto field created_at but let's populate for past days
	RealPrice int `json:"realprice"`
	PredictedPrice int `json:"predictedprice"`
	Sentiment int `json:"sentiment"`
}