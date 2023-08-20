package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseDevProfile(t *testing.T) {
	UseProfile("dev")
	Init("../../config.ini")
	assert.Equal(t, "127.0.0.1:8080", C.Server.Addr)
}
