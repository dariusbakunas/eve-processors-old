package esi

import (
	"context"
	"testing"
)

func TestEsi(t *testing.T) {
	ctx := context.Background()
	m := PubSubMessage{Data: []byte{'g', 'o', 'l', 'a', 'n', 'g'}}

	if got := Esi(ctx, m); got != nil {
		t.Errorf("Esi() = %q, want nil", got)
	}
}