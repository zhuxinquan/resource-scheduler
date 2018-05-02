package controller

import (
	"controller/models"
	"encoding/json"
	"github.com/cihub/seelog"
)

type CGroups struct{}

//获取一个Group所有的指标
func (this CGroups) ReadAllCgroupMetric(groupPath string) (string, error) {
	return models.CGroups{}.ReadAllCgroupMetric(groupPath)
}

func (this CGroups) Exec(req string) (string, error) {
	var cGExecReq models.CGExecReq
	err := json.Unmarshal([]byte(req), &cGExecReq)
	if err != nil {
		seelog.Errorf("json格式化错误, %v", err)
		return "failed", err
	}
	return "success", models.CGroups{}.Exec(cGExecReq)
}

func InitCgroupMountPath() {
	models.CgroupMountPath = models.GetCgroupMountPath()
	if models.CgroupMountPath == "" {
		seelog.Critical("CGroup未挂载，请检查")
	} else {
		seelog.Info("Cgroup 已挂载，path:", models.CgroupMountPath)
	}
	subSystemList, err := models.GetAllSubsystemList()
	if err != nil {
		seelog.Errorf("获取子系统列表失败, err:%v", err)
	}
	models.CGroupSubSystemList = subSystemList
}
