package main

import (
	"database/sql"
	"testTask/data"
)

const (
	DriverName = "sqlite3"
	TableName  = "Queries"
)

type SqliteSaver struct {
	db *sql.DB
}

func NewSqliteSaver(path string) (*SqliteSaver, error) {
	db, err := initDB(path)
	if err != nil {
		return nil, err
	}
	return &SqliteSaver{db}, nil
}

func initDB(path string) (*sql.DB, error) {
	db, err := sql.Open(DriverName, path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (ds *SqliteSaver) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS ` + TableName + `
	(
		OperationID INTEGER,
		Action TEXT,
		QueryPermited INTEGER
	);`

	_, err := ds.db.Exec(query)
	return err
}

func (ds *SqliteSaver) Close() error {
	return ds.db.Close()
}

func (ds *SqliteSaver) SaveData(data data.Query) error {
	query := `INSERT INTO ` + TableName + `(
	OperationID,
	Action,
	QueryPermited
	) values (?, ?, ?)`

	prQuery, err := ds.db.Prepare(query)
	if err != nil {
		return err
	}
	defer prQuery.Close()

	_, err = prQuery.Exec(data.OperationID, data.Action, data.QueryPermited)
	return err
}
