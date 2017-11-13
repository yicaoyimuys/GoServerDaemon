package main

import (
	"tools/process"
	"tools/timer"
	. "tools"
	"tools/cfg"
	"tools/logger"
	"strings"
	"os/exec"
	"tools/mail"
	"runtime"
	"strconv"
	"global"
	"tools/logfile"
)

var (
	config cfg.ConfigInfo
)

func main() {
	//配置数据
	config = cfg.Get()

	//log设置
	logger.StartLogger(config.Daemon_log)

	//开启定时器
	timer.Do(config.LoopTime_Process, 0, processDaemon)
	timer.Do(config.LoopTime_Mem, 0, memDaemon)
	timer.Do(config.LoopTime_Cpu, 0, cpuDaemon)
	if config.LogPath != "" {
		logfile.InitLogFile(config.LogPath)
		timer.Do(config.LoopTime_Log, 0, logDaemon)
	}

	//系统信息
	INFO("系统：", runtime.GOOS)
	INFO("执行路径：", global.GetCurrPath())

	//保持进程
	run()
}

// 保持进程
func run() {
	temp := make(chan int32, 10)
	for {
		select {
		case <-temp:
		}
	}
}

//log监控
func logDaemon() {
	sumUpdateTexts := logfile.CheckLogFileUpdate(config.LogPath)
	if sumUpdateTexts != "" {
		script := ""
		for _, keyword := range(config.LogKeywords){
			if strings.Contains(sumUpdateTexts, keyword.Key) {
				mailTitle := "[LOG报警]" + config.ServerName
				mailContent := "触发Log：" + keyword.Key + "\nLog内容：" + sumUpdateTexts + "\n"
				if script == "" {
					script = keyword.RestartScript
					cmd := exec.Command("/bin/sh", global.GetCurrPath() + "../scripts/" + script)
					err := cmd.Run()
					if err != nil {
						mailContent += "脚本执行失败，" + err.Error()
					} else {
						mailContent += "脚本执行成功"
					}
				}
				WARN(mailContent)
				go sendMail(mailTitle, mailContent)
			}
		}
	}
}

//cpu监控
func cpuDaemon() {
	idleCpu := process.GetIdleCpu()
	mailContent := "可用CPU：" + strconv.FormatFloat(idleCpu * 100, 'f', 2, 64) + "%"
	if idleCpu < 0.1 {
		mailTitle := "[CPU报警]" + config.ServerName
		WARN(mailContent)
		go sendMail(mailTitle, mailContent)
	} else {
		INFO(mailContent)
	}
}

//内存监控
func memDaemon() {
	idleMem := process.GetIdleMems()
	mailContent := "可用内存：" + strconv.FormatFloat(idleMem * 100, 'f', 2, 64) + "%"
	if idleMem < 0.1 {
		mailTitle := "[内存报警]" + config.ServerName
		WARN(mailContent)
		go sendMail(mailTitle, mailContent)
	} else {
		INFO(mailContent)
	}
}

//进程监控
func processDaemon() {
	processes := process.GetAllProcess()
	for _, s := range(config.Servers) {
		exists := false
		for _, p := range(processes) {
			if p.User == s.User && strings.Contains(p.Command, s.Flag) {
				exists = true
				break;
			}
		}
		if exists == false {
			go restartServer(s)
			break
		}
	}
}

//重启服务器
func restartServer(server cfg.ServerInfo) {
	mailTitle := "[宕机]" + config.ServerName
	mailContent := "重启[" + server.Flag + "], "
	if server.RestartScript != "" {
		cmd := exec.Command("/bin/sh", global.GetCurrPath() + "../scripts/" + server.RestartScript)
		err := cmd.Run()
		if err != nil {
			mailContent += "失败，" + err.Error()
			ERR(mailContent)
		} else {
			mailContent += "成功"
			INFO(mailContent)

			//5秒后，检测是否还存在宕机进程
			timer.SetTimeOut(5, processDaemon)
		}
	} else {
		mailContent += "失败，不存在重启脚本"
		ERR(mailContent)
	}

	//发送宕机邮件
	go sendMail(mailTitle, mailContent)
}

//发送邮件
func sendMail(title string, content string) {
	list := strings.Split(config.ReceiveMailList, ";")
	for _, drr := range(list) {
		err := mail.SendToMail(config.SendMailUser, config.SendMailPwd, config.SendMailHost, drr, title, content, "text")
		if err != nil {
			ERR(title + "，邮件发送到" + drr + "失败：", err)
		} else {
			INFO(title + "，邮件发送到" + drr + "成功")
		}
	}

}
