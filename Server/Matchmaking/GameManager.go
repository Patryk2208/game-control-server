package Matchmaking

import (
	"Server/Communication"
	"Server/Database"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"fmt"
	"github.com/gorilla/websocket"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

type GameContainerAddress struct {
	Ip   string
	Port int
}

type GameInstance struct {
	Game              Database.GameDB
	GameInfo          Match
	ControlConnection *websocket.Conn
	GameAddress       GameContainerAddress
}

type MatchPlayer struct {
	Player       *Database.PlayerDB
	ReplyChannel *chan Communication.Reply
	ReplyMutex   *sync.Mutex
}

type Match struct {
	Capacity int
	Players  []*MatchPlayer
}

type GameManager struct {
	MatchingMutex  sync.Mutex
	ActiveMutex    sync.Mutex
	WaitingMatches []*Match
	ActiveGames    []*GameInstance
	DbPool         *Database.DBConnectionPool
	ContainerInfo  ContainerizationInfo
}

func NewGameManager(dbp *Database.DBConnectionPool) (*GameManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Agones client: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	fmt.Println("correct agones and kubernetes client creation")
	const maxMatchCount = 10000
	return &GameManager{
		MatchingMutex: sync.Mutex{},
		ActiveMutex:   sync.Mutex{},
		DbPool:        dbp,
		ContainerInfo: ContainerizationInfo{
			AgonesClient:     agonesClient,
			KubernetesClient: kubeClient,
			Namespace:        "game",
			AllocationTemplate: &allocationv1.GameServerAllocation{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "alloc-",
					Namespace:    "game",
				},
				Spec: allocationv1.GameServerAllocationSpec{
					Selectors: []allocationv1.GameServerSelector{{
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"agones.dev/fleet": "game-server-fleet",
							},
						},
					}},
				},
			},
		},
		WaitingMatches: make([]*Match, 0, maxMatchCount),
		ActiveGames:    make([]*GameInstance, 0, maxMatchCount),
	}, nil
}
