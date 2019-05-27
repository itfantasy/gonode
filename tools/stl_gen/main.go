package main

import (
	"fmt"
	"strings"

	"github.com/itfantasy/gonode/utils/args"
	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/strs"
)

func main() {
	conf := args.Parser().
		AddArg("t", "", "you can use a type to create a list, or use a pair values of key-value to create a dictionary.")
	tmp, b := conf.Get("t")
	if !b || tmp == "" {
		fmt.Println("the type setting (-t) is necessary, and it could not be an empty string!")
		return
	}
	infos := strings.Split(tmp, "-")
	l := len(infos)
	if l != 1 && l != 2 {
		fmt.Println("the type setting (-t) must be a single type of a pair types of k-v")
		return
	}
	if l == 1 { // list
		txt, err := io.LoadFile(io.CurDir() + "tmp/list.tmp")
		if err != nil {
			fmt.Println("can not find the tmp file of list!" + err.Error())
			return
		}
		t := infos[0]
		txt = strings.Replace(txt, "List<T>", "List"+strs.UcFirst(t), -1)
		txt = strings.Replace(txt, "<T>", t, -1)
		fileName := "list" + strs.UcFirst(t) + ".go"
		if err := io.SaveFile(io.CurDir()+"dst/"+fileName, txt); err != nil {
			fmt.Println(fileName + "created failed! " + err.Error())
			return
		}
		fmt.Println(fileName + "has been created!")
	} else { // dict
		txt, err := io.LoadFile(io.CurDir() + "tmp/dict.tmp")
		if err != nil {
			fmt.Println("can not find the tmp file of dict!" + err.Error())
			return
		}
		k := infos[0]
		v := infos[1]
		txt = strings.Replace(txt, "<K>", k, -1)
		txt = strings.Replace(txt, "<V>", v, -1)
		txt = strings.Replace(txt, "<K,V>", strs.UcFirst(k)+strs.UcFirst(v), -1)
		fileName := "dict" + strs.UcFirst(k) + strs.UcFirst(v) + ".go"
		if err := io.SaveFile(io.CurDir()+"dst/"+fileName, txt); err != nil {
			fmt.Println(fileName + "created failed! " + err.Error())
			return
		}
		fmt.Println(fileName + "has been created!")
	}
}
