package router

import (
	"deploy_server/controller"
	"deploy_server/pkg/cache"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
)

func RegistryRouters(r *core.Mux, db db.Repo, cache cache.Repo) {
	controller.NewUserController(db, cache).RegistryRouter(r)
	controller.NewAuthController(db, cache).RegistryRouter(r)
	controller.NewProjectController(db, cache).RegistryRouter(r)
}
