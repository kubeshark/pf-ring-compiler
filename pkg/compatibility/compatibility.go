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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c.logger.Infof("creating %s/%s daemonset", namespace, dsName)
	_, err = c.createUnameDaemonSet(dsName, namespace)
	if err != nil {
		return err
	}
	c.logger.Infof("daemonset %s/%s created", namespace, dsName)

	podNames, err := c.waitForDaemonSetPodsRunning(ctx, dsName, namespace)
	if err != nil {
		return err
	}

	fmt.Println(podNames)

	c.logger.Infof("cleaning up %s/%s daemonset", namespace, dsName)
	deleteOptions := metav1.DeleteOptions{}
	err = c.clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), dsName, deleteOptions)
	if err != nil {
		return err
	}
	c.logger.Infof("daemonset %s/%s cleaned up", namespace, dsName)

	return nil
}
