package cmd

import (
	"context"
	"io"
	"path/filepath"

	"github.com/kubeshark/pfring-compiler/pkg/compiler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	var compileRunner *compiler.Compiler
	{
		c := compiler.Config{
			Target:    r.flag.Target,
			Clientset: clientset,
			Logger:    r.logger,
		}

		compileRunner, err = compiler.New(c)
		if err != nil {
			return err
		}
	}

	err = compileRunner.Compile()
	if err != nil {
		return err
	}

	return nil
}

// GetKubernetesClient creates and returns a new Kubernetes client
func getKubernetesClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
