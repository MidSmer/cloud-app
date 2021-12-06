package db

import (
	"errors"
	"github.com/MidSmer/cloud-app/common"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func InitArticle() error {
	txn, err := DB.Begin()
	if err != nil {
		return err
	}

	sql := `
	create table if not exists article (
	    id serial primary key,
	    key varchar not null,
	    content varchar not null
	);
	`
	if _, err := txn.Exec(sql); err != nil {
		txn.Rollback()
		Logger.Fatal(err)
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func CreateArticle(content string) (string, error) {
	key, err := common.NewShortId()
	if err != nil {
		return "", err
	}

	txn, err := DB.Begin()
	if err != nil {
		return "", err
	}

	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert("article").
		Cols("key", "content").
		Vals(
			goqu.Vals{key, content},
		)
	sql, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	_, err = txn.Exec(sql, args...)
	if err != nil {
		return "", err
	}

	if err = txn.Commit(); err != nil {
		return "", err
	}

	return key, nil
}

func UpdateArticle(key, content string) (string, error) {
	txn, err := DB.Begin()
	if err != nil {
		return "", err
	}

	dialect := goqu.Dialect("postgres")
	ds := dialect.Update("article").Set(goqu.Record{
		"content": content,
	}).Where(goqu.Ex{"key": key})
	sql, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	rows, err := txn.Query(sql, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if err = txn.Commit(); err != nil {
		return "", err
	}

	return key, nil
}

func GetArticleForKey(key string) (string, error) {
	txn, err := DB.Begin()
	if err != nil {
		return "", err
	}

	dialect := goqu.Dialect("postgres")
	ds := dialect.From("article").Select("content").Where(goqu.Ex{"key": key}).Limit(1)
	sql, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	rows, err := txn.Query(sql, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	contents := make([]string, 0)

	for rows.Next() {
		var content string

		err := rows.Scan(&content)
		if err != nil {
			return "", err
		}

		contents = append(contents, content)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	if err = txn.Commit(); err != nil {
		return "", err
	}

	if len(contents) != 1 {
		return "", errors.New("not record")
	}

	return contents[0], nil
}
