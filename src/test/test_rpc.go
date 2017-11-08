package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply string
	args := 7
	err = client.Call("Hello.EchoHello", args, &reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("reply: ", reply)
}
