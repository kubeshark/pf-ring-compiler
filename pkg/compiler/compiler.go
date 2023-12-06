package compiler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	Target    string
	Clientset *kubernetes.Clientset

	Logger *logrus.Logger
}

type Compiler struct {
	target    string
	clientset *kubernetes.Clientset

	logger *logrus.Logger
}

func New(c Config) (*Compiler, error) {
	return &Compiler{
		clientset: c.Clientset,
		target:    c.Target,
		logger:    c.Logger,
	}, nil
}

func (c *Compiler) Compile() error {
	var err error

	namespace := "default"
	jobName := fmt.Sprintf("%s-pf-ring-compiler", c.target)

	c.logger.Infof("creating compile job %s/%s", namespace, jobName)
	_, err = c.CreateCompileJob(jobName, namespace)
	if err != nil {
		return fmt.Errorf("Compile: %w", err)
	}
	c.logger.Infof("compile job %s/%s created", namespace, jobName)

	c.logger.Infof("waiting for compile job to start")
	err = c.WaitForJobStart(jobName, namespace)
	if err != nil {
		panic(err.Error())
	}
	c.logger.Infof("compile job started")

	return nil
}
