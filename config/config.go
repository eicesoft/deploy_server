package config

import (
	"deploy_server/pkg/env"
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var config *Config
var once sync.Once

type Config struct {
	Server struct {
		Name string `toml:"name"`
	} `toml:"Server"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
	} `toml:"redis"`

	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
		} `toml:"base"`
	} `toml:"mysql"`

	JWT struct {
		Secret         string        `toml:"secret"`
		ExpireDuration time.Duration `toml:"expireDuration"`
	} `toml:"jwt"`
}

// Get 获得配置
func Get() Config {
	once.Do(func() {
		config = &Config{}
		viper.SetConfigName(env.Env + "_config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(env.ConfigPath)

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})

	return *config
}

func ProjectLogFile() string {
	return fmt.Sprintf("./storage/logs/server-%s.log", env.Env)
}
