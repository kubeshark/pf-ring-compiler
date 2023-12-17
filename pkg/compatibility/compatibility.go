package compatibility

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	Clientset *kubernetes.Clientset

	Logger *logrus.Logger
}

type Compatibility struct {
	clientset *kubernetes.Clientset

	logger *logrus.Logger
}

func New(c Config) (*Compatibility, error) {
	return &Compatibility{
		clientset: c.Clientset,
		logger:    c.Logger,
	}, nil
}

func (c *Compatibility) Run() error {
	var err error

	dsName := "kubeshark-kernel-version"
	namespace := "default"

	_, err = c.createUnameDaemonSet(dsName, namespace)
	if err != nil {
		return err
	}

	return nil
}
