package models

import (
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"strings"
)

type Metrics struct{}

func (this Metrics) SetCGroupMetric(subSystemMetrics []SubSystemMetric, path string) error {
	seelog.Info("开始写入子系统参数")
	for _, subSystemMetric := range subSystemMetrics {
		subSystem := subSystemMetric.SubSystem
		if subSystem == "memory" {
			cgroupPath := Paths{}.JoinSubSystemPath(path, subSystem, CgroupMountPath)
			metricPath := fmt.Sprintf("%s/%s", cgroupPath, "memory.limit_in_bytes")
			r, err := ioutil.ReadFile(metricPath)
			metricStr := strings.Trim(string(r), "\n")
			metricStr = strings.Trim(metricStr, " ")
			if metricStr == "9223372036854771712" {
				err = ioutil.WriteFile(metricPath, []byte("1000M"), 0644)
				if err != nil {
					return fmt.Errorf("写入cgroup文件失败：%s", err)
				}
			}
			metricPath = fmt.Sprintf("%s/%s", cgroupPath, "memory.memsw.limit_in_bytes")
			err = ioutil.WriteFile(metricPath, []byte("9223372036854771712"), 0644)
			if err != nil {
				return fmt.Errorf("写入cgroup文件失败：%s", err)
			}
		}
		for k, v := range subSystemMetric.Metric {
			cgroupPath := Paths{}.JoinSubSystemPath(path, subSystem, CgroupMountPath)
			metricPath := fmt.Sprintf("%s/%s", cgroupPath, k)
			err := ioutil.WriteFile(metricPath, []byte(v), 0644)
			if err != nil {
				return fmt.Errorf("写入cgroup文件失败：%s", err)
			}
		}
	}
	return nil
}
