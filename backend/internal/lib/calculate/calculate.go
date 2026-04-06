package calculate

import (
	"errors"
	"fmt"
	"math"

	"github.com/stepan41k/grpc-test/internal/model"
)

func CalculateTopN(rates []model.RateItem, n int) (float64, error) {
	if n < 0 || n >= len(rates) {
		return 0, errors.New("index out of range")
	}
	return rates[n].Price, nil
}

func CalculateAvgNM(rates []model.RateItem, n, m int) (float64, error) {
	limit := len(rates)
	
	if n < 0 || m < 0 ||m >= limit || n >= limit || n > m {
		return 0, fmt.Errorf("indices out of range")
	}
	var sum float64
	for i := n; i <= m; i++ {
		sum += rates[i].Price
	}
	
	return round(sum / float64(m-n+1)), nil
}

func round(expression float64) float64 {
	ratio := math.Pow(10, float64(2))
	return math.Round(expression * ratio) / ratio
}