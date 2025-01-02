package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type Storage struct {
  DB *sql.DB
}

var (
  AliasExists = errors.New("alias is already exists")
)

func Init(dbPassword string) (storage *Storage, err error) {
  db, err := sql.Open("postgres", fmt.Sprintf("host=db port=5432 user=postgres password=%s sslmode=disable", dbPassword))
  if err != nil {
    return storage, err
  }

  storage = &Storage{
    DB: db,
  }

  err = db.Ping()
  if err != nil {
    return storage, err 
  }

  query := `CREATE TABLE IF NOT EXISTS urls(
    id SERIAL PRIMARY KEY,
    alias VARCHAR NOT NULL,
    link VARCHAR NOT NULL,
    views INTEGER DEFAULT 0
  )`
  _, err = db.Exec(query)
  if err != nil {
    return storage, err 
  }

  query = `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL
  )`

  _, err = db.Exec(query)
  if err != nil {
    return storage, err 
  }

  slog.Info("connected to database")
  return storage, nil
}

func (s Storage) Close() {
  s.DB.Close()
  slog.Debug("database was close")
}

func (s Storage) Check(dbname, column, value string) (bool, error) {
  query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", dbname, column)
  var count int 
  err := s.DB.QueryRow(query, value).Scan(&count)
  if err != nil {
    return false, err 
  }
  
  slog.Debug("checked value")
  return count > 0, nil
}

func (s Storage) Save(dbname, column1, column2, value1, value2 string) error {
  x, err := s.Check(dbname, column1, value1)
  if x {
    return AliasExists
  }

  query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)", dbname, column1, column2)
  _, err = s.DB.Exec(query, value1, value2)
  if err != nil {
    return err
  }

  if column1 == "alias" {
    slog.Info("saved link", "alias", value1, "link", value2)
  } else if column1 == "username" {
    slog.Info("add user", "username", value1, "password", value2)
  }
  return nil
}

func (s Storage) Get(column, dbname, value1, value2 string) (string, error) { 
  query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", column, dbname, value1)
  var result string
  err := s.DB.QueryRow(query, value2).Scan(&result)
  if err != nil {
    return "", err
  }

  if dbname == "urls" {
    query = `UPDATE urls SET views = views + 1 WHERE alias = $1`
    _, err = s.DB.Exec(query, value2)
    if err != nil {
      return "", err
    }
  }

  return result, nil
}
