package test

import (
	"controller/models"
	"encoding/json"
	"fmt"
	"testing"
	"thrift/gen-go/cgroupRpc"
)

func TestCGroups_GetCpuAndMemStats(t *testing.T) {
	s, err := cgroupRpc.GetCpuAndMemStats("127.0.0.1", 10085)
	fmt.Println(s, err)
}

func TestCGroup_ReadAllCgroupMetric(t *testing.T) {
	s, err := cgroupRpc.ReadAllCgroupMetric("temp", "127.0.0.1", 10085)
	subSystemMetric := make([]models.SubSystemMetric, 0)
	json.Unmarshal([]byte(s), &subSystemMetric)
	fmt.Println(subSystemMetric, err)
}
