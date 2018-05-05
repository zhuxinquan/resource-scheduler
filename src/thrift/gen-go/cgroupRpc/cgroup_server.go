package cgroupRpc

import (
	"controller"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/cihub/seelog"
)

const (
	NetworkAddr = "0.0.0.0:10085"
)

type rpcService struct{}

func (rs *rpcService) ReadAllCgroupMetric(req string) (string, error) {
	return controller.CGroups{}.ReadAllCgroupMetric(req)
}

func (rs *rpcService) ReadSingleSubsytemCgroupMetric(path, subSystem string) (string, error) {
	return controller.CGroups{}.ReadSingleSubsytemCgroupMetric(path, subSystem)
}

func (rs *rpcService) Exec(req string) (res string, err error) {
	return controller.CGroups{}.Exec(req)
}

func (rs *rpcService) SetMetric(req string) (res string, err error) {
	return controller.Metrics{}.SetMetric(req)
}

func (rs *rpcService) GetCpuAndMemStats() (res string, err error) {
	return "ReadAllCgroupMetric", nil
}

func (rs *rpcService) GetGroupList() (res string, err error) {
	return "", nil
}

var server *thrift.TSimpleServer

func StartRpcServer() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		seelog.Criticalf("thrift start failed:%v", err)
		return
	}

	handler := &rpcService{}
	processor := NewRpcServiceProcessor(handler)
	server = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	seelog.Infof("start thrift server on : %v", NetworkAddr)
	server.Serve()
}
