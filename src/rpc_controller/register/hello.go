package register

import (
	"encoding/json"
	"fmt"
)

type Hello struct{}

func (this *Hello) EchoHello(args *int64, result *string) error {
	responseInfo := &OperateResponseInfo{
		Status: "success",
		Data:   "Hello resource-scheduler",
	}
	resultByte, _ := json.Marshal(responseInfo)
	*result = string(resultByte)
	fmt.Println("return a string success")
	return nil
}
