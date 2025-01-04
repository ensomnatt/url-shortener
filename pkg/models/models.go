package models

type User struct {
  Username string 
  Password string
}

type RequestSave struct {
  Alias string `json:"alias"`
  Link string `json:"link"`
}

type ResponseSave struct {
  Link string `json:"link"`
}
