package config

import (
	"gitee.com/wxlao/config-client"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	viper *viper.Viper
}

// New 新配置
func New(viper *viper.Viper) *Config {
	return &Config{
		viper: viper,
	}
}

// Load 加载配置
func (c *Config) Load() error {
	c.viper.SetConfigType("env")

	if err := config.LoadFile(".env"); err != nil {
		log.Error().Str("err", err.Error()).Msg("Error reading config file")
	}

	c.setDefault()

	return nil
}

const (
	SESSION_ID_KEY  = "SESSION_ID_KEY"
	REQUEST_ID_KEY  = "REQUEST_ID_KEY"
	USER_ID_KEY     = "USER_ID_KEY"
	APP_ID          = "APP_ID"
	ENV_TYPE        = "ENV_TYPE"
	LOG_LEVEL       = "LOG_LEVEL"
	LOG_FILENAME    = "LOG_FILENAME"
	LOG_MAX_SIZE    = "LOG_MAX_SIZE"
	LOG_MAX_AGE     = "LOG_MAX_AGE"
	LOG_MAX_BACKUPS = "LOG_MAX_BACKUPS"
	LOG_JSON        = "LOG_JSON"
	LOG_CONSOLE     = "LOG_CONSOLE"
)

func (c *Config) setDefault() {
	c.viper.SetDefault(REQUEST_ID_KEY, "logcontext-requestid")
	c.viper.SetDefault(USER_ID_KEY, "logcontext-userid")
	c.viper.SetDefault(SESSION_ID_KEY, "Ayg-Sessionid")
	c.viper.SetDefault(APP_ID, "an-app")
	c.viper.SetDefault(ENV_TYPE, "pro")

	c.viper.SetDefault(LOG_LEVEL, 0)                                   // 日志等级 Debug:0 Info:1 Warn:2 Error:3 Fatal:4 Panic:5 No:6 Disabled:7
	c.viper.SetDefault(LOG_FILENAME, c.viper.GetString(APP_ID)+".log") // 日志文件名
	c.viper.SetDefault(LOG_MAX_SIZE, 100)                              // 日志单文件大小 mb
	c.viper.SetDefault(LOG_MAX_AGE, 7)                                 // 日志天数
	c.viper.SetDefault(LOG_MAX_BACKUPS, 10)                            // 日志文件数
	c.viper.SetDefault(LOG_JSON, false)                                // 日志保存JSON
	c.viper.SetDefault(LOG_CONSOLE, false)                             // 日志输出到stdout
}
