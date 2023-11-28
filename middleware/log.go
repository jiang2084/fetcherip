package middleware

import (
	"bytes"
	"encoding/json"
	"fetchip/utils"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func InitLog(setting *Logs) {
	file, logFilePath := getLogFile(setting.FilePath, setting.FileName)
	mode := strings.Replace(strings.ToLower(setting.Mode), " ", "", -1)
	switch mode {
	case "console":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logrus.SetOutput(os.Stdout)
	case "file":
		logrus.SetOutput(file)
	case "console,file":
		// 日志输出到文件同时输出到控制台
		multiWriter := io.MultiWriter(os.Stdout, file)
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logrus.SetOutput(multiWriter)
	default:
		logrus.SetOutput(file)
	}
	logrus.SetLevel(logrus.InfoLevel)
	// 设置rotatelogs日志分割Hook
	logrus.AddHook(NewLfsHook(logFilePath))
}

func NewLfsHook(filePath string) logrus.Hook {
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		filePath+".%Y%m%d.log",
		// 生成软链，指向最新日志文件, WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(filePath),
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// 设置文件清理前的最长保存时间, 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置文件清理前最多保存的个数
		//rotatelogs.WithRotationCount(5),
		// 设置日志分割的时间,这里设置为一天分割一次, 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		logrus.Errorf("config logger error: %v\n", err)
	}

	writeMap := lfshook.WriterMap{
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfsHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return lfsHook
}

func getLogFile(logDir string, logFileName string) (*os.File, string) {
	// 目录不存在创建一个
	if !utils.PathExists(logDir) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Warnf("create directory failed: %v\n", err)
			os.Exit(-1)
		}
	}

	// 文件路径
	logFilePath := path.Join(logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("Open File err: %v\n", err)
		os.Exit(-1)
	}
	return file, logFilePath
}

// 日志自定义格式
type LogFormatter struct{}

// 格式详情
func (logFormat *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	level := strings.ToUpper(entry.Level.String())
	content, _ := json.Marshal(entry.Data)
	msg := fmt.Sprintf("%s [%s] [GOID:%d] [%s:%d] #msg:%s #content:%v\n",
		timestamp,
		level,
		getGID(),
		file,
		line,
		entry.Message,
		content,
	)
	return []byte(msg), nil
}

// 获取当前协程id
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
