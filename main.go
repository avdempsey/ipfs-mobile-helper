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

const (
	UploadPathEnvVar           = "IPFS_MOBILE_HELPER_UPLOAD_PATH"
	UploadPathSingleNodeSuffix = "single"
	UploadPathClusterSuffix    = "cluster"
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

	uploadPath := os.Getenv(UploadPathEnvVar)
	if uploadPath == "" {
		log.Fatalf("missing env var: %s", UploadPathEnvVar)
	}

	clusterCfg := &ipfsCluster.Config{}
	cluster, err := ipfsCluster.NewDefaultClient(clusterCfg)
	if err != nil {
		log.Fatalf("could not create ipfs-cluster client: %v", err)
	}

	addHandler := NewAddHandler(ipfs, uploadPath)
	http.Handle("/add", addHandler)

	clusterAddHandler := NewClusterAddHandler(cluster, uploadPath)
	http.Handle("/cluster/add", clusterAddHandler)

	log.Fatal(http.ListenAndServe(":9999", nil))
}

type ClusterAddHandler struct {
	client     ipfsCluster.Client
	uploadPath string
}

func NewClusterAddHandler(client ipfsCluster.Client, uploadPath string) ClusterAddHandler {
	return ClusterAddHandler{
		client:     client,
		uploadPath: filepath.Join(uploadPath, UploadPathClusterSuffix),
	}
}

func (h ClusterAddHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	var addPaths []string
	err := filepath.Walk(h.uploadPath, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}
		addPaths = append(addPaths, path)
		return nil
	})
	if err != nil {
		log.Printf("error walking upload path %q: %v", h.uploadPath, err)
		w.WriteHeader(500)
		return
	}

	ctx := context.Background()
	addChan := make(chan *api.AddedOutput)
	errChan := make(chan error)
	var results []*api.AddedOutput
	go func() {
		errChan <- h.client.Add(ctx, addPaths, api.DefaultAddParams(), addChan)
	}()
	for {
		select {
		case res := <-addChan:
			if res == nil {
				continue
			}
			log.Printf("cluster add res: %+v", res)
			results = append(results, res)
		case err := <-errChan:
			if err != nil {
				log.Printf("error calling cluster Add: %v", err)
				w.WriteHeader(500)
				return
			}
			resBytes, err := json.Marshal(results)
			if err != nil {
				log.Printf("error marshaling json %v in ClusterAddHandler: %v", results, err)
				w.WriteHeader(500)
				return
			}
			n, err := fmt.Fprint(w, string(resBytes))
			if err != nil {
				log.Printf("error writing ClusterAddHandler respsonse after %d bytes written: %v", n, err)
			}
			return
		}
	}
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
	return AddHandler{
		ipfs:       ipfs,
		uploadPath: filepath.Join(uploadPath, UploadPathSingleNodeSuffix),
	}
}

func (h AddHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	var results []AddResult
	err := filepath.Walk(h.uploadPath, newFileAdder(h.ipfs, &results))
	if err != nil {
		log.Printf("error walking upload path %q: %v", h.uploadPath, err)
		w.WriteHeader(500)
		return
	}
	resBytes, err := json.Marshal(results)
	if err != nil {
		log.Printf("error marshaling json %v in AddHandler: %v", results, err)
		w.WriteHeader(500)
		return
	}
	n, err := fmt.Fprint(w, string(resBytes))
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
