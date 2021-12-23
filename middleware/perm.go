package middleware

import (
	"deploy_server/config"
	"deploy_server/pkg/casbin"
	"deploy_server/pkg/context"
	"deploy_server/pkg/token"
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

var (
	auth *casbin.Casbin
	once sync.Once
)

var (
	noHeaderErr = errors.New("Header 中缺少 Authorization 参数")
	noPermErr   = errors.New("验证权限失败, 您无权限访问")
)

func GetCasbin() *casbin.Casbin {
	once.Do(func() {
		auth = casbin.New()
		auth.Load()
	})

	return auth
}

// PermissionMiddleWare 基于Casbin权限系统的 JWT Token中间件
func PermissionMiddleWare() context.HandlerFunc {
	return func(c context.Context) {
		authorization := c.GetHeader("Authorization")
		c.Logger().Debug(fmt.Sprintf("Authorization: %s", authorization))
		if authorization == "" {
			c.Fail(10001, noHeaderErr)

			return
		}

		cfg := config.Get().JWT
		claims, err := token.New(cfg.Secret).JwtParse(authorization)

		if err != nil {
			c.Fail(10001, err)

			return
		}

		//c.Logger().Debug(fmt.Sprintf("%v", claims))
		ok := GetCasbin().Enforce(claims.Role, c.Path(), c.Method())

		if ok {
			return
		} else {
			c.Fail(10002, noPermErr)
		}
	}
}
