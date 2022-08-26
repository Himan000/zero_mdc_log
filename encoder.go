package zero_mdc_log

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// 这个为什么叫Lowercase呢，默认就是小写？
func LowercaseLevelEncoder() func(l zerolog.Level) string {
	return func(l zerolog.Level) string {
		return l.String()
	}
}

func UppercaseLevelEncoder() func(l zerolog.Level) string {
	return func(l zerolog.Level) string {
		switch l {
		case zerolog.TraceLevel:
			return "TRACE"
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARN"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "FATAL"
		case zerolog.PanicLevel:
			return "PANIC"
		case zerolog.NoLevel:
			return ""
		}
		return ""
	}
}

// 拼接文件和行号
func FullCallerEncoder() func(file string, line int) string {
	return func(file string, line int) string {
		return file + ":" + strconv.Itoa(line)
	}
}

// 只有文件没有路径
func ShortCallerEncoder() func(file string, line int) string {
	return func(file string, line int) string {
		return TrimmedPath(file) + ":" + strconv.Itoa(line)
	}
}

// 只有文件没有路径
func TrimmedPath(file string) string {
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file
	}
	return file[idx+1:]
}
