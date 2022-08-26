package zero_mdc_log

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	config_pack "gitee.com/aiyuangong_group/zero_mdc_log/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Logger *zerolog.Logger // 需要传
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC            bool           // 默认即可
	SkipPath       []string       // 默认即可
	SkipPathRegexp *regexp.Regexp // 默认即可
	requestIdKey   string         // 从header获取
	userIdKey      string         // 从header获取
	sessionIdKey   string         // 默认没有
	appId          string         // 需要传
	envType        string         // 需要传
}

type ConfigItem interface {
	applyConfigItem(*Config)
}

type configFunc func(*Config)

func (f configFunc) applyConfigItem(c *Config) {
	f(c)
}

// ConfigEnvtype 环境
func SetEnvtype(envType string) ConfigItem {
	return configFunc(func(c *Config) {
		c.envType = envType
	})
}

func SetAppId(appId string) ConfigItem {
	return configFunc(func(c *Config) {
		c.appId = appId
	})
}

func SetSessionid(sessionIdKey string) ConfigItem {
	return configFunc(func(c *Config) {
		c.sessionIdKey = sessionIdKey
	})
}

func SetRequestIdKey(requestIdKey string) ConfigItem {
	return configFunc(func(c *Config) {
		c.requestIdKey = requestIdKey
	})
}

func SetUserIdKey(userIdKey string) ConfigItem {
	return configFunc(func(c *Config) {
		c.userIdKey = userIdKey
	})
}

func SetLogger(configItems ...ConfigItem) gin.HandlerFunc {
	cfg := Config{
		Logger:       GetZeroLogger(),
		requestIdKey: viper.GetString(config_pack.REQUEST_ID_KEY),
		userIdKey:    viper.GetString(config_pack.USER_ID_KEY),
		sessionIdKey: viper.GetString(config_pack.SESSION_ID_KEY),
		appId:        viper.GetString(config_pack.APP_ID),
		envType:      viper.GetString(config_pack.ENV_TYPE),
	}

	for _, configItem := range configItems {
		configItem.applyConfigItem(&cfg)
	}

	var skip map[string]struct{}
	if length := len(cfg.SkipPath); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range cfg.SkipPath {
			skip[path] = struct{}{}
		}
	}

	var sublog zerolog.Logger
	if cfg.Logger == nil {
		sublog = log.Logger
	} else {
		sublog = *cfg.Logger
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		track := true

		if _, ok := skip[path]; ok {
			track = false
		}

		if track &&
			cfg.SkipPathRegexp != nil &&
			cfg.SkipPathRegexp.MatchString(path) {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if cfg.UTC {
				end = end.UTC()
			}

			userID := c.GetHeader(cfg.userIdKey)
			if userID == "" {
				if u, ok := c.Get("userid"); ok {
					userID = fmt.Sprintf("%v", u)
				}

				if userID == "" {
					userID = "nouserid"
				}

				c.Header(cfg.userIdKey, userID)
			}

			requestID := c.GetHeader(cfg.requestIdKey)
			if requestID == "" {
				requestID = uuid.New().String()
				c.Header(cfg.requestIdKey, requestID)
			}

			sessionID := c.GetHeader(cfg.sessionIdKey)
			if sessionID == "" {
				sessionID = "nosessionid"
				c.Header(cfg.sessionIdKey, sessionID)
			}

			msg := "Request"
			if len(c.Errors) > 0 {
				msg = c.Errors.String()
			}

			dumplogger := sublog.With().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("requestid", requestID).
				Str("userid", userID).
				Str("sessionid", sessionID).
				Str("appid", cfg.appId).
				Str("envtype", cfg.envType).
				Logger()

			switch {
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				{
					dumplogger.Warn().
						Msg(msg)
				}
			case c.Writer.Status() >= http.StatusInternalServerError:
				{
					dumplogger.Error().
						Msg(msg)
				}
			default:
				dumplogger.Info().
					Msg(msg)
			}
		}

	}
}
