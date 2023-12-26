package compatibility

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	timeout := 2 * time.Minute
	jobRunId := fmt.Sprintf("job-run-%d", time.Now().Unix())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c.logger.Debugf("creating %s/%s daemonset", namespace, dsName)
	_, err = c.createUnameDaemonSet(dsName, namespace, jobRunId)
	if err != nil {
		return err
	}
	c.logger.Debugf("daemonset %s/%s created", namespace, dsName)

	c.logger.Debugf("waiting for pods in %s/%s daemonset to start", namespace, dsName)
	pods, err := c.waitForDaemonSetPodsRunning(ctx, dsName, namespace, jobRunId)
	if err != nil {
		return err
	}
	c.logger.Debugf("all pods in %s/%s daemonset started", namespace, dsName)

	kernelVersions, err := c.getKernelVersions(pods, namespace)
	if err != nil {
		return err
	}

	var reportItems []ReportData
	{
		for i, kernelVersion := range kernelVersions {
			isSupported, err := isSupportedKernelVersion(kernelVersion)
			if err != nil {
				return err
			}
			reportItems = append(reportItems, ReportData{
				NodeName:      pods[i].Spec.NodeName,
				KernelVersion: kernelVersion,
				IsSupported:   isSupported,
			})
		}
	}

	printReportTable(reportItems)

	c.logger.Debugf("cleaning up %s/%s daemonset", namespace, dsName)
	deleteOptions := metav1.DeleteOptions{}
	err = c.clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), dsName, deleteOptions)
	if err != nil {
		return err
	}
	c.logger.Debugf("daemonset %s/%s cleaned up", namespace, dsName)

	return nil
}
