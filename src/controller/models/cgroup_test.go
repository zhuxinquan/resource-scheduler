package models

import (
	"testing"
	"time"
)

//func TestCGroups_ExecCommand(t *testing.T) {
//	pid, err := CGroups{}.ExecCommand("/home/zxq/tmptest/cgroup/a.out", "")
//	fmt.Println(pid)
//	if err != nil {
//		fmt.Println("1233", err)
//	}
//	ioutil.WriteFile("/sys/fs/cgroup/cpu/rs/temp/tasks", []byte(fmt.Sprintf("%d", pid)), 0644)
//	time.Sleep(200 * time.Second)
//}

func TestCGroups_Exec(t *testing.T) {
	CgroupMountPath = "/sys/fs/cgroup"
	var cGExec CGExecReq
	cGExec.User = ""
	cGExec.Path = "/agent"
	cGExec.Cmd = "/home/zxq/tmptest/cgroup/a.out"
	subSystemMetrics := make([]SubSystemMetric, 0)
	metrics := make(map[string]string)
	metrics["cpu.cfs_quota_us"] = "2000"
	metrics["cpu.cfs_period_us"] = "10000"
	subSystemMetrics = append(subSystemMetrics, SubSystemMetric{
		SubSystem: "cpu",
		Metric: metrics,
	})
	cGExec.SubSystemMetric = subSystemMetrics
	CGroups{}.Exec(cGExec)
	time.Sleep(1000* time.Second)
}