package authservices

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/modaniru/cards-auth-service/sqlc/db"
)

type VKAuth struct {
	vkapi   *api.VK
	queries *db.Queries
	db      *sql.DB
}

func NewVKAuth(token string, queries *db.Queries, db *sql.DB) *VKAuth {
	return &VKAuth{vkapi: api.NewVK(token), queries: queries, db: db}
}

type vkToken struct {
	Token string `json:"token"`
}

func (v *VKAuth) SignIn(c context.Context, credentials []byte) (int, error) {
	var token vkToken
	err := json.Unmarshal(credentials, &token)
	if err != nil {
		return 0, err
	}
	//TODO check what's error when token is not correct or expire
	response, err := v.vkapi.SecureCheckToken(api.Params{"token": token.Token})
	if err != nil {
		return 0, err
	}

	userId, err := v.queries.GetUserByAuthTypeAndAuthId(c, db.GetUserByAuthTypeAndAuthIdParams{
		AuthType: sql.NullString{String: v.Key(), Valid: true},
		AuthID:   sql.NullString{String: strconv.Itoa(response.UserID), Valid: true},
	})

	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("user not found, register.", "vk_id", response.UserID)
		tx, err := v.db.Begin()
		if err != nil {
			return 0, err
		}
		defer tx.Rollback()

		q := v.queries.WithTx(tx)
		id, err := q.CreateEmptyUser(c)
		if err != nil {
			return 0, err
		}

		err = q.AddUserAuthType(c, db.AddUserAuthTypeParams{
			UserID:   sql.NullInt32{Int32: id, Valid: true},
			AuthType: sql.NullString{String: v.Key(), Valid: true},
			AuthID:   sql.NullString{String: strconv.Itoa(response.UserID), Valid: true},
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

	if err != nil {
		return 0, err
	}

	return int(userId.Int32), nil
}

// return service key
func (v *VKAuth) Key() string {
	return "vk"
}
