# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gess android ios gess-cross swarm evm all test clean
.PHONY: gess-linux gess-linux-386 gess-linux-amd64 gess-linux-mips64 gess-linux-mips64le
.PHONY: gess-linux-arm gess-linux-arm-5 gess-linux-arm-6 gess-linux-arm-7 gess-linux-arm64
.PHONY: gess-darwin gess-darwin-386 gess-darwin-amd64
.PHONY: gess-windows gess-windows-386 gess-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gess:
	build/env.sh go run build/ci.go install ./cmd/gess
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gess\" to launch gess."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gess.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gess.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gess-cross: gess-linux gess-darwin gess-windows gess-android gess-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gess-*

gess-linux: gess-linux-386 gess-linux-amd64 gess-linux-arm gess-linux-mips64 gess-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-*

gess-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gess
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep 386

gess-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gess
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gesslinux-* | grep amd64

gess-linux-arm: gess-linux-arm-5 gess-linux-arm-6 gess-linux-arm-7 gess-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep arm

gess-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gess
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep arm-5

gess-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gess
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep arm-6

gess-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gess
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep arm-7

gess-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gess
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep arm64

gess-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gess
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep mips

gess-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gess
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep mipsle

gess-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gess
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep mips64

gess-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gess
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gess-linux-* | grep mips64le

gess-darwin: gess-darwin-386 gess-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gess-darwin-*

gess-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gess
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gess-darwin-* | grep 386

gess-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gess
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gess-darwin-* | grep amd64

gess-windows: gess-windows-386 gess-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gess-windows-*

gess-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gess
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gess-windows-* | grep 386

gess-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gess
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gess-windows-* | grep amd64
