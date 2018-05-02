package models

import (
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
)

type Metrics struct{}

func (this Metrics) SetCGroupMetric(subSystemMetrics []SubSystemMetric, path string) error {
	seelog.Info("开始写入子系统参数")
	for _, subSystemMetric := range subSystemMetrics {
		subSystem := subSystemMetric.SubSystem
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
