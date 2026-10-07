package realtime

import (
	"context"
	"strconv"
	"time"
)

type FakeIssuer struct{}

func (f *FakeIssuer) Issue(ctx context.Context, cfg SessionConfig) (Secret, error) {
	return Secret{
		Value: "a__fake_" + strconv.FormatInt(time.Now().UnixNano(), 36),
		ExpiresAt: time.Now().Add(time.Minute),
	}, nil
}