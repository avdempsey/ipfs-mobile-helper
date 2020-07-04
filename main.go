package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		log.Fatalf("could not make NewMultiaddr: %v", err)
	}

	ipfs, err := httpapi.NewApi(addr)
	if err != nil {
		log.Fatalf("could not create ipfs client: %v", err)
	}

	filepath := os.Getenv("IPFS_MOBILE_HELPER_TEST_FILE")
	if filepath == "" {
		log.Fatalf("missing env var: IPFS_MOBILE_HELPER_TEST_FILE")
	}

	stat, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("could not stat file: %q because of err: %v", filepath, err)
	}

	node, err := files.NewSerialFile(filepath, false, stat)
	if err != nil {
		log.Fatalf("could not create new serial file: %v", err)
	}

	ctx := context.Background()
	res, err := ipfs.Unixfs().Add(ctx, node)
	if err != nil {
		log.Fatalf("error adding file: %v", err)
	}
	fmt.Println(res.Cid())
}
