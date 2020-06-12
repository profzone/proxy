package storage

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"longhorn/proxy/internal/global"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type StorageMysql struct {
	sync.Mutex
	db *sql.DB
}

func NewDBMysql(config global.DBConfig) (*StorageMysql, error) {
	db := &StorageMysql{}
	err := db.init(config)
	return db, err
}

func (s *StorageMysql) init(config global.DBConfig) error {
	var conf = mysql.NewConfig()
	conf.User = config.UserName
	conf.Passwd = config.Password
	conf.Addr = config.Endpoints[0]
	conf.DBName = config.DatabaseName

	var err error
	s.db, err = sql.Open("mysql", conf.FormatDSN())
	return err
}

func (s *StorageMysql) Close() error {
	return s.db.Close()
}

func (s *StorageMysql) Create(prefix string, e Element) (uint64, error) {
	panic("implement me")
}

func (s *StorageMysql) Update(prefix string, condition *Condition, e Element) error {
	panic("implement me")
}

func (s *StorageMysql) Delete(prefix string, condition *Condition) error {
	panic("implement me")
}

func (s *StorageMysql) Get(prefix string, idField string, idVal uint64, target Element) error {
	panic("implement me")
}

func (s *StorageMysql) Walk(prefix string, condition *Condition, startField string, start uint64, limit int64, elementFactory func() Element, walking func(e Element) error) (uint64, error) {
	panic("implement me")
}
