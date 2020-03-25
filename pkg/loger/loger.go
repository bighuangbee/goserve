package loger

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"goserve/pkg/file"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var loger *logrus.Logger
var fileRootPath string

func Setup(path string){
	fileRootPath = path

	if err := file.IsNotExistMkDir(fileRootPath); err != nil {
		fmt.Errorf("Loger Setup Error ###", err.Error())
		return
	}

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Errorf("Loger File Open Error ###", err.Error())
	}

	loger = logrus.New()
	loger.SetReportCaller(true)

	//输出到文件
	loger.Out = src

	logerFormatter := new(LogerFormatter)
	loger.SetFormatter(logerFormatter)

	loger.AddHook(newLocalFileLogHook(logrus.ErrorLevel, logerFormatter))
	loger.AddHook(newLocalFileLogHook(logrus.InfoLevel, logerFormatter))

	loger.SetOutput(os.Stdout)

	Info("Loger SetUp Success...")
}


/**
	写入本地日志文件， 按日期、日志级别分割为不同的文件
 */
func newLocalFileLogHook(level logrus.Level, formatter logrus.Formatter) logrus.Hook {

	fileName := filepath.Join(fileRootPath, level.String() + "_%Y%m%d%H.log")

	//文件分割
	writer, err := rotatelogs.New(
		fileName,
		// 最大保存时间(30天)
		rotatelogs.WithMaxAge(30*24*time.Hour),
		// 日志分割间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		fmt.Errorf("config local file system for loger error: %v", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		level: writer,
	}, formatter)

}

func Infof(format string, args ...interface{}){
	loger.Infof(format, args)
}

func Info(args ...interface{}){
	loger.Info(args)
}

func Error(args ...interface{}){
	loger.Error(args)
}


type LogerFormatter struct{}

/*
	日志输出格式
 */
func (s *LogerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] %s %s \n",  strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	return []byte(msg), nil
}