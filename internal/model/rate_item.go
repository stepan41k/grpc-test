package model

type RateItem struct {
	Price  float64 `json:"price,string"`
	Volume float64 `json:"volume,string"`
	Amount float64 `json:"amount,string"`
}