package models

import (
	"fmt"
	"strings"
	"testing"
)

//func TestCGroups_ExecCommand(t *testing.T) {
//	pid, err := CGroups{}.ExecCommand("/home/zxq/tmptest/cgroup/a.out", "")
//	fmt.Println(pid)
//	if err != nil {
//		fmt.Println("1233", err)
//	}
//	ioutil.WriteFile("/sys/fs/cgroup/cpu/rs/temp/tasks", []byte(fmt.Sprintf("%d", pid)), 0644)
//	time.Sleep(200 * time.Second)
//}

func TestCGroups_GetGroupList(t *testing.T) {
	s := "1  2   3"
	sl := strings.Split(s, " ")
	fmt.Println(len(sl))
	fmt.Println(sl)
}

func TestCGroups_GetProcessInfo(t *testing.T) {
	CgroupMountPath = GetCgroupMountPath()
	w := Weight{}.GetAllCpuWeight()
	fmt.Println(w)
}
