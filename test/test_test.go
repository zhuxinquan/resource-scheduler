package test

import (
	"controller/models"
	"encoding/json"
	"fmt"
	"testing"
)

func Test_GetAllSubSystem(t *testing.T) {
	models.CgroupMountPath = models.GetCgroupMountPath()
	fmt.Println(models.CgroupMountPath)
	l, _ := models.GetAllSubsystemList()
	models.CGroupSubSystemList = l
	fmt.Println(models.CGroupSubSystemList)
}

func Test_Test(t *testing.T) {
	var setMetricData models.SetMetricData
	setMetricData.Path = "/test"
	subSystemMetrics := make([]models.SubSystemMetric, 0)
	metrics := make(map[string]string)
	metrics["cpu.cfs_quota_us"] = "2000"
	subSystemMetrics = append(subSystemMetrics, models.SubSystemMetric{
		SubSystem: "cpu",
		Metric:    metrics,
	})
	setMetricData.SubSystemMetrics = subSystemMetrics
	s, _ := json.Marshal(setMetricData)
	fmt.Println(string(s))

}
