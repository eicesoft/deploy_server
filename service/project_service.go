package service

import (
	"context"
	"deploy_server/model/project"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
)

type WrapBuilderFunc func(*project.QueryBuilder) *project.QueryBuilder

type ProjectService struct {
	service
}

// Get 利用warp函数进行Builder构造条件获得单条数据
func (s *ProjectService) Get(ctx context.Context, warp WrapBuilderFunc) (p *project.Project, err error) {
	p, err = warp(project.NewQueryBuilder()).
		QueryOne(s.GetDBReader(ctx))

	return
}

// GetAll 利用warp函数进行Builder构造条件获得列表数据
func (s *ProjectService) GetAll(ctx context.Context, warp WrapBuilderFunc) ([]*project.Project, error) {
	list, err := warp(project.NewQueryBuilder()).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// GetPages 利用warp函数进行Builder构造条件获得分页数据
func (s *ProjectService) GetPages(ctx context.Context, page int, pageSize int, warp WrapBuilderFunc) ([]*project.Project, error) {
	list, err := warp(project.NewQueryBuilder()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// Create 创建数据
func (s *ProjectService) Create(ctx context.Context, request interface{}) (id int32, err error) {
	projectModel := project.NewModel()
	core.StructCopy(projectModel, request)
	id, err = projectModel.Create(s.GetDBWriter(ctx))

	return
}

// Update 利用warp函数进行Builder构造更新条件进行更新
func (s *ProjectService) Update(ctx context.Context, data map[string]interface{}, warp WrapBuilderFunc) (err error) {
	err = warp(project.NewQueryBuilder()).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// UpdateById 按ID更新数据
func (s *ProjectService) UpdateById(ctx context.Context, data map[string]interface{}, id int32) (err error) {
	err = project.NewQueryBuilder().
		WhereId(project.EqualPredicate, id).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// Delete 按ID删除数据
func (s *ProjectService) Delete(ctx context.Context, id int32) (err error) {
	projectModel := project.NewModel()
	projectModel.Id = id
	err = projectModel.Delete(s.GetDBWriter(ctx))
	return
}

// NewProjectService ProjectService构造函数
func NewProjectService(db db.Repo) *ProjectService {
	return &ProjectService{
		service{
			db,
		},
	}
}
