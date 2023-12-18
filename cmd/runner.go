package cmd

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runner struct {
	flag   *flag
	logger *log.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return err
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	err := cmd.Help()
	if err != nil {
		return err
	}

	return nil
}
