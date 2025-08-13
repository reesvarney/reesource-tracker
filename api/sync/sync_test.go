package sync_test

import (
	"reesource-tracker/api/sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastEvent_NoClients(t *testing.T) {
	// Should not panic or error if no clients are registered
	assert.NotPanics(t, func() {
		sync.BroadcastEvent("test_event", map[string]string{"foo": "bar"})
	})
}

// Add more tests for EventStream and client registration as needed.
