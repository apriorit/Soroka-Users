package config

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

type LoggerConfig struct {
	Logger log.Logger
}

type PrivilegesConfig struct {
	Privileges map[string]int
}

var loggerConfig *LoggerConfig
var privilegesConfig *PrivilegesConfig

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

func GetPrivileges() *PrivilegesConfig {
	if privilegesConfig == nil {
		privMap := make(map[string]int)
		privMap["changeRole"] = stub.ChangeRole
		privMap["userList"] = stub.QueryUsers
		privMap["userProfile"] = stub.GetProfile
		privMap["changeUserStatus"] = stub.ChangeUserStatus
		privilegesConfig = &PrivilegesConfig{
			Privileges: privMap,
		}
	}

	return privilegesConfig
}
