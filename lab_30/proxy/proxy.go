package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

const proxyAddr = ":9000"
var(
	counter = 0
	appNumber = 0
	ports []string
) 

func main() {
	flag.Parse()
	ports = flag.Args()
	appNumber = len(ports)

	http.HandleFunc("/", handleProxy)
	log.Printf("Proxy was started on %v!\n", proxyAddr)
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
	url := r. + ports[counter]
	log.Printf("url = %v\n", url)
	requestMethod := r.Method
	if requestMethod == http.MethodPost {
		resp, err = http.Post("http://" + url, r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			log.Fatalln(err)
		}
	} else if requestMethod == http.MethodGet {
		resp, err = http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
	}


	textBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	log.Printf("Request %v send to http//:%v\n", string(content), ports[counter])
	log.Printf("Response %v get from http//:%v\n", string(textBytes), ports[counter])

	if counter == appNumber {
		counter = 0
	} else {
		counter++
	}
}