package main

import (
	"bufio"
	//	"encoding/json"
	"fmt"
	//	"io/ioutil"
	"os"

	//"path/filepath"
	//	"strings"
	//"sync"

	"github.com/ovcharovvladimir/essentiaHybrid/essclient"
	//	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"github.com/ovcharovvladimir/essentiaHybrid/rpc"
)

func makePanel(conn string, typ ConnectionEnum) *panel {

	return &panel{
		path:       conn,
		in:         bufio.NewReader(os.Stdin),
		connection: typ,
	}
}

// run displays some useful infos to the user, starting on the journey of
// setting up a new or managing an existing Ethereum private network.
func (w *panel) run() {

	// Set up RPC client
	var rpcClient *rpc.Client
	var err error

	//var txOps *bind.TransactOpts

	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| " + CYAN + "Wellcome to the ESSENTIA Supernode Helper" + RESET + "                 |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()
	switch w.connection {
	case Ipc:
		fmt.Println(GREEN + "ESSENTIA ipc channel : " + WHITE + w.path + RESET)

	case Rpc:
		fmt.Println(GREEN + "ESSENTIA rpc channel : " + WHITE + w.path + RESET)
	}

	rpcClient, err = rpc.Dial(w.path)
	Error(err)
	rpcClient.Close()
	client := essclient.NewClient(rpcClient)
	if client == nil {
		Error(err)
	}

	// Basics done, loop ad infinitum about what to do
	for {
		fmt.Println()
		fmt.Println("What would you like to do?")
		fmt.Println(" 1. Generate Supernode public key")
		fmt.Println(" 2. Became a VOTER")
		fmt.Println(" q. QUIT")
		choice := w.read()
		switch {
		case choice == "1":
			fmt.Println("1")
			res, err := Key()
			log.Info("cfg", "c=", w.connection)
			if err != nil {
				log.Error("Fatal", "err", err)
			}
			fmt.Printf(GREEN+"Public key:"+RESET+" 0x%s \n", res)
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Println("NOW you may run supernode:")
			fmt.Println("> ./sness ---vrcaddr {contract_address}  --pubkey {public_key} --enable-powchain")
			fmt.Println("--------------------------------------------------------------------------------")
		case choice == "2":
			fmt.Println("2")
		case choice == "q":
			fmt.Println(CYAN + "BYE" + RESET)
			return
		}
	}

}
