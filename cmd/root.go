package cmd

import (
	"context"
	"deploy_server/config"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
	"deploy_server/pkg/env"
	log "deploy_server/pkg/logger"
	"deploy_server/pkg/server"
	"deploy_server/pkg/shutdown"
	"deploy_server/router"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

var rootCmd = &cobra.Command{
	Use:               "deploy_server",
	Short:             "deploy_server is web server",
	Long:              "deploy_server is web server",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run:               Run,
}

var (
	logger *zap.Logger
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&env.DebugFlag, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&env.Env, "env", "e", "dev", "server run mode")
	rootCmd.PersistentFlags().StringVarP(&env.ServerPort, "port", "p", "9105", "server run port")
	rootCmd.PersistentFlags().StringVarP(&env.ConfigPath, "config", "c", "./config", "server config file path")
}

// Execute 根命令执行
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func Run(cmd *cobra.Command, args []string) {
	logger := initLogger()
	// 初始化ORM
	dbRepo := initDB()
	// 初始化 Cache
	cacheRepo := initCache()

	core.DebugLog("* listen port %s", env.ServerPort)
	handleServer, err := server.NewServer(logger, dbRepo, cacheRepo, &core.Option{
		Env:     env.Env,
		IsDebug: env.DebugFlag,
		Port:    env.ServerPort,
	})

	if err != nil {
		logger.Fatal("Create server handler error: " + err.Error())
	}

	router.RegistryRouters(handleServer, dbRepo, cacheRepo)

	handleServer.Run()

	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			if err := handleServer.Shutdown(ctx); err != nil {
				logger.Error("server shutdown err", zap.Error(err))
			} else {
				logger.Info("server shutdown success")
			}
		},
		func() {
			if err := dbRepo.DbWClose(); err != nil {
				logger.Error("dbw close err", zap.Error(err))
			} else {
				logger.Info("dbw close success")
			}

			if err := dbRepo.DbRClose(); err != nil {
				logger.Error("dbr close err", zap.Error(err))
			} else {
				logger.Info("dbr close success")
			}
		},
	)
}

// initDB 初始化DB连接
func initDB() db.Repo {
	dbRepo, err := db.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	return dbRepo
}

// initLogger 初始化Logger
func initLogger() *zap.Logger {
	logger, err := log.NewJSONLogger(
		log.WithDebugLevel(),
		log.WithField("app", fmt.Sprintf("%s[%s]", config.Get().Server.Name, env.Env)),
		log.WithTimeLayout("2006-01-02 15:04:05"),
		log.WithFileP(config.ProjectLogFile()),
	)

	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger
}

func initCache() cache.Repo {
	cacheRepo, err := cache.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}

	return cacheRepo
}
