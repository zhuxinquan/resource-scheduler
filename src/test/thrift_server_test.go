package test

import (
	"testing"
	"thrift/gen-go/cgroupRpc"
)

func TestCGroups_StartRpcServer(t *testing.T) {
	cgroupRpc.StartRpcServer()
}
