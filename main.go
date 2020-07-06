package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/ipfs-cluster/api"
	ipfsCluster "github.com/ipfs/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"
)

const UploadPathEnvVar = "IPFS_MOBILE_HELPER_UPLOAD_PATH"

func main() {
	clusterCfg := &ipfsCluster.Config{}
	cluster, err := ipfsCluster.NewDefaultClient(clusterCfg)
	if err != nil {
		log.Fatalf("could not create ipfs-cluster client: %v", err)
	}

	ctx := context.Background()
	addChan := make(chan *api.AddedOutput)
	testPaths := []string{"cluster_test.txt"}
	go cluster.Add(ctx, testPaths, api.DefaultAddParams(), addChan)
	select {
	case res := <-addChan:
		log.Printf("cluster add res: %+v", res)
	}

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

	addHandler := NewAddHandler(ipfs, uploadPath)
	http.Handle("/add", addHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

type AddResult struct {
	Path string `json:"path"`
	CID  string `json:"cid"`
	Err  error  `json:"err,omitempty"`
}

type AddHandler struct {
	ipfs       *httpapi.HttpApi
	uploadPath string
}

func NewAddHandler(ipfs *httpapi.HttpApi, uploadPath string) AddHandler {
	return AddHandler{ipfs, uploadPath}
}

func (h AddHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	var results []AddResult
	err := filepath.Walk(h.uploadPath, newFileAdder(h.ipfs, &results))
	if err != nil {
		log.Printf("error walking upload path: %v", err)
	}
	resBytes, err := json.Marshal(results)
	if err != nil {
		log.Printf("error marshaling json %v in AddHandler: %v", results, err)
	}
	n, err := fmt.Fprintf(w, "%s", resBytes)
	if err != nil {
		log.Printf("error writing AddHandler response after %d bytes written: %v", n, err)
	}
}

func newFileAdder(ipfs *httpapi.HttpApi, results *[]AddResult) func(path string, info os.FileInfo, err error) error {
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

		*results = append(*results, AddResult{path, res.Cid().String(), nil})
		return nil
	}
}
