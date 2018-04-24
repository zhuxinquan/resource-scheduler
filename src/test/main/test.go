package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dirs, _ := ioutil.ReadDir("/sys/fs/cgroup")
	for _, f := range dirs {
		subSysPath := f.Name()
		fmt.Println(subSysPath)
	}
}
