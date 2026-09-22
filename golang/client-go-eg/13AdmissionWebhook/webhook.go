package admissionwebhook

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Address  = ":9443"
	CertFile = "13AdmissionWebhook/tls/tls.crt"
	KeyFile  = "13AdmissionWebhook/tls/tls.key"
	CAFile   = "13AdmissionWebhook/tls/ca.crt"
)

func Run(ctx context.Context) error {
	if err := ensureCertificate(CertFile, KeyFile, CAFile); err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", admissionHandler(mutatePod))
	mux.HandleFunc("/validate", admissionHandler(validatePod))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	server := &http.Server{Addr: Address, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	log.Printf("admission webhook listening on https://host.docker.internal%s", Address)
	err := server.ListenAndServeTLS(CertFile, KeyFile)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

type admitFunc func(*admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse

func admissionHandler(admit admitFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var review admissionv1.AdmissionReview
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil || review.Request == nil {
			http.Error(w, "invalid AdmissionReview", http.StatusBadRequest)
			return
		}
		response := admit(review.Request)
		response.UID = review.Request.UID
		result := admissionv1.AdmissionReview{TypeMeta: metav1.TypeMeta{APIVersion: "admission.k8s.io/v1", Kind: "AdmissionReview"}, Response: response}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("write AdmissionReview: %v", err)
		}
	}
}

func mutatePod(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	if request.Kind.Kind != "Pod" {
		return allowed()
	}
	var pod corev1.Pod
	if err := json.Unmarshal(request.Object.Raw, &pod); err != nil {
		return denied("cannot decode Pod: " + err.Error())
	}
	if pod.Labels != nil {
		if _, exists := pod.Labels["mini-webhook"]; exists {
			return allowed()
		}
	}
	var patch []map[string]interface{}
	if pod.Labels == nil {
		patch = []map[string]interface{}{{"op": "add", "path": "/metadata/labels", "value": map[string]string{"mini-webhook": "mutated"}}}
	} else {
		patch = []map[string]interface{}{{"op": "add", "path": "/metadata/labels/mini-webhook", "value": "mutated"}}
	}
	data, err := json.Marshal(patch)
	if err != nil {
		return denied(err.Error())
	}
	patchType := admissionv1.PatchTypeJSONPatch
	return &admissionv1.AdmissionResponse{Allowed: true, Patch: data, PatchType: &patchType}
}

func validatePod(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	if request.Kind.Kind != "Pod" {
		return allowed()
	}
	var pod corev1.Pod
	if err := json.Unmarshal(request.Object.Raw, &pod); err != nil {
		return denied("cannot decode Pod: " + err.Error())
	}
	for _, container := range pod.Spec.Containers {
		limits := container.Resources.Limits
		if limits.Cpu().IsZero() || limits.Memory().IsZero() {
			return denied(fmt.Sprintf("container %q must set cpu and memory limits", container.Name))
		}
	}
	return allowed()
}

func allowed() *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{Allowed: true}
}

func denied(message string) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{Allowed: false, Result: &metav1.Status{Message: message}}
}

func ensureCertificate(certPath, keyPath, caPath string) error {
	if err := os.MkdirAll(filepath.Dir(certPath), 0700); err != nil {
		return err
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1")}
	if addresses, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addresses {
			if network, ok := address.(*net.IPNet); ok && network.IP.To4() != nil {
				ipAddresses = append(ipAddresses, network.IP)
			}
		}
	}
	template := &x509.Certificate{
		SerialNumber: serial, Subject: pkix.Name{CommonName: "mini-admission-webhook"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IsCA: true, BasicConstraintsValid: true,
		DNSNames: []string{"host.docker.internal", "localhost"}, IPAddresses: ipAddresses,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err := os.WriteFile(certPath, certPEM, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return err
	}
	return os.WriteFile(caPath, certPEM, 0644)
}
