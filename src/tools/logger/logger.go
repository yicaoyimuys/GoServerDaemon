package logger

import (
	"log"
	"os"
	"strings"
	"global"
)

//系统日志
func StartLogger(path string) {
	if !strings.HasPrefix(path, "/") {
		path = global.GetCurrPath() + "../logs/" + path
	}

	// 打开日志文件
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("cannot open logfile %v\n", err)
	}

	// 创建MUX
	var r Repeater
	r.out1 = os.Stdout
	r.out2 = file
	log.SetOutput(&r)
}