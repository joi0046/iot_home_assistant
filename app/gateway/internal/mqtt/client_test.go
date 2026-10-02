package mqtt

import (
	"testing"
	"time"
)

func TestNewFailsWhenBrokerUnavailable(t *testing.T) {
	done := make(chan error, 1)

	go func() {
		_, err := New("tcp://127.0.0.1:1")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("New() error = nil, want error")
		}

	case <-time.After(8 * time.Second):
		t.Fatal("New() did not return; connection retry is blocking")
	}
}
