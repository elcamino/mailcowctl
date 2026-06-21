package cmd

import (
	"io"
	"strings"
	"testing"
)

func TestBccCreateRejectsInvalidType(t *testing.T) {
	isolateConfig(t)
	cmd := NewRootCmd(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{
		"--host", "http://127.0.0.1:1", "--api-key", "k",
		"bcc", "create",
		"--local-dest", "example.org",
		"--bcc-dest", "a@example.org",
		"--type", "bogus",
	})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --type, got nil")
	}
	if !strings.Contains(err.Error(), "type") {
		t.Fatalf("expected error message to mention type, got: %s", err.Error())
	}
}
