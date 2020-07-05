package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	ma "github.com/multiformats/go-multiaddr"
)

const UploadPathEnvVar = "IPFS_MOBILE_HELPER_UPLOAD_PATH"

func main() {
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		log.Fatalf("could not make NewMultiaddr: %v", err)
	}

	ipfs, err := httpapi.NewApi(addr)
	if err != nil {
		log.Fatalf("could not create ipfs client: %v", err)
	}

	uploadPath := os.Getenv(UploadPathEnvVar)
	if uploadPath == "" {
		log.Fatalf("missing env var: %s", UploadPathEnvVar)
	}

	err = filepath.Walk(uploadPath, newFileAdder(ipfs))
	if err != nil {
		log.Fatalf("error walking upload path: %v", err)
	}
}

func newFileAdder(ipfs *httpapi.HttpApi) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %w", path, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		node, err := files.NewSerialFile(path, false, info)
		if err != nil {
			return fmt.Errorf("could not create new serial file from %q: %w", path, err)
		}

		ctx := context.Background()
		res, err := ipfs.Unixfs().Add(ctx, node)
		if err != nil {
			return fmt.Errorf("error adding file %q: %w", path, err)
		}

		fmt.Println(path, res.Cid())
		return nil
	}
}
