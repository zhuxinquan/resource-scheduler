package models

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type CGroups struct{}

//获取某个Group所有子系统的指标
func (this CGroups) ReadAllCgroupMetric(path string) (string, error) {
	subSystemsMetric := make([]SubSystemMetric, 0)
	var existGroup bool = false
	for _, subSys := range CGroupSubSystemList {
		subSystemPath := Paths{}.JoinSubSystemPath(path, subSys, CgroupMountPath)
		exist, _ := common.PathExists(subSystemPath)
		if exist {
			existGroup = true
			metric := make(map[string]string)
			var tmpSubSys SubSystemMetric
			tmpSubSys.SubSystem = subSys
			files, err := ioutil.ReadDir(subSystemPath)
			if err != nil {
				seelog.Errorf("获取目录下文件失败，err: %v", err)
			}
			for _, f := range files {
				if f.Name() == "cgroup.event_control" {
					continue
				}
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
	err = Metrics{}.SetCGroupMetric(subSystemMetrics, path)
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
	tasksPath := fmt.Sprintf("%s/memory/rs/tasks", CgroupMountPath)
	err := ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(int(pid))), os.ModeAppend)
	if err != nil {
		seelog.Errorf("通用mem子系统写入PID失败, err:%v", err)
		return err
	}
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

func (this CGroups) CreatGroup(subSystems []string, path string) error {
	for _, s := range subSystems {
		path := Paths{}.JoinSubSystemPath(path, s, CgroupMountPath)
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
	cmd := exec.Command("/bin/bash", "-c", "sleep 1 && "+command)
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

//获取所有的GroupList
func (this CGroups) ReadAllSubsystemList() ([]string, error) {
	//groupMaps := make(map[string]int)
	//groups := make([]string, 0)
	//dirs, err := ioutil.ReadDir(CgroupMountPath)
	//if err != nil {
	//	seelog.Errorf("%v", err)
	//	return groups, err
	//}
	//for _, subSys := range dirs {
	//	subSys
	//}
	return nil, nil
}
