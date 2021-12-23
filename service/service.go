package service

import (
	"context"
	"deploy_server/pkg/db"

	"gorm.io/gorm"
)

var _ Service = (*service)(nil)

type service struct {
	Db db.Repo
}

type Service interface {
	GetDBReader(ctx context.Context) *gorm.DB
	GetDBWriter(ctx context.Context) *gorm.DB
}

func (h service) GetDBWriter(ctx context.Context) *gorm.DB {
	return h.Db.GetDbW().WithContext(ctx)
}

func (h service) GetDBReader(ctx context.Context) *gorm.DB {
	return h.Db.GetDbR().WithContext(ctx)
}
