package syncutil

import (
	"testing"
)

func TestStopper(t *testing.T) {
	s := NewStopper()
	s.Close()
	s.Close()
	s.Wait()
}
