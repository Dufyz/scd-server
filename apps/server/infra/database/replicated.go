package db

import "database/sql"

type ReplicatedDB struct {
	Write *sql.DB
	Read  *sql.DB
}

func NewReplicatedDB(write, read *sql.DB) *ReplicatedDB {
	return &ReplicatedDB{Write: write, Read: read}
}

func (d *ReplicatedDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.Read.Query(query, args...)
}

func (d *ReplicatedDB) QueryRow(query string, args ...any) *sql.Row {
	return d.Read.QueryRow(query, args...)
}

func (d *ReplicatedDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.Write.Exec(query, args...)
}

func (d *ReplicatedDB) Begin() (*sql.Tx, error) {
	return d.Write.Begin()
}
