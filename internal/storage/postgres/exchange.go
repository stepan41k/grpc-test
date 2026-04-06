package postgres

import (
	"context"
	"fmt"
	"time"
)

func (pgs *PGStorage) SaveRate(ctx context.Context, askPrice, bidPrice float64, timestamp time.Time) (err error) {
	const path = "repository.postgres.exchange.SaveRate"
	
	tx, err := pgs.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if err != nil {
				err = fmt.Errorf("%s: %w", path, rollbackErr)
			}
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", path, commitErr)
		}
	}()
	
	// We'll take the timestamp from the API request, not the transaction start time.
	_, err = tx.Exec(ctx, `
			INSERT INTO usdt_rates(ask_price, bid_price, created_at)
			VALUES($1, $2, $3)
			RETURNING id;
		`, askPrice, bidPrice, timestamp)
	
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	
	return nil
}
