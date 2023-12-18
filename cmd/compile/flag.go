package compile

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagTarget = "target"

	flagLogLevel = "log-level"
)

type flag struct {
	Target string

	// logging
	LogLevel string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Target, flagTarget, "", "Target")

	cmd.Flags().StringVar(&f.LogLevel, flagLogLevel, "info", "Log level")

}

func (f *flag) Validate() error {
	if f.Target == "" {
		return fmt.Errorf("--%s can't be empty", flagTarget)
	}
	supportedTargets := []string{"al2", "rhel9", "rockylinux9", "ubuntu"}
	if !contains(supportedTargets, f.Target) {
		return fmt.Errorf("supported targets: %v, got: %s", supportedTargets, f.Target)
	}

	_, err := logrus.ParseLevel(f.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level %#q, using default info", f.LogLevel)
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
