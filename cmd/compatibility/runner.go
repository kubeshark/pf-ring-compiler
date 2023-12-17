package compatibility

import (
	"context"
	"io"

	"github.com/kubeshark/pfring-compiler/pkg/compatibility"
	"github.com/kubeshark/pfring-compiler/pkg/k8sclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runner struct {
	flag   *flag
	logger *logrus.Logger
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
	var err error

	logLevel, _ := logrus.ParseLevel(r.flag.LogLevel) // validated in init already
	r.logger.SetLevel(logLevel)

	clientset, err := k8sclient.New()
	if err != nil {
		return err
	}

	var compatibilityRunner *compatibility.Compatibility
	{
		c := compatibility.Config{
			Clientset: clientset,
			Logger:    r.logger,
		}

		compatibilityRunner, err = compatibility.New(c)
		if err != nil {
			return err
		}
	}

	err = compatibilityRunner.Run()
	if err != nil {
		return err
	}

	return nil
}
