package Session

import "Server/Communication"

/*type Context int

const (
	normal Context = iota
	authenticated
	playing
)*/

type UserConnectionContext interface {
	GetHandler(request *Communication.Request) (RequestHandler, error)
}
