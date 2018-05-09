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

func (this CGroups) GroupDelete(path string) (string, error) {
	return models.CGroups{}.GroupDelete(path)
}

//获取某个Group指定子系统的指标
func (this CGroups) ReadSingleSubsytemCgroupMetric(path, subSystem string) (string, error) {
	return models.CGroups{}.ReadSingleSubsytemCgroupMetric(path, subSystem)
}

//获取所有Group列表
func (this CGroups) GetGroupList() (string, error) {
	return models.CGroups{}.GetGroupList()
}

//创建Group
func (this CGroups) NewGroup(path, subSystems string) (string, error) {
	return models.CGroups{}.NewGroup(path, subSystems)
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
