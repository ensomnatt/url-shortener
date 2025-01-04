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
  UsernameExists = errors.New("username is already exists")
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
    views INTEGER DEFAULT 0,
    owner VARCHAR NOT NULL
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

func (s Storage) Check(alias string) (bool, error) {
  query := `"SELECT COUNT(*) FROM urls WHERE alias = alias`
  var count int 
  err := s.DB.QueryRow(query, alias).Scan(&count)
  if err != nil {
    return false, err 
  }
  
  slog.Debug("checked alias")
  return count > 0, nil
}

func (s Storage) Save(alias, link, username string) error {
  x, err := s.Check(alias)
  if x {
    return AliasExists
  }

  query := `INSERT INTO urls (alias, link, owner) VALUES ($1, $2, $3)`
  _, err = s.DB.Exec(query, alias, link, username)
  if err != nil {
    return err
  }

  slog.Info("added link", "alias", alias, "link", link, "username", username)
  return nil
}

func (s Storage) Get(alias string) (string, string, error) { 
  query := `SELECT link, owner FROM urls WHERE alias = $1`
  var link, owner string
  err := s.DB.QueryRow(query, alias).Scan(&link, &owner)
  if err != nil {
    return "", "", err
  }
 
  query = `UPDATE urls SET views = views + 1 WHERE alias = $1`
  _, err = s.DB.Exec(query, alias)
  if err != nil {
    return "", "", err
  }

  return link, owner, nil
}

func (s Storage) Delete(alias string) error {
  query := `DELETE FROM urls WHERE alias = $1`
  _, err := s.DB.Exec(query, alias)
  if err != nil {
    return err
  }

  slog.Info("deleted one link", "alias", alias)
  return nil
}

func (s Storage) CheckUser(username string) (bool, error) {
  query := `SELECT COUNT(*) FROM users WHERE username = $1`
  var count int 
  err := s.DB.QueryRow(query, username).Scan(&count)
  if err != nil {
    return false, err 
  }
  
  slog.Debug("checked username")
  return count > 0, nil
}

func (s Storage) SaveUser(username, password string) error {
  x, err := s.CheckUser(username)
  if x {
    return UsernameExists
  }

  query := `INSERT INTO urls (username, password) VALUES ($1, $2)`
  _, err = s.DB.Exec(query, username, password)
  if err != nil {
    return err
  }

  slog.Info("added user", "username", username)
  return nil
}

func (s Storage) GetUser(username string) (string, error){ 
  query := `SELECT password FROM users WHERE username = $1`
  var result string
  err := s.DB.QueryRow(query, username).Scan(&result)
  if err != nil {
    return "", err
  }

  return result, nil
}
