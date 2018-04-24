package main

import (
	"common/slog"
	"controller"
	"controller/models"
	"github.com/cihub/seelog"
	"thrift/gen-go/cgroupRpc"
)

func InitAll() {
	slog.InitSeelog()
	controller.InitCgroupMountPath()
	models.InitAllSubSystemRsPath()
}

func main() {
	defer seelog.Flush()
	InitAll()
	cgroupRpc.StartRpcServer()
}
