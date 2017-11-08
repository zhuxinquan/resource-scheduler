package rpc_controller

import (
	"log"
	"net/http"
	"net/rpc"
	"rpc_controller/register"
)

type Rpc struct{}

func (this Rpc) RpcRegister() {
	rpc.Register(new(register.Hello))
	rpc.HandleHTTP()
	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		log.Fatal("listen error:", err)
	}
}
