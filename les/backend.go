// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Ethereum Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/essentiaHybrid/common/hexutil"
	"github.com/ovcharovvladimir/essentiaHybrid/consensus"
	"github.com/ovcharovvladimir/essentiaHybrid/core"
	"github.com/ovcharovvladimir/essentiaHybrid/core/bloombits"
	"github.com/ovcharovvladimir/essentiaHybrid/core/rawdb"
	"github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/ess"
	"github.com/ovcharovvladimir/essentiaHybrid/ess/downloader"
	"github.com/ovcharovvladimir/essentiaHybrid/ess/filters"
	"github.com/ovcharovvladimir/essentiaHybrid/ess/gasprice"
	"github.com/ovcharovvladimir/essentiaHybrid/event"
	"github.com/ovcharovvladimir/essentiaHybrid/internal/essapi"
	"github.com/ovcharovvladimir/essentiaHybrid/light"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"github.com/ovcharovvladimir/essentiaHybrid/node"
	"github.com/ovcharovvladimir/essentiaHybrid/p2p"
	"github.com/ovcharovvladimir/essentiaHybrid/p2p/discv5"
	"github.com/ovcharovvladimir/essentiaHybrid/params"
	rpc "github.com/ovcharovvladimir/essentiaHybrid/rpc"
)

type LightEthereum struct {
	lesCommons

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	peers      *peerSet
	txPool     *light.TxPool
	blockchain *light.LightChain
	serverPool *serverPool
	reqDist    *requestDistributor
	retriever  *retrieveManager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *essapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *ess.Config) (*LightEthereum, error) {
	chainDb, err := ess.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	less := &LightEthereum{
		lesCommons: lesCommons{
			chainDb: chainDb,
			config:  config,
		},
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		peers:          peers,
		reqDist:        newRequestDistributor(peers, quitSync),
		accountManager: ctx.AccountManager,
		engine:         ess.CreateConsensusEngine(ctx, chainConfig, &config.Ethash, nil, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   ess.NewBloomIndexer(chainDb, light.BloomTrieFrequency, light.HelperTrieConfirmations),
	}

	less.relay = NewLesTxRelay(peers, less.reqDist)
	less.serverPool = newServerPool(chainDb, quitSync, &less.wg)
	less.retriever = newRetrieveManager(peers, less.reqDist, less.serverPool)

	less.odr = NewLesOdr(chainDb, less.retriever)
	less.chtIndexer = light.NewChtIndexer(chainDb, true, less.odr)
	less.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, true, less.odr)
	less.odr.SetIndexers(less.chtIndexer, less.bloomTrieIndexer, less.bloomIndexer)

	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if less.blockchain, err = light.NewLightChain(less.odr, less.chainConfig, less.engine); err != nil {
		return nil, err
	}
	// Note: AddChildIndexer starts the update process for the child
	less.bloomIndexer.AddChildIndexer(less.bloomTrieIndexer)
	less.chtIndexer.Start(less.blockchain)
	less.bloomIndexer.Start(less.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		less.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	less.txPool = light.NewTxPool(less.chainConfig, less.blockchain, less.relay)
	if less.protocolManager, err = NewProtocolManager(less.chainConfig, true, config.NetworkId, less.eventMux, less.engine, less.peers, less.blockchain, nil, chainDb, less.odr, less.relay, less.serverPool, quitSync, &less.wg); err != nil {
		return nil, err
	}
	less.ApiBackend = &LesApiBackend{less, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	less.ApiBackend.gpo = gasprice.NewOracle(less.ApiBackend, gpoParams)
	return less, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Etherbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightEthereum) APIs() []rpc.API {
	return append(essapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "ess",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightEthereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightEthereum) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightEthereum) TxPool() *light.TxPool              { return s.txPool }
func (s *LightEthereum) Engine() consensus.Engine           { return s.engine }
func (s *LightEthereum) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightEthereum) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightEthereum) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightEthereum) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *LightEthereum) Start(srvr *p2p.Server) error {
	s.startBloomHandlers()
	log.Warn("Light client mode is an experimental feature")
	s.netRPCService = essapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *LightEthereum) Stop() error {
	s.odr.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.engine.Close()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
