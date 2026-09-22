package cmController

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const labelKey = "app-controller"

type Controller struct {
	client       kubernetes.Interface
	configMaps   corelisters.ConfigMapLister
	deployments  appslisters.DeploymentLister
	configSynced cache.InformerSynced
	deploySynced cache.InformerSynced
	queue        workqueue.TypedRateLimitingInterface[string]
}

func New(client kubernetes.Interface, factory informers.SharedInformerFactory) *Controller {
	configInformer := factory.Core().V1().ConfigMaps()
	deployInformer := factory.Apps().V1().Deployments()
	c := &Controller{
		client:       client,
		configMaps:   configInformer.Lister(),
		deployments:  deployInformer.Lister(),
		configSynced: configInformer.Informer().HasSynced,
		deploySynced: deployInformer.Informer().HasSynced,
		queue:        workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]()),
	}

	// Event handlers do little work: they only put a key on the queue.
	configInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueConfigMap,
		UpdateFunc: func(_, obj interface{}) { c.enqueueConfigMap(obj) },
		DeleteFunc: c.enqueueConfigMap,
	})
	deployInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueDeploymentOwner,
		UpdateFunc: func(_, obj interface{}) { c.enqueueDeploymentOwner(obj) },
		DeleteFunc: c.enqueueDeploymentOwner,
	})
	return c
}

func (c *Controller) enqueueConfigMap(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Printf("could not get ConfigMap key: %v", err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) enqueueDeploymentOwner(obj interface{}) {
	if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		obj = tombstone.Obj
	}
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return
	}
	owner := metav1.GetControllerOf(deployment)
	if owner != nil && owner.APIVersion == "v1" && owner.Kind == "ConfigMap" {
		c.queue.Add(deployment.Namespace + "/" + owner.Name)
	}
}

func (c *Controller) Run(ctx context.Context) error {
	defer c.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), c.configSynced, c.deploySynced) {
		return fmt.Errorf("informer caches did not sync")
	}
	log.Print("caches synced; controller is watching ConfigMaps")
	go func() {
		<-ctx.Done()
		c.queue.ShutDown()
	}()
	for c.processNext(ctx) {
	}
	return nil
}

func (c *Controller) processNext(ctx context.Context) bool {
	key, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(key)
	if err := c.reconcile(ctx, key); err != nil {
		log.Printf("reconcile %s failed: %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
	}
	return true
}

func (c *Controller) reconcile(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// Listers read the informers' local cache, not the API server.
	cm, err := c.configMaps.ConfigMaps(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		cm = nil
	} else if err != nil {
		return err
	}
	deployment, err := c.deployments.Deployments(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		deployment = nil
	} else if err != nil {
		return err
	}

	// Removing the label or ConfigMap removes only a Deployment we own.
	if cm == nil || cm.Labels[labelKey] != "true" {
		if deployment != nil && ownedByConfigMap(deployment, name, "") {
			log.Printf("deleting Deployment %s", key)
			err := c.client.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		return nil
	}

	image := cm.Data["image"]
	if image == "" {
		return fmt.Errorf("ConfigMap %s needs data.image", key)
	}
	replicas := int32(1)
	if raw := cm.Data["replicas"]; raw != "" {
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || n < 0 {
			return fmt.Errorf("ConfigMap %s has invalid data.replicas %q", key, raw)
		}
		replicas = int32(n)
	}
	desired := desiredDeployment(cm, image, replicas)
	if deployment == nil {
		log.Printf("creating Deployment %s", key)
		_, err = c.client.AppsV1().Deployments(namespace).Create(ctx, desired, metav1.CreateOptions{})
		return err
	}
	if !ownedByConfigMap(deployment, name, string(cm.UID)) {
		return fmt.Errorf("Deployment %s exists but is not owned by this ConfigMap", key)
	}
	currentImage := ""
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		currentImage = deployment.Spec.Template.Spec.Containers[0].Image
	}
	if reflect.DeepEqual(deployment.Spec.Replicas, desired.Spec.Replicas) &&
		reflect.DeepEqual(deployment.Spec.Template.Labels, desired.Spec.Template.Labels) &&
		currentImage == image {
		return nil
	}
	updated := deployment.DeepCopy()
	updated.Spec.Replicas = desired.Spec.Replicas
	updated.Spec.Template.Labels = desired.Spec.Template.Labels
	if len(updated.Spec.Template.Spec.Containers) == 0 {
		updated.Spec.Template.Spec.Containers = desired.Spec.Template.Spec.Containers
	} else {
		updated.Spec.Template.Spec.Containers[0].Image = image
	}
	log.Printf("updating Deployment %s", key)
	_, err = c.client.AppsV1().Deployments(namespace).Update(ctx, updated, metav1.UpdateOptions{})
	return err
}

func ownedByConfigMap(deployment *appsv1.Deployment, name string, uid string) bool {
	owner := metav1.GetControllerOf(deployment)
	return owner != nil && owner.APIVersion == "v1" && owner.Kind == "ConfigMap" &&
		owner.Name == name && (uid == "" || string(owner.UID) == uid)
}

func desiredDeployment(cm *corev1.ConfigMap, image string, replicas int32) *appsv1.Deployment {
	selector := map[string]string{"app": cm.Name}
	labels := map[string]string{"app": cm.Name, "managed-by": "configmap-controller"}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: cm.Name, Namespace: cm.Namespace, Labels: labels,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(cm, corev1.SchemeGroupVersion.WithKind("ConfigMap"))},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "app", Image: image}}},
			},
		},
	}
}
