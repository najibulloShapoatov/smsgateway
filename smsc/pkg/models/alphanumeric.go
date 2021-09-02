package models

import (
	"context"
	"smsc/pkg/db"

	"time"
)

type Alphanumeric struct {
	ID             int64
	Name           string
	OrganizationID int64
	CustomerID     int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (a *Alphanumeric) Get(ctx context.Context) (*Alphanumeric, error) {
	db := db.GetDB()
	query := `select * from alphanumerics where customer_id = $1 and name = $2`
	err := db.QueryRowContext(ctx, query, a.CustomerID, a.Name).Scan(
		&a.ID,
		&a.Name,
		&a.OrganizationID,
		&a.CustomerID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	return a, err
}
