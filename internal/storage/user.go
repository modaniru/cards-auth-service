package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/modaniru/cards-auth-service/sqlc/db"
)

type UserStorage struct {
	db      *sql.DB
	queries *db.Queries
}

func NewUserStorage(db *sql.DB, queries *db.Queries) *UserStorage {
	return &UserStorage{
		db:      db,
		queries: queries,
	}
}

func (u UserStorage) createUserWithOAuth(ctx context.Context, authType string, authId string) (int, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	q := u.queries.WithTx(tx)
	id, err := q.CreateEmptyUser(ctx)
	if err != nil {
		return 0, err
	}

	err = q.AddUserAuthType(ctx, db.AddUserAuthTypeParams{
		UserID:   sql.NullInt32{Int32: id, Valid: true},
		AuthType: sql.NullString{String: authType, Valid: true},
		AuthID:   sql.NullString{String: authId, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (u UserStorage) getUserIdByAuthTypeAndAuthId(ctx context.Context, authType string, authId string) (int, error) {
	userId, err := u.queries.GetUserByAuthTypeAndAuthId(ctx, db.GetUserByAuthTypeAndAuthIdParams{
		AuthType: sql.NullString{String: authType, Valid: true},
		AuthID:   sql.NullString{String: authId, Valid: true},
	})
	return int(userId.Int32), err
}

func (u UserStorage) GetOrCreateUserIdByAuthType(ctx context.Context, authType string, authId string) (int, error) {
	id, err := u.getUserIdByAuthTypeAndAuthId(ctx, authType, authId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = u.createUserWithOAuth(ctx, authType, authId)
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
