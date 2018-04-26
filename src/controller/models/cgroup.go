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
	"os"
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
func (this CGroups) Exec(cGExecReq CGExecReq) error {
	command := cGExecReq.Cmd
	path := cGExecReq.Path
	path = JoinCommonPath(path)
	subSystemMetrics := cGExecReq.SubSystemMetric
	userName := cGExecReq.User
	subSystems := make([]string, 0)
	for _, s := range subSystemMetrics {
		subSystems = append(subSystems, s.SubSystem)
	}
	err := this.CreatGroup(subSystems, path)
	if err != nil {
		seelog.Errorf("创建Group失败, err：%v", err)
		return err
	}
	err = this.SetCGroupMetric(subSystemMetrics, path)
	if err != nil {
		seelog.Errorf("设置子系统参数失败, err：%v", err)
		return err
	}
	pid, err := this.ExecCommand(command, userName)
	if err != nil {
		seelog.Errorf("执行Command失败, err:%v", err)
		return err
	}
	err = this.WritePidToTasks(pid, path, subSystems)
	if err != nil {
		seelog.Errorf("写入Pid失败,err:%v", err)
		return err
	}
	return nil
}

func (this CGroups) WritePidToTasks(pid int64, path string, subSystems []string) error {
	for _, s := range subSystems {
		tasksPath := fmt.Sprintf("%s/%s/%s/tasks", CgroupMountPath, s, path)
		err := ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(int(pid))), os.ModeAppend)
		if err != nil {
			seelog.Errorf("写入PID失败, err:%v", err)
			return err
		}
	}
	return nil
}

func (this CGroups) SetCGroupMetric(subSystemMetrics []SubSystemMetric, path string) error {
	seelog.Info("开始写入子系统参数")
	for _, subSystemMetric := range subSystemMetrics {
		subSystem := subSystemMetric.SubSystem
		for k, v := range subSystemMetric.Metric {
			cgroupPath := this.JoinSubSystemPath(path, subSystem, CgroupMountPath)
			metricPath := fmt.Sprintf("%s/%s", cgroupPath, k)
			err := ioutil.WriteFile(metricPath, []byte(v), 0644)
			if err != nil {
				return fmt.Errorf("写入cgroup文件失败：%s", err)
			}
		}
	}
	return nil
}


func (this CGroups) CreatGroup(subSystems []string, path string) error {
	for _, s := range subSystems {
		path := this.JoinSubSystemPath(path, s, CgroupMountPath)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			seelog.Errorf("创建Group失败, path: %s, err:%s", path, err)
			return err
		}
		seelog.Infof("创建Group完成, path: %s", path)
	}
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
			seelog.Info("服务开始执行", cmd.Process.Pid)
			return int64(cmd.Process.Pid), nil
		}
		time.Sleep(1 * time.Nanosecond)
	}
}

//获取子系统中Group的路径
func (this CGroups) JoinSubSystemPath(path, subSystem, cgroupMountPath string) string {
	path = JoinCommonPath(path)
	subSystem = strings.TrimPrefix(subSystem, "/")
	subSystem = strings.TrimSuffix(subSystem, "/")
	return fmt.Sprintf("%s/%s/%s", cgroupMountPath, subSystem, path)
}

//给Group路径上添加统一路径 rs, 最终返回结果不包含 '/'
func JoinCommonPath(path string) string {
	if !strings.HasPrefix(path, "/rs/") && !strings.HasPrefix(path, "rs/") {
		path = strings.TrimPrefix(path, "/")
		path = fmt.Sprintf("%s/%s", COMMON_CGROUP_PATH, path)
		strings.TrimSuffix(path, "/")
	}
	strings.TrimPrefix(path, "/")
	return path
}
