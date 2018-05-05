package sysinfo

import (
	"fmt"
	"testing"
)

func TestGetSysInfo(t *testing.T) {
	s, err := GetSysInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}
