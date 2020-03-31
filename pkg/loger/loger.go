package loger

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"goserve/pkg/file"
	"os"
	"path/filepath"
	"runtime"
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
	setPrefix("Infof")
	loger.Infof(format, args)
}

func Info(args ...interface{}){
	setPrefix("Info")
	loger.Info(args)
}

func Error(args ...interface{}){
	setPrefix("Error")
	loger.Error(args)
}


type LogerFormatter struct{}


// setPrefix set the prefix of the log output
func setPrefix(level string) {

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.TrimPrefix(filepath.Ext(funcName), ".")
		timestamp := time.Now().Local().Format("2006-01-02 15:04:05")

		fmt.Printf("[%s][%s][%s:%d:%s]", strings.ToUpper(level), timestamp, filepath.Base(file), line, funcName)
	}
}

/*
	日志输出格式
 */
func (s *LogerFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	msg := fmt.Sprintf(" %s \n", entry.Message)
	return []byte(msg), nil
}