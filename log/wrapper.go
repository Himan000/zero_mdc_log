package log

import (
	"io"
	"os"
	"strings"
	"time"

	"gitee.com/aiyuangong_group/zero_mdc_log/config"

	over "gitee.com/aiyuangong_group/zero_mdc_log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

const TRACE_ID = "traceID"

func Init() {
	//项目配置
	c := config.New(viper.GetViper())
	_ = c.Load()

	fileOut := &lumberjack.Logger{
		Filename:   viper.GetString("LOG_FILENAME"), // 文件
		MaxSize:    viper.GetInt("LOG_MAX_SIZE"),    // megabytes
		MaxBackups: viper.GetInt("LOG_MAX_BACKUPS"), // MaxBackups
		MaxAge:     viper.GetInt("LOG_MAX_AGE"),     // days
		LocalTime:  true,                            // 这个需要设置, 不然日志文件的名字就是UTC时间
	}

	var out io.Writer
	if viper.GetBool("LOG_CONSOLE") {
		out = io.MultiWriter(os.Stdout, fileOut)
	} else {
		out = fileOut
	}

	if viper.GetBool("LOG_JSON") {
		log.Logger = log.Output(out)
	} else {
		writer := zerolog.NewConsoleWriter()
		writer.NoColor = true
		writer.TimeFormat = time.RFC3339
		writer.FormatLevel = func(i interface{}) string {
			if ll, ok := i.(string); ok {
				return strings.ToUpper("[" + ll + "]")
			}
			return "[???]" // level为空
		}

		writer.Out = out
		log.Logger = log.Output(writer)
	}
	log.Level(zerolog.Level(viper.GetUint("LOG_LEVEL")))
	over.New(log.Logger)
	over.AddGlobalFields(TRACE_ID)
}

func MDC() *over.MdcAdapter {
	return over.MDC()
}

func Log() *over.Overlog {
	return over.Log()
}

func Info() *zerolog.Event {
	return over.Log().Info()
}

func Debug() *zerolog.Event {
	return over.Log().Debug()
}

func Error() *zerolog.Event {
	return over.Log().Error()
}
