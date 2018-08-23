.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/ovcharovvladimir/essentiaHybrid/cmd/swarm
	go build -o bin/geth     github.com/ovcharovvladimir/essentiaHybrid/cmd/geth
	go build -o bin/bootnode github.com/ovcharovvladimir/essentiaHybrid/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/ovcharovvladimir/essentiaHybrid/swarm/...
