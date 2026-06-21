package cmd

import (
	"io"
	"strings"
	"testing"
)

// isolateConfig prevents the test from picking up any real mailcow config or
// environment variables so that opts.client() fails fast with a config error
// rather than attempting a network connection.
func isolateConfig(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("MAILCOW_HOST", "")
	t.Setenv("MAILCOW_API_KEY", "")
	t.Setenv("MAILCOW_PROFILE", "")
}

func TestSyncjobCreateRejectsInvalidEnc1(t *testing.T) {
	isolateConfig(t)
	cmd := NewRootCmd(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{
		"syncjob", "create",
		"--mailbox", "a@x.org",
		"--host1", "h",
		"--user1", "u",
		"--enc1", "bogus",
		"--password1-stdin",
	})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --enc1, got nil")
	}
	if !strings.Contains(err.Error(), "enc1") {
		t.Fatalf("expected error message to mention enc1, got: %s", err.Error())
	}
}

func TestSyncjobCreateAcceptsValidEnc1(t *testing.T) {
	for _, enc := range []string{"SSL", "TLS", "PLAIN"} {
		enc := enc
		t.Run(enc, func(t *testing.T) {
			// Validation must fail before any network call: we expect a config
			// error (no host configured), not an enc1 error.
			isolateConfig(t)
			cmd := NewRootCmd(strings.NewReader("pw\n"), io.Discard, io.Discard)
			cmd.SetArgs([]string{
				"syncjob", "create",
				"--mailbox", "a@x.org",
				"--host1", "h",
				"--user1", "u",
				"--enc1", enc,
				"--password1-stdin",
			})
			err := cmd.Execute()
			if err != nil && strings.Contains(err.Error(), "enc1") {
				t.Fatalf("valid enc1 %q should not produce an enc1 error: %s", enc, err.Error())
			}
		})
	}
}
