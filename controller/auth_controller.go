package controller

import (
	"deploy_server/config"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/context"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
	"deploy_server/pkg/token"
	"deploy_server/service"
	"github.com/pkg/errors"
	"time"
)

// authRequest 登录请求
type authRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// authData 登录请求响应
type authData struct {
	Authorization string `json:"authorization"` // 签名
	ExpireTime    int64  `json:"expire_time"`   // 过期时间
}

type authHandler struct {
	cache      cache.Repo
	userServer *service.UserService
}

// Auth 登录验证获取Token
// @Summary 登录验证获取Token
// @Description 登录验证获取Token
// @Tags Auth
// @Produce  json
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {string} json "{"code":200,"data":{"authorization": "", "expire_time": 0}, "message":""}"
// @Failure 400 {object} core.Failure
// @Router /auth [get]
// @Router /auth [post]
func (handle *authHandler) Auth() context.HandlerFunc {
	return func(c context.Context) {
		req := new(authRequest)
		resp := new(authData)
		if c.Method() == "GET" {
			if err := c.ShouldBindQuery(req); err != nil {
				c.Fail(10001, err)
				return
			}
		} else {
			if err := c.ShouldBindForm(req); err != nil {
				c.Fail(10001, err)
				return
			}
		}

		user := handle.userServer.Check(c.RequestContext(), req.Username, req.Password)
		if user == nil {
			c.Fail(10001, errors.New("用户名密码错误"))
			return
		}

		cfg := config.Get().JWT
		tokenString, err := token.New(cfg.Secret).JwtSign(user.Id, user.Name, user.Role, time.Hour*cfg.ExpireDuration)
		if err != nil {
			c.Fail(10001, err)
			return
		}
		resp.Authorization = tokenString
		resp.ExpireTime = time.Now().Add(time.Hour * cfg.ExpireDuration).Unix()

		c.Success(resp)
	}
}

func (handle *authHandler) RegistryRouter(m *core.Mux) {
	authGroup := m.Group("/auth")
	{
		authGroup.POST("", handle.Auth())
		authGroup.GET("", handle.Auth())
	}
}

func NewAuthController(db db.Repo, cache cache.Repo) *authHandler {
	return &authHandler{
		cache,
		service.NewUserService(db),
	}
}
