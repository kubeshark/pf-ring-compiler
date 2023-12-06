package compiler

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Compiler) CreateCompileJob(jobName, namespace string) (*batchv1.Job, error) {
	compileContainerImage := getCompileContainerImage(c.target)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "compile-container",
							Image: compileContainerImage,
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

func (c *Compiler) CheckJobLogsForString(jobName, namespace string) error {
	timeout := time.Minute * 5
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	searchString := "Kernel module is ready"

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout waiting for string '%s' in job logs", searchString)
		case <-ticker.C:
			pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "job-name=" + jobName,
			})
			if err != nil {
				return err
			}

			for _, pod := range pods.Items {
				req := c.clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
				logs, err := req.Stream(context.TODO())
				if err != nil {
					return err
				}
				defer logs.Close()

				buf := new(strings.Builder)
				_, err = io.Copy(buf, logs)
				if err != nil {
					return err
				}

				if strings.Contains(buf.String(), searchString) {
					fmt.Println(buf.String())
					return nil
				}
			}
		}
	}
}

func getCompileContainerImage(target string) string {
	containers := map[string]string{
		"al2": "corest/build:kubeshark-pf-ring-al2-builder",
	}

	// TODO: possible missing target
	return containers[target]
}
