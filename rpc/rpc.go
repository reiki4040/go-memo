package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type RPCArgs struct {
	A, B int
}

type RPCFuncType int

func (t *RPCFuncType) Multiply(args *RPCArgs, reply *int) error {
	*reply = args.A * args.B

	// sleep for Acync call checking.
	time.Sleep(1 * time.Second)
	return nil
}

func startServer() {
	// RPC function
	rpcFunc := new(RPCFuncType)
	rpc.Register(rpcFunc)

	// set Default HTTP server
	rpc.HandleHTTP()

	// init listener
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	// run Default HTTP server
	http.Serve(l, nil)
}

func main() {
	// launch RPC server
	go startServer()

	// init RPC Client
	client, err := rpc.DialHTTP("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("failed dial:", err)
	}

	//Sync call
	args := &RPCArgs{7, 8}
	var reply int
	err = client.Call("RPCFuncType.Multiply", args, &reply)
	if err != nil {
		log.Fatal("failed call multiply:", err)
	}

	fmt.Printf("RPCFuncType: %d*%d=%d\n", args.A, args.B, reply)

	// Async call
	args2 := &RPCArgs{3, 6}
	asyncCall := client.Go("RPCFuncType.Multiply", args2, &reply, nil)

	// Block
	// if does not wait Done, reply is still 56. not 18.
	_ = <-asyncCall.Done

	if asyncCall.Error == nil {
		if a, ok := asyncCall.Args.(*RPCArgs); ok && a != nil {
			if r, ok := asyncCall.Reply.(*int); ok && r != nil {
				fmt.Printf("RPCFuncType: %d*%d=%d\n", a.A, a.B, *r)
			}
		}
	} else {
		log.Fatal("Async failed:", asyncCall.Error)
	}
}
