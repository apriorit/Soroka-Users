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

type SessionSecret struct {
	Secret string
}

type publicKey struct {
	pubKey []byte
}

var loggerConfig *LoggerConfig
var priveledgeConfig *PriveledgeConfig
var secret *SessionSecret
var pkey *publicKey

func GetSecretString() string {
	if secret == nil {
		secret = &SessionSecret{
			Secret: "secret",
		}
	}
	return secret.Secret
}

func SetSecretString(s string) {
	if secret == nil {
		secret = &SessionSecret{
			Secret: s,
		}
	}
}

func GetPublicKey() []byte {
	if pkey == nil {
		pkey = &publicKey{
			pubKey: make([]byte, 0),
		}
	}
	return pkey.pubKey
}

func SetPublicKey(key []byte) {
	if pkey == nil {
		pkey = &publicKey{
			pubKey: key,
		}
	}
}

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
		privMap["changeUserStatus"] = c.ChangeUserStatus
		priveledgeConfig = &PriveledgeConfig{
			Priveledges: privMap,
		}
	}

	return priveledgeConfig
}
