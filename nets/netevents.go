package nets

type INetEventListener interface {
	OnConn(string)                       // 获得新链接时
	OnMsg(string, []byte)                // 有新消息时
	OnClose(string)                      // 链接断开时
	OnError(string, error)               // 链接异常时
	CheckUrlLegal(string) (string, bool) // 检测Url的合法性，并返还合法的节点id
}
