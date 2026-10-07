package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pranav7002/voice-server/internal/controller"
	"github.com/pranav7002/voice-server/internal/realtime"
	"github.com/pranav7002/voice-server/internal/services"
)

type application struct {
	config config

	sessionController *controller.SessionController
}

type config struct {
	addr      string
	openAIKey string
	session   realtime.SessionConfig
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/session", app.sessionController.CreateSessionHandler)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:              app.config.addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server has started", "addr", app.config.addr)
	return srv.ListenAndServe()
}


func (app *application) hydrate() {
	var issuer realtime.Issuer
	if app.config.openAIKey == "" {
		slog.Warn("OPENAI_API_KEY not set: using FakeIssuer (fake keys, no OpenAI calls)")
		issuer = &realtime.FakeIssuer{}
	} else {
		issuer = &realtime.OpenAIIssuer{
			APIKey: app.config.openAIKey,
			Client: &http.Client{Timeout: 10 * time.Second},
		}
	}

	sessionService := &services.SessionService{
		Issuer: issuer,
		Config: app.config.session,
	}
	app.sessionController = &controller.SessionController{SessionService: sessionService}
}