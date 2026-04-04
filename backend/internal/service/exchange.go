package service

import (
	"errors"
	"fmt"
)

type RateItem struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

// CalculateTopN returns price at specific index
func CalculateTopN(rates []RateItem, n int) (float64, error) {
	if n < 0 || n >= len(rates) {
		return 0, errors.New("index out of range")
	}
	return rates[n].Price, nil
}

// CalculateAvgNM returns average price in range [n, m]
func CalculateAvgNM(rates []RateItem, n, m int) (float64, error) {
	if n < 0 || m >= len(rates) || n > m {
		return 0, errors.New("invalid range")
	}
	var sum float64
	for i := n; i <= m; i++ {
		sum += rates[i].Price
	}
	return sum / float64(m-n+1), nil
}