package authservices

import (
	"context"
	"encoding/json"
	"github.com/modaniru/cards-auth-service/internal/storage"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
)

type VKAuth struct {
	vkapi       *api.VK
	userStorage storage.IUser
}

// NewVKAuth return Auth implementation
func NewVKAuth(token string, userStorage storage.IUser) *VKAuth {
	return &VKAuth{
		vkapi:       api.NewVK(token),
		userStorage: userStorage,
	}
}

// credentials json type must look as vkToken
type vkToken struct {
	Token string `json:"token"`
}

// SignIn is vk oauth implementation of Auth
// credentials must be look as
//
//	{
//		"token": "vk_oauth_token"
//	}
//
// IMPORTANT: user oauth token must be generated using app token. app token in client side must be equal with server side
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

	id, err := v.userStorage.GetOrCreateUserIdByAuthType(c, v.Key(), strconv.Itoa(response.UserID))
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Key return service key
func (v *VKAuth) Key() string {
	return "vk"
}
