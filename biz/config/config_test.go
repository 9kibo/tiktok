package config

import (
	"encoding/json"
	"testing"
)

func TestInitConfig(t *testing.T) {
	Init("../../config.ini")
	marshal, err := json.MarshalIndent(C, "  ", "    ")
	if err != nil {
		panic(err)
	}
	t.Log(string(marshal))
}
