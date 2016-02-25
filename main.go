package main

import (
	"fmt"
	"os"

	"github.com/prime/cmd"
	"github.com/codegangsta/cli"
	"github.com/wrench/setting"
)

func main() {

	if err := setting.SetConfig("conf/containerops.conf"); err != nil {
		fmt.Printf("Read config error: %s", err.Error())
	}

	app := cli.NewApp()

	app.Name = setting.AppName
	app.Usage = setting.Usage
	app.Version = setting.Version
	app.Author = setting.Author
	app.Email = setting.Email

	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
