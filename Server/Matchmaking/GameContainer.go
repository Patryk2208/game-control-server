package Matchmaking

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ContainerizationInfo struct {
	AgonesClient       *versioned.Clientset
	KubernetesClient   *kubernetes.Clientset
	Namespace          string
	AllocationTemplate *allocationv1.GameServerAllocation
}

func (gm *GameManager) AllocateGameServer(ctx context.Context) (string, int, error) {
	info, err := gm.ContainerInfo.AgonesClient.AgonesV1().GameServers("game").List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", 0, err
	}
	for i := range info.Items {
		fmt.Println(info.Items[i].Name, info.Items[i].Status.State)
	}
	alloc, err := gm.ContainerInfo.AgonesClient.AllocationV1().GameServerAllocations(gm.ContainerInfo.Namespace).Create(ctx, gm.ContainerInfo.AllocationTemplate, metav1.CreateOptions{})
	if err != nil {
		return "", 0, fmt.Errorf("failed to create GameServerAllocation: %w", err)
	}
	fmt.Println("game server allocated with status:", alloc.Status, ", now getting address and port")
	address := alloc.Status.Address
	p := alloc.Status.Ports
	if len(p) == 0 {
		return "", 0, fmt.Errorf("no ports")
	}
	port := int(alloc.Status.Ports[0].Port)
	fmt.Println("allocated GameServer, address:", address, "port:", port)
	return address, port, nil
}

func (gm *GameManager) WatchContainerState(ctx context.Context) {
	watcher, err := gm.ContainerInfo.AgonesClient.AgonesV1().GameServers(gm.ContainerInfo.Namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("watcher error:", err)
		return
	}
	defer watcher.Stop()

	fmt.Println("watcher started")
	for e := range watcher.ResultChan() {
		gs, ok := e.Object.(*agonesv1.GameServer)
		if !ok {
			continue
		}
		fmt.Println("watcher registered an event: " + e.Type)
		if gs.Status.State == agonesv1.GameServerStateError ||
			gs.Status.State == agonesv1.GameServerStateUnhealthy ||
			gs.Status.State == agonesv1.GameServerStateShutdown {
			fmt.Println("game server is in state " + gs.Status.State + " thus ending...")
			return
		}
	}

}
