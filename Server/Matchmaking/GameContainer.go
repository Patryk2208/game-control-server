package Matchmaking

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type ContainerizationInfo struct {
	AgonesClient     *versioned.Clientset
	KubernetesClient *kubernetes.Clientset
	Namespace        string
	StdGSTemplate    *agonesv1.GameServer
}

func (gm *GameManager) CreateGameServer(ctx context.Context) (string, int, error) {
	result, err := gm.ContainerInfo.AgonesClient.AgonesV1().GameServers(gm.ContainerInfo.Namespace).Create(ctx, gm.ContainerInfo.StdGSTemplate, metav1.CreateOptions{})
	if err != nil {
		return "", 0, fmt.Errorf("failed to create GameServer: %w", err)
	}

	return gm.WaitForAllocation(ctx, result)
}

func (gm *GameManager) WatchContainerState(ctx context.Context, wgs *agonesv1.GameServer) {
	watcher, err := gm.ContainerInfo.AgonesClient.AgonesV1().GameServers(wgs.Namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	defer watcher.Stop()

	for e := range watcher.ResultChan() {
		gs, ok := e.Object.(*agonesv1.GameServer)
		if !ok || gs.Name != wgs.Name {
			continue
		}
		if gs.Status.State == agonesv1.GameServerStateError ||
			gs.Status.State == agonesv1.GameServerStateUnhealthy ||
			gs.Status.State == agonesv1.GameServerStateShutdown {
			return
		}
	}

}

func (gm *GameManager) WaitForAllocation(ctx context.Context, rgs *agonesv1.GameServer) (string, int, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			gs, err := gm.ContainerInfo.AgonesClient.AgonesV1().GameServers(gm.ContainerInfo.Namespace).Get(ctx, rgs.Name, metav1.GetOptions{})
			if err != nil {
				return "", 0, err
			}

			if gs.Status.State == agonesv1.GameServerStateAllocated {
				if len(gs.Status.Ports) == 0 {
					return "", 0, fmt.Errorf("no ports allocated")
				}

				return gs.Status.Address, int(gs.Status.Ports[0].Port), nil
			}

			if gs.Status.State == agonesv1.GameServerStateUnhealthy ||
				gs.Status.State == agonesv1.GameServerStateShutdown {
				return "", 0, fmt.Errorf("server entered %s state", gs.Status.State)
			}

		case <-timeout:
			return "", 0, fmt.Errorf("timed out waiting for allocation")

		case <-ctx.Done():
			return "", 0, ctx.Err()
		}
	}
}
