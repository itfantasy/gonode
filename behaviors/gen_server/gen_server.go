package gen_server

type GenServer interface {
	SelfNodeInfo() *NodeInfo    // 自身节点信息（别名的自纠正，用于网关，自动负载均衡）
	IsInterestedIn(string) bool // 是否对新节点感兴趣

	Start()                 // 初始化完毕时
	Update()                // 轮询更新时
	OnConn(string)          // 获得新链接时
	OnMsg(string, []byte)   // 有新消息时
	OnClose(string)         // 链接断开时
	OnShell(string, string) // 有来自redis的订阅消息时
	OnReload(string) error  // 重新加载业务逻辑时
	CreateConnId() string   // 外来客户端随机id
}
