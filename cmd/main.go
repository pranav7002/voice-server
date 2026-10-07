package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/pranav7002/voice-server/internal/realtime"
)

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	cfg := config{
		addr:      ":" + getEnv("PORT", "8080"),
		openAIKey: os.Getenv("OPENAI_API_KEY"),
		session: realtime.SessionConfig{
			Model:        getEnv("REALTIME_MODEL", "gpt-realtime-1.5"),
			Voice:        getEnv("REALTIME_VOICE", "marin"),
			Instructions: "You are a friendly voice assistant. Keep answers short and conversational.",
		},
	}

	app := &application{config: cfg}

	// Inject all dependencies
	app.hydrate()

	if err := app.run(app.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}