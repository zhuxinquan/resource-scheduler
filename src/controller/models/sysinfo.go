package models

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"strconv"
	"strings"
)

type SysInfos struct{}

func (this SysInfos) GetSysInfo() (string, error) {
	var sysInfo SysInfo
	//CPU us
	cmd := common.NewShell("top -bn 1 |grep Cpu | cut -d \",\" -f 1 | cut -d \":\" -f 2")
	r, err := cmd.CombinedOutput()
	seelog.Info("cpu us:", string(r))
	if err != nil {
		seelog.Errorf("获取CPU us错误[%v]", err)
		return "", err
	}
	us := string(r)
	us = strings.Trim(us, " ")
	us = strings.Split(us, " ")[0]
	sysInfo.CpuUserUse = us
	//CPU sy
	cmd = common.NewShell("top -bn 1 | grep Cpu | cut -d \",\" -f 2")
	r, err = cmd.CombinedOutput()
	seelog.Info("cpu sy:", string(r))
	if err != nil {
		seelog.Errorf("获取CPU sy错误[%v]", err)
		return "", err
	}
	sy := string(r)
	sy = strings.Trim(sy, " ")
	sy = strings.Split(sy, " ")[0]
	sysInfo.CpuSysUse = sy
	//MEM and Swap Total
	cmd = common.NewShell("top -bn 1 |grep Mem | cut -d \",\" -f 1 | cut -d \":\" -f 2")
	r, err = cmd.CombinedOutput()
	seelog.Info("cpu Total:", string(r))
	if err != nil {
		seelog.Errorf("获取CPU us错误[%v]", err)
		return "", err
	}
	swapMemStr := string(r)
	swapMemSplit := strings.Split(swapMemStr, "\n")
	for i, s := range swapMemSplit {
		if i == 0 {
			tmp := strings.Trim(s, " ")
			memTotal, err := strconv.Atoi(strings.Split(tmp, " ")[0])
			if err != nil {
				seelog.Errorf("获取MemTotal错误[%v]", err)
				return "", err
			}
			sysInfo.MemTotal = int64(memTotal)
		}
		if i == 1 {
			tmp := strings.Trim(s, " ")
			swapTotal, err := strconv.Atoi(strings.Split(tmp, " ")[0])
			if err != nil {
				seelog.Errorf("获取swapTotal错误[%v]", err)
				return "", err
			}
			sysInfo.SwapTotal = int64(swapTotal)
		}
	}
	//Mem and Swap Used
	cmd = common.NewShell("top -bn 1 |grep Mem | cut -d \",\" -f 3")
	r, err = cmd.CombinedOutput()
	if err != nil {
		seelog.Errorf("获取CPU us错误[%v]", err)
		return "", err
	}
	swapMemStr = string(r)
	swapMemSplit = strings.Split(swapMemStr, "\n")
	for i, s := range swapMemSplit {
		if i == 0 {
			tmp := strings.Trim(s, " ")
			memUsed, err := strconv.Atoi(strings.Split(tmp, " ")[0])
			if err != nil {
				seelog.Errorf("获取MemTotal错误[%v]", err)
				return "", err
			}
			sysInfo.MemUsed = int64(memUsed)
		}
		if i == 1 {
			tmp := strings.Trim(s, " ")
			tmp = strings.Split(tmp, ".")[0]
			tmp = strings.Split(tmp, " ")[0]
			swapUsed, err := strconv.Atoi(tmp)
			if err != nil {
				seelog.Errorf("获取MemTotal错误[%v]", err)
				return "", err
			}
			sysInfo.SwapUsed = int64(swapUsed)
		}
	}
	//Swap and Mem used total - free = used
	sysInfo.MemFree = sysInfo.MemTotal - sysInfo.MemUsed
	sysInfo.SwapFree = sysInfo.SwapTotal - sysInfo.SwapUsed
	sysInfo.MemRate = fmt.Sprintf("%.3f", float64(sysInfo.MemUsed)/float64(sysInfo.MemTotal))
	sysInfo.SwapRate = fmt.Sprintf("%.3f", float64(sysInfo.SwapUsed)/float64(sysInfo.SwapTotal))
	result, err := json.Marshal(sysInfo)
	return string(result), nil
}
