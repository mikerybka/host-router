package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	hostsFile := os.Getenv("HOSTS_FILE")
	if hostsFile == "" {
		log.Fatal("Environment variable HOSTS_FILE is not set")
	}
	b, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Println(err)
	}
	hosts := map[string]string{}
	err = json.Unmarshal(b, &hosts)
	if err != nil {
		fmt.Println(err)
	}

	handlers := map[string]http.Handler{}
	for host, target := range hosts {
		u, err := url.Parse(target)
		if err != nil {
			log.Fatalf("invalid target url: %s: %v", target, err)
		}
		handlers[host] = httputil.NewSingleHostReverseProxy(u)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h, ok := handlers[r.Host]
		if !ok {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
		fmt.Println(r.Host)
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Environment variable PORT is not set")
	}
	addr := ":" + port
	log.Fatal(http.ListenAndServe(addr, nil))
}
