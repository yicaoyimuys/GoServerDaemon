package process

import (
	"os/exec"
	"bytes"
	. "tools"
	"strings"
	"strconv"
	"global"
)

type Process struct {
	User string
	Pid int
	Cpu float64
	Mem float64
	Command string
}

func GetAllProcess() []*Process {
	processes := make([]*Process, 0)

	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		ERR(err)
		return processes
	}

	for {
		line, err := out.ReadString('\n')
		if err!=nil {
			break;
		}

		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range(tokens) {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}

		user := ft[0]

		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			continue
		}
		mem, err := strconv.ParseFloat(ft[3], 64)
		if err != nil {
			continue
		}
		command := ""
		for i := 10; i<len(ft); i++ {
			command = command + ft[i] + " "
		}
		processes = append(processes, &Process{user, pid, cpu, mem, command})
	}

	return processes
}

func GetIdleMems() float64 {
	//此命令仅在linux下使用
	if global.IsLinux() == false {
		return 1
	}

	cmd := exec.Command("/bin/sh", "-c", `free -m | grep 'Mem' | awk '{print $2" "$3" "$4" "$5" "$6" "$7}'`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 1
	}

	line, err := out.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	if err != nil || line == ""{
		return 1
	}

	tokens := strings.Split(line, " ")

	total, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 1
	}
	idle1, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return 1
	}
	idle2, err := strconv.ParseFloat(tokens[4], 64)
	if err != nil {
		return 1
	}
	idle3, err := strconv.ParseFloat(tokens[5], 64)
	if err != nil {
		return 1
	}

	return (idle1 + idle2 + idle3) / total
}

func GetIdleCpu() float64 {
	command := ""
	if global.IsLinux() {
		command = `uptime | awk '{print $10}'`
	} else if global.IsMac() {
		command = `uptime | awk '{print $8}'`
	} else {
		return 1
	}

	cmd := exec.Command("/bin/sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 1
	}

	line, err := out.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	if err != nil || line == ""{
		return 1
	}

	tokens := strings.Split(line, " ")

	use, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 1
	}

	if global.IsLinux() {
		return (1 - use) / 1
	} else if global.IsMac() {
		return (100 - use) / 100
	} else {
		return 1
	}
}
