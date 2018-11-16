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
bootnode:
	build/env.sh go run build/ci.go install ./cmd/bootnode
	@echo "Done building BOOTNODE."
	@echo "Run \"$(GOBIN)/bootnode\" to launch bootnode."
esskey:
	build/env.sh go run build/ci.go install ./cmd/esskey
	@echo "Done building ESSKEY."
	@echo "Run \"$(GOBIN)/esskey\" to launch esskey"

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

gess-linux: gess-linux-386 gess-linux-amd64 gess-linux-arm #gess-linux-mips64 gess-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/linux/386/gess
	@ls -ld $(GOBIN)/linux/amd64/gess
	@ls -ld $(GOBIN)/linux/arm5/gess
	@ls -ld $(GOBIN)/linux/arm6/gess
	@ls -ld $(GOBIN)/linux/arm7/gess

gess-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/386 --targets=linux/386 -v ./cmd/gess 
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/386 --targets=linux/386 -v ./cmd/snhelper
	@echo "Linux 386 cross compilation done:"
	@sudo mv $(GOBIN)/linux/386/gess-linux-386 $(GOBIN)/linux/386/gess
	@sudo rm -fv  $(GOBIN)/linux/386/gess-linux-386
	@sudo mv $(GOBIN)/linux/386/snhelper-linux-386 $(GOBIN)/linux/386/snhelper
	@sudo rm -fv  $(GOBIN)/linux/386/snhelper-linux-386

gess-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/amd64 -targets=linux/amd64 -v ./cmd/gess 
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/amd64 -targets=linux/amd64 -v ./cmd/snhelper
	@echo "Linux amd64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/amd64/gess-linux-amd64 $(GOBIN)/linux/amd64/gess
	@sudo rm -fv  $(GOBIN)/linux/amd64/gess-linux-amd64
	@sudo mv $(GOBIN)/linux/amd64/snhelper-linux-amd64 $(GOBIN)/linux/amd64/snhelper
	@sudo rm -fv  $(GOBIN)/linux/amd64/snhelper-linux-amd64

gess-linux-arm: gess-linux-arm-5 gess-linux-arm-6 gess-linux-arm-7 gess-linux-arm64
	@echo "Linux ARM cross compilation done"

gess-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm5 --targets=linux/arm-5 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm5 --targets=linux/arm-5 -v ./cmd/snhelper
	@echo "Linux ARMv5 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm5/gess-linux-arm-5 $(GOBIN)/linux/arm5/gess
	@sudo rm -fv  $(GOBIN)/linux/arm5/gess-linux-arm-5
	@sudo mv $(GOBIN)/linux/arm5/snhelper-linux-arm-5 $(GOBIN)/linux/arm5/snhelper
	@sudo rm -fv  $(GOBIN)/linux/arm5/snhelper-linux-arm-56

gess-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm6 --targets=linux/arm-6 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm6 --targets=linux/arm-6 -v ./cmd/snhelper
	@echo "Linux ARMv6 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm6/gess-linux-arm-6 $(GOBIN)/linux/arm6/gess
	@sudo rm -fv  $(GOBIN)/linux/arm6/gess-linux-arm-6
	@sudo mv $(GOBIN)/linux/arm6/snhelper-linux-arm-6 $(GOBIN)/linux/arm6/snhelper
	@sudo rm -fv  $(GOBIN)/linux/arm6/snhelper-linux-arm-6

gess-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm7 --targets=linux/arm-7 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm7 --targets=linux/arm-7 -v ./cmd/snhelper
	@echo "Linux ARMv7 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm7/gess-linux-arm-7 $(GOBIN)/linux/arm7/gess
	@sudo rm -fv  $(GOBIN)/linux/arm7/gess-linux-arm-7
	@sudo mv $(GOBIN)/linux/arm7/snhelper-linux-arm-7 $(GOBIN)/linux/arm7/snhelper
	@sudo rm -fv  $(GOBIN)/linux/arm7/snhelper-linux-arm-7

gess-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm64 --targets=linux/arm64 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm64 --targets=linux/arm64 --ldflags '-extldflags "-static"' -v ./cmd/snhelper
	@echo "Linux ARM64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm64/gess-linux-arm64 $(GOBIN)/linux/arm64/gess
	@sudo rm -fv  $(GOBIN)/linux/arm64/gess-linux-arm64
	@sudo mv $(GOBIN)/linux/arm64/snhelper-linux-arm64 $(GOBIN)/linux/arm64/snhelper
	@sudo rm -fv  $(GOBIN)/linux/arm64/snhelper-linux-amd64

gess-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/snhelper
	@echo "Linux MIPS cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips/gess-linux-mips $(GOBIN)/linux/mips/gess
	@sudo rm -fv  $(GOBIN)/linux/mips/gess-linux-mips
	@sudo mv $(GOBIN)/linux/mips/snhelper-linux-mips $(GOBIN)/linux/mips/snhelper
	@sudo rm -fv  $(GOBIN)/linux/mips/snhelper-linux-mips


gess-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mipsle  --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mipsle  --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/snhelper
	@echo "Linux MIPSle cross compilation done:"
	@sudo mv $(GOBIN)/linux/mipsle/gess-linux-mipsle $(GOBIN)/linux/mipsle/gess
	@sudo rm -fv  $(GOBIN)/linux/mipsle/gess-linux-mipsle
	@sudo mv $(GOBIN)/linux/mipsle/snhelper-linux-mipsle $(GOBIN)/linux/mipsle/snhelper
	@sudo rm -fv  $(GOBIN)/linux/mipsle/snhelper-linux-mipsle

gess-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips64   --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips64   --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/snhelper
	@echo "Linux MIPS64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips64/gess-linux-mips64 $(GOBIN)/linux/mips64/gess
	@sudo rm -fv  $(GOBIN)/linux/mips64/gess-linux-mips64
	@sudo mv $(GOBIN)/linux/mips64/snhelper-linux-mips64 $(GOBIN)/linux/mips64/snhelper
	@sudo rm -fv  $(GOBIN)/linux/mips64/snhelper-linux-mips64

gess-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO)  --dest=$(GOBIN)/linux/mips64le  --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO)  --dest=$(GOBIN)/linux/mips64le  --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/snhelper
	@echo "Linux MIPS64le cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips64le/gess-linux-mips64le $(GOBIN)/linux/mips64le/gess
	@sudo rm -fv  $(GOBIN)/linux/mips64le/gess-linux-mips64le
	@sudo mv $(GOBIN)/linux/mips64le/snhelper-linux-mips64le $(GOBIN)/linux/mips64le/snhelper
	@sudo rm -fv  $(GOBIN)/linux/mips64le/snhelper-linux-mips64le


gess-darwin: gess-darwin-amd64 #gess-darwin-386 
	@echo "Darwin cross compilation done:"
	#@ls -ld $(GOBIN)/darwin/386/gess
	@ls -ld $(GOBIN)/darwin/amd64/gess

gess-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO)  --dest=$(GOBIN)/darwin/386 --targets=darwin/386 -v  ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO)  --dest=$(GOBIN)/darwin/386 --targets=darwin/386 -v  ./cmd/snhelper
	@echo "Darwin 386 cross compilation done:"
	@sudo mv $(GOBIN)/darwin/386/gess-darwin-10.6-386 $(GOBIN)/darwin/386/gess
	@sudo rm -fv  $(GOBIN)/darwin/386/gess-darwin-10.6-386
	@sudo mv $(GOBIN)/darwin/386/snhelper-darwin-10.6-386 $(GOBIN)/darwin/386/snhelper
	@sudo rm -fv  $(GOBIN)/darwin/386/snhelper-darwin-10.6-386

gess-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/darwin/amd64 --targets=darwin/amd64 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/darwin/amd64 --targets=darwin/amd64 -v ./cmd/snhelper
	@echo "Darwin amd64 cross compilation done:"
	@sudo mv $(GOBIN)/darwin/amd64/gess-darwin-10.6-amd64 $(GOBIN)/darwin/amd64/gess
	@sudo rm -fv  $(GOBIN)/darwin/amd64/gess-darwin-10.6-amd64
	@sudo mv $(GOBIN)/darwin/amd64/snhelper-darwin-10.6-amd64 $(GOBIN)/darwin/amd64/snhelper
	@sudo rm -fv  $(GOBIN)/darwin/amd64/snhelper-darwin-10.6-amd64

gess-windows: gess-windows-386 gess-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld  $(GOBIN)/windows/386/gess.exe
	@ls -ld  $(GOBIN)/windows/amd64/gess.exe

gess-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/386 --targets=windows/386 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/386 --targets=windows/386 -v ./cmd/snhelper
	@echo "Windows 386 cross compilation done:"
	@sudo mv $(GOBIN)/windows/386/gess-windows-4.0-386.exe $(GOBIN)/windows/386/gess.exe
	@sudo rm -fv  $(GOBIN)/windows/386/gess-windows-4.0-386.exe
	@sudo mv $(GOBIN)/windows/386/snhelper-windows-4.0-386.exe $(GOBIN)/windows/386/snhelper.exe
	@sudo rm -fv  $(GOBIN)/windows/386/snhelper-windows-4.0-386.exe

gess-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/amd64 --targets=windows/amd64 -v ./cmd/gess
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/amd64 --targets=windows/amd64 -v ./cmd/snhelper
	@echo "Windows amd64 cross compilation done:"
	@sudo mv $(GOBIN)/windows/amd64/gess-windows-4.0-amd64.exe $(GOBIN)/windows/amd64/gess.exe
	@sudo rm -fv  $(GOBIN)/windows/amd64/gess-windows-4.0-amd64.exe
	@sudo mv $(GOBIN)/windows/amd64/snhelper-windows-4.0-amd64.exe $(GOBIN)/windows/amd64/snhelper.exe
	@sudo rm -fv  $(GOBIN)/windows/amd64/snhelper-windows-4.0-amd64.exe
