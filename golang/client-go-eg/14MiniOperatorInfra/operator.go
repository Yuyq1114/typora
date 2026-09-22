package minioperator

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	webappcrd "client-go-eg/2WebAppCrd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func Run(ctx context.Context, client kubernetes.Interface, dynamicClient dynamic.Interface, namespace string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	identity := hostname + "_" + string(uuid.NewUUID())
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{Name: "webapp-controller-leader", Namespace: namespace},
		Client:    client.CoordinationV1(), LockConfig: resourcelock.ResourceLockConfig{Identity: identity},
	}
	leaderCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	elector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock: lock, LeaseDuration: 15 * time.Second, RenewDeadline: 10 * time.Second, RetryPeriod: 2 * time.Second,
		ReleaseOnCancel: true,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(runCtx context.Context) {
				log.Printf("%s became leader", identity)
				kubeFactory := informers.NewSharedInformerFactory(client, 0)
				webFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 0)
				controller := webappcrd.NewController(client, kubeFactory, webFactory)
				kubeFactory.Start(runCtx.Done())
				webFactory.Start(runCtx.Done())
				if runErr := controller.Run(runCtx); runErr != nil && runCtx.Err() == nil {
					select {
					case errCh <- runErr:
					default:
					}
					cancel()
				}
			},
			OnStoppedLeading: func() { log.Printf("%s stopped leading", identity) },
			OnNewLeader:      func(newLeader string) { log.Printf("leader is %s", newLeader) },
		},
	})
	if err != nil {
		return fmt.Errorf("create leader elector: %w", err)
	}
	elector.Run(leaderCtx)
	select {
	case runErr := <-errCh:
		return runErr
	default:
		return nil
	}
}
