package service

import (
	"context"
	"deploy_server/model"
	"deploy_server/model/config"
	"deploy_server/pkg/db"
)

var _ ConfigService = (*configService)(nil)

type ConfigService interface {
	Get(ctx context.Context, id int32) *config.Config
	List(ctx context.Context) ([]*config.Config, error)
	Create(ctx context.Context, request interface{}) (id int32, err error)
	Update(ctx context.Context, id int32, request map[string]interface{}) (err error)
}

type configService struct {
	Service
}

func NewConfigService(db db.Repo) *configService {
	return &configService{
		service{
			db,
		},
	}
}

func (s *configService) Update(ctx context.Context, id int32, data map[string]interface{}) (err error) {
	err = config.NewQueryBuilder().
		WhereId(model.EqualPredicate, id).
		Updates(s.GetDBWriter(ctx), data)
	return
}

func (s *configService) Create(ctx context.Context, request interface{}) (id int32, err error) {
	configModel := config.NewModel()
	configModel.Assign(request)
	id, err = configModel.Create(s.GetDBWriter(ctx))

	return
}

func (s *configService) Get(ctx context.Context, id int32) *config.Config {
	cfg, err := config.
		NewQueryBuilder().
		WhereId(model.EqualPredicate, id).
		QueryOne(s.GetDBReader(ctx))

	if err != nil {

	}

	return cfg
}

func (s *configService) List(ctx context.Context) ([]*config.Config, error) {
	configList, err := config.
		NewQueryBuilder().
		Limit(10).
		Offset(0).
		WhereIsDelete(model.EqualPredicate, 0).
		QueryAll(s.GetDBReader(ctx))

	return configList, err
}
