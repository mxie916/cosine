// Copyright 2016 mxie916@163.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cosine

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

// 日期格式（用于文件名）
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

// 文件大小
const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

// 日志结构体
type Logger struct {
	logLevel        LEVEL
	consoleFlag     bool
	rollingFileFlag bool
	dailyFileFlag   bool
	maxFileSize     int64
	currentFileDate *time.Time
	logObj          *_FILE
	lg              *log.Logger
}

// 实例化Logger对象
func newLogger() *Logger {
	logger := new(Logger)
	if os.Getenv("log.level") == "" {
		logger.SetLevel(ALL)
	} else {
		switch strings.ToLower(os.Getenv("log.level")) {
		case "all":
			logger.SetLevel(ALL)
		case "debug":
			logger.SetLevel(DEBUG)
		case "info":
			logger.SetLevel(INFO)
		case "warn":
			logger.SetLevel(WARN)
		case "error":
			logger.SetLevel(ERROR)
		case "fatal":
			logger.SetLevel(FATAL)
		case "off":
			logger.SetLevel(OFF)
		default:
			panic("日志级别配置错误")
		}
	}
	if os.Getenv("log.console") != "" {
		switch strings.ToLower(os.Getenv("log.console")) {
		case "true":
			logger.SetConsole(true)
		case "false":
			logger.SetConsole(false)
		default:
			panic("控制台日志输出配置错误")
		}
	} else {
		logger.SetConsole(true)
	}
	if strings.ToLower(os.Getenv("log.rollingfile")) == "true" {
		maxSize, err := strconv.ParseInt(os.Getenv("log.maxsize"), 10, 64)
		if err != nil {
			panic("日志滚动大小配置错误")
		}
		var _unit UNIT
		switch strings.ToLower(os.Getenv("log.sizeunit")) {
		case "kb":
			_unit = KB
		case "mb":
			_unit = MB
		case "gb":
			_unit = GB
		case "tb":
			_unit = TB
		default:
			panic("日志大小单位配置错误")
		}
		logger.SetRollingFile(os.Getenv("log.dir"), os.Getenv("log.file"), maxSize, _unit)
	}
	if strings.ToLower(os.Getenv("log.dailyfile")) == "true" {
		logger.SetDailyFile(os.Getenv("log.dir"), os.Getenv("log.file"))
	}

	return logger
}

// 设置日志级别
func (self *Logger) SetLevel(level LEVEL) {
	self.logLevel = level
}

// 设置是否在控制台打印
func (self *Logger) SetConsole(isConsole bool) {
	self.consoleFlag = isConsole
	if isConsole {
		self.lg = log.New(os.Stdout, "[Cosine] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
}

// 设置是否采用文件打印（按文件大小滚动）
func (self *Logger) SetRollingFile(fileDir, fileName string, maxSize int64, _unit UNIT) {
	self.rollingFileFlag = true
	self.dailyFileFlag = false

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

	self.chkFile4Size()

	// 1s检测一次文件是否需要滚动
	go func(logger *Logger) {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				logger.chkFile4Size()
			}
		}
	}(self)
}

// 设置是否采用文件打印（按文件日期滚动）
func (self *Logger) SetDailyFile(fileDir, fileName string) {
	self.dailyFileFlag = true
	self.rollingFileFlag = false

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

	self.chkFile4Daily()

	// 1s检测一次文件是否需要滚动
	go func(logger *Logger) {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				logger.chkFile4Daily()
			}
		}
	}(self)
}

// 获取设置的日志级别
func (self *Logger) GetLevel() LEVEL {
	return self.logLevel
}

// 打印DEBUG级别日志
func (self *Logger) Debug(v ...interface{}) {
	self.filePrint(DEBUG, "[DEBUG]", v)
}

// 打印INFO级别日志
func (self *Logger) Info(v ...interface{}) {
	self.filePrint(INFO, "[INFO ]", v)
}

// 打印WARN级别日志
func (self *Logger) Warn(v ...interface{}) {
	self.filePrint(WARN, "[WARN ]", v)
}

// 打印ERROR级别日志
func (self *Logger) Error(v ...interface{}) {
	self.filePrint(ERROR, "[ERROR]", v)
}

// 打印FATAL级别日志
func (self *Logger) Fatal(v ...interface{}) {
	self.filePrint(FATAL, "[FATAL]", v)
}

// 文件打印
func (self *Logger) filePrint(l LEVEL, level string, v ...interface{}) {
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
	if l >= self.logLevel {
		// 将需要打印的内容拼接成字符串
		msg := ""
		for i := 0; i < len(v); i++ {
			if i > 0 {
				msg += " "
			}
			msg += fmt.Sprint(v[i])
		}

		// 打印
		if self.logObj != nil {
			self.logObj.lg.Output(3, level+" "+msg)
		}
		self.consolePrint(level, msg)
	}
}

// 控制台打印
func (self *Logger) consolePrint(level, msg string) {
	if self.consoleFlag {
		self.lg.Output(4, level+" "+msg)
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
func (self *Logger) chkFile4Size() {
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
		self.logObj.lg = log.New(self.logObj.logfile, "[Cosine] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	} else {
		if self.logObj.logfile == nil || self.logObj.lg == nil {
			// 打开之前的日志文件
			self.logObj.logfile, _ = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			self.logObj.lg = log.New(self.logObj.logfile, "[Cosine] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		}
	}
}

// 检测及创建日志文件（日期滚动）
func (self *Logger) chkFile4Daily() {
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
		self.logObj.lg = log.New(self.logObj.logfile, "[Cosine] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

		// 更新日期
		self.currentFileDate = &t
	} else {
		if self.logObj.logfile == nil || self.logObj.lg == nil {
			// 打开之前的日志文件
			self.logObj.logfile, _ = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			self.logObj.lg = log.New(self.logObj.logfile, "[Cosine] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		}
	}
}

// 判断文件是否存在
func (self *Logger) isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
