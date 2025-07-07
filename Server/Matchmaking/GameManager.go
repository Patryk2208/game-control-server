package Matchmaking

import (
	"Server/Database"
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type Match struct {
	Capacity int
	Players  []*Database.PlayerDB
}

type GameManager struct {
	MatchingMutex  *sync.Mutex
	ActiveMutex    *sync.Mutex
	WaitingMatches []*Match
	ActiveGames    []*GameInstance
	DbPool         *Database.DBConnectionPool
	ContainerInfo  ContainerizationInfo
}

func NewGameManager(dbp *Database.DBConnectionPool) *GameManager {
	const maxMatchCount = 10000
	return &GameManager{
		MatchingMutex: new(sync.Mutex),
		ActiveMutex:   new(sync.Mutex),
		DbPool:        dbp,
		ContainerInfo: ContainerizationInfo{
			StdGSTemplate: &agonesv1.GameServer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "game",
				},
				Spec: agonesv1.GameServerSpec{
					Ports: []agonesv1.GameServerPort{{
						Name:          "game",
						PortPolicy:    agonesv1.Dynamic,
						ContainerPort: 5555,
						Protocol:      corev1.ProtocolTCP,
					},
						{
							Name:          "agones-sdk",
							ContainerPort: 9357,
							Protocol:      corev1.ProtocolTCP,
						}},
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  "game-server",
								Image: "rpg-game-server:local",
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("3000"),
										corev1.ResourceMemory: resource.MustParse("1Gi"),
									},
								},
								Env: []corev1.EnvVar{{
									Name: "SERVER_PORT",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.annotations['agones.dev/containerPort-game']",
										},
									},
								}},
							}},
						},
					},
				},
			},
			Namespace: "game",
		},
		WaitingMatches: make([]*Match, 0, maxMatchCount),
		ActiveGames:    make([]*GameInstance, 0, maxMatchCount),
	}
}
