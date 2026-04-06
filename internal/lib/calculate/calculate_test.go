package calculate

import (
	"testing"

	"github.com/stepan41k/grpc-test/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAvgNM(t *testing.T) {
	rates := []model.RateItem{
		{Price: 10.0},
		{Price: 20.0},
		{Price: 30.0},
		{Price: 40.0},
	}

	tests := []struct {
		name    string
		n, m    int
		want    float64
		wantErr bool
	}{
		{"Valid range", 0, 1, 15.0, false},
		{"Single element", 2, 2, 30.0, false},
		{"Full range", 0, 3, 25.0, false},
		{"N > M", 2, 1, 0, true},
		{"Out of bounds", 0, 5, 0, true},
		{"Negative index", -1, 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err :=CalculateAvgNM(rates, tt.n, tt.m)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCalculateTopN(t *testing.T) {
	rates := []model.RateItem{
		{Price: 10.0},
		{Price: 20.0},
		{Price: 30.0},
		{Price: 40.0},
	}

	tests := []struct {
		name    string
		n    int
		want    float64
		wantErr bool
	}{
		{"Valid N", 1, 20.0, false},
		{"Left limit", 0, 10.0, false},
		{"Right limit", 3, 40.0, false},
		{"N > len(rates)", 5, 0, true},
		{"Negative index", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err :=CalculateTopN(rates, tt.n)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}