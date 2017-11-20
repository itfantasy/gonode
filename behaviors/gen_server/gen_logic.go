package gen_server

type GenLogic interface {
	OnMsg(string, []byte)   // 有新消息时
	OnShell(string, string) // 有来自redis的订阅消息时
	OnReload(string) error  // 重新加载业务逻辑时
}
