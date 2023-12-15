package dto

import "gin-use-demo/model"

type UserDto struct {
	Name     string `json:"name"`
	Number   string `json:"number"`
	Password string `json:"password"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:     user.Name,
		Number:   user.Number,
		Password: user.Password,
	}
}
