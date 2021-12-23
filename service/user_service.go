package service

import (
	"context"
	"crypto/sha256"
	"deploy_server/model"
	"deploy_server/model/user"
	"deploy_server/pkg/db"
	"fmt"
)

type UserService struct {
	service
}

func EncodePass(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (u UserService) Check(ctx context.Context, username string, password string) *user.User {
	userModel, err := user.
		NewQueryBuilder().
		WhereName(model.EqualPredicate, username).
		WherePassword(model.EqualPredicate, EncodePass(password)).
		QueryOne(u.GetDBReader(ctx))

	if err != nil {
		return nil
	}

	return userModel
}

func NewUserService(db db.Repo) *UserService {
	return &UserService{
		service{
			db,
		},
	}
}
