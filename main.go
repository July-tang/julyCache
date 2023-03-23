package main

import (
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

func main() {
	julyCache.NewGroup("scores", 2<<10, julyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := julyCache.NewHttpPool(addr)
	log.Println("julycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
