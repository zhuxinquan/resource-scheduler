package controller

import (
	"controller/models"
	"encoding/json"
	"github.com/cihub/seelog"
)

type CGroups struct{}

//获取一个Group所有的指标
func (this CGroups) ReadAllCgroupMetric(groupPath string) (string, error) {
	return models.CGroups{}.ReadAllCgroupMetric(groupPath, models.CgroupMountPath)
}

func (this CGroups) Exec(req string) (string, error) {
	var cGExecReq models.CGExecReq
	err := json.Unmarshal([]byte(req), &cGExecReq)
	if err != nil {
		seelog.Errorf("json格式化错误, %v", err)
		return "", err
	}
	return models.CGroups{}.Exec(cGExecReq)
}

func InitCgroupMountPath() {
	models.CgroupMountPath = models.GetCgroupMountPath()
	if models.CgroupMountPath == "" {
		seelog.Critical("CGroup未挂载，请检查")
	} else {
		seelog.Info("Cgroup 已挂载，path:", models.CgroupMountPath)
	}
}
