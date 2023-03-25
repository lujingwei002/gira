package gira

import (
	"database/sql"
)

type MysqlClient interface {
	GetMysqlClient() *sql.DB
}
