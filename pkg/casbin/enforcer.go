package casbin

import (
	"deploy_server/config"
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func New() *Casbin {
	return &Casbin{}
}

type Casbin struct {
	enforcer *casbin.Enforcer
}

func (c *Casbin) Load() {
	cfg := config.Get().MySQL.Read
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Pass, cfg.Addr, cfg.Name)
	adapter := gormadapter.NewAdapter("mysql", dsn, true)
	basePath, _ := os.Getwd()
	c.enforcer = casbin.NewEnforcer(basePath+"/config/rbac2.conf", adapter)
	err := c.enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}

	// 日志记录
	c.enforcer.EnableLog(false)
}

func (c *Casbin) Enforce(vals ...interface{}) bool {
	return c.enforcer.Enforce(vals...)
}
