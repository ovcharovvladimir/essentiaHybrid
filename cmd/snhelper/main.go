package main

import (
	"os"
	"path/filepath"

	"github.com/ovcharovvladimir/essentiaHybrid/cmd/snhelper/util"

	"github.com/mattn/go-colorable"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
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
		cli.StringFlag{
			Name:  "ws",
			Usage: "WS address",
			Value: "ws://127.0.0.1:8546",
		},
		cli.StringFlag{
			Name:  "datadir",
			Usage: "UTC file",
			Value: "~/.essentia",
		},
		cli.StringFlag{
			Name:  "pass",
			Usage: "password file",
			Value: "./password.txt",
		},
	}

	app.Action = func(c *cli.Context) error {

		// Set up the logger
		log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat(true))))

		var home string
		home = util.GetUserHomePath()

		if c.NumFlags() > 4 {
			cli.ShowAppHelpAndExit(c, 0)

		}
		if c.IsSet("passfile") {
			if _, err := os.Stat(c.GlobalString("passfile")); os.IsNotExist(err) {
				//File does not exist
				c.Set("passfile", "")
			}
		} else {
			if _, err := os.Stat(c.GlobalString("passfile")); os.IsNotExist(err) {
				//File does not exist
				c.Set("passfile", "")
			}
		}
		if !c.IsSet("ws") {
			c.Set("ws", "ws://127.0.0.1:8546")
		}
		if c.IsSet("datadir") {
			if _, err := os.Stat(c.GlobalString("datadir")); os.IsNotExist(err) {
				//File does not exist
				c.Set("datadir", "")
				log.Crit("DataDir Not exists")
			} else {

				keystorePath := filepath.Join(home, ".essentia", "keystore")
				c.Set("datadir", keystorePath)
			}
		} else {
			keystorePath := filepath.Join(home, ".essentia", "keystore")
			c.Set("datadir", keystorePath)
		}
		if c.IsSet("rpc") {
			makePanel(c.GlobalString("rpc"), c.GlobalString("ws"), Rpc, c.GlobalString("datadir"), c.GlobalString("passfile")).run()
		} else {
			c.Set("rpc", "http://localhost:8545")
			makePanel(c.GlobalString("rpc"), c.GlobalString("ws"), Rpc, c.GlobalString("datadir"), c.GlobalString("passfile")).run()
		}
		if c.IsSet("ipc") {
			makePanel(c.GlobalString("ipc"), c.GlobalString("ws"), Ipc, c.GlobalString("datadir"), c.GlobalString("passfile")).run()
		}

		return nil
	}

	app.Run(os.Args)

}
