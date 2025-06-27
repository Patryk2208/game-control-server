package Session

/*type Context int

const (
	normal Context = iota
	authenticated
	playing
)*/

type UserConnectionContext interface {
	GetHandler(request *Request) (RequestHandler, error)
}
