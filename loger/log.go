package log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type CustomTag string
type Method string

const (
	None      CustomTag = ""
	NormalTag CustomTag = "NORMAL"
	FatalTag  CustomTag = "FATAL"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		NoColors:        true,
		NoFieldsColors:  true,
	})
	if len(os.Args) > 1 {
		mode := os.Args[1]
		if strings.EqualFold(mode, "k8dev") ||
			strings.EqualFold(mode, "k8stg") ||
			strings.EqualFold(mode, "k8online") {
			logrus.Info("init hook")
			initHook()
		}
	}
}

func initHook() {
	//修改gin日志输出
	writer, err := rotatelogs.New(
		"/data/log/fuxi-data-web-qb.log.%Y%m%d",
		rotatelogs.WithLinkName("/data/log/fuxi-data-web-qb.log"),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithRotationSize(1024*1024*500),
	)
	if nil != err {
		logrus.Info("initHook fail:", err.Error())
		return
	}
	//修改gin日志输出
	w := io.MultiWriter(os.Stdout, writer)
	gin.DefaultWriter = w

	hk := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		NoColors:        true,
		NoFieldsColors:  true,
	})
	logrus.AddHook(hk)
}
