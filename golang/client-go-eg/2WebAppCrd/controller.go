package webappcrd

import (
	"context"
	"fmt"
	"log"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var webAppGVR = schema.GroupVersionResource{Group: "demo.example.com", Version: "v1", Resource: "webapps"}

type Controller struct {
	client    kubernetes.Interface
	webApps   cache.GenericLister
	deploys   appslisters.DeploymentLister
	services  corelisters.ServiceLister
	webSynced cache.InformerSynced
	depSynced cache.InformerSynced
	svcSynced cache.InformerSynced
	queue     workqueue.TypedRateLimitingInterface[string]
}

func NewController(client kubernetes.Interface, kubeFactory informers.SharedInformerFactory, webFactory dynamicinformer.DynamicSharedInformerFactory) *Controller {
	webInformer := webFactory.ForResource(webAppGVR)
	depInformer := kubeFactory.Apps().V1().Deployments()
	svcInformer := kubeFactory.Core().V1().Services()
	c := &Controller{
		client: client, webApps: webInformer.Lister(), deploys: depInformer.Lister(), services: svcInformer.Lister(),
		webSynced: webInformer.Informer().HasSynced, depSynced: depInformer.Informer().HasSynced, svcSynced: svcInformer.Informer().HasSynced,
		queue: workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]()),
	}
	webInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueWebApp,
		UpdateFunc: func(_, obj interface{}) { c.enqueueWebApp(obj) },
		DeleteFunc: c.enqueueWebApp,
	})
	childHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueOwner,
		UpdateFunc: func(_, obj interface{}) { c.enqueueOwner(obj) },
		DeleteFunc: c.enqueueOwner,
	}
	depInformer.Informer().AddEventHandler(childHandler)
	svcInformer.Informer().AddEventHandler(childHandler)
	return c
}

func (c *Controller) enqueueWebApp(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		c.queue.Add(key)
	}
}

func (c *Controller) enqueueOwner(obj interface{}) {
	if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		obj = tombstone.Obj
	}
	child, ok := obj.(metav1.Object)
	if !ok {
		return
	}
	owner := metav1.GetControllerOf(child)
	if owner != nil && owner.APIVersion == "demo.example.com/v1" && owner.Kind == "WebApp" {
		c.queue.Add(child.GetNamespace() + "/" + owner.Name)
	}
}

func (c *Controller) Run(ctx context.Context) error {
	defer c.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), c.webSynced, c.depSynced, c.svcSynced) {
		return fmt.Errorf("informer caches did not sync; is the WebApp CRD installed?")
	}
	log.Print("WebApp controller ready")
	go func() {
		<-ctx.Done()
		c.queue.ShutDown()
	}()
	for {
		key, shutdown := c.queue.Get()
		if shutdown {
			return nil
		}
		if err := c.reconcile(ctx, key); err != nil {
			log.Printf("reconcile %s: %v", key, err)
			c.queue.AddRateLimited(key)
		} else {
			c.queue.Forget(key)
		}
		c.queue.Done(key)
	}
}

func (c *Controller) reconcile(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	obj, err := c.webApps.ByNamespace(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		return nil // Kubernetes garbage collection removes owned children.
	}
	if err != nil {
		return err
	}
	webApp := obj.(*unstructured.Unstructured)
	image, _, err := unstructured.NestedString(webApp.Object, "spec", "image")
	if err != nil || image == "" {
		return fmt.Errorf("WebApp %s needs spec.image", key)
	}
	replicas, found, err := unstructured.NestedInt64(webApp.Object, "spec", "replicas")
	if err != nil || replicas < 0 || replicas > 1000 {
		return fmt.Errorf("WebApp %s has invalid spec.replicas: %v", key, err)
	}
	if !found {
		replicas = 1
	}
	port, found, err := unstructured.NestedInt64(webApp.Object, "spec", "port")
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("WebApp %s has invalid spec.port: %v", key, err)
	}
	if !found {
		port = 80
	}
	if port == 0 {
		return fmt.Errorf("WebApp %s needs a nonzero spec.port", key)
	}
	if err := c.ensureDeployment(ctx, webApp, image, int32(replicas), int32(port)); err != nil {
		return err
	}
	return c.ensureService(ctx, webApp, int32(port))
}

func (c *Controller) ensureDeployment(ctx context.Context, webApp *unstructured.Unstructured, image string, replicas, port int32) error {
	namespace, name := webApp.GetNamespace(), webApp.GetName()
	labels := map[string]string{"app": name, "managed-by": "webapp-controller"}
	selector := map[string]string{"app": name}
	desired := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, OwnerReferences: ownerReference(webApp)},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "app", Image: image, Ports: []corev1.ContainerPort{{ContainerPort: port}}}}},
			},
		},
	}
	current, err := c.deploys.Deployments(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		log.Printf("creating Deployment %s/%s", namespace, name)
		_, err = c.client.AppsV1().Deployments(namespace).Create(ctx, desired, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}
	if !ownedBy(current, webApp) {
		return fmt.Errorf("Deployment %s/%s is not owned by this WebApp", namespace, name)
	}
	container := current.Spec.Template.Spec.Containers[0]
	if reflect.DeepEqual(current.Spec.Replicas, desired.Spec.Replicas) &&
		reflect.DeepEqual(current.Spec.Template.Labels, labels) &&
		container.Image == image && len(container.Ports) > 0 && container.Ports[0].ContainerPort == port {
		return nil
	}
	updated := current.DeepCopy()
	updated.Spec.Replicas = desired.Spec.Replicas
	updated.Spec.Template.Labels = labels
	updated.Spec.Template.Spec.Containers[0].Image = image
	updated.Spec.Template.Spec.Containers[0].Ports = desired.Spec.Template.Spec.Containers[0].Ports
	log.Printf("updating Deployment %s/%s", namespace, name)
	_, err = c.client.AppsV1().Deployments(namespace).Update(ctx, updated, metav1.UpdateOptions{})
	return err
}

func (c *Controller) ensureService(ctx context.Context, webApp *unstructured.Unstructured, port int32) error {
	namespace, name := webApp.GetNamespace(), webApp.GetName()
	selector := map[string]string{"app": name}
	ports := []corev1.ServicePort{{Port: port, TargetPort: intstr.FromInt32(port), Protocol: corev1.ProtocolTCP}}
	desired := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, OwnerReferences: ownerReference(webApp)},
		Spec:       corev1.ServiceSpec{Selector: selector, Ports: ports},
	}
	current, err := c.services.Services(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		log.Printf("creating Service %s/%s", namespace, name)
		_, err = c.client.CoreV1().Services(namespace).Create(ctx, desired, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}
	if !ownedBy(current, webApp) {
		return fmt.Errorf("Service %s/%s is not owned by this WebApp", namespace, name)
	}
	if reflect.DeepEqual(current.Spec.Selector, selector) && len(current.Spec.Ports) == 1 &&
		current.Spec.Ports[0].Port == port && current.Spec.Ports[0].TargetPort == ports[0].TargetPort {
		return nil
	}
	updated := current.DeepCopy() // Preserve the assigned ClusterIP.
	updated.Spec.Selector = selector
	updated.Spec.Ports = ports
	log.Printf("updating Service %s/%s", namespace, name)
	_, err = c.client.CoreV1().Services(namespace).Update(ctx, updated, metav1.UpdateOptions{})
	return err
}

func ownerReference(webApp *unstructured.Unstructured) []metav1.OwnerReference {
	gvk := schema.GroupVersionKind{Group: "demo.example.com", Version: "v1", Kind: "WebApp"}
	return []metav1.OwnerReference{*metav1.NewControllerRef(webApp, gvk)}
}

func ownedBy(child metav1.Object, webApp *unstructured.Unstructured) bool {
	owner := metav1.GetControllerOf(child)
	return owner != nil && owner.APIVersion == "demo.example.com/v1" && owner.Kind == "WebApp" &&
		owner.Name == webApp.GetName() && owner.UID == webApp.GetUID()
}
