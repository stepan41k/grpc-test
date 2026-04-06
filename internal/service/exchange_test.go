package service

import (
	"context"
	"testing"
	"time"

	"github.com/stepan41k/grpc-test/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveRate(ctx context.Context, ask, bid float64, ts time.Time) error {
	args := m.Called(ctx, ask, bid, ts)
	return args.Error(0)
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) FetchRates(ctx context.Context) (*model.GrinexResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GrinexResponse), args.Error(1)
}

func TestGetAndProcessRates_Success(t *testing.T) {
	repo := new(MockRepo)
	api := new(MockClient)
	log := zap.NewNop()

	fakeResponse := &model.GrinexResponse{
		Timestamp: 1712345678,
		Asks: []model.RateItem{
			{Price: 80.0, Volume: 1.0},
			{Price: 81.0, Volume: 2.0},
			{Price: 82.0, Volume: 3.0},
		},
		Bids: []model.RateItem{
			{Price: 79.0, Volume: 1.0},
		},
	}

	api.On("FetchRates", mock.Anything).Return(fakeResponse, nil)
	repo.On("SaveRate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := New(log, api, repo)

	// topN=1, n=0, m=2
	res, err := svc.GetAndProcessRates(context.Background(), 1, 0, 2)

	assert.NoError(t, err)
	assert.Equal(t, 81.0, res.TopPrice)
	assert.Equal(t, 81.0, res.AvgPrice)
	repo.AssertExpectations(t)
}
