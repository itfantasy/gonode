package nets

type INetWorker interface {
	Listen(url string) error
	Connect(url string, origin string) error
	Send(id string, msg []byte) error
	SendAsync(id string, msg []byte)
	BindEventListener(eventListener INetEventListener) error
	IsIdExists(string) bool
	GetAllConnIds() []string
	Close(string) error
}

type Enum_NetWorkerType string

const (
	None Enum_NetWorkerType = "none"
	WS   Enum_NetWorkerType = "ws"
	TCP  Enum_NetWorkerType = "tcp"
	KCP  Enum_NetWorkerType = "kcp"
	HTTP Enum_NetWorkerType = "http"
)
