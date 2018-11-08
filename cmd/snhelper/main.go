package main

import (
	"os"
	"path/filepath"

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
			Name:  "utc",
			Usage: "UTC file",
			Value: "~/.essentia/keystore/*",
		},
		cli.StringFlag{
			Name:  "passfile",
			Usage: "password file",
			Value: "./password.txt",
		},
	}

	app.Action = func(c *cli.Context) error {
		// Set up the logger to print everything and the random generator
		log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
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
		if c.IsSet("utc") {
			if _, err := os.Stat(c.GlobalString("utc")); os.IsNotExist(err) {
				//File does not exist
				c.Set("utc", "")
			}
		} else {
			// get first UTC file from keystore dir
			// doesnt work!!
			files, err := filepath.Glob(c.GlobalString("utc"))
			if err == nil {
				//	fmt.Println("no error")
				//fmt.Printf("%v", files)
				log.Info("UTC Keysote files", "file", files)
				for _, myfile := range files {
					str := myfile
					log.Info("UTC Keysote files", "file", str)
				}
				//			c.Set("utc", str)
			} else {
				c.Set("utc", "")
			}
		}
		if c.IsSet("rpc") {
			makePanel(c.GlobalString("rpc"), Rpc, c.GlobalString("utc"), c.GlobalString("passfile")).run()
		} else {
			c.Set("rpc", "http://localhost:8545")
			makePanel(c.GlobalString("rpc"), Rpc, c.GlobalString("utc"), c.GlobalString("passfile")).run()
		}
		if c.IsSet("ipc") {
			makePanel(c.GlobalString("ipc"), Ipc, c.GlobalString("utc"), c.GlobalString("passfile")).run()
		}

		return nil
	}

	app.Run(os.Args)

}
