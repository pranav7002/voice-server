package realtime

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)


const clientSecretsURL = "https://api.openai.com/v1/realtime/client_secrets"

type OpenAIIssuer struct {
	APIKey string 
	Client *http.Client
}

type clientSecretRequest struct {
	Session sessionBody `json:"session"`
}

type sessionBody struct {
	Type         string    `json:"type"`
	Model        string    `json:"model"`
	Instructions string    `json:"instructions,omitempty"`
	Audio        audioBody `json:"audio"`
}

type audioBody struct {
	Output outputBody `json:"output"`
}

type outputBody struct {
	Voice string `json:"voice"`
}

type clientSecretResponse struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expires_at"`
}

func (o *OpenAIIssuer) Issue(ctx context.Context, cfg SessionConfig) (Secret, error) {
	body, err := json.Marshal(clientSecretRequest{
		Session: sessionBody{
			Type:         "realtime",
			Model:        cfg.Model,
			Instructions: cfg.Instructions,
			Audio:        audioBody{Output: outputBody{Voice: cfg.Voice}},
		},
	})
	if err != nil {
		return Secret{}, fmt.Errorf("encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, clientSecretsURL, bytes.NewReader(body))
	if err != nil {
		return Secret{}, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return Secret{}, fmt.Errorf("calling openai: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return Secret{}, fmt.Errorf("openai returned %d: %s", resp.StatusCode, detail)
	}

	var out clientSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Secret{}, fmt.Errorf("decoding response: %w", err)
	}
	if out.Value == "" {
		return Secret{}, errors.New("openai response had no key")
	}

	return Secret{Value: out.Value, ExpiresAt: time.Unix(out.ExpiresAt, 0)}, nil
}