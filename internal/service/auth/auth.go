package auth

import (
	"context"
	"encoding/json"
)

type Auth interface {
	SignIn(c context.Context, credentials []byte) (int, error)
	//return service key
	Key() string
}

type AuthStub struct {
}

type Request struct {
	Login string `json:"login"`
}

var users map[int]string = make(map[int]string)
var logins map[string]int = make(map[string]int)
var id int = 1

func (a *AuthStub) SignIn(c context.Context, credentials []byte) (int, error) {
	req := Request{}
	err := json.Unmarshal(credentials, &req)
	if err != nil {
		return 0, err
	}

	if id, ok := logins[req.Login]; ok {
		return id, nil
	}

	users[id] = req.Login
	logins[req.Login] = id
	id++

	return logins[req.Login], nil
}

// return service key
func (a *AuthStub) Key() string {
	return "stub"
}
