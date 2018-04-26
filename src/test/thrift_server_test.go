package test

import (
	"testing"
	"thrift/gen-go/cgroupRpc"
	"common/slog"
	"controller"
	"controller/models"
)

func TestCGroups_StartRpcServer(t *testing.T) {
	slog.InitSeelog("/home/zxq/develop/resource-scheduler-agent/conf/seelog.xml")
	controller.InitCgroupMountPath()
	models.InitAllSubSystemRsPath()
	cgroupRpc.StartRpcServer()
}
