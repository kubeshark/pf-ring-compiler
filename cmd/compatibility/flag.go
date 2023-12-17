package compatibility

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagLogLevel = "log-level"
)

type flag struct {
	// logging
	LogLevel string
}

func (f *flag) Init(cmd *cobra.Command) {

	cmd.Flags().StringVar(&f.LogLevel, flagLogLevel, "info", "Log level")

}

func (f *flag) Validate() error {

	_, err := logrus.ParseLevel(f.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level %#q, using default info", f.LogLevel)
	}

	return nil
}
