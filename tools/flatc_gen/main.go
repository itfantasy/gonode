package main

import (
	"fmt"

	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"

	"os/exec"
)

func main() {

	config, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
	}

	srcPath := config.Get("path", "src")
	csPath := config.Get("path", "cs")
	goPath := config.Get("path", "go")

	files, err := io.ListDir(srcPath, ".idl")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		arg := "flatc -n -o " + csPath + " " + srcPath + file
		cmd := exec.Command("bin/flatbuffers/flatc", "-n", "-o", csPath, srcPath+file)
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(arg)
		}
		arg2 := "flatc -g -o " + goPath + " " + srcPath + file
		cmd2 := exec.Command("bin/flatbuffers/flatc", "-g", "-o", goPath, srcPath+file)
		err2 := cmd2.Start()
		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println(arg2)
		}
	}
	fmt.Println("done!")

}
