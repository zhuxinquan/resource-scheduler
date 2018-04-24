package cgroupRpc

import (
	"log"
	"net"
	"fmt"
	"os"
	"time"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func GetCpuAndMemStats(host string, port int) (res string, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(net.JoinHostPort(host, fmt.Sprintf("%d", port)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		return res, err
	}

	useTransport, erru := transportFactory.GetTransport(transport)
	if erru != nil {
		fmt.Fprintln(os.Stderr, "error get transport: ", erru)
		return res, err
	}
	client := NewRpcServiceClientFactory(useTransport, protocolFactory)

	for i := 0; i < 3; i++ {
		err = transport.Open()
		if err == nil {
			break
		} else if i == 2 {
			fmt.Fprintln(os.Stderr, "Error opening socket to "+host+":"+fmt.Sprintf("%d", port), " ", err)
			return res, err
		}
		time.Sleep(1 * time.Second)
	}
	defer transport.Close()

	res, err = client.GetCpuAndMemStats()
	if err != nil {
		return res, err
	}
	return res, nil
}
