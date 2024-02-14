#!/usr/bin/make -f

celestia-kyve-rpc:
	go build -mod=readonly -o ./build/celestia-kyve-rpc ./cmd/celestia-kyve-rpc/main.go