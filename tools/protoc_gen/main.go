package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
)

func main() {

	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
	}

	srcPath := conf.Get("path", "src")
	csPath := conf.Get("path", "cs")
	goPath := conf.Get("path", "go")

	files, err := io.ListDir(srcPath, ".proto")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		args := "bin\\protobuf-net\\ProtoGen\\protogen" + " -i:" + srcPath + file + " -o:" + csPath + file + ".cs" + " -ns:proto"
		dstFile := strings.Replace(file, "proto", "pb", -1)
		dstFile += ".cs"
		cmd := exec.Command("bin/protobuf-net/ProtoGen/protogen", "-i:"+srcPath+file, "-o:"+csPath+dstFile, "-ns:proto")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(args)
		}

		args2 := "bin\\protoc-gen-go\\protoc" + " --proto_path=" + srcPath + " --go_out=" + goPath + " " + file
		cmd2 := exec.Command("bin/protoc-gen-go/protoc", "--proto_path="+srcPath, "--go_out="+goPath, file)
		err2 := cmd2.Start()
		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println(args2)
		}
	}
}
