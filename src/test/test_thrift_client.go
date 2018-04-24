package main

import (
	"thrift/gen-go/cgroupRpc"
	"fmt"
)

func main() {
	s, err := cgroupRpc.GetCpuAndMemStats("127.0.0.1", 10085)
	fmt.Println(s, err)
}
