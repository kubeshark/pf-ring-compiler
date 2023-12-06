package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kubeshark/pfring-compiler/cmd"
)

func main() {
	err := mainE(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func mainE(ctx context.Context) error {
	var err error

	// init logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	var rootCommand *cobra.Command
	{
		c := cmd.Config{
			Logger: logger,
		}

		rootCommand, err = cmd.New(c)
		if err != nil {
			return err
		}
	}

	err = rootCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}
