package controller

import (
	"net/http"

	"github.com/pranav7002/voice-server/internal/services"
)

type SessionController struct {
	SessionService *services.SessionService
}

type sessionResponse struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expiresAt"`
}

func (c *SessionController) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	secret, err := c.SessionService.CreateSession(r.Context())
	if err != nil {
		WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, sessionResponse{
		Value:     secret.Value,
		ExpiresAt: secret.ExpiresAt.Unix(),
	})
}