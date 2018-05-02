package cgroupRpc

import (
	"controller/models"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"testing"
)

func TestCGroups_GetCpuAndMemStats(t *testing.T) {
	s, err := GetCpuAndMemStats("127.0.0.1", 10085)
	fmt.Println(s, err)
}

func TestCGroup_ReadAllCgroupMetric(t *testing.T) {
	s, err := ReadAllCgroupMetric("temp", "127.0.0.1", 10085)
	subSystemMetric := make([]models.SubSystemMetric, 0)
	json.Unmarshal([]byte(s), &subSystemMetric)
	fmt.Println(subSystemMetric, err)
}

func TestCGroup_Exec(t *testing.T) {
	req := `{"path":"/agent","subSystemMetric":[{"subSystem":"cpu","metric":{"cpu.cfs_period_us":"10000","cpu.cfs_quota_us":"2000"}}],"cmd":"/home/zxq/tmptest/cgroup/a.out","user":""}`
	s, err := Exec(req, "127.0.0.1", 10085)
	if err != nil {
		seelog.Info(err)
	}
	fmt.Println(string(s))
}
