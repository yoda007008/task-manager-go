package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(databaseUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseUrl)

	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных %w", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(25) // максимальное кол-во открытых соединений
	db.SetMaxIdleConns(5)  // максимальное кол-во простаивающих соединений

	return db, nil
}
