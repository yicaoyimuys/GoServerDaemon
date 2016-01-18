package logfile

import (
	"os"
	. "tools"
	"path/filepath"
	"bufio"
)

var (
	LogFiles map[string]uint64 = make(map[string]uint64)
)

func InitLogFile(dirPath string) {
	err := filepath.Walk(dirPath, func(filePath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		LogFiles[filePath] = GetFileLine(filePath)
		return nil
	})

	if err != nil{
		ERR("遍历log文件夹错误：", err)
	}
}

func GetFileLine(filePath string) uint64 {
	f, err := os.Open(filePath)
	if err != nil {
		ERR(filePath, err)
		return 0
	}

	defer f.Close()

	var sumLine uint64 = 0
	buf := bufio.NewReader(f)
	for {
		_, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		sumLine += 1
	}

	return sumLine
}

func CheckLogFileUpdate(dirPath string) string {
	var updateTexts string = ""
	err := filepath.Walk(dirPath, func(filePath string, fInfo os.FileInfo, err error) error {
		if fInfo == nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}

		f, err := os.Open(filePath)
		if err != nil {
			ERR(filePath, err)
			return err
		}

		defer f.Close()

		oldLine, exists := LogFiles[filePath]
		if exists == false {
			oldLine = 0
		}

		var sumLine uint64 = 0
		buf := bufio.NewReader(f)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			sumLine += 1
			if sumLine > oldLine {
				oldLine += 1
				updateTexts += line
			}
		}

		LogFiles[filePath] = sumLine

		return nil
	})

	if err != nil{
		ERR("遍历log文件夹错误：", err)
	}

	return updateTexts
}