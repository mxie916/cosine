package cosine

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type LEVEL int32
type UNIT int64

// 日志文件结构体
type _FILE struct {
	dir     string
	name    string
	logfile *os.File
	mu      *sync.RWMutex
	lg      *log.Logger
}

// 日期格式
const DATEFORMAT = "2006-01-02"

// 日志级别
const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

// 文件大小但闻
const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

type Logger struct {
	logLevel        LEVEL
	consoleFlag     bool
	rollingFileFlag bool
	dailyFileFlag   bool
	maxFileSize     int64
	currentFileDate *time.Time
	logObj          *_FILE
}

// 设置日志级别
func (self *Logger) SetLevel(level LEVEL) {
	self.logLevel = level
}

// 设置是否在控制台打印
func (self *Logger) SetConsole(isConsole bool) {
	self.consoleFlag = isConsole
}

// 设置是否采用文件打印（按文件大小滚动）
func (self *Logger) SetRollingFile(fileDir, fileName string, maxSize int64, _unit UNIT) {
	self.rollingFileFlag = true

	// 设置文件最大值
	self.maxFileSize = maxSize * int64(_unit)

	// 创建日志目录&日志文件对象
	self.mklogdir(fileDir)
	if !self.isExist(fileDir + "/" + fileName) {
		os.Create(fileDir + "/" + fileName)
	}
	self.logObj = &_FILE{dir: fileDir, name: fileName, mu: new(sync.RWMutex)}

	// 锁定文件操作
	self.logObj.mu.Lock()
	defer self.logObj.mu.Unlock()

	self.handleFile4Size()

	// 1s检测一次文件是否需要滚动
	go func(logger *Logger) {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				logger.handleFile4Size()
			}
		}
	}(self)
}

// 设置是否采用文件打印（按文件日期滚动）
func (self *Logger) SetDailyFile(fileDir, fileName string) {
	self.dailyFileFlag = true

	// 设置当前日期
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	self.currentFileDate = &t

	// 创建日志目录&日志文件对象
	self.mklogdir(fileDir)
	if !self.isExist(fileDir + "/" + fileName) {
		os.Create(fileDir + "/" + fileName)
	}
	self.logObj = &_FILE{dir: fileDir, name: fileName, mu: new(sync.RWMutex)}

	// 锁定文件操作
	self.logObj.mu.Lock()
	defer self.logObj.mu.Unlock()

	self.handleFile4Daily()

	// 1s检测一次文件是否需要滚动
	go func(logger *Logger) {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				logger.handleFile4Daily()
			}
		}
	}(self)
}

// 打印DEBUG级别日志
func (self *Logger) Debug(v ...interface{}) {
	self.filePrint(DEBUG, "debug", v)
}

// 打印INFO级别日志
func (self *Logger) Info(v ...interface{}) {
	self.filePrint(INFO, "info", v)
}

// 打印WARN级别日志
func (self *Logger) Warn(v ...interface{}) {
	self.filePrint(WARN, "warn", v)
}

// 打印ERROR级别日志
func (self *Logger) Error(v ...interface{}) {
	self.filePrint(ERROR, "error", v)
}

// 打印FATAL级别日志
func (self *Logger) Fatal(v ...interface{}) {
	self.filePrint(FATAL, "fatal", v)
}

// 文件打印
func (self *Logger) filePrint(level LEVEL, levelStr string, v ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	// 给日志文件加锁
	if self.logObj != nil {
		self.logObj.mu.RLock()
		defer self.logObj.mu.RUnlock()
	}

	// 写日志
	if self.logLevel >= level {
		if self.logObj != nil {
			self.logObj.lg.Output(3, fmt.Sprintln(levelStr, v))
		}
		self.consolePrint(levelStr, v)
	}
}

// 控制台打印
func (self *Logger) consolePrint(v ...interface{}) {
	if self.consoleFlag {
		// 获取文件名&行号
		_, file, line, _ := runtime.Caller(3)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short

		// 控制台打印
		log.Println(file, strconv.Itoa(line), v)
	}
}

// 创建日志目录
func (self *Logger) mklogdir(dir string) {
	_, err := os.Stat(dir)
	if !(err == nil || os.IsExist(err)) {
		if err = os.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
	}
}

// 检测及创建日志文件（大小滚动）
func (self *Logger) handleFile4Size() {
	// 获取当前日志文件全路径
	filePath := self.logObj.dir + "/" + self.logObj.name

	// 获取文件信息
	file, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	// 文件大小达到设定最大值，创建新文件
	if file.Size() >= self.maxFileSize {
		// 日志文件正在使用，先关闭
		if self.logObj.logfile != nil {
			self.logObj.logfile.Close()
		}

		// 将原日志文件该名
		suffix := fmt.Sprintf(".%d", time.Now().Unix())
		if self.isExist(filePath + suffix) {
			os.Remove(filePath + suffix)
		}
		os.Rename(filePath, filePath+suffix)

		// 创建新的日志文件
		self.logObj.logfile, _ = os.Create(filePath)
		self.logObj.lg = log.New(self.logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// 打开之前的日志文件
		self.logObj.logfile, _ = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		self.logObj.lg = log.New(self.logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// 检测及创建日志文件（日期滚动）
func (self *Logger) handleFile4Daily() {
	// 获取当前日志文件全路径
	filePath := self.logObj.dir + "/" + self.logObj.name

	// 获取当前时间
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))

	if t.After(*self.currentFileDate) {
		// 日志文件正在使用，先关闭
		if self.logObj.logfile != nil {
			self.logObj.logfile.Close()
		}

		// 将原日志文件该名
		suffix := "." + self.currentFileDate.Format(DATEFORMAT)
		if self.isExist(filePath + suffix) {
			os.Remove(filePath + suffix)
		}
		os.Rename(filePath, filePath+suffix)

		// 创建新的日志文件
		self.logObj.logfile, _ = os.Create(filePath)
		self.logObj.lg = log.New(self.logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

		// 更新日期
		self.currentFileDate = &t
	} else {
		// 打开之前的日志文件
		self.logObj.logfile, _ = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		self.logObj.lg = log.New(self.logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// 判断文件是否存在
func (self *Logger) isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
