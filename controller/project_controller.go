package controller

import (
	"deploy_server/middleware"
	"deploy_server/model/project"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/context"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
	"deploy_server/service"
	"fmt"
	"strconv"
)

type projectHandler struct {
	cache         cache.Repo
	projectServer *service.ProjectService
}

func (handle *projectHandler) List() context.HandlerFunc {
	return func(ctx context.Context) {
		page, err := strconv.ParseInt(ctx.Params("page"), 10, 32)
		if err != nil {
			page = DefaultPage
		}

		pageSize, err := strconv.ParseInt(ctx.Params("page_size"), 10, 32)
		if err != nil {
			pageSize = DefaultPageSize
		}
		projects, err := handle.projectServer.GetPages(ctx.RequestContext(), int(page), int(pageSize), func(build *project.QueryBuilder) *project.QueryBuilder {
			return build.WhereGitType(project.EqualPredicate, 1)
		})

		if err != nil {
			ctx.Fail(10002, err)
			return
		}

		ctx.Success(projects)
	}
}

// storeProjectRequest 存储项目请求
type storeProjectRequest struct {
	Title     string `form:"title"`
	Desc      string `form:"desc"`
	GitUrl    string `form:"git_url"`
	GitType   int32  `form:"git_type"`
	GitBranch string `form:"git_branch"`
}

// Create 创建项目数据
// @Summary 创建项目数据
// @Description 创建项目数据
// @Tags Project
// @Produce  json
// @Param title query string true "项目名称"
// @Param desc query string true "项目描述"
// @Param git_url query string true "Git地址"
// @Param git_type query string true "Git类型(1 => 分支, 2 => Tag)"
// @Param git_branch query string true "Git分支或者Tag名称"
// @Param Authorization header string true "验证Token"
// @Success 200 {object} SimpleResponse
// @Failure 400 {object} core.Failure
// @Router /project/create [post]
func (handle *projectHandler) Create() context.HandlerFunc {
	return func(ctx context.Context) {
		req := new(storeProjectRequest)

		if err := ctx.ShouldBindForm(&req); err != nil {
			ctx.Fail(10002, err)
		}
		id, err := handle.projectServer.Create(ctx.RequestContext(), req)

		if err != nil {
			ctx.Fail(10002, err)
		}

		ctx.Success(id)
	}
}

func (handle *projectHandler) Get() context.HandlerFunc {
	return func(c context.Context) {
		var req = new(IdRequest)
		if err := c.ShouldBindQuery(req); err != nil {
			c.Fail(10002, err)
			return
		}
		p, err := handle.projectServer.Get(c.RequestContext(), func(builder *project.QueryBuilder) *project.QueryBuilder {
			builder.WhereId(project.EqualPredicate, req.Id)
			return builder
		})
		if err != nil {
			c.Fail(10003, err)
		} else {
			c.Success(p)
		}
	}
}

// Update 更新项目数据
// @Summary 更新项目数据
// @Description 更新项目数据
// @Tags Project
// @Produce  json
// @Param id path string true "项目ID"
// @Param title query string true "项目名称"
// @Param desc query string true "项目描述"
// @Param git_url query string false "Git地址"
// @Param git_type query string false "Git类型(1 => 分支, 2 => Tag)"
// @Param git_branch query string false "Git分支或者Tag名称"
// @Param Authorization header string true "验证Token"
// @Success 200 {object} SimpleResponse
// @Failure 400 {object} core.Failure
// @Router /project/update/{id} [post]
func (handle *projectHandler) Update() context.HandlerFunc {
	return func(c context.Context) {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 0)

		req := make(map[string]interface{})
		if err := c.ShouldBindFormToMap(req); err != nil {
			c.Fail(10001, err)
			return
		}

		err := handle.projectServer.UpdateById(c.RequestContext(), req, int32(id))
		if err != nil {
			c.Fail(10002, err)
			return
		}

		c.Success("success")
	}
}

func (handle *projectHandler) Delete() context.HandlerFunc {
	return func(c context.Context) {
		var req = new(IdRequest)
		if err := c.ShouldBindQuery(req); err != nil {
			c.Fail(10002, err)
			return
		}
		c.Logger().Debug(fmt.Sprintf("%v", req))
		err := handle.projectServer.Delete(c.RequestContext(), req.Id)
		if err != nil {
			c.Fail(10003, err)
		} else {
			c.Success("success")
		}
	}
}

func (handle *projectHandler) RegistryRouter(m *core.Mux) {
	authGroup := m.Group("/project", middleware.PermissionMiddleWare())
	{
		authGroup.GET("/list", handle.List())
		authGroup.POST("/create", handle.Create())
		authGroup.POST("/update/:id", handle.Update())
		authGroup.GET("/delete", handle.Delete())
		authGroup.GET("/get", handle.Get())
	}
}

func NewProjectController(db db.Repo, cache cache.Repo) *projectHandler {
	return &projectHandler{
		cache,
		service.NewProjectService(db),
	}
}
