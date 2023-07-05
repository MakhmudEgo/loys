package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Birthdate  time.Time `json:"birthdate"`
	Gender     string    `json:"gender"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
	Password   string    `json:"-"`
}

type UserCreateReq struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Birthdate  string `json:"birthdate"  validate:"required,birthdate"`
	Gender     string `json:"gender" validate:"required,gender"`
	Biography  string `json:"biography"`
	City       string `json:"city" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (u *UserCreateReq) ToUser() *User {
	birthdate, _ := time.Parse(time.DateOnly, u.Birthdate)
	return &User{
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		Birthdate:  birthdate,
		Gender:     u.Gender,
		Biography:  u.Biography,
		City:       u.City,
		Password:   u.Password,
	}
}

type UserCreateResp struct {
	ID string `json:"user_id"`
}

type UserAuthReq struct {
	ID       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type UserAuthResp struct {
	AccessToken string `json:"access_token" validate:"required"`
}
