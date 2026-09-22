package crdclient

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var webApps = schema.GroupVersionResource{
	Group: "demo.example.com", Version: "v1", Resource: "webapps",
}

func Run(ctx context.Context, client dynamic.Interface) error {
	const namespace = "test1"
	const name = "webapp-demo"
	webAppClient := client.Resource(webApps).Namespace(namespace)

	// A custom resource has no built-in clientset method. The dynamic client
	// returns an Unstructured object whose fields can be read by name.
	webApp, err := webAppClient.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get WebApp %s/%s: %w", namespace, name, err)
	}
	image, _, err := unstructured.NestedString(webApp.Object, "spec", "image")
	if err != nil {
		return err
	}
	replicas, _, err := unstructured.NestedInt64(webApp.Object, "spec", "replicas")
	if err != nil {
		return err
	}
	port, _, err := unstructured.NestedInt64(webApp.Object, "spec", "port")
	if err != nil {
		return err
	}
	fmt.Printf("Get: %s/%s image=%s replicas=%d port=%d\n", namespace, name, image, replicas, port)

	list, err := webAppClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("list WebApps in %s: %w", namespace, err)
	}
	fmt.Printf("List: %d WebApp(s) in %s\n", len(list.Items), namespace)
	return nil
}
