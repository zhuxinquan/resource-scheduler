package models

import (
	"testing"
	"io/ioutil"
	"fmt"
	"time"
)

func TestCGroups_ExecCommand(t *testing.T) {
	pid, err := CGroups{}.ExecCommand("/home/zxq/tmptest/cgroup/a.out", "")
	fmt.Println(pid)
	if err != nil {
		fmt.Println("1233", err)
	}
	ioutil.WriteFile("/sys/fs/cgroup/cpu/rs/temp/tasks", []byte(fmt.Sprintf("%d", pid)), 0644)
	time.Sleep(200 * time.Second)
}
