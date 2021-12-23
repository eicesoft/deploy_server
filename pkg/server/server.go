package server

import (
	"deploy_server/pkg/cache"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
)

func NewServer(logger *zap.Logger, db db.Repo, cache cache.Repo, opt *core.Option) (*core.Mux, error) {
	if opt.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	s := core.NewServer(logger, db, cache, opt)

	return s, nil
}
