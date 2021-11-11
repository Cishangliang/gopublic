package log

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"path"

	"github.com/sirupsen/logrus"
)

const (
	requestIDKey    = "requestID"
	requestRouteKey = "requestRout"
	customTagKey    = "customTag"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	FieldsOrder     []string // default: fields sorted alphabetically
	TimestampFormat string   // default: time.StampMilli = "Jan _2 15:04:05.000"
	HideKeys        bool     // show [fieldValue] instead of [fieldKey:fieldValue]
	NoColors        bool     // disable colors
	NoFieldsColors  bool     // color only level, default is level + fields
	ShowFullLevel   bool     // true to show full level [WARNING] instead [WARN]
	TrimMessages    bool     // true to trim whitespace on messages
}

// Format an log entry
//[时间]|[级别]|[自定义标记]|[线程]|[类.函数]|[行数]|[请求方法 请求路径]|[终端请求编号]|自定义消息
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	baseF := "[%s]|[%s]|[%s]|[%s]|[%s]|[%d]|[%s]|[%s]|%s\n"
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}
	t := entry.Time.Format(timestampFormat)
	level := strings.ToUpper(entry.Level.String())
	cTag := ""
	if v, ok := entry.Data[customTagKey]; ok {
		cTag = string(v.(CustomTag))
	} else {
		if strings.EqualFold(level, "Fatal") {
			cTag = string(FatalTag)
		} else {
			cTag = ""
		}
	}
	_, file := path.Split(entry.Caller.File)

	// output buffer
	b := &bytes.Buffer{}
	_, _ = fmt.Fprintf(b, baseF, t, level, cTag, "", file, entry.Caller.Line, "", "", entry.Message)
	return b.Bytes(), nil
}

// Format an log entry
//func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
//	levelColor := getColorByLevel(entry.Level)
//
//	timestampFormat := f.TimestampFormat
//	if timestampFormat == "" {
//		timestampFormat = time.StampMilli
//	}
//
//	// output buffer
//	b := &bytes.Buffer{}
//
//	// write time
//	b.WriteString(entry.Time.Format(timestampFormat))
//
//	// write level
//	level := strings.ToUpper(entry.Level.String())
//
//	if !f.NoColors {
//		fmt.Fprintf(b, "\x1b[%dm", levelColor)
//	}
//
//	b.WriteString(" [")
//	if f.ShowFullLevel {
//		b.WriteString(level)
//	} else {
//		b.WriteString(level[:4])
//	}
//	b.WriteString("] ")
//
//	if !f.NoColors && f.NoFieldsColors {
//		b.WriteString("\x1b[0m")
//	}
//
//	if entry.HasCaller() {
//		_, file := path.Split(entry.Caller.File)
//		//funs := strings.Split(entry.Caller.Function, ".")
//		//len := len(funs)
//		fmt.Fprintf(
//			b,
//			"[%s:%d] ",
//			file,
//			entry.Caller.Line,
//		)
//	}
//
//	// write fields
//	if f.FieldsOrder == nil {
//		f.writeFields(b, entry)
//	} else {
//		f.writeOrderedFields(b, entry)
//	}
//
//	if !f.NoColors && !f.NoFieldsColors {
//		b.WriteString("\x1b[0m")
//	}
//
//	// write message
//	if f.TrimMessages {
//		b.WriteString(strings.TrimSpace(entry.Message))
//	} else {
//		b.WriteString(entry.Message)
//	}
//
//	b.WriteByte('\n')
//
//	return b.Bytes(), nil
//}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
