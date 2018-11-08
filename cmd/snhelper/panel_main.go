package main

import (
	"bufio"

	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/ovcharovvladimir/essentiaHybrid"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/keystore"
	"github.com/ovcharovvladimir/essentiaHybrid/cmd/snhelper/util"
	"github.com/ovcharovvladimir/essentiaHybrid/common"

	//"github.com/ovcharovvladimir/essentiaHybrid/core/types"

	"github.com/ovcharovvladimir/essentiaHybrid/crypto"

	"github.com/ovcharovvladimir/essentiaHybrid/essclient"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"github.com/ovcharovvladimir/essentiaHybrid/rpc"
)

func makePanel(conn string, typ ConnectionEnum, utc string, pass string) *panel {

	return &panel{
		path:       conn,
		in:         bufio.NewReader(os.Stdin),
		connection: typ,
		utcfile:    utc,
		passfile:   pass,
	}
}

// run displays some useful infos to the user, starting on the journey of
// setting up a new or managing aresultn existing Ethereum private network.
func (w *panel) run() {

	// Set up RPC client
	var rpcClient *rpc.Client
	var err error
	var modules map[string]string
	var account common.Address
	var msg ethereum.CallMsg
	var result []byte

	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("|       Wellcome to the ESSENTIA Supernode Helper           |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

	switch w.connection {
	case Ipc:
		log.Info("ESSENTIA ipc channel", "addr=", w.path, "utc", w.utcfile, "pass", w.passfile)
	case Rpc:
		log.Info("ESSENTIA rpc channel", "addr", w.path, "utc", w.utcfile, "pass", w.passfile)
	}

	rpcClient, err = rpc.Dial(w.path)
	if err != nil {
		log.Crit("RPC Dial Failed", "err=", err)
	} else {

		modules, err = rpcClient.SupportedModules()
		if len(modules) == 0 {
			log.Crit("Can't dial with ESS node")
			return
		} else {
			log.Info("RPC Connected", "modules", modules)
		}

	}
	client := essclient.NewClient(rpcClient)
	var conn = false
	if client != nil {
		conn = true
	}
	var id *big.Int
	id, err = client.NetworkID(context.Background())
	if err == nil {
		log.Info("ESS Client", "connection", conn, "netId=", id)
	} else {
		log.Error("ESS Client", "err", err)
	}
	fmt.Println("Enter account address:")
	input := w.read()
	account = common.HexToAddress(input)
	log.Info("Account", "adr", account)
	// Basics done, loop ad infinitum about what to do
	for {

		fmt.Println("What would you like to do?")
		fmt.Println(" 1. Generate Supernode public key")
		fmt.Println(" 2. Became a VOTER")
		fmt.Println(" 3. Deploy contract")
		fmt.Println(" 4. Account")
		fmt.Println(" q. QUIT")
		choice := w.read()
		switch {
		case choice == "1":
			log.Info("Generating Supernode public key ..")
			res, err := KeyGen()
			if err != nil {
				log.Error("Fatal", "err", err)
			}
			log.Info("Result", "public key", res)
		case choice == "2":

			fmt.Println("Enter contract address:")
			input := w.read()
			contract := common.HexToAddress(input)
			log.Info("Contract", "addr", contract)

			msg.From = account
			msg.To = &contract
			msg.Value = util.ToWei(32, 18)

			//
			//client.SendTransaction(context.Background(), tx)
			result, err = client.CallContract(context.Background(), msg, nil)
			log.Info("res", "r", result)

		case choice == "3":
			privKey, error := KeysLoader(w)
			if error != nil {
				log.Crit(error.Error())
			}
			txOps := bind.NewKeyedTransactor(privKey.PrivateKey)
			txOps.Value = big.NewInt(0)
			// Deploy validator registration contract
			addr, tx, _, err := DeployValidatorRegistration(txOps, client)
			if err != nil {
				log.Error("Error when deploy contract")
			}

			// Wait for contract to mine
			for pending := true; pending; _, pending, err = client.TransactionByHash(context.Background(), tx.Hash()) {
				if err != nil {
					log.Error("Error when pending")
				}
				time.Sleep(1 * time.Second)
			}
			log.Info("New contract deployed", "addr", addr.Hex())

		case choice == "4":
			if w.passfile == "" || w.utcfile == "" {
				log.Error("error")
			} else {
				file, err := os.Open(w.passfile)
				if err != nil {
					log.Error("Error when open password file")
				}

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanWords)
				scanner.Scan()
				password := scanner.Text()

				keyJSON, err := ioutil.ReadFile(w.utcfile)
				if err != nil {
					log.Error("Error when read utc file")
				}
				privKey, err := keystore.DecryptKey(keyJSON, password)
				if err != nil {
					log.Error("Error when get decrypt key")
				}

				publicKey := privKey.PrivateKey.PublicKey
				address := crypto.PubkeyToAddress(publicKey).Hex()
				log.Info("Address: %v", address)

			}

		case choice == "q":
			log.Info("BYE")
			return
		}
	}
	client.Close()
	rpcClient.Close()
}

func KeysLoader(w *panel) (*keystore.Key, error) {
	if w.passfile == "" || w.utcfile == "" {
		return nil, errors.New("No key or passphrase")
	} else {
		file, err := os.Open(w.passfile)
		if err != nil {
			log.Crit("Error when open password file")
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		password := scanner.Text()

		keyJSON, err := ioutil.ReadFile(w.utcfile)
		if err != nil {
			log.Crit("Error when read utc file")
			return nil, err
		}
		privKey, err := keystore.DecryptKey(keyJSON, password)
		if err != nil {
			log.Crit("Error when get decrypt key")
			return nil, err
		}

		return privKey, nil

	}
}
