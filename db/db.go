package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jordi-reinsma/bagel/model"
	_ "github.com/mattn/go-sqlite3"
)

const (
	connString   = "file:%s?_foreign_keys=on"
	sqliteDBPath = "./db/sql/bagel.sqlite"
	dbSchemaPath = "./db/sql/schema.sql"
)

type DB struct {
	*sql.DB
}

func MustConnect(reset bool) DB {
	if reset {
		fmt.Println("Resetting database")
		os.Remove(sqliteDBPath)
	}

	conn, err := sql.Open("sqlite3", fmt.Sprintf(connString, sqliteDBPath))
	if err != nil {
		panic(err)
	}

	db := DB{conn}

	err = db.createTables()
	if err != nil {
		panic(err)
	}

	return db
}

func (db DB) createTables() error {
	schema, err := os.ReadFile(dbSchemaPath)
	if err != nil {
		return err
	}

	_, err = db.Exec((string(schema)))
	return err
}

func (db DB) AddAndGetUsers(users []model.User) ([]model.User, error) {
	query := "INSERT INTO users (uuid) VALUES (?) ON CONFLICT (uuid) DO UPDATE SET uuid=EXCLUDED.uuid RETURNING id"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for i := range users {
		err = stmt.QueryRow(users[i].UUID).Scan(&users[i].ID)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (db DB) AddAndGetPairs(pairs []model.Match) ([]model.Match, error) {
	query := "INSERT INTO matches (a, b) VALUES (?, ?) ON CONFLICT (a, b) DO UPDATE SET a=EXCLUDED.a RETURNING id, freq"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for i := range pairs {
		err = stmt.QueryRow(pairs[i].A.ID, pairs[i].B.ID).Scan(&pairs[i].ID, &pairs[i].Freq)
		if err != nil {
			return nil, err
		}
	}
	return pairs, nil
}

func (db DB) UpdateMatch(match model.Match) error {
	query := "UPDATE matches SET freq = freq + 1 WHERE id = ?"
	_, err := db.Exec(query, match.ID)
	return err
}

func (db DB) GetLastExecutionDate() (time.Time, error) {
	var date time.Time
	query := "SELECT date FROM executions ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(query).Scan(&date)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	return date, err
}

func (db DB) SaveExecution() error {
	query := "INSERT INTO executions DEFAULT VALUES"
	_, err := db.Exec(query)
	return err
}
