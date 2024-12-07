package svc

import (
	"go-zero-learn/models"
	"go-zero-learn/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModal models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		UserModal: models.NewUsersModel(SqlConn, c.Cache),
	}
}
