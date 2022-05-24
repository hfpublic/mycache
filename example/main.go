package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hfpublic/mycache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *mycache.Group {
	return mycache.NewGroup("scores", 1024, mycache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, has := db[key]; has {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, cache *mycache.Group) {
	peers := mycache.NewHTTPPool(addr)
	peers.Set(addrs...)
	cache.RegisterPeers(peers)
	log.Println("mycache is running at", addr)
	fmt.Println(addr[strings.Index(addr, ":"):])
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, cache *mycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			v, err := cache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(v.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

// go build -o mycache-example
// ./mycache-example -port=8001 &
// ./mycache-example -port=8002 &
// ./mycache-example -port=8003 -api=1 &
func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8003, "MyCache server port")
	flag.BoolVar(&api, "api", true, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	addrs := make([]string, len(addrMap))
	addrIndex := 0
	for _, v := range addrMap {
		addrs[addrIndex] = v
		addrIndex++
	}

	cache := createGroup()
	if api {
		go startAPIServer(apiAddr, cache)
	}
	startCacheServer(addrMap[port], addrs, cache)
}
