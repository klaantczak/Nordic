package hpstest

import (
	"encoding/json"
	nm "hps/networkmachine"
	"strings"
	"testing"
)

func Test_NetworkMachineToJson(t *testing.T) {
	m := nm.NewMachine("test", "default")

	jsonBin, _ := json.MarshalIndent(m, "", " ")
	jsonStr := string(jsonBin)
	expected := strings.Join([]string{
		"{",
		" \"name\": \"test\",",
		" \"machine\": \"default\",",
		" \"content\": {",
		"  \"machines\": {}",
		" },",
		" \"properties\": {}",
		"}",
	}, "\n")
	if jsonStr != expected {
		t.Errorf("Unexpected JSON serialisation")
	}
}
