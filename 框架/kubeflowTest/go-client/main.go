package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	namespace = "default"
	jobName   = "stage4-go-client"
)

var trainJobGVR = schema.GroupVersionResource{
	Group: "trainer.kubeflow.org", Version: "v1alpha1", Resource: "trainjobs",
}

func main() {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	must(err)

	dynamicClient, err := dynamic.NewForConfig(config)
	must(err)
	coreClient, err := kubernetes.NewForConfig(config)
	must(err)

	jobs := dynamicClient.Resource(trainJobGVR).Namespace(namespace)
	_ = jobs.Delete(ctx, jobName, metav1.DeleteOptions{})
	must(wait.PollUntilContextTimeout(ctx, time.Second, time.Minute, true, func(ctx context.Context) (bool, error) {
		_, err := jobs.Get(ctx, jobName, metav1.GetOptions{})
		return apierrors.IsNotFound(err), nil
	}))

	created, err := jobs.Create(ctx, newTrainJob(), metav1.CreateOptions{})
	must(err)
	fmt.Printf("created TrainJob %s\n", created.GetName())

	must(wait.PollUntilContextTimeout(ctx, 2*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
		job, err := jobs.Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		conditions, _, _ := unstructured.NestedSlice(job.Object, "status", "conditions")
		for _, item := range conditions {
			condition := item.(map[string]interface{})
			conditionType, _, _ := unstructured.NestedString(condition, "type")
			status, _, _ := unstructured.NestedString(condition, "status")
			if conditionType == "Failed" && status == "True" {
				return false, fmt.Errorf("TrainJob failed")
			}
			if conditionType == "Complete" && status == "True" {
				return true, nil
			}
		}
		fmt.Println("waiting for TrainJob...")
		return false, nil
	}))

	pods, err := coreClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "jobset.sigs.k8s.io/jobset-name=" + jobName,
	})
	must(err)
	if len(pods.Items) == 0 {
		must(fmt.Errorf("no Pod found for %s", jobName))
	}

	var logs []byte
	must(wait.PollUntilContextTimeout(ctx, time.Second, 30*time.Second, true, func(ctx context.Context) (bool, error) {
		var err error
		logs, err = coreClient.CoreV1().Pods(namespace).
			GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{}).
			DoRaw(ctx)
		if err != nil {
			fmt.Printf("training finished; waiting for logs: %v\n", err)
			return false, nil
		}
		return true, nil
	}))
	fmt.Printf("\n--- training logs ---\n%s", logs)
}

func newTrainJob() *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "trainer.kubeflow.org/v1alpha1",
		"kind":       "TrainJob",
		"metadata": map[string]interface{}{
			"name": jobName,
			"labels": map[string]interface{}{
				"learning.kubeflow.org/stage": "4",
			},
		},
		"spec": map[string]interface{}{
			"runtimeRef": map[string]interface{}{
				"kind": "TrainingRuntime",
				"name": "single-gpu-runtime",
			},
			"trainer": map[string]interface{}{
				"numNodes": int64(1),
				"image":    "python:3.11-slim",
				"command": []interface{}{
					"python", "/workspace/train.py",
				},
				"env": []interface{}{
					map[string]interface{}{"name": "STAGE", "value": "04-go-client"},
				},
				"resourcesPerNode": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu": "500m", "memory": "1Gi", "nvidia.com/gpu": "1",
					},
					"limits": map[string]interface{}{
						"memory": "2Gi", "nvidia.com/gpu": "1",
					},
				},
			},
			"podTemplateOverrides": []interface{}{
				map[string]interface{}{
					"targetJobs": []interface{}{map[string]interface{}{"name": "node"}},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": "node",
								"volumeMounts": []interface{}{
									map[string]interface{}{
										"name": "training-code", "mountPath": "/workspace", "readOnly": true,
									},
								},
							},
						},
						"volumes": []interface{}{
							map[string]interface{}{
								"name":      "training-code",
								"configMap": map[string]interface{}{"name": "kubeflow-workflow-code"},
							},
						},
					},
				},
			},
		},
	}}
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
