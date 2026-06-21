package cmd

import (
	"io"
	"strings"
	"testing"
)

func TestTransportCreatePasswordSourceMutualExclusion(t *testing.T) {
	isolateConfig(t)
	cmd := NewRootCmd(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{
		"--host", "http://127.0.0.1:1", "--api-key", "k",
		"transport", "create",
		"--destination", "example.org",
		"--nexthop", "[relay.example.org]:587",
		"--password-env", "SOMEVAR",
		"--password-stdin",
	})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when both --password-env and --password-stdin are set, got nil")
	}
	if !strings.Contains(err.Error(), "--password-env") || !strings.Contains(err.Error(), "--password-stdin") {
		t.Fatalf("expected error to mention both --password-env and --password-stdin, got: %s", err.Error())
	}
}
