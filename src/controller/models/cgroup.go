package models

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
	"strconv"
	"syscall"
	"os/user"
)

type CGroups struct{}

//获取某个Group所有子系统的指标
func (this CGroups) ReadAllCgroupMetric(path, cgroupMountPath string) (string, error) {
	subSystemsMetric := make([]SubSystemMetric, 0)
	var existGroup bool = false
	dirs, err := ioutil.ReadDir(cgroupMountPath)
	if err != nil {
		seelog.Errorf("%v", err)
	}
	for _, subSys := range dirs {
		subSystemPath := this.JoinSubSystemPath(path, subSys.Name(), cgroupMountPath)
		exist, _ := common.PathExists(subSystemPath)
		if exist {
			existGroup = true
			metric := make(map[string]string)
			var tmpSubSys SubSystemMetric
			tmpSubSys.SubSystem = subSys.Name()
			files, err := ioutil.ReadDir(subSystemPath)
			if err != nil {
				seelog.Errorf("获取目录下文件失败，err: %v", err)
			}
			for _, f := range files {
				filePath := fmt.Sprintf("%s/%s", subSystemPath, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					metric[f.Name()] = fmt.Sprintf("%s", err)
				} else {
					content := string(b)
					metric[f.Name()] = content
				}
			}
			tmpSubSys.Metric = metric
			subSystemsMetric = append(subSystemsMetric, tmpSubSys)
		}
	}
	if existGroup == false {
		return "", fmt.Errorf("Group 不存在")
	}
	s, err := json.Marshal(subSystemsMetric)
	if err != nil {
		seelog.Errorf("json字符串化失败")
		return "", err
	}
	return string(s), nil

}

//执行一个服务
func (this CGroups) Exec(cGExecReq CGExecReq) (string, error) {
	command := cGExecReq.Cmd

	return "", nil
}

func (this CGroups) CreatGroup(subSystem, path string) error {
	return nil
}

func (this CGroups) ExecCommand(command, userName string) (int64, error) {
	cmd := exec.Command("/bin/bash", "-c", "sleep 1 && " + command)
	if userName != "" {
		if userInfo, err := user.Lookup(userName); err != nil {
			return -1, err
		} else if uid, err := strconv.Atoi(userInfo.Uid); err != nil {
			return -1, err
		} else {
			fmt.Println(uid)
			cmd.SysProcAttr = &syscall.SysProcAttr{}
			cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(uid)}
		}
	}
	result := []byte{}
	var err error
	go func() {
		result, err = cmd.CombinedOutput()
		if err != nil {
			seelog.Errorf("命令执行出错，err:%v", err)
		}
	}()
	for {
		if cmd.Process != nil {
			seelog.Info("程序开始执行", cmd.Process.Pid)
			return int64(cmd.Process.Pid), nil
		}
		time.Sleep(1 * time.Nanosecond)
	}
}

func (this CGroups) SetCGroupMetric(subSystemMetric SubSystemMetric, path string) error {
	subSystem := subSystemMetric.SubSystem
	for k, v := range subSystemMetric.Metric {
		cgroupPath := this.JoinSubSystemPath(path, subSystem, CgroupMountPath)
		metricPath := fmt.Sprintf("%s/%s", cgroupPath, k)
		err := ioutil.WriteFile(metricPath, []byte(v), 0644)
		if err != nil {
			return fmt.Errorf("写入cgroup文件失败：%s", err)
		}
	}
	return nil
}

//获取子系统中Group的路径
func (this CGroups) JoinSubSystemPath(path, subSystem, cgroupMountPath string) string {
	path = JoinCommonPath(path)
	subSystem = strings.TrimPrefix(subSystem, "/")
	subSystem = strings.TrimSuffix(subSystem, "/")
	return fmt.Sprintf("%s/%s/%s", cgroupMountPath, subSystem, path)
}

//给Group路径上添加通一路径 rs, 最终返回结果不包含 '/'
func JoinCommonPath(path string) string {
	if !strings.HasPrefix(path, "/rs/") && !strings.HasPrefix(path, "rs/") {
		path = strings.TrimPrefix(path, "/")
		path = fmt.Sprintf("%s/%s", COMMON_CGROUP_PATH, path)
		strings.TrimSuffix(path, "/")
	}
	strings.TrimPrefix(path, "/")
	return path
}
