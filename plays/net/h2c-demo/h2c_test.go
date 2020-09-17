package main_test

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h2s := &http2.Server{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
	})

	server := &http.Server{
		Addr:    "0.0.0.0:1010",
		Handler: h2c.NewHandler(handler, h2s),
	}
	fmt.Printf("Listening [0.0.0.0:1010]...\n")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func TestHttp2OnlyServer(t *testing.T) {
	server := http2.Server{}
	l, err := net.Listen("tcp", "0.0.0.0:1010")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening [0.0.0.0:1010]...\n")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		server.ServeConn(conn, &http2.ServeConnOpts{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
			}),
		})
	}
}
func TestH2cClient(t *testing.T) {
	client := http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP: true,
			// Pretend we are dialing a TLS endpoint.
			// Note, we ignore the passed tls.Config
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	resp, _ := client.Get("http://localhost:1010")
	fmt.Printf("Client Proto: %d\n", resp.ProtoMajor)
}

/*
# Upgrade
curl -v --http2 http://localhost:1010
# No upgrade
curl -v --http2-prior-knowledge http://localhost:1010
*/
