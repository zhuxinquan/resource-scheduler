package models

//
//import (
//	"fmt"
//	"github.com/cihub/seelog"
//	"io/ioutil"
//	"strconv"
//	"strings"
//)
//
//type Weight struct{}
//
//func (this Weight) GetAllCpuWeight() int64 {
//	var weight int64 = 0
//	cpuSubSystemPath := fmt.Sprintf("%s/cpu/rs", CgroupMountPath)
//	dirs, err := ioutil.ReadDir(cpuSubSystemPath)
//	if err != nil {
//		seelog.Errorf("读取CPU子系统目录错误[%v]", err)
//		return 0
//	}
//	for _, f := range dirs {
//		if f.IsDir() {
//			path := fmt.Sprintf("%s/%s/cpu.shares", cpuSubSystemPath, f.Name())
//			out, err := ioutil.ReadFile(path)
//			if err != nil {
//				seelog.Errorf("读取CPU share 错误[%v]", err)
//				return weight
//			}
//			tmpWeight := strings.Trim(string(out), " ")
//			tmpWeight = strings.Trim(string(out), "\n")
//			w, err := strconv.Atoi(tmpWeight)
//			if err != nil {
//				seelog.Errorf("Atoi Failed[%v]", err)
//			}
//			weight += int64(w)
//		}
//	}
//	return weight
//}
//
//func (this Weight) GetAllBlkioWeight() int64 {
//	var weight int64 = 0
//	memSubSystemPath := fmt.Sprintf("%s/blkio/rs", CgroupMountPath)
//	dirs, err := ioutil.ReadDir(memSubSystemPath)
//	if err != nil {
//		seelog.Errorf("读取CPU子系统目录错误[%v]", err)
//		return 0
//	}
//	for _, f := range dirs {
//		if f.IsDir() {
//			path := fmt.Sprintf("%s/%s/blkio.weight", memSubSystemPath, f.Name())
//			out, err := ioutil.ReadFile(path)
//			if err != nil {
//				seelog.Errorf("读取Blkio Weight 错误[%v]", err)
//				return weight
//			}
//			tmpWeight := strings.Trim(string(out), " ")
//			tmpWeight = strings.Trim(string(out), "\n")
//			w, err := strconv.Atoi(tmpWeight)
//			if err != nil {
//				seelog.Errorf("Atoi Failed[%v]", err)
//			}
//			weight += int64(w)
//		}
//	}
//	return weight
//}
