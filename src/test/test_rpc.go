package main

import (
	"encoding/json"
	"fmt"
)

//func main() {
//	client, err := rpc.DialHTTP("tcp", "127.0.0.1:6666")
//	if err != nil {
//		log.Fatal("dialing:", err)
//	}
//	var reply string
//	args := 7
//	err = client.Call("Hello.EchoHello", args, &reply)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("reply: ", reply)
//}


func main(){
	m := make(map[string]string)
	//m["1"] = "1"
	//m["2"] = "2"
	//s, _ :=json.Marshal(m)
	json.Unmarshal([]byte("{\"1\":\"1\",\"2\":\"2\"}"), &m)
	fmt.Println(len(m))
}