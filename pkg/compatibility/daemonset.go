package compatibility

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (c *Compatibility) createUnameDaemonSet(dsName, namespace, jobRunId string) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dsName,
			Namespace: namespace,
			Labels: map[string]string{
				"job-run-id": jobRunId,
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        dsName,
					"job-run-id": jobRunId,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        dsName,
						"job-run-id": jobRunId,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "alpine",
							Image:   "alpine:3.18",
							Command: []string{"sh", "-c", "uname -r && sleep infinity"},
						},
					},
				},
			},
		},
	}

	ds, err := c.clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), ds, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("createUnameDaemonSet: failed to create daemonset: %w", err)
	}

	return ds, nil
}

func (c *Compatibility) waitForDaemonSetPodsRunning(ctx context.Context, dsName, namespace, jobRunId string) ([]corev1.Pod, error) {
	var dsPods []corev1.Pod
	timeout := 5 * time.Second
	immediate := true
	labelSelector := fmt.Sprintf("app=%s,job-run-id=%s", dsName, jobRunId)

	err := wait.PollUntilContextCancel(ctx, timeout, immediate, func(ctx context.Context) (bool, error) {
		pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return false, err
		}

		allRunning := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				allRunning = false
				break
			}
			dsPods = append(dsPods, pod)
		}

		return allRunning, nil
	})

	if err != nil {
		return nil, err
	}
	return dsPods, nil
}

func (c *Compatibility) getKernelVersions(pods []corev1.Pod, namespace string) ([]string, error) {
	var logs []string
	for _, pod := range pods {
		podLogOpts := corev1.PodLogOptions{}
		req := c.clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &podLogOpts)
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			return nil, err
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return nil, err
		}

		log := strings.TrimRight(buf.String(), "\r\n")
		logs = append(logs, log)
	}
	return logs, nil
}
