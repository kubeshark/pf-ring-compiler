package cmd

import (
	"fmt"
	"io"

	"os"

	"github.com/kubeshark/pfring-compiler/cmd/compatibility"
	"github.com/kubeshark/pfring-compiler/cmd/compile"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	name        = "pfring-compiler"
	description = "tool to build pf_ring kernel module"
)

type Config struct {
	Logger *log.Logger
	Stderr io.Writer
	Stdout io.Writer
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, fmt.Errorf("%T.Logger must not be empty", config)
	}
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}

	var err error

	var compileCmd *cobra.Command
	{
		c := compile.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		compileCmd, err = compile.New(c)
		if err != nil {
			return nil, err
		}
	}

	var compatibilityCmd *cobra.Command
	{
		c := compatibility.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		compatibilityCmd, err = compatibility.New(c)
		if err != nil {
			return nil, err
		}
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		RunE:         r.Run,
		SilenceUsage: true,
	}

	f.Init(c)

	c.AddCommand(compileCmd)
	c.AddCommand(compatibilityCmd)

	return c, nil
}
