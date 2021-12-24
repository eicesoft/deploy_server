package service

import (
	"context"
	"deploy_server/model/role"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
)

type WrapRoleBuilderFunc func(*role.QueryBuilder) *role.QueryBuilder

type RoleService struct {
	service
}

// Get 利用warp函数进行Builder构造条件获得单条数据
func (s *RoleService) Get(ctx context.Context, warp WrapRoleBuilderFunc) (p *role.Role, err error) {
	p, err = warp(role.NewQueryBuilder()).
		QueryOne(s.GetDBReader(ctx))

	return
}

// GetAll 利用warp函数进行Builder构造条件获得列表数据
func (s *RoleService) GetAll(ctx context.Context, warp WrapRoleBuilderFunc) ([]*role.Role, error) {
	list, err := warp(role.NewQueryBuilder()).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// GetPages 利用warp函数进行Builder构造条件获得分页数据
func (s *RoleService) GetPages(ctx context.Context, page int, pageSize int, warp WrapRoleBuilderFunc) ([]*role.Role, error) {
	list, err := warp(role.NewQueryBuilder()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// Create 创建数据
func (s *RoleService) Create(ctx context.Context, request interface{}) (id int32, err error) {
	roleModel := role.NewModel()
	core.StructCopy(roleModel, request)
	id, err = roleModel.Create(s.GetDBWriter(ctx))

	return
}

// Update 利用warp函数进行Builder构造更新条件进行更新
func (s *RoleService) Update(ctx context.Context, data map[string]interface{}, warp WrapRoleBuilderFunc) (err error) {
	err = warp(role.NewQueryBuilder()).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// UpdateById 按ID更新数据
func (s *RoleService) UpdateById(ctx context.Context, data map[string]interface{}, id int32) (err error) {
	err = role.NewQueryBuilder().
		WhereId(role.EqualPredicate, id).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// Delete 按ID删除数据
func (s *RoleService) Delete(ctx context.Context, id int32) (err error) {
	roleModel := role.NewModel()
	roleModel.Id = id
	err = roleModel.Delete(s.GetDBWriter(ctx))
	return
}

// NewRoleService RoleService构造函数
func NewRoleService(db db.Repo) *RoleService {
	return &RoleService{
		service{
			db,
		},
	}
}
