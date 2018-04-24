package common

import (
	"errors"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

type Shell struct {
	cmd  string
	user string
	Cmd  *exec.Cmd
}

func NewShell(cmd string) *Shell {
	return &Shell{
		cmd: cmd,
	}
}

func (s *Shell) SetUser(username string) {
	s.user = username
}

func (s *Shell) CombinedOutput() (result []byte, err error) {
	result = []byte{}
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	cmd := exec.Command("/bin/bash", "-c", s.cmd)
	if s.user != "" {
		if userInfo, err := user.Lookup(s.user); err != nil {
			return result, err
		} else if uid, err := strconv.Atoi(userInfo.Uid); err != nil {
			return result, err
		} else {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
			cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(uid)}
		}
	}

	doneChan := make(chan int)

	go func() {
		result, err = cmd.CombinedOutput()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		if err != nil {
			err = errors.New(string(result) + " " + err.Error())
		}
	}
	return result, err
}

func (s *Shell) Run() ([]byte, error) {
	output, err := s.CombinedOutput()
	return output, err
}
