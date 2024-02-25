package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
)

func (s *Server) initRouter() {
	s.router.Use(s.logger)
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		n := rand.Intn(5)
		if n == 0 {
			w.WriteHeader(400)
			return
		}
		if n == 1 {
			w.WriteHeader(500)
			return
		}
		w.Write([]byte("ping"))
	})
	secure := s.router.With(s.jwtMiddleware)
	//TODO get user by jwt
	secure.Get("/secure-ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("ping %d", r.Context().Value("id"))))
	})

	s.router.Post("/sign-up", s.SignUp)
}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	bodyReader := r.Body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	way := r.URL.Query().Get("way")
	service, ok := s.authServices[way]
	if !ok {
		err = errors.New("service doesnt exist")
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := service.SignIn(r.Context(), body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//user id to jwt and refresh

	token, err := s.jwtService.GetJwt(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"jwt": token,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(data)
}

//TODO
/*


 */

// default error message
type ErrorMessage struct {
	ErrorMessage string `json:"error_message"`
	StatusCode   int    `json:"status_code"`
}

// func thats write error
func writeError(w http.ResponseWriter, status int, errorMessage string) {
	slog.Error("error message", "err", errorMessage)
	w.WriteHeader(status)
	response := ErrorMessage{ErrorMessage: errorMessage, StatusCode: status}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal("fatal message", "err", err.Error())
	}
	w.Write(jsonResponse)
}
