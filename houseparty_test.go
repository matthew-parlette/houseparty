package houseparty

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	_ = GetEnv("HOME", "")
}
