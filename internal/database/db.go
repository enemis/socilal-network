package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"social-network-otus/internal/config"
)

type DatabaseStack struct {
	master *sqlx.DB
	slave1 *sqlx.DB
}

func NewDatabaseStack(config *config.Config) *DatabaseStack {
	return &DatabaseStack{
		master: setupDbConnection(config.DBHost, config.DBPort, config.DBUsername, config.DBName, config.DBPassword, config.DBSSLMode),
		slave1: setupDbConnection(config.DBHostReplica1, config.DBPortReplica1, config.DBUsernameReplica1, config.DBNameReplica1, config.DBPasswordReplica1, config.DBSSLModeReplica1),
	}
}

func (stack *DatabaseStack) Slave() *sqlx.DB {
	return stack.slave1
}

func (stack *DatabaseStack) Master() *sqlx.DB {
	return stack.master
}

func (stack *DatabaseStack) GetReadConnection() *sqlx.DB {
	return stack.slave1
}

func (stack *DatabaseStack) GetWriteConnection() *sqlx.DB {
	return stack.master
}

func setupDbConnection(host string, port uint, user, dbname, password, sslmode string) *sqlx.DB {
	connectString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	logrus.Debug("db connect string: ")
	logrus.Debugln(connectString)
	db := sqlx.MustConnect("postgres", connectString)
	db.SetMaxOpenConns(80)

	return db
}
