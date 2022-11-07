package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
)

const proxyAddr = "localhost:9000"

var (
	nodeCounter = 0
	appNumber   = 0
	ports       []string
)

//go run proxy.go 8080 8081

func main() {
	flag.Parse()
	ports = flag.Args()
	appNumber = len(ports)
	if appNumber == 0 {
		panic("Param 'ports' is absent! Please specify ports for your app.")
	}
	http.HandleFunc("/", handleProxy)
	log.Printf("Proxy was started on %v with apps on ports: %v!\n", proxyAddr, ports)
	err := http.ListenAndServe(proxyAddr, nil)
	log.Fatal(err)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	var resp *http.Response
	url := "http://localhost:" + ports[nodeCounter] + r.RequestURI
	log.Printf("url = %v", url)
	requestMethod := r.Method
	if requestMethod == http.MethodPost {
		resp, err = http.Post(url, r.Header.Get("Content-Type"), bytes.NewReader(content))
		if err != nil {
			log.Fatalln(err)
		}
	} else if requestMethod == http.MethodGet {
		resp, err = http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		client := &http.Client{}

		req, err := http.NewRequest(requestMethod, url, bytes.NewReader(content))
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if resp != nil {
		textBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		log.Printf("Request %v send to %v\n", string(content), url)
		log.Printf("Response %v get from %v\n", string(textBytes), url)
	}
	nodeCounter++
	log.Printf("counter = %v, appNumber = %v", nodeCounter, appNumber)
	if nodeCounter >= appNumber {
		nodeCounter = 0
	}
}
