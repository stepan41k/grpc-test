package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGrinexClient_FetchTicker(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который отдает ваш реальный JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// JSON со строками, как на вашем скриншоте
		_, _ = w.Write([]byte(`{
			"timestamp": 1775383658,
			"asks": [{"price": "80.9", "volume": "100.0", "amount": "8090.0"}],
			"bids": [{"price": "80.8", "volume": "50.0", "amount": "4040.0"}]
		}`))
	}))
	defer server.Close()

	// Инициализируем клиент, указывая адрес тестового сервера
	client := NewGrinexClient(server.URL)

	// Выполняем запрос
	res, err := client.FetchRates(context.Background())

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, int64(1775383658), res.Timestamp)
	assert.Equal(t, 80.9, res.Asks[0].Price) // Проверка, что "80.9" стало 80.9
	assert.Equal(t, 80.8, res.Bids[0].Price)
}