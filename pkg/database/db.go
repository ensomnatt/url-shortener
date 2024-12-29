package database

import (
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/lib/pq"
)

type Storage struct {
  DB *sql.DB
}

var (
  AliasExists = errors.New("alias is already exists")
)

func Init(connStr string) (storage *Storage, err error) {
  db, err := sql.Open("postgres", connStr)
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

  slog.Info("connected to database")
  return storage, nil
}

func (s Storage) Close() {
  s.DB.Close()
  slog.Debug("database was close")
}

func (s Storage) Check(alias string) (bool, error) {
  query := `SELECT COUNT(*) FROM urls WHERE alias = $1`
  var count int 
  err := s.DB.QueryRow(query, alias).Scan(&count)
  if err != nil {
    return false, err 
  }
  
  slog.Debug("checked alias")
  return count > 0, nil
}

func (s Storage) Save(alias, link string) error {
  x, err := s.Check(alias)
  if x {
    return AliasExists
  }

  query := `INSERT INTO urls (alias, link) VALUES ($1, $2)`
  _, err = s.DB.Exec(query, alias, link)
  if err != nil {
    return err
  }

  slog.Info("saved link", "alias", alias, "link", link)
  return nil
}

func (s Storage) Get(alias string) (string, error) { 
  query := `SELECT link FROM urls WHERE alias = $1`
  var link string
  err := s.DB.QueryRow(query, alias).Scan(&link)
  if err != nil {
    return "", err
  }

  query = `UPDATE urls SET views = views + 1 WHERE alias = $1`
  _, err = s.DB.Exec(query, alias)
  if err != nil {
    return "", err
  }

  return link, nil
}
