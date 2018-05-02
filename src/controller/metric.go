package controller

import (
	"controller/models"
	"encoding/json"
	"github.com/cihub/seelog"
)

type Metrics struct{}

//对已有的Group设置参数
func (this Metrics) SetMetric(data string) (string, error) {
	var setMetricData models.SetMetricData
	err := json.Unmarshal([]byte(data), &setMetricData)
	if err != nil {
		seelog.Errorf("原Json串反解失败[%v]", err)
		return "", err
	}
	err = models.Metrics{}.SetCGroupMetric(setMetricData.SubSystemMetrics, setMetricData.Path)
	if err != nil {
		seelog.Errorf("设置Cgroup参数失败[%v]", err)
		return "", err
	}
	return "success", nil
}
