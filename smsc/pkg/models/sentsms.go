package models

import (
	"context"
	"smsc/pkg/db"
	"time"
)

type SentSms struct {
	ID           int64
	Status       int
	Alphanumeric string
	Phone        string
	Content      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *SentSms) Insert(ctx context.Context) (*SentSms, error) {
	db := db.GetDB()

	query := `INSERT INTO "public"."sent_sms"( "alphanumeric", "phone", "content") VALUES 
	( $1, $2, $3);`

	err := db.QueryRowContext(ctx, query,
		s.Alphanumeric,
		s.Phone,
		s.Content,
	).Scan(
		&s.ID,
		&s.Status,
		&s.Alphanumeric,
		&s.Phone,
		&s.Content,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	return s, err
}

func (s *SentSms) Update(ctx context.Context) (*SentSms, error) {
	db := db.GetDB()

	query := `UPDATE "public"."sent_sms" SET
	 "status" = $2,  "updated_at" = $3 WHERE "id" = $1;`

	err := db.QueryRowContext(ctx, query,
		s.ID,
		s.Status,
		s.UpdatedAt,
	).Scan(
		&s.ID,
		&s.Status,
		&s.Alphanumeric,
		&s.Phone,
		&s.Content,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	return s, err
}
