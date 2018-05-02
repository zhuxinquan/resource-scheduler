package models

import (
	"fmt"
	"github.com/cihub/seelog"
	"strings"
)

type Paths struct{}

//获取子系统中Group的路径
func (this Paths) JoinSubSystemPath(path, subSystem, cgroupMountPath string) string {
	path = JoinCommonPath(path)
	subSystem = strings.TrimPrefix(subSystem, "/")
	subSystem = strings.TrimSuffix(subSystem, "/")
	seelog.Infof("path: %s, subSystem:%s, cgroupMountPath %s", path, subSystem, cgroupMountPath)
	return fmt.Sprintf("%s/%s/%s", cgroupMountPath, subSystem, path)
}

//给Group路径上添加统一路径 rs, 最终返回结果不包含 '/'
func JoinCommonPath(path string) string {
	if !strings.HasPrefix(path, "/rs/") && !strings.HasPrefix(path, "rs/") {
		path = strings.TrimPrefix(path, "/")
		path = fmt.Sprintf("%s/%s", COMMON_CGROUP_PATH, path)
		strings.TrimSuffix(path, "/")
	}
	path = strings.TrimPrefix(path, "/")
	return path
}
