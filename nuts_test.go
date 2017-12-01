package nuts

import (
	"testing"
)

func TextVersion(t testing.T) {
	version = `v0.0.0`
	if Version != version {
		t.Errorf("expected %s, got %s", version, Version)
	}
}
