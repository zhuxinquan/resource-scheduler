package cgroupRpc

import (
	"common/slog"
	"controller"
	"controller/models"
	"testing"
	"thrift/gen-go/cgroupRpc"
)

func TestCGroups_StartRpcServer(t *testing.T) {
	slog.InitSeelog("/home/zxq/develop/resource-scheduler-agent/conf/seelog.xml")
	controller.InitCgroupMountPath()
	models.InitAllSubSystemRsPath()
	cgroupRpc.StartRpcServer()
}
