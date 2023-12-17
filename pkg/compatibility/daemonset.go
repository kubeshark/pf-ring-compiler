package compatibility

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (c *Compatibility) createUnameDaemonSet(dsName, namespace string) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dsName,
			Namespace: namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": dsName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": dsName,
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

func (c *Compatibility) waitForDaemonSetPodsRunning(ctx context.Context, dsName, namespace string) ([]string, error) {
	var podNames []string
	timeout := 5 * time.Second
	immediate := true

	err := wait.PollUntilContextCancel(ctx, timeout, immediate, func(ctx context.Context) (bool, error) {
		pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
			LabelSelector: "app=" + dsName,
		})
		if err != nil {
			return false, err
		}

		podNames = make([]string, 0)
		allRunning := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				allRunning = false
				break
			}
			podNames = append(podNames, pod.Name)
		}

		return allRunning, nil
	})

	if err != nil {
		return nil, err
	}
	return podNames, nil
}
