package utils

import (
	"chat/globals"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ReadConf() {
	viper.SetConfigFile(configFile)

	if !IsFileExist(configFile) {
		fmt.Printf("[service] config.yaml not found, creating one from template: %s\n", configExampleFile)
		if err := CopyFile(configExampleFile, configFile); err != nil {
			fmt.Println(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	secret := viper.GetString("secret")
	if len(secret) < 32 {
		globals.Warn(fmt.Sprintf("[service] invalid secret length: got %d, expected at least 32 bytes; starting in 10 seconds, please set a stronger `secret` in config or environment; future versions may panic on weak secrets", len(secret)))
		time.Sleep(10 * time.Second)
	}

	if timeout := viper.GetInt("max_timeout"); timeout > 0 {
		globals.HttpMaxTimeout = time.Second * time.Duration(timeout)
		globals.Debug(fmt.Sprintf("[service] http client timeout set to %ds from env", timeout))
	}
}

// trustedProxies restricts which upstream hop c.ClientIP() (and therefore
// every per-IP rate limiter and IP ban list) is allowed to trust for
// X-Forwarded-For / X-Real-IP. Without this, Gin trusts those headers from
// any direct client, letting an attacker set an arbitrary value and rotate
// it per request to defeat IP-based throttling and IP ban lists entirely
// (e.g. brute-forcing the 6-digit /reset email code). Only loopback is
// trusted, matching the bundled sample nginx.conf, which proxies from the
// same host; a request that didn't actually come through that proxy falls
// back to its real socket address instead of an attacker-supplied header.
var trustedProxies = []string{"127.0.0.1", "::1"}

func NewEngine() *gin.Engine {
	if viper.GetBool("debug") {
		engine := gin.Default()
		_ = engine.SetTrustedProxies(trustedProxies)
		return engine
	}

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	_ = engine.SetTrustedProxies(trustedProxies)
	return engine
}
