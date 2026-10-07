package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/pranav7002/voice-server/internal/realtime"
)

type SessionService struct {
	Issuer realtime.Issuer // interface can be fake or real, chosen in hydrate()
	Config realtime.SessionConfig
}

func (s *SessionService) CreateSession(ctx context.Context) (realtime.Secret, error) {
	secret, err := s.Issuer.Issue(ctx, s.Config)
	if err != nil {
		slog.Error("issuing realtime key failed", "error", err)
		return realtime.Secret{}, errors.New("could not create session")
	}
	return secret, nil
}