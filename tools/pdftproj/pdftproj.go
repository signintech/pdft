package main

import (
	"flag"
	"fmt"

	"github.com/signintech/pdft/tools/pdftproj/core"
)

var flagCreate = flag.String("create", "", "create project")

func main() {
	flag.Parse()
	err := checkFlagCreate()
	if err != nil {
		echoErr(err)
		return
	}
}

func echoErr(err error) {
	//log.Panicf("%s", err.Error())
	fmt.Printf("error %s\n", err.Error())
}

func checkFlagCreate() error {

	if *flagCreate == "" {
		return nil
	}

	cSubCmd := core.CreateSubCmd{
		ProjectPath: *flagCreate,
	}

	err := cSubCmd.Create()
	if err != nil {
		return err
	}

	return nil
}
