package lobby

type LobbyServer interface {
	Start()  // when start
	Update() // timer update
	OnMsg(string, []byte)
	OnReload(string) error
}
