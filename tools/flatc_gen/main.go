package main

import (
	"fmt"
	"gonode/utils/ini"
	"gonode/utils/io"

	"log"
	"os"
	"os/exec"
)

func main() {

	os.Remove("flatc_gen.cmd")
	flog, err := os.OpenFile("flatc_gen.cmd", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer flog.Close()

	config, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		fmt.Println("配置文件缺失..[conf.ini]")
	}

	srcPath := config.Get("path", "src")
	netPath := config.Get("path", "net")
	goPath := config.Get("path", "go")

	files, err := io.ListDir(srcPath, ".idl")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 日志
	l := log.New(flog, "", os.O_APPEND)
	for _, file := range files {
		arg := "flatc -n -o " + netPath + " " + srcPath + file
		cmd := exec.Command("flatc", "-n", "-o", netPath, srcPath+file)
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(arg)
			l.Println(arg)
		}
		arg2 := "flatc -g -o " + goPath + " " + srcPath + file
		cmd2 := exec.Command("flatc", "-g", "-o", goPath, srcPath+file)
		err2 := cmd2.Start()
		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println(arg2)
			l.Println(arg2)
		}
	}
	fmt.Println("done!")

}
