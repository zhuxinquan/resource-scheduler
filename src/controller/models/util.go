package models

import (
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"strings"
)

var (
	COMMON_CGROUP_PATH = "rs"
	CgroupMountPath    = ""
)

// 返回的结果不包含/,已过滤  如:/sys/fs/cgroup
func GetCgroupMountPath() string {
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
	cgroupPath = strings.TrimSuffix(cgroupPath, "/")
	return cgroupPath
}

//初始化所有子系统，添加rs路径
func InitAllSubSystemRsPath() {
	dirs, err := ioutil.ReadDir(CgroupMountPath)
	if err != nil {
		seelog.Errorf("读取CGroupMountPath错误, err: %v", err)
		return
	}
	for _, f := range dirs {
		path := fmt.Sprintf("%s/%s/%s", CgroupMountPath, f.Name(), COMMON_CGROUP_PATH)
		_, err := os.Stat(path)
		if err != nil {
			seelog.Info("添加默认路径,", path)
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	//通一设置 memory 限制, 包括内存使用上限和oom设置
	metricPath := fmt.Sprintf("%s/memory/rs/memory.memsw.limit_in_bytes", CgroupMountPath)
	err = ioutil.WriteFile(metricPath, []byte("2G"), 0644)
	metricPath = fmt.Sprintf("%s/memory/rs/memory.oom_control", CgroupMountPath)
	err = ioutil.WriteFile(metricPath, []byte("0"), 0644)
	if err != nil {
		seelog.Info("统一设置内存限制出错，err:", err)
		panic("统一设置内存限制出错")
		return
	}
}
