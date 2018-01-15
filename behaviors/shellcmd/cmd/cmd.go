package cmd

/*
 * 内部协议
 * 固定格式:
 * gonode [cmd] [sur] [dist] (最少4字段)
 */

const (
	NEWNODE string = "newnode"
	PRINT   string = "print"
	STATUS  string = "status"
	RELOAD  string = "reload"
)

// 新节点上报
func NewNode(from string) string {
	return sprinf(NEWNODE, from, "*", "..")
}

// 控制台打印
func Print(from string, to string, txt string) string {
	return sprinf(PRINT, from, to, txt)
}

// 获取目标节点的链接信息
func Status(from string, to string) string {
	return sprinf(STATUS, from, to, "..")
}

// 目标节点重新加载程序集(热更新)
func Reload(from string, to string) string {
	return sprinf(RELOAD, from, to, "..")
}

func sprinf(cmdstr string, from string, to string, txt string) string {
	return "gonode " + cmdstr + " " + from + " " + to + " " + txt
}
