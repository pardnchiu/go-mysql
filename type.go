package mysqlPool

import (
	"database/sql"

	goLogger "github.com/pardnchiu/go-logger"
)

const (
	defaultLogPath      = "./logs/cron.log"
	defaultLogMaxSize   = 16 * 1024 * 1024
	defaultLogMaxBackup = 5
)

type Log = goLogger.Log
type Logger = goLogger.Logger

type Pool struct {
	db     *sql.DB
	Logger *Logger
}

type PoolList struct {
	Read   *Pool
	Write  *Pool
	Logger *Logger
}

type DBConfig struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	Charset    string `json:"charset,omitempty"`
	Connection int    `json:"connection,omitempty"`
}

type Config struct {
	Read  *DBConfig
	Write *DBConfig
	Log   *Log `json:"log,omitempty"`
}

type builder struct {
	db          *sql.DB
	dbName      *string
	table       *string
	selectList  []string
	joinList    []string
	whereList   []string
	bindingList []interface{}
	orderList   []string
	setList     []string
	limit       *int
	offset      *int
	withTotal   bool
	logger      *Logger
}
