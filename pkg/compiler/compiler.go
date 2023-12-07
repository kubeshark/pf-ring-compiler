package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"

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
	// to make sure logs from old runs aren't retrieved
	jobRunId := fmt.Sprintf("job-run-%d", time.Now().Unix())

	// Setup signal handling for SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func(logger *logrus.Logger) {
		<-sigChan
		logger.Info("SIGINT received, cleaning compile job")
		err := c.CleanupJob(jobName, namespace)
		if err != nil {
			logger.Errorf("failed to cleanup compile job: %v", err)
		}
		os.Exit(1)
	}(c.logger)

	// Check if the job already exists
	exists, err := c.JobExists(jobName, namespace)
	if err != nil {
		return err
	}

	if exists {
		c.logger.Info("compile job already exists, attempting to delete it")
		err = c.CleanupJob(jobName, namespace)
		if err != nil {
			return err
		}
		c.logger.Info("cleanup up existing compile job")
	}

	// Create compiler job
	c.logger.Infof("creating compile job %s/%s", namespace, jobName)
	_, err = c.CreateCompileJob(jobName, namespace, jobRunId)
	if err != nil {
		return fmt.Errorf("Compile: %w", err)
	}
	c.logger.Infof("compile job %s/%s created", namespace, jobName)

	// Wait for compiler job to start
	c.logger.Infof("waiting for compile job to start")
	err = c.WaitForJobStart(jobName, namespace)
	if err != nil {
		return fmt.Errorf("Compile: %w", err)
	}
	c.logger.Infof("compile job started")

	// Wait for compiler job to reach completed status
	c.logger.Infof("waiting for compile job to complete pf-ring module compilation")
	fileName, err := c.CheckJobLogsForString(jobName, namespace, jobRunId)
	if err != nil {
		return fmt.Errorf("Compile: %w", err)
	}
	c.logger.Infof("pf-ring module compilation completed")

	// Copy file to local fs
	c.logger.Info("copying kernel module to local fs")
	podName, err := c.GetPodNameFromJob(jobName, namespace, jobRunId)
	if err != nil {
		return err
	}
	// TODO: configurable copy path?
	filePath := path.Join("/tmp", fileName)
	err = copyFileFromPod(podName, namespace, filePath, fileName)
	if err != nil {
		return err
	}
	c.logger.Infof("kernel module copied to %s", fileName)

	// Cleanup job after completion
	c.logger.Infof("cleaning up compile job")
	err = c.CleanupJob(jobName, namespace)
	if err != nil {
		return fmt.Errorf("Compile: %w", err)
	}
	c.logger.Infof("compile job clean up completed")

	return nil
}

func copyFileFromPod(podName, namespace, filePath, localPath string) error {
	cmd := exec.Command("kubectl", "cp", fmt.Sprintf("%s/%s:%s", namespace, podName, filePath), localPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy file from pod: %v", err)
	}
	return nil
}

func getCompileContainerImage(target string) string {
	containers := map[string]string{
		"al2": "corest/build:kubeshark-pf-ring-al2-builder",
	}

	// TODO: possible missing target
	return containers[target]
}
