package controller

import (
	"testing"
	"fmt"
)

func TestCGroups_GetCgroupMountPath(t *testing.T) {
	path := CGroups{}.GetCgroupMountPath()
	fmt.Println(path)
}
