package config

import (
	"encoding/json"
	"os"
	"testing"
	"tiktok/pkg/constant"
)

func TestInitConfig(t *testing.T) {
	os.Setenv(constant.ProfileEnv, "dev")
	Init("../../config.ini")
	marshal, err := json.MarshalIndent(C, "  ", "    ")
	if err != nil {
		panic(err)
	}
	t.Log(string(marshal))
}
