package controller

import (
	"deploy_server/middleware"
	"deploy_server/model/config"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/context"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
	errno "deploy_server/pkg/error"
	"deploy_server/service"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
)

var _ ConfigHandler = (*configHandler)(nil)

type configHandler struct {
	cache        cache.Repo
	configServer service.ConfigService
}

type ConfigHandler interface {
	RegistryRouter(m *core.Mux)
	List() context.HandlerFunc
	Get() context.HandlerFunc
	Create() context.HandlerFunc
	Update() context.HandlerFunc
	Test() context.HandlerFunc
}

type DataLoadFn func(args ...interface{})

func CacheGet(cache cache.Repo, key string, data interface{}) {
	val, _ := cache.Get(key)
	if val == "" {

	} else {
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			panic(err)
		}
	}
}

func (h *configHandler) Test() context.HandlerFunc {
	return func(c context.Context) {
		req := new(IdRequest)

		if err := c.ShouldBindQuery(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest, 10001, err.Error(),
			))
			return
		}

		key := fmt.Sprintf("cfg:%d", req.Id)

		val, _ := h.cache.Get(key, cache.WithTrace(c.Trace()))
		var cfg *config.Config
		if val == "" {
			cfg = h.configServer.Get(c.RequestContext(), req.Id)

			byte, _ := json.Marshal(cfg)
			h.cache.Set(key, string(byte), redis.KeepTTL, cache.WithTrace(c.Trace()))
		} else {
			fmt.Printf("%s\n", val)

			err := json.Unmarshal([]byte(val), &cfg)
			if err != nil {
				panic(err)
			}
		}
		c.Payload(&SimpleResponse{http.StatusOK, "", cfg})
	}
}

func (h *configHandler) Update() context.HandlerFunc {
	return func(c context.Context) {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 0)

		req := make(map[string]interface{})

		if err := c.ShouldBindFormToMap(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest, 10001, err.Error(),
			))
			return
		}

		h.configServer.Update(c.RequestContext(), int32(id), req)

		fmt.Printf("sdfsa%v\n", req)
		c.Payload(&SimpleResponse{
			http.StatusOK,
			"",
			req,
		})
	}
}

type storeConfigRequest struct {
	Title string `form:"title"`
	Key   string `form:"key"`
	Value string `form:"value"`
	Type  int32  `form:"type"`
}

func (h *configHandler) Create() context.HandlerFunc {
	return func(c context.Context) {
		req := new(storeConfigRequest)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest, 10001, err.Error(),
			))
			return
		}

		id, err := h.configServer.Create(c.RequestContext(), req)

		if err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest, 10001, err.Error(),
			))
			return
		}

		c.Payload(&SimpleResponse{
			http.StatusOK,
			"",
			id,
		})
	}
}

func (h *configHandler) Get() context.HandlerFunc {
	return func(c context.Context) {
		req := new(IdRequest)

		if err := c.ShouldBindQuery(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest, 10001, err.Error(),
			))
			return
		}

		cfg := h.configServer.Get(c.RequestContext(), req.Id)
		c.Payload(&SimpleResponse{
			http.StatusOK,
			"",
			cfg,
		})
	}
}

func (h *configHandler) List() context.HandlerFunc {
	return func(c context.Context) {
		configList, err := h.configServer.List(c.RequestContext())

		if err != nil {
			c.Logger().Error(err.Error())
		}

		c.Payload(&SimpleResponse{
			http.StatusOK,
			"",
			configList,
		})
	}
}

func (h *configHandler) RegistryRouter(r *core.Mux) {
	userGroup := r.Group("/config", middleware.PermissionMiddleWare())
	{
		userGroup.GET("/list", h.List())
		userGroup.GET("/get", h.Get())
		userGroup.POST("/create", h.Create())
		userGroup.POST("/update/:id", h.Update())
		userGroup.GET("/test", h.Test())
	}
}

func NewUserController(db db.Repo, cache cache.Repo) ConfigHandler {
	return &configHandler{
		cache,
		service.NewConfigService(db),
	}
}
