package cgroupRpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	NetworkAddr = "0.0.0.0:10085"
)

type rpcService struct{}

func (rs * rpcService) ReadAllCgroupMetric(req string) (res string, err error) {
	return "ReadAllCgroupMetric", nil
}

func (rs * rpcService) Exec(req string) (res string, err error) {
	return "ReadAllCgroupMetric", nil
}

func (rs * rpcService) SetMetric(req string) (res string, err error) {
	return "ReadAllCgroupMetric", nil
}

func (rs * rpcService) GetCpuAndMemStats() (res string, err error) {
	return "ReadAllCgroupMetric", nil
}

var server *thrift.TSimpleServer

func StartRpcServer() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		//seelog.Criticalf("thrift start failed:%v", err)
		return
	}

	handler := &rpcService{}
	processor := NewRpcServiceProcessor(handler)
	server = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	//seelog.Infof("start thrift server on : %v", NetworkAddr)
	server.Serve()
}
