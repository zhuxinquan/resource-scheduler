package controller

import "controller/models"

type SysInfos struct{}

func (this SysInfos) GetSysInfo() (string, error) {
	return models.SysInfos{}.GetSysInfo()
}
