package cfg

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"global"
)

type ConfigInfo struct  {
	Debug   			bool    		`json:"debug"`
	Daemon_log 			string 			`json:"daemon_log"`
	LoopTime_Process 	int64			`json:"loopTime_process"`
	LoopTime_Cpu		int64				`json:"loopTime_cpu"`
	LoopTime_Mem		int64				`json:"loopTime_mem"`
	LoopTime_Log		int64				`json:"loopTime_log"`
	Servers				[]ServerInfo 	`json:"servers"`
	SendMailUser 		string   		`json:"send_mail_user"`
	SendMailPwd 		string   		`json:"send_mail_pwd"`
	SendMailHost		string			`json:"send_mail_host"`
	ReceiveMailList 	string			`json:"receive_mail_list"`
	ServerName			string			`json:"server_name"`
	LogPath				string			`json:"logs_path"`
	LogKeywords			[]KeywordInfo	`json:"logs_keywords"`
}

type ServerInfo struct {
	Flag 			string `json:"flag"`
	User 			string `json:"user"`
	RestartScript 	string `json:"restart_script"`
}

type KeywordInfo struct {
	Key 			string `json:"key"`
	RestartScript 	string `json:"restart_script"`
}

var (
	info ConfigInfo
)


func init() {
	var configPath string = global.GetCurrPath() + "../data/config.json"
	_load_config(configPath)

}

func Get() ConfigInfo {
	return info
}

func _load_config(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(path, err)
		return
	}

	err = json.Unmarshal(bytes, &info)
	if err != nil {
		log.Println(path, err)
		return
	}
}
