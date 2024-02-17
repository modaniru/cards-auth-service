package storage

import (
	"context"
	"database/sql"
	"github.com/modaniru/cards-auth-service/sqlc/db"
)

type Storage struct {
	IUser
}

type IUser interface {
	GetOrCreateUserIdByAuthType(ctx context.Context, authType string, authId string) (int, error)
}

func NewStorage(db *sql.DB, queries *db.Queries) *Storage {
	return &Storage{
		IUser: NewUserStorage(db, queries),
	}
}
