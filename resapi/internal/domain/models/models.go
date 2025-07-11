package models

type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Login string `json:"log"`
	Password string `json:"pass"`
}