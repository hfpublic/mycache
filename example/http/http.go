package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hfpublic/mycache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// curl localhost:9999/_mycache/example/Tom
// curl localhost:9999/_mycache/example/Ang
func main() {
	mycache.NewGroup("example", 2048, mycache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, has := db[key]; has {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not found", key)
	}))

	addr := "localhost:9999"
	pool := mycache.NewHTTPPool(addr)
	log.Fatalln(http.ListenAndServe(addr, pool))
}
