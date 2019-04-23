package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//定义日志级别
const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

//定义日志结构体
type logger struct {
	level     int
	file      *os.File
	loggerf   *log.Logger
	writeFile bool
}

// logging 接口
type logging interface {
	Init(path, logname string, level int) error
	fileWrite(time time.Time, msg interface{}) error
	SetLevel(level string)
	GetLevel() int
}

// 获取等级
func (llog *logger) GetLevel() int {
	return llog.level
}

//等级设置
func (llog *logger) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		llog.level = LevelDebug
	case "info":
		llog.level = LevelInfo
	case "warning":
		llog.level = LevelWarning
	case "error":
		llog.level = LevelError
	}
}

//时间转字符串
func formatTime(t *time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func NewLogger() *logger {
	loger := new(logger)
	loger.SetLevel("debug")
	loger.loggerf = log.New(os.Stdout, "Logger_Stdout_", log.Llongfile|log.Ltime|log.Ldate)
	return loger
}

//当前文件
func getwd() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.New(fmt.Sprintln("Getwd Error:", err))
	}
	filepathstr := filepath.Join(path, "logs")
	err = os.MkdirAll(filepathstr, 0755)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Create Dirs Error:", err))
	}
	return filepathstr, nil
}

//创建日志文件
//制定创建路径和日志名
func Init(path, logname string, level string) (*logger, error) {
	var logging *os.File
	var err error
	if path == "" {
		filePath, _ := getwd()
		filepathstr := filepath.Join(filePath, logname)
		logging, err = os.Create(filepathstr)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Create File Error:", err))
		}
	} else if logname == "" {
		filepathstr := path
		err := os.MkdirAll(filepathstr, 0644)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Create Dirs Error:", err))
		}
		logname = filepath.Join(filepathstr, "logs.log")
		logging, err = os.Create(logname)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Create File Error:", err))
		}
	} else {
		filepathstr := path
		err := os.MkdirAll(filepathstr, 0644)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Create Dirs Error:", err))
		}
		logname = filepath.Join(filepathstr, logname)
		logging, err = os.Create(logname)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("Create File Error:", err))
		}
	}
	llog := new(logger)
	llog.file = logging
	llog.SetLevel(level)
	//log.Llongfile|log.Ltime|log.Ldata 返回文件信息，时间和代码信息不需要重新写函数获取
	//初始化一个*log.Logger
	llog.loggerf = log.New(logging, "Logger_File_", log.Llongfile|log.Ltime|log.Ldate)
	return llog, nil
}

//
func (llog *logger) SetPut(boolean bool) {
	llog.writeFile = boolean
}

//set output
func (llog *logger) SetWrite() {
	if llog.writeFile == true {
		llog.loggerf.SetOutput(llog.file)
	} else {
		llog.loggerf.SetOutput(os.Stdout)
	}
}

//按照格式打印日志
func (llog *logger) Debugf(format string, v ...interface{}) {
	if llog.level > LevelDebug {
		//log里的函数都自带有mutex
		//这里获取一次锁
		return
	}
	llog.loggerf.Printf("[DEBUG] "+format, v...)
}

//Infof
func (llog *logger) Infof(format string, v ...interface{}) {
	if llog.level > LevelInfo {
		return
	}
	llog.loggerf.Printf("[INFO] "+format, v...)
}

//Warningf
func (llog *logger) Warningf(format string, v ...interface{}) {
	if llog.level > LevelWarning {
		return
	}
	llog.loggerf.Printf("[WARNING] "+format, v...)
}

//Errorf
func (llog *logger) Errorf(format string, v ...interface{}) {
	if llog.level > LevelError {
		return
	}
	llog.loggerf.Printf("[ERROR] "+format, v...)
}

//标准输出
//Println存在输出会换行
//这里采用Print输出方式
func (llog *logger) Debug(v ...interface{}) {
	if llog.level > LevelDebug {
		return
	}
	llog.loggerf.Print("[DEBUG] " + fmt.Sprintln(v...))
}

//Info
func (llog *logger) Info(v ...interface{}) {
	if llog.level > LevelInfo {
		return
	}
	llog.loggerf.Print("[INFO] " + fmt.Sprintln(v...))
}

//Warning
func (llog *logger) Warning(v ...interface{}) {
	if llog.level > LevelWarning {
		return
	}
	llog.loggerf.Print("[WARNING] " + fmt.Sprintln(v...))
}

//Error
func (llog *logger) Error(v ...interface{}) {
	if llog.level > LevelError {
		return
	}
	llog.loggerf.Print("[ERROR] " + fmt.Sprintln(v...))
}
