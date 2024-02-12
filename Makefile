#!/usr/bin/make -f

ksync:
	go build -mod=readonly -o ./build/celestia-kyve-rpc ./cmd/celestia-kyve-rpc/main.go