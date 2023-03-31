package main

import (
	"flag"
	"fmt"
	"julyCache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":   "444",
	"Jack":  "555",
	"Susan": "666",
	"July":  "777",
}

func createGroup() *julyCache.Group {
	return julyCache.NewGroup("scores", 2<<10, julyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, july *julyCache.Group) {
	peers := julyCache.NewHttpPool(addr)
	peers.Set(addrs...)
	july.RegisterPeers(peers)
	log.Println("julycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, july *julyCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := july.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("font-end server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "julycache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	july := createGroup()
	if api {
		go startAPIServer(apiAddr, july)
	}
	startCacheServer(addrMap[port], addrs, july)
}
