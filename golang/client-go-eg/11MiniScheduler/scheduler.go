package minischeduler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const SchedulerName = "mini-scheduler"

type Scheduler struct {
	client     kubernetes.Interface
	pods       corelisters.PodLister
	nodes      corelisters.NodeLister
	podsSynced cache.InformerSynced
	nodeSynced cache.InformerSynced
	queue      workqueue.TypedRateLimitingInterface[string]
}

func New(client kubernetes.Interface, factory informers.SharedInformerFactory) *Scheduler {
	podInformer := factory.Core().V1().Pods()
	nodeInformer := factory.Core().V1().Nodes()
	s := &Scheduler{
		client: client, pods: podInformer.Lister(), nodes: nodeInformer.Lister(),
		podsSynced: podInformer.Informer().HasSynced, nodeSynced: nodeInformer.Informer().HasSynced,
		queue: workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]()),
	}
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    s.enqueue,
		UpdateFunc: func(_, obj interface{}) { s.enqueue(obj) },
	})
	return s
}

func (s *Scheduler) enqueue(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok || pod.Spec.SchedulerName != SchedulerName || pod.Spec.NodeName != "" || pod.DeletionTimestamp != nil {
		return
	}
	key, err := cache.MetaNamespaceKeyFunc(pod)
	if err == nil {
		s.queue.Add(key)
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	defer s.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), s.podsSynced, s.nodeSynced) {
		return fmt.Errorf("scheduler caches did not sync")
	}
	log.Print("mini scheduler ready")
	go func() { <-ctx.Done(); s.queue.ShutDown() }()
	for {
		key, shutdown := s.queue.Get()
		if shutdown {
			return nil
		}
		err := s.schedule(ctx, key)
		if err != nil {
			log.Printf("schedule %s: %v", key, err)
			s.queue.AddRateLimited(key)
		} else {
			s.queue.Forget(key)
		}
		s.queue.Done(key)
	}
}

func (s *Scheduler) schedule(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	pod, err := s.pods.Pods(namespace).Get(name)
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil || pod.Spec.NodeName != "" {
		return err
	}
	nodes, err := s.nodes.List(labels.Everything())
	if err != nil {
		return err
	}
	allPods, err := s.pods.List(labels.Everything())
	if err != nil {
		return err
	}
	var selected *corev1.Node
	bestScore := int64(-1)
	for _, node := range nodes {
		if !fits(pod, node, allPods) {
			continue
		}
		score := scoreNode(pod, node, allPods)
		if score > bestScore {
			selected, bestScore = node, score
		}
	}
	if selected == nil {
		return fmt.Errorf("no node fits pod")
	}
	binding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{Name: pod.Name, Namespace: pod.Namespace, UID: pod.UID},
		Target:     corev1.ObjectReference{APIVersion: "v1", Kind: "Node", Name: selected.Name},
	}
	if err := s.client.CoreV1().Pods(namespace).Bind(ctx, binding, metav1.CreateOptions{}); err != nil {
		return err
	}
	log.Printf("bound %s to %s (score=%d)", key, selected.Name, bestScore)
	return nil
}

func fits(pod *corev1.Pod, node *corev1.Node, pods []*corev1.Pod) bool {
	if node.Spec.Unschedulable || !nodeReady(node) || !labels.SelectorFromSet(pod.Spec.NodeSelector).Matches(labels.Set(node.Labels)) {
		return false
	}
	if !toleratesNode(pod, node) || !matchesRequiredAffinity(pod, node) {
		return false
	}
	wantCPU, wantMemory := requests(pod)
	usedCPU, usedMemory := usedOnNode(node.Name, pods)
	return usedCPU+wantCPU <= node.Status.Allocatable.Cpu().MilliValue() && usedMemory+wantMemory <= node.Status.Allocatable.Memory().Value()
}

func scoreNode(pod *corev1.Pod, node *corev1.Node, pods []*corev1.Pod) int64 {
	wantCPU, wantMemory := requests(pod)
	usedCPU, usedMemory := usedOnNode(node.Name, pods)
	allocCPU, allocMemory := node.Status.Allocatable.Cpu().MilliValue(), node.Status.Allocatable.Memory().Value()
	if allocCPU == 0 || allocMemory == 0 {
		return 0
	}
	return (allocCPU-usedCPU-wantCPU)*100/allocCPU + (allocMemory-usedMemory-wantMemory)*100/allocMemory
}

func usedOnNode(node string, pods []*corev1.Pod) (int64, int64) {
	var cpu, memory int64
	for _, pod := range pods {
		if pod.Spec.NodeName != node || pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		podCPU, podMemory := requests(pod)
		cpu, memory = cpu+podCPU, memory+podMemory
	}
	return cpu, memory
}

func requests(pod *corev1.Pod) (int64, int64) {
	var cpu, memory int64
	for _, container := range pod.Spec.Containers {
		cpu += container.Resources.Requests.Cpu().MilliValue()
		memory += container.Resources.Requests.Memory().Value()
	}
	var initCPU, initMemory int64
	for _, container := range pod.Spec.InitContainers {
		if value := container.Resources.Requests.Cpu().MilliValue(); value > initCPU {
			initCPU = value
		}
		if value := container.Resources.Requests.Memory().Value(); value > initMemory {
			initMemory = value
		}
	}
	if initCPU > cpu {
		cpu = initCPU
	}
	if initMemory > memory {
		memory = initMemory
	}
	if pod.Spec.Overhead != nil {
		cpu += pod.Spec.Overhead.Cpu().MilliValue()
		memory += pod.Spec.Overhead.Memory().Value()
	}
	return cpu, memory
}

func nodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func toleratesNode(pod *corev1.Pod, node *corev1.Node) bool {
	for _, taint := range node.Spec.Taints {
		if taint.Effect != corev1.TaintEffectNoSchedule && taint.Effect != corev1.TaintEffectNoExecute {
			continue
		}
		tolerated := false
		for _, toleration := range pod.Spec.Tolerations {
			if toleration.Effect != "" && toleration.Effect != taint.Effect {
				continue
			}
			if toleration.Key == taint.Key && (toleration.Operator == corev1.TolerationOpExists || toleration.Value == taint.Value) {
				tolerated = true
				break
			}
		}
		if !tolerated {
			return false
		}
	}
	return true
}

func matchesRequiredAffinity(pod *corev1.Pod, node *corev1.Node) bool {
	affinity := pod.Spec.Affinity
	if affinity == nil || affinity.NodeAffinity == nil || affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		return true
	}
	for _, term := range affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
		if matchesTerm(term, node) {
			return true
		}
	}
	return false
}

func matchesTerm(term corev1.NodeSelectorTerm, node *corev1.Node) bool {
	for _, expression := range term.MatchExpressions {
		if !matchesRequirement(node.Labels[expression.Key], expression, node.Labels) {
			return false
		}
	}
	for _, field := range term.MatchFields {
		value := ""
		if field.Key == "metadata.name" {
			value = node.Name
		}
		if !matchesRequirement(value, field, map[string]string{field.Key: value}) {
			return false
		}
	}
	return true
}

func matchesRequirement(value string, requirement corev1.NodeSelectorRequirement, values map[string]string) bool {
	_, exists := values[requirement.Key]
	switch requirement.Operator {
	case corev1.NodeSelectorOpIn:
		return exists && contains(requirement.Values, value)
	case corev1.NodeSelectorOpNotIn:
		return !exists || !contains(requirement.Values, value)
	case corev1.NodeSelectorOpExists:
		return exists
	case corev1.NodeSelectorOpDoesNotExist:
		return !exists
	case corev1.NodeSelectorOpGt, corev1.NodeSelectorOpLt:
		if len(requirement.Values) != 1 || !exists {
			return false
		}
		left, err1 := strconv.Atoi(value)
		right, err2 := strconv.Atoi(requirement.Values[0])
		if err1 != nil || err2 != nil {
			return false
		}
		return requirement.Operator == corev1.NodeSelectorOpGt && left > right || requirement.Operator == corev1.NodeSelectorOpLt && left < right
	default:
		return false
	}
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
