package main

import (
	"bufio"

	"path/filepath"
	"strings"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts"

	"context"
	//"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/keystore"
	"golang.org/x/crypto/ssh/terminal"

	//	"github.com/ovcharovvladimir/essentiaHybrid/cmd/utils"

	//"github.com/ovcharovvladimir/essentiaHybrid/cmd/snhelper/util"
	//	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/common"

	//	"github.com/ovcharovvladimir/essentiaHybrid/internal/essapi"

	"github.com/ovcharovvladimir/essentiaHybrid/core/types"

	//"github.com/ovcharovvladimir/essentiaHybrid/crypto"

	"github.com/ovcharovvladimir/essentiaHybrid/essclient"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"github.com/ovcharovvladimir/essentiaHybrid/rpc"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/crypto"
)

func makePanel(conn string, typ ConnectionEnum, dir string, address string) *panel {

	return &panel{
		path:       conn,
		in:         bufio.NewReader(os.Stdin),
		connection: typ,
		datadir:    dir,
		address:    address,
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

	//	var msg ethereum.CallMsg
	//	var result []byte

	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("|       Wellcome to the ESSENTIA Supernode Helper           |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

	switch w.connection {
	case Ipc:
		log.Info("ESSENTIA ipc channel", "addr=", w.path, "utc", w.datadir)
	case Rpc:
		log.Info("ESSENTIA rpc channel", "addr", w.path, "utc", w.datadir)
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

	// Basics done, loop ad infinitum about what to do
	for {

		fmt.Println("What would you like to do?")
		fmt.Println(" 1. Generate Supernode public key")
		fmt.Println(" 2. Became a VOTER")
		fmt.Println(" 3. Deploy contract")
		fmt.Println(" 4. Account")
		fmt.Println(" 5. Transfer funds")
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
			fmt.Println("Enter account address:")
			input := w.read()
			account = common.HexToAddress(input)
			log.Info("Account", "adr", account)
			fmt.Println("Enter contract address:")
			input = w.read()
			contract := common.HexToAddress(input)
			log.Info("Contract", "addr", contract)

			//			contract := bind.NewBoundContract(
			//				addr,
			//				Containers.Containers[w.Container].Contracts[w.Contract].Abi,
			//				Client,
			//				Client,
			//			)
		case choice == "3":
			fmt.Println("Enter account address:")
			input := w.read()
			address := common.HexToAddress(input)
			w.address = address.Hex()
			fmt.Println("Enter passphrase:")
			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			if err != nil {
				log.Crit(err.Error())
			}
			privKey, err := KeysLoader(w, string(password))
			if err != nil {
				log.Crit(err.Error())
			}
			txOps := bind.NewKeyedTransactor(privKey.PrivateKey)
			txOps.Value = big.NewInt(0)

			// Deploy validator registration contract
			log.Info("Wait for contract to mine")
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

			fmt.Println("Enter account address:")
			input := w.read()
			address := common.HexToAddress(input)
			fmt.Println("Enter passphrase:")
			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			var acc accounts.Account
			acc.Address = address

			ks := keystore.NewKeyStore(w.datadir, keystore.StandardScryptN, keystore.StandardScryptP)
			signAcc, err := ks.Find(acc)
			if err != nil {
				log.Crit(err.Error())
			}
			var res bool
			res, err = ks.TimedUnlock(signAcc, string(password), time.Duration(6000))

			if err == nil {
				log.Info("Unlock", "result", res, "account", signAcc.Address.String())

				balance, err := client.BalanceAt(context.Background(), address, nil)
				if err == nil {

					log.Info("Balance", "amount", balance)
				} else {
					log.Info("Balance", "err", err.Error())
				}

			} else {
				log.Info("Unlock", "result", res, "err", err.Error())
			}

		case choice == "5":

			fmt.Println("Enter account (Source) address:")
			input := w.read()
			source := common.HexToAddress(input)
			fmt.Println("Enter passphrase:")
			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			fmt.Println("Enter account (Destination) address:")
			input = w.read()
			dest := common.HexToAddress(input)
			//			fmt.Println("Enter passphrase:")
			//			password1, err := terminal.ReadPassword(int(os.Stdin.Fd()))

			fmt.Println("Enter amount (in wei):")
			input = w.read()
			amount := input

			//example from:
			//https://ethereum.stackexchange.com/questions/50775/how-to-send-signed-transaction-to-ropsten-through-infura-in-golang
			//TODO: remove below lines after debugging
			//start
			dest = common.HexToAddress("0x325a01232291c820d167feb2d7a1bfe3d8401003")
			source = common.HexToAddress("0x8071eebbd56263d11f465567a45d0cf71cddeb67")
			//end
			w.address = source.String()

			privKey, err := KeysLoader(w, string(password))
			if err != nil {
				log.Info("Keys", "addr", w.address, "pK", privKey)
				log.Crit(err.Error())
			}
			privateKey := privKey.PrivateKey
			log.Info("Priv key ok")

			nonce, err := client.PendingNonceAt(context.Background(), source)
			if err != nil {
				log.Error("Fatal", "err", err)
			}

			log.Info("PendingNonce", "val", nonce)

			value := new(big.Int)

			amount = "10000000000000000000"
			value.SetString(amount, 10) // in wei (10 eth)
			gasLimit := uint64(21000)   // in units
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Error("Fatal", "err", err)
			}
			log.Info("Suggest Gas Price", "price", gasPrice, "limit", gasLimit)

			var data []byte

			tx := types.NewTransaction(nonce, dest, value, gasLimit, gasPrice, data)
			signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
			if err != nil {
				log.Error("Fatal", "err", err, "gas", gasPrice, "nonce", nonce)
			}

			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Error("Fatal", "err", err, "gas", gasPrice, "nonce", nonce)
			}

			log.Info("Transfer ", "tx sent: %s", signedTx.Hash().Hex())

		case choice == "q":
			log.Info("BYE")
			return
		}
	}
	client.Close()
	rpcClient.Close()
}

func KeysLoader(w *panel, password string) (*keystore.Key, error) {
	var fp string
	// the function that handles each file or dir
	var ff = func(pathX string, infoX os.FileInfo, errX error) error {

		// first thing to do, check error. and decide what to do about it
		if errX != nil {
			log.Info("Error", "msg", errX, "path", pathX)
			return errX
		}
		// find out if it's a dir or file, if file, print info
		if !infoX.IsDir() {
			cs := strings.Replace(w.address, "0x", "", -1)
			log.Info("File", "path", pathX, "name", infoX.Name())
			if strings.Contains(strings.ToLower(infoX.Name()), strings.ToLower(cs)) {
				fp = pathX
				return nil
			}

		}

		return nil
	}
	err := filepath.Walk(w.datadir, ff)

	if err != nil {
		log.Crit("Error when get decrypt key")
		return nil, err
	}
	if err != nil {
		log.Info("error walking the path ", w.datadir)
	}
	keyJSON, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Crit("Error when read utc file")
		return nil, err
	}
	privKey, err := keystore.DecryptKey(keyJSON, password)

	return privKey, nil

}
