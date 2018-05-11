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
	"strings"
	"syscall"
	"time"
)

type CGroups struct{}

func (this CGroups) GroupDelete(path string) (string, error) {
	seelog.Info(path)
	groupInfoStr, err := this.GetGroupList()
	if err != nil {
		seelog.Errorf("获取GroupInfo错误[%v]", err)
		return "failed", err
	}
	groupInfos := make([]GroupInfo, 0)
	err = json.Unmarshal([]byte(groupInfoStr), &groupInfos)
	if err != nil {
		return "failed", seelog.Errorf("获取GroupInfo错误Json反解失败[%v]", err)
	}
	subSystems := make([]string, 0)
	for _, tmpInfos := range groupInfos {
		seelog.Info(tmpInfos.GroupPath)
		if tmpInfos.GroupPath == path {
			subSystems = append(subSystems, tmpInfos.SubSystems...)
		}
	}
	for _, s := range subSystems {
		tmpPath := Paths{}.JoinSubSystemPath(path, s, CgroupMountPath)
		exist, _ := common.PathExists(tmpPath)
		if exist {
			err := this.DeletePathDir(tmpPath)
			if err != nil {
				return "failed", seelog.Errorf("Group[%s/%s]删除失败[%v]", s, path, err)
			}
		}
	}
	return "success", nil
}

func (this CGroups) DeletePathDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return seelog.Errorf("ReadDir 错误[%v]", err)
	}
	for _, f := range files {
		if f.IsDir() {
			childDir := fmt.Sprintf("%s/%s", path, f.Name())
			this.DeletePathDir(childDir)
		}
	}
	pidByte, err := ioutil.ReadFile(fmt.Sprintf("%s/tasks", path))
	pidStr := string(pidByte)
	pids := strings.Join(strings.Split(pidStr, "\n"), " ")
	if pids != "" {
		err = this.KillPids(pids)
		if err != nil {
			return seelog.Errorf("Kill错误[%v]", err)
		}
	}
	time.Sleep(500 * time.Microsecond)
	exist, _ := common.PathExists(path)
	if exist {
		cmd := common.NewShell(fmt.Sprintf("rmdir %s", path))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("RemoveAll Err[%v]Out[%s]", err, string(out))
		}
	}

	return nil
}

func (this CGroups) KillPids(pids string) error {
	cmdStr := fmt.Sprintf("kill -9 %s", pids)
	cmd := common.NewShell(cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return seelog.Errorf("Kill错误[%v][%s]", err, string(out))
	}
	return nil
}

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

//获取某个Group指定子系统的指标
func (this CGroups) ReadSingleSubsytemCgroupMetric(path, subSystem string) (string, error) {
	var subSystemMetric SubSystemMetric
	var existGroup bool = false
	subSystemPath := Paths{}.JoinSubSystemPath(path, subSystem, CgroupMountPath)
	exist, _ := common.PathExists(subSystemPath)
	if exist {
		existGroup = true
		metric := make(map[string]string)
		subSystemMetric.SubSystem = subSystem
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
		subSystemMetric.Metric = metric
	}
	if existGroup == false {
		return "", fmt.Errorf("Group 不存在")
	}
	s, err := json.Marshal(subSystemMetric)
	if err != nil {
		seelog.Errorf("json字符串化失败")
		return "", err
	}
	return string(s), nil
}

//获取所有Group列表
func (this CGroups) GetGroupList() (string, error) {
	groupInfos := make([]GroupInfo, 0)
	groupMap := make(map[string][]string)
	cmd := common.NewShell("lscgroup | grep :/rs")
	r, err := cmd.CombinedOutput()
	if err != nil {
		e := seelog.Errorf("获取所有Group列表失败[%s]", err)
		return "", e
	}
	out := string(r)
	out = strings.Trim(out, "\n")
	out = strings.Trim(out, " ")
	ss := strings.Split(out, "\n")
	for _, s := range ss {
		sp := strings.Split(s, ":")
		if sp[1] == "/rs" {
			continue
		}
		if value, ok := groupMap[sp[1]]; ok {
			value = append(value, strings.Split(sp[0], ",")...)
			groupMap[sp[1]] = value
		} else {
			value = make([]string, 0)
			value = append(value, strings.Split(sp[0], ",")...)
			groupMap[sp[1]] = value
		}
	}
	for k, v := range groupMap {
		k = strings.TrimPrefix(k, "/rs/")
		k = strings.Replace(k, "---", "/", -1)
		groupInfos = append(groupInfos, GroupInfo{
			GroupPath:  k,
			SubSystems: v,
		})
	}
	str, err := json.Marshal(groupInfos)
	if err != nil {
		e := seelog.Errorf("Json 转化失败[%v]", err)
		return "", e
	}
	return string(str), nil
}

//执行一个服务
func (this CGroups) Exec(cGExecReq CGExecReq) error {
	command := cGExecReq.Cmd
	path := cGExecReq.Path
	groupInfoStr, err := this.GetGroupList()
	if err != nil {
		seelog.Errorf("获取GroupInfo错误[%v]", err)
		return err
	}
	groupInfos := make([]GroupInfo, 0)
	err = json.Unmarshal([]byte(groupInfoStr), &groupInfos)
	if err != nil {
		seelog.Errorf("获取GroupInfo错误Json反解失败[%v]", err)
		return err
	}
	subSystems := make([]string, 0)
	for _, tmpInfos := range groupInfos {
		seelog.Info(tmpInfos.GroupPath)
		if tmpInfos.GroupPath == path {
			subSystems = append(subSystems, tmpInfos.SubSystems...)
		}
	}
	userName := cGExecReq.User
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
	seelog.Info("开始写入PId:", pid, "path:", path, "subSystems:", subSystems)
	//tasksPath := fmt.Sprintf("%s/memory/rs/tasks", CgroupMountPath)
	//err := ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(int(pid))), os.ModeAppend)
	//if err != nil {
	//	seelog.Errorf("通用mem子系统写入PID失败, err:%v", err)
	//	return err
	//}
	for _, s := range subSystems {
		tasksPath := fmt.Sprintf("%s/%s/rs/%s/tasks", CgroupMountPath, s, path)
		seelog.Info(tasksPath)
		err := ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(int(pid))), os.ModeAppend)
		if err != nil {
			seelog.Errorf("写入PID失败, err:%v", err)
			return err
		}
	}
	return nil
}

func (this CGroups) NewGroup(path, subSystems, weight string) (string, error) {
	path = JoinCommonPath(path)
	subs := strings.Split(subSystems, ",")
	err := this.CheckSubSystem(subs)
	if err != nil {
		return "failed", err
	}
	return this.CreatGroup(subs, path, weight)
}

func (this CGroups) CheckSubSystem(subs []string) error {
	for _, s := range subs {
		exist := false
		for _, tmp := range CGroupSubSystemList {
			if s == tmp {
				exist = true
				break
			}
		}
		if exist == false {
			return fmt.Errorf("子系统[%s]不存在", s)
		}
	}
	return nil
}

//获取所有进程信息top
func (this CGroups) GetProcessInfo() (string, error) {
	processInfos := make([]ProcessInfo, 0)
	cmd := common.NewShell("top -bcwn 1 > /tmp/topoutput")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", seelog.Errorf("top 执行失败")
	}
	out, err := ioutil.ReadFile("/tmp/topoutput")
	outStr := string(out)
	lineSlice := strings.Split(outStr, "\n")
	lineSlice = lineSlice[7:]
	for _, line := range lineSlice {
		if line == "" {
			continue
		}
		var processInfo ProcessInfo
		strSlice := strings.Split(line, " ")
		var i int64 = 1
		for k, v := range strSlice {
			if v == "" {
				continue
			}
			switch i {
			case 1:
				processInfo.Pid = v
			case 2:
				processInfo.User = v
			case 8:
				processInfo.State = v
			case 9:
				processInfo.Cpu = v
			case 10:
				processInfo.Mem = v
			case 12:
				processInfo.Cmd = strings.Join(strSlice[k:], " ")
			}
			i++
		}
		processInfos = append(processInfos, processInfo)
	}
	result, err := json.Marshal(processInfos)
	if err != nil {
		return "", seelog.Errorf("json转化失败")
	}
	return string(result), nil
}

func (this CGroups) CreatGroup(subSystems []string, path, weight string) (string, error) {
	for _, s := range subSystems {
		path := Paths{}.JoinSubSystemPath(path, s, CgroupMountPath)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			seelog.Errorf("创建Group失败, path: %s, err:%s", path, err)
			return "failed", err
		}
		if s == "cpu" || s == "blkio" {
			this.WriteWeight(s, path, weight)
		}
		seelog.Infof("创建Group完成, path: %s", path)
	}
	return "success", nil
}

func (this CGroups) WriteWeight(subSystem, path, weight string) error {
	subSystemPath := fmt.Sprintf("%s/cpu.shares", path)
	if subSystem == "blkio" {
		subSystemPath = fmt.Sprintf("%s/blkio.weight", path)
	}
	w, err := strconv.Atoi(weight)
	if err != nil {
		return seelog.Errorf("Atoi 错误[%v]", err)
	}
	weight = strconv.Itoa(w * 100)
	err = ioutil.WriteFile(subSystemPath, []byte(weight), 0644)
	if err != nil {
		return seelog.Errorf("写入权重失败[%v]", err)
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
