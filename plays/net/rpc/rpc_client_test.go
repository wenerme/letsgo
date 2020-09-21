package main_test

import (
	"github.com/wenerme/letsgo/rpcutil"
	"log"
	"net/rpc"
	"testing"
)

func TestMakeCallClient(t *testing.T) {
	TestServer(t)

	c := &ArithClient{}

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}

	rpcutil.MakeCallClient(client.Call, "Arith", c)

	if rel, err := c.Multiply(&Args{A: 10, B: 2}); err != nil {
		panic(err)
	} else {
		log.Printf("R1: %v", rel)
	}
	if rel, err := c.Divide(&Args{A: 10, B: 2}); err != nil {
		panic(err)
	} else {
		log.Printf("R2: %v", rel)
	}
}

type ArithClient struct {
	Multiply func(args *Args) (int, error)
	Divide   func(args *Args) (Quotient, error)
}
