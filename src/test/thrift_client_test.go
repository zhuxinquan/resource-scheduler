package test

import (
	"controller/models"
	"encoding/json"
	"fmt"
	"testing"
	"thrift/gen-go/cgroupRpc"
	"github.com/cihub/seelog"
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

func TestCGroup_Exec(t *testing.T) {
	req := `{"path":"/agent","subSystemMetric":[{"subSystem":"cpu","metric":{"cpu.cfs_period_us":"10000","cpu.cfs_quota_us":"2000"}}],"cmd":"/home/zxq/tmptest/cgroup/a.out","user":""}`
	s, err := cgroupRpc.Exec(req, "127.0.0.1", 10085)
	if err != nil {
		seelog.Info(err)
	}
	fmt.Println(string(s))
}

