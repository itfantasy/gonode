package gonode

import (
	"fmt"
	"gonode/cmd"
	"strings"
)

func (this *GoNode) onShell(channel string, msg string) {
	defer this.autoRecover()
	if channel == GONODE_PUB_CHAN {
		// handle the cmd
		fmt.Println("recive the shellcmd:" + msg)
		msgInfos := strings.Split(msg, " ")
		if len(msgInfos) < 5 || msgInfos[0] != "gothis" {
			// illegel pub channel msg
			return
		}

		// ----------parse the params------------
		cmdstr := msgInfos[1]
		from := msgInfos[2]
		to := msgInfos[3]
		txt := msgInfos[4]
		num := len(msgInfos)
		if num > 5 {
			for i := 5; i < num; i++ {
				txt += " "
				txt += msgInfos[i]
			}
		}
		// -----------------------------
		if from == this.info.Id { // ignore self msg
			return
		}
		if to != this.info.Id && to != "*" { // if the cmd target is not this node
			return
		}

		switch cmdstr {
		case cmd.NEWNODE:
			this.logger.Info(this.sprinfLog("reive an new this cmd:" + from))
			this.checkNewNode(from)
		case cmd.STATUS:
			msg := cmd.Print(this.info.Id, from, this.getNodeStatusInfo())
			this.PublishMsg(msg)
		case cmd.PRINT:
			fmt.Println(from + ":")
			fmt.Println(txt)
		case cmd.RELOAD:
			err := this.onReload(txt)
			msg := ""
			if err != nil {
				msg = cmd.Print(this.info.Id, from, err.Error())
			} else {
				msg = cmd.Print(this.info.Id, from, "reload succeed!")
			}
			this.PublishMsg(msg)
		default:
			// not shellcmd
			this.onDefaultShell(channel, msg)
		}
	} else {
		// other channel
		this.onDefaultShell(channel, msg)
	}
}

func (this *GoNode) onDefaultShell(channel string, msg string) {
	this.behavior.OnShell(channel, msg)
}

func (this *GoNode) onReload(tag string) error {
	this.logger.Info(this.sprinfLog("node reloading..."))
	return this.behavior.OnReload(tag)
}

func (this *GoNode) getNodeStatusInfo() string {
	return "status"
}
