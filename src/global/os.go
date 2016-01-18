package global

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"runtime"
)


//获取当前执行路径
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "/")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	return splitstring[0]
}

//当前是否是Linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

//当前是否是Mac系统
func IsMac() bool {
	return runtime.GOOS == "darwin"
}
