package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func changedBool(cmd *cobra.Command, onFlag, offFlag string) (*bool, error) {
	on := cmd.Flags().Changed(onFlag)
	off := cmd.Flags().Changed(offFlag)
	if on && off {
		return nil, fmt.Errorf("--%s and --%s cannot be used together", onFlag, offFlag)
	}
	if on {
		v := true
		return &v, nil
	}
	if off {
		v := false
		return &v, nil
	}
	return nil, nil
}

func confirm(in io.Reader, errOut io.Writer, yes bool, subject string) error {
	if yes {
		return nil
	}
	fmt.Fprintf(errOut, "Delete %s? Type yes to continue: ", subject)
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if strings.TrimSpace(line) != "yes" {
		return errors.New("delete cancelled; pass --yes to confirm in scripts")
	}
	return nil
}

func passwordFromFlags(cmd *cobra.Command, in io.Reader, passwordEnv string, passwordStdin bool) (string, error) {
	if passwordEnv != "" && passwordStdin {
		return "", errors.New("--password-env and --password-stdin cannot be used together")
	}
	switch {
	case passwordEnv != "":
		value := strings.TrimRight(getenv(passwordEnv), "\r\n")
		if value == "" {
			return "", fmt.Errorf("environment variable %s is empty", passwordEnv)
		}
		return value, nil
	case passwordStdin:
		data, err := io.ReadAll(in)
		if err != nil {
			return "", err
		}
		value := strings.TrimRight(string(data), "\r\n")
		if value == "" {
			return "", errors.New("password from stdin is empty")
		}
		return value, nil
	default:
		return "", errors.New("password is required; use --password-stdin or --password-env")
	}
}

var getenv = os.Getenv

func splitMailbox(address string) (string, string, error) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("mailbox address must be local@domain")
	}
	return parts[0], parts[1], nil
}

func intFlagAttr(cmd *cobra.Command, name string, attrName string, attrs map[string]any) {
	if cmd.Flags().Changed(name) {
		value, _ := cmd.Flags().GetInt(name)
		attrs[attrName] = strconv.Itoa(value)
	}
}

func stringFlagAttr(cmd *cobra.Command, name string, attrName string, attrs map[string]any) {
	if cmd.Flags().Changed(name) {
		value, _ := cmd.Flags().GetString(name)
		attrs[attrName] = value
	}
}
