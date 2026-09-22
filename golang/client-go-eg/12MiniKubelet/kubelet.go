package minikubelet

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const NodeName = "mini-node"

type Kubelet struct {
	pods      corelisters.PodLister
	podSynced cache.InformerSynced
	queue     workqueue.TypedRateLimitingInterface[string]
	running   map[string]string
}

func New(factory informers.SharedInformerFactory) *Kubelet {
	podInformer := factory.Core().V1().Pods()
	k := &Kubelet{
		pods: podInformer.Lister(), podSynced: podInformer.Informer().HasSynced,
		queue:   workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]()),
		running: map[string]string{},
	}
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    k.enqueue,
		UpdateFunc: func(_, obj interface{}) { k.enqueue(obj) },
		DeleteFunc: k.enqueue,
	})
	return k
}

func (k *Kubelet) enqueue(obj interface{}) {
	if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		obj = tombstone.Obj
	}
	pod, ok := obj.(*corev1.Pod)
	if !ok || pod.Spec.NodeName != NodeName {
		return
	}
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(pod)
	if err == nil {
		k.queue.Add(key)
	}
}

func (k *Kubelet) Run(ctx context.Context) error {
	defer k.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), k.podSynced) {
		return fmt.Errorf("kubelet cache did not sync")
	}
	log.Printf("mini kubelet simulator ready for node %q", NodeName)
	go func() { <-ctx.Done(); k.queue.ShutDown() }()
	for {
		key, shutdown := k.queue.Get()
		if shutdown {
			return nil
		}
		if err := k.syncPod(key); err != nil {
			log.Printf("sync pod %s: %v", key, err)
			k.queue.AddRateLimited(key)
		} else {
			k.queue.Forget(key)
		}
		k.queue.Done(key)
	}
}

func (k *Kubelet) syncPod(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	pod, err := k.pods.Pods(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		if uid, exists := k.running[key]; exists {
			log.Printf("CRI StopPodSandbox pod=%s uid=%s", key, uid)
			log.Printf("CRI RemovePodSandbox pod=%s uid=%s", key, uid)
			delete(k.running, key)
		}
		return nil
	}
	if err != nil {
		return err
	}
	uid := string(pod.UID)
	if k.running[key] == uid {
		return nil
	}
	log.Printf("CRI RunPodSandbox pod=%s uid=%s", key, uid)
	for _, container := range pod.Spec.InitContainers {
		log.Printf("CRI CreateContainer pod=%s container=%s image=%s (init)", key, container.Name, container.Image)
		log.Printf("CRI StartContainer pod=%s container=%s (init)", key, container.Name)
	}
	for _, container := range pod.Spec.Containers {
		log.Printf("CRI CreateContainer pod=%s container=%s image=%s", key, container.Name, container.Image)
		log.Printf("CRI StartContainer pod=%s container=%s", key, container.Name)
	}
	k.running[key] = uid
	log.Printf("status simulation: pod=%s phase=Running", key)
	return nil
}
