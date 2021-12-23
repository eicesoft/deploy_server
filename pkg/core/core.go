package core

import (
	stdctx "context"
	_ "deploy_server/docs"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/color"
	"deploy_server/pkg/context"
	"deploy_server/pkg/db"
	"deploy_server/pkg/env"
	"deploy_server/pkg/trace"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

const (
	ServerHeader = "Server"
	ServerName   = "GServer 1.0"
)

type Option struct {
	Env     string
	Port    string
	IsDebug bool
}

// StructCopy 从value复制结构数据到 binding 中
func StructCopy(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem()
	vVal := reflect.ValueOf(value).Elem()
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok { //目标结构中存在字段
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface())) //
		}
	}
}

func wrapHandlers(handlers ...context.HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {
			ctx := context.NewContext(c)
			defer context.ReleaseContext(ctx)

			handler(ctx)
		}
	}

	return funcs
}

type Failure struct {
	Code         int    `json:"code"`          // 业务码
	BusinessCode int    `json:"business_code"` // 业务码
	Message      string `json:"message"`       // 描述信息
}

func DisableTrace(ctx context.Context) {
	ctx.DisableTrace()
}

type RouterGroup interface {
	Group(string, ...context.HandlerFunc) RouterGroup
	IRoutes
}

var _ IRoutes = (*router)(nil)

// IRoutes 包装gin的IRoutes
type IRoutes interface {
	Any(string, ...context.HandlerFunc)
	GET(string, ...context.HandlerFunc)
	POST(string, ...context.HandlerFunc)
	DELETE(string, ...context.HandlerFunc)
	PATCH(string, ...context.HandlerFunc)
	PUT(string, ...context.HandlerFunc)
	OPTIONS(string, ...context.HandlerFunc)
	HEAD(string, ...context.HandlerFunc)
}

type router struct {
	group *gin.RouterGroup
}

func (r *router) Group(relativePath string, handlers ...context.HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(handlers...)...)
	return &router{group: group}
}

func (r *router) Any(relativePath string, handlers ...context.HandlerFunc) {
	r.group.Any(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) GET(relativePath string, handlers ...context.HandlerFunc) {
	r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) POST(relativePath string, handlers ...context.HandlerFunc) {
	r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) DELETE(relativePath string, handlers ...context.HandlerFunc) {
	r.group.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PATCH(relativePath string, handlers ...context.HandlerFunc) {
	r.group.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PUT(relativePath string, handlers ...context.HandlerFunc) {
	r.group.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...context.HandlerFunc) {
	r.group.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...context.HandlerFunc) {
	r.group.HEAD(relativePath, wrapHandlers(handlers...)...)
}

var _ MuxServer = (*Mux)(nil)

type MuxServer interface {
	http.Handler
	Group(relativePath string, handlers ...context.HandlerFunc) RouterGroup
	init()
	Run()
	Shutdown(ctx stdctx.Context) error
}

type Mux struct {
	server *http.Server
	engine *gin.Engine
	logger *zap.Logger
	db     db.Repo
	cache  cache.Repo
	opt    *Option
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.engine.ServeHTTP(w, req)
}

func (m *Mux) Group(relativePath string, handlers ...context.HandlerFunc) RouterGroup {
	return &router{
		group: m.engine.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

func (m *Mux) init() {
	logger := m.logger

	withoutTracePaths := map[string]bool{
		"/swagger":                  true,
		"/metrics":                  true,
		"/debug/pprof/":             true,
		"/debug/pprof/cmdline":      true,
		"/debug/pprof/profile":      true,
		"/debug/pprof/symbol":       true,
		"/debug/pprof/trace":        true,
		"/debug/pprof/allocs":       true,
		"/debug/pprof/block":        true,
		"/debug/pprof/goroutine":    true,
		"/debug/pprof/heap":         true,
		"/debug/pprof/mutex":        true,
		"/debug/pprof/threadcreate": true,
		"/favicon.ico":              true,
		"/system/health":            true,
	}

	if env.Env != "prod" {
		m.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
		DebugLog("* register swagger router")
	}

	m.engine.Use(func(ctx *gin.Context) {
		ts := time.Now()
		//核心处理
		newCtx := context.NewContext(ctx)
		defer context.ReleaseContext(newCtx)

		newCtx.Init(m.logger)

		if !withoutTracePaths[ctx.Request.URL.Path] {
			//trace id 前端Header传递该值, 方便调试
			if traceId := newCtx.GetHeader(trace.Header); traceId != "" {
				newCtx.SetTrace(trace.New(traceId))
			} else {
				newCtx.SetTrace(trace.New(""))
			}
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Error("Http request Error:", zap.Any("err", err))

				//context.AbortWithError(errno.NewError(
				//	http.StatusInternalServerError,
				//	message.ServerError,
				//	message.Text(message.ServerError)),
				//)
			}

			if ctx.Writer.Status() == http.StatusNotFound {
				return
			}
			var (
				response        interface{}
				businessCode    int
				businessCodeMsg string
				abortErr        error
				//traceId         string
			)

			newCtx.SetHeader(ServerHeader, ServerName)

			if ctx.IsAborted() { //
				if err := newCtx.AbortError(); err != nil {
					response = err
					businessCode = err.GetBusinessCode()
					businessCodeMsg = err.GetMsg()

					if x := newCtx.Trace(); x != nil {
						newCtx.SetHeader(trace.Header, x.ID())
						//traceId = x.ID()
					}

					ctx.JSON(err.GetHttpCode(), &Failure{
						Code:         ctx.Writer.Status(),
						BusinessCode: businessCode,
						Message:      businessCodeMsg,
					})
				}
			} else {
				response = newCtx.GetPayload()
				if response != nil {
					if x := newCtx.Trace(); x != nil {
						newCtx.SetHeader(trace.Header, x.ID()) //设置Trace Id
						//traceId = x.ID()
					}
					ctx.JSON(http.StatusOK, response)
				}
			}

			var t *trace.Trace
			if x := newCtx.Trace(); x != nil {
				t = x.(*trace.Trace)
			} else {
				return
			}
			decodedURL, _ := url.QueryUnescape(ctx.Request.URL.RequestURI())
			t.WithRequest(&trace.Request{
				TTL:        "un-limit",
				Method:     ctx.Request.Method,
				DecodedURL: decodedURL,
				//Header:     ctx.Request.Header,
				Body: string(newCtx.RawData()),
			})

			var responseBody interface{}

			if response != nil {
				responseBody = response
			}

			t.WithResponse(&trace.Response{
				Header:          ctx.Writer.Header(),
				HttpCode:        ctx.Writer.Status(),
				HttpCodeMsg:     http.StatusText(ctx.Writer.Status()),
				BusinessCode:    businessCode,
				BusinessCodeMsg: businessCodeMsg,
				Body:            responseBody,
				CostSeconds:     time.Since(ts).Seconds(),
			})

			t.Success = !ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK
			t.CostSeconds = time.Since(ts).Seconds()

			logger.Info("router-interceptor",
				zap.Any("method", ctx.Request.Method),
				zap.Any("path", decodedURL),
				zap.Any("http_code", ctx.Writer.Status()),
				//zap.Any("business_code", businessCode),
				zap.Any("success", t.Success),
				zap.Any("cost_seconds", t.CostSeconds),
				zap.Any("trace_id", t.Identifier),
				zap.Any("trace_info", t),
				zap.Error(abortErr),
			)
		}()
		ctx.Next()
	})

	m.engine.NoMethod(wrapHandlers(DisableTrace)...)
	m.engine.NoRoute(wrapHandlers(DisableTrace)...)

	system := m.Group("/system")
	{
		system.GET("/health", func(ctx context.Context) {
			resp := &struct {
				Timestamp   time.Time `json:"timestamp"`
				Environment string    `json:"environment"`
				Host        string    `json:"host"`
				Status      string    `json:"status"`
			}{
				Timestamp:   time.Now(),
				Environment: m.opt.Env,
				Host:        ctx.Host(),
				Status:      "ok",
			}
			ctx.Payload(resp)
		})
	}
}

func (m *Mux) Run() {
	m.server = &http.Server{
		Addr:           ":" + m.opt.Port,
		Handler:        m.engine,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 28, //256M
	}

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Fatal("http server startup err", zap.Error(err))
		}
	}()
}

func (m *Mux) Shutdown(ctx stdctx.Context) error {
	return m.server.Shutdown(ctx)
}

func NewServer(logger *zap.Logger, db db.Repo, cache cache.Repo, opt *Option) *Mux {
	gin.DisableConsoleColor()
	gin.DefaultWriter = ioutil.Discard

	mux := &Mux{
		engine: gin.Default(),
		logger: logger,
		db:     db,
		cache:  cache,
		opt:    opt,
	}
	mux.init()

	return mux
}

func DebugLog(msg string, a ...interface{}) {
	fmt.Println(color.Green(fmt.Sprintf(msg, a...)))
}

func ErrorLog(msg string, a ...interface{}) {
	fmt.Println(color.Red(fmt.Sprintf(msg, a...)))
}
