package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

type Proxy struct{}

func NewProxy() *Proxy {
	return &Proxy{}
}

// ServeHTTP is the main handler for all requests.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("Received request from %s\n%s", req.RemoteAddr, string(dump))

	if req.Method != "CONNECT" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("This is a http tunnel proxy, only CONNECT method is allowed."))
		return
	}

	// Step 1
	host := req.URL.Host
	hij, ok := rw.(http.Hijacker)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("HTTP Server does not support hijacking"))
		return
	}

	client, _, err := hij.Hijack()
	if err != nil {
		return
	}

	// Step 2
	server, err := net.Dial("tcp", host)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Dial failed"))
		return
	}
	log.Printf("Dial %s", host)

	client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))

	// Step 3
	go io.Copy(server, client)
	io.Copy(client, server)
}

func main() {
	proxy := NewProxy()
	log.Printf("Listen on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", proxy))
}

/*
nc 127.0.0.1 8080
CONNECT icanhazip.com:80 HTTP/1.1
Host: icanhazip.com:80


GET http://icanhazip.com HTTP/1.1
Host: icanhazip.com


*/
