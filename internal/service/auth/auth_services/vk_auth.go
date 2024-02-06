package authservices

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/modaniru/cards-auth-service/sqlc/db"
)

type VKAuth struct{
	vkapi *api.VK
	db *db.Queries
}

func NewVKAuth(token string, db *db.Queries) *VKAuth{
	return &VKAuth{vkapi: api.NewVK(token), db: db}
}

type vkToken struct{
	Token string `json:"token"`
}

func (v *VKAuth) SignIn(c context.Context, credentials []byte) (int, error) {
	var token vkToken
	err := json.Unmarshal(credentials, &token)
	if err != nil{
		return 0, err
	}
	//TODO check what's error when token is not correct or expire
	response, err := v.vkapi.SecureCheckToken(api.Params{"token": token.Token})
	if err != nil{
		return 0, err
	}

	userId, err := v.db.GetUserByAuthTypeAndAuthId(c, db.GetUserByAuthTypeAndAuthIdParams{
		AuthType: sql.NullString{String: v.Key(), Valid: true},
		AuthID: sql.NullString{String: strconv.Itoa(response.UserID), Valid: true},
	})
	if err != nil{
		return 0, err
	}


	slog.Info("info", "vk user id", userId)
	return 0, nil
}

// return service key
func (v *VKAuth) Key() string {
	return "vk"
}

