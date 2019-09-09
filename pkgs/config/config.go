package config

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	c "github.com/Soroka-EDMS/svc/users/pkgs/constants"
)

type LoggerConfig struct {
	Logger log.Logger
}

type PriveledgeConfig struct {
	Priveledges map[string]int
}

var loggerConfig *LoggerConfig
var priveledgeConfig *PriveledgeConfig

func GetLogger() *LoggerConfig {
	if loggerConfig == nil {
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			logger = level.NewFilter(logger, level.AllowAll())
			logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
			logger = log.With(logger, "service", "users")
		}
		loggerConfig = &LoggerConfig{
			Logger: logger,
		}
	}
	return loggerConfig
}

func GetPriveledges() *PriveledgeConfig {
	if priveledgeConfig == nil {
		privMap := make(map[string]int)
		privMap["changeRole"] = c.ChangeRole
		privMap["userList"] = c.QueryUsers
		privMap["userProfile"] = c.GetProfile
		privMap["disableUser"] = c.ChangeUserStatus
		privMap["enableUser"] = c.ChangeUserStatus
		priveledgeConfig = &PriveledgeConfig{
			Priveledges: privMap,
		}
	}

	return priveledgeConfig
}
