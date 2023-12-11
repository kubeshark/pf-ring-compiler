package compiler

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Compiler) CreateCompileJob(jobName, namespace, jobRunId string) (*batchv1.Job, error) {
	compileContainerImage := getCompileContainerImage(c.target)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"job-run-id": jobRunId,
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"job-run-id": jobRunId,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "compile-container",
							Image:           compileContainerImage,
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return c.clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
}

func (c *Compiler) WaitForJobStart(jobName, namespace string) error {
	timeout := time.Minute * 3
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout waiting for job %s/%s to start", namespace, jobName)
		case <-ticker.C:
			job, err := c.clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if job.Status.Active > 0 {
				return nil
			}
		}
	}
}

func (c *Compiler) WaitForContainerToStart(jobName, namespace, jobRunId string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.Minute * 2
	timeoutChan := time.After(timeout)
	labelSelector := fmt.Sprintf("job-name=%s,job-run-id=%s", jobName, jobRunId)

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout waiting for container in job %s/%s to start",
				namespace, jobName)
		case <-ticker.C:
			pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: labelSelector,
			})
			if err != nil {
				return err
			}

			for _, pod := range pods.Items {
				if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].State.Running != nil {
					return nil
				}
			}
		}
	}
}

func (c *Compiler) CheckJobLogsForString(jobName, namespace, jobRunId string) (string, error) {
	timeout := time.Minute * 5
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	searchString := "Kernel module is ready"

	for {
		select {
		case <-timeoutChan:
			return "", fmt.Errorf("timeout waiting for string '%s' in job logs", searchString)
		case <-ticker.C:
			podName, err := c.GetPodNameFromJob(jobName, namespace, jobRunId)
			if err != nil {
				return "", err
			}

			req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
			logs, err := req.Stream(context.TODO())
			if err != nil {
				return "", err
			}
			defer logs.Close()

			buf := new(strings.Builder)
			_, err = io.Copy(buf, logs)
			if err != nil {
				return "", err
			}

			logString := buf.String()
			if strings.Contains(logString, searchString) {
				// Extract file path from the log string
				parts := strings.Split(logString, "/")
				filePath := strings.TrimSpace(parts[len(parts)-1])
				return filePath, nil
			}

		}
	}
}

func (c *Compiler) CleanupJob(jobName, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := c.clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}

	// Polling to confirm deletion
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(1 * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout reached while waiting for job deletion")
		case <-ticker.C:
			_, err := c.clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				c.logger.Errorf("error checking job status: %v", err)
			}
		}
	}
}

func (c *Compiler) JobExists(jobName, namespace string) (bool, error) {
	_, err := c.clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Compiler) GetPodNameFromJob(jobName, namespace, jobRunId string) (string, error) {
	labelSelector := fmt.Sprintf("job-name=%s,job-run-id=%s", jobName, jobRunId)

	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found for job %s/%s", namespace, jobName)
	}

	// TODO issues with using direct array index?
	return pods.Items[0].Name, nil
}
