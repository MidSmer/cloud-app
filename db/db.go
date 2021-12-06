package db

import (
	"database/sql"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var (
	Logger *logrus.Logger
	DB     *sql.DB
)

func Init(connStr string) error {
	var err error
	DB, err = GetDB(connStr)
	if err != nil {
		Logger.Fatal(err)
	}

	if err = InitArticle(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if err := DB.Close(); err != nil {
		return err
	}

	return nil
}

func GetDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
