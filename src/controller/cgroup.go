package controller

import (
	"io/ioutil"
	"strings"
	"fmt"
	"os"

)

var (
	COMMON_CGROUP_PATH = "rs"
)

type CGroups struct{}

//func (this CGroups) ReadCgroupMetric(path string) (models.CgroupMetric, error) {
//	path = strings.TrimPrefix(path, "/")
//	var cgroupMetric models.CgroupMetric
//	cgroupMountPath := this.GetCgroupMountPath()
//	if cgroupMountPath == "" {
//		return cgroupMetric, fmt.Errorf("CGroup未挂载")
//	}
//	// 遍历各个子系统
//	subSystems, err := ioutil.ReadDir(cgroupMountPath)
//	for _, sub := range subSystems {
//		// 组合Path
//		cgroupPath := fmt.Sprintf("%s/%s/%s/%s", cgroupMountPath, sub, COMMON_CGROUP_PATH, path)
//		files, err := ioutil.ReadDir(cgroupPath)
//		if err != nil {
//			//info := fmt.Sprintf("获取目录下[%s]的文件失败. err:%s", agentCgroupPath, err)
//		}
//	}
//
//	files, err := ioutil.ReadDir(cgroupPath)
//	if err != nil {
//		//info := fmt.Sprintf("获取目录下[%s]的文件失败. err:%s", agentCgroupPath, err)
//	}
//
//	ret := make(map[string]string)
//	ret["metricPath"] = agentCgroupPath
//	for _, f := range files {
//		fileName := fmt.Sprintf("%s/%s", agentCgroupPath, f.Name())
//		b, err := ioutil.ReadFile(fileName)
//		if err != nil {
//			ret[f.Name()] = fmt.Sprintf("%s", err)
//		} else {
//			content := string(b)
//			ret[f.Name()] = content
//		}
//	}
//
//	return ret
//}

// 检查一个子系统下是否存在某个path
func (this CGroups) CheckCgroupMount(path string) (bool, error) {
	cgroupPath := this.GetCgroupMountPath()
	if cgroupPath == "" {
		return false, fmt.Errorf("Cgroup未挂载")
	}
	metricPath := fmt.Sprintf("%s/%s", cgroupPath, path)
	_, err := os.Stat(metricPath)
	if err != nil {
		return false, fmt.Errorf("[%s]路径不存在", metricPath)
	}
	return true, nil
}

// 返回的结果包含/  如:/sys/fs/cgroup/
func (this CGroups) GetCgroupMountPath() string {
	b, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return ""
	}
	content := string(b)
	lines := strings.Split(content, "\n")
	cgroupPath := ""
	for _, line := range lines {
		if strings.Contains(line, "/cgroup/memory") {
			parts := strings.Split(line, " ")
			memCgroupPath := parts[1]
			cgroupPath = strings.Replace(memCgroupPath, "memory", "", -1)
			break
		}
	}
	return cgroupPath
}
