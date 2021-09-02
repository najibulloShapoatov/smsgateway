package models

import (
	"context"
	"smsc/pkg/db"
)

type Customer struct {
	ID     int64
	APIKey string
}

func (c *Customer) Get(ctx context.Context) (*Customer, error) {

	db := db.GetDB()

	query := `select * from customers where api_key = $1`

	err := db.QueryRowContext(ctx, query, c.APIKey).Scan(
		&c.ID,
		&c.APIKey,
	)

	return c, err

}
