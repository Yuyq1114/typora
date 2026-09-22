package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	minischeduler "client-go-eg/11MiniScheduler"
	minikubelet "client-go-eg/12MiniKubelet"
	admissionwebhook "client-go-eg/13AdmissionWebhook"
	minioperator "client-go-eg/14MiniOperatorInfra"
	webappcrd "client-go-eg/2WebAppCrd"
	crdclient "client-go-eg/3ClientGoCrd"
	"client-go-eg/cmController"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	controllerChoice := 11 // Change this number to run a different stage.
	if value := os.Getenv("CONTROLLER_CHOICE"); value != "" {
		selected, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("invalid CONTROLLER_CHOICE %q: %v", value, err)
		}
		controllerChoice = selected
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig")
	flag.Parse()

	var config *rest.Config
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err == nil {
			// GoLand's proxy must not intercept the local Docker Desktop API.
			config.Proxy = func(req *http.Request) (*url.URL, error) {
				if req.URL.Hostname() == "kubernetes.docker.internal" {
					return nil, nil
				}
				return http.ProxyFromEnvironment(req)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	factory := informers.NewSharedInformerFactory(client, 0)
	var runErr error
	switch controllerChoice {
	case 1:
		controller := cmController.New(client, factory)
		factory.Start(ctx.Done())
		runErr = controller.Run(ctx)
	case 2:
		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		webAppFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 0)
		controller := webappcrd.NewController(client, factory, webAppFactory)
		factory.Start(ctx.Done())
		webAppFactory.Start(ctx.Done())
		runErr = controller.Run(ctx)
	case 3:
		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		runErr = crdclient.Run(ctx, dynamicClient)
	case 11:
		scheduler := minischeduler.New(client, factory)
		factory.Start(ctx.Done())
		runErr = scheduler.Run(ctx)
	case 12:
		kubelet := minikubelet.New(factory)
		factory.Start(ctx.Done())
		runErr = kubelet.Run(ctx)
	case 13:
		runErr = admissionwebhook.Run(ctx)
	case 14:
		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		namespace := os.Getenv("POD_NAMESPACE")
		if namespace == "" {
			namespace = "mini-operator-system"
		}
		runErr = minioperator.Run(ctx, client, dynamicClient, namespace)
	default:
		log.Fatalf("unknown controllerChoice: %d", controllerChoice)
	}
	if runErr != nil && ctx.Err() == nil {
		log.Fatal(runErr)
	}
}
