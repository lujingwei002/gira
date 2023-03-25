package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
)

type AdminDbClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *sql.DB
	config     gira.AdminDbConfig
}

func NewAdminDbClient() *AdminDbClient {
	return &AdminDbClient{}
}

func (self *AdminDbClient) GetMysqlClient() *sql.DB {
	return self.client
}

// 初始化并启动
func (self *AdminDbClient) Start(ctx context.Context, config gira.AdminDbConfig) error {
	self.config = config
	self.ctx, self.cancelFunc = context.WithCancel(ctx)

	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Db)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	self.client = db
	log.Info("connect admindb success")
	return nil
}
