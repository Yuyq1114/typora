package admissionwebhook

import (
	"encoding/json"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func requestFor(t *testing.T, pod corev1.Pod) *admissionv1.AdmissionRequest {
	t.Helper()
	raw, err := json.Marshal(pod)
	if err != nil {
		t.Fatal(err)
	}
	return &admissionv1.AdmissionRequest{Kind: metav1GroupVersionKind("Pod"), Object: runtime.RawExtension{Raw: raw}}
}

func metav1GroupVersionKind(kind string) metav1.GroupVersionKind {
	return metav1.GroupVersionKind{Group: "", Version: "v1", Kind: kind}
}

func TestMutateAddsLabel(t *testing.T) {
	response := mutatePod(requestFor(t, corev1.Pod{}))
	if !response.Allowed || len(response.Patch) == 0 {
		t.Fatalf("expected mutation patch: %#v", response)
	}
}

func TestValidateLimits(t *testing.T) {
	pod := corev1.Pod{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "app"}}}}
	if validatePod(requestFor(t, pod)).Allowed {
		t.Fatal("expected missing limits to be rejected")
	}
	pod.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
		corev1.ResourceCPU: resource.MustParse("100m"), corev1.ResourceMemory: resource.MustParse("64Mi"),
	}
	if !validatePod(requestFor(t, pod)).Allowed {
		t.Fatal("expected limits to be accepted")
	}
}
