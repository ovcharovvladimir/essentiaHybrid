package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

type ConnectionEnum int

const (
	Ipc ConnectionEnum = 1 << iota
	Rpc
)

func main() {
	app := cli.NewApp()
	app.Name = "shelper"
	app.Version = "0.1.0"
	app.Usage = "Supernode initialization helper"
	app.ArgsUsage = "[arguments...]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rpc",
			Usage: "Masternode rpc path",
		},
		cli.StringFlag{
			Name:  "ipc",
			Usage: "Masternode ipc path",
		},
	}

	app.Action = func(c *cli.Context) error {

		if c.NumFlags() > 1 {
			cli.ShowAppHelpAndExit(c, 0)
		}
		if c.IsSet("rpc") {
			makePanel(c.GlobalString("rpc"), Rpc).run()
		} else {
			c.Set("rpc", "http://localhost:8545")
			makePanel(c.GlobalString("rpc"), Rpc).run()
		}
		if c.IsSet("ipc") {
			makePanel(c.GlobalString("ipc"), Ipc).run()
		}

		return nil
	}

	app.Run(os.Args)

}
