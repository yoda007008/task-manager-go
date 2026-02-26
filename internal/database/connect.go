package database

import (
	"fmt"
	//"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных %w", err)
	}

	db.SetMaxOpenConns(25) // максимальное кол-во открытых соединений
	db.SetMaxIdleConns(5)  // максимальное кол-во простаивающих соединений

	return db, nil
}
