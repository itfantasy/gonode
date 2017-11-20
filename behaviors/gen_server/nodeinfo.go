package gen_server

// 节点信息
type NodeInfo struct {
	Tag string // 节点标签
	Id  string // 节点实例名(不指定时自动生成)
	Url string // 节点url

	RedCore string // 核心redis服务器url
	RedPool int    // redis连接池最大连接数
	RedDB   int    // redis目标数据库
	RedAuth string // redis认证信息

	AutoDetect bool // 是否自检新节点

	Net         string // 网络
	AllowOrigin string // 作用域

	UserData map[string]string // 用户数据
}
